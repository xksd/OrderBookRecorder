package data

import (
	"fmt"
	"log"
	"strconv"
)

// Count how many numbers after decimal point the string-number has
func countNumberOfDigitsAfterComma(stringNumber string) int {
	count := 0
	decimalPassed := false
	for _, r := range stringNumber {
		if decimalPassed {
			count++
		}
		if r == '.' {
			decimalPassed = true
		}
	}
	return count
}

// When converting float64 price or qty back to string for CSV,
// we want to maintain the format when each number has exact same
// number of zeroes after decimal. For example: "1.00", "13.010"
func createNumberString(number float64, numberOfDigitsAfterComma int) (outputString string) {
	switch numberOfDigitsAfterComma {
	case 0:
		outputString = fmt.Sprintf("%v", int(number))
	case 1:
		outputString = fmt.Sprintf("%.1f", number)
	case 2:
		outputString = fmt.Sprintf("%.2f", number)
	case 3:
		outputString = fmt.Sprintf("%.3f", number)
	case 4:
		outputString = fmt.Sprintf("%.4f", number)
	case 5:
		outputString = fmt.Sprintf("%.5f", number)
	case 6:
		outputString = fmt.Sprintf("%.6f", number)
	case 7:
		outputString = fmt.Sprintf("%.7f", number)
	case 8:
		outputString = fmt.Sprintf("%.8f", number)
	}
	return
}

// Convert String to Float64 for Price and Quantity
// returns "[price, quantity]"
func convertPriceQuantityFromStringToFloat(price string, quantity string) (float64, float64) {
	// Convert string price and quantity to float64
	p, err := strconv.ParseFloat(price, 64)
	if err != nil {
		log.Fatal("Error in data.UpdateFromWsTick() > converting Price")
		log.Fatal(err)
	}
	q, err := strconv.ParseFloat(quantity, 64)
	if err != nil {
		log.Fatal("Error in data.UpdateFromWsTick() > converting Price")
		log.Fatal(err)
	}
	return p, q
}
