param(
    [string]$BaseUrl = 'http://127.0.0.1:18080',
    [string]$ContainerName = 'sub2api',
    [string]$EnvFile = (Join-Path $PSScriptRoot '.env'),
    [string]$ApiKey = '',
    [string]$Model = 'gpt-5.4',
    [int]$Concurrency = 4,
    [int]$RequestsPerWorker = 5,
    [int]$RequestTimeoutSec = 120,
    [int]$StatsIntervalSec = 2,
    [switch]$NonStream
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

function Invoke-ContainerPsqlScalar {
    param(
        [Parameter(Mandatory = $true)]
        [string]$TargetContainer,
        [Parameter(Mandatory = $true)]
        [hashtable]$EnvMap,
        [Parameter(Mandatory = $true)]
        [string]$Query
    )

    foreach ($requiredKey in @('POSTGRES_PASSWORD', 'POSTGRES_USER', 'POSTGRES_DB')) {
        if (-not $EnvMap.ContainsKey($requiredKey) -or [string]::IsNullOrWhiteSpace($EnvMap[$requiredKey])) {
            throw "Missing $requiredKey in $EnvFile"
        }
    }

    $queryEscaped = $Query.Replace('"', '\"')
    $command = "PGPASSWORD=$($EnvMap['POSTGRES_PASSWORD']) psql -t -A -h postgres -U $($EnvMap['POSTGRES_USER']) -d $($EnvMap['POSTGRES_DB']) -c `"$queryEscaped`""
    $result = & docker exec $TargetContainer /bin/sh -lc $command
    if ($LASTEXITCODE -ne 0) {
        throw "docker exec psql failed for container $TargetContainer"
    }
    return ($result | Out-String).Trim()
}

function Resolve-ApiKey {
    param(
        [string]$ExplicitApiKey,
        [Parameter(Mandatory = $true)]
        [string]$TargetContainer,
        [Parameter(Mandatory = $true)]
        [hashtable]$EnvMap
    )

    if (-not [string]::IsNullOrWhiteSpace($ExplicitApiKey)) {
        return [PSCustomObject]@{
            Value  = $ExplicitApiKey.Trim()
            Source = 'parameter'
        }
    }

    $query = "select key from api_keys where status='active' and deleted_at is null and key like 'sk-%' order by last_used_at desc nulls last, id desc limit 1;"
    $resolved = Invoke-ContainerPsqlScalar -TargetContainer $TargetContainer -EnvMap $EnvMap -Query $query
    if ([string]::IsNullOrWhiteSpace($resolved)) {
        throw 'No active local api key could be auto-resolved. Pass -ApiKey explicitly.'
    }

    return [PSCustomObject]@{
        Value  = $resolved
        Source = 'local_db'
    }
}

function Get-ContainerInspectSummary {
    param(
        [Parameter(Mandatory = $true)]
        [string]$TargetContainer
    )

    $raw = & docker inspect $TargetContainer --format '{{json .}}'
    if ($LASTEXITCODE -ne 0) {
        throw "docker inspect failed for container $TargetContainer"
    }
    $inspect = $raw | ConvertFrom-Json
    return [PSCustomObject]@{
        RestartCount = [int]$inspect.RestartCount
        Status       = $inspect.State.Status
        Running      = [bool]$inspect.State.Running
        OOMKilled    = [bool]$inspect.State.OOMKilled
        ExitCode     = [int]$inspect.State.ExitCode
        StartedAt    = $inspect.State.StartedAt
        FinishedAt   = $inspect.State.FinishedAt
    }
}

function Invoke-StabilityRequest {
    param(
        [Parameter(Mandatory = $true)]
        [string]$Url,
        [Parameter(Mandatory = $true)]
        [string]$ResolvedApiKey,
        [Parameter(Mandatory = $true)]
        [string]$TargetModel,
        [Parameter(Mandatory = $true)]
        [int]$WorkerIndex,
        [Parameter(Mandatory = $true)]
        [int]$RequestIndex,
        [Parameter(Mandatory = $true)]
        [bool]$Stream,
        [Parameter(Mandatory = $true)]
        [int]$TimeoutSec
    )

    $bodyObject = @{
        model = $TargetModel
        input = @(
            @{
                role = 'user'
                content = @(
                    @{
                        type = 'input_text'
                        text = "Stability probe worker $WorkerIndex request $RequestIndex. Reply with OK only."
                    }
                )
            }
        )
    }
    if ($Stream) {
        $bodyObject.stream = $true
    }

    $body = $bodyObject | ConvertTo-Json -Depth 6 -Compress
    $headers = @{
        Authorization = "Bearer $ResolvedApiKey"
        'Content-Type' = 'application/json'
    }

    $stopwatch = [System.Diagnostics.Stopwatch]::StartNew()
    try {
        $response = Invoke-WebRequest -Method Post -Uri $Url -Headers $headers -Body $body -TimeoutSec $TimeoutSec -UseBasicParsing
        $stopwatch.Stop()
        return [PSCustomObject]@{
            Worker      = $WorkerIndex
            Request     = $RequestIndex
            Passed      = ($response.StatusCode -eq 200 -and $response.Content -match 'OK')
            StatusCode  = $response.StatusCode
            DurationMs  = [int]$stopwatch.ElapsedMilliseconds
            Error       = ''
        }
    } catch {
        $stopwatch.Stop()
        $statusCode = 0
        if ($_.Exception.Response -and $_.Exception.Response.StatusCode) {
            $statusCode = [int]$_.Exception.Response.StatusCode
        }
        return [PSCustomObject]@{
            Worker      = $WorkerIndex
            Request     = $RequestIndex
            Passed      = $false
            StatusCode  = $statusCode
            DurationMs  = [int]$stopwatch.ElapsedMilliseconds
            Error       = $_.Exception.Message
        }
    }
}

$repoRoot = Split-Path -Parent $PSScriptRoot
$outputRoot = Join-Path $repoRoot 'output\stability'
$runStamp = Get-Date -Format 'yyyyMMdd-HHmmss'
$runDir = Join-Path $outputRoot $runStamp
$statsFile = Join-Path $runDir 'docker-stats.log'
$eventsFile = Join-Path $runDir 'docker-events.log'
$logsFile = Join-Path $runDir 'docker-logs.log'
$resultsFile = Join-Path $runDir 'request-results.json'
$summaryFile = Join-Path $runDir 'summary.json'
$stopFile = Join-Path $runDir 'stop.signal'
$streamEnabled = -not $NonStream.IsPresent

New-Item -ItemType Directory -Force -Path $runDir | Out-Null
$envMap = Read-EnvFile -Path $EnvFile
$apiKeyResolution = Resolve-ApiKey -ExplicitApiKey $ApiKey -TargetContainer $ContainerName -EnvMap $envMap
$resolvedApiKey = $apiKeyResolution.Value

$inspectBefore = Get-ContainerInspectSummary -TargetContainer $ContainerName
$startUtc = (Get-Date).ToUniversalTime().ToString('o')
$base = $BaseUrl.TrimEnd('/')
$requestUrl = $base + '/v1/responses'

$statsJob = Start-Job -ScriptBlock {
    param($TargetContainer, $TargetPath, $TargetStopFile, $IntervalSec)
    while (-not (Test-Path -LiteralPath $TargetStopFile)) {
        $timestamp = (Get-Date).ToString('yyyy-MM-dd HH:mm:ss')
        $line = & docker stats --no-stream --format '{{.Name}}|{{.CPUPerc}}|{{.MemUsage}}|{{.MemPerc}}|{{.NetIO}}' $TargetContainer 2>$null
        "$timestamp|$line" | Out-File -FilePath $TargetPath -Append -Encoding utf8
        Start-Sleep -Seconds $IntervalSec
    }
} -ArgumentList $ContainerName, $statsFile, $stopFile, $StatsIntervalSec

$workerScript = {
    param($Url, $ResolvedApiKey, $TargetModel, $WorkerIndex, $WorkerRequests, $StreamEnabled, $TimeoutSec)

    $items = New-Object 'System.Collections.Generic.List[object]'
    for ($requestIndex = 1; $requestIndex -le $WorkerRequests; $requestIndex++) {
        $bodyObject = @{
            model = $TargetModel
            input = @(
                @{
                    role = 'user'
                    content = @(
                        @{
                            type = 'input_text'
                            text = "Stability probe worker $WorkerIndex request $RequestIndex. Reply with OK only."
                        }
                    )
                }
            )
        }
        if ($StreamEnabled) {
            $bodyObject.stream = $true
        }

        $body = $bodyObject | ConvertTo-Json -Depth 6 -Compress
        $headers = @{
            Authorization = "Bearer $ResolvedApiKey"
            'Content-Type' = 'application/json'
        }

        $stopwatch = [System.Diagnostics.Stopwatch]::StartNew()
        try {
            $response = Invoke-WebRequest -Method Post -Uri $Url -Headers $headers -Body $body -TimeoutSec $TimeoutSec -UseBasicParsing
            $stopwatch.Stop()
            $items.Add([PSCustomObject]@{
                Worker      = $WorkerIndex
                Request     = $requestIndex
                Passed      = ($response.StatusCode -eq 200 -and $response.Content -match 'OK')
                StatusCode  = $response.StatusCode
                DurationMs  = [int]$stopwatch.ElapsedMilliseconds
                Error       = ''
            })
        } catch {
            $stopwatch.Stop()
            $statusCode = 0
            if ($_.Exception.Response -and $_.Exception.Response.StatusCode) {
                $statusCode = [int]$_.Exception.Response.StatusCode
            }
            $items.Add([PSCustomObject]@{
                Worker      = $WorkerIndex
                Request     = $requestIndex
                Passed      = $false
                StatusCode  = $statusCode
                DurationMs  = [int]$stopwatch.ElapsedMilliseconds
                Error       = $_.Exception.Message
            })
        }
    }

    return $items
}

$jobs = 1..$Concurrency | ForEach-Object {
    Start-Job -ScriptBlock $workerScript -ArgumentList $requestUrl, $resolvedApiKey, $Model, $_, $RequestsPerWorker, $streamEnabled, $RequestTimeoutSec
}

$requestResults = @($jobs | Wait-Job | Receive-Job)
New-Item -ItemType File -Force -Path $stopFile | Out-Null
Wait-Job $statsJob | Out-Null
$null = Receive-Job $statsJob

$endUtc = (Get-Date).ToUniversalTime().ToString('o')
$events = & docker events --since $startUtc --until $endUtc --filter "container=$ContainerName" 2>$null
$logs = & docker logs --since $startUtc $ContainerName 2>&1
$inspectAfter = Get-ContainerInspectSummary -TargetContainer $ContainerName

$events | Out-File -FilePath $eventsFile -Encoding utf8
$logs | Out-File -FilePath $logsFile -Encoding utf8
$requestResults | ConvertTo-Json -Depth 5 | Out-File -FilePath $resultsFile -Encoding utf8

$passedResults = @($requestResults | Where-Object { $_.Passed })
$failedResults = @($requestResults | Where-Object { -not $_.Passed })
$avgMs = 0
$maxMs = 0
if ($requestResults.Count -gt 0) {
    $avgMs = [int](($requestResults | Measure-Object -Property DurationMs -Average).Average)
    $maxMs = [int](($requestResults | Measure-Object -Property DurationMs -Maximum).Maximum)
}

$summary = [PSCustomObject]@{
    started_utc            = $startUtc
    finished_utc           = $endUtc
    base_url               = $base
    request_url            = $requestUrl
    container_name         = $ContainerName
    model                  = $Model
    stream                 = $streamEnabled
    concurrency            = $Concurrency
    requests_per_worker    = $RequestsPerWorker
    total_requests         = $requestResults.Count
    passed_requests        = $passedResults.Count
    failed_requests        = $failedResults.Count
    average_duration_ms    = $avgMs
    max_duration_ms        = $maxMs
    api_key_source         = $apiKeyResolution.Source
    restart_count_before   = $inspectBefore.RestartCount
    restart_count_after    = $inspectAfter.RestartCount
    oom_killed_after       = [bool]$inspectAfter.OOMKilled
    container_status_after = $inspectAfter.Status
    output_dir             = $runDir
    artifacts              = [PSCustomObject]@{
        summary_json  = $summaryFile
        results_json  = $resultsFile
        stats_log     = $statsFile
        events_log    = $eventsFile
        container_log = $logsFile
    }
}

$summary | ConvertTo-Json -Depth 6 | Out-File -FilePath $summaryFile -Encoding utf8

Write-Host "=== stability summary ==="
$summary | Format-List | Out-String | Write-Host

if ($failedResults.Count -gt 0 -or $inspectAfter.RestartCount -gt $inspectBefore.RestartCount -or $inspectAfter.OOMKilled) {
    throw "Local stability check failed. See artifacts under $runDir"
}

Write-Host "Artifacts written to: $runDir"
