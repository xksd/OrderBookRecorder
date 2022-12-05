package exchange

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"runtime"
)

// Getting first Orderbook snapshot
// using http-response to Binance API

// Struct to unmarshall response from Binance
type ObResponse struct {
	Symbol    string
	Timestamp int64      `json:"E"`
	Bids      [][]string `json:"bids"`
	Asks      [][]string `json:"asks"`
}

// Request Orderbook snapshot from Binance
// params: symbol string, limitOfObLevels (0 = no limits)
func GetObResponse(symbol string, limitOfObLevels int) ObResponse {
	// Define URL for Get request
	var url string = binanceUrls["api_base"] + binanceUrls["order_book"] +
		"?symbol=" + symbol
	if limitOfObLevels > 1 {
		url = url + "&limit=" + fmt.Sprintf("%v", limitOfObLevels)
	}

	// Make a Get request
	res, err := http.Get(url)
	runtime.Gosched()
	if err != nil {
		fmt.Println("Error in GET > exchange.GetOrederBook()")
		log.Fatal(err)
	}
	defer res.Body.Close()

	// Unmarshalling
	resData, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var obResp ObResponse
	json.Unmarshal(resData, &obResp)

	obResp.Symbol = symbol

	return obResp
}
