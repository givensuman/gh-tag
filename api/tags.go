// Package api handles interactions with the GitHub API.
package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/cli/go-gh/v2/pkg/api"
)

// Tag represents a GitHub tag as returned by the list endpoint.
type Tag struct {
	Name   string `json:"name"`
	Commit struct {
		SHA string `json:"sha"`
		URL string `json:"url"`
	} `json:"commit"`
	ZipballURL string `json:"zipball_url"`
	TarballURL string `json:"tarball_url"`
}

// TagObject is the object returned by
// GET /repos/{owner}/{repo}/git/tags/{sha}
type TagObject struct {
	Tag     string `json:"tag"`
	SHA     string `json:"sha"`
	Message string `json:"message"`
	Tagger  struct {
		Name  string `json:"name"`
		Email string `json:"email"`
		Date  string `json:"date"`
	} `json:"tagger"`
	Object struct {
		SHA  string `json:"sha"`
		Type string `json:"type"`
		URL  string `json:"url"`
	} `json:"object"`
}

func jsonBody(v any) (*bytes.Reader, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(b), nil
}

// ListTags returns up to `limit` tags, optionally filtering
// by `search` substring.
func ListTags(owner, repo string, limit int, search string) ([]Tag, error) {
	client, err := api.DefaultRESTClient()
	if err != nil {
		return nil, err
	}

	var all []Tag
	page := 1
	perPage := 100

	for len(all) < limit {
		var pageTags []Tag
		path := fmt.Sprintf("repos/%s/%s/tags?per_page=%d&page=%d", owner, repo, perPage, page)
		if err := client.Get(path, &pageTags); err != nil {
			return nil, err
		}

		if len(pageTags) == 0 {
			break
		}

		for _, t := range pageTags {
			if search == "" || strings.Contains(t.Name, search) {
				all = append(all, t)
				if len(all) >= limit {
					break
				}
			}
		}

		page++
	}

	return all, nil
}

// GetDefaultBranchSHA returns the HEAD commit SHA of the default branch.
func GetDefaultBranchSHA(owner, repo string) (string, error) {
	client, err := api.DefaultRESTClient()
	if err != nil {
		return "", err
	}

	var repoInfo struct {
		DefaultBranch string `json:"default_branch"`
	}

	if err := client.Get(fmt.Sprintf("repos/%s/%s", owner, repo), &repoInfo); err != nil {
		return "", err
	}

	var branch struct {
		Commit struct {
			SHA string `json:"sha"`
		} `json:"commit"`
	}

	if err := client.Get(fmt.Sprintf("repos/%s/%s/branches/%s", owner, repo, repoInfo.DefaultBranch), &branch); err != nil {
		return "", err
	}

	return branch.Commit.SHA, nil
}

// CreateLightweightTag creates a lightweight tag.
// https://git-scm.com/book/en/v2/Git-Basics-Tagging
func CreateLightweightTag(owner, repo, name, sha string) error {
	client, err := api.DefaultRESTClient()
	if err != nil {
		return err
	}

	body, err := jsonBody(map[string]string{
		"ref": "refs/tags/" + name,
		"sha": sha,
	})
	if err != nil {
		return err
	}

	var result map[string]any
	return client.Post(fmt.Sprintf("repos/%s/%s/git/refs", owner, repo), body, &result)
}

// CreateAnnotatedTag creates an annotated tag.
// https://git-scm.com/book/en/v2/Git-Basics-Tagging
func CreateAnnotatedTag(owner, repo, name, sha, message string) error {
	client, err := api.DefaultRESTClient()
	if err != nil {
		return err
	}

	tagBody, err := jsonBody(map[string]string{
		"tag":     name,
		"message": message,
		"object":  sha,
		"type":    "commit",
	})
	if err != nil {
		return err
	}

	var tagObj struct {
		SHA string `json:"sha"`
	}

	if err := client.Post(fmt.Sprintf("repos/%s/%s/git/tags", owner, repo), tagBody, &tagObj); err != nil {
		return err
	}

	refBody, err := jsonBody(map[string]string{
		"ref": "refs/tags/" + name,
		"sha": tagObj.SHA,
	})
	if err != nil {
		return err
	}

	var result map[string]any
	return client.Post(fmt.Sprintf("repos/%s/%s/git/refs", owner, repo), refBody, &result)
}

// DeleteTag removes a tag ref from the remote.
func DeleteTag(owner, repo, name string) error {
	client, err := api.DefaultRESTClient()
	if err != nil {
		return err
	}

	return client.Delete(fmt.Sprintf("repos/%s/%s/git/refs/tags/%s", owner, repo, name), nil)
}

// GetTag fetches a specific tag by name. Returns the Tag
// and, if annotated, its TagObject.
func GetTag(owner, repo, name string) (Tag, *TagObject, error) {
	client, err := api.DefaultRESTClient()
	if err != nil {
		return Tag{}, nil, err
	}

	page := 1
	for {
		var tags []Tag
		path := fmt.Sprintf("repos/%s/%s/tags?per_page=100&page=%d", owner, repo, page)

		if err := client.Get(path, &tags); err != nil {
			return Tag{}, nil, err
		}

		if len(tags) == 0 {
			break
		}

		for _, t := range tags {
			if t.Name == name {
				var tagObj TagObject
				err = client.Get(fmt.Sprintf("repos/%s/%s/git/tags/%s", owner, repo, t.Commit.SHA), &tagObj)
				if err != nil {
					var httpErr *api.HTTPError
					if errors.As(err, &httpErr) && httpErr.StatusCode == 404 {
						return t, nil, nil
					}

					return Tag{}, nil, err
				}

				return t, &tagObj, nil
			}
		}

		page++
	}

	return Tag{}, nil, fmt.Errorf("tag %q not found", name)
}
