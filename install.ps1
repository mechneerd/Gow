# install.ps1 — Install the latest gow CLI for Windows

$repo = "mechneerd/gow"
$installDir = "$env:ProgramFiles\gow"
$binaryName = "gow.exe"

Write-Host "Installing gow from $repo..."

# Detect architecture (x86_64, arm64)
$arch = if ($env:PROCESSOR_ARCHITECTURE -eq "ARM64") { "arm64" } elseif ([Environment]::Is64BitOperatingSystem) { "x86_64" } else { "i386" }

$url = "https://github.com/$repo/releases/latest/download/gow_Windows_${arch}.zip"

$tmpZip = "$env:TEMP\gow.zip"
$tmpDir = "$env:TEMP\gow_extract"

Invoke-WebRequest -Uri $url -OutFile $tmpZip

if (Test-Path $tmpDir) { Remove-Item -Recurse -Force $tmpDir }
Expand-Archive -Path $tmpZip -DestinationPath $tmpDir -Force

if (-not (Test-Path $installDir)) {
    New-Item -ItemType Directory -Path $installDir | Out-Null
}

Move-Item -Force "$tmpDir\gow.exe" "$installDir\$binaryName"

Remove-Item -Recurse -Force $tmpDir
Remove-Item -Force $tmpZip

$envPath = [Environment]::GetEnvironmentVariable("Path", "Machine")
if ($envPath -notlike "*$installDir*") {
    [Environment]::SetEnvironmentVariable("Path", "$envPath;$installDir", "Machine")
}

Write-Host "gow installed successfully!"
& "$installDir\$binaryName" --version
