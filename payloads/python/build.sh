#!/bin/bash

encoded=$(base64 src/main.py)
encoded=$(echo "$encoded" | sed -z 's/\n/_NL_/g')

mkdir -pv dist

cat > dist/main.py << EOF
e=str("$encoded").replace("_NL_", "\n")
exec(__import__("base64").b64decode(e.encode()))
EOF

final=$(base64 < dist/main.py | sed -z 's/\n/_NL_/g')
echo "from base64 import b64decode; exec(b64decode(str('$final').replace('_NL_', '\n')))" > dist/main.py
