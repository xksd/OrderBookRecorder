package main

import (
	"log"
	"sync"
	"time"

	"github.com/xksd/OrderBookRecorder/data"
	"github.com/xksd/OrderBookRecorder/exchange"
	"github.com/xksd/OrderBookRecorder/storage"
	"github.com/xksd/OrderBookRecorder/utils"
)

func main() {
	// Initialize the settings
	settings := utils.ReadSettings()
	// List of supported Symbols
	symbolsList := data.CreateListOfSymbols(settings.UserSymbolsList...)
	// Warm Introduction
	utils.PrintIntroduction(symbolsList, settings.Prod)

	// Fill the Buffer with the first OrderBook Snapshots for each symbol,
	// it will used for the following WebSocket updates.
	var allObBuffer data.AllObBuffer

	for _, symbol := range symbolsList {
		// Get Starting OrderBook Snapshots
		obResp := exchange.GetObResponse(symbol, 0)

		// Convert ObResponse to SymbolObSnapshot
		ob := data.CreateSymbolObSnapshot(&obResp)

		// Add snapshot of current symbol to the Buffer
		allObBuffer.Append(ob)

	}

	// Loop ticker
	ticker := time.NewTicker(time.Duration(settings.Timeframe) * time.Second)
	// Channels for Websocket tick updates and close command
	wsCh := make(chan exchange.WsObResponse)
	wsChDone := make(chan bool)
	defer close(wsCh)
	defer close(wsChDone)

	var wg sync.WaitGroup
	wg.Add(1)

	// WebSocket Tick Data Update.
	// Using allObBuffer to buffer incoming ws updates,
	// so the buffer is always ready to be saved to csv when ticker signals.

	// Create URL for Websocket stream, to tell what symbols to receive.
	wsUrl := exchange.CreateWsUrl(symbolsList)

	// Start WebSocket stream to get updates in a loop,
	// each update sent to channel.
	go func() {
		exchange.GetUpdate(wsCh, wsChDone, wsUrl)
	}()

	// Receive signals from channels:
	// update from websocket or stop signal
	go func() {
		for {
			select {
			case <-wsChDone:
				ticker.Stop()
				return
			case wsTick := <-wsCh:
				// Update the Buffer with new tick data
				err := allObBuffer.UpdateFromWsTick(wsTick)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}()

	// Write to CSV-file from AllObBuffer single ob, at ticker signal.
	go func() {
		for {
			select {
			case <-ticker.C:
				for _, ob := range allObBuffer {
					if err := storage.SaveToCSV(&ob, settings.CSVFolder); err != nil {
						log.Fatal("Function SaveToCSV Failed from MAIN file")
						log.Fatal("Reason:", err)
					}
				}
			}
		}
	}()

	wg.Wait()

}
