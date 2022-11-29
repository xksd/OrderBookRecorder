package main

import (
	"fmt"
	"time"

	"github.com/xksd/OrderBookRecorder/exchange"
	"github.com/xksd/OrderBookRecorder/utils"
)

func main() {
	// Initialize the settings
	// .Prod, .UserSymbolsList []string, .ApiKey, .ApiSecret
	// settings.
	settings := utils.GetSettings()
	// The build is PROD or DEV?
	PROD := settings.Prod
	// List of supported Symbols
	symbolsList := exchange.CreateListOfSymbols(settings.UserSymbolsList...)

	// Get number of cpu cores to calculate possible number of get requests
	// couCoresCount := runtime.NumCPU()
	// utils.PrintIntroduction(symbolsList, couCoresCount, PROD)

	// allOrderBooksAggregated
	// Prepare empty structs for ObSnapshots for all symbols for 1 day
	var allObAggr exchange.AllDailyAggregations

	for _, symbol := range symbolsList {
		var symbolDailySnapshots exchange.ObSnapshots_SymbolDailyAggr
		symbolDailySnapshots.Symbol = symbol
		allObAggr = append(allObAggr, symbolDailySnapshots)
	}

	// Start execution
	// Loop ticker
	ticker := time.NewTicker(5 * time.Second)
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-done:
				return
			case t := <-ticker.C:
				fmt.Println("Tick at", t)
				// aggregate ObSnaphots for symbol for current day
				for _, symbol := range symbolsList {
					s := symbol
					go func() {
						// var ob exchange.ObSnapshot
						// Get Order Book
						ob := exchange.GetObSnapshot(s, 5)
						// Add order book to all
						allObAggr.Add(s, ob)
					}()
				}
			}
		}
	}()

	if !PROD {
		time.Sleep(10 * time.Second)
	}

	ticker.Stop()
	done <- true

	fmt.Println("\nAFTER in main():")
	allObAggr.Print()
	fmt.Println()

}
