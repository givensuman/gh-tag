package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/cli/go-gh/v2/pkg/tableprinter"
	"github.com/cli/go-gh/v2/pkg/term"
	"github.com/spf13/cobra"
	"github.com/givensuman/gh-tag/api"
)

func newListCmd() *cobra.Command {
	var limit int
	var search string
	var jsonFields string

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

			if jsonFields != "" {
				return printJSON(tags, repo, jsonFields)
			}

			terminal := term.FromEnv()
			printer := tableprinter.New(os.Stdout, terminal.IsTerminalOutput(), 120)
			printer.AddHeader([]string{"TAG NAME", "COMMIT SHA", "URL"})

			for _, t := range tags {
				printer.AddField(t.Name)
				printer.AddField(t.Commit.SHA)
				printer.AddField(fmt.Sprintf("https://github.com/%s/%s/releases/tag/%s", repo.Owner, repo.Repo, t.Name))
				printer.EndRow()
			}

			return printer.Render()
		},
	}

	cmd.Flags().IntVarP(&limit, "limit", "L", 30, "Maximum number of tags to fetch")
	cmd.Flags().StringVarP(&search, "search", "S", "", "Filter tags by name pattern")
	cmd.Flags().StringVar(&jsonFields, "json", "", "Output in JSON format (fields: name,sha,url)")

	return cmd
}

var validJSONFields = map[string]bool{"name": true, "sha": true, "url": true}

func printJSON(tags []api.Tag, repo RepoContext, fields string) error {
	requestedFields := strings.Split(fields, ",")

	for _, f := range requestedFields {
		f = strings.TrimSpace(f)
		if f != "" && !validJSONFields[f] {
			return fmt.Errorf("unknown field %q; valid fields: name, sha, url", f)
		}
	}

	fieldSet := make(map[string]bool)
	for _, f := range requestedFields {
		fieldSet[strings.TrimSpace(f)] = true
	}

	type outTag map[string]string
	out := make([]outTag, 0, len(tags))

	for _, t := range tags {
		entry := make(outTag)
		if fieldSet["name"] {
			entry["name"] = t.Name
		}
		if fieldSet["sha"] {
			entry["sha"] = t.Commit.SHA
		}
		if fieldSet["url"] {
			entry["url"] = fmt.Sprintf("https://github.com/%s/%s/releases/tag/%s", repo.Owner, repo.Repo, t.Name)
		}

		out = append(out, entry)
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")

	return enc.Encode(out)
}
