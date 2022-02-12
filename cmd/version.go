package cmd

import (
	"fmt"

	"github.com/akaahmedkamal/go-cli/v1"
)

type Version struct {
	Name  string
	Help  string
	Alias string
}

func NewVersionCmd() *Version {
	return &Version{
		Name:  "version",
		Help:  "print version and build time",
		Alias: "v",
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
