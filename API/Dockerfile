# syntax=docker/dockerfile:1

## Build
FROM golang:1.19-buster AS build

WORKDIR /app

COPY ./go.mod ./
COPY ./go.sum ./
COPY ./Library ./Library
COPY ./API ./API

RUN go mod download

WORKDIR /app/API/cmd/server

RUN go build -o /api

## Deploy
FROM gcr.io/distroless/base-debian11:debug

WORKDIR /

COPY --from=build /api /api
COPY ./.env ./

EXPOSE 8080

USER root:root

ENTRYPOINT ["/api"]
