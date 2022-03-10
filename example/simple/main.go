package main

import "github.com/skyareas/go-cli"

func main() {
	cli.NewApp().
		Register(NewGreetingCmd()).
		Run()
}
