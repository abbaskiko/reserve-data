package cmd

import (
	"runtime"

	"github.com/KyberNetwork/reserve-data/common/config"
	"github.com/KyberNetwork/reserve-data/data/storage"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	newDB string
)

func migrateDB(_ *cobra.Command, _ []string) {
	numCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPU)

	w := configLog(stdoutLog)
	s, f, err := newSugaredLogger(w)
	if err != nil {
		panic(err)
	}
	defer f()
	zap.ReplaceGlobals(s.Desugar())
	var ac = config.DefaultAppConfig()
	if err = config.LoadConfig(configFile, &ac); err != nil {
		s.Panicw("failed to load config file", "err", err, "file", configFile)
	}
	currentStorage, err := storage.NewBoltStorage(ac.DataDB)
	if err != nil {
		panic(err)
	}
	newStorage, err := storage.NewBoltStorage(newDB)
	if err != nil {
		panic(err)
	}

	if err := currentStorage.Migrate(newStorage); err != nil {
		s.Panicw("failed to migrate data", "err", err)
	}
	s.Infow("migrate db successfully", "to", newDB)
}
