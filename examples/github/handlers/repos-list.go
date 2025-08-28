package handlers

import (
	"context"
	"flag"
	"fmt"

	"github.com/google/go-github/v56/github"
	"github.com/pix-xip/go-command"
)

func ReposListHandler(ghClient *github.Client) command.Handler {
	return func(ctx context.Context, flagSet *flag.FlagSet, _ []string) error {
		user := command.Lookup[string](flagSet, "user")
		if user == "" {
			return fmt.Errorf("missing required flag: user")
		}

		repos, _, err := ghClient.Repositories.List(ctx, user, nil)
		if err != nil {
			return fmt.Errorf("failed to list repositories: %v", err)
		}

		fmt.Printf("Repositories of %s:\n", user)
		for _, repo := range repos {
			fmt.Printf("- %s\n", *repo.Name)
		}

		return nil
	}
}
