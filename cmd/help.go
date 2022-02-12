package cmd

import (
	"fmt"

	"github.com/akaahmedkamal/go-cli/v1"
)

// Help command to print help message.
type Help struct {
	Name string `cli:"name"`
	Help string `cli:"help"`
}

func NewHelpCmd() *Help {
	return &Help{
		Name: "help",
		Help: "print help message",
	}
}

func (h *Help) Run(app *cli.App) {
	var commands string

	for _, cmd := range app.Commands() {
		commands += fmt.Sprintf("\t%s\t%s\n", cmd.Name, cmd.Help)
	}

	fmt.Printf(
		"Usage: %s [command] [options...]\n\ncommands:\n%s",
		cli.AppName,
		commands,
	)
}
