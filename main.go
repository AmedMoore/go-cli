package main

import (
	"os"

	"github.com/akaahmedkamal/go-cli/cmd"
	"github.com/akaahmedkamal/go-cli/v1"
)

func main() {
	// create app instance
	app := cli.NewApp(os.Args[1:])

	// use standard logger by default
	app.SetLogger(cli.NewStdLogger())

	// register commands
	app.Register(cmd.NewHelpCmd())
	app.Register(cmd.NewVersionCmd())

	// start the app
	app.Run()
}
