#!/usr/bin/env pwsh

[CmdletBinding()]
param (
  [Parameter()]
  [ValidateSet('build', 'quick', 'package')]
  [string]
  $Target = 'build'
)
$Env:WORKINGDIR = $PSScriptRoot

$arch = go env GOHOSTARCH
$platform = go env GOHOSTOS
$BuildFolderPath = Join-Path $PSScriptRoot 'dist' "flfa_${platform}_${arch}"
$NotelBuildFolderPath = Join-Path $PSScriptRoot 'dist' "notel_flfa_${platform}_${arch}"
$DocsSource = Join-Path $PSScriptRoot 'docs' 'content' '*'
$DataSource = Join-Path $PSScriptRoot 'modules' '*'

switch ($Target) {
  'build' {
    # Set goreleaser to build for current platform only
    # Add environment variables for honeycomb if not already loaded
    if (!(Test-Path ENV:\HONEYCOMB_API_KEY)) {
      $ENV:HONEYCOMB_API_KEY = 'not_set'
    }
    if (!(Test-Path ENV:\HONEYCOMB_DATASET)) {
      $ENV:HONEYCOMB_DATASET = 'not_set'
    }
    goreleaser build --snapshot --rm-dist --single-target
    # Copy docs content to package folder
    foreach ($BuildFolderDocs in @("$BuildFolderPath/docs", "$NotelBuildFolderPath/docs")) {
      $null = New-Item -Path $BuildFolderDocs -ItemType Directory
      $null = Copy-Item -Path $DocsSource -Destination $BuildFolderDocs -Recurse
    }
    # Copy docs content to package folder
    foreach ($BuildFolderData in @("$BuildFolderPath/modules", "$NotelBuildFolderPath/modules")) {
      $null = New-Item -Path $BuildFolderData -ItemType Directory
      $null = Copy-Item -Path $DataSource -Destination $BuildFolderData -Recurse
    }
  }
  'quick' {
    If ($Env:OS -match '^Windows') {
      go build -o "$BuildFolderPath/flfa.exe" -tags telemetry
      go build -o "$NotelBuildFolderPath/flfa.exe"
    } else {
      go build -o "$BuildFolderPath/flfa" -tags telemetry
      go build -o "$NotelBuildFolderPath/flfa"
    }
  }
  'package' {
    goreleaser --skip-publish --snapshot --rm-dist
  }
}