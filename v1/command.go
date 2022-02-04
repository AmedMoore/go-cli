package cli

type Command interface {
	Run(app *App)
}

type CommandEntry struct {
	Name  string
	Alias string
	Help  string
	cmd   Command
}
