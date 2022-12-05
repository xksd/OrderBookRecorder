package exchange

var binanceUrls = map[string]string{
	"api_base":        "https://fapi.binance.com/fapi/v1",
	"account_status":  "/sapi/v1/account/status",
	"get_server_time": "/time",

	// FUTURES Market Data
	"symbol_price_ticker": "/ticker/price",
	"24hr_ticker_stats":   "/ticker/24hr",
	"candles":             "/klines",
	"order_book":          "/depth",
}
