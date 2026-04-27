package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/cli/go-gh/v2/pkg/tableprinter"
	"github.com/cli/go-gh/v2/pkg/term"
	"github.com/spf13/cobra"
	"github.com/givensuman/gh-tag/api"
)

func newListCmd() *cobra.Command {
	var limit int
	var search string
	var json bool

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List tags in the repository",
		RunE: func(cmd *cobra.Command, args []string) error {
			if limit <= 0 {
				return fmt.Errorf("--limit must be greater than 0")
			}

			repo, err := CurrentRepo()
			if err != nil {
				return err
			}

			tags, err := api.ListTags(repo.Owner, repo.Repo, limit, search)
			if err != nil {
				return err
			}

			if json {
				return printJSON(tags, repo)
			}

			terminal := term.FromEnv()
			printer := tableprinter.New(os.Stdout, terminal.IsTerminalOutput(), 120)
			printer.AddHeader([]string{"TAG NAME", "COMMIT SHA", "URL"})

			for _, t := range tags {
				printer.AddField(t.Name)
				printer.AddField(t.Commit.SHA)
				printer.AddField(formatURL(repo, t.Name))
				printer.EndRow()
			}

			return printer.Render()
		},
	}

	cmd.Flags().IntVarP(&limit, "limit", "L", 30, "Maximum number of tags to fetch")
	cmd.Flags().StringVarP(&search, "search", "S", "", "Filter tags by name pattern")
	cmd.Flags().BoolVarP(&json, "json", "J", false, "Output in JSON format")

	return cmd
}

func printJSON(tags []api.Tag, repo RepoContext) error {
	type outTag map[string]string
	out := make([]outTag, 0, len(tags))

	for _, t := range tags {
		entry := make(outTag)
		entry["name"] = t.Name
		entry["sha"] = t.Commit.SHA
		entry["url"] = formatURL(repo, t.Name)

		out = append(out, entry)
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")

	return enc.Encode(out)
}

func formatURL(repo RepoContext, tagName string) string {
	return fmt.Sprintf("https://github.com/%s/%s/releases/tag/%s", repo.Owner, repo.Repo, tagName)
}
