package main

import (
	"context"
	"flag"

	"github.com/google/go-github/v56/github"
	"github.com/pix-xip/go-command"
	"github.com/pix-xip/go-command/examples/github/handlers"
)

func main() {
	client := github.NewClient(nil)

	root := command.Root().Flags(func(flagSet *flag.FlagSet) {
		flagSet.Bool("verbose", false, "Enable verbose output")
	}).Help("Example command")

	reposCommand := root.SubCommand("repos").Help("Manage GitHub repositories")
	{
		reposCommand.SubCommand("list").Action(handlers.ReposListHandler(client)).Flags(func(flagSet *flag.FlagSet) {
			flagSet.String("user", "", "GitHub user")
		}).Help("List repositories of a GitHub user")
	}

	if err := root.Execute(context.Background()); err != nil {
		panic(err)
	}
}
