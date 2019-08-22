package postgres

import (
	"encoding/json"
	"time"

	"github.com/jmoiron/sqlx"
)

type stableTokenParamsDb struct {
	ID        uint64          `db:"id"`
	Timepoint time.Time       `db:"timepoint"`
	Data      json.RawMessage `db:"data"`
}

func (s *Storage) GetStableTokenParams() (map[string]interface{}, error) {
	var r stableTokenParamsDb
	var data = make(map[string]interface{})
	//var jsonMsg json.RawMessage

	err := s.stmts.getStableTokenParam.Get(&r)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(r.Data, &data)
	if err != nil {
		return nil, err
	}
	return data, err
}

func (s *Storage) updateStableTokenParams(tx *sqlx.Tx, data map[string]interface{}) error {
	jsonB, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, err = tx.Stmtx(s.stmts.newStableTokenParam).Exec(jsonB)
	return err
}
