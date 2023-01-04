package main

import "github.com/ahmedmkamal/go-cli"

func main() {
	cli.NewApp().
		Register(NewGreetingCmd()).
		Run()
}
