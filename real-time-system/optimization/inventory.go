package optimization

// Inventory structure to represent the market maker's inventory
type Inventory struct {
	cashBalance   float64 // USD balance
	cryptoBalance float64 // BTC balance
	tradingFee    float64 // Trading fee as a percentage
}

// NewInventory creates and initializes a new Inventory object
func NewInventory(cashBalance, cryptoBalance, tradingFee float64) *Inventory {
	return &Inventory{
		cashBalance:   cashBalance,
		cryptoBalance: cryptoBalance,
		tradingFee:    tradingFee,
	}
}

// GetBalances returns the current cash and crypto balances
func (inv *Inventory) GetBalances() (float64, float64) {
	return inv.cashBalance, inv.cryptoBalance
}

// UpdateBalance updates the balances based on a trade
func (inv *Inventory) UpdateBalance(isBuy bool, quantity, price float64) {
	// Calculate the trading fee
	fee := inv.tradingFee * price * quantity / 100.0

	if isBuy {
		inv.cashBalance -= (price * quantity) + fee
		inv.cryptoBalance += quantity
	} else {
		inv.cashBalance += (price * quantity) - fee
		inv.cryptoBalance -= quantity
	}
}

// TradeExecuted executes a trade based on the current price and optimal bid and ask prices.
// It updates the inventory balances accordingly.
func (inv *Inventory) TradeExecuted(currentPrice, optimalBid, optimalAsk float64) {
	if currentPrice <= optimalBid {
		// Buy logic: Increase crypto balance, decrease cash balance
		amountToBuy := inv.cashBalance / currentPrice * (1 - inv.tradingFee)
		inv.cryptoBalance += amountToBuy
		inv.cashBalance -= amountToBuy * currentPrice
	} else if currentPrice >= optimalAsk {
		// Sell logic: Decrease crypto balance, increase cash balance
		amountToSell := inv.cryptoBalance * (1 - inv.tradingFee)
		inv.cashBalance += amountToSell * currentPrice
		inv.cryptoBalance -= amountToSell
	}
}
