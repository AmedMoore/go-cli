package cmd

import (
	"fmt"

	"github.com/akaahmedkamal/go-cli/v1"
)

type Help struct{}

func (h *Help) Name() string {
	return "help"
}

func (h *Help) Desc() string {
	return "display help message"
}

func (h *Help) Run(app *cli.App) {
	var commands string

	for _, cmd := range app.Commands() {
		commands += fmt.Sprintf("\t%s\t%s\n", cmd.Name(), cmd.Desc())
	}

	fmt.Printf(
		"Usage: %s [command] [options...]\n\ncommands:\n%s",
		cli.AppName,
		commands,
	)
}
