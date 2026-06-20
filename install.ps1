# one-liner: irm https://raw.githubusercontent.com/hyperpuncher/yd-dl/main/install.ps1 | iex
$dest = "$env:USERPROFILE\AppData\Local\Microsoft\WindowsApps\yd-dl.exe"
$url = "https://github.com/hyperpuncher/yd-dl/releases/latest/download/yd-dl-windows-amd64.exe"
Write-Host "→ $url"
Invoke-WebRequest -Uri $url -OutFile $dest
Write-Host "→ yd-dl installed to $dest"
