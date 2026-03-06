package command

import (
	"flag"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	ansiBold   = "\x1b[1m"
	ansiReset  = "\x1b[0m"
	ansiPurple = "\x1b[35m"
	ansiBlue   = "\x1b[34m"
	ansiYellow = "\x1b[33m"
)

func colourize(code string, value string) string {
	return code + value + ansiReset
}

func purple(value string) string {
	return colourize(ansiPurple, value)
}

func blue(value string) string {
	return colourize(ansiBlue, value)
}

func yellow(value string) string {
	return colourize(ansiYellow, value)
}

func bold(value string) string {
	return colourize(ansiBold, value)
}

func formatFlag(f *flag.Flag) string {
	name, usage := flag.UnquoteUsage(f)

	var builder strings.Builder
	builder.WriteString("  -")
	builder.WriteString(purple(f.Name))
	if name != "" {
		builder.WriteString(" ")
		builder.WriteString(blue(name))
	}

	if usage != "" {
		builder.WriteString("\n    ")
		builder.WriteString(usage)
	}

	if defaultValue, ok := formatDefaultValue(f); ok {
		builder.WriteString(" (default ")
		builder.WriteString(yellow(defaultValue))
		builder.WriteString(")")
	}

	return builder.String()
}

func formatSubcommand(name string, help string) string {
	var builder strings.Builder
	builder.WriteString("  ")
	builder.WriteString(purple(name))
	builder.WriteString(":")
	if help != "" {
		builder.WriteString("\n    ")
		builder.WriteString(help)
	}

	return builder.String()
}

func formatDefaultValue(f *flag.Flag) (string, bool) {
	getter, ok := f.Value.(flag.Getter)
	if !ok {
		return "", false
	}

	switch value := getter.Get().(type) {
	case bool:
		if !value {
			return "", false
		}
	case string:
		if value == "" {
			return "", false
		}
		return strconv.Quote(value), true
	case int:
		if value == 0 {
			return "", false
		}
	case int64:
		if value == 0 {
			return "", false
		}
	case uint:
		if value == 0 {
			return "", false
		}
	case uint64:
		if value == 0 {
			return "", false
		}
	case float64:
		if value == 0 {
			return "", false
		}
	case time.Duration:
		if value == 0 {
			return "", false
		}
	default:
		if f.DefValue == "" {
			return "", false
		}
	}

	return fmt.Sprint(f.DefValue), true
}
