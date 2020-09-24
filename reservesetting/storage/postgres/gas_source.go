package postgres

import (
	"encoding/json"

	v3 "github.com/KyberNetwork/reserve-data/reservesetting/common"
)

const (
	keyPreferGasSource = "prefer-gas-source"
)

// GetPreferGasSource ...
func (s *Storage) GetPreferGasSource() (v3.PreferGasSource, error) {
	var result v3.PreferGasSource
	preferGasSourceData, err := s.GetGeneralData(keyPreferGasSource)
	if err != nil {
		return result, err
	}
	err = json.Unmarshal([]byte(preferGasSourceData.Value), &result)
	return result, err
}

// SetPreferGasSource ...
func (s *Storage) SetPreferGasSource(data v3.PreferGasSource) error {
	byteData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, err = s.SetGeneralData(v3.GeneralData{
		Key:   keyPreferGasSource,
		Value: string(byteData),
	})
	return err
}
