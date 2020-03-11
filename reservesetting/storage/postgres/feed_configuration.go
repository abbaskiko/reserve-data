package postgres

import (
	"database/sql"
	"fmt"

	pgutil "github.com/KyberNetwork/reserve-data/common/postgres"
	"github.com/KyberNetwork/reserve-data/reservesetting/common"
	"github.com/KyberNetwork/reserve-data/world"
	"github.com/jmoiron/sqlx"
)

type setFeedConfigurationParams struct {
	Name                 string   `db:"name"`
	SetRate              string   `db:"set_rate"`
	Enabled              *bool    `db:"enabled"`
	BaseVolatilitySpread *float64 `db:"base_volatility_spread"`
	NormalSpread         *float64 `db:"normal_spread"`
}

func (s *Storage) initFeedData() error {
	// init all feed as enabled
	query := `INSERT INTO "feed_configurations" (name, set_rate, enabled, base_volatility_spread, normal_spread) 
				VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING;`
	for feed := range world.AllFeeds().USD {
		if _, err := s.db.Exec(query, feed, common.USDFeed.String(), true, 0, 0); err != nil {
			return err
		}
	}
	for feed := range world.AllFeeds().BTC {
		if _, err := s.db.Exec(query, feed, common.BTCFeed.String(), true, 0, 0); err != nil {
			return err
		}
	}
	for feed := range world.AllFeeds().Gold {
		if _, err := s.db.Exec(query, feed, common.GoldFeed.String(), true, 0, 0); err != nil {
			return err
		}
	}
	return nil
}

// UpdateFeedStatus update feed status
func (s *Storage) UpdateFeedStatus(name string, setRate common.SetRate, enabled bool) error {
	return s.setFeedConfiguration(nil, common.SetFeedConfigurationEntry{
		Name:    name,
		Enabled: common.BoolPointer(enabled),
		SetRate: setRate,
	})
}

func (s *Storage) setFeedConfiguration(tx *sqlx.Tx, feedConfiguration common.SetFeedConfigurationEntry) error {
	var sts = s.stmts.setFeedConfiguration
	if tx != nil {
		sts = tx.NamedStmt(s.stmts.setFeedConfiguration)
	}
	var feedName string
	err := sts.Get(&feedName, setFeedConfigurationParams{
		Name:                 feedConfiguration.Name,
		SetRate:              feedConfiguration.SetRate.String(),
		Enabled:              feedConfiguration.Enabled,
		BaseVolatilitySpread: feedConfiguration.BaseVolatilitySpread,
		NormalSpread:         feedConfiguration.NormalSpread,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return common.ErrFeedConfiguration
		}
		return fmt.Errorf("failed to set feed config, err=%s,", err)
	}
	return nil
}

// GetFeedConfigurations return all feed configuration
func (s *Storage) GetFeedConfigurations() ([]common.FeedConfiguration, error) {
	var result []common.FeedConfiguration
	tx, err := s.db.Beginx()
	if err != nil {
		return nil, err
	}
	defer pgutil.RollbackUnlessCommitted(tx)

	if err := tx.Stmtx(s.stmts.getFeedConfigurations).Select(&result); err != nil {
		return nil, err
	}
	return result, nil
}

// GetFeedConfiguration return feed configuration by name
func (s *Storage) GetFeedConfiguration(name string, setRate common.SetRate) (common.FeedConfiguration, error) {
	var result common.FeedConfiguration
	tx, err := s.db.Beginx()
	if err != nil {
		return result, err
	}
	defer pgutil.RollbackUnlessCommitted(tx)

	if err := tx.Stmtx(s.stmts.getFeedConfiguration).Get(&result, name, setRate.String()); err != nil {
		if err == sql.ErrNoRows {
			return result, common.ErrNotFound
		}
		return result, err
	}
	return result, nil
}
