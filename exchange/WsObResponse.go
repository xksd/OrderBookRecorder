package exchange

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gorilla/websocket"
)

// Websocket responds with Orderbook Snapshots
// tick by tick emmited by Binance

// WsSymbolObResponse - stores data for "combined streams" from Binance WS
// If we would pick only one symbol, ws response would be a single WsObResponseData,
// but payload for "combined streams" has a bit different format,
// that is WsSymbolObResponse.

// Combined stream for multiple symbols payload
type WsObResponse struct {
	Symbol string `json:"stream"`
	Data   WsObResponseData
}

// Individual Symbol payload for Orderbook snapshot
type WsObResponseData struct {
	Symbol               string     `json:"s"`
	EventType            string     `json:"e"`
	Timestamp            int64      `json:"E"`
	TimestampTransaction int64      `json:"T"`
	IdFirst              int        `json:"U"`
	IdFinal              int        `json:"u"`
	IdFinalLastStream    int        `json:"pu"`
	Bids                 [][]string `json:"b"`
	Asks                 [][]string `json:"a"`
}

// URL for Websocket request
// wss://fstream.binance.com/stream?streams=btcusdt@depth/ethusdt@depth/
const wsUrlBase = "wss://fstream.binance.com"
const wsUrlDepth = "/stream?streams="

func CreateWsUrl(symbolsList []string) (url string) {
	url = wsUrlBase + wsUrlDepth
	for i, symbol := range symbolsList {
		symbolLower := strings.ToLower(symbol)
		if i == 0 {
			url = url + symbolLower + "@depth"
		} else if i > 0 {
			url = url + "/" + symbolLower + "@depth"
		}
	}
	return
}

// Start WebSocket stream to get updates in a loop,
// each update sent to channel.
func GetUpdate(wsCh chan<- WsObResponse, wsChDone <-chan bool, wsUrl string) {
	// Establish connection and deal with different outcomes from the connection
	c, resp, err := websocket.DefaultDialer.Dial(wsUrl, nil)
	if err != nil {
		fmt.Printf("handshake failed with status %d\n", resp.StatusCode)
		fmt.Println("dial:", err)
	}
	// When the program closes, close the connection
	defer c.Close()

	// Read message
	// Each message is a tick update for a single symbol
	go func() {
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				fmt.Println("read:", err)
				return
			}

			var res WsObResponse
			json.Unmarshal(message, &res)

			res.Symbol = strings.ToUpper(strings.Replace(res.Symbol, "@depth", "", 1))
			// fmt.Println("ws tick")
			wsCh <- res
		}
	}()

	// Exit after signal from channel wsChDone
	i := <-wsChDone
	fmt.Println(">>> GetUpdate() Interrupted", i)
	{
		err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		if err != nil {
			fmt.Println("Error on Write close:", err)
			return
		}
	}

}
