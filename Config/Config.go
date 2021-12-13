package Config

import "koinfolio/Utils"

var (
	CoinMarketCapApiUrl = Utils.GetEnv(
		"COIN_MARKET_CAP_API_URL", "https://pro-api.coinmarketcap.com/v1/tools/price-conversion")
	CoinMarketCapApiKey = Utils.GetEnv(
		"COIN_MARKET_CAP_API_KEY", "30593477-b629-4f2b-bbb5-6a95d8e88211")
)
