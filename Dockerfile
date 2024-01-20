FROM golang:1.21.6

WORKDIR /app

COPY *.go ./
COPY *.mod ./
COPY *.sum ./
COPY log ./log
COPY routes ./routes
COPY shell ./shell
COPY static ./static
COPY util ./util
COPY views ./views

EXPOSE 5566 
RUN go build .

ARG PASSWORD
ENV PASSWORD $PASSWORD

ENTRYPOINT ["/app/ezcat"]
