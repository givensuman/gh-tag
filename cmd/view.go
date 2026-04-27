package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/givensuman/gh-tag/api"
)

func newViewCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "view <name>",
		Short: "Display detailed information about a tag",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			repo, err := CurrentRepo()
			if err != nil {
				return err
			}

			tag, tagObj, err := api.GetTag(repo.Owner, repo.Repo, name)
			if err != nil {
				return err
			}

			fmt.Printf("Tag:    %s\n", tag.Name)
			fmt.Printf("Commit: %s\n", tag.Commit.SHA)

			if tagObj != nil {
				fmt.Printf("Author: %s <%s>\n", tagObj.Tagger.Name, tagObj.Tagger.Email)
				fmt.Printf("Date:   %s\n", tagObj.Tagger.Date)

				if tagObj.Message != "" {
					fmt.Printf("\n%s\n", tagObj.Message)
				}
			}

			return nil
		},
	}
}
