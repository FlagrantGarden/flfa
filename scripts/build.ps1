#!/usr/bin/env pwsh

[CmdletBinding()]
param (
  [Parameter()]
  [ValidateSet('build', 'quick', 'package')]
  [string]$Target = 'build'
)
$Env:WORKINGDIR = Split-Path -Parent $PSScriptRoot

$arch = go env GOHOSTARCH
$platform = go env GOHOSTOS
$BuildFolderPath = Join-Path $Env:WORKINGDIR 'dist' "flfa_${platform}_${arch}"
$NotelBuildFolderPath = Join-Path $Env:WORKINGDIR 'dist' "notel_flfa_${platform}_${arch}"

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