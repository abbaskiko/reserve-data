package main

import (
	"fmt"
	"log"
	"os"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/pkg/errors"
	"github.com/urfave/cli"

	"github.com/KyberNetwork/httpsign-utils/authenticator"
	"github.com/KyberNetwork/reserve-data/gateway/http"
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

	v3EndpointFlag = "v3-endpoint"
)

var (
	v3EndpointDefaultValue = fmt.Sprint("http://127.0.0.1:8000")
)

func main() {
	app := cli.NewApp()
	app.Name = "HTTP gateway for reserve core"
	app.Action = run
	app.Flags = append(app.Flags,
		cli.StringFlag{
			Name:   readAccessKeyFlag,
			Usage:  "key for access read paths",
			EnvVar: "READ_ACCESS_KEY",
		},
		cli.StringFlag{
			Name:   readSecretKeyFlag,
			Usage:  "secret key for read paths",
			EnvVar: "READ_SECRET_KEY",
		},
		cli.StringFlag{
			Name:   writeAccessKeyFlag,
			Usage:  "key for access write paths",
			EnvVar: "WRITE_ACCESS_KEY",
		},
		cli.StringFlag{
			Name:   writeSecretKeyFlag,
			Usage:  "secret key for access write paths",
			EnvVar: "WRITE_SECRET_KEY",
		},
		cli.StringFlag{
			Name:   confirmAccessKeyFlag,
			Usage:  "access key for access confirm paths",
			EnvVar: "CONFIRM_ACCESS_KEY",
		},
		cli.StringFlag{
			Name:   confirmSecretKeyFlag,
			Usage:  "secret key for access confirm paths",
			EnvVar: "CONFIRM_SECRET_KEY",
		},
		cli.StringFlag{
			Name:   rebalanceAccessKeyFlag,
			Usage:  "access key for access rebalance paths",
			EnvVar: "REBALANCE_ACCESS_KEY",
		},
		cli.StringFlag{
			Name:   rebalanceSecretKeyFlag,
			Usage:  "secret key to access rebalance paths",
			EnvVar: "REBALANCE_SECRET_KEY",
		},
		cli.StringFlag{
			Name:   v3EndpointFlag,
			Usage:  "v3 endpoint url",
			EnvVar: "V3_ENDPOINT",
			Value:  v3EndpointDefaultValue,
		},
	)

	app.Flags = append(app.Flags, httputil.NewHTTPCliFlags(httputil.GatewayPort)...)

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {

	err := validation.Validate(c.String(v3EndpointFlag),
		validation.Required,
		is.URL)
	if err != nil {
		return errors.Wrapf(err, "app names API URL error: %s", c.String(v3EndpointFlag))
	}
	if err := validation.Validate(c.String(writeAccessKeyFlag), validation.Required); err != nil {
		return errors.Wrap(err, "write access key error")
	}

	if err := validation.Validate(c.String(writeSecretKeyFlag), validation.Required); err != nil {
		return errors.Wrap(err, "secret key error")
	}
	keyPairs := []authenticator.KeyPair{
		{
			AccessKeyID:     c.String(readAccessKeyFlag),
			SecretAccessKey: c.String(readSecretKeyFlag),
		},
		{
			AccessKeyID:     c.String(writeAccessKeyFlag),
			SecretAccessKey: c.String(writeSecretKeyFlag),
		},
	}
	auth, err := authenticator.NewAuthenticator(keyPairs...)
	if err != nil {
		return errors.Wrap(err, "authentication object creation error")
	}
	perm, err := http.NewPermissioner(
		c.String(readAccessKeyFlag),
		c.String(writeAccessKeyFlag),
		c.String(confirmAccessKeyFlag),
		c.String(rebalanceAccessKeyFlag))
	if err != nil {
		return errors.Wrap(err, "permission object creation error")
	}

	svr, err := http.NewServer(httputil.NewHTTPAddressFromContext(c),
		auth,
		perm,
		http.WithV3Endpoint(c.String(v3EndpointFlag)),
	)
	if err != nil {
		return errors.Wrap(err, "create new server error")
	}
	return svr.Start()
}
