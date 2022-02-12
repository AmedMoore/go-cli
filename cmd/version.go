package cmd

import (
	"fmt"

	"github.com/akaahmedkamal/go-cli/v1"
)

// Version command to print version and build time.
type Version struct {
	Name string `cli:"name"`
	Help string `cli:"help"`
}

// NewVersionCmd initializes new version command.
func NewVersionCmd() *Version {
	return &Version{
		Name: "version",
		Help: "print version and build time",
	}
}

// Run executes the command's logic.
func (h *Version) Run(_ *cli.App) {
	fmt.Printf(
		"%s %s (build %s).\n",
		cli.AppName,
		cli.AppVersion,
		cli.AppBuild,
	)
}
