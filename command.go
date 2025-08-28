package command

import (
	"context"
	"flag"
	"fmt"
	"os"
	"slices"
	"strings"
)

// Handler represents a command function called by [Command.Execute].
// The command flags can be accessed from the FlagSet parameter using [Lookup] or [flag.Lookup].
type Handler func(context.Context, *flag.FlagSet, []string) error

// Middleware represents a function used to wrap a [Handler]. They can be used to make actions that will execute before or after the command.
// They are also inherited by subcommands, unlike command actions.
type Middleware func(Handler) Handler

// Command represents any command or subcommand of the application.
type Command interface {
	// SubCommand adds a new subcommand to an existing command.
	SubCommand(string) Command

	// Middlewares adds a list of middlewares to the command and its subcommands.
	Middlewares(...Middleware) Command

	// Action sets the action to execute when calling the command.
	Action(Handler) Command

	// Execute runs the command using [os.Args]. It should normally be called on the root command.
	Execute(context.Context) error

	// Help sets the help message of a command.
	Help(string) Command

	// Flags is used to declare the flags of a command.
	Flags(func(*flag.FlagSet)) Command
}

type command struct {
	name        string
	help        string
	middlewares []Middleware
	handler     Handler
	subCommands map[string]*command
	flagSet     *flag.FlagSet
	parent      *command
}

// Root creates a new root command.
func Root() Command {
	command := command{
		name:        os.Args[0],
		subCommands: map[string]*command{},
		flagSet:     flag.CommandLine,
	}

	flag.CommandLine.Usage = command.usage

	return &command
}

func (c *command) SubCommand(name string) Command {
	c.subCommands[name] = &command{
		name:        name,
		subCommands: map[string]*command{},
		flagSet:     flag.NewFlagSet(name, flag.ExitOnError),
		parent:      c,
	}

	c.subCommands[name].flagSet.Usage = c.subCommands[name].usage

	return c.subCommands[name]
}

func (c *command) Middlewares(middlewares ...Middleware) Command {
	c.middlewares = append(c.middlewares, middlewares...)
	return c
}

func (c *command) Action(handler Handler) Command {
	c.handler = handler
	return c
}

func (c *command) Execute(ctx context.Context) error {
	command, args := c, os.Args[1:]
	middlewares := slices.Clone(c.middlewares)
	for {
		if err := command.flagSet.Parse(args); err != nil {
			// This should never occur because the flag sets use flag.ExitOnError
			os.Exit(2) // Use 2 to mimick the behavior of flag.ExitOnError
		}

		args = command.flagSet.Args()
		if len(args) == 0 {
			break
		}

		subCommand, ok := command.subCommands[args[0]]
		if !ok {
			break
		}

		command.flagSet.VisitAll(func(f *flag.Flag) {
			subCommand.flagSet.Var(f.Value, f.Name, f.Usage)
		})

		command = subCommand
		args = args[1:]
		middlewares = append(middlewares, subCommand.middlewares...)
	}

	if command.handler == nil {
		if len(args) > 0 {
			command.flagSet.SetOutput(os.Stderr)
			fmt.Fprintf(command.flagSet.Output(), "command provided but not defined: %s\n", args[0])
			command.usage()
			os.Exit(2) // Use 2 to mimick the behavior of flag.ExitOnError
		}

		command.usage()
		os.Exit(0)
	}

	handler := command.handler
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}

	return handler(ctx, command.flagSet, args)
}

func (c *command) Help(help string) Command {
	c.help = help
	return c
}

func (c *command) Flags(flags func(*flag.FlagSet)) Command {
	flags(c.flagSet)
	return c
}

func (c *command) usage() {
	var builder strings.Builder
	output := c.flagSet.Output()
	c.flagSet.SetOutput(&builder)

	fullCommand := []string{c.name}
	for command := c.parent; command != nil; command = command.parent {
		fullCommand = append([]string{command.name}, fullCommand...)
	}

	var nbFlags int
	c.flagSet.VisitAll(func(*flag.Flag) {
		nbFlags++
	})

	optionsHint := ""
	if nbFlags > 0 {
		optionsHint = " [OPTIONS]"
	}

	subCommandHint := ""
	if len(c.subCommands) > 0 {
		subCommandHint = " [COMMAND]"
		if c.handler == nil {
			subCommandHint = " COMMAND"
		}
	}

	builder.WriteString("Usage: ")
	builder.WriteString(strings.Join(fullCommand, " "))
	builder.WriteString(optionsHint)
	builder.WriteString(subCommandHint)
	builder.WriteString("\n")

	if c.help != "" {
		builder.WriteString("\n")
		builder.WriteString(c.help)
		builder.WriteString("\n")
	}

	if nbFlags > 0 {
		builder.WriteString("\n")
		builder.WriteString("Options:\n")
		c.flagSet.PrintDefaults()
	}

	if len(c.subCommands) > 0 {
		builder.WriteString("\n")
		builder.WriteString("Subcommands:")

		for name, subCommand := range c.subCommands {
			builder.WriteString("\n  ")
			builder.WriteString(name)
			if subCommand.help != "" {
				builder.WriteString("\n\t")
				builder.WriteString(subCommand.help)
			}
		}
	}

	fmt.Fprintln(output, builder.String())
}
