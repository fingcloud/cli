#!/usr/bin/env pwsh
# Copyright 2018 the Fing authors. All rights reserved. MIT license.
# TODO(everyone): Keep this script simple and easily auditable.

$ErrorActionPreference = 'Stop'

if ($v) {
  $Version = "v${v}"
}
if ($args.Length -eq 1) {
  $Version = $args.Get(0)
}

$FingInstall = $env:FING_INSTALL
$BinDir = if ($FingInstall) {
  "$FingInstall\bin"
} else {
  "$Home\.fing\bin"
}

$FingZip = "$BinDir\fing.zip"
$FingExe = "$BinDir\fing.exe"
$Target = 'windows-amd64'

# GitHub requires TLS 1.2
[Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12

$FingUri = if (!$Version) {
  "https://github.com/fingcloud/cli/releases/latest/download/fing-${Target}.zip"
} else {
  "https://github.com/fingcloud/cli/releases/download/${Version}/fing-${Target}.zip"
}

if (!(Test-Path $BinDir)) {
  New-Item $BinDir -ItemType Directory | Out-Null
}

Invoke-WebRequest $FingUri -OutFile $FingZip -UseBasicParsing

if (Get-Command Expand-Archive -ErrorAction SilentlyContinue) {
  Expand-Archive $FingZip -Destination $BinDir -Force
} else {
  if (Test-Path $FingExe) {
    Remove-Item $FingExe
  }
  Add-Type -AssemblyName System.IO.Compression.FileSystem
  [IO.Compression.ZipFile]::ExtractToDirectory($FingZip, $BinDir)
}

Remove-Item $FingZip

$User = [EnvironmentVariableTarget]::User
$Path = [Environment]::GetEnvironmentVariable('Path', $User)
if (!(";$Path;".ToLower() -like "*;$BinDir;*".ToLower())) {
  [Environment]::SetEnvironmentVariable('Path', "$Path;$BinDir", $User)
  $Env:Path += ";$BinDir"
}

Write-Output "Fing CLI was installed successfully to $FingExe"
Write-Output "Run 'fing help' to get started"