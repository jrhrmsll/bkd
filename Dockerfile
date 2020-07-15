FROM       golang:alpine as build
RUN        apk add --no-cache ca-certificates
WORKDIR    /app
COPY       . .
RUN        go build -o bkd ./cmd/bkd

FROM       alpine:latest
RUN        apk add --no-cache ca-certificates
COPY       --from=build /app/bkd /usr/local/bin/bkd
EXPOSE     8080
ENTRYPOINT [ "bkd" ]
