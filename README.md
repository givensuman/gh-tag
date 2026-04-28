<div align="center">
  <img src="./assets/logo.png" alt="gh-tag" width="200" />
</div>

# gh-tag

A GitHub CLI extension for managing tags.

## Installation

```sh
gh extension install givensuman/gh-tag
```

## Usage

### List tags

```sh
gh tag list
gh tag list --limit 50
gh tag list --search v1.
gh tag list --json
```

### Create a tag

```sh
gh tag create v1.2.3
gh tag create v1.2.3 --commit abc1234
gh tag create v1.2.3 --message "Release v1.2.3"
```

### Delete a tag

```sh
gh tag delete v1.2.3
gh tag delete v1.2.3 --confirm
gh tag delete v1.2.3 --both
```

### View a tag

```sh
gh tag view v1.2.3
```

### Override repository

All commands accept `--repo owner/repo` to target a repository other than the one
detected from the current directory:

```sh
gh tag list --repo cli/cli
```

## License

[MIT](./LICENSE)
