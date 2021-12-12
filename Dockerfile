# syntax=docker/dockerfile:1

FROM golang:1.17.5

COPY . /app
WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download
