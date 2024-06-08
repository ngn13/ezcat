#!/bin/sh

# update the API_URL for the frontend APP
if [ ! -z "$API_URL" ]; then
  API_URL=$(echo $API_URL | sed 's/\//\\\//g')
  find ./static -type f -exec sed -i -e "s/http:\/\/127.0.0.1:5566/$API_URL/g" {} \;
fi

# run the server
./server
