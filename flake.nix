{
  description = "Shared theme-state detector: resolve, parse and watch the theme-state file";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-parts.url = "github:hercules-ci/flake-parts";
    treefmt-nix = {
      url = "github:numtide/treefmt-nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };
    git-hooks-nix = {
      url = "github:cachix/git-hooks.nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs = inputs @ {flake-parts, ...}:
    flake-parts.lib.mkFlake {inherit inputs;} {
      imports = [
        inputs.treefmt-nix.flakeModule
        inputs.git-hooks-nix.flakeModule
      ];

      systems = ["x86_64-linux" "aarch64-linux" "x86_64-darwin" "aarch64-darwin"];

      perSystem = {
        pkgs,
        config,
        ...
      }: let
        nilaway = pkgs.buildGoModule rec {
          name = "nilaway";
          src = pkgs.fetchFromGitHub {
            owner = "uber-go";
            repo = "nilaway";
            rev = "acb8859b9031";
            hash = "sha256-GvDZ5tlvOrTI93tYcIcLd45ZHdqwFopVtoBffD/kbuM=";
          };
          vendorHash = "sha256-qVmvDneq6V/q5UHZ/Cjjqd5/XPPNfvVGoxwg9nz4/Ds=";
          subPackages = ["cmd/nilaway"];
          doCheck = false;
        };
      in {
        treefmt = {
          projectRootFile = "flake.nix";
          programs = {
            alejandra.enable = true;
            gofmt.enable = true;
          };
        };

        pre-commit.settings.hooks = {
          statix.enable = true;
          deadnix.enable = true;
          alejandra.enable = true;
          typos.enable = true;
          check-merge-conflicts.enable = true;
          trim-trailing-whitespace.enable = true;
        };

        # No packages.default: this is a library, so there is no main package to build.
        devShells.default = pkgs.mkShell {
          inherit (config.pre-commit) shellHook;
          packages =
            config.pre-commit.settings.enabledPackages
            ++ [
              pkgs.go
              pkgs.gopls
              pkgs.gotools
              pkgs.golangci-lint
              nilaway
              config.treefmt.build.wrapper
            ];
        };

        # Static-analysis and test gates; run via `nix flake check`.
        checks = {
          golangci-lint-run =
            pkgs.runCommand "golangci-lint-run" {
              nativeBuildInputs = [pkgs.golangci-lint pkgs.go];
              src = ./.;
            } ''
              export HOME="$TMPDIR"
              cd "$src"
              golangci-lint run ./... >&2
              touch "$out"
            '';

          nilaway-check =
            pkgs.runCommand "nilaway-check" {
              nativeBuildInputs = [nilaway pkgs.go];
              src = ./.;
            } ''
              export HOME="$TMPDIR"
              cd "$src"
              nilaway -include-pkgs="github.com/noamsto/themestate" ./... >&2
              touch "$out"
            '';

          go-test-race =
            pkgs.runCommand "go-test-race" {
              nativeBuildInputs = [pkgs.go pkgs.stdenv.cc];
              src = ./.;
            } ''
              export HOME="$TMPDIR"
              cd "$src"
              go test -race ./... >&2
              touch "$out"
            '';
        };
      };
    };
}
