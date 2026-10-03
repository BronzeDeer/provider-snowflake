{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
    nixpkgs-terraform.url = "github:stackbuilders/nixpkgs-terraform"; # Allows us to fetch old pre-BSL terraform binaries
    nixpkgs-terraform.inputs.nixpkgs.follows = "nixpkgs";

    # Using git+https over https url rather than gitlab api, since the gitlab: scheme chokes on correctly parsing sub-groups in owners
    nono-nix.url = "git+https://gitlab.com/gitlab-org/nix/nono.nix.git";
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
      nixpkgs-terraform,
      nono-nix,
    }:
    # Nix supports many different os/arch combinations, but in practice we mostly care about support x86_64 and arm for linux and darwin (macOS)
    # The list maintained by as "defaultSystems" is exactly this, i.e. ["x86_64-linux" "aarch64-linux" "x86_64-darwin" "aarch64-darwin"]
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        terraformVersion = "1.5.6";
        pkgs = import nixpkgs {
          inherit system;
        };
        nono = nono-nix.lib.nono pkgs;
        lib = pkgs.lib;
      in
      rec {
        formatter = pkgs.nixfmt-tree;
        packages.sandboxed-llm = nono "opencode" pkgs.opencode (
          c:
          with c;
          [
            #(extends-profile "opencode")
            (add-runtime ''
              mkdir -p \
              "$HOME/.opencode" \
              "$HOME/.config/opencode" \
              "$HOME/.cache/opencode" \
              "$HOME/.local/share/opencode" \
              "$HOME/.local/share/opentui" \
              "$HOME/.local/state/opencode"

            '')
            (try-readwrite (noescape "~/.opencode"))
            (try-readwrite (noescape "~/.npm"))
            (try-readwrite (noescape "~/.config/opencode"))
            (try-readwrite (noescape "~/.cache/opencode"))
            (try-readwrite (noescape "~/.local/share/opencode"))
            (try-readwrite (noescape "~/.local/share/opentui"))
            (try-readwrite (noescape "~/.local/state/opencode"))
            (readonly "/tmp") # Necessary to compensate for buns adhoc so loading
            (groups-include "user_caches_linux")
            (groups-include "node_runtime")
            (groups-include "rust_runtime")
            (groups-include "nix_runtime")
            (groups-include "git_config")
            (readwrite ".") # allow-CWD doesn't work because nono.nix also passes a profile file. using (noescape "\"${DIRENV_DIR:-.}\"") would be preferrable, but shell check rejects this with and without the escaped quotes since nono.nix forgot to quote the path correctly
          ]
          ++ (lib.optionals (lib.strings.hasSuffix "-linux" system) [ (groups-include "unlink_protection") ])
        );

        devShells.default = pkgs.mkShell {
          packages =
            with pkgs;
            [

              go
              gopls
              govulncheck
              gotools
              golangci-lint

              # importing the package directly rather than through the overlay, due to a bug in overlay hostPlatform detection on NixOs with stud-ld:
              # https://github.com/NixOS/nixpkgs/issues/325318
              nixpkgs-terraform.packages.${system}."terraform-${terraformVersion}"
              delve

              packages.sandboxed-llm
            ]
            # Allow timetravel debugging on linux on all modern intel and many amd ryzen cpus
            # Note that `sysctl kernel.perf_event_paranoid` needs to be <= 1 in order to record a trace
            ++ (lib.optionals (lib.strings.hasSuffix "-linux" system) [ rr ]);

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
