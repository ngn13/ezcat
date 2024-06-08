#!/bin/bash -e

level=$((1 + RANDOM % 10))

for i in {1 .$level}; do
  if [ -z "$encoded" ]; then
    encoded=$(base64 src/main.sh)
  else
    encoded=$(echo "$encoded" | base64 | sed -z 's/\n/_NL_/g')
  fi
done

mkdir -pv dist

cat > dist/main.sh << EOF
for i in {1 .$level}; do
  if [ -z "\$p" ]; then
    p=\$(echo '$encoded' | sed -z 's/_NL_/\\n/g'| base64 -d)
  else
    p=\$(echo "\$p" | base64 -d)
  fi
done
echo "\$p" | bash
EOF

final="$(base64 < dist/main.sh | sed -z 's/\n/_NL_/g')"
echo "echo \"$final\" | sed -z 's/_NL_/\\n/g' | base64 -d | bash" > dist/main.sh
