class ClaudeStatusline < Formula
  desc "Fast, minimal statusline generator for Claude Code"
  homepage "https://github.com/jpdarago/claude-statusline"
  url "https://github.com/jpdarago/claude-statusline/archive/refs/tags/v0.1.1.tar.gz"
  sha256 "12bd0bc2e2e3dda57d8235273eeb0a1888bed62e065fb7f3ab97003866312c2e"
  license "MIT"
  head "https://github.com/jpdarago/claude-statusline.git", branch: "main"

  depends_on "go" => :build

  def install
    system "go", "build", *std_go_args(ldflags: "-s -w")
  end

  test do
    assert_match "Opus 4.6",
      pipe_output("#{bin}/claude-statusline", '{"model":{"display_name":"Opus 4.6"}}')
  end
end
