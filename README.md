# Koinfolio

A rest-api that works with golang as coin portfolio

# Technical Details
- Golang/Gin is used as application framework
- MongoDB is database

## Usage

Fist of all, clone the repo with the command below. You must have golang installed on your computer to run the project.

```shell
git clone https://github.com/kursadbilgin/koinfolio
```

+ ### With Docker

If you have docker installed on your system, you can run the project with docker. Run command below on terminal.

````shell
docker-compose -f docker-compose.yml build  
docker-compose up -d
````

Now go to `localhost:8090` in your browser. Here you can send your requests to api.

or you can test the api on the URL below.

## Endpoint Table

|       Endpoints        |      Descriptions       | Methods |                                                                                  cURL example                                                                                   |
|:----------------------:|:-----------------------:|:-------:|:-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------:|
|       /currency        |   Add Coin Endpoint.    |  POST   |     curl -X POST http://localhost:8090/currency -H 'cache-control: no-cache' -H 'content-type: application/json' -d '{ "amount": <int_amount>, "coin_code": <coin_code>"}''     |
|    /currencies/:id     |  Get Coin Record By ID  |   GET   |                                                curl -X GET http://localhost:8090/currencies/<db_id> -H 'cache-control: no-cache'                                                |
|      /currencies       |  Get All Coin Record.   |   GET   |                                                   curl -X GET http://localhost:8090/currencies/ -H 'cache-control: no-cache'                                                    |
|    /currencies/:id     |  Edit Currency Record.  |  PATCH  | curl -X PATCH http://localhost:8090/currencies/<db_id> -H 'cache-control: no-cache' -H 'content-type: application/json' -d '{"amount": <int_amount>, "coin_code": <coin_code>}' |
|    /currencies/:id     | Delete Currency Record. | DELETE  |                                              curl -X DELETE http://localhost:8090/currencies/<db_id> -H 'cache-control: no-cache'                                               |


