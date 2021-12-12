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
	resp := CoinMarketCapAPI(request.Amount, request.CoinCode)
	if !Utils.ValidateResponse(resp.Status) {
		c.JSON(http.StatusBadRequest, resp.Status)
		return
	}
	dbRecord := Models.DbCoinRecord{
		ID:       uuid.New().String()[:8],
		Amount:   resp.Data.Amount,
		CoinCode: resp.Data.Symbol,
		Price:    fmt.Sprintf("%f", resp.Data.Quote.USD.Price),
	}

	ctx := context.Background()
	_, err := Models.Collection.InsertOne(ctx, dbRecord)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ERROR": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dbRecord)
}

func GetCoinByID(c *gin.Context) {
	id := c.Param("id")

	ctx := context.Background()
	var coins []Models.DbCoinRecord
	searchResult, err := Models.Collection.Find(ctx, bson.M{"id": id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ERROR": err.Error()})
		return
	}

	defer searchResult.Close(ctx)

	if err = searchResult.All(ctx, &coins); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ERROR": err.Error()})
		return
	}

	c.JSON(http.StatusOK, coins)
}

func GetCoins(c *gin.Context) {

	ctx := context.Background()
	var coins []Models.DbCoinRecord

	cursor, err := Models.Collection.Find(context.TODO(), bson.M{})
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
	var deletedCoin Models.DbCoinRecord
	err := Models.Collection.FindOneAndDelete(ctx, bson.M{"id": id}).Decode(&deletedCoin)
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
	resp := CoinMarketCapAPI(request.Amount, request.CoinCode)
	newData := Models.DbCoinRecord{
		ID:       id,
		Amount:   request.Amount,
		CoinCode: resp.Data.Symbol,
		Price:    fmt.Sprintf("%f", resp.Data.Quote.USD.Price),
	}
	ctx := context.Background()
	var updatedCoin Models.DbCoinRecord
	opts := options.FindOneAndUpdate()
	filter := bson.M{"id": id}
	update := bson.D{{Key: "$set", Value: newData}}
	err = Models.Collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updatedCoin)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotModified, gin.H{"ERROR": err.Error()})
			return

		}
		c.JSON(http.StatusInternalServerError, gin.H{"ERROR": err.Error()})
		return
	}

	c.JSON(http.StatusOK, newData)
}
