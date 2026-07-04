param(
    [Parameter(Mandatory = $true)]
    [ValidateSet('dev', 'prod', 'status', 'stop', 'open')]
    [string]$Action,
    [string]$ContainerName = 'sub2api',
    [string]$EnvFile = (Join-Path $PSScriptRoot '.env'),
    [string]$BindHost = '127.0.0.1',
    [int]$DevPort = 8080,
    [int]$BackendPort = 18080,
    [int]$HealthTimeoutSec = 120,
    [switch]$SkipSmoke
)

Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'

$repoRoot = Split-Path -Parent $PSScriptRoot
$tmpDir = Join-Path $repoRoot '.tmp'
$vitePidFile = Join-Path $tmpDir 'local-entry-vite.pid'
$viteOutLog = Join-Path $repoRoot 'frontend-dev-fixed-8080.out.log'
$viteErrLog = Join-Path $repoRoot 'frontend-dev-fixed-8080.err.log'
$browserUrl = "http://${BindHost}:${DevPort}"
$backendUrl = "http://${BindHost}:${BackendPort}"
$composeFile = Join-Path $PSScriptRoot 'docker-compose.local.yml'
$localOpsScript = Join-Path $PSScriptRoot 'local-ops.ps1'

New-Item -ItemType Directory -Force -Path $tmpDir | Out-Null

function Get-RunningProcessById {
    param([int]$Id)

    try {
        return Get-Process -Id $Id -ErrorAction Stop
    } catch {
        return $null
    }
}

function Get-ChildProcessIds {
    param([int]$ParentId)

    $children = Get-CimInstance Win32_Process -Filter "ParentProcessId = $ParentId" -ErrorAction SilentlyContinue
    $all = @()
    foreach ($child in $children) {
        $all += [int]$child.ProcessId
        $all += Get-ChildProcessIds -ParentId ([int]$child.ProcessId)
    }
    return $all
}

function Stop-ProcessTree {
    param([int]$RootId)

    $ids = @($RootId) + @(Get-ChildProcessIds -ParentId $RootId)
    $ids = $ids | Select-Object -Unique | Sort-Object -Descending
    foreach ($id in $ids) {
        try {
            Stop-Process -Id $id -Force -ErrorAction Stop
        } catch {
            # already exited
        }
    }
}

function Get-VitePid {
    if (-not (Test-Path -LiteralPath $vitePidFile)) {
        return $null
    }

    $content = Get-Content -LiteralPath $vitePidFile -ErrorAction SilentlyContinue | Select-Object -First 1
    if ([string]::IsNullOrWhiteSpace($content)) {
        return $null
    }

    $viteProcessId = 0
    if (-not [int]::TryParse($content.Trim(), [ref]$viteProcessId)) {
        return $null
    }

    return $viteProcessId
}

function Get-ViteProcess {
    $viteProcessId = Get-VitePid
    if ($null -eq $viteProcessId) {
        return $null
    }

    return Get-RunningProcessById -Id $viteProcessId
}

function Stop-Vite {
    $process = Get-ViteProcess
    if ($null -ne $process) {
        Stop-ProcessTree -RootId $process.Id
    }

    if (Test-Path -LiteralPath $vitePidFile) {
        Remove-Item -LiteralPath $vitePidFile -Force
    }
}

function Wait-HttpReady {
    param(
        [Parameter(Mandatory = $true)]
        [string]$Url,
        [int]$TimeoutSec = 60
    )

    $deadline = (Get-Date).AddSeconds($TimeoutSec)
    while ((Get-Date) -lt $deadline) {
        try {
            $response = Invoke-WebRequest -UseBasicParsing -Uri $Url -TimeoutSec 5
            if ($response.StatusCode -ge 200 -and $response.StatusCode -lt 500) {
                return
            }
        } catch {
            Start-Sleep -Milliseconds 800
            continue
        }
        Start-Sleep -Milliseconds 800
    }

    throw "Timed out waiting for $Url"
}

function Invoke-ComposeUp {
    param([int]$PublishedPort)

    $previousBindHost = $env:BIND_HOST
    $previousServerPort = $env:SERVER_PORT
    try {
        $env:BIND_HOST = $BindHost
        $env:SERVER_PORT = [string]$PublishedPort
        & docker compose --env-file $EnvFile -f $composeFile up -d
        if ($LASTEXITCODE -ne 0) {
            throw "docker compose up failed for published port $PublishedPort."
        }
    } finally {
        if ($null -eq $previousBindHost) {
            Remove-Item Env:BIND_HOST -ErrorAction SilentlyContinue
        } else {
            $env:BIND_HOST = $previousBindHost
        }

        if ($null -eq $previousServerPort) {
            Remove-Item Env:SERVER_PORT -ErrorAction SilentlyContinue
        } else {
            $env:SERVER_PORT = $previousServerPort
        }
    }
}

