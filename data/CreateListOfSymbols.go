package data

import "sort"

func CreateListOfSymbols(symbols ...string) []string {
	sort.Strings(symbols)

	for i := range symbols {
		symbols[i] += "USDT"
	}

	return symbols
}
