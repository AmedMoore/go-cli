package main

import (
	"fmt"
	"github.com/akaahmedkamal/go-cli/v1"
	"strings"
)

type Greeting struct {
	Name     string `cli:"name"`
	Help     string `cli:"help"`
	Alias    string `cli:"alias"`
	Username string `cli:"option" optName:"name" optAlias:"n" optHelp:"Username to say hi to"`
	AllCaps  bool   `cli:"option" optName:"all-caps" optHelp:"Print message in uppercase form"`
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

	msg := fmt.Sprintf("Hello, %s!\n", g.Username)

	if g.AllCaps {
		msg = strings.ToUpper(msg)
	}

	app.Log().Println(msg)
}
