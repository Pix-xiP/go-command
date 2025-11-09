package main

import (
	"context"
	"flag"
	"fmt"
	"strings"

	"github.com/pix-xip/go-command"
)

func main() {
	root := command.Root().Flags(func(flagSet *flag.FlagSet) {
		flagSet.Bool("verbose", false, "Enable verbose output")
	}).Help("Example command")

	root.SubCommand("echo").Action(EchoHandler).Flags(func(flagSet *flag.FlagSet) {
		flagSet.String("case", "", "Case to use (upper, lower)")
	})

	if err := root.Execute(context.Background()); err != nil {
		panic(err)
	}
}

func EchoHandler(ctx context.Context, flagSet *flag.FlagSet, args []string) error {
	verbose := command.Lookup[bool](flagSet, "verbose")
	textCase := command.Lookup[string](flagSet, "case")

	if verbose {
		fmt.Println("command echo called with case: " + textCase)
	}

	switch textCase {
	case "upper":
		fmt.Println(strings.ToUpper(strings.Join(args, " ")))

	case "lower":
		fmt.Println(strings.ToLower(strings.Join(args, " ")))

	default:
		fmt.Println(strings.Join(args, " "))
	}

	return nil
}
