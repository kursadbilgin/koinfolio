package Handler

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"koinfolio/Models"
	"koinfolio/Utils"
	"net/http"
)

func AddCoin(c *gin.Context) {
	var request Models.AddCoinRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		return
	}
	resp, errApi := CoinMarketCapAPI(request.Amount, request.CoinCode)
	if errApi != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": errApi.Error()})
		return
	}
	if !Utils.ValidateResponse(resp.Status) {
		c.JSON(http.StatusBadRequest, gin.H{"Error:": "Not valid", "Message": resp.Status.ErrorMessage})
		return
	}
	dbRecord := Models.DbCoinRecord{
		ID:       uuid.New().String()[:8],
		Amount:   resp.Data.Amount,
		CoinCode: resp.Data.Symbol,
		Price:    fmt.Sprintf("%f", resp.Data.Quote.USD.Price),
	}

	ctx := context.Background()
	_, err := Models.CoinPortfolioCollection.InsertOne(ctx, dbRecord)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ERROR": err.Error()})
		return
	}

	historyNew := Models.CoinHistory{
		ID:   dbRecord.ID,
		Code: resp.Data.Symbol,
		History: []Models.History{
			{
				Amount: dbRecord.Amount,
				Price: Models.Price{
					Old:     dbRecord.Price,
					Current: dbRecord.Price,
				},
			},
			{
				Amount: dbRecord.Amount,
				Price: Models.Price{
					Old:     dbRecord.Price,
					Current: dbRecord.Price,
				},
			},
		},
	}
	_, err = Models.HistoryCollection.InsertOne(ctx, historyNew)
	if err != nil {
		fmt.Println("asd")
		c.JSON(http.StatusInternalServerError, gin.H{"ERROR": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dbRecord)
}

func GetCoinByID(c *gin.Context) {
	id := c.Param("id")

	ctx := context.Background()
	var coins []Models.CoinHistory
	searchResult, err := Models.HistoryCollection.Find(ctx, bson.M{"id": id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ERROR": err.Error()})
		return
	}

	defer searchResult.Close(ctx)
	if err = searchResult.All(ctx, &coins); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ERROR": err.Error()})
		return
	}

	for idx, _ := range coins[0].History {
		respOld, _ := CoinMarketCapAPI(coins[0].History[idx].Amount, coins[0].Code)
		coins[0].History[idx].Price.Current = fmt.Sprintf("%f", respOld.Data.Quote.USD.Price)
	}

	c.JSON(http.StatusOK, coins)
}

func GetCoins(c *gin.Context) {

	ctx := context.Background()
	var coins []Models.DbCoinRecord

	cursor, err := Models.CoinPortfolioCollection.Find(context.TODO(), bson.M{})
	defer cursor.Close(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ERROR": err.Error()})
		return
	}

	if err = cursor.All(ctx, &coins); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ERROR": err.Error()})
		return
	}

	c.JSON(http.StatusOK, coins)
}

func DeleteCoin(c *gin.Context) {
	id := c.Param("id")

	ctx := context.Background()
	var deletedCoin *Models.DbCoinRecord
	err := Models.CoinPortfolioCollection.FindOneAndDelete(ctx, bson.M{"id": id}).Decode(deletedCoin)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotModified, gin.H{"ERROR": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"ERROR": err.Error()})
		return
	}

	c.JSON(http.StatusOK, nil)
}

func EditCoin(c *gin.Context) {
	var request Models.AddCoinRequest
	id := c.Param("id")

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"ERROR": err.Error()})
		return
	}

	byteUser, err := bson.Marshal(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ERROR": err.Error()})
		return
	}

	var bCoin bson.M
	if err = bson.Unmarshal(byteUser, &bCoin); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ERROR": err.Error()})
		return
	}
	resp, err := CoinMarketCapAPI(request.Amount, request.CoinCode)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}
	newData := Models.DbCoinRecord{
		ID:       id,
		Amount:   request.Amount,
		CoinCode: resp.Data.Symbol,
		Price:    fmt.Sprintf("%f", resp.Data.Quote.USD.Price),
	}
	ctx := context.Background()
	var getCoin Models.DbCoinRecord
	filter := bson.M{"id": id}
	err = Models.CoinPortfolioCollection.FindOne(ctx, filter).Decode(&getCoin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ERROR": err.Error()})
		return
	}

	opts := options.Update().SetUpsert(false)
	update := bson.D{{Key: "$set", Value: newData}}
	_, err = Models.CoinPortfolioCollection.UpdateOne(ctx, filter, update, opts)
	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusNotModified, gin.H{"ERROR": err.Error()})
		return
	}
	var history Models.CoinHistory

	err = Models.HistoryCollection.FindOne(ctx, bson.M{"id": id}).Decode(&history)
	if err != nil {
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"ERROR": err.Error()})
			return
		}
	}

	historyNew := Models.CoinHistory{
		ID:   id,
		Code: resp.Data.Symbol,
		History: []Models.History{
			{
				Amount: getCoin.Amount,
				Price: Models.Price{
					Old:     getCoin.Price,
					Current: "",
				},
			},
			{
				Amount: newData.Amount,
				Price: Models.Price{
					Old:     newData.Price,
					Current: "",
				},
			},
		},
	}

	updateHistory := bson.D{{Key: "$set", Value: historyNew}}
	_, err = Models.HistoryCollection.UpdateOne(ctx, filter, updateHistory, opts)
	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusNotModified, gin.H{"ERROR": err.Error()})
		return
	}

	c.JSON(http.StatusOK, newData)
}
