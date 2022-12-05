package storage_test

import (
	"testing"

	"github.com/xksd/OrderBookRecorder/data"
	"github.com/xksd/OrderBookRecorder/storage"
)

func TestSaveToCSV(t *testing.T) {
	obBids := data.BidsAsks{
		data.BidAskPriceLevel{Price: 16905.10, Quantity: 2.785},
		data.BidAskPriceLevel{Price: 16904.10, Quantity: 3.214},
		data.BidAskPriceLevel{Price: 16903.10, Quantity: 6.123},
		data.BidAskPriceLevel{Price: 16902.10, Quantity: 0.123},
		data.BidAskPriceLevel{Price: 16901.10, Quantity: 2.123},
	}
	obAsks := data.BidsAsks{
		data.BidAskPriceLevel{Price: 16905.90, Quantity: 1.112},
		data.BidAskPriceLevel{Price: 16906.40, Quantity: 0.100},
		data.BidAskPriceLevel{Price: 16907.70, Quantity: 4.311},
		data.BidAskPriceLevel{Price: 16908.80, Quantity: 1.110},
		data.BidAskPriceLevel{Price: 16909.25, Quantity: 12.020},
	}
	ob := data.SymbolObSnapshot{
		Symbol:           "BTCUSDT",
		Timestamp:        1670155864616,
		DecimalsPrice:    2,
		DecimalsQuantity: 3,
		Bids:             obBids,
		Asks:             obAsks,
	}

	if err := storage.SaveToCSV(&ob, "./test_csv_files"); err != nil {
		t.Error("Writing funciton failed!")
		t.Error("ERR:", err)
	}

}

func TestOpenCSVFile(t *testing.T) {
	if !storage.FileExists("./test_csv_files/2022_12_4__BTCUSDT.csv") {
		t.Error("File doesn't exist. First create file, then try to test opening it :-)")
	}

	csvFile, err := storage.OpenCSVFile("./test_csv_files", "2022_12_4__BTCUSDT.csv")
	if err != nil {
		t.Error("Opening File FAILED!")
		t.Error("ERROR: ", err)
	}
	defer csvFile.Close()

}
