package bybitconnector

// OrderBookSnapshot represents the full snapshot of the order book
type OrderBookSnapshot struct {
	Topic     string        `json:"topic"`
	Type      string        `json:"type"`
	Timestamp int64         `json:"ts"`
	Data      OrderBookData `json:"data"`
}

// OrderBookData represents the order book data within the snapshot
type OrderBookData struct {
	Symbol   string     `json:"s"`
	Bids     [][]string `json:"b"`
	Asks     [][]string `json:"a"`
	UpdateID int64      `json:"u"`
	Seq      int64      `json:"seq"`
}

// Trade represents a trade that has occurred
type TradeData struct {
	Topic     string `json:"topic"`
	Type      string `json:"type"`
	Timestamp int64  `json:"ts"`
	Data      []struct {
		TradeTimestamp int64  `json:"T"`
		Symbol         string `json:"s"`
		Direction      string `json:"S"`
		Volume         string `json:"v"`
		Price          string `json:"p"`
		Liquidation    bool   `json:"BT"`
	} `json:"data"`
}

// Ticker represents the ticker information
type Ticker struct {
	Topic     string     `json:"topic"`
	Type      string     `json:"type"`
	Data      TickerData `json:"data"`
	CheckSum  int64      `json:"cs"`
	Timestamp int64      `json:"ts"`
}

// TickerData represents the ticker data within the ticker snapshot
type TickerData struct {
	Symbol            string `json:"symbol"`
	TickDirection     string `json:"tickDirection"`
	Price24hPcnt      string `json:"price24hPcnt"`
	LastPrice         string `json:"lastPrice"`
	PrevPrice24h      string `json:"prevPrice24h"`
	HighPrice24h      string `json:"highPrice24h"`
	LowPrice24h       string `json:"lowPrice24h"`
	PrevPrice1h       string `json:"prevPrice1h"`
	MarkPrice         string `json:"markPrice"`
	IndexPrice        string `json:"indexPrice"`
	OpenInterest      string `json:"openInterest"`
	OpenInterestValue string `json:"openInterestValue"`
	Turnover24h       string `json:"turnover24h"`
	Volume24h         string `json:"volume24h"`
	NextFundingTime   string  `json:"nextFundingTime"`
	FundingRate       string `json:"fundingRate"`
	Bid1Price         string `json:"bid1Price"`
	Bid1Size          string `json:"bid1Size"`
	Ask1Price         string `json:"ask1Price"`
	Ask1Size          string `json:"ask1Size"`
}
