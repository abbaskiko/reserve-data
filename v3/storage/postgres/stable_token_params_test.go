package postgres

import (
	"fmt"
	"github.com/KyberNetwork/reserve-data/common/testutil"
	"github.com/KyberNetwork/reserve-data/v3/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestStorage_GetStableTokenParams(t *testing.T) {
	db, tearDown := testutil.MustNewDevelopmentDB()
	defer func() {
		assert.NoError(t, tearDown())
	}()

	s, err := NewStorage(db)
	assert.NoError(t, err)

	id, err := s.CreateSettingChange(common.ChangeCatalogStableToken, common.SettingChange{
		ChangeList: []common.SettingChangeEntry{
			{
				Type: common.ChangeTypeUpdateStableTokenParams,
				Data: common.UpdateStableTokenParamsEntry{
					Data: map[string]interface{}{
						"a": "1",
						"b": 2,
					},
				},
			},
		},
	})
	require.NoError(t, err)
	require.NoError(t, s.ConfirmSettingChange(id, true))
	data, err := s.GetStableTokenParams()
	require.NoError(t, err)
	fmt.Println(data)
}
