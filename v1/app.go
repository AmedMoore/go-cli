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

func (ref *App) GetCommand(name string) *CommandEntry {
	for _, cmd := range ref.cmds {
		if cmd.Name == name || cmd.Alias == name {
			return &cmd
		}
	}
	return nil
}

func (ref *App) LookupCommand(name string) (cmd *CommandEntry, exists bool) {
	cmd = ref.GetCommand(name)
	exists = cmd != nil
	return
}

func (ref *App) Register(cmd Command) {
	typ := reflect.TypeOf(cmd).Elem()
	val := reflect.ValueOf(cmd).Elem()

	var helpField reflect.Value
	var nameField reflect.Value
	var aliasField reflect.Value

	for i := 0; i < val.NumField(); i++ {
		fieldType := typ.Field(i)
		fieldValue := val.Field(i)
		cliTag := fieldType.Tag.Get("cli")

		if cliTag == "name" {
			nameField = fieldValue
		}

		if cliTag == "help" {
			helpField = fieldValue
		}

		if cliTag == "alias" {
			aliasField = fieldValue
		}
	}

	if helpField.Type().Name() != "string" {
		ref.LogOrDefault().Fatalln("Invalid command, command must have \"Help\" field")
	}

	if nameField.Type().Name() != "string" {
		ref.LogOrDefault().Fatalln("Invalid command, command must have \"Name\" field")
	}

	entry := CommandEntry{}
	entry.Help = helpField.String()
	entry.Name = nameField.String()
	entry.Alias = aliasField.String()
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

	cmd, exists := ref.LookupCommand(cmdName)
	if !exists {
		ref.LogOrDefault().Println("nothing to do!")
	} else {
		cmd.cmd.Run(ref)
	}
}

func (ref *App) LogOrDefault() *log.Logger {
	if ref.Log() != nil {
		return ref.Log().info
	} else {
		return log.Default()
	}
}
