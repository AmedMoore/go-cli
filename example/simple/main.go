package main

import (
	"os"

	"github.com/ahmedmkamal/go-cli/cmd"
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

	// register greeting command
	app.Register(NewGreetingCmd())

	// start the app
	app.Run()
}
