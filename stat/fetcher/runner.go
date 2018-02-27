package fetcher

import (
	"time"
)

type FetcherRunner interface {
	GetBlockTicker() <-chan time.Time
	Start() error
	Stop() error
}
