# syntax=docker/dockerfile:1

FROM golang:1.18-alpine

# go cn proxy
ENV GOPROXY "https://goproxy.cn,direct"

WORKDIR /app

COPY . .
RUN go mod download

RUN go build cmd/serve/main.go

EXPOSE 10000
EXPOSE 10001
EXPOSE 10002

CMD [ "/app/main" ]
