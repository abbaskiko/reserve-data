package testutil

import (
	"go.uber.org/zap"
)

// NewExampleSugar return an example sugar log, use in testing purpose
func NewExampleSugar() *zap.SugaredLogger {
	l := zap.NewExample()
	return l.Sugar()
}
