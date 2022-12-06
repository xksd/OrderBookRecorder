package data

import (
	"fmt"
	"testing"

	"github.com/xksd/OrderBookRecorder/exchange"
)

type testNumber struct {
	number   float64
	position int
}

var asks = BidsAsks{
	BidAskPriceLevel{16901.10, 2.123}, // 0
	BidAskPriceLevel{16902.10, 0.123}, // 1
	BidAskPriceLevel{16903.10, 6.123}, // 2
	BidAskPriceLevel{16904.10, 3.214}, // 3
	BidAskPriceLevel{16905.10, 2.785}, // 4
	BidAskPriceLevel{16906.10, 2.785}, // 5
	BidAskPriceLevel{16907.10, 2.785}, // 6
	BidAskPriceLevel{16908.10, 2.785}, // 7
}

var bids = BidsAsks{
	BidAskPriceLevel{16908.10, 2.785}, // 0
	BidAskPriceLevel{16907.10, 2.785}, // 1
	BidAskPriceLevel{16906.10, 2.785}, // 2
	BidAskPriceLevel{16905.10, 2.785}, // 3
	BidAskPriceLevel{16904.10, 3.214}, // 4
	BidAskPriceLevel{16903.10, 6.123}, // 5
	BidAskPriceLevel{16902.10, 0.123}, // 6
	BidAskPriceLevel{16901.10, 2.123}, // 7
}

func copyBidsAsksSlice(oldSlice *BidsAsks) *BidsAsks {
	// Copy to get separate slice
	newSlice := make(BidsAsks, len(*oldSlice))
	copy(newSlice, *oldSlice)
	return &newSlice
}

// Testing of updating existing Orderbook Snapshot BidsAsks fields
// from received WebSocket update.
func TestUpdatePriceLevelsFromWsTick(t *testing.T) {

	// Test cases:
	// 1) Update quantity on existing price level
	// 2) Add new price level
	// 3) Delete the price level

	// Testing 3 types of cases: index start (0), index middle, index end.

	// CASE 1: Update quantity on existing price level

	{
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
					{"16905.10", "95.100"},
					{"16901.10", "55.001"},
				},
				Asks: [][]string{
					{},
				},
			},
		}

		// Shadowing to intialize local variable for
		// existing orderbook snapshot BidsAsks fields
		bids := copyBidsAsksSlice(&bids)
		err := bids.UpdatePriceLevelsFromWsTick(&wsTick.Data.Bids, "bid")
		if err != nil {
			t.Error(err)
		}
		if (*bids)[0].Quantity != 99.010 || (*bids)[3].Quantity != 95.100 || (*bids)[7].Quantity != 55.001 {
			t.Error("Records are NOT properly updated: Quantity is wrong!")
		}

	}

	// CASE 2: Add new price level
	{
		// Initialize a variable containing a fake update from WebSocket
		wsTick := exchange.WsObResponse{
			Symbol: "btcusdt@depth",
			Data: exchange.WsObResponseData{
				Symbol:               "BTCUSDT",
				EventType:            "depthUpdate",
				Timestamp:            1670155864616,
				TimestampTransaction: 1670155864616,
				IdFirst:              1231,
				IdFinal:              1232,
				IdFinalLastStream:    1233,
				Bids: [][]string{
					{"16912.32", "99.010"},
					{"16904.35", "95.100"},
					{"16899.00", "135.001"},
				},
				Asks: [][]string{
					{},
				},
			},
		}

		// Shadowing to intialize local variable for
		// existing orderbook snapshot BidsAsks fields
		bids := copyBidsAsksSlice(&bids)
		err := bids.UpdatePriceLevelsFromWsTick(&wsTick.Data.Bids, "bid")
		if err != nil {
			t.Error(err)
		}
		if len(*bids) != 11 {
			t.Error("Slice is not complete - some elements are missing.")
			t.Errorf("Length should be: %v | Instead is: %v", 11, len(*bids))
		}
	}

	// CASE 3: Delete the price level
	{
		// Initialize a variable containing a fake update from WebSocket
		wsTick := exchange.WsObResponse{
			Symbol: "btcusdt@depth",
			Data: exchange.WsObResponseData{
				Symbol:               "BTCUSDT",
				EventType:            "depthUpdate",
				Timestamp:            1670155864616,
				TimestampTransaction: 1670155864616,
				IdFirst:              1231,
				IdFinal:              1232,
				IdFinalLastStream:    1233,
				Bids: [][]string{
					{"16908.10", "0.000"},
					{"16905.10", "0.000"},
					{"16901.10", "0.000"},
				},
				Asks: [][]string{
					{},
				},
			},
		}

		// Shadowing to intialize local variable for
		// existing orderbook snapshot BidsAsks fields
		bids := copyBidsAsksSlice(&bids)
		err := bids.UpdatePriceLevelsFromWsTick(&wsTick.Data.Bids, "bid")
		if err != nil {
			t.Error(err)
		}
		if len(*bids) != 5 {
			t.Error("Slice is not complete - some elements are missing.")
			t.Errorf("Length should be: %v | Instead is: %v", 11, len(*bids))
		}
	}
}

