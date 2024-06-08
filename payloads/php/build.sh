#!/bin/bash

encoded=$(base64 src/main.php)
encoded=$(echo "$encoded" | sed -z 's/\n/_NL_/g')

mkdir -pv dist

cat > dist/main.php << EOF
<?php eval(base64_decode(str_replace("_NL_", "\n", "$encoded"))) ?>
EOF
