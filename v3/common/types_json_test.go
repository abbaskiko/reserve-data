package common

import (
	"encoding/json"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestSettingChangeEntry_UnmarshalJSON(t *testing.T) {

	var settingChange = SettingChange{
		ChangeList: []SettingChangeEntry{
			{
				Type: ChangeTypeChangeAssetAddr,
				Data: &ChangeAssetAddressEntry{
					ID:      1,
					Address: common.HexToAddress("0x000001"),
				},
			},
		},
	}

	data, err := json.Marshal(settingChange)
	require.NoError(t, err)
	var newSettingChange SettingChange
	err = json.Unmarshal(data, &newSettingChange)
	require.NoError(t, err)
	require.Equal(t, newSettingChange, settingChange)

}
