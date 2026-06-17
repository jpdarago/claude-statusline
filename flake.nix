{
  description = "A fast, minimal statusline generator for Claude Code";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
        claude-statusline = pkgs.buildGoModule {
          pname = "claude-statusline";
          version = "0.1.0";
          src = ./.;
          # No external dependencies; only the standard library is used.
          vendorHash = null;
          meta = {
            description = "A fast, minimal statusline generator for Claude Code";
            homepage = "https://github.com/jpdarago/claude-statusline";
            license = pkgs.lib.licenses.mit;
            mainProgram = "claude-statusline";
          };
        };
      in
      {
        packages = {
          inherit claude-statusline;
          default = claude-statusline;
        };

        apps.default = flake-utils.lib.mkApp {
          drv = claude-statusline;
        };
      }
    );
}
