package exchange

// CoinbaseInterface contains the methods to interact with Coinbase centralized exchange.
type CoinbaseInterface interface {
	GetOnePairOrderBook(baseID, quoteID string) (CoinbaseResp, error)
}
