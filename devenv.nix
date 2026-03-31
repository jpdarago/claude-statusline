{ pkgs, ... }:

{
  languages.go.enable = true;

  pre-commit.hooks = {
    govet.enable = true;
    gotest.enable = true;
  };
}
