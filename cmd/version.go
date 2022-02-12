package cmd

import (
	"fmt"

	"github.com/akaahmedkamal/go-cli/v1"
)

type Version struct {
	Name string `cli:"name"`
	Help string `cli:"help"`
}

func NewVersionCmd() *Version {
	return &Version{
		Name: "version",
		Help: "print version and build time",
	}
}

func (h *Version) Run(app *cli.App) {
	fmt.Printf(
		"%s %s (build %s).\n",
		cli.AppName,
		cli.AppVersion,
		cli.AppBuild,
	)
}
