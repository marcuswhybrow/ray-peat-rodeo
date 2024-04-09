{
  description = "The engine that builds Ray Peat Rodeo from markdown to HTML";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
    gomod2nix.url = "github:nix-community/gomod2nix";
    tailwind-scrollbar.url = "github:marcuswhybrow/tailwind-scrollbar";
  };

  outputs = inputs: with inputs; flake-utils.lib.eachDefaultSystem (system: let
    pkgs = import inputs.nixpkgs {
      overlays = [inputs.gomod2nix.overlays.default];
      inherit system;
    };
  in {
    # https://github.com/nix-community/gomod2nix/blob/master/docs/nix-reference.md
    packages = rec {
      ray-peat-rodeo = pkgs.buildGoApplication {
        name = "ray-peat-rodeo";
        pwd = ./.;
        src = ./.;
        modules = ./gomod2nix.toml;

        buildPhase = ''
          mkdir -p $out/bin
          ${pkgs.templ}/bin/templ generate
          go build ./cmd/ray-peat-rodeo
          mv ray-peat-rodeo $out/bin/ray-peat-rodeo
        '';
      };

      build = pkgs.stdenv.mkDerivation {
        pname = "build";
        version = "unstable";
        src = ./.;

        buildInputs = [
          inputs.tailwind-scrollbar.packages.x86_64-linux.default
          pkgs.nodejs_20
        ];

        buildPhase = ''
          ${self.packages.${system}.ray-peat-rodeo}/bin/ray-peat-rodeo
          ${pkgs.pagefind}/bin/pagefind --site ./build
          ${pkgs.nodePackages.tailwindcss}/bin/tailwindcss \
            --config ./tailwind.config.js \
            --minify \
            --output ./build/assets/tailwind.css
          cp -r ./internal/assets/* ./build/assets
          mv ./build $out
        '';
      };

      whisper-json2md = pkgs.buildGoApplication {
        name = "whisper-json2md";
        pwd = ./.;
        src = ./.;
        modules = ./gomod2nix.toml;

        buildPhase = ''
          mkdir -p $out/bin
          go build ./cmd/whisper-json2md
          mv whisper-json2md $out/bin/whisper-json2md
        '';
      };

      default = build;
    };

    # Run `nix develop` to enter a shell containing all dependencies.
    # One may use nix-direnv to auto load said shell on cd into project.
    devShells.default = pkgs.mkShell {
      name = "ray-peat-rodeo-devshell";
      packages = with pkgs; [
        (pkgs.writeScriptBin "build" ''
          # Echo commands to stdout before running
          set -o xtrace

          templ generate && \
          go run ./cmd/ray-peat-rodeo && \
          pagefind --site ./build && \
          tailwind \
            --config ./tailwind.config.js \
            --minify \
            --output ./build/assets/tailwind.css && \
          cp -r ./internal/assets/* ./build/assets
        '')

        # Add "go" command with correct modules in environment
        # https://github.com/nix-community/gomod2nix/blob/master/docs/nix-reference.md
        (mkGoEnv { 
          pwd = ./.; # wordking directory
          modules = ./gomod2nix.toml;
        })

        # Translates go.mod packages into a nix expression.
        gomod2nix

        # Compiles .templ files into .go files
        templ

        # Builds JS search API by inspecting HTML build by this package
        pagefind 

        # NodeJS is needed to for Tailwind plugins to be found
        nodejs_20

        # Scrollbar styling plugin for TaildindCSS
        inputs.tailwind-scrollbar.packages.x86_64-linux.default

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

        # AI transcription of audio files
        openai-whisper

        # For download's audio files from any URL
        yt-dlp
      ];
    };
  });
}
