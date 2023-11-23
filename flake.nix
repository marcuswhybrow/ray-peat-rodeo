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
      overlays = [inputs.gomod2nix.overlays.default];
      inherit system;
    };
  in {
    apps = rec {
      build = inputs.flake-utils.lib.mkApp {
        drv = pkgs.writeScriptBin "build" ''
          # Echo commands to stdout before running
          set -o xtrace

          ${pkgs.nodejs_20}/bin/npm --version
          ${inputs.templ.packages.${system}.templ}/bin/templ generate
          ${self.packages.${system}.ray-peat-rodeo}/bin/ray-peat-rodeo
          ${pkgs.pagefind}/bin/pagefind --site ./build
          ${pkgs.nodePackages.tailwindcss}/bin/tailwindcss \
            --config ./tailwind.config.js \
            --minify \
            --output ./build/assets/tailwind.css
          cp -r ./internal/assets/* ./build/assets
        '';
      };
      default = build;
    };

    # https://github.com/nix-community/gomod2nix/blob/master/docs/nix-reference.md
    packages.ray-peat-rodeo = pkgs.buildGoApplication {
      name = "ray-peat-rodeo";
      pwd = ./.;
      src = ./.;
      modules = ./gomod2nix.toml;

      buildPhase = ''
        mkdir -p $out/bin
        ${inputs.templ.packages.${system}.templ}/bin/templ generate
        go build ./cmd/ray-peat-rodeo
        mv ray-peat-rodeo $out/bin/ray-peat-rodeo
      '';
    };

    packages.tailwind-scrollbar = pkgs.buildNpmPackage rec {
      pname = "tailwind-scrollbar";
      version = "3.0.5";

      src = pkgs.fetchFromGitHub {
        owner = "adoxography";
        repo = pname;
        rev = "v${version}";
        hash = "sha256-i3tWZmchE+jYoPwOkyUR3j1d7imJNdN+fzC3ainJj8A=";
      };

      npmDepsHash = "sha256-iht2umjqANBwkZR57Y8P+KtH/JkvNTOLj8tR9m91eKo=";

      dontNpmBuild = true;

      # The prepack script runs the build script, which we'd rather do in the build phase.
      #npmPackFlags = [ "--ignore-scripts" ];

      NODE_OPTIONS = "";

      meta = with pkgs.lib; {
        description = "Scrollbar plugin for Tailwind CSS";
        homepage = "https://github.com/adoxography/tailwind-scrollbar";
        license = licenses.mit;
        maintainers = [ "Marcus Whybrow <marcus@whybrow.uk>" ];
      };
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
        inputs.templ.packages.${system}.templ

        # Builds JS search API by inspecting HTML build by this package
        pagefind 

        # NodeJS is needed to for Tailwind plugins to be found
        nodejs_20

        # Scrollbar styling plugin for TaildindCSS
        inputs.self.packages.${system}.tailwind-scrollbar

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
