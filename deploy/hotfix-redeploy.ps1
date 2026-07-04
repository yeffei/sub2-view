param(
    [string]$ContainerName = 'sub2api',
    [string]$BuilderImage = 'golang:1.26.4-alpine',
    [string]$BaseUrl = 'http://127.0.0.1:8080',
    [string]$EnvFile = (Join-Path $PSScriptRoot '.env'),
    [int]$HealthTimeoutSec = 120,
    [switch]$SkipTypecheck,
    [switch]$SkipFrontendBuild,
    [switch]$SkipSmoke
)

Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'

function Read-EnvFile {
    param(
        [Parameter(Mandatory = $true)]
        [string]$Path
    )

    $map = @{}
    if (-not (Test-Path -LiteralPath $Path)) {
        return $map
    }

    foreach ($line in Get-Content -LiteralPath $Path) {
        $trimmed = $line.Trim()
        if ([string]::IsNullOrWhiteSpace($trimmed) -or $trimmed.StartsWith('#')) {
            continue
        }

        $parts = $trimmed.Split('=', 2)
        if ($parts.Count -ne 2) {
            continue
        }

        $map[$parts[0].Trim()] = $parts[1].Trim()
    }

    return $map
}

function Invoke-Step {
    param(
        [Parameter(Mandatory = $true)]
        [string]$Name,
        [Parameter(Mandatory = $true)]
        [scriptblock]$Action
    )

    Write-Host "==> $Name"
    & $Action
}

function Wait-ForHealthy {
    param(
        [Parameter(Mandatory = $true)]
        [string]$ContainerName,
        [int]$TimeoutSec = 120
    )

    $deadline = (Get-Date).AddSeconds($TimeoutSec)
    while ((Get-Date) -lt $deadline) {
        $status = docker inspect $ContainerName --format '{{.State.Status}} {{if .State.Health}}{{.State.Health.Status}}{{else}}none{{end}}'
        if ($LASTEXITCODE -eq 0) {
            $parts = $status.Trim().Split(' ', [System.StringSplitOptions]::RemoveEmptyEntries)
            if ($parts.Count -ge 2 -and $parts[0] -eq 'running' -and $parts[1] -eq 'healthy') {
                return
            }
        }
        Start-Sleep -Seconds 3
    }

    throw "Container $ContainerName did not become healthy within $TimeoutSec seconds."
}

$repoRoot = Split-Path -Parent $PSScriptRoot
$outputDir = Join-Path $repoRoot 'output'
$binaryPath = Join-Path $outputDir 'sub2api-hotfix'
$goModCache = Join-Path $repoRoot '.tmp\go-mod-cache'
$goBuildCache = Join-Path $repoRoot '.tmp\go-build-cache'

New-Item -ItemType Directory -Force -Path $outputDir, $goModCache, $goBuildCache | Out-Null

if (-not $SkipTypecheck) {
    Invoke-Step -Name 'frontend typecheck' -Action {
        & pnpm --dir frontend typecheck
        if ($LASTEXITCODE -ne 0) { throw 'frontend typecheck failed.' }
    }
}

if (-not $SkipFrontendBuild) {
    Invoke-Step -Name 'frontend build' -Action {
        & pnpm --dir frontend build
        if ($LASTEXITCODE -ne 0) { throw 'frontend build failed.' }
    }
}

Invoke-Step -Name 'linux hotfix binary build' -Action {
    & docker run --rm `
        -e GOPROXY=https://goproxy.cn,direct `
        -e GOSUMDB=sum.golang.google.cn `
        -v "${repoRoot}:/app" `
        -v "${goModCache}:/go/pkg/mod" `
        -v "${goBuildCache}:/root/.cache/go-build" `
        -w /app/backend `
        $BuilderImage `
        sh -lc "CGO_ENABLED=0 GOOS=linux /usr/local/go/bin/go build -tags embed -ldflags='-s -w -X main.BuildType=release' -trimpath -o /app/output/sub2api-hotfix ./cmd/server"
    if ($LASTEXITCODE -ne 0) { throw 'linux hotfix binary build failed.' }
}

if (-not (Test-Path -LiteralPath $binaryPath)) {
    throw "Expected binary not found: $binaryPath"
}

$backupSuffix = Get-Date -Format 'yyyy-MM-dd-HHmmss'

Invoke-Step -Name 'copy binary into container and restart' -Action {
    & docker cp $binaryPath "${ContainerName}:/app/sub2api.new"
    if ($LASTEXITCODE -ne 0) { throw 'docker cp failed.' }

    & docker exec $ContainerName sh -lc "cp /app/sub2api /app/sub2api.bak-$backupSuffix && mv /app/sub2api.new /app/sub2api && chmod +x /app/sub2api"
    if ($LASTEXITCODE -ne 0) { throw 'container binary swap failed.' }

    & docker restart $ContainerName
    if ($LASTEXITCODE -ne 0) { throw 'container restart failed.' }
}

Invoke-Step -Name 'wait for healthy container' -Action {
    Wait-ForHealthy -ContainerName $ContainerName -TimeoutSec $HealthTimeoutSec
}

if (-not $SkipSmoke) {
    $envMap = Read-EnvFile -Path $EnvFile
    $email = if ($envMap.ContainsKey('ADMIN_EMAIL')) { $envMap['ADMIN_EMAIL'] } else { '' }
    $password = if ($envMap.ContainsKey('ADMIN_PASSWORD')) { $envMap['ADMIN_PASSWORD'] } else { '' }

    Invoke-Step -Name 'site smoke check' -Action {
        & (Join-Path $PSScriptRoot 'site-smoke-check.ps1') -BaseUrl $BaseUrl -Email $email -Password $password
        if ($LASTEXITCODE -ne 0) { throw 'site smoke check failed.' }
    }
}

Write-Host "Hotfix redeploy completed."
