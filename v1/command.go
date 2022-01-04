package cli

type Command interface {
	Name() string
	Desc() string
	Run(app *App)
}
