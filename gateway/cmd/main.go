package main

import (
	"fmt"
	"log"
	"os"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
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
		return fmt.Errorf("app names API URL: %s", c.String(v3EndpointFlag))
	}
	if err := validation.Validate(c.String(writeAccessKeyFlag), validation.Required); err != nil {
		return fmt.Errorf("access key error: %s", err.Error())
	}

	if err := validation.Validate(c.String(writeSecretKeyFlag), validation.Required); err != nil {
		return fmt.Errorf("secret key error: %s", err.Error())
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
		return fmt.Errorf("authentication object creation error: %s", err)
	}
	perm, err := http.NewPermissioner(
		c.String(readAccessKeyFlag),
		c.String(writeAccessKeyFlag),
		c.String(confirmAccessKeyFlag),
		c.String(rebalanceAccessKeyFlag))
	if err != nil {
		return fmt.Errorf("permission object creation error: %s", err)
	}

	svr, err := http.NewServer(httputil.NewHTTPAddressFromContext(c),
		auth,
		perm,
		http.WithV3Endpoint(c.String(v3EndpointFlag)),
	)
	if err != nil {
		return err
	}
	return svr.Start()
}
