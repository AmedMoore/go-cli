package cli

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/akaahmedkamal/go-args"
)

var AppName string
var AppVersion string
var AppBuild string

type App struct {
	args   *args.ArgsParser
	cmds   []Command
	vars   map[string]interface{}
	exeDir string
	log    *Logger
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

func NewApp(rawArgs []string, config ...Config) *App {
	initAppInfo()

	app := new(App)

	// initialize args parser
	app.args = args.NewParser(rawArgs)
	if err := app.args.Parse(); err != nil {
		log.Fatal(err)
	}

	// set the app exeDir path
	exe, err := os.Executable()
	if err != nil {
		panic(err)
	}
	app.exeDir = filepath.Dir(exe)

	app.cmds = make([]Command, 0)
	app.vars = make(map[string]interface{})

	if len(config) > 0 {
		cfg := config[0]
		if cfg.Logger != nil {
			app.log = cfg.Logger
		}
	}

	if app.log == nil {
		if Mode == AppModeRelease {
			fname := filepath.Join(app.exeDir, fmt.Sprintf("%s.log", AppName))
			app.log = NewFileLogger(fname)
		} else {
			app.log = NewStdLogger()
		}
	}

	return app
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

func (a *App) Vars() map[string]interface{} {
	return a.vars
}

func (a *App) Set(name string, value interface{}) {
	a.vars[name] = value
}

func (a *App) Get(name string) interface{} {
	return a.vars[name]
}

func (a *App) ExeDir() string {
	return a.exeDir
}

func (a *App) Resolve(path ...string) string {
	return filepath.Join(append([]string{a.exeDir}, path...)...)
}

func (a *App) Log() *Logger {
	return a.log
}

func (a *App) SetLogger(logger *Logger) {
	a.log = logger
}

func (a *App) Mode() AppMode {
	return Mode
}

func (a *App) Run() {
	cmdName := strings.Join(a.args.Positional(), "/")

	cmd := a.Command(cmdName)

	if cmd == nil {
		log.Fatal("nothing to do!")
	}

	cmd.Run(a)
}
