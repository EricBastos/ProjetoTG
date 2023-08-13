# syntax=docker/dockerfile:1

## Build
FROM node:18 AS build

WORKDIR /usr/src/app

COPY ./Frontend/package*.json ./

RUN npm install

COPY ./Frontend ./

RUN npm run build --omit=dev

## Deploy
FROM node:18-slim

WORKDIR /

COPY --from=build /usr/src/app/dist ./
RUN npm install --global http-server

EXPOSE 8080

USER root:root

CMD http-server
