package data

import (
	"fmt"
	"testing"

	"github.com/xksd/OrderBookRecorder/exchange"
)

func TestAppend(t *testing.T) {
	// Init
	obBids1 := BidsAsks{
		BidAskPriceLevel{16905.10, 2.785},
		BidAskPriceLevel{16904.10, 3.214},
		BidAskPriceLevel{16903.10, 6.123},
		BidAskPriceLevel{16902.10, 0.123},
		BidAskPriceLevel{16901.10, 2.123},
	}
	obAsks1 := BidsAsks{
		BidAskPriceLevel{16905.90, 1.112},
		BidAskPriceLevel{16906.40, 0.100},
		BidAskPriceLevel{16907.70, 4.311},
		BidAskPriceLevel{16908.80, 1.110},
		BidAskPriceLevel{16909.25, 12.020},
	}
	ob1 := SymbolObSnapshot{
		Symbol:           "BTCUSDT",
		Timestamp:        1670155864616,
		DecimalsPrice:    2,
		DecimalsQuantity: 3,
		Bids:             obBids1,
		Asks:             obAsks1,
	}
	obBids2 := BidsAsks{
		BidAskPriceLevel{905.510, 2.785},
		BidAskPriceLevel{904.240, 3.214},
		BidAskPriceLevel{903.110, 6.123},
		BidAskPriceLevel{902.001, 0.123},
		BidAskPriceLevel{901.000, 2.123},
	}
	obAsks2 := BidsAsks{
		BidAskPriceLevel{905.901, 1.112},
		BidAskPriceLevel{906.400, 0.100},
		BidAskPriceLevel{907.970, 4.311},
		BidAskPriceLevel{908.280, 1.110},
		BidAskPriceLevel{909.225, 12.020},
	}
	ob2 := SymbolObSnapshot{
		Symbol:           "ETHUSDT",
		Timestamp:        1670155864616,
		DecimalsPrice:    2,
		DecimalsQuantity: 3,
		Bids:             obBids2,
		Asks:             obAsks2,
	}

	// Testing
	var allObBuffer AllObBuffer

	allObBuffer.Append(&ob1)
	allObBuffer.Append(&ob2)

	if len(allObBuffer) != 2 {
		t.Error("AllObBuffer slice length is incorrect, some symbols were not added!")
	}
	for _, ob := range allObBuffer {
		if ob.Symbol != "BTCUSDT" && ob.Symbol != "ETHUSDT" {
			t.Error("Symbol is missing: " + ob.Symbol)
		}
		if len(ob.Asks) != 5 {
			t.Error("Asks length is incorrect!")
		}
		if len(ob.Bids) != 5 {
			t.Error("Bids length is incorrect!")
		}
	}

}

