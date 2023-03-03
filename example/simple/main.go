package main

import "github.com/amedmoore/go-cli"

func main() {
	cli.NewApp().
		Register(NewGreetingCmd()).
		Run()
}
