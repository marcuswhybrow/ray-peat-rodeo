{
  description = "The engine that builds Ray Peat Rodeo from markdown to HTML";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
    templ.url = "github:a-h/templ";
    gomod2nix.url = "github:nix-community/gomod2nix";
  };

  outputs = inputs: with inputs; flake-utils.lib.eachDefaultSystem (system: let
    pkgs = import inputs.nixpkgs {
      overlays = [
        inputs.gomod2nix.overlays.default
        inputs.templ.overlays.default
      ];
      inherit system;
    };
  in {
    apps = rec {
      build = inputs.flake-utils.lib.mkApp {
        drv = pkgs.writeScriptBin "build" ''
          # Echo commands to stdout before running
          set -o xtrace

          ${inputs.templ.packages.${system}.templ}/bin/templ generate
          ${self.packages.${system}.ray-peat-rodeo}/bin/ray-peat-rodeo
          ${pkgs.pagefind}/bin/pagefind --site ./build
          ${pkgs.nodePackages.tailwindcss}/bin/tailwindcss \
            --config ./tailwind.config.js \
            --minify \
            --output ./build/assets/tailwind.css \
        '';
      };
      default = build;
    };

    packages = {
      # https://github.com/nix-community/gomod2nix/blob/master/docs/nix-reference.md
      ray-peat-rodeo = pkgs.buildGoApplication {
        name = "ray-peat-rodeo";
        pwd = ./.;
        src = ./.;
        modules = ./gomod2nix.toml;

        buildPhase = ''
          mkdir -p $out/bin
          ${inputs.templ.packages.${system}.templ}/bin/templ generate
          go build ./cmd/ray-peat-rodeo
          mv ray-peat-rodeo $out/bin
        '';
      };

      # Pagefind builds a JS search API by inspect HTML files.
      # It's not in nixpkgs, so this manually packages it.
      # Using precompiled downloads from GitHub, as building from source takes
      # a while and fails for some reason, anyway this works.
      pagefind = pkgs.stdenv.mkDerivation rec {
        pname = "pagefind";
        version = "1.0.3";
        src = ./.;

        meta = {
          description = "Pagefind is a fully static search library that aims to perform well on large sites, while using as little of your users' bandwidth as possible.";
          homepage = "https://pagefind.app";
        };

        installPhase = let
          binary = let
            type = "pagefind"; # "pagefind_extended" includes Chinese language support

            # https://github.com/CloudCannon/pagefind/releases/tag/v0.12.0
            translations = {
              "aarch64-darwin" = {
                system = "aarch64-apple-darwin";
                sha256 = "sha256:0bsc57cbfymfadxa27a64321g4a9zh3mz8yxbm2l7k0f1a62ysv9";
              };
              "aarch64-linux" = {
                system = "aarch64-unknown-linux-musl";
                sha256 = "sha256:0hikvdjafajjcdlix46chi4w7c7j57g579ssgggc0klx4yjvmxg9";
              };
              "x86_64-darwin" = {
                system = "x86_64-apple-darwin";
                sha256 = "sha256:0p84g2h4khnpahq0r7phbdkw9acy6k6gj2kpdxi4vi08wpnkhlil";
              };
              "x86_64-linux" = {
                system = "x86_64-unknown-linux-musl";
                sha256 = "sha256:0l4fnf8ad2cif2lvsxb9nfw7a2mqzi8bdn0i3b8wv33hzh9az2ak";
              };
            };
          in fetchTarball {
            url = "https://github.com/CloudCannon/pagefind/releases/download/v${version}/${type}-v${version}-${translations.${system}.system}.tar.gz";
            sha256 = translations.${system}.sha256;
          };
        in ''
          mkdir -p $out/bin;
          cp ${binary} $out/bin/pagefind
        '';
      };
    };

    devShells.default = pkgs.mkShell {
      name = "ray-peat-rodeo-devshell";
      packages = with pkgs; [

        # Add "go" command with correct modules in environment
        # https://github.com/nix-community/gomod2nix/blob/master/docs/nix-reference.md
        (mkGoEnv { 
          pwd = ./.; # wordking directory
          modules = ./gomod2nix.toml;
        })

        # Translates go.mod packages into a nix expression.
        gomod2nix

        # Compiles .templ files into .go files
        inputs.templ.packages.${system}.templ

        # Builds JS search API by inspecting HTML build by this package
        self.packages.${system}.pagefind

        # Builds CSS utility classes by inspecting template source code
        nodePackages.tailwindcss

        # Modd should be running tailwind for us, but "watching" doesnt work
        (pkgs.writeScriptBin "tailwind-watch" ''
          tailwind --watch --output ./build/assets/tailwind.css
        '')

        # Dev tools to watch the files system and rerun (above) commands
        modd 

        # Dev HTTP server with auto page reload on file changes
        devd 

      ];
    };
  });
}
