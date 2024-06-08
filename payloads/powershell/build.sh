#!/bin/bash -e

plain=$(cat src/main.ps1 | sed -z 's/\n/;/g')
encoded=$(echo "$plain" | base64 | sed -z 's/\n/_NL_/g')

mkdir -pv dist

cat > dist/main.ps1 << EOF
Set-ExecutionPolicy Bypass -scope Process -Force; \$d=[System.Text.Encoding]::UTF8.GetString([System.Convert]::FromBase64String("$encoded".Replace("_NL_", "\`n"))); Invoke-Expression "\$d"
EOF
