# syntax=docker/dockerfile:1

## Build
FROM golang:1.19-buster AS build

WORKDIR /app

COPY ./go.mod ./
COPY ./go.sum ./
COPY ./Library ./Library
COPY ./BankWebhook ./BankWebhook

RUN go mod download

WORKDIR /app/BankWebhook/cmd/server

RUN go build -o /bankWebhook

## Deploy
FROM gcr.io/distroless/base-debian11:debug

WORKDIR /

COPY --from=build /bankWebhook /bankWebhook
COPY ./.env ./

EXPOSE 8080

USER root:root

ENTRYPOINT ["/bankWebhook"]
