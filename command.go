package cli

// Command is the base type to be implemented by cli.App command.
type Command interface {
	Run(app *App)
}
