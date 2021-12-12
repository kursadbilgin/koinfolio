package Models

import (
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

var CoinPortfolioCollection *mongo.Collection
var HistoryCollection *mongo.Collection

type CoinHistory struct {
	ID      string    `json:"id"`
	Code    string    `json:"code"`
	History []History `json:"history"`
}

type History struct {
	Amount int   `json:"amount"`
	Price  Price `json:"price"`
}

type Price struct {
	Old     string `json:"old"`
	Current string `json:"current"`
}

type AddCoinRequest struct {
	Amount   int    `json:"amount"`
	CoinCode string `json:"coin_code"`
}

type DbCoinRecord struct {
	ID       string `bson:"id" json:"id"`
	Amount   int    `bson:"amount" json:"amount"`
	CoinCode string `bson:"coin_code" json:"coin_code"`
	Price    string `bson:"price" json:"price"`
}

type APIResponse struct {
	Status Status `json:"status"`
	Data   Data   `json:"data"`
}

type Quote struct {
	USD USD `json:"USD"`
}

type USD struct {
	Price       float64   `json:"price"`
	LastUpdated time.Time `json:"last_updated"`
}

type Data struct {
	ID          int       `json:"id"`
	Symbol      string    `json:"symbol"`
	Name        string    `json:"name"`
	Amount      int       `json:"amount"`
	LastUpdated time.Time `json:"last_updated"`
	Quote       Quote     `json:"quote"`
}

type Status struct {
	Timestamp    time.Time   `json:"timestamp"`
	ErrorCode    int         `json:"error_code"`
	ErrorMessage interface{} `json:"error_message"`
	Elapsed      int         `json:"elapsed"`
	CreditCount  int         `json:"credit_count"`
	Notice       interface{} `json:"notice"`
}

type Response struct {
	ID     string `json:"id"`
	Code   string `json:"code"`
	Amount int    `json:"amount"`
	Price  string `json:"price"`
}
