# install.ps1 — Install the latest gow CLI for Windows

$repo = "mechneerd/gow"
$installDir = "$env:ProgramFiles\gow"
$binaryName = "gow.exe"

Write-Host "Installing gow from $repo..."

$os = "windows"
$arch = if ([Environment]::Is64BitOperatingSystem) { "amd64" } else { "386" }

$url = "https://github.com/$repo/releases/latest/download/gow-$os-$arch.exe"

$tmp = "$env:TEMP\gow.exe"
Invoke-WebRequest -Uri $url -OutFile $tmp

if (-not (Test-Path $installDir)) {
    New-Item -ItemType Directory -Path $installDir | Out-Null
}

Move-Item -Force $tmp "$installDir\$binaryName"

$envPath = [Environment]::GetEnvironmentVariable("Path", "Machine")
if ($envPath -notlike "*$installDir*") {
    [Environment]::SetEnvironmentVariable("Path", "$envPath;$installDir", "Machine")
}

Write-Host "✅ gow installed successfully!"
& "$installDir\$binaryName" --version
