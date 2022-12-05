package exchange_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/xksd/OrderBookRecorder/exchange"
)

func TestGetObResponse(t *testing.T) {
	symbol := "BTCUSDT"
	r := exchange.GetObResponse(symbol, 0)

	timestampMin := time.Now().UnixMilli() - int64(time.Second)

	fmt.Println(r.Bids)

	if r.Symbol != symbol {
		t.Error("Incorrect: Wrong symbol returned")
	}
	if r.Timestamp < timestampMin {
		t.Error("Incorrect: Timestamp is wrong, maybe delay is too big.")
	}
	if len(r.Bids) < 10 {
		t.Error("Incorrect: OrderBook Bids are almost empty!")
	}
	if len(r.Asks) < 10 {
		t.Error("Incorrect: OrderBook Asks are almost empty!")
	}

}
