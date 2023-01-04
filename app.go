package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/ahmedmkamal/go-args"
)

// AppName is the name of the App. Defaults to the binary file name.
// This can be set at build time using
//
//	-ldflags -X github.com/ahmedmkamal/go-cli.AppName=[NAME]
var AppName string

// AppVersion is the version of the App. Defaults to 0.1.0.
// This can be set at build time using
//
//	-ldflags -X github.com/ahmedmkamal/go-cli.AppVersion=[VERSION]
var AppVersion string

// AppBuild is the build time of the App. Defaults to modified time of the binary file.
// This can be set at build time using
//
//	-ldflags -X github.com/ahmedmkamal/go-cli.AppBuild=[BUILD]
var AppBuild string

// cmdOptionNamePrefix is the prefix used for the cli option name (i.e. --help).
var cmdOptionNamePrefix = "--"

// cmdOptionAliasPrefix is the prefix used for the cli option alias (i.e. -h).
var cmdOptionAliasPrefix = "-"

// App is the main structure of a cli application.
// Do NOT create new instance of the App directly,
// use cli.NewApp() instead.
type App struct {
	args       *args.Parser
	cfg        *Config
	commands   []*CommandEntry
	defaultCmd *CommandEntry
	exe        string
	exeDir     string
}

// NewApp creates a new cli App with some defaults
// for name, version, build time, and help command,
// optionally accepting Config.
func NewApp(config ...*Config) *App {
	app := &App{}
	app.setExePath()
	app.setConfig(config...)
	app.setName()
	app.setVersion()
	app.setBuildTime()
	app.setArgsParser()
	app.commands = make([]*CommandEntry, 0)
	app.Register(NewHelpCmd())
	app.Register(NewVersionCmd())
	app.SetDefaultCommand("help")
	return app
}

// setExePath sets the app executable path and executable dir path.
func (a *App) setExePath() {
	exe, err := os.Executable()
	if err != nil {
		panic(err)
	}
	a.exe = exe
	a.exeDir = filepath.Dir(exe)
}

// setConfig sets the app configuration, if no config passed
// a default empty config struct will be assigned.
func (a *App) setConfig(config ...*Config) {
	if len(config) > 0 && config[0] != nil {
		a.cfg = config[0]
	} else {
		a.cfg = &Config{}
	}
}

// setName sets the app name from config, if empty, then
// the AppName will be used, if empty, defaults to the
// executable/binary file name.
func (a *App) setName() {
	if a.cfg.Name != "" {
		AppName = a.cfg.Name
		return
	}
	if AppName != "" {
		a.cfg.Name = AppName
		return
	}
	AppName = filepath.Base(a.exe)
	a.cfg.Name = AppName
}

// setVersion sets the app version from config, if empty, then
// the AppVersion will be used, if empty, defaults to 0.1.0.
func (a *App) setVersion() {
	if a.cfg.Version != "" {
		AppVersion = a.cfg.Version
		return
	}
	if AppVersion != "" {
		a.cfg.Version = AppVersion
		return
	}
	AppVersion = "0.1.0"
	a.cfg.Version = AppVersion
}

// setBuildTime sets the app build time from config, if empty, then
// the AppBuild will be used, if empty, try to use the modified time
// of the executable/binary file, in case of an error, defaults to
// the current time.
func (a *App) setBuildTime() {
	if a.cfg.BuildTime != "" {
		AppBuild = a.cfg.BuildTime
		return
	}
	if AppBuild != "" {
		a.cfg.BuildTime = AppBuild
		return
	}
	info, err := os.Stat(a.exe)
	if err != nil {
		AppBuild = time.Now().Format("2006-01-02 15:04")
	}
	AppBuild = info.ModTime().Format("2006-01-02 15:04")
	a.cfg.BuildTime = AppBuild
}

// setArgsParser sets the app argument parser.
func (a *App) setArgsParser() {
	a.args = args.NewParser(os.Args[1:])
}

// Args returns the argument parser of the App.
func (a *App) Args() *args.Parser {
	return a.args
}

// Commands returns all commands registered on App.
func (a *App) Commands() []*CommandEntry {
	return a.commands
}

// SubCommands returns all subcommands registered on App.
func (a *App) SubCommands(name string) []*CommandEntry {
	commands := make([]*CommandEntry, 0)
	for _, cmd := range a.commands {
		if strings.HasPrefix(cmd.Name, name) {
			commands = append(commands, cmd)
		}
	}
	return commands
}

// Command returns the named command on App. Returns nil if the command does not exist.
func (a *App) Command(name string) *CommandEntry {
	if name == "" {
		return nil
	}
	for _, c := range a.commands {
		if c.Name == name || c.Alias == name {
			return c
		}
	}
	return nil
}

// LookupCommand returns the named command on App and a bool indicates whether the command exists.
func (a *App) LookupCommand(name string) (*CommandEntry, bool) {
	c := a.Command(name)
	return c, c != nil
}

// ensureUnique looks up the command entry on App and exit with error if found.
func (a *App) ensureUnique(entry *CommandEntry) {
	if _, exists := a.LookupCommand(entry.Name); exists {
		a.exitWithError("command with name \"%s\" already exists\n", entry.Name)
	}
	if _, exists := a.LookupCommand(entry.Alias); entry.Alias != "" && exists {
		a.exitWithError("command with alias \"%s\" already exists\n", entry.Alias)
	}
}

