package exchange

import (
	"math/big"

	ethereum "github.com/ethereum/go-ethereum/common"

	"github.com/KyberNetwork/reserve-data/common"
	commonv3 "github.com/KyberNetwork/reserve-data/v3/common"
)

// BinanceInterface contains the methods to interact with Binance centralized exchange.
type BinanceInterface interface {
	GetDepthOnePair(baseID, quoteID string) (Binaresp, error)

	OpenOrdersForOnePair(pair commonv3.TradingPairSymbols) (Binaorders, error)

	GetInfo() (Binainfo, error)

	GetExchangeInfo() (BinanceExchangeInfo, error)

	GetDepositAddress(tokenID string) (Binadepositaddress, error)

	GetAccountTradeHistory(baseSymbol, quoteSymbol, fromID string) (BinaAccountTradeHistory, error)

	Withdraw(
		token common.Token,
		amount *big.Int,
		address ethereum.Address) (string, error)

	Trade(
		tradeType string,
		base, quote common.Token,
		rate, amount float64) (Binatrade, error)

	CancelOrder(symbol string, id uint64) (Binacancel, error)

	DepositHistory(startTime, endTime uint64) (Binadeposits, error)

	WithdrawHistory(startTime, endTime uint64) (Binawithdrawals, error)

	OrderStatus(symbol string, id uint64) (Binaorder, error)
}
