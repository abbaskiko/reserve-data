package deployment

import (
	"fmt"

	"github.com/urfave/cli"
)

const deploymentFlag = "deployment"

// Deployment is a deployment of Kyber Network contract, can be in different networks.
//go:generate stringer -type=Deployment -linecomment
type Deployment int

const (
	// TODO: should we remove Development?
	// Development is the development deployment of Kyber Network contracts on mainnet.
	Development Deployment = iota // develop
	// Production is the production deployment of Kyber Network contracts on mainnet.
	Production // production
	// Simulation is the special deployment of Kyber Network contracts that supports simulation testing.
	Simulation // simulation
)

// NewCliFlag returns new cli flag from context.
func NewCliFlag() cli.Flag {
	return cli.StringFlag{
		Name:   deploymentFlag,
		Usage:  "Kyber Network deployment name",
		EnvVar: "DEPLOYMENT",
		Value:  Development.String(),
	}
}

// NewDeploymentFromContext returns the running mode from context.
func NewDeploymentFromContext(c *cli.Context) (Deployment, error) {
	deploymentStr := c.GlobalString(deploymentFlag)
	switch deploymentStr {
	case Development.String():
		return Development, nil
	case Production.String():
		return Production, nil
	case Simulation.String():
		return Simulation, nil
	}

	return 0, fmt.Errorf("unknown deployment %s", deploymentStr)
}
