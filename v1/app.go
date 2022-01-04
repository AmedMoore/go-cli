package cli

import (
	"log"
	"os"
	"path/filepath"

	"github.com/akaahmedkamal/go-args"
)

var AppName string
var AppVersion string
var AppBuild string

type App struct {
	args *args.ArgsParser
	cmds []Command
}

func initAppInfo() {
	if AppName == "" {
		exe, err := os.Executable()
		if err != nil {
			panic(err)
		}
		AppName = filepath.Base(exe)
	}

	if AppVersion == "" {
		AppVersion = "Unspecified"
	}

	if AppBuild == "" {
		AppBuild = "Unspecified"
	}
}

func NewApp(rawArgs []string) *App {
	initAppInfo()

	args := args.NewParser(rawArgs)
	args.Parse()

	return &App{
		args: args,
		cmds: make([]Command, 0),
	}
}

func (a *App) Args() *args.ArgsParser {
	return a.args
}

func (a *App) Commands() []Command {
	return a.cmds
}

func (a *App) Command(name string) Command {
	for _, cmd := range a.cmds {
		if cmd.Name() == name {
			return cmd
		}
	}
	return nil
}

func (a *App) Register(cmd Command) {
	a.cmds = append(a.cmds, cmd)
}

func (a *App) Run() {
	cmdName, exists := a.args.At(0)
	cmd := a.Command(cmdName)

	if !exists || cmd == nil {
		log.Fatal("nothing to do!")
	}

	cmd.Run(a)
}
