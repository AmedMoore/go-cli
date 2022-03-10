package cli

import (
	"fmt"
	"strings"
)

const helpMessagePadding = "    "
const commandOptionAliasSeparator = ", "

// Help command to print help message.
type Help struct {
	Name string `cli:"name"`
	Help string `cli:"help"`
	app  *App
}

// NewHelpCmd initializes new help command.
func NewHelpCmd() *Help {
	return &Help{
		Name: "help",
		Help: "print app or command help message",
	}
}

// Run executes the command's logic.
func (h *Help) Run(app *App) {
	h.app = app

	args := app.Args().Positional()
	if len(args) > 0 && args[len(args)-1] == "help" {
		args = args[:len(args)-1]
	}

	cmdName := strings.Join(args, "/")
	cmd, exists := app.LookupCommand(cmdName)
	if cmdName == "" {
		h.printHelp()
	} else if exists {
		h.printHelpForCmd(*cmd)
	} else {
		h.printHelpForModule(cmdName)
	}
}

// longestCmdName returns the longest command name registered on cli.App
func (h *Help) longestCmdName() string {
	var name string
	for _, cmd := range h.app.Commands() {
		if len(cmd.Name) > len(name) {
			name = cmd.Name
		}
	}
	return name
}

// joinCommandNames returns command names and help messages of the
// given CommandEntry slice joined with a newline.
func (h *Help) joinCommandNames(commands []*CommandEntry) string {
	var names string
	longestCmdName := h.longestCmdName()
	for _, cmd := range commands {
		name := cmd.Name
		for len(name) < len(longestCmdName) {
			name += " "
		}
		names += fmt.Sprintf("  %s%s%s\n", name, helpMessagePadding, cmd.Help)
	}
	return names
}

// printHelp prints general help message of the cli.App
func (h *Help) printHelp() {
	commands := h.joinCommandNames(h.app.Commands())
	if commands != "" {
		fmt.Printf("Usage: %s [command] [options...]\n\nCommands:\n%s", AppName, commands)
	}
}

// printHelpForModule prints module/subcommands help message.
func (h *Help) printHelpForModule(name string) {
	commands := h.joinCommandNames(h.app.SubCommands(name))
	if commands != "" {
		fmt.Printf("Usage: %s %s [command] [options...]\n\nCommands:\n%s", AppName, name, commands)
	}
}

// longestOptionName returns the longest option name of Command
func (h *Help) longestOptionName(options []Option) string {
	var name string
	for _, opt := range options {
		nameAndAlias := opt.Name
		if opt.Alias != "" {
			nameAndAlias += commandOptionAliasSeparator + opt.Alias
		}
		if len(nameAndAlias) > len(name) {
			name = nameAndAlias
		}
	}
	return name
}

// joinOptionNames returns command names and help messages of the
// given CommandEntry slice joined with a newline.
func (h *Help) joinOptionNames(options []Option) string {
	var names string
	longestOptName := h.longestOptionName(options)
	for _, opt := range options {
		nameAndAlias := opt.Name
		if opt.Alias != "" {
			nameAndAlias += commandOptionAliasSeparator + opt.Alias
		}
		for len(nameAndAlias) < len(longestOptName) {
			nameAndAlias += " "
		}
		names += fmt.Sprintf("  %s%s%s\n", nameAndAlias, helpMessagePadding, opt.Help)
	}
	return names
}

// printHelpForCmd prints help message for Command
func (h *Help) printHelpForCmd(cmd CommandEntry) {
	optionNames := h.joinOptionNames(h.app.getCmdOptions(cmd.cmd))
	msg := "Usage: %s %s [options...]\n\n%s\n"
	if optionNames == "" {
		fmt.Printf(msg, AppName, cmd.Name, cmd.Help)
	} else {
		fmt.Printf(msg+"\nOptions:\n%s", AppName, cmd.Name, cmd.Help, optionNames)
	}
}
