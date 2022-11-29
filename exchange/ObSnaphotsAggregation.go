package exchange

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

//

type AllDailyAggregations []ObSnapshots_SymbolDailyAggr

// Single File of Snapshots
// 1 file = 1 day, 1 symbol
type ObSnapshots_SymbolDailyAggr struct {
	Symbol    string
	T_Starts  int64
	T_Ends    int64
	Snapshots []ObSnapshot
}

func (allObAggr AllDailyAggregations) Add(symbol string, ob ObSnapshot) {
	for i := 0; i < len(allObAggr); i++ {
		// find the right symbol in obAggr
		if allObAggr[i].Symbol == symbol {
			// append new snapshot to that symbol
			allObAggr[i].Snapshots = append(allObAggr[i].Snapshots, ob)
		}
	}
}

func (allObAggr *AllDailyAggregations) Print() {
	fmt.Println()

	// Find the longest symbol name
	longestSymbolName := 0
	for _, symbol := range *allObAggr {
		counted := utf8.RuneCountInString(symbol.Symbol)

		if counted > longestSymbolName {
			longestSymbolName = counted
		}
	}

	for _, symbol := range *allObAggr {
		symbolLength := utf8.RuneCountInString(symbol.Symbol)
		fmt.Println("SYMBOL: ", symbol.Symbol, strings.Repeat(" ", longestSymbolName-symbolLength), "   |   LEN: ", len(symbol.Snapshots))
		// for _, s := range symbol.Snapshots {
		// 	fmt.Println("TIME:", time.UnixMilli(s.Timestamp))
		// 	// fmt.Println("BIDS:", s.Bids)
		// 	// fmt.Println("ASKS:", s.Asks)
		// }
	}
}
