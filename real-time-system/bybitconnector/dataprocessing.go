package bybitconnector

import (
	"container/list" // For using list as a deque
	"log"
	"math"
	"strconv"
	"sync"

	"github.com/369geofreeman/inventory-control/real-time-system/optimization"
)

// Constants
const recentTradePeriod = 50 // Number of recent trades to consider for volatility
const liquidityRange = 0.01 // 0.5% of mid-price for liquidity calculation

// Data structures
var (
	MidPrice       float64
	Volatility     float64
	Liquidity      float64
	OrderBookDepth float64
	Mutex          sync.RWMutex

	RecentTrades = list.New() // Deque to hold recent trades

)

// ProcessTicker processes the ticker data and updates the mid-price
func ProcessTicker(ticker Ticker) {
	Mutex.Lock()
	defer Mutex.Unlock()

	bidPrice := parseFloat(ticker.Data.Bid1Price)
	askPrice := parseFloat(ticker.Data.Ask1Price)

	MidPrice = (bidPrice + askPrice) / 2

	// log.Printf("MidPrice: %f", MidPrice)
}

// ProcessTrade processes trade data and updates the volatility
func ProcessTrade(trade TradeData) {
	Mutex.Lock()
	defer Mutex.Unlock()

	// Push the new trade into the RecentTrades deque
	price := parseFloat(trade.Data[0].Price)
	RecentTrades.PushFront(price)

	// Remove the oldest trade if we have more than our period
	if RecentTrades.Len() > recentTradePeriod {
		RecentTrades.Remove(RecentTrades.Back())
	}

	// Calculate volatility as the standard deviation of recent trade prices
	var sum float64
	var sumOfSquares float64
	for e := RecentTrades.Front(); e != nil; e = e.Next() {
		val := e.Value.(float64)
		sum += val
		sumOfSquares += val * val
	}
	mean := sum / float64(RecentTrades.Len())
	Volatility = math.Sqrt(sumOfSquares/float64(RecentTrades.Len()) - mean*mean)

	optimization.AdjustEmaFactorBasedOnVolatility(Volatility)

	// log.Printf("Volatility: %f", Volatility)
}

// ProcessOrderBook processes order book data to calculate liquidity and order book depth
func ProcessOrderBook(orderBook OrderBookData) {
	Mutex.Lock()
	defer Mutex.Unlock()

	// Calculate liquidity and order book depth
	var liquidityBid, liquidityAsk float64
	var depthBid, depthAsk int

	for _, bid := range orderBook.Bids {
		price := parseFloat(bid[0])
		quantity := parseFloat(bid[1])

		if price >= MidPrice*(1-liquidityRange) {
			liquidityBid += quantity
		}
		if price >= MidPrice*0.95 { // considering 5% below mid-price for depth
			depthBid++
		}
	}

	for _, ask := range orderBook.Asks {
		price := parseFloat(ask[0])
		quantity := parseFloat(ask[1])

		if price <= MidPrice*(1+liquidityRange) {
			liquidityAsk += quantity
		}
		if price <= MidPrice*1.05 { // considering 5% above mid-price for depth
			depthAsk++
		}
	}

	Liquidity = (liquidityBid + liquidityAsk) / 2
	OrderBookDepth = float64(depthBid+depthAsk) / 2

	// log.Printf("Liquidity: %f, OrderBookDepth: %f", Liquidity, OrderBookDepth)
}

func parseFloat(s string) float64 {
	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		log.Println("Error parsing float:", err)
		return 0
	}
	return val
}
