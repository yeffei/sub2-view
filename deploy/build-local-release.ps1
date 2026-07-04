param(
    [string]$ImageName = 'sub2api-local',
    [string]$ContainerName = 'sub2api',
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
        [string]$TargetContainer,
        [int]$TimeoutSec = 120
    )

    $deadline = (Get-Date).AddSeconds($TimeoutSec)
    while ((Get-Date) -lt $deadline) {
        $status = docker inspect $TargetContainer --format '{{.State.Status}} {{if .State.Health}}{{.State.Health.Status}}{{else}}none{{end}}'
        if ($LASTEXITCODE -eq 0) {
            $parts = $status.Trim().Split(' ', [System.StringSplitOptions]::RemoveEmptyEntries)
            if ($parts.Count -ge 2 -and $parts[0] -eq 'running' -and $parts[1] -eq 'healthy') {
                return
            }
        }
        Start-Sleep -Seconds 3
    }

    throw "Container $TargetContainer did not become healthy within $TimeoutSec seconds."
}

$repoRoot = Split-Path -Parent $PSScriptRoot
$outputDir = Join-Path $repoRoot 'output'
$binaryPath = Join-Path $outputDir 'sub2api-hotfix'
$goModCache = Join-Path $repoRoot '.tmp\go-mod-cache'
$goBuildCache = Join-Path $repoRoot '.tmp\go-build-cache'
$baseCacheTag = "${ImageName}:basecache"
$releaseTimestamp = Get-Date -Format 'yyyyMMdd-HHmmss'
$releaseTag = "${ImageName}:release-$releaseTimestamp"

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
        golang:1.26.4-alpine `
        sh -lc "CGO_ENABLED=0 GOOS=linux /usr/local/go/bin/go build -tags embed -ldflags='-s -w -X main.BuildType=release' -trimpath -o /app/output/sub2api-hotfix ./cmd/server"
    if ($LASTEXITCODE -ne 0) { throw 'linux hotfix binary build failed.' }
}

if (-not (Test-Path -LiteralPath $binaryPath)) {
    throw "Expected binary not found: $binaryPath"
}

Invoke-Step -Name 'retag current runtime image as basecache' -Action {
    & docker image inspect "${ImageName}:latest" | Out-Null
    if ($LASTEXITCODE -ne 0) {
        throw "Base image ${ImageName}:latest not found."
    }
    & docker tag "${ImageName}:latest" $baseCacheTag
    if ($LASTEXITCODE -ne 0) { throw 'failed to tag basecache image.' }
}

Invoke-Step -Name 'build persistent local release image' -Action {
    & docker build `
        -f (Join-Path $PSScriptRoot 'Dockerfile.local-release') `
        -t $releaseTag `
        -t "${ImageName}:latest" `
        --build-arg "BASE_IMAGE=$baseCacheTag" `
        $repoRoot
    if ($LASTEXITCODE -ne 0) { throw 'local release image build failed.' }
}

Invoke-Step -Name 'recreate service from image' -Action {
    & docker compose --env-file $EnvFile -f (Join-Path $PSScriptRoot 'docker-compose.local.yml') up -d --force-recreate $ContainerName
    if ($LASTEXITCODE -ne 0) { throw 'docker compose recreate failed.' }
}

Invoke-Step -Name 'wait for healthy container' -Action {
    Wait-ForHealthy -TargetContainer $ContainerName -TimeoutSec $HealthTimeoutSec
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

Write-Host "Local release image completed: $releaseTag"
