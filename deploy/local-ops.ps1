param(
    [Parameter(Mandatory = $true)]
    [ValidateSet('check-entry', 'entry', 'open-entry', 'smoke', 'redeploy', 'release-image', 'stability')]
    [string]$Action,
    [string]$BaseUrl = 'http://127.0.0.1:8080',
    [string]$Email = '',
    [string]$Password = '',
    [string]$ContainerName = 'sub2api',
    [string]$EnvFile = '',
    [string]$ApiKey = '',
    [string]$Model = 'gpt-5.4',
    [int]$Concurrency = 4,
    [int]$RequestsPerWorker = 5,
    [int]$RequestTimeoutSec = 120,
    [int]$StatsIntervalSec = 2,
    [switch]$NonStream,
    [switch]$SkipTypecheck,
    [switch]$SkipFrontendBuild,
    [switch]$SkipSmoke
)

Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'

$scriptRoot = Split-Path -Parent $MyInvocation.MyCommand.Path

switch ($Action) {
    'check-entry' {
        $checkParams = @{
            ContainerName = $ContainerName
        }
        if (-not [string]::IsNullOrWhiteSpace($EnvFile)) {
            $checkParams.EnvFile = $EnvFile
        }
        & (Join-Path $scriptRoot 'check-entry.ps1') @checkParams
        break
    }
    'entry' {
        $checkParams = @{
            ContainerName = $ContainerName
        }
        if (-not [string]::IsNullOrWhiteSpace($EnvFile)) {
            $checkParams.EnvFile = $EnvFile
        }
        $entryInfo = & (Join-Path $scriptRoot 'check-entry.ps1') @checkParams | ConvertFrom-Json
        Write-Output $entryInfo.RecommendedBrowserUrl
        break
    }
    'open-entry' {
        $checkParams = @{
            ContainerName = $ContainerName
        }
        if (-not [string]::IsNullOrWhiteSpace($EnvFile)) {
            $checkParams.EnvFile = $EnvFile
        }
        $entryInfo = & (Join-Path $scriptRoot 'check-entry.ps1') @checkParams | ConvertFrom-Json
        Start-Process $entryInfo.RecommendedBrowserUrl
        Write-Host "Opened $($entryInfo.RecommendedBrowserUrl)"
        break
    }
    'smoke' {
        $params = @{
            BaseUrl = $BaseUrl
        }
        if (-not [string]::IsNullOrWhiteSpace($Email)) {
            $params.Email = $Email
        }
        if (-not [string]::IsNullOrWhiteSpace($Password)) {
            $params.Password = $Password
        }
        & (Join-Path $scriptRoot 'site-smoke-check.ps1') @params
        break
    }
    'redeploy' {
        $params = @{
            BaseUrl = $BaseUrl
        }
        if (-not [string]::IsNullOrWhiteSpace($ContainerName)) {
            $params.ContainerName = $ContainerName
        }
        if (-not [string]::IsNullOrWhiteSpace($EnvFile)) {
            $params.EnvFile = $EnvFile
        }
        if ($SkipTypecheck) {
            $params.SkipTypecheck = $true
        }
        if ($SkipFrontendBuild) {
            $params.SkipFrontendBuild = $true
        }
        if ($SkipSmoke) {
            $params.SkipSmoke = $true
        }
        & (Join-Path $scriptRoot 'hotfix-redeploy.ps1') @params
        break
    }
    'release-image' {
        $params = @{
            BaseUrl = $BaseUrl
        }
        if (-not [string]::IsNullOrWhiteSpace($ContainerName)) {
            $params.ContainerName = $ContainerName
        }
        if (-not [string]::IsNullOrWhiteSpace($EnvFile)) {
            $params.EnvFile = $EnvFile
        }
        if ($SkipTypecheck) {
            $params.SkipTypecheck = $true
        }
        if ($SkipFrontendBuild) {
            $params.SkipFrontendBuild = $true
        }
        if ($SkipSmoke) {
            $params.SkipSmoke = $true
        }
        & (Join-Path $scriptRoot 'build-local-release.ps1') @params
        break
    }
    'stability' {
        $params = @{
            BaseUrl = if ($PSBoundParameters.ContainsKey('BaseUrl')) { $BaseUrl } else { 'http://127.0.0.1:18080' }
            ContainerName = $ContainerName
            Model = $Model
            Concurrency = $Concurrency
            RequestsPerWorker = $RequestsPerWorker
            RequestTimeoutSec = $RequestTimeoutSec
            StatsIntervalSec = $StatsIntervalSec
        }
        if (-not [string]::IsNullOrWhiteSpace($EnvFile)) {
            $params.EnvFile = $EnvFile
        }
        if (-not [string]::IsNullOrWhiteSpace($ApiKey)) {
            $params.ApiKey = $ApiKey
        }
        if ($NonStream) {
            $params.NonStream = $true
        }
        & (Join-Path $scriptRoot 'local-stability-check.ps1') @params
        break
    }
}
