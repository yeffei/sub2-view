param(
    [string]$BaseUrl = 'http://127.0.0.1:8080',
    [string]$Email = '',
    [string]$Password = '',
    [int]$TimeoutSec = 20
)

Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'

function Add-CheckResult {
    param(
        [System.Collections.Generic.List[object]]$Results,
        [string]$Name,
        [bool]$Passed,
        [string]$Detail
    )

    $Results.Add([PSCustomObject]@{
        Name   = $Name
        Passed = $Passed
        Detail = $Detail
    })
}

function Invoke-TextRequest {
    param(
        [Parameter(Mandatory = $true)]
        [ValidateSet('Get', 'Post')]
        [string]$Method,
        [Parameter(Mandatory = $true)]
        [string]$Url,
        [hashtable]$Headers,
        [string]$Body,
        [int]$TimeoutSec = 20
    )

    $params = @{
        Method      = $Method
        Uri         = $Url
        TimeoutSec  = $TimeoutSec
        UseBasicParsing = $true
    }
    if ($Headers) {
        $params.Headers = $Headers
    }
    if ($Body) {
        $params.ContentType = 'application/json'
        $params.Body = $Body
    }

    return Invoke-WebRequest @params
}

$results = New-Object 'System.Collections.Generic.List[object]'
$base = $BaseUrl.TrimEnd('/')
$routes = @(
    '/',
    '/home',
    '/login',
    '/register',
    '/dashboard',
    '/keys',
    '/usage',
    '/profile',
    '/admin/dashboard',
    '/admin/settings'
)

$rootResponse = $null
foreach ($route in $routes) {
    try {
        $response = Invoke-TextRequest -Method Get -Url ($base + $route) -TimeoutSec $TimeoutSec
        Add-CheckResult -Results $results -Name "route:$route" -Passed ($response.StatusCode -eq 200) -Detail "status=$($response.StatusCode) length=$($response.Content.Length)"
        if ($route -eq '/') {
            $rootResponse = $response
        }
    } catch {
        Add-CheckResult -Results $results -Name "route:$route" -Passed $false -Detail $_.Exception.Message
    }
}

if ($null -ne $rootResponse) {
    $assetMatches = [regex]::Matches($rootResponse.Content, '(?:src|href)="(?<path>/(?:assets/[^"]+|logo\.png))"')
    $assetPaths = $assetMatches |
        ForEach-Object { $_.Groups['path'].Value } |
        Select-Object -Unique

    foreach ($assetPath in $assetPaths) {
        try {
            $assetResponse = Invoke-TextRequest -Method Get -Url ($base + $assetPath) -TimeoutSec $TimeoutSec
            Add-CheckResult -Results $results -Name "asset:$assetPath" -Passed ($assetResponse.StatusCode -eq 200) -Detail "status=$($assetResponse.StatusCode) length=$($assetResponse.Content.Length)"
        } catch {
            Add-CheckResult -Results $results -Name "asset:$assetPath" -Passed $false -Detail $_.Exception.Message
        }
    }
}

$token = $null
if (-not [string]::IsNullOrWhiteSpace($Email) -and -not [string]::IsNullOrWhiteSpace($Password)) {
    try {
        $loginBody = @{
            email    = $Email
            password = $Password
        } | ConvertTo-Json
        $login = Invoke-RestMethod -Method Post -Uri ($base + '/api/v1/auth/login') -ContentType 'application/json' -Body $loginBody -TimeoutSec $TimeoutSec
        $token = $login.data.access_token
        Add-CheckResult -Results $results -Name 'api:login' -Passed ($login.code -eq 0 -and -not [string]::IsNullOrWhiteSpace($token)) -Detail "code=$($login.code) message=$($login.message)"
    } catch {
        Add-CheckResult -Results $results -Name 'api:login' -Passed $false -Detail $_.Exception.Message
    }
}

if ($token) {
    $headers = @{ Authorization = "Bearer $token" }

    try {
        $me = Invoke-RestMethod -Method Get -Uri ($base + '/api/v1/auth/me') -Headers $headers -TimeoutSec $TimeoutSec
        Add-CheckResult -Results $results -Name 'api:auth/me' -Passed ($me.code -eq 0) -Detail "code=$($me.code) message=$($me.message)"
    } catch {
        Add-CheckResult -Results $results -Name 'api:auth/me' -Passed $false -Detail $_.Exception.Message
    }

    $apiKeys = $null
    try {
        $apiKeys = Invoke-RestMethod -Method Get -Uri ($base + '/api/v1/keys?page=1&page_size=1') -Headers $headers -TimeoutSec $TimeoutSec
        $count = if ($apiKeys.data.items) { $apiKeys.data.items.Count } else { 0 }
        Add-CheckResult -Results $results -Name 'api:keys' -Passed ($apiKeys.code -eq 0) -Detail "code=$($apiKeys.code) items=$count"
    } catch {
        Add-CheckResult -Results $results -Name 'api:keys' -Passed $false -Detail $_.Exception.Message
    }

    $firstKeyId = $null
    if ($apiKeys -and $apiKeys.data -and $apiKeys.data.items -and $apiKeys.data.items.Count -gt 0) {
        $firstKeyId = $apiKeys.data.items[0].id
    }

    if ($firstKeyId) {
        try {
            $workbenchBody = @{ api_key_ids = @($firstKeyId) } | ConvertTo-Json
            $workbench = Invoke-RestMethod -Method Post -Uri ($base + '/api/v1/usage/dashboard/api-keys-workbench') -Headers $headers -ContentType 'application/json' -Body $workbenchBody -TimeoutSec $TimeoutSec
            $hasStats = $workbench.code -eq 0 -and $null -ne $workbench.data.stats
            Add-CheckResult -Results $results -Name 'api:workbench' -Passed $hasStats -Detail "code=$($workbench.code) api_key_id=$firstKeyId"
        } catch {
            Add-CheckResult -Results $results -Name 'api:workbench' -Passed $false -Detail $_.Exception.Message
        }
    } else {
        Add-CheckResult -Results $results -Name 'api:workbench' -Passed $true -Detail 'skipped: no api keys available for this user'
    }

    try {
        $adminSettings = Invoke-RestMethod -Method Get -Uri ($base + '/api/v1/admin/settings') -Headers $headers -TimeoutSec $TimeoutSec
        Add-CheckResult -Results $results -Name 'api:admin/settings' -Passed ($adminSettings.code -eq 0) -Detail "code=$($adminSettings.code) message=$($adminSettings.message)"
    } catch {
        Add-CheckResult -Results $results -Name 'api:admin/settings' -Passed $false -Detail $_.Exception.Message
    }
}

$results | Format-Table -AutoSize | Out-String | Write-Host

$failed = @($results | Where-Object { -not $_.Passed })
if ($failed.Count -gt 0) {
    Write-Error ("Smoke check failed: " + (($failed | ForEach-Object { $_.Name }) -join ', '))
}