func TestCreateBidsAsks(t *testing.T) {
	obResponseBidsAsks := [][]string{
		{"16908.10", "2.782"}, // 0
		{"16907.10", "2.145"}, // 1
		{"16906.10", "3.322"}, // 2
		{"16905.10", "4.711"}, // 3
		{"16904.10", "3.214"}, // 4
		{"16903.10", "6.123"}, // 5
		{"16902.10", "0.123"}, // 6
		{"16901.10", "2.123"}, // 7
	}

	ba := CreateBidsAsks(&obResponseBidsAsks)

	if ba[0].Price != 16908.10 || ba[4].Price != 16904.10 || ba[7].Price != 16901.10 {
		t.Error("ba.Price don't match with original source")
	}
	if ba[0].Quantity != 2.782 || ba[4].Quantity != 3.214 || ba[7].Quantity != 2.123 {
		t.Error("ba.Price don't match with original source")
	}
}

func Test_shiftFromIndex(t *testing.T) {
	{
		obAsks := copyBidsAsksSlice(&asks)
		obAsks.shiftFromIndex(0)

		if (*obAsks)[0] != (*obAsks)[1] {
			t.Error("Shifted elements not match, they should be equal.")
		}
		if len(*obAsks) != 9 {
			t.Error("Slice must end up being 9 elements long")
			t.Error("Result length:", len(*obAsks))
		}
	}

	{
		obAsks := copyBidsAsksSlice(&asks)
		obAsks.shiftFromIndex(4)
		if (*obAsks)[4] != (*obAsks)[5] {
			t.Error("Shifted elements not match, they should be equal.")
		}
		if len(*obAsks) != 9 {
			t.Error("Slice must end up being 9 elements long")
			t.Error("Result length:", len(*obAsks))
		}
	}

	{
		obAsks := copyBidsAsksSlice(&asks)
		obAsks.shiftFromIndex(7)
		if (*obAsks)[7] != (*obAsks)[8] {
			t.Error("Shifted elements not match, they should be equal.")
		}
		if len(*obAsks) != 9 {
			t.Error("Slice must end up being 9 elements long")
			t.Error("Result length:", len(*obAsks))
		}
	}
	{
		obBids := copyBidsAsksSlice(&bids)
		obBids.shiftFromIndex(0)

		if (*obBids)[0] != (*obBids)[1] {
			t.Error("Shifted elements not match, they should be equal.")
		}
		if len(*obBids) != 9 {
			t.Error("Slice must end up being 9 elements long")
			t.Error("Result length:", len(*obBids))
		}
	}

	{
		obBids := copyBidsAsksSlice(&bids)
		obBids.shiftFromIndex(4)
		if (*obBids)[4] != (*obBids)[5] {
			t.Error("Shifted elements not match, they should be equal.")
		}
		if len(*obBids) != 9 {
			t.Error("Slice must end up being 9 elements long")
			t.Error("Result length:", len(*obBids))
		}
	}

	{
		obBids := copyBidsAsksSlice(&bids)
		obBids.shiftFromIndex(7)
		if (*obBids)[7] != (*obBids)[8] {
			t.Error("Shifted elements not match, they should be equal.")
		}
		if len(*obBids) != 9 {
			t.Error("Slice must end up being 9 elements long")
			t.Error("Result length:", len(*obBids))
		}
	}
}

