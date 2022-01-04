package main

import (
	"fmt"

	"github.com/akaahmedkamal/go-cli/v1"
)

type Greeting struct{}

func (h *Greeting) Name() string {
	return "greeting"
}

func (h *Greeting) Desc() string {
	return "say hi!"
}

func (h *Greeting) Run(app *cli.App) {
	name, exists := app.Args().GetString("--name")
	if !exists {
		fmt.Println("please provide a name using --name")
		return
	}
	fmt.Printf("Hello, %s!\n", name)
}
