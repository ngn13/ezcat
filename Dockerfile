####################
# build the server #
####################
FROM golang:1.23.5 as go

COPY server /server
WORKDIR /server

RUN CGO_ENABLED=0 go build

#####################
# build the web app #
#####################
FROM node as node

COPY app /app
WORKDIR /app

ENV VITE_API_URL_DEV "http://127.0.0.1:5566"
RUN npm install
RUN npm run build

#####################
# the actual runner #
#####################
FROM alpine as main

RUN apk add sed bash build-base dumb-init gcc mingw-w64-gcc

WORKDIR /ezcat

COPY --from=node /app/build     ./static
COPY --from=go   /server/server ./

COPY payloads       ./payloads
COPY docker/init.sh ./

RUN chmod +x "init.sh"
ENV EZCAT_STATIC_DIR "./static"
ENV EZCAT_PAYLOAD_DIR "./payloads"

ENTRYPOINT ["dumb-init", "./init.sh"]
