FROM golang:1.19 AS build
WORKDIR /
COPY . .

RUN go build -o /app /main.go

FROM ubuntu:20.04

ENV GIN_MODE release

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /app /app
ENTRYPOINT [ "/app" ]

