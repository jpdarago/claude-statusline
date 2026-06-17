package main

import (
	"testing"
	"time"
)

func stubBranch(branch string) func(string) GitInfo {
	return func(string) GitInfo { return GitInfo{Branch: branch} }
}

func stubGit(info GitInfo) func(string) GitInfo {
	return func(string) GitInfo { return info }
}

func TestFormatStatusline_Empty(t *testing.T) {
	got := formatStatusline(Input{}, stubBranch(""))
	if got != "" {
		t.Errorf("expected empty string, got %q", got)
	}
}

func TestFormatStatusline_ModelOnly(t *testing.T) {
	var input Input
	input.Model.DisplayName = "Opus 4.6"
	got := formatStatusline(input, stubBranch(""))
	if got != "Opus 4.6" {
		t.Errorf("got %q", got)
	}
}

func TestFormatStatusline_GitBranch(t *testing.T) {
	var input Input
	input.Model.DisplayName = "Sonnet"
	input.CWD = "/tmp"
	got := formatStatusline(input, stubBranch("main"))
	want := "Sonnet | main"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFormatStatusline_WorkspaceCurrentDirPreferred(t *testing.T) {
	var input Input
	input.CWD = "/should-not-use"
	input.Workspace.CurrentDir = "/use-this"

	var calledWith string
	gitFn := func(cwd string) GitInfo {
		calledWith = cwd
		return GitInfo{Branch: "feat"}
	}
	formatStatusline(input, gitFn)
	if calledWith != "/use-this" {
		t.Errorf("expected branchFn called with /use-this, got %q", calledWith)
	}
}

func TestFormatStatusline_Worktree(t *testing.T) {
	var input Input
	input.Worktree.Name = "fix-123"
	got := formatStatusline(input, stubBranch(""))
	if got != "wt:fix-123" {
		t.Errorf("got %q", got)
	}
}

func TestFormatStatusline_ContextWindowFull(t *testing.T) {
	var input Input
	input.ContextWindow.ContextWindowSize = 200000
	input.ContextWindow.RemainingPercentage = 75
	input.ContextWindow.CurrentUsage.InputTokens = 50000
	got := formatStatusline(input, stubBranch(""))
	want := "ctx:75% (150000 left)"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFormatStatusline_ContextWindowPercentageOnly(t *testing.T) {
	var input Input
	input.ContextWindow.RemainingPercentage = 42
	got := formatStatusline(input, stubBranch(""))
	if got != "ctx:42%" {
		t.Errorf("got %q", got)
	}
}

func TestFormatStatusline_RateLimitUsedOnly(t *testing.T) {
	var input Input
	input.RateLimits.FiveHour.UsedPercentage = 30
	got := formatStatusline(input, stubBranch(""))
	if got != "limit:70% left" {
		t.Errorf("got %q", got)
	}
}

func TestFormatStatusline_RateLimitWithReset(t *testing.T) {
	var input Input
	input.RateLimits.FiveHour.UsedPercentage = 55
	resetTime := time.Date(2026, 3, 30, 14, 32, 0, 0, time.Local)
	input.RateLimits.FiveHour.ResetsAt = float64(resetTime.Unix())
	got := formatStatusline(input, stubBranch(""))
	want := "limit:45% left (resets 14:32)"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFormatStatusline_AllParts(t *testing.T) {
	var input Input
	input.Model.DisplayName = "Opus 4.6"
	input.CWD = "/repo"
	input.Worktree.Name = "wt1"
	input.ContextWindow.ContextWindowSize = 100000
	input.ContextWindow.RemainingPercentage = 80
	input.ContextWindow.CurrentUsage.InputTokens = 20000
	input.RateLimits.FiveHour.UsedPercentage = 10

	got := formatStatusline(input, stubBranch("develop"))
	want := "Opus 4.6 | develop | wt:wt1 | ctx:80% (80000 left) | limit:90% left"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFormatStatusline_EmptyBranchNotShown(t *testing.T) {
	var input Input
	input.Model.DisplayName = "Haiku"
	input.CWD = "/repo"
	got := formatStatusline(input, stubBranch(""))
	if got != "Haiku" {
		t.Errorf("got %q", got)
	}
}

func TestFormatStatusline_GitDirtyAndAheadBehind(t *testing.T) {
	var input Input
	input.CWD = "/repo"
	got := formatStatusline(input, stubGit(GitInfo{Branch: "main", Dirty: true, Ahead: 2, Behind: 1}))
	want := "main* ↑2 ↓1"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFormatStatusline_GitAheadOnly(t *testing.T) {
	var input Input
	input.CWD = "/repo"
	got := formatStatusline(input, stubGit(GitInfo{Branch: "main", Ahead: 3}))
	if got != "main ↑3" {
		t.Errorf("got %q", got)
	}
}

func TestFormatStatusline_Cost(t *testing.T) {
	var input Input
	input.Cost.TotalCostUSD = 0.4231
	got := formatStatusline(input, stubBranch(""))
	if got != "$0.42" {
		t.Errorf("got %q", got)
	}
}

func TestFormatStatusline_LinesChanged(t *testing.T) {
	var input Input
	input.Cost.TotalLinesAdded = 120
	input.Cost.TotalLinesRemoved = 30
	got := formatStatusline(input, stubBranch(""))
	if got != "+120 -30" {
		t.Errorf("got %q", got)
	}
}

func TestFormatStatusline_Version(t *testing.T) {
	var input Input
	input.Version = "1.0.30"
	got := formatStatusline(input, stubBranch(""))
	if got != "v1.0.30" {
		t.Errorf("got %q", got)
	}
}

func TestFormatStatusline_AllPartsExtended(t *testing.T) {
	var input Input
	input.Model.DisplayName = "Opus 4.6"
	input.CWD = "/repo"
	input.Worktree.Name = "wt1"
	input.ContextWindow.ContextWindowSize = 100000
	input.ContextWindow.RemainingPercentage = 80
	input.ContextWindow.CurrentUsage.InputTokens = 20000
	input.RateLimits.FiveHour.UsedPercentage = 10
	input.Cost.TotalCostUSD = 1.5
	input.Cost.TotalLinesAdded = 10
	input.Cost.TotalLinesRemoved = 5
	input.Version = "1.0.30"

	got := formatStatusline(input, stubGit(GitInfo{Branch: "develop", Dirty: true, Ahead: 1}))
	want := "Opus 4.6 | develop* ↑1 | wt:wt1 | ctx:80% (80000 left) | limit:90% left | $1.50 | +10 -5 | v1.0.30"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestParseGitStatus(t *testing.T) {
	out := "# branch.oid abc123\n# branch.head main\n# branch.upstream origin/main\n# branch.ab +2 -1\n1 .M N... 100644 100644 100644 abc def file.go\n"
	got := parseGitStatus(out)
	want := GitInfo{Branch: "main", Dirty: true, Ahead: 2, Behind: 1}
	if got != want {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

func TestParseGitStatus_Clean(t *testing.T) {
	out := "# branch.oid abc123\n# branch.head main\n# branch.ab +0 -0\n"
	got := parseGitStatus(out)
	want := GitInfo{Branch: "main"}
	if got != want {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

func TestParseGitStatus_Detached(t *testing.T) {
	out := "# branch.oid abc123\n# branch.head (detached)\n"
	got := parseGitStatus(out)
	if got.Branch != "" {
		t.Errorf("expected empty branch for detached HEAD, got %q", got.Branch)
	}
}
