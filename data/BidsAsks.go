package data

import (
	"errors"
	"fmt"
	"log"
)

// One side of the OrderBook
// either Bids or Asks
type BidsAsks []BidAskPriceLevel

// One price level in Orderbook
type BidAskPriceLevel struct {
	Price    float64
	Quantity float64
}

func CreateBidsAsks(obRespBidsAsks *[][]string) BidsAsks {
	var bidsAsks BidsAsks
	// Loop over the OB Price Level of the obResponse struct
	for _, ObRespPriceLevel := range *obRespBidsAsks {
		pl := BidAskPriceLevel{}

		// Convert number simply to float64,
		// as we know number or decimals later will be easy to convert back to string
		pl.Price, pl.Quantity = convertPriceQuantityFromStringToFloat(ObRespPriceLevel[0], ObRespPriceLevel[1])

		bidsAsks = append(bidsAsks, pl)
	}

	return bidsAsks
}

// Update Price Levels From WebSocket Tick
//
// Three Types of updates:
// 1) Update quantity on existing price level
// 2) Add new price level
// 3) Delete the price level
//
// Once we receive an update from WebSocket, it will contain the following fields:
//   - fields with Price existing in our orderbook snapshot and
//     Quanity == 0 : these price levels must be deleted
//   - fields with Price existing in our orderbook snapshot and
//     Quantity > 0 : these price levels must be updated (i.e. overwritten)
//   - fields with Price not existed before in our orderbook snapshot and some Quantity:
//     these price levels are new and should be added into the orderbook
//
// (i) Bids are Descending order, Asks are Ascending order
func (ba *BidsAsks) UpdatePriceLevelsFromWsTick(wsTickBidsAsks *[][]string, side string) error {
	for i := range *wsTickBidsAsks {
		// Convert strings from Websocket response to float64 numbers
		newPrice, newQuantity := convertPriceQuantityFromStringToFloat((*wsTickBidsAsks)[i][0], (*wsTickBidsAsks)[i][1])
		// Find index of price level that needs an update
		priceLevelIndex := ba.findPriceLevelIndex(newPrice, side)

		// Update required Price Level in BidsAsks
		// Check if priceLevelIndex is inside the range of existing "ba"-slice
		if priceLevelIndex < len(*ba) {
			// If price level exists in OrderBook => update quantity (overwrite)
			if (*ba)[priceLevelIndex].Price == newPrice {
				// If Quanity is not 0, update the Quantity
				if newQuantity != 0.0 {
					ba.rewritePriceLevel(priceLevelIndex, newPrice, newQuantity)
					// If Quantity == 0 => price level must be deleted
				} else if newQuantity == 0.0 {
					ba.deletePriceLevel(priceLevelIndex)
				}
				// If Price NOT existed in orderbook before -> add new price level
			} else if (*ba)[priceLevelIndex].Price != newPrice {
				// Check to avoid inserting empty price level if it was already deleted before
				if newQuantity != 0.0 {
					err := ba.shiftFromIndex(priceLevelIndex)
					if err != nil {
						return errors.New("shiftFromIndex error:" + fmt.Sprintln(err))
					}
					err = ba.rewritePriceLevel(priceLevelIndex, newPrice, newQuantity)
					if err != nil {
						return errors.New("shiftFromIndex error:" + fmt.Sprintln(err))
					}
				}
			}
		} else {
			// If priceLevelIndex is greater than slice range => new element should be added
			// Check to avoid inserting empty price level
			if newQuantity != 0.0 {
				(*ba) = append((*ba), BidAskPriceLevel{newPrice, newQuantity})
			}
		}
	}
	return nil
}

// Insert new price level before the bigger one,
// or overwrite with new quantity
func (ba *BidsAsks) rewritePriceLevel(index int, newPrice float64, newQuantity float64) error {
	if index < 0 || index > len(*ba)-1 {
		return errors.New("rewritePriceLevel(): Index out of range")
	}

	(*ba)[index].Price = newPrice
	(*ba)[index].Quantity = newQuantity

	return nil
}

func (ba *BidsAsks) findPriceLevelIndex(newPrice float64, side string) (index int) {
	// Check if "side" is right
	if side != "bid" && side != "bids" && side != "ask" && side != "asks" {
		log.Fatal("ERROR: Wrong 'side' parameter!")
	}

	// Set closures to work with ascending and descending order
	var f1 func(float64, float64) bool
	var f2 func(float64, float64) bool
	if side == "bid" || side == "bids" {
		f1 = func(p float64, m float64) bool { return p > m }
		f2 = func(p float64, m float64) bool { return p < m }
	} else if side == "ask" || side == "asks" {
		f1 = func(p float64, m float64) bool { return p < m }
		f2 = func(p float64, m float64) bool { return p > m }
	}

	// Binary search with returning closest value if not exist
	lowerBound := 0
	upperBound := len((*ba)) - 1
	p := newPrice

	for lowerBound <= upperBound {
		midpoint := (upperBound + lowerBound) / 2
		index = midpoint
		valueAtMidpoint := (*ba)[midpoint].Price

		if p == valueAtMidpoint {
			index = midpoint
			return
		} else if f1(p, valueAtMidpoint) {
			upperBound = midpoint - 1
		} else {
			lowerBound = midpoint + 1
		}

		if upperBound-lowerBound == 0 {
			index = lowerBound
			if f2(p, (*ba)[lowerBound].Price) {
				index = index + 1
			}
			return
		}

	}
	return
}

// Deleting price level means shifting slice backwards
// one element in place of index
func (ba *BidsAsks) deletePriceLevel(index int) error {
	if index < 0 || index > len(*ba)-1 {
		return errors.New("rewritePriceLevel(): Index out of range")
	}

	// Shift slice backwards one element in a place of index
	for i := index; i < len(*ba); i++ {
		if i == len(*ba)-1 {
			(*ba) = (*ba)[:len(*ba)-1]
			return nil
		}
		(*ba)[i] = (*ba)[i+1]
	}

	return nil
}

// Shift elements to index+1
func (ba *BidsAsks) shiftFromIndex(index int) error {
	if index < 0 || index > len(*ba)-1 {
		return errors.New("rewritePriceLevel(): Index out of range")
	}

	for i := len(*ba) - 1; i >= index; i-- {
		if i == len(*ba)-1 {
			*ba = append(*ba, (*ba)[i])
		}
		if i > index {
			(*ba)[i] = (*ba)[i-1]
		}
	}

	return nil
}
