FROM golang:1.20-alpine AS build

ARG GOOS=linux
ARG GOARCH=amd64
ARG HTTP_PORT=8080
ARG APP_PKG_NAME=item-service

ENV CGO_ENABLED=0 \
    GOOS=$GOOS \
    GOARCH=$GOARCH \
    GO111MODULE=on \
    APP_PKG_NAME=$APP_PKG_NAME

RUN apk update && apk add --no-cache bash

WORKDIR /go/src/$APP_PKG_NAME
COPY . .

RUN go mod tidy
RUN go test ./...
RUN go build -o /app -ldflags "-s -w -extldflags '-static'" ./cmd/*.go

# Copy to Alpine image
FROM alpine:3.12
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=build /app /app

EXPOSE $HTTP_PORT
CMD ["./app"]
