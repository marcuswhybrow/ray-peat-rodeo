{
  description = "The engine that builds Ray Peat Rodeo from markdown to HTML";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  };

  outputs = inputs: let
    pkgs = import inputs.nixpkgs {
      system ="x86_64-linux";
    };
  in {
    packages.x86_64-linux.ray-peat-rodeo-node = pkgs.buildNpmPackage {
      pname = "ray-peat-rodeo-node";
      dontNpmBuild = true;
      version = "0.0.0";
      src = ./.;
      npmDepsHash = "sha256-Upe+DtmfEmjVAIKMV0bbjYgLeSsPC8EEOdNHblAKVRI=";
      # npmDepsHash = pkgs.lib.fakeHash;
    };

    packages.x86_64-linux.build = pkgs.stdenv.mkDerivation {
      pname = "build";
      version = "0.0.0";
      src = ./.;
      buildInputs = [ pkgs.nodejs ];
      buildPhase = ''
        mkdir $out
        node ${inputs.self.packages.x86_64-linux.ray-peat-rodeo-node}/lib/node_modules/ray-peat-rodeo/src/app.js
        mv ./build/* $out
      '';
    };

    packages.x86_64-linux.transcribe = pkgs.writeScriptBin "transcribe" /* bash */ ''
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
    packages.x86_64-linux.ufuzzy = pkgs.fetchFromGitHub {
      owner = "leeoniya";
      repo = "uFuzzy";
      rev = "1.0.14";
      hash = "sha256-g70bBIYc2CWMXVGKKXd1EgcomOJ0CnS3wTYAQWQS0fg=";
    };

    # https://github.com/GoogleChromeLabs/text-fragments-polyfill
    packages.x86_64-linux.text-fragments-polyfill = pkgs.fetchFromGitHub {
      owner = "GoogleChromeLabs";
      repo = "text-fragments-polyfill";
      rev = "53375fea08665bac009bb0aa01a030e065c3933d"; # 2024-01-09
      hash = "sha256-iKIuA10f/oDPj0AVUZOSuI7z+YpHsL1SUVal/hdBBOM=";
    };

    packages.x86_64-linux.set-zero-timeout = pkgs.fetchFromGitHub {
      owner = "shahyar";
      repo = "setZeroTimeout-js";
      rev = "5547e33b873d535ebd69f489be7102912e889eaf";
      hash = "sha256-K42Tz3xN6lf2XKeLlNUSVAGt3hcQZRoNItf71i88z3o=";
    };

    packages.x86_64-linux.copy-static = pkgs.writeShellScriptBin "copy-static" ''
      OUT="$1"

      echo "Copying ./src/client"
      cp -rf ./src/client/* "$OUT"

      echo "Copying ./src/public"
      mkdir --parents "$OUT/public"
      cp -rf ./src/public/* "$OUT/public"

      echo "Copying text-fragments-polyfill"
      mkdir --parents "$OUT/public/scripts/text-fragments-polyfill"
      cp -f ${inputs.self.packages.x86_64-linux.text-fragments-polyfill}/src/* \
      "$OUT/public/scripts/text-fragments-polyfill"

      echo "Copying uFuzzy"
      cp -f ${inputs.self.packages.x86_64-linux.ufuzzy}/dist/uFuzzy.iife.min.js \
        "$OUT/public/scripts"

      echo "Copying setZeroTimeout.js"
      cp -f ${inputs.self.packages.x86_64-linux.set-zero-timeout}/setZeroTimeout.min.js \
        "$OUT/public/scripts"
    '';


    packages.x86_64-linux.pull-types = pkgs.writeShellScriptBin "pull-types" ''
      out="./types"
      echo "Copying Pagefind type definitions to $out";
      mkdir --parents $out
      cp -f ${inputs.self.packages.x86_64-linux.pagefind}/pagefind_web_js/types/index.d.ts $out
    '';

    packages.x86_64-linux.default = inputs.self.packages.x86_64-linux.build;

    devShells.x86_64-linux.default = pkgs.mkShell {
      name = "ray-peat-rodeo-node-devshell";
      packages = [
        # NodeJS is needed to for Tailwind plugins to be found
        pkgs.nodejs

        # Dev tools to watch the files system and rerun (above) commands
        pkgs.modd 

        # AI transcription of audio files
        pkgs.openai-whisper

        # For download's audio files from any URL
        pkgs.yt-dlp

        inputs.self.packages.x86_64-linux.copy-static

        # Convenience bash script using yt-dlp, whisper & whisper-json2md to 
        # transcribe and update assets with a `source.url` in the frontmatter.
        # inputs.self.packages.x86_64-linux.transcribe

        # Get text for PDF assets that don't have it
        # pkgs.ocrmypdf
      ];
    };
  };
}
