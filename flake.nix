{
  description = "Markdown transcripts of Ray Peat interviews. Built to html with Golang.";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
    gomod2nix.url = "github:nix-community/gomod2nix";
  };

  outputs = { self, nixpkgs, gomod2nix, flake-utils }: flake-utils.lib.eachDefaultSystem (system: let
    inherit (builtins) fetchTarball;

    pkgs = import nixpkgs {
      inherit system;
      overlays = [
        gomod2nix.overlays.default
      ];
    };
  in {
    packages = rec {
      default = ray-peat-rodeo;

      ray-peat-rodeo = pkgs.runCommand "ray-peat-rodeo" {} ''
        mkdir build
        ${self.packages.${system}.builder}/bin/builder build
        ${self.packages.${system}.pagefind}/bin/pagefind --source ./build
        mkdir $out
        cp --recursive build $out
      '';

      builder = pkgs.buildGoApplication {
        name = "builder";
        pwd = ./.;
        src = ./.;
        modules = ./gomod2nix.toml;
      };

      pagefind = pkgs.stdenv.mkDerivation rec {
        pname = "pagefind";
        version = "0.12.0";
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
                sha256 = "sha256:0ch9vasiqassgm77v0g2fcgz87ax606zq3ixj9k5lsngi06m27ps";
              };
              "aarch64-linux" = {
                system = "aarch64-unknown-linux-musl";
                sha256 = "sha256:0ch9vasiqassgm77v0g2fcgz87ax606zq3ixj9k5lsngi06m27ps";
              };
              "x86_64-darwin" = {
                system = "x86_64-apple-darwin";
                sha256 = "sha256:0ch9vasiqassgm77v0g2fcgz87ax606zq3ixj9k5lsngi06m27ps";
              };
              "x86_64-linux" = {
                system = "x86_64-unknown-linux-musl";
                sha256 = "sha256:1wg8i5vqicz90vqkdi8j5p8wjk6wk2616h5vg9w7afcb67hp712v";
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

      dev = pkgs.runCommand "dev" {
        nativeBuildInputs = [ pkgs.makeWrapper ];
      } ''
        mkdir --parents $out/bin
        echo "${pkgs.devd}/bin/devd ./build" > $out/bin/dev
        chmod +x $out/bin/dev
      '';
    };

    devShells = {
      default = pkgs.mkShell {
        name = "ray-peat-rodeo";
        packages = [
          (pkgs.mkGoEnv { pwd = ./.; })
          pkgs.gomod2nix
          pkgs.modd
          pkgs.devd
          self.packages.${system}.pagefind
        ];
      };
    };
  });
}
