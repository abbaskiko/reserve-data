package exchange

import (
	"github.com/KyberNetwork/reserve-data/common"
)

type Setting interface {
	GetInternalTokenByID(tokenID string) (common.Token, error)
	GetAllTokens() ([]common.Token, error)
	GetTokenByID(tokenID string) (common.Token, error)
	GetFee(ex common.ExchangeName) (common.ExchangeFees, error)
	GetMinDeposit(ex common.ExchangeName) (common.ExchangesMinDeposit, error)
	GetDepositAddresses(ex common.ExchangeName) (common.ExchangeAddresses, error)
	UpdateDepositAddress(name common.ExchangeName, addrs common.ExchangeAddresses, timestamp uint64) error
	GetExchangeInfo(ex common.ExchangeName) (common.ExchangeInfo, error)
	UpdateExchangeInfo(ex common.ExchangeName, exInfo common.ExchangeInfo, timestamp uint64) error
	ETHToken() common.Token
}
