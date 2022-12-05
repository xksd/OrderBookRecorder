package utils_test

import (
	"strings"
	"testing"

	"github.com/xksd/OrderBookRecorder/utils"
)

// Timeframe: 1
// Prod: false
// CSVFolder: ./csv_files
// UserSymbolsList: [BTC SOL]
// ApiKey: c6b1yJn7IXMoF6elATCzi5F5pMKe6kkjjpHDMdXltrKjWLJ9gGfSpHMwLfP7iQkf
// M1BBuHKrm0RWZAvDeN2EU145o2whLoPh2kytatuk0S2n3nyhqBN7X0CvJ7pSoVdo
func TestReadSettings(t *testing.T) {
	s := utils.ReadSettings()

	if s.Timeframe == 0 {
		t.Error("Timeframe can't be 0 seconds!")
	}
	if !strings.Contains(s.CSVFolder, "csv") {
		t.Error("Check your CSV Folder path, it doesn't look right!")
	}
	if len(s.UserSymbolsList) == 0 {
		t.Error("Symbols list is EMPTY!")
	}
	if len(s.ApiKey) != 64 {
		t.Error("ApiKey is wrong, it should be 64 symbols long")
		t.Errorf("Instead it's %v symbols long.", len(s.ApiKey))
	}
	if len(s.ApiSecret) != 64 {
		t.Error("ApiSecret is wrong, it should be 64 symbols long.")
		t.Errorf("Instead it's %v symbols long.", len(s.ApiSecret))
	}
}
