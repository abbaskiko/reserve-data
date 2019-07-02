package http

import (
	ethereum "github.com/ethereum/go-ethereum/common"

	"github.com/KyberNetwork/reserve-data/common"
)

type Setting interface {
	GetInternalTokenByID(tokenID string) (common.Token, error)
	GetActiveTokenByID(tokenID string) (common.Token, error)
	GetTokenByID(tokenID string) (common.Token, error)
	GetInternalTokens() ([]common.Token, error)
	GetAllTokens() ([]common.Token, error)
	NewTokenPairFromID(base, quote string) (common.TokenPair, error)
	GetFee(ex common.ExchangeName) (common.ExchangeFees, error)
	UpdateFee(ex common.ExchangeName, data common.ExchangeFees, timestamp uint64) error
	GetMinDeposit(ex common.ExchangeName) (common.ExchangesMinDeposit, error)
	UpdateMinDeposit(ex common.ExchangeName, minDeposit common.ExchangesMinDeposit, timestamp uint64) error
	GetDepositAddresses(ex common.ExchangeName) (common.ExchangeAddresses, error)
	UpdateDepositAddress(ex common.ExchangeName, addrs common.ExchangeAddresses, timestamp uint64) error
	GetExchangeInfo(ex common.ExchangeName) (common.ExchangeInfo, error)
	UpdateExchangeInfo(ex common.ExchangeName, exInfo common.ExchangeInfo, timestamp uint64) error
	GetExchangeStatus() (common.ExchangesStatus, error)
	UpdateExchangeStatus(data common.ExchangesStatus) error
	UpdatePendingTokenUpdates(map[string]common.TokenUpdate) error
	ApplyTokenWithExchangeSetting([]common.Token, map[common.ExchangeName]*common.ExchangeSetting, uint64) error
	GetPendingTokenUpdates() (map[string]common.TokenUpdate, error)
	RemovePendingTokenUpdates() error
	GetTokenVersion() (uint64, error)
	GetExchangeVersion() (uint64, error)
	GetActiveTokens() ([]common.Token, error)
	GetTokenByAddress(ethereum.Address) (common.Token, error)
}
