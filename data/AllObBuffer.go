package data

import (
	"errors"
	"fmt"
	"unicode/utf8"

	"github.com/xksd/OrderBookRecorder/exchange"
)

// Each slice element is a single Write-batch to CSV file
// i.e. a single snapshot update to be recorded.
type AllObBuffer []SymbolObSnapshot

func (allObBuffer *AllObBuffer) Append(ob *SymbolObSnapshot) {
	*allObBuffer = append(*allObBuffer, *ob)
}

// UpdateFromWsTick(wsTick) (insert and search ordered array)
func (allObBuffer *AllObBuffer) UpdateFromWsTick(wsTick exchange.WsObResponse) error {
	// Find corresponding ob.Symbol in Buffer
	for i := range *allObBuffer {
		if (*allObBuffer)[i].Symbol == wsTick.Symbol {
			// Save current SymbolObSnapshot to a local variable
			ob := (*allObBuffer)[i]

			(*allObBuffer)[i].Timestamp = wsTick.Data.Timestamp

			// Update Ob with passed Ws update
			if err := ob.Bids.UpdatePriceLevelsFromWsTick(&wsTick.Data.Bids, "bid"); err != nil {
				return errors.New("Error updating Bids" + fmt.Sprint(err))
			}
			if err := ob.Asks.UpdatePriceLevelsFromWsTick(&wsTick.Data.Bids, "ask"); err != nil {
				return errors.New("Error updating Asks" + fmt.Sprint(err))

			}

		} else {
			// log.Fatal("error in allObBuffer.UpdateFromWsTick()")
			// log.Fatal("(!) ERROR: Can't find the symbol!")
			// return errors.New("(!) allObBuffer.UpdateFromWsTick(): Can't find the symbol!")
		}
	}

	return nil
}

func (allObBuffer *AllObBuffer) Print() {
	fmt.Println()

	// Find the longest symbol name
	longestSymbolName := 0
	for _, symbol := range *allObBuffer {
		counted := utf8.RuneCountInString(symbol.Symbol)

		if counted > longestSymbolName {
			longestSymbolName = counted
		}
	}

}
