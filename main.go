package main

import (
	"os"

	"github.com/ahmedmkamal/go-cli/cmd"
	"github.com/ahmedmkamal/go-cli/v1"
)

func main() {
	// create app instance
	app := cli.NewApp(os.Args[1:])

	// use standard logger by default
	app.SetLogger(cli.NewStdLogger())

	// register default help command
	app.RegisterDefaultHelp()

	// register version command
	app.Register(cmd.NewVersionCmd())

	// start the app
	app.Run()
}
