package data

import (
	"testing"
)

func Test_countNumberOfDigitsAfterComma(t *testing.T) {
	desired := []int{7, 2, 3}
	numbersToTest := []string{"0.0004000", "132.00", "1932.100"}

	for i, number := range numbersToTest {
		result := countNumberOfDigitsAfterComma(number)
		if result != desired[i] {
			t.Error("Incorrect!")
			t.Errorf("Result:  %d digits after decimal for %v", desired[i], number)
			t.Errorf("Must be: %d digits after the decimal", result)
		}
	}
}

func Test_createNumberString(t *testing.T) {
	type testNumber struct {
		number                    float64
		numbersOfDigitsAfterComma int
	}
	set := []testNumber{
		{0.0004000, 7},
		{132.00, 2},
		{1932.100, 3},
	}
	desired := []string{"0.0004000", "132.00", "1932.100"}

	for i, num := range set {
		result := createNumberString(num.number, num.numbersOfDigitsAfterComma)
		if result != desired[i] {
			t.Error("Incorrect!")
			t.Error("Result:  " + result)
			t.Error("Must be: " + desired[i])
		}
	}
}

func Test_convertPriceQuantityFromStringToFloat(t *testing.T) {
	set := [][2]string{{"0.0004000", "32323.00"}, {"132.00", "329.21"}, {"1932.100", "111"}}
	desired := [][2]float64{{0.0004000, 32323}, {132.00, 329.21}, {1932.100, 111}}

	for i, line := range set {
		price, quantity := convertPriceQuantityFromStringToFloat(line[0], line[1])
		if price != desired[i][0] {
			t.Error("Price Incorrect!")
			t.Errorf("Result:  %v", price)
			t.Errorf("Must be: %v", desired[i][0])
		}
		if quantity != desired[i][1] {
			t.Error("Quantity Incorrect!")
			t.Errorf("Result:  %v", quantity)
			t.Errorf("Must be: %v", desired[i][1])
		}
	}
}
