package command

import (
	"bytes"
	"flag"
	"strings"
	"testing"
)

func TestUsageColorsFlagsAndSubcommands(t *testing.T) {
	t.Parallel()

	var output bytes.Buffer

	root := &command{
		name:        "go-command",
		subCommands: map[string]*command{},
		flagSet:     flag.NewFlagSet("go-command", flag.ContinueOnError),
	}
	root.flagSet.SetOutput(&output)
	root.flagSet.String("level", "info", "Minimum level to display")
	root.flagSet.Bool("verbose", false, "Enable verbose output")
	root.subCommands["serve"] = &command{
		name:        "serve",
		help:        "Run the server",
		subCommands: map[string]*command{},
		flagSet:     flag.NewFlagSet("serve", flag.ContinueOnError),
		parent:      root,
	}

	root.usage()

	usage := output.String()

	if !strings.Contains(usage, bold("Usage:")+" go-command [OPTIONS] COMMAND") {
		t.Fatalf("usage missing bold usage header:\n%s", usage)
	}

	if !strings.Contains(usage, bold("Options:")+"\n  -"+purple("level")+" "+blue("string")) {
		t.Fatalf("usage missing colored string flag:\n%s", usage)
	}

	if !strings.Contains(usage, "Minimum level to display (default "+yellow(`"info"`)+")") {
		t.Fatalf("usage missing string default:\n%s", usage)
	}

	if !strings.Contains(usage, "  -"+purple("verbose")+"\n    Enable verbose output") {
		t.Fatalf("usage missing colored bool flag:\n%s", usage)
	}

	if !strings.Contains(usage, bold("Subcommands:")+"\n  "+purple("serve")+":\n    Run the server") {
		t.Fatalf("usage missing aligned subcommand help:\n%s", usage)
	}
}

func TestFormatDefaultValueHandlesDuration(t *testing.T) {
	t.Parallel()

	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	fs.Duration("timeout", 2, "Timeout")

	formatted := formatFlag(fs.Lookup("timeout"))
	if !strings.Contains(formatted, blue("duration")) {
		t.Fatalf("formatted flag missing colored duration type: %s", formatted)
	}

	if !strings.Contains(formatted, "(default "+yellow("2ns")+")") {
		t.Fatalf("formatted flag missing duration default: %s", formatted)
	}
}

func TestFormatSubcommandAlignsHelp(t *testing.T) {
	t.Parallel()

	formatted := formatSubcommand("version", "Display the version")
	if formatted != "  "+purple("version")+":\n    Display the version" {
		t.Fatalf("unexpected subcommand formatting: %q", formatted)
	}
}
