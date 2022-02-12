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
	args        *args.Parser
	cmds        []CommandEntry
	vars        map[string]interface{}
	exeDir      string
	log         *Logger
	defaultHelp bool
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

func (a *App) Args() *args.Parser {
	return a.args
}

func (a *App) Commands() []CommandEntry {
	return a.cmds
}

func (a *App) GetCommand(name string) *CommandEntry {
	for _, cmd := range a.cmds {
		if cmd.Name == name || cmd.Alias == name {
			return &cmd
		}
	}
	return nil
}

func (a *App) LookupCommand(name string) (cmd *CommandEntry, exists bool) {
	cmd = a.GetCommand(name)
	exists = cmd != nil
	return
}

func (a *App) Register(cmd Command) {
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
		a.LogOrDefault().Fatalln("Invalid command, command must have \"Help\" field")
	}

	if nameField.Type().Name() != "string" {
		a.LogOrDefault().Fatalln("Invalid command, command must have \"Name\" field")
	}

	entry := CommandEntry{}
	entry.Help = helpField.String()
	entry.Name = nameField.String()
	entry.Alias = aliasField.String()
	entry.cmd = cmd

	a.cmds = append(a.cmds, entry)
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

	cmd, exists := a.LookupCommand(cmdName)

	if a.defaultHelp {
		if a.Args().HasOption("-h", "--help") {
			if cmdName == "" {
				a.PrintHelp()
			} else if exists {
				a.PrintHelpForCmd(*cmd)
			} else {
				a.PrintHelpForModule(cmdName)
			}
			return
		}
	}

	if !exists {
		a.LogOrDefault().Println("nothing to do!")
	} else {
		a.assignCmdOptions(cmd.cmd)
		cmd.cmd.Run(a)
	}
}

func (a *App) LogOrDefault() *log.Logger {
	if a.Log() != nil {
		return a.Log().info
	} else {
		return log.Default()
	}
}

func (a *App) PrintHelp() {
	var longestCmdName string
	for _, cmd := range a.Commands() {
		if len(cmd.Name) > len(longestCmdName) {
			longestCmdName = cmd.Name
		}
	}

	cmdNamePadding := "    "

	var commands string
	for _, cmd := range a.Commands() {
		cmdName := cmd.Name
		for len(cmdName) < len(longestCmdName) {
			cmdName += " "
		}
		commands += fmt.Sprintf("  %s%s%s\n", cmdName, cmdNamePadding, cmd.Help)
	}

	msg := "Usage: %s [command] [options...]\n\nCommands:\n%s"

	if commands != "" {
		fmt.Printf(msg, AppName, commands)
	}
}

func (a *App) PrintHelpForModule(name string) {
	var longestCmdName string
	for _, cmd := range a.Commands() {
		if len(cmd.Name) > len(longestCmdName) {
			longestCmdName = cmd.Name
		}
	}

	cmdNamePadding := "    "

	var commands string
	for _, cmd := range a.Commands() {
		if strings.HasPrefix(cmd.Name, name) {
			cmdName := cmd.Name
			for len(cmdName) < len(longestCmdName) {
				cmdName += " "
			}
			commands += fmt.Sprintf("  %s%s%s\n", cmdName, cmdNamePadding, cmd.Help)
		}
	}

	msg := "Usage: %s %s [command] [options...]\n\nCommands:\n%s"

	if commands != "" {
		fmt.Printf(msg, AppName, name, commands)
	}
}

func (a *App) PrintHelpForCmd(cmd CommandEntry) {
	options := a.getCmdOptions(cmd.cmd)

	var longestOptName string
	for _, opt := range options {
		if len(opt.Name) > len(longestOptName) {
			longestOptName = opt.Name
		}
	}

	optNamePadding := "    "

	var optionsStr string
	for _, opt := range options {
		optName := opt.Name
		for len(optName) < len(longestOptName) {
			optName += " "
		}
		optionsStr += fmt.Sprintf("  %s%s%s\n", optName, optNamePadding, opt.Help)
	}

	msg := "Usage: %s %s [options...]\n\n%s\n\nOptions:\n%s"

	if optionsStr == "" {
		optionsStr = "None"
	}

	fmt.Printf(msg, AppName, cmd.Name, cmd.Help, optionsStr)
}

func (a *App) RegisterDefaultHelp() {
	a.defaultHelp = true
}

func (a *App) getCmdOptions(cmd Command) []Option {
	typ := reflect.TypeOf(cmd).Elem()

	options := make([]Option, 0)

	for i := 0; i < typ.NumField(); i++ {
		fieldType := typ.Field(i)
		cliTag := fieldType.Tag.Get("cli")

		if cliTag == "option" {
			opt := Option{}
			opt.Name = fieldType.Tag.Get("optName")
			opt.Alias = fieldType.Tag.Get("optAlias")
			opt.Help = fieldType.Tag.Get("optHelp")
			options = append(options, opt)
		}
	}

	return options
}

func (a *App) assignCmdOptions(cmd Command) {
	typ := reflect.TypeOf(cmd).Elem()
	val := reflect.ValueOf(cmd).Elem()

	for i := 0; i < typ.NumField(); i++ {
		fieldType := typ.Field(i)
		cliTag := fieldType.Tag.Get("cli")

		if cliTag == "option" {
			fieldValue := val.Field(i)

			optName := fieldType.Tag.Get("optName")
			optAlias := fieldType.Tag.Get("optAlias")

			optVal, _ := a.Args().GetString("--"+optName, "-"+optAlias)

			fieldValue.SetString(optVal)
		}
	}
}
