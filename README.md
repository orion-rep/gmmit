# GMMIT

A CLI tool that uses Google Gemini AI to generate [Conventional Commits](https://www.conventionalcommits.org)-compliant commit messages and pull request descriptions from your staged changes.

![gmmit command example](images/commit.gif)

---

## Installation

### Linux and macOS

```bash
curl -fsSL https://raw.githubusercontent.com/orion-rep/gmmit/main/scripts/install.sh | bash
```

The script detects your OS and architecture, downloads the latest release binary, and installs it to `/usr/local/bin`.

### Manual

Download the binary for your platform from the [releases page](https://github.com/orion-rep/gmmit/releases), extract it, and move it to a directory on your `PATH` (e.g. `/usr/local/bin`).

### Requirements

- A Gemini API key — get one [here](https://geminiforwork.gwaddons.com/setup-api-keys/create-geminiai-api-key) (free, requires a Google account)

On first run, gmmit will prompt you for the API key and save it to `~/.gmenv`.

---

## Usage

### Generate a commit message

Stage your changes and run:

```bash
git add <files>
gmmit
```

gmmit reads your staged diff and branch name, generates a Conventional Commit message, and asks what to do:

- `y` — create the commit
- `r` — regenerate the message
- `N` — cancel

The generated message follows the pattern `<type>[scope]: <description> (#<ticket-id>)`, where the ticket ID is extracted automatically from your branch name (e.g. `feat/123-my-feature` → `#123`).

### Generate a pull request

```bash
gmmit --pr
```

gmmit diffs your branch against the default branch, generates a PR title and description, and — if you're on GitHub or
Bitbucket — offers to create the PR directly via the API and open it in your browser. Otherwise it copies the content
to your clipboard.

To create PRs via the API you'll need an access token. See [docs/git-tokens.md](docs/git-tokens.md) for setup instructions.

![gmmit pull request example](images/pull-request.gif)

### Push after committing

```bash
gmmit --pu
```

Commits and immediately runs `git push` to the remote origin.

### Skip pre-commit hooks

```bash
gmmit --no-verify
```

Passes `--no-verify` to `git commit`, skipping any configured pre-commit hooks.

### Non-interactive mode

```bash
gmmit -y
```

Skips all confirmation prompts and automatically accepts the first generated message. Useful for scripts and CI pipelines.

### Options reference

| Option | Description |
| ------ | ----------- |
| `--pr` | Generate a PR title and description instead of a commit message |
| `--pu` | Automatically push to remote origin after committing |
| `--no-verify` | Skip pre-commit hooks when creating the commit |
| `-y` | Auto-confirm all prompts, run in non-interactive mode |

---

## Configuration

All configuration is stored in `~/.gmenv`. gmmit prompts for required values on first use and saves them automatically.

| Variable | Default | Description |
| -------- | ------- | ----------- |
| `GMMIT_API_KEY` | — | **Required.** Gemini API key |
| `GMMIT_MODEL` | `gemini-2.5-flash-lite` | Gemini model to use |
| `GMMIT_COMMIT_PATTERN` | `<type>[optional scope]: <description> (#<ticket-id>)` | Commit message format |
| `GMMIT_GH_USER` | — | GitHub username (for PR creation) |
| `GMMIT_GH_PASS` | — | GitHub personal access token (for PR creation) |
| `GMMIT_BB_USER` | — | Bitbucket username (for PR creation) |
| `GMMIT_BB_PASS` | — | Bitbucket app password (for PR creation) |
| `GMMIT_MAX_RETRIES` | `5` | Retries on API 500 errors |
| `GMMIT_RETRY_DELAY` | `5` | Seconds between retries |
| `GMMIT_DEBUG` | `false` | Enable debug logging |

See [docs/configuration.md](docs/configuration.md) for full details.

---

## Troubleshooting

Run gmmit with debug logging to get detailed output:

```bash
GMMIT_DEBUG=true gmmit
```

### FinishReasonSafety — message blocked

```text
<date-time> blocked: candidate: FinishReasonSafety
```

The AI model flagged the diff content as potentially sensitive. Review your staged changes and try again.

### Error 429 — quota exceeded

```text
<date-time> googleapi: Error 429:
```

You've exceeded the Gemini API rate limit. Wait a moment and try again.

### Error 500 — unknown error

```text
<date-time> googleapi: Error 500:
```

Intermittent error from the Gemini API. gmmit retries automatically; if it keeps failing, try again in a few seconds.

### PR fails — ambiguous argument

```text
fatal: ambiguous argument 'remotes/origin/HEAD': unknown revision or path not in the working tree.
```

Your local repo is missing the `origin/HEAD` ref. Fix it with:

```bash
git remote set-head origin --auto
```

---

## Contributing

Contributions are welcome! Fork the repo, make your changes on a branch, and open a pull request.

Every commit in this repo should follow the [Conventional Commits](https://www.conventionalcommits.org) standard — use gmmit for that.

### Build from source

You'll need Go 1.22+, golangci-lint, and goimports:

```bash
brew install go golangci-lint
go install golang.org/x/tools/cmd/goimports@latest
```

Run locally:

```bash
go run ./cmd/gmmit/
```

Build binary:

```bash
make build
```

This repository follows the [Golang standard project layout](https://github.com/golang-standards/project-layout?tab=readme-ov-file#go-directories).

---

## License

This project is under [Apache License v2](LICENSE.md).