func Test_deletePriceLevel(t *testing.T) {
	{
		obAsks := copyBidsAsksSlice(&asks)
		i := 0
		obAsks.deletePriceLevel(i)
		if (*obAsks)[i].Price != 16902.10 {
			t.Error("Deleting Price Level FAILED. Desired: ", 16902.10, "| Result:", (*obAsks)[i].Price)
		}
	}

	{
		obAsks := copyBidsAsksSlice(&asks)
		i := 3
		obAsks.deletePriceLevel(i)
		if (*obAsks)[i].Price != 16905.10 {
			t.Error("Deleting Price Level FAILED. Desired: ", 16905.10, "| Result:", (*obAsks)[i].Price)
		}
	}

	{
		obAsks := copyBidsAsksSlice(&asks)
		i := 7
		obAsks.deletePriceLevel(i)
		if len((*obAsks))-1 > 7 {
			t.Error("Deleting Price Level FAILED. Price Level was NOT deleted.")
		}
	}

	{
		obBids := copyBidsAsksSlice(&bids)
		i := 0
		obBids.deletePriceLevel(i)
		if (*obBids)[i].Price != 16907.10 {
			t.Error("Deleting Price Level FAILED. Desired: ", 16907.10, "| Result:", (*obBids)[i].Price)
		}
	}

	{
		obBids := copyBidsAsksSlice(&bids)
		i := 3
		obBids.deletePriceLevel(i)
		if (*obBids)[i].Price != 16904.10 {
			t.Error("Deleting Price Level FAILED. Desired: ", 16904.10, "| Result:", (*obBids)[i].Price)
		}
	}

	{
		obBids := copyBidsAsksSlice(&bids)
		i := 7
		obBids.deletePriceLevel(i)
		if len((*obBids))-1 > 7 {
			t.Error("Deleting Price Level FAILED. Price Level was NOT deleted.")
		}
	}
}

func Test_rewritePriceLevel(t *testing.T) {
	// Copy to get separate slice
	obAsks := make(BidsAsks, len(asks))
	copy(obAsks, asks)

	index := obAsks.findPriceLevelIndex(16908.10, "asks")
	obAsks.rewritePriceLevel(index, 111.00, 77.000)

	// Test: Regular
	if obAsks[index].Price != 111.00 || obAsks[index].Quantity != 77.000 {
		t.Error("Incorrect: Rewriting price level failed!")
	}

	// Test: Corner cases
	indexOutOfRange1 := obAsks.findPriceLevelIndex(25000.01, "asks")
	indexOutOfRange2 := -1
	if err := obAsks.rewritePriceLevel(indexOutOfRange1, 111.00, 77.000); err == nil {
		fmt.Println(err)
		t.Error("Error checking for Index Out of Range case FAILED")
	}
	if err := obAsks.rewritePriceLevel(indexOutOfRange2, 111.00, 77.000); err == nil {
		fmt.Println(err)
		t.Error("Error checking for Index Out of Range case FAILED")
	}

}

func Test_findPriceLevelIndex(t *testing.T) {

	var asksSearchNumbers = []testNumber{
		{16801.00, 0}, {16900.50, 0}, {16901.10, 0},
		{16901.50, 1}, {16902.10, 1},
		{16902.50, 2}, {16903.10, 2},
		{16903.50, 3}, {16904.10, 3},
		{16904.50, 4}, {16905.10, 4},
		{16905.50, 5}, {16906.10, 5},
		{16906.50, 6}, {16907.10, 6},
		{16907.50, 7}, {16908.10, 7},
		{16910.50, 8}, {16912.50, 8}, {25000.10, 8},
	}

	var bidsSearchNumbers = []testNumber{
		{16910.50, 0}, {16912.50, 0}, {25000.10, 0},
		{16908.10, 0},
		{16907.50, 1}, {16907.10, 1},
		{16906.50, 2}, {16906.10, 2},
		{16905.50, 3}, {16905.10, 3},
		{16904.50, 4}, {16904.10, 4},
		{16903.50, 5}, {16903.10, 5},
		{16902.50, 6}, {16902.10, 6},
		{16901.50, 7}, {16901.10, 7},
		{16801.00, 8}, {16900.50, 8},
	}

	for _, searchSetAsks := range asksSearchNumbers {
		foundIndex := asks.findPriceLevelIndex(searchSetAsks.number, "ask")
		if foundIndex != searchSetAsks.position {
			t.Error("Desired:", searchSetAsks.position, " | Test result:", foundIndex, "| Number:", searchSetAsks.number)
		}
	}

	for _, searchSetBids := range bidsSearchNumbers {
		foundIndex := bids.findPriceLevelIndex(searchSetBids.number, "bid")

		if foundIndex != searchSetBids.position {
			t.Error("Desired:", searchSetBids.position, " | Test result:", foundIndex, "| Number:", searchSetBids.number)
		}
	}
}
