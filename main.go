package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"os/exec"
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
	CWD        string `json:"cwd"`
	RateLimits struct {
		FiveHour struct {
			UsedPercentage float64 `json:"used_percentage"`
			ResetsAt       float64 `json:"resets_at"`
		} `json:"five_hour"`
	} `json:"rate_limits"`
}

func main() {
	var input Input
	if err := json.NewDecoder(os.Stdin).Decode(&input); err != nil {
		fmt.Print("statusline: parse error")
		return
	}

	var parts []string

	// Model
	if input.Model.DisplayName != "" {
		parts = append(parts, input.Model.DisplayName)
	}

	// Git branch
	cwd := input.Workspace.CurrentDir
	if cwd == "" {
		cwd = input.CWD
	}
	if cwd != "" {
		cmd := exec.Command("git", "-C", cwd, "symbolic-ref", "--short", "HEAD")
		cmd.Env = append(os.Environ(), "GIT_OPTIONAL_LOCKS=0")
		if out, err := cmd.Output(); err == nil {
			branch := strings.TrimSpace(string(out))
			if branch != "" {
				parts = append(parts, branch)
			}
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

	fmt.Print(strings.Join(parts, " | "))
}
