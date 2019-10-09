package configuration

import (
	"path/filepath"

	"github.com/urfave/cli"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-data/common"
)

const secretConfigFileFlag = "secret-config"

// NewSecretConfigCliFlag returns the cli flag to configure secret config file flag.
func NewSecretConfigCliFlag() cli.Flag {
	return cli.StringFlag{
		Name:   secretConfigFileFlag,
		Usage:  "location of secret config file",
		EnvVar: "SECRET_CONFIG_FILE",
		Value:  filepath.Join(common.CmdDirLocation(), "config.json"),
	}
}

// NewSecretConfigFileFromContext returns the configured secret config file location.
func NewSecretConfigFileFromContext(c *cli.Context, l *zap.SugaredLogger) string {
	secretConfigFile := c.GlobalString(secretConfigFileFlag)
	l.Infof("using secret config file %s", secretConfigFile)
	return secretConfigFile
}
