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

        meta = {
          description = "Custom static-site-generator. Ran from this repo it consumes markdown files in `./assets` and produces HTML files in `./build`.";
          maintainers = [
            "Marcus Whybrow <marcus@whybrow.uk>"
          ];
          homepage = "https://raypeat.rodeo";
        };
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
          cp -r ./web/static/* ./build/assets

          # mkdir --parents ./build/assets/static/scripts
          # cp ${self.packages.${system}.ufuzzy}/src/uFuzzy.js ./build/assets/static/scripts/uFuzzy.js

          mv ./build $out
        '';

        meta = {
          description = "Creates the final website deployment by running ray-peat-rodeo, pagefind static search, tailwind CSS processing, and copying raw assets into place.";
          maintainers = [
            "Marcus Whybrow <marcus@whybrow.uk>"
          ];
          homepage = "https://github.com/marcuswhybrow/ray-peat-rodeo";
        };
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

        meta = {
          description = "Takes a Whisper AI JSON file and your name and outputs markdown to stdout appropriate to append to Ray Peat Rodeo markdown file.";
          homepage = "https://github.com/marcuswhybrow/ray-peat-rodeo";
          maintainers = [
            "Marcus Whybrow <marcus@whybrow.uk>"
          ];
        };
      };

      transcribe = pkgs.writeScriptBin "transcribe" ''
        set -o xtrace

        asset_path="$1"
        author="$2"
        start="''${3:-0}"

        asset_name=$(basename "$asset_path")
        source_url=$(${pkgs.yq-go}/bin/yq ".source.url | select(.)" "$asset_path")

        tmp_dir_audio=$(mktemp --directory)
        audio_path="$tmp_dir_audio/$asset_name"

        ${pkgs.yt-dlp}/bin/yt-dlp -x "$source_url" -o "$audio_path"
        audio_name_actual=$(ls -AU "$tmp_dir_audio" | head -1)
        audio_path_actual="$tmp_dir_audio/$audio_name_actual"

        if [ "$start" != "0" ] then
          ${pkgs.ffmpeg}/bin/ffmpeg -ss "$start" -i "$audio_path_actual" "$audio_path_actual-trimmed"
          audio_path_actual="$audio_path_actual-trimmed"
        fi

        ls "$tmp_dir_audio"

        tmp_dir_json=$(mktemp --directory)
        ${pkgs.openai-whisper}/bin/whisper --language English --fp16 False --output_format json --output_dir "$tmp_dir_json" "$audio_path_actual"
        json_name=$(ls -AU "$tmp_dir_json" | head -1)
        json_path="$tmp_dir_json/$json_name"

        today=$(date +"%Y-%m-%d")
        yq="${pkgs.yq-go}/bin/yq --front-matter process --inplace"
        $yq ".transcription.date = \"$today\"" "$asset_path"
        $yq ".transcription.author = \"Whisper AI\"" "$asset_path"
        $yq ".transcription.kind = \"auto-generated\"" "$asset_path"
        $yq ".added.author = \"$author\"" "$asset_path"
        $yq ".added.date = \"$today\"" "$asset_path"
        $yq ".completion.content = true" "$asset_path"
        ${inputs.self.packages.x86_64-linux.whisper-json2md}/bin/whisper-json2md "$json_path" "$start" >> "$asset_path"

        rm -r "$tmp_dir_audio"
        rm -r "$tmp_dir_json"
      '';

      # https://github.com/leeoniya/uFuzzy
      ufuzzy = pkgs.fetchFromGitHub {
        owner = "leeoniya";
        repo = "uFuzzy";
        rev = "1.0.14";
        hash = "sha256-g70bBIYc2CWMXVGKKXd1EgcomOJ0CnS3wTYAQWQS0fg=";
      };

      # https://github.com/GoogleChromeLabs/text-fragments-polyfill
      text-fragments-polyfill = pkgs.fetchFromGitHub {
        owner = "GoogleChromeLabs";
        repo = "text-fragments-polyfill";
        rev = "53375fea08665bac009bb0aa01a030e065c3933d"; # 2024-01-09
        hash = "sha256-iKIuA10f/oDPj0AVUZOSuI7z+YpHsL1SUVal/hdBBOM=";
      };

      set-zero-timeout = pkgs.fetchFromGitHub {
        owner = "shahyar";
        repo = "setZeroTimeout-js";
        rev = "5547e33b873d535ebd69f489be7102912e889eaf";
        hash = "sha256-K42Tz3xN6lf2XKeLlNUSVAGt3hcQZRoNItf71i88z3o=";
      };

      copy-assets = pkgs.writeShellScriptBin "copy-assets" ''
        echo "Copying ./web/static"
        mkdir --parents ./build/assets
        cp -rf ./web/static/* ./build/assets

        echo "Copying uFuzzy"
        mkdir --parents ./build/assets/static/scripts
        cp -f ${self.packages.${system}.ufuzzy}/dist/uFuzzy.iife.min.js \
          ./build/assets/scripts

        echo "Copying text-fragments-polyfill"
        mkdir --parents ./build/assets/scripts/text-fragments-polyfill
        cp -f ${self.packages.${system}.text-fragments-polyfill}/src/* \
        ./build/assets/scripts/text-fragments-polyfill

        echo "Copying setZeroTimeout.js"
        mkdir --parents ./build/assets/scripts 
        cp -f ${self.packages.${system}.set-zero-timeout}/setZeroTimeout.min.js \
          ./build/assets/scripts
      '';

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
          cp -r ./web/static/* ./build/assets
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

        # Dev tools to watch the files system and rerun (above) commands
        modd 

        # Dev HTTP server with auto page reload on file changes
        devd 

        # AI transcription of audio files
        openai-whisper

        # For download's audio files from any URL
        yt-dlp

        # Custom tool to convert Whisper JSON output to our markdown format
        inputs.self.packages.x86_64-linux.whisper-json2md

        # Convenience bash script using yt-dlp, whisper & whisper-json2md to 
        # transcribe and update assets with a `source.url` in the frontmatter.
        inputs.self.packages.x86_64-linux.transcribe

        # Allows modd to copy assets from nix packages
        inputs.self.packages.x86_64-linux.copy-assets

        # Get text for PDF assets that don't have it
        # ocrmypdf
      ];
    };
  });
}
