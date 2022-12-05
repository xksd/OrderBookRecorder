package data

import (
	"testing"
)

func TestCreateListOfSymbols(t *testing.T) {
	symbols := []string{"1000SHIB", "BTC", "ETH"}
	sl := CreateListOfSymbols(symbols...)

	if sl[0] != "1000SHIBUSDT" || sl[1] != "BTCUSDT" || sl[2] != "ETHUSDT" {
		t.Error("The symbol is wrong")
	}
}
