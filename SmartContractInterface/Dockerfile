# syntax=docker/dockerfile:1

## Build
FROM golang:1.19-buster AS build

WORKDIR /app

COPY ./go.mod ./
COPY ./go.sum ./
COPY ./Library ./Library
COPY ./SmartContractInterface ./SmartContractInterface

RUN go mod download

WORKDIR /app/SmartContractInterface/cmd/server

RUN go build -o /smartContractInterface

## Deploy
FROM gcr.io/distroless/base-debian11:debug

WORKDIR /

COPY --from=build /smartContractInterface /smartContractInterface
COPY ./.env ./

EXPOSE 50051

USER root:root

ENTRYPOINT ["/smartContractInterface"]
