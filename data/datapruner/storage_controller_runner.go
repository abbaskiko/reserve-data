package datapruner

import (
	"time"
)

// StorageControllerRunner is the controller interface of data pruner jobs.
type StorageControllerRunner interface {
	GetAuthBucketTicker() <-chan time.Time
	Start() error
}

type ControllerTickerRunner struct {
	authDuration time.Duration
	authClock    *time.Ticker
	signal       chan bool
}

func (c *ControllerTickerRunner) GetAuthBucketTicker() <-chan time.Time {
	if c.authClock == nil {
		<-c.signal
	}
	return c.authClock.C
}

func (c *ControllerTickerRunner) Start() error {
	c.authClock = time.NewTicker(c.authDuration)
	c.signal <- true
	return nil
}

func NewStorageControllerTickerRunner(
	authDuration time.Duration) *ControllerTickerRunner {
	return &ControllerTickerRunner{
		authDuration,
		nil,
		make(chan bool, 1),
	}
}
