package cli

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/akaahmedkamal/go-args"
)

var AppName string
var AppVersion string
var AppBuild string

type App struct {
	args   *args.ArgsParser
	cmds   []CommandEntry
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

	app.cmds = make([]CommandEntry, 0)
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

func (ref *App) Args() *args.ArgsParser {
	return ref.args
}

func (ref *App) Commands() []CommandEntry {
	return ref.cmds
}

func (ref *App) Command(name string) *CommandEntry {
	for _, cmd := range ref.cmds {
		if cmd.Name == name {
			return &cmd
		}
	}
	return nil
}

func (ref *App) Register(cmd Command) {
	t := reflect.TypeOf(cmd).Elem()
	helpField, helpFieldExist := t.FieldByName("help")
	if !helpFieldExist {
		ref.LogOrDefault().Fatalln("Invalid command, command must \"have\" help tag")
	}

	nameField, nameFieldExist := t.FieldByName("name")
	if !nameFieldExist {
		ref.LogOrDefault().Fatalln("Invalid command, command must \"name\" help tag")
	}

	aliasField, aliasFieldExist := t.FieldByName("alias")
	if !aliasFieldExist {
		ref.LogOrDefault().Fatalln("Invalid command, command must \"alias\" help tag")
	}

	entry := CommandEntry{}
	entry.Help = string(helpField.Tag)
	entry.Name = string(nameField.Tag)
	entry.Alias = string(aliasField.Tag)
	entry.cmd = cmd

	ref.cmds = append(ref.cmds, entry)
}

func (ref *App) Vars() map[string]interface{} {
	return ref.vars
}

func (ref *App) Set(name string, value interface{}) {
	ref.vars[name] = value
}

func (ref *App) Get(name string) interface{} {
	return ref.vars[name]
}

func (ref *App) ExeDir() string {
	return ref.exeDir
}

func (ref *App) Resolve(path ...string) string {
	return filepath.Join(append([]string{ref.exeDir}, path...)...)
}

func (ref *App) Log() *Logger {
	return ref.log
}

func (ref *App) SetLogger(logger *Logger) {
	ref.log = logger
}

func (ref *App) Mode() AppMode {
	return Mode
}

func (ref *App) Run() {
	cmdName := strings.Join(ref.args.Positional(), "/")

	cmd := ref.Command(cmdName)

	if cmd == nil {
		ref.LogOrDefault().Println("nothing to do!")
	}

	cmd.cmd.Run(ref)
}

func (ref *App) LogOrDefault() *log.Logger {
	if ref.Log() != nil {
		return ref.Log().info
	} else {
		return log.Default()
	}
}
