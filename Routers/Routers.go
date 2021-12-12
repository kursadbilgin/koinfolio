package Routers

import (
	"github.com/gin-gonic/gin"
	"koinfolio/Handler"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.POST("currency", Handler.AddCoin)
	r.GET("currencies/:id", Handler.GetCoinByID)
	r.GET("currencies", Handler.GetCoins)
	r.PATCH("currencies/:id", Handler.EditCoin)
	r.DELETE("currencies/:id", Handler.DeleteCoin)

	return r
}
