package core

import (
	ethereum "github.com/ethereum/go-ethereum/common"

	"github.com/KyberNetwork/reserve-data/settings"
)

type Setting interface {
	GetAddress(settings.AddressName) (ethereum.Address, error)
}
