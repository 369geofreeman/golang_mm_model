package main

import (
	"fmt"
	"log"
	"time"

	"github.com/369geofreeman/inventory-control/real-time-system/bybitconnector"
	"github.com/369geofreeman/inventory-control/real-time-system/optimization"
)

func main() {
	fmt.Println("Starting the real-time system...")

	// Create an Inventory object with initial cash balance, crypto balance, and trading fee
	tradingFee := 0.0200 // Bybit's trading fee
	initialCashBalance := 1000.0
	initialCryptoBalance := 0.12345
	inventory := optimization.NewInventory(initialCashBalance, initialCryptoBalance, tradingFee)

	// Establish connection to Bybit in a Goroutine
	go func() {
		err := bybitconnector.ConnectToBybit()
		if err != nil {
			log.Fatalf("Error connecting to Bybit: %v", err)
		}
	}()

	for {
		// Wait until we have received at least one of each type of message
		if !bybitconnector.IsOrderBookReady || !bybitconnector.IsTradeReady || !bybitconnector.IsTickerReady {
			time.Sleep(5 * time.Second) // Wait for 5 seconds before checking again
			continue
		}

		// Fetch market data
		currentPrice := bybitconnector.MidPrice
		volatility := bybitconnector.Volatility
		liquidity := bybitconnector.Liquidity
		orderBookDepth := bybitconnector.OrderBookDepth

		// Optimize spread, passing the inventory object
		optimalAsk, optimalBid, _ := optimization.OptimizeSpread(currentPrice, inventory, volatility, liquidity, orderBookDepth)
		fmt.Printf("Optimal Bid: %f\nOptimal Ask: %f\nPrice: %f\n", optimalBid, optimalAsk, currentPrice)

		// Execute trades based on current price and optimal bid/ask
		inventory.TradeExecuted(currentPrice, optimalBid, optimalAsk)

		// Sleep for the determined time before the next optimization
		sleepTime := optimization.GetOptimizationFrequency()
		fmt.Printf("Time till next optimisation: %d\n", optimization.GetOptimizationFrequency())
		time.Sleep(time.Duration(sleepTime) * time.Minute)
	}
}
