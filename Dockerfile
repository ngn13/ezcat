####################
# build the server #
####################
FROM golang:1.22.4 as go

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
ENV STATIC_DIR "./static"
ENV PAYLOAD_DIR "./payloads"

ARG API_URL
ENV API_URL $API_URL

ENTRYPOINT ["dumb-init", "./init.sh"]
