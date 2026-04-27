package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/givensuman/gh-tag/api"
)

func newDeleteCmd() *cobra.Command {
	var confirm bool
	var both bool

	cmd := &cobra.Command{
		Use:   "delete <name>",
		Short: "Delete a tag from the remote repository",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			repo, err := CurrentRepo()
			if err != nil {
				return err
			}

			if !confirm {
				fmt.Fprintf(os.Stderr, "Delete tag %q from %s/%s? [y/N] ", name, repo.Owner, repo.Repo)
				var response string
				_, err := fmt.Fscan(os.Stdin, &response)
				if err != nil {
					return fmt.Errorf("failed to read confirmation response: %w", err)
				}

				if response != "y" && response != "Y" {
					fmt.Println("Aborted.")
					return nil
				}
			}

			if err := api.DeleteTag(repo.Owner, repo.Repo, name); err != nil {
				return fmt.Errorf("failed to delete remote tag: %w", err)
			}

			fmt.Printf("Deleted remote tag %s\n", name)

			if both {
				gitCmd := exec.Command("git", "tag", "-d", name)
				gitCmd.Stdout = os.Stdout
				gitCmd.Stderr = os.Stderr
				if err := gitCmd.Run(); err != nil {
					fmt.Fprintf(os.Stderr, "Warning: could not delete local tag %q: %v\n", name, err)
				}
			}

			return nil
		},
	}

	cmd.Flags().BoolVarP(&confirm, "confirm", "y", false, "Skip confirmation prompt")
	cmd.Flags().BoolVar(&both, "both", false, "Also delete the local tag via git")

	return cmd
}
