#!/bin/bash
echo ">> downloading ezcat binary..."
wget https://github.com/ngn13/ezcat/releases/latest/download/ezcat
echo ">> download complete"
chmod +x ezcat 
mv ezcat /usr/bin/ezcat
echo ">> ezcat has been installed!"
