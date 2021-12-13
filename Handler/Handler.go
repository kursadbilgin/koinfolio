package Handler

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"koinfolio/Logger"
	"koinfolio/Models"
	"koinfolio/Utils"
	"net/http"
)

func AddCoin(c *gin.Context) {
	var request Models.AddCoinRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		return
	}
	resp, errApi := Utils.CoinMarketCapAPI(request.Amount, request.CoinCode)
	if errApi != nil {
		Logger.Error.Println(errApi)
		c.JSON(http.StatusBadRequest, gin.H{"Error[1]": "Bad request"})
		return
	}
	if !Utils.ValidateResponse(resp.Status) {
		Logger.Error.Println(resp.Status.ErrorMessage)
		c.JSON(http.StatusBadRequest, gin.H{"Error[2]:": "Bad request"})
		return
	}
	dbRecord := Models.DbCoinRecord{
		ID:       uuid.New().String()[:8],
		Amount:   resp.Data.Amount,
		CoinCode: resp.Data.Symbol,
		Price:    fmt.Sprintf("%f", resp.Data.Quote.USD.Price),
	}
	ctx := context.Background()

	searchResult, err := Models.HistoryCollection.Find(ctx, bson.M{"coin_code": dbRecord.CoinCode})
	if err != nil {
		Logger.Error.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"Error[3]": "server error"})
		return
	}
	if searchResult.Current.String() != "" {
		c.JSON(http.StatusForbidden, gin.H{"Error": "Currency already exists"})
		return
	}

	_, err = Models.CoinPortfolioCollection.InsertOne(ctx, dbRecord)
	if err != nil {
		Logger.Error.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"Error[4]": "server error"})
		return
	}

	err = Utils.HistoryCreate(&ctx, dbRecord.ID, resp.Data.Symbol, &dbRecord, nil)
	if err != nil {
		Logger.Error.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "server error"})
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
		Logger.Error.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "server error"})
		return
	}

	defer searchResult.Close(ctx)
	if err = searchResult.All(ctx, &coins); err != nil {
		Logger.Error.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "server error"})
		return
	}
	if len(coins) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"Error": "Currency with that id does not exist"})
		return
	}
	for idx, _ := range coins[0].History {
		respOld, err := Utils.CoinMarketCapAPI(coins[0].History[idx].Amount, coins[0].Code)
		if err != nil {
			Logger.Error.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "server error"})
			return
		}
		coins[0].History[idx].Price.Current = fmt.Sprintf("%f", respOld.Data.Quote.USD.Price)
	}

	c.JSON(http.StatusOK, coins)
}

func GetCoins(c *gin.Context) {
	ctx := context.Background()
	var coins []Models.CoinHistory

	cursor, err := Models.HistoryCollection.Find(context.TODO(), bson.M{})
	defer cursor.Close(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "server error"})
		return
	}

	if err = cursor.All(ctx, &coins); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "server error"})
		return
	}
	if 0 == len(coins) {
		c.JSON(http.StatusNotFound, gin.H{"Error": "currency record not found"})
		return
	}

	for id, _ := range coins {
		for idx, _ := range coins[id].History {
			respOld, err := Utils.CoinMarketCapAPI(coins[id].History[idx].Amount, coins[id].Code)
			if err != nil {
				Logger.Error.Println(err)
				c.JSON(http.StatusInternalServerError, gin.H{"Error": "server error"})
				return
			}
			coins[id].History[idx].Price.Current = fmt.Sprintf("%f", respOld.Data.Quote.USD.Price)
		}
	}

	c.JSON(http.StatusOK, coins)
}

func DeleteCoin(c *gin.Context) {
	id := c.Param("id")

	ctx := context.Background()
	var deletedCoin Models.DbCoinRecord
	err := Models.CoinPortfolioCollection.FindOneAndDelete(ctx, bson.M{"id": id}).Decode(&deletedCoin)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			Logger.Error.Println(err)
			c.JSON(http.StatusNotFound, gin.H{"Error": "Currency with that id does not exist"})
			return
		}
		Logger.Error.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"ERROR": "server error"})
		return
	}

	var deletedCoinHistory Models.CoinHistory
	err = Models.HistoryCollection.FindOneAndDelete(ctx, bson.M{"id": id}).Decode(&deletedCoinHistory)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			Logger.Error.Println(err)
			c.JSON(http.StatusNotFound, gin.H{"Error": "Currency with that id does not exist"})
			return
		}
		Logger.Error.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"ERROR": "server error"})
		return
	}

	c.JSON(http.StatusOK, nil)
}

func EditCoin(c *gin.Context) {
	var request Models.AddCoinRequest
	id := c.Param("id")

	if err := c.ShouldBindJSON(&request); err != nil {
		Logger.Error.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"ERROR": "Bad request"})
		return
	}

	byteUser, err := bson.Marshal(request)
	if err != nil {
		Logger.Error.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"ERROR": "server error"})
		return
	}

	var bCoin bson.M
	if err = bson.Unmarshal(byteUser, &bCoin); err != nil {
		Logger.Error.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"ERROR": "server error"})
		return
	}
	resp, err := Utils.CoinMarketCapAPI(request.Amount, request.CoinCode)
	if err != nil {
		Logger.Error.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"Error": "bad request"})
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
	err = Models.CoinPortfolioCollection.FindOne(ctx, bson.M{"id": id}).Decode(&getCoin)
	if err != nil {
		Logger.Error.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "server error"})
		return
	}

	opts := options.Update().SetUpsert(false)
	_, err = Models.CoinPortfolioCollection.UpdateOne(
		ctx, bson.M{"id": id}, bson.D{{Key: "$set", Value: newData}}, opts)
	if err == mongo.ErrNoDocuments {
		Logger.Error.Println(err)
		c.JSON(http.StatusNotFound, gin.H{"Error": "Currency with that id does not exist"})
		return
	}

	err = Utils.HistoryCreate(&ctx, id, resp.Data.Symbol, &getCoin, &newData)
	if err != nil {
		Logger.Error.Println(err)
		c.JSON(http.StatusNotFound, gin.H{"ERROR": "Currency with that id does not exist"})
		return
	}

	c.JSON(http.StatusOK, newData)
}
