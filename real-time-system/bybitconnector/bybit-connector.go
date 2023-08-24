package bybitconnector

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	bybitWSURL = "wss://stream.bybit.com/v5/public/linear"

	backoffFactor = 2 // Multiplier for each subsequent reconnection delay
)

var (
	reconnectDelay       = 5 // Delay in seconds before attempting a reconnection
	maxReconnectAttempts = 5 // Maximum number of reconnection attempts

	IsOrderBookReady bool = false
	IsTradeReady     bool = false
	IsTickerReady    bool = false
)

func ConnectToBybit() error {
	attempts := 0
	for attempts < maxReconnectAttempts {
		err := connectAndListen()
		if err == nil {
			return nil // Exit if successful
		}

		log.Printf("Error connecting or listening: %v", err)
		log.Printf("Attempt %d/%d. Reconnecting in %d seconds...", attempts+1, maxReconnectAttempts, reconnectDelay)
		time.Sleep(time.Duration(reconnectDelay) * time.Second)

		// Increase the delay for the next attempt, with a backoff factor
		reconnectDelay *= backoffFactor
		attempts++
	}

	log.Println("Max reconnection attempts reached. Exiting...")
	return errors.New("max reconnection attempts reached")
}

func connectAndListen() error {
	// Establish a WebSocket connection
	conn, _, err := websocket.DefaultDialer.Dial(bybitWSURL, nil)
	if err != nil {
		return err
	}
	defer conn.Close()

	// Subscribe to the necessary channels
	channels := []string{
		"orderbook.50.BTCUSDT",
		"publicTrade.BTCUSDT",
		"tickers.BTCUSDT",
	}

	for _, channel := range channels {
		err := conn.WriteJSON(map[string]interface{}{
			"op":   "subscribe",
			"args": []string{channel},
		})
		if err != nil {
			log.Printf("Failed to subscribe to channel %s: %v", channel, err)
		}
	}

	maxReconnectAttempts = 5

	// Listen for incoming messages
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			return err
		}
		handleMessage(message)
	}
}

func handleMessage(message []byte) {
	var topic struct {
		Topic string `json:"topic"`
	}

	err := json.Unmarshal(message, &topic)
	if err != nil {
		log.Println("Error parsing topic:", err)
		return
	}

	switch topic.Topic {
	case "orderbook.50.BTCUSDT":
		var orderBook OrderBookSnapshot

		err := json.Unmarshal(message, &orderBook)
		if err != nil {
			log.Println("Error parsing order book:", err)
			return
		}
		ProcessOrderBook(orderBook.Data)

		if RecentTrades.Len() >= recentTradePeriod {
			IsOrderBookReady = true
		}

	case "publicTrade.BTCUSDT":
		var trade TradeData
		err := json.Unmarshal(message, &trade)
		if err != nil {
			log.Println("Error parsing trade:", err)
			return
		}
		ProcessTrade(trade)
		IsTradeReady = true

	case "tickers.BTCUSDT":
		var ticker Ticker
		err := json.Unmarshal(message, &ticker)
		if err != nil {
			fmt.Println(string(message))
			log.Println("Error parsing ticker:", err)
			return
		}
		ProcessTicker(ticker)
		IsTickerReady = true
	}
}
