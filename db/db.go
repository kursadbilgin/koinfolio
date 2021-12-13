package db

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"koinfolio/Models"
	"log"
	"time"
)

func InitMongoDB() {
	newClient, err := mongo.NewClient(options.Client().ApplyURI("mongodb://mongo:27017"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = newClient.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = newClient.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	quickstartDatabase := newClient.Database("portfolio")
	Models.CoinPortfolioCollection = quickstartDatabase.Collection("coin_portfolio")
	Models.HistoryCollection = quickstartDatabase.Collection("history")
}
