# claude-statusline

A fast, minimal statusline generator for [Claude Code](https://docs.anthropic.com/en/docs/claude-code). Written in Go, it reads Claude Code's JSON state from stdin and outputs a formatted statusline string.

## Example output

```
Opus 4.6 | main | ctx:75% (150000 left) | limit:90% left (resets 14:32)
```

The statusline shows (when available):

- **Model name** — the active Claude model
- **Git branch** — current branch in the working directory
- **Worktree** — active git worktree name (`wt:<name>`)
- **Context window** — remaining percentage and token count
- **Rate limit** — remaining 5-hour quota and reset time

## Installation

```bash
go install github.com/jpdarago/claude-statusline@latest
```

Or build from source:

```bash
go build -o claude-statusline .
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
