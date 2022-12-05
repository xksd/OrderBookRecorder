package exchange_test

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/xksd/OrderBookRecorder/exchange"
)

func TestGetUpdate(t *testing.T) {
	// Initial variables initialization
	symbolsList := []string{"BTCUSDT", "ETHUSDT", "SOLUSDT"}
	wsUrl := exchange.CreateWsUrl(symbolsList)

	// Create channel to exit the program
	// interrupt := make(chan os.Signal, 1)
	// defer close(interrupt)
	// signal.Notify(interrupt, os.Interrupt)

	wsCh := make(chan exchange.WsObResponse)
	wsChDone := make(chan bool)
	defer close(wsCh)
	defer close(wsChDone)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		exchange.GetUpdate(wsCh, wsChDone, wsUrl)
	}()

	updatesReceived := 0
	go func() {
		for {
			if updatesReceived == 5 {
				wsChDone <- true
				wg.Done()
			}

			upd := <-wsCh
			updatesReceived++

			timestampMin := time.Now().UnixMilli() - int64(time.Second)

			fmt.Println()
			fmt.Println("Updates received:", updatesReceived)

			fmt.Println(upd.Symbol, len(upd.Data.Asks), len(upd.Data.Bids))

			symbolMainCondition := func() bool {
				for _, symbol := range symbolsList {
					if upd.Symbol == symbol {
						return true
					}
				}
				return false
			}()
			symbolSecondaryCondition := func() bool {
				for _, symbol := range symbolsList {
					if upd.Data.Symbol == symbol {
						return true
					}
				}
				return false
			}()
			if !symbolMainCondition {
				t.Error("Incorrect: Wrong symbol in top struct level.", upd.Symbol)
			}

			if !symbolSecondaryCondition {
				t.Error("Incorrect: Wrong symbol in struct.Data level.", upd.Data.Symbol)
			}

			if upd.Data.Timestamp < timestampMin {
				t.Error("Incorrect: Timestamp is wrong, maybe delay is too big.")
			}

			if len(upd.Data.Bids) == 0 {
				t.Error("Incorrect: OrderBook Bids are empty!")
			}
			if len(upd.Data.Asks) == 0 {
				t.Error("Incorrect: OrderBook Asks are empty!")
			}

		}

	}()

	wg.Wait()
}