function Get-ContainerPublishedPort {
    param([string]$Name)

    $json = docker inspect $Name --format '{{json .HostConfig.PortBindings}}'
    if ($LASTEXITCODE -ne 0 -or [string]::IsNullOrWhiteSpace($json)) {
        return $null
    }

    $bindings = $json | ConvertFrom-Json
    $entry = $bindings.'8080/tcp'
    if (-not $entry -or $entry.Count -eq 0) {
        return $null
    }

    return [int]$entry[0].HostPort
}

function Get-PortListener {
    param([int]$Port)

    return Get-NetTCPConnection -State Listen -ErrorAction SilentlyContinue |
        Where-Object { $_.LocalPort -eq $Port } |
        Select-Object -First 1
}

function Start-Vite {
    Stop-Vite

    if (Test-Path -LiteralPath $viteOutLog) {
        Remove-Item -LiteralPath $viteOutLog -Force
    }
    if (Test-Path -LiteralPath $viteErrLog) {
        Remove-Item -LiteralPath $viteErrLog -Force
    }

    $command = @"
`$env:VITE_DEV_PROXY_TARGET = '${backendUrl}'
`$env:VITE_DEV_PORT = '${DevPort}'
Set-Location '${repoRoot}'
pnpm --dir frontend exec vite --host ${BindHost} --port ${DevPort} --strictPort
"@

    $process = Start-Process -FilePath 'pwsh' `
        -ArgumentList @('-NoProfile', '-Command', $command) `
        -WorkingDirectory $repoRoot `
        -RedirectStandardOutput $viteOutLog `
        -RedirectStandardError $viteErrLog `
        -WindowStyle Hidden `
        -PassThru

    Set-Content -LiteralPath $vitePidFile -Value $process.Id
}

function Get-StatusObject {
    $viteProcess = Get-ViteProcess
    $containerPort = Get-ContainerPublishedPort -Name $ContainerName
    $listener8080 = Get-PortListener -Port $DevPort
    $listener18080 = Get-PortListener -Port $BackendPort

    $mode = 'unknown'
    if ($null -ne $viteProcess -and $containerPort -eq $BackendPort) {
        $mode = 'dev'
    } elseif ($null -eq $viteProcess -and $containerPort -eq $DevPort) {
        $mode = 'prod'
    } elseif ($null -eq $viteProcess -and $containerPort -eq $BackendPort) {
        $mode = 'backend-only'
    }

    return [PSCustomObject]@{
        Mode = $mode
        BrowserUrl = $browserUrl
        BackendUrl = $backendUrl
        ViteRunning = ($null -ne $viteProcess)
        VitePid = if ($null -ne $viteProcess) { $viteProcess.Id } else { $null }
        ViteOutLog = $viteOutLog
        ViteErrLog = $viteErrLog
        ContainerPublishedPort = $containerPort
        Port8080Listener = if ($null -ne $listener8080) { "$($listener8080.LocalAddress):$($listener8080.LocalPort)" } else { $null }
        Port18080Listener = if ($null -ne $listener18080) { "$($listener18080.LocalAddress):$($listener18080.LocalPort)" } else { $null }
    }
}

switch ($Action) {
    'dev' {
        Invoke-ComposeUp -PublishedPort $BackendPort
        Wait-HttpReady -Url "${backendUrl}/health" -TimeoutSec $HealthTimeoutSec
        Start-Vite
        Wait-HttpReady -Url "${browserUrl}/home" -TimeoutSec $HealthTimeoutSec
        Get-StatusObject | ConvertTo-Json -Depth 4
        break
    }
    'prod' {
        Stop-Vite
        Invoke-ComposeUp -PublishedPort $DevPort
        Wait-HttpReady -Url "${browserUrl}/home" -TimeoutSec $HealthTimeoutSec
        if (-not $SkipSmoke) {
            & pwsh -File $localOpsScript smoke -BaseUrl $browserUrl
            if ($LASTEXITCODE -ne 0) {
                throw 'Smoke check failed after switching to prod mode.'
            }
        }
        Get-StatusObject | ConvertTo-Json -Depth 4
        break
    }
    'status' {
        Get-StatusObject | ConvertTo-Json -Depth 4
        break
    }
    'stop' {
        Stop-Vite
        Get-StatusObject | ConvertTo-Json -Depth 4
        break
    }
    'open' {
        Start-Process $browserUrl
        Write-Output $browserUrl
        break
    }
}
