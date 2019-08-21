package profiler

import "github.com/urfave/cli"

const (
	pprofFlag = "enable-profiler"
)

func NewCliFlags() []cli.Flag {
	return []cli.Flag{
		cli.BoolFlag{
			Name:   pprofFlag,
			EnvVar: "ENABLE_PROFILER",
			Usage:  "enable profiler for debugging",
		},
	}
}

func IsEnableProfilerFromContext(c *cli.Context) bool {
	return c.Bool(pprofFlag)
}
