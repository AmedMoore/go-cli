package main

import (
	"github.com/akaahmedkamal/go-cli/v1"
)

type Greeting struct {
	Name     string `cli:"name"`
	Help     string `cli:"help"`
	Alias    string `cli:"alias"`
	Username string `cli:"option" optName:"name" optAlias:"n" optHelp:"Username to say hi to"`
}

func NewGreetingCmd() *Greeting {
	return &Greeting{
		Name:  "greeting",
		Help:  "Say hi!",
		Alias: "g",
	}
}

func (g *Greeting) Run(app *cli.App) {
	if g.Username == "" {
		app.Log().Fatalln("Option --name is missing")
	}
	app.Log().Printf("Hello, %s!\n", g.Username)
}
