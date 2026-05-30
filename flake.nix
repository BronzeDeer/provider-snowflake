{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs";
    flake-utils.url = "github:numtide/flake-utils";
    nixpkgs-terraform.url = "github:stackbuilders/nixpkgs-terraform"; # Allows us to fetch old pre-BSL terraform binaries
  };

  outputs = {
    self,
    nixpkgs,
    flake-utils,
    nixpkgs-terraform,
  }:
  # Nix supports many different os/arch combinations, but in practice we mostly care about support x86_64 and arm for linux and darwin (macOS)
  # The list maintained by as "defaultSystems" is exactly this, i.e. ["x86_64-linux" "aarch64-linux" "x86_64-darwin" "aarch64-darwin"]
  flake-utils.lib.eachDefaultSystem (
    system: let
      terraformVersion = "1.5.6";
      pkgs = import nixpkgs{
        inherit system;
        overlays = [nixpkgs-terraform.overlays.default];
      };
    in rec {
      devShells.default = pkgs.mkShell {
        packages = with pkgs; [

          go_1_24
          gopls
          govulncheck
          gotools
          golangci-lint

          # For testing
          vcluster
          snowflake-cli
          yq-go
          # Needed to generate temporary rsa key pairs
          openssl

          pkgs."terraform-${terraformVersion}"
          delve

        ]
        # Allow timetravel debugging on linux on all modern intel and many amd ryzen cpus
        # Note that `sysctl kernel.perf_event_paranoid` needs to be <= 1 in order to record a trace
        ++ (lib.optionals ( lib.strings.hasSuffix "-linux" system) [rr]);

      hardeningDisable = [ "fortify" ]; # Necessary to allow compiling go builds with debug information

      shellHook = ''
        # Set version of our installed terraform to ensure it gets taken up by the Makefile if it differs from the default
        export TERRAFORM_VERSION=${terraformVersion}
        # Read relevant other terraform provider flags from Makefile to allow us to easily execute the provider directly through delve/IDE
        source <(make print-env)
      '';

      };
      defaultPackage = devShells.default; # Allow nix build to also pick up the shell by default
    }
  );
}
