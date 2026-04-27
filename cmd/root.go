// Package cmd defines the root command and shared
// functionality for the gh-tag CLI tool.
package cmd

import (
	"fmt"
	"os"

	"github.com/cli/go-gh/v2/pkg/repository"
	"github.com/spf13/cobra"
)

var rootRepo string

var rootCmd = &cobra.Command{
	Use:   "gh-tag",
	Short: "Manage GitHub tags remotely",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if rootRepo != "" {
			return os.Setenv("GH_REPO", rootRepo)
		}

		return nil
	},
}

// RepoContext holds resolved owner/repo for subcommands.
type RepoContext struct {
	Owner string
	Repo  string
}

// CurrentRepo resolves the owner and repo from the current git remote,
// or from GH_REPO environment variable.
func CurrentRepo() (RepoContext, error) {
	r, err := repository.Current()
	if err != nil {
		return RepoContext{}, fmt.Errorf("could not determine repository: %w", err)
	}

	return RepoContext{Owner: r.Owner, Repo: r.Name}, nil
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&rootRepo, "repo", "", "Repository to use (owner/repo), overrides git remote detection")
	rootCmd.AddCommand(newListCmd())
	rootCmd.AddCommand(newCreateCmd())
	rootCmd.AddCommand(newDeleteCmd())
	rootCmd.AddCommand(newViewCmd())
}
