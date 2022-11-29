package exchange

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"runtime"
	"time"
)

type ObSnapshot struct {
	Timestamp int64      `json:"E"`
	Bids      [][]string `json:"bids"`
	Asks      [][]string `json:"asks"`
}

func GetObSnapshot(symbol string, limitOfObLevels int) ObSnapshot {
	fmt.Println("Started ObSnaphost for   --- ", symbol)
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

	var resUnmarshalled ObSnapshot
	json.Unmarshal(resData, &resUnmarshalled)

	fmt.Println("Returned ObSnapshot for  --- ", symbol)
	// fmt.Println(resUnmarshalled.Bids[0])

	return resUnmarshalled

}

func (ob *ObSnapshot) Print() {
	fmt.Println()
	fmt.Println("Timestamp         ", time.UnixMilli(ob.Timestamp))
	fmt.Println("Timestamp (unix)  ", ob.Timestamp)

	fmt.Println("\nBids\n", len(ob.Bids))
	fmt.Println(ob.Bids)
	fmt.Println("\nAsks\n", len(ob.Asks))
	fmt.Println(ob.Asks)
}
