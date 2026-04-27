package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/givensuman/gh-tag/api"
)

func newCreateCmd() *cobra.Command {
	var commitSHA string
	var message string

	cmd := &cobra.Command{
		Use:   "create <name>",
		Short: "Create a new tag on the remote repository",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			repo, err := CurrentRepo()
			if err != nil {
				return err
			}

			sha := commitSHA
			if sha == "" {
				sha, err = api.GetDefaultBranchSHA(repo.Owner, repo.Repo)
				if err != nil {
					return fmt.Errorf("could not determine default branch HEAD: %w", err)
				}
			}

			if message != "" {
				if err := api.CreateAnnotatedTag(repo.Owner, repo.Repo, name, sha, message); err != nil {
					return fmt.Errorf("failed to create annotated tag: %w", err)
				}

				fmt.Printf("Created annotated tag %q at %s\n", name, sha)
			} else {
				if err := api.CreateLightweightTag(repo.Owner, repo.Repo, name, sha); err != nil {
					return fmt.Errorf("failed to create tag: %w", err)
				}

				fmt.Printf("Created tag %q at %s\n", name, sha)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&commitSHA, "commit", "c", "", "Commit SHA to tag")
	cmd.Flags().StringVarP(&message, "message", "m", "", "Create an annotated tag with this message")

	return cmd
}
