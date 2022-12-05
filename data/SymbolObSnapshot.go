package data

import (
	"errors"
	"strconv"

	"github.com/xksd/OrderBookRecorder/exchange"
)

// SymbolObSnapshot is a snapshot for single symbol
// that will be in future an element of AllObBuffer slice

// First SymbolObSnapshot is received by converting
// ObResponse from http-call into it.

type SymbolObSnapshot struct {
	Symbol           string
	Timestamp        int64
	DecimalsPrice    int
	DecimalsQuantity int
	Bids             BidsAsks
	Asks             BidsAsks
}

// Create SymbolObSnapshot from ObResponse
func CreateSymbolObSnapshot(obResp *exchange.ObResponse) *SymbolObSnapshot {
	// init variable
	var ob SymbolObSnapshot

	// Set the Symbol name
	ob.Symbol = obResp.Symbol

	// Findout the number of digits after the decimal in Price and Quantity
	ob.DecimalsPrice = countNumberOfDigitsAfterComma(obResp.Bids[0][0])
	ob.DecimalsQuantity = countNumberOfDigitsAfterComma(obResp.Asks[0][1])

	// Set DateTime for current Snapshot
	// in UnixMilli because it's handy to use the same as Binance
	ob.Timestamp = obResp.Timestamp

	// Convert Bids and Asks to required format.
	ob.Bids = CreateBidsAsks(&obResp.Bids)
	ob.Asks = CreateBidsAsks(&obResp.Asks)

	return &ob
}

// CSV Write function should receive a prepared data structure
// -- merged Bids and Asks into single slice
// -- price floats converted to strings to keep 0 after decimal
// format: (timestamp, side, price, quantity)

// Convert (ob SymbolObSnapshot) to CSV format:
func (ob *SymbolObSnapshot) PrepareForCsv() ([][]string, error) {
	var csvLines [][]string
	if len(ob.Bids) == 0 {
		return nil, errors.New("PrepareForCSV(): BIDS are empty")
	}
	if len(ob.Asks) == 0 {
		return nil, errors.New("PrepareForCSV(): ASKS are empty")
	}

	for i := range ob.Bids {
		// Convert float64 to string with zeroes after decimal.
		stringPrice := createNumberString(ob.Bids[i].Price, ob.DecimalsPrice)
		stringQuantity := createNumberString(ob.Bids[i].Quantity, ob.DecimalsPrice)
		// Compose a CSV line, which is a slice of strings
		csvLine := []string{
			strconv.FormatInt(ob.Timestamp, 10), "bid", stringPrice, stringQuantity,
		}
		// Append csvLine to all Csv Lines
		csvLines = append(csvLines, csvLine)
	}

	for i := range ob.Asks {
		// Convert float64 to string with zeroes after decimal.
		stringPrice := createNumberString(ob.Asks[i].Price, ob.DecimalsPrice)
		stringQuantity := createNumberString(ob.Asks[i].Quantity, ob.DecimalsPrice)
		// Compose a CSV line, which is a slice of strings
		csvLine := []string{
			strconv.FormatInt(ob.Timestamp, 10), "ask", stringPrice, stringQuantity,
		}
		// Append csvLine to all Csv Lines
		csvLines = append(csvLines, csvLine)
	}

	return csvLines, nil
}
