package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/skyareas/go-cli"
)

// Greeting command to print greeting message.
type Greeting struct {
	Name     string `cli:"name"`
	Help     string `cli:"help"`
	Alias    string `cli:"alias"`
	Username string `cli:"option" optName:"name" optAlias:"n" optHelp:"Username to say hi to"`
	AllCaps  bool   `cli:"option" optName:"all-caps" optHelp:"Print message in uppercase form"`
}

// NewGreetingCmd initializes new greeting command.
func NewGreetingCmd() *Greeting {
	return &Greeting{
		Name:  "greeting",
		Help:  "Say hi!",
		Alias: "g",
	}
}

// Run executes the command's logic.
func (g *Greeting) Run(_ *cli.App) {
	if g.Username == "" {
		fmt.Println("Option --name is missing")
		os.Exit(1)
	}
	if g.AllCaps {
		fmt.Println(strings.ToUpper(fmt.Sprintf("Hello, %s!\n", g.Username)))
	} else {
		fmt.Println(fmt.Sprintf("Hello, %s!\n", g.Username))
	}
}
