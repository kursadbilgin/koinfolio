package Utils

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"koinfolio/Models"
)

func HistoryCreate(ctx *context.Context, id, symbol string, dbRecord, newRecord *Models.DbCoinRecord) (err error) {
	var historyNew Models.CoinHistory

	if newRecord == nil {
		historyNew = Models.CoinHistory{
			ID:   id,
			Code: symbol,
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
		_, err = Models.HistoryCollection.InsertOne(*ctx, historyNew)
		if err != nil {
			return err
		}

		return nil
	} else {
		var history Models.History
		err = Models.HistoryCollection.FindOne(*ctx, bson.M{"id": id}).Decode(&history)
		if err != nil {
			return err
		}

		historyNew = Models.CoinHistory{
			ID:   id,
			Code: symbol,
			History: []Models.History{
				{
					Amount: dbRecord.Amount,
					Price: Models.Price{
						Old:     dbRecord.Price,
						Current: "",
					},
				},
				{
					Amount: newRecord.Amount,
					Price: Models.Price{
						Old:     newRecord.Price,
						Current: "",
					},
				},
			},
		}
		opts := options.Update().SetUpsert(false)
		_, err = Models.HistoryCollection.UpdateOne(*ctx, bson.M{"id": id}, bson.D{{Key: "$set", Value: historyNew}}, opts)
		if err == mongo.ErrNoDocuments {
			return err
		}

		return nil
	}

	return nil
}
