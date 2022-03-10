package cli

import (
	"fmt"
	"os"
	"reflect"
)

// CommandEntry is the command entry registered on cli.App.
type CommandEntry struct {
	Name  string
	Alias string
	Help  string
	cmd   Command
}

// setName sets the CommandEntry name from the given Command field value.
func (c *CommandEntry) setName(val reflect.Value) {
	if !val.IsValid() && val.IsZero() {
		fmt.Println("invalid command, command must have \"Name\" field")
		os.Exit(1)
	}
	if val.Type().Name() != "string" {
		fmt.Println("invalid command, command \"Name\" must be string")
		os.Exit(1)
	}
	c.Name = val.String()
}

// setAlias sets the CommandEntry alias from the given Command field value.
func (c *CommandEntry) setAlias(val reflect.Value) {
	if val.IsValid() && !val.IsZero() {
		if val.Type().Name() != "string" {
			fmt.Println("invalid command, command \"Alias\" must be string")
			os.Exit(1)
		}
		c.Alias = val.String()
	}
}

// setHelp sets the CommandEntry help from the given Command field value.
func (c *CommandEntry) setHelp(val reflect.Value) {
	if !val.IsValid() && val.IsZero() {
		fmt.Println("invalid command, command must have \"Help\" field")
		os.Exit(1)
	}
	if val.Type().Name() != "string" {
		fmt.Println("invalid command, command \"Help\" must be string")
		os.Exit(1)
	}
	c.Help = val.String()
}
