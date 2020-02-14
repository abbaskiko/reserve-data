package cmd

import (
	"runtime"

	"github.com/KyberNetwork/reserve-data/data/storage"
	"github.com/spf13/cobra"
)

var (
	currentDB, newDB string
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
	currentStorage, err := storage.NewBoltStorage(currentDB)
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
	s.Infow("migrate db successfully", "from", currentDB, "to", newDB)
}
