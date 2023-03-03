Package `cli` implements command line apps library.

## Usage

A basic example of `go-cli` usage is as follows

```go
package main

import (
	"fmt"

	"github.com/amedmoore/go-cli"
)

// Greeting command to print greeting message.
type Greeting struct {
	Name     string `cli:"name"`
	Help     string `cli:"help"`
	Username string `cli:"option" optName:"name" optAlias:"n" optHelp:"Username to say hi to"`
}

// Run executes the command's logic.
func (g *Greeting) Run(_ *cli.App) {
	fmt.Printf("Hello, %s!\n", g.Username)
}

func main() {
	// create a new app instance
	app := cli.NewApp()
	// register greeting command
	app.Register(&Greeting{Name:  "greeting", Help:  "Say hi!"})
	// run the app
	app.Run()
}
```

### App

App is the main structure of your cli application.

An app instance will always have two commands registered by default `help` and `version` commands.

`version` command will print the name, version, and build time of your application.

`help` command will print help message for the app or a specific command, given the example above, help command can be used in two forms.

1. `$ myapp help` that will print the following message
    
    ```shell
    $ myapp help
    Usage: myapp [command] [options...]
    
    Commands:
      help        print app or command help message
      version     print version and build time
      greeting    Say hi!
    ```

2. `$ myapp greeting help` that will print the following message 
    
    ```shell
    $ myapp greeting help
    Usage: myapp greeting [options...]
    
    Say hi!
    
    Options:
      --name, --n    Username to say hi to
    ```

### Command

A command struct can be any Go struct that has at least these requirements:

1. A string field that has ``cli:"name"`` tag representing the command name.
2. A string field that has ``cli:"help"`` tag representing the command description.
3. A function `Run(app *cli.App)` that will be executed at runtime.

A command can have an optional alias tag, i.e.

```go
type Command struct {
    Alias string `cli:"alias"` // field name will not make a difference
}
```

A command can have any number of option fields identified by the ``cli:"option"`` tag, see example above.

An option (command option field) must have at least ``optName:"[OPTION_NAME]"``, and ``optHelp:"[OPTION_HELP]"`` tags.

An option can have an optional alias tag, i.e.

```go
type Command struct {
    MyOption string `cli:"option" optName:"myOption" optAlias:"myOptionAlias" optHelp:"Option with alias name"`
}
```

## License

This package is licensed under the [MIT License][license] feel free to use it as you want!

[license]: https://github.com/amedmoore/go-cli/blob/main/LICENSE
