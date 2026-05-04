![gh-tag logo](./assets/logo.png)

# gh-tag

A GitHub CLI extension for managing tags.

## Installation

```bash
gh extension install givensuman/gh-tag
```

## Usage

### List tags

```bash
gh tag list
gh tag list --limit 50
gh tag list --search v1.
gh tag list --json
```

### Create a tag

```bash
gh tag create v1.2.3
gh tag create v1.2.3 --commit abc1234
gh tag create v1.2.3 --message "Release v1.2.3"
```

### Delete a tag

```bash
gh tag delete v1.2.3
gh tag delete v1.2.3 --confirm
gh tag delete v1.2.3 --both
```

### View a tag

```bash
gh tag view v1.2.3
```

### Override repository

All commands accept `--repo owner/repo` to target a repository other than the one
detected from the current directory:

```bash
gh tag list --repo cli/cli
```

## License

[MIT](./LICENSE)
