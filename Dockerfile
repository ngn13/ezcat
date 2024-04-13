FROM golang:1.22.2 as build

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
RUN CGO_ENABLED=0 go build .

FROM alpine as main
COPY --from=build /app /app

ARG PASSWORD
ENV PASSWORD $PASSWORD
WORKDIR /app

ENTRYPOINT ["/app/ezcat"]