// Register registers a new command on App. command
// name and alias must be unique, app will exit with
// error when attempt to register command that already exists.
func (a *App) Register(cmd Command) *App {
	entry := a.makeCommandEntry(cmd)
	a.ensureUnique(entry)
	a.commands = append(a.commands, entry)
	return a
}

// RegisterDefault registers the default command
// that gets executed when the command name passed
// to the CLI doesn't match any registered command.
func (a *App) RegisterDefault(cmd Command) *App {
	entry := a.makeCommandEntry(cmd)
	a.ensureUnique(entry)
	a.defaultCmd = entry
	return a
}

// SetDefaultCommand sets the default command to the
// command registered on App with the given name.
func (a *App) SetDefaultCommand(name string) *App {
	cmd, exists := a.LookupCommand(name)
	if !exists {
		a.exitWithError("command named \"%s\" not found\n", name)
	}
	a.defaultCmd = cmd
	return a
}

// Exe returns the App executable path.
func (a *App) Exe() string {
	return a.exe
}

// ExeDir returns the App executable directory path.
func (a *App) ExeDir() string {
	return a.exeDir
}

// Mode returns the App mode, also accessible from cli.Mode directly.
func (a *App) Mode() AppMode {
	return Mode
}

// Run is the entry point to the cli app. Parses the cli
// arguments and executes the proper command.
func (a *App) Run() {
	if err := a.args.Parse(); err != nil {
		a.exitWithError(err.Error())
	}

	cliArgs := a.args.Positional()

	var cmd *CommandEntry
	var exists bool

	if len(cliArgs) != 0 && cliArgs[len(cliArgs)-1] == "help" || a.args.HasOption("--help") {
		cmd, exists = a.LookupCommand("help")
	} else {
		cmdName := strings.Join(cliArgs, "/")
		cmd, exists = a.LookupCommand(cmdName)
	}

	if exists {
		a.assignCmdOptions(cmd.cmd)
		cmd.cmd.Run(a)
		return
	}

	if a.defaultCmd != nil {
		a.assignCmdOptions(a.defaultCmd.cmd)
		a.defaultCmd.cmd.Run(a)
		return
	}

	fmt.Println("nothing to do!")
}

// makeCommandEntry creates a new command entry for the given command.
func (a *App) makeCommandEntry(cmd Command) *CommandEntry {
	typ := reflect.TypeOf(cmd).Elem()
	val := reflect.ValueOf(cmd).Elem()

	var nameField reflect.Value
	var aliasField reflect.Value
	var helpField reflect.Value

	for i := 0; i < val.NumField(); i++ {
		fieldType := typ.Field(i)
		cliTag := fieldType.Tag.Get("cli")
		fieldValue := val.Field(i)

		switch cliTag {
		case "name":
			nameField = fieldValue
		case "alias":
			aliasField = fieldValue
		case "help":
			helpField = fieldValue
		}
	}

	entry := &CommandEntry{}
	entry.setName(nameField)
	entry.setAlias(aliasField)
	entry.setHelp(helpField)
	entry.cmd = cmd
	return entry
}

// getCmdOptions gets option fields info for the given command.
func (a *App) getCmdOptions(cmd Command) []Option {
	options := make([]Option, 0)
	typ := reflect.TypeOf(cmd).Elem()
	for i := 0; i < typ.NumField(); i++ {
		fieldType := typ.Field(i)
		cliTag := fieldType.Tag.Get("cli")
		if cliTag == "option" {
			opt := Option{}
			opt.setName(fieldType)
			opt.setAlias(fieldType)
			opt.setHelp(fieldType)
			options = append(options, opt)
		}
	}
	return options
}

// assignCmdOptions assigns values to the option fields of the given command.
func (a *App) assignCmdOptions(cmd Command) {
	typ := reflect.TypeOf(cmd).Elem()
	val := reflect.ValueOf(cmd).Elem()
	for i := 0; i < typ.NumField(); i++ {
		fieldType := typ.Field(i)
		if fieldType.Tag.Get("cli") == "option" {
			fieldValue := val.Field(i)
			a.setCommandOptionValue(&fieldType, &fieldValue)
		}
	}
}

// setCommandOptionValue assigns the proper value to command option field.
func (a *App) setCommandOptionValue(fieldType *reflect.StructField, fieldValue *reflect.Value) {
	optName := cmdOptionNamePrefix + fieldType.Tag.Get("optName")
	optAlias := cmdOptionAliasPrefix + fieldType.Tag.Get("optAlias")
	typeName := fieldType.Type.Name()
	switch typeName {
	case "string":
		fieldValue.SetString(a.Args().GetString(optName, optAlias))
	case "bool":
		fieldValue.SetBool(a.Args().HasOption(optName, optAlias))
	case "uint":
		fieldValue.SetInt(a.Args().GetInt(optName, optAlias))
	case "int":
		fieldValue.SetInt(a.Args().GetInt(optName, optAlias))
	default:
		a.exitWithError("unsupported option type \"%s\"\n", typeName)
	}
}

// exitWithError prints out the formatted message and calls os.Exit(1)
func (a *App) exitWithError(format string, params ...interface{}) {
	fmt.Printf(format, params...)
	os.Exit(1)
}
