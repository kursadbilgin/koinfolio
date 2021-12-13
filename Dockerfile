# syntax=docker/dockerfile:1

FROM golang:1.17.5

ENV COIN_MARKET_CAP_API_URL="https://pro-api.coinmarketcap.com/v1/tools/price-conversion"
ENV COIN_MARKET_CAP_API_KEY="30593477-b629-4f2b-bbb5-6a95d8e88211"

COPY . /app
WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download
RUN go build -o /koinfolio

CMD [ "/koinfolio" ]
