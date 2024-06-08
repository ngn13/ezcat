#!/bin/bash

rm -f stage

if [ "$1" == "windows_amd64" ]; then
  make CFLAGS="-O3 -s -static -lws2_32" CC=x86_64-w64-mingw32-gcc
  mv stage.exe stage
  strip --strip-unneeded stage
  exit $?
elif [ "$1" == "linux_amd64" ]; then
  make CFLAGS="-O3 -s -static" CC=gcc
  strip --strip-unneeded stage
  exit $?
fi

echo "Unknown build target"
exit 1
