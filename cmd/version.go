package cmd

import (
	"fmt"

	"github.com/akaahmedkamal/go-cli/v1"
)

type Version struct{}

func (h *Version) Name() string {
	return "version"
}

func (h *Version) Desc() string {
	return "display version and build time"
}

func (h *Version) Run(app *cli.App) {
	fmt.Printf(
		"%s %s (build %s).\n",
		cli.AppName,
		cli.AppVersion,
		cli.AppBuild,
	)
}
