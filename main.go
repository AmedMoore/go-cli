package main

import (
	"os"

	"github.com/akaahmedkamal/go-cli/cmd"
	"github.com/akaahmedkamal/go-cli/v1"
)

func main() {
	// create app instance
	app := cli.NewApp(os.Args[1:])

	// register commands
	app.Register(&cmd.Help{})
	app.Register(&cmd.Version{})

	// start the app
	app.Run()
}
