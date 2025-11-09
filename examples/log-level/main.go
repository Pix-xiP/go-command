package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"strings"

	"github.com/pix-xip/go-command"
)

func main() {
	root := command.Root().Flags(func(flagSet *flag.FlagSet) {
		flagSet.String("level", "info", "Minimum level of logs to display")
	}).Middlewares(LevelMiddleware)

	root.SubCommand("info").Action(InfoHandler)
	root.SubCommand("error").Action(ErrorHandler)

	if err := root.Execute(context.Background()); err != nil {
		log.Fatal(err)
	}
}

func LevelMiddleware(next command.Handler) command.Handler {
	return func(ctx context.Context, flagSet *flag.FlagSet, args []string) error {
		switch level := command.Lookup[string](flagSet, "level"); level {
		case "debug":
			slog.SetLogLoggerLevel(slog.LevelDebug)

		case "info":
			slog.SetLogLoggerLevel(slog.LevelInfo)

		case "warn":
			slog.SetLogLoggerLevel(slog.LevelWarn)

		case "error":
			slog.SetLogLoggerLevel(slog.LevelError)

		default:
			return fmt.Errorf("unknown level: %s", level)
		}

		return next(ctx, flagSet, args)
	}
}

func InfoHandler(ctx context.Context, _ *flag.FlagSet, args []string) error {
	slog.InfoContext(ctx, strings.Join(args, " "))
	return nil
}

func ErrorHandler(ctx context.Context, _ *flag.FlagSet, args []string) error {
	slog.ErrorContext(ctx, strings.Join(args, " "))
	return nil
}
