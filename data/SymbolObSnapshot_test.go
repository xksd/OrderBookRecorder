package data

import (
	"testing"
	"time"

	"github.com/xksd/OrderBookRecorder/exchange"
)

func Test_CreateSymbolObSnapshot(t *testing.T) {
	// Init
	obResp := exchange.ObResponse{
		Symbol:    "BTCUSDT",
		Timestamp: 1670155864616,
		Bids: [][]string{
			{"16905.10", "2.785"},
			{"16904.10", "3.214"},
			{"16903.10", "6.123"},
			{"16902.10", "0.123"},
			{"16901.10", "2.123"},
		},
		Asks: [][]string{
			{"16905.90", "1.112"},
			{"16906.40", "0.100"},
			{"16907.70", "4.311"},
			{"16908.80", "1.110"},
			{"16909.25", "12.020"},
		},
	}
	ob := CreateSymbolObSnapshot(&obResp)

	// Testing
	if ob.Symbol != obResp.Symbol {
		t.Error("Incorrect: Wrong symbol returned")
	}
	if ob.Timestamp != obResp.Timestamp {
		t.Error("Incorrect: Timestamp is wrong!")
	}
	if ob.DecimalsPrice != countNumberOfDigitsAfterComma(obResp.Bids[0][0]) {
		t.Error("DecimalsPrice counted Incorrectly!")
	}
	if ob.DecimalsQuantity != countNumberOfDigitsAfterComma(obResp.Bids[0][1]) {
		t.Error("DecimalsQuantity counted Incorrectly!")
	}
	for i, bid := range ob.Bids {
		tPrice, tQty := convertPriceQuantityFromStringToFloat(obResp.Bids[i][0], obResp.Bids[i][1])
		if bid.Price != tPrice {
			t.Error("Incorrect Price: Wrong result of conversion!")
		}

		if bid.Quantity != tQty {
			t.Error("Incorrect Qty: Wrong result of conversion!")
		}
	}

}

func Test_PrepareForCsv(t *testing.T) {
	// Init
	timestampNow := time.Now().UnixMilli()
	obBids := BidsAsks{
		BidAskPriceLevel{16905.10, 2.785},
		BidAskPriceLevel{16904.10, 3.214},
		BidAskPriceLevel{16903.10, 6.123},
		BidAskPriceLevel{16902.10, 0.123},
		BidAskPriceLevel{16901.10, 2.123},
	}
	obAsks := BidsAsks{
		BidAskPriceLevel{16905.90, 1.112},
		BidAskPriceLevel{16906.40, 0.100},
		BidAskPriceLevel{16907.70, 4.311},
		BidAskPriceLevel{16908.80, 1.110},
		BidAskPriceLevel{16909.25, 12.020},
	}
	ob := SymbolObSnapshot{
		Symbol:           "BTCUSDT",
		Timestamp:        timestampNow,
		DecimalsPrice:    2,
		DecimalsQuantity: 3,
		Bids:             obBids,
		Asks:             obAsks,
	}

	// Testing
	csvLines, err := ob.PrepareForCsv()
	if err != nil {
		t.Error("Test failed")
		t.Error(err)
	}

	// Check if Bids or Asks are not empty
	bids := 0
	asks := 0
	for _, l := range csvLines {
		if l[1] == "bid" {
			bids++
		} else if l[1] == "ask" {
			asks++
		}
	}
	if bids == 0 {
		t.Error("ZERO Bids: Something ain't right!")
	}
	if asks == 0 {
		t.Error("ZERO Asks: Something ain't right!")
	}
}
