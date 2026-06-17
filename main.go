package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type Input struct {
	Model struct {
		DisplayName string `json:"display_name"`
	} `json:"model"`
	ContextWindow struct {
		RemainingPercentage float64 `json:"remaining_percentage"`
		ContextWindowSize   int     `json:"context_window_size"`
		CurrentUsage        struct {
			InputTokens int `json:"input_tokens"`
		} `json:"current_usage"`
	} `json:"context_window"`
	Worktree struct {
		Name string `json:"name"`
	} `json:"worktree"`
	Workspace struct {
		CurrentDir string `json:"current_dir"`
	} `json:"workspace"`
	CWD     string `json:"cwd"`
	Version string `json:"version"`
	Cost    struct {
		TotalCostUSD      float64 `json:"total_cost_usd"`
		TotalLinesAdded   int     `json:"total_lines_added"`
		TotalLinesRemoved int     `json:"total_lines_removed"`
	} `json:"cost"`
	RateLimits struct {
		FiveHour struct {
			UsedPercentage float64 `json:"used_percentage"`
			ResetsAt       float64 `json:"resets_at"`
		} `json:"five_hour"`
	} `json:"rate_limits"`
}

// GitInfo summarizes the repository state shown in the statusline.
type GitInfo struct {
	Branch string
	Dirty  bool
	Ahead  int
	Behind int
}

// parseGitStatus extracts branch, dirty state, and ahead/behind counts from
// the output of `git status --porcelain=v2 --branch`.
func parseGitStatus(out string) GitInfo {
	var info GitInfo
	for line := range strings.SplitSeq(out, "\n") {
		if line == "" {
			continue
		}
		if !strings.HasPrefix(line, "#") {
			// Any non-header line is a changed/untracked entry.
			info.Dirty = true
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}
		switch fields[1] {
		case "branch.head":
			if fields[2] != "(detached)" {
				info.Branch = fields[2]
			}
		case "branch.ab":
			if len(fields) >= 4 {
				info.Ahead, _ = strconv.Atoi(strings.TrimPrefix(fields[2], "+"))
				info.Behind, _ = strconv.Atoi(strings.TrimPrefix(fields[3], "-"))
			}
		}
	}
	return info
}

func gitInfo(cwd string) GitInfo {
	cmd := exec.Command("git", "-C", cwd, "status", "--porcelain=v2", "--branch")
	cmd.Env = append(os.Environ(), "GIT_OPTIONAL_LOCKS=0")
	out, err := cmd.Output()
	if err != nil {
		return GitInfo{}
	}
	return parseGitStatus(string(out))
}

func formatStatusline(input Input, gitFn func(string) GitInfo) string {
	var parts []string

	// Model
	if input.Model.DisplayName != "" {
		parts = append(parts, input.Model.DisplayName)
	}

	// Git branch + dirty/ahead/behind indicators
	cwd := input.Workspace.CurrentDir
	if cwd == "" {
		cwd = input.CWD
	}
	if cwd != "" {
		if git := gitFn(cwd); git.Branch != "" {
			seg := git.Branch
			if git.Dirty {
				seg += "*"
			}
			if git.Ahead > 0 {
				seg += fmt.Sprintf(" ↑%d", git.Ahead)
			}
			if git.Behind > 0 {
				seg += fmt.Sprintf(" ↓%d", git.Behind)
			}
			parts = append(parts, seg)
		}
	}

	// Worktree
	if input.Worktree.Name != "" {
		parts = append(parts, "wt:"+input.Worktree.Name)
	}

	// Context window
	if input.ContextWindow.ContextWindowSize > 0 {
		remaining := input.ContextWindow.ContextWindowSize - input.ContextWindow.CurrentUsage.InputTokens
		parts = append(parts, fmt.Sprintf("ctx:%.0f%% (%d left)", input.ContextWindow.RemainingPercentage, remaining))
	} else if input.ContextWindow.RemainingPercentage > 0 {
		parts = append(parts, fmt.Sprintf("ctx:%.0f%%", input.ContextWindow.RemainingPercentage))
	}

	// Rate limit
	if input.RateLimits.FiveHour.UsedPercentage > 0 || input.RateLimits.FiveHour.ResetsAt > 0 {
		remaining := math.Round(100 - input.RateLimits.FiveHour.UsedPercentage)
		limitStr := fmt.Sprintf("limit:%.0f%% left", remaining)
		if input.RateLimits.FiveHour.ResetsAt > 0 {
			t := time.Unix(int64(input.RateLimits.FiveHour.ResetsAt), 0)
			limitStr += fmt.Sprintf(" (resets %s)", t.Format("15:04"))
		}
		parts = append(parts, limitStr)
	}

	// Session cost
	if input.Cost.TotalCostUSD > 0 {
		parts = append(parts, fmt.Sprintf("$%.2f", input.Cost.TotalCostUSD))
	}

	// Lines changed this session
	if input.Cost.TotalLinesAdded > 0 || input.Cost.TotalLinesRemoved > 0 {
		parts = append(parts, fmt.Sprintf("+%d -%d", input.Cost.TotalLinesAdded, input.Cost.TotalLinesRemoved))
	}

	// Claude Code version
	if input.Version != "" {
		parts = append(parts, "v"+input.Version)
	}

	return strings.Join(parts, " | ")
}

func main() {
	var input Input
	if err := json.NewDecoder(os.Stdin).Decode(&input); err != nil {
		fmt.Print("statusline: parse error")
		return
	}
	fmt.Print(formatStatusline(input, gitInfo))
}
