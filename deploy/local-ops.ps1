param(
    [Parameter(Mandatory = $true)]
    [ValidateSet('check-entry', 'entry', 'open-entry', 'smoke', 'redeploy', 'release-image')]
    [string]$Action,
    [string]$BaseUrl = 'http://127.0.0.1:8080',
    [string]$Email = '',
    [string]$Password = '',
    [string]$ContainerName = 'sub2api',
    [string]$EnvFile = '',
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
}
