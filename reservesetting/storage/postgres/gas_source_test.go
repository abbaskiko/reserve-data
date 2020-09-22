package postgres

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/KyberNetwork/reserve-data/common/testutil"
	"github.com/KyberNetwork/reserve-data/reservesetting/common"
)

func TestSetGeneralData(t *testing.T) {
	db, teardown := testutil.MustNewDevelopmentDB(migrationPath)
	defer func() {
		assert.NoError(t, teardown())
	}()

	ps, err := NewStorage(db)
	assert.NoError(t, err)

	data := common.PreferGasSource{
		Name: "etherscan",
	}
	err = ps.SetPreferGasSource(data)
	assert.NoError(t, err)
	pgs, err := ps.GetPreferGasSource()
	assert.NoError(t, err)
	assert.Equal(t, data, pgs)
}