func TestUpdateFromWsTick(t *testing.T) {
	// Init
	obBids1 := BidsAsks{
		BidAskPriceLevel{16905.10, 2.785},
		BidAskPriceLevel{16904.10, 3.214},
		BidAskPriceLevel{16903.10, 6.123},
		BidAskPriceLevel{16902.10, 0.123},
		BidAskPriceLevel{16901.10, 2.123},
	}
	obAsks1 := BidsAsks{
		BidAskPriceLevel{16905.90, 1.112},
		BidAskPriceLevel{16906.40, 0.100},
		BidAskPriceLevel{16907.70, 4.311},
		BidAskPriceLevel{16908.80, 1.110},
		BidAskPriceLevel{16909.25, 12.020},
	}
	ob1 := SymbolObSnapshot{
		Symbol:           "BTCUSDT",
		Timestamp:        1670155864616,
		DecimalsPrice:    2,
		DecimalsQuantity: 3,
		Bids:             obBids1,
		Asks:             obAsks1,
	}
	obBids2 := BidsAsks{
		BidAskPriceLevel{905.510, 2.785},
		BidAskPriceLevel{904.240, 3.214},
		BidAskPriceLevel{903.110, 6.123},
		BidAskPriceLevel{902.001, 0.123},
		BidAskPriceLevel{901.000, 2.123},
	}
	obAsks2 := BidsAsks{
		BidAskPriceLevel{905.901, 1.112},
		BidAskPriceLevel{906.400, 0.100},
		BidAskPriceLevel{907.970, 4.311},
		BidAskPriceLevel{908.280, 1.110},
		BidAskPriceLevel{909.225, 12.020},
	}
	ob2 := SymbolObSnapshot{
		Symbol:           "ETHUSDT",
		Timestamp:        1670155864616,
		DecimalsPrice:    2,
		DecimalsQuantity: 3,
		Bids:             obBids2,
		Asks:             obAsks2,
	}
	// Testing
	var allObBuffer AllObBuffer

	allObBuffer.Append(&ob1)
	allObBuffer.Append(&ob2)

	for _, ob := range allObBuffer {
		if ob.Symbol == "BTCUSDT" {
			// fmt.Println(ob.Bids)
			fmt.Println(ob.Asks)
		}
	}

	// Initialize a variable containing a fake update from WebSocket
	wsTick := exchange.WsObResponse{
		Symbol: "BTCUSDT",
		Data: exchange.WsObResponseData{
			Symbol:               "BTCUSDT",
			EventType:            "depthUpdate",
			Timestamp:            1670155864616,
			TimestampTransaction: 1670155864616,
			IdFirst:              1231,
			IdFinal:              1232,
			IdFinalLastStream:    1233,
			Bids: [][]string{
				{"16908.10", "99.010"},
				{"16905.10", "99.100"},
				{"16901.10", "99.001"},
			},
			Asks: [][]string{
				{"16902.45", "99.001"},
				{"16906.40", "99.100"},
				{"16907.10", "99.010"},
				{"16911.20", "99.000"},
			},
		},
	}

	// Testing
	if err := allObBuffer.UpdateFromWsTick(wsTick); err != nil {
		t.Error("Error while running UpdateFromWsTick()")
		t.Error(err)
	}

	for _, ob := range allObBuffer {
		if ob.Symbol == "BTCUSDT" {
			// fmt.Println(ob.Bids)
			fmt.Println(ob.Asks)
		}
	}

	// Check
	for _, ob := range allObBuffer {
		if ob.Symbol != "BTCUSDT" && ob.Symbol != "ETHUSDT" {
			t.Error("Symbol is missing: " + ob.Symbol)
		}
		if ob.Symbol == "BTCUSDT" {
			if len(ob.Asks) != 8 {
				t.Error("Asks length is incorrect!")
			}
			if len(ob.Bids) != 6 {
				t.Error("Bids length is incorrect!")
			}
			if ob.Bids[0].Price != 16908.1 && ob.Bids[0].Quantity != 99.01 {
				t.Error("Bids 0 element is incorrect.")
				t.Error("Desired: ", 16908.1, 99.01)
				t.Error("Result: ", ob.Bids[0].Price, ob.Bids[0].Quantity)
			}
			if ob.Bids[1].Price != 16905.1 && ob.Bids[0].Quantity != 99.1 {
				t.Error("Bids 0 element is incorrect.")
				t.Error("Desired: ", 16905.1, 99.1)
				t.Error("Result: ", ob.Bids[1].Price, ob.Bids[1].Quantity)
			}
			if ob.Bids[5].Price != 16901.1 && ob.Bids[0].Quantity != 55.001 {
				t.Error("Bids 0 element is incorrect.")
				t.Error("Desired: ", 16905.1, 99.1)
				t.Error("Result: ", ob.Bids[1].Price, ob.Bids[1].Quantity)
			}
			if ob.Asks[0].Price != 16902.45 && ob.Asks[0].Quantity != 99.000 {
				t.Error("Asks 0 element is incorrect.")
				t.Error("Desired: ", 16902.45, 99.000)
				t.Error("Result: ", ob.Asks[0].Price, ob.Asks[0].Quantity)
			}
			if ob.Asks[3].Price != 16907.10 && ob.Asks[3].Quantity != 99.000 {
				t.Error("Asks 0 element is incorrect.")
				t.Error("Desired: ", 16907.10, 99.000)
				t.Error("Result: ", ob.Asks[3].Price, ob.Asks[3].Quantity)
			}
			if ob.Asks[7].Price != 16911.20 && ob.Asks[7].Quantity != 99.000 {
				t.Error("Asks 0 element is incorrect.")
				t.Error("Desired: ", 16911.20, 99.000)
				t.Error("Result: ", ob.Asks[7].Price, ob.Asks[7].Quantity)
			}
		}
	}
}
