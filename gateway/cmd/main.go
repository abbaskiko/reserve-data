package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-contrib/httpsign"
	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/pkg/errors"
	"github.com/urfave/cli"

	"github.com/KyberNetwork/httpsign-utils/authenticator"
	"github.com/KyberNetwork/reserve-data/cmd/mode"
	"github.com/KyberNetwork/reserve-data/gateway/http"
	libapp "github.com/KyberNetwork/reserve-data/lib/app"
	"github.com/KyberNetwork/reserve-data/lib/httputil"
)

const (
	readAccessKeyFlag      = "read-access-key"
	readSecretKeyFlag      = "read-secret-key"
	writeAccessKeyFlag     = "write-access-key"
	writeSecretKeyFlag     = "write-secret-key"
	confirmAccessKeyFlag   = "confirm-access-key"
	confirmSecretKeyFlag   = "confirm-secret-key"
	rebalanceAccessKeyFlag = "rebalance-access-key"
	rebalanceSecretKeyFlag = "rebalance-secret-key"

	coreEndpointFlag    = "core-endpoint"
	settingEndpointFlag = "setting-endpoint"

	noAuthFlag = "no-auth"
)

var (
	coreEndpointDefaultValue    = fmt.Sprint("http://127.0.0.1:8000")
	settingEndpointDefaultValue = fmt.Sprint("http://127.0.0.1:8002")
)

func main() {
	app := cli.NewApp()
	app.Name = "HTTP gateway for reserve core"
	app.Action = run
	app.Flags = append(app.Flags,
		cli.StringSliceFlag{
			Name:   readAccessKeyFlag,
			Usage:  "key for access read paths",
			EnvVar: "READ_ACCESS_KEY",
		},
		cli.StringSliceFlag{
			Name:   readSecretKeyFlag,
			Usage:  "secret key for read paths",
			EnvVar: "READ_SECRET_KEY",
		},
		cli.StringSliceFlag{
			Name:   writeAccessKeyFlag,
			Usage:  "key for access write paths",
			EnvVar: "WRITE_ACCESS_KEY",
		},
		cli.StringSliceFlag{
			Name:   writeSecretKeyFlag,
			Usage:  "secret key for access write paths",
			EnvVar: "WRITE_SECRET_KEY",
		},
		cli.StringSliceFlag{
			Name:   confirmAccessKeyFlag,
			Usage:  "access key for access confirm paths",
			EnvVar: "CONFIRM_ACCESS_KEY",
		},
		cli.StringSliceFlag{
			Name:   confirmSecretKeyFlag,
			Usage:  "secret key for access confirm paths",
			EnvVar: "CONFIRM_SECRET_KEY",
		},
		cli.StringSliceFlag{
			Name:   rebalanceAccessKeyFlag,
			Usage:  "access key for access rebalance paths",
			EnvVar: "REBALANCE_ACCESS_KEY",
		},
		cli.StringSliceFlag{
			Name:   rebalanceSecretKeyFlag,
			Usage:  "secret key to access rebalance paths",
			EnvVar: "REBALANCE_SECRET_KEY",
		},
		cli.StringFlag{
			Name:   coreEndpointFlag,
			Usage:  "core endpoint url",
			EnvVar: "CORE_ENDPOINT",
			Value:  coreEndpointDefaultValue,
		},
		cli.StringFlag{
			Name:   settingEndpointFlag,
			Usage:  "setting endpoint url",
			EnvVar: "SETTING_ENDPOINT",
			Value:  settingEndpointDefaultValue,
		},
		cli.BoolFlag{
			Name:   noAuthFlag,
			Usage:  "no authenticate",
			EnvVar: "NO_AUTH",
		},
	)

	app.Flags = append(app.Flags, httputil.NewHTTPCliFlags(httputil.GatewayPort)...)
	app.Flags = append(app.Flags, mode.NewCliFlag())

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func getKeyList(c *cli.Context, accessKeyFlag, secretKeyFlag string) ([]authenticator.KeyPair, error) {
	var keyPairs []authenticator.KeyPair
	accessKeys := c.StringSlice(accessKeyFlag)
	secretKeys := c.StringSlice(secretKeyFlag)
	if len(accessKeys) != len(secretKeys) {
		return nil, errors.Errorf("length read access keys (%d) and read secret keys (%d)", len(accessKeys), len(secretKeys))
	}

	for index := range accessKeys {
		keyPairs = append(keyPairs, authenticator.KeyPair{
			AccessKeyID:     accessKeys[index],
			SecretAccessKey: secretKeys[index],
		})
	}
	return keyPairs, nil
}

func run(c *cli.Context) error {
	var (
		keyPairs, readKeys, writeKeys, confirmKeys, rebalanceKeys []authenticator.KeyPair
		err                                                       error
		auth                                                      *httpsign.Authenticator
		perm                                                      gin.HandlerFunc
	)
	logger, err := libapp.NewLogger(c)
	if err != nil {
		return err
	}
	defer libapp.NewFlusher(logger)()
	if err := validation.Validate(c.String(coreEndpointFlag),
		validation.Required,
		is.URL); err != nil {
		return errors.Wrapf(err, "app names API URL error: %s", c.String(coreEndpointFlag))
	}

	if err := validation.Validate(c.String(settingEndpointFlag),
		validation.Required,
		is.URL); err != nil {
		return errors.Wrapf(err, "app names API URL error: %s", c.String(settingEndpointFlag))
	}

	noAuth := c.Bool(noAuthFlag)
	if !noAuth {
		readKeys, err = getKeyList(c, readAccessKeyFlag, readSecretKeyFlag)
		if err != nil {
			return errors.Wrap(err, "failed to get read keys")
		}
		keyPairs = append(keyPairs, readKeys...)
		writeKeys, err = getKeyList(c, writeAccessKeyFlag, writeSecretKeyFlag)
		if err != nil {
			return errors.Wrap(err, "failed to get write keys")
		}
		keyPairs = append(keyPairs, writeKeys...)
		confirmKeys, err = getKeyList(c, confirmAccessKeyFlag, confirmSecretKeyFlag)
		if err != nil {
			return errors.Wrap(err, "failed to get confirm keys")
		}
		keyPairs = append(keyPairs, confirmKeys...)
		rebalanceKeys, err = getKeyList(c, rebalanceAccessKeyFlag, rebalanceSecretKeyFlag)
		if err != nil {
			return errors.Wrap(err, "failed to get rebalance keys")
		}
		keyPairs = append(keyPairs, rebalanceKeys...)

		if err := validation.Validate(c.String(writeAccessKeyFlag), validation.Required); err != nil {
			return errors.Wrap(err, "write access key error")
		}

		if err := validation.Validate(c.String(writeSecretKeyFlag), validation.Required); err != nil {
			return errors.Wrap(err, "secret key error")
		}
		auth, err = authenticator.NewAuthenticator(keyPairs...)
		if err != nil {
			return errors.Wrap(err, "authentication object creation error")
		}
		perm, err = http.NewPermissioner(readKeys, writeKeys, confirmKeys, rebalanceKeys)
		if err != nil {
			return errors.Wrap(err, "permission object creation error")
		}
	}

	svr, err := http.NewServer(httputil.NewHTTPAddressFromContext(c),
		auth,
		perm,
		noAuth,
		logger,
		http.WithCoreEndpoint(c.String(coreEndpointFlag)),
		http.WithSettingEndpoint(c.String(settingEndpointFlag)),
	)
	if err != nil {
		return errors.Wrap(err, "create new server error")
	}
	return svr.Start()
}
