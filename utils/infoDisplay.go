package utils

import "fmt"

func PrintIntroduction(symbolsList []string, PROD bool) {
	// create string for State (production or development)
	var state string
	if state = "development"; PROD {
		state = "production"
	}

	fmt.Println("\n----- started -----", state, "----- ")
	fmt.Println("SYMBOLS: ", len(symbolsList), "   |  ")
	// fmt.Println(symbolsList)
	fmt.Println("-------------------------------------")
	fmt.Println()
}
