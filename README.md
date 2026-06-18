# claude-statusline

A fast, minimal statusline generator for [Claude Code](https://docs.anthropic.com/en/docs/claude-code). Written in Go, it reads Claude Code's JSON state from stdin and outputs a formatted statusline string.

## Example output

```
Opus 4.6 | main* ↑2 ↓1 | ctx:75% (150000 left) | limit:90% left (resets 14:32) | $0.42 | +120 -30 | v1.0.30
```

The statusline shows (when available):

- **Model name** — the active Claude model
- **Git branch** — current branch in the working directory, with a `*` when the
  tree is dirty and `↑N`/`↓N` for commits ahead of/behind the upstream
- **Worktree** — active git worktree name (`wt:<name>`)
- **Context window** — remaining percentage and token count
- **Rate limit** — remaining 5-hour quota and reset time
- **Session cost** — total USD spent this session (`$0.42`)
- **Lines changed** — lines added/removed this session (`+120 -30`)
- **Version** — the running Claude Code version (`v1.0.30`)

## Installation

```bash
go install github.com/jpdarago/claude-statusline@latest
```

Or build from source:

```bash
go build -o claude-statusline .
```

### Homebrew

```bash
brew tap jpdarago/claude-statusline https://github.com/jpdarago/claude-statusline
brew install claude-statusline
```

The formula lives in this repository's `Formula/` directory, so the tap URL
points straight at it.

### Nix

This repository is a flake. Run it directly:

```bash
nix run github:jpdarago/claude-statusline
```

Install it into your profile:

```bash
nix profile install github:jpdarago/claude-statusline
```

Or reference it from your own flake to add the package to your system or
home-manager configuration:

```nix
{
  inputs.claude-statusline.url = "github:jpdarago/claude-statusline";
  # ... then add inputs.claude-statusline.packages.${system}.default
  # to environment.systemPackages or home.packages.
}
```

## Configuration

Add the statusline to your Claude Code settings (`~/.claude/settings.json`):

```json
{
  "statusline": {
    "command": "claude-statusline"
  }
}
```

## How it works

Claude Code pipes a JSON object to the statusline command's stdin on each render. The JSON contains the current model, context window usage, rate limits, workspace path, and worktree info. `claude-statusline` parses this and outputs a pipe-separated status string.

## Development

This project uses [devenv](https://devenv.sh) and [direnv](https://direnv.net) for development. On NixOS or with Nix installed:

```bash
# Allow direnv (one-time)
direnv allow

# Run tests
go test ./...

# Build
go build -o claude-statusline .
```

Pre-commit hooks for `go vet` and `go test` are configured automatically via devenv.

## License

MIT
