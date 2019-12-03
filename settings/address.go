package settings

import (
	"errors"

	ethereum "github.com/ethereum/go-ethereum/common"
)

var (
	// ErrNoAddr error cannot find the address
	ErrNoAddr = errors.New("cannot find the address")
)

// GetAddress return address from address name
func (s *Settings) GetAddress(name AddressName) (ethereum.Address, error) {
	return s.Address.GetAddress(name)
}

// GetAddress return address from address name
func (addrSetting *AddressSetting) GetAddress(name AddressName) (ethereum.Address, error) {
	addr, ok := addrSetting.Addresses[name]
	if !ok {
		return ethereum.Address{}, ErrNoAddr
	}
	return addr, nil
}

// GetAllAddresses return all the address setting in cores.
func (s *Settings) GetAllAddresses() (map[string]interface{}, error) {
	allAddress := make(map[string]interface{})
	for name, addr := range s.Address.Addresses {
		allAddress[name.String()] = addr
	}
	return allAddress, nil
}
