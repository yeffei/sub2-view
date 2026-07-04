param(
    [string]$ContainerName = 'sub2api',
    [string]$EnvFile = (Join-Path $PSScriptRoot '.env')
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

function Get-LiveBinding {
    param(
        [Parameter(Mandatory = $true)]
        [string]$ContainerName
    )

    $json = docker inspect $ContainerName --format '{{json .HostConfig.PortBindings}}'
    if ($LASTEXITCODE -ne 0 -or [string]::IsNullOrWhiteSpace($json)) {
        throw "Failed to inspect container port bindings for $ContainerName."
    }

    $bindings = $json | ConvertFrom-Json
    $entry = $bindings.'8080/tcp'
    if (-not $entry -or $entry.Count -eq 0) {
        throw "Container $ContainerName does not expose 8080/tcp."
    }

    return $entry[0]
}

$envMap = Read-EnvFile -Path $EnvFile
$configuredBindHost = if ($envMap.ContainsKey('BIND_HOST')) { $envMap['BIND_HOST'] } else { '0.0.0.0' }
$configuredPort = if ($envMap.ContainsKey('SERVER_PORT')) { $envMap['SERVER_PORT'] } else { '8080' }

$liveBinding = Get-LiveBinding -ContainerName $ContainerName
$liveHostIp = if ([string]::IsNullOrWhiteSpace($liveBinding.HostIp)) { '0.0.0.0' } else { $liveBinding.HostIp }
$liveHostPort = $liveBinding.HostPort

$ipv4List = Get-NetIPAddress -AddressFamily IPv4 |
    Where-Object {
        $_.IPAddress -notmatch '^127\.' -and
        $_.IPAddress -notmatch '^169\.254\.' -and
        $_.PrefixOrigin -ne 'WellKnown'
    } |
    Select-Object -ExpandProperty IPAddress -Unique

$scope = if ($liveHostIp -eq '127.0.0.1') { 'local-only' } else { 'all-interfaces' }
$canonicalLocalUrl = "http://127.0.0.1:$liveHostPort"
$urls = @()
if ($scope -eq 'local-only') {
    $urls += $canonicalLocalUrl
} else {
    $urls += $canonicalLocalUrl
    foreach ($ip in $ipv4List) {
        $urls += "http://${ip}:$liveHostPort"
    }
}

$result = [PSCustomObject]@{
    ContainerName       = $ContainerName
    EnvFile             = $EnvFile
    ConfiguredBindHost  = $configuredBindHost
    ConfiguredPort      = $configuredPort
    LiveHostIp          = $liveHostIp
    LiveHostPort        = $liveHostPort
    CanonicalLocalUrl   = $canonicalLocalUrl
    RecommendedBrowserUrl = $canonicalLocalUrl
    ExposureScope       = $scope
    ReachableUrls       = $urls
    AddressPolicy       = if ($liveHostPort -eq '8080') { 'fixed-local-entry' } else { 'nonstandard-local-port' }
    AddressNote         = if ($liveHostPort -eq '8080') {
        'Use this localhost URL as the single local browser entry. Other dev ports such as 3000 are internal-only and should not be treated as the site address.'
    } else {
        'Local site is not running on the canonical 8080 entry. Restore SERVER_PORT=8080 if you want a single stable local address.'
    }
    RecreateCommand     = "docker compose --env-file deploy/.env -f deploy/docker-compose.local.yml up -d"
    ExternalAccessNote  = if ($scope -eq 'local-only') {
        'Current binding is local-only. Set BIND_HOST=0.0.0.0 in deploy/.env and recreate the service if LAN or public access is required.'
    } else {
        'Current binding is not limited to localhost. Confirm host firewall, reverse proxy, and TLS before exposing this service beyond a trusted network.'
    }
}

$result | ConvertTo-Json -Depth 4
