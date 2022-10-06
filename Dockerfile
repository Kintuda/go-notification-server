FROM golang:1.17.7-alpine AS build_base

COPY . /myapp
WORKDIR /myapp
ENV GO111MODULE=on
RUN CGO_ENABLED=0 GOOS=linux go build -o myapp

FROM alpine:latest

RUN apk add --no-cache tzdata
WORKDIR /root/

COPY --from=build_base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build_base /myapp .

EXPOSE 8080

CMD ["./myapp", "server"]