{
  description = "Markdown transcripts of Ray Peat interviews. Built to html with Golang.";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }: flake-utils.lib.eachDefaultSystem (system: let
    pkgs = import nixpkgs { inherit system; };
    inherit (builtins) fetchTarball;

    pagefind = pkgs.callPackage ({ stdenv, lib, extended ? false }: stdenv.mkDerivation rec {
      pname = "pagefind";
      version = "0.12.0";
      src = ./.;

      meta = with lib; {
        description = "Pagefind is a fully static search library that aims to perform well on large sites, while using as little of your users' bandwidth as possible.";
        homepage = "https://pagefind.app";
      };

      installPhase = let
        binary = let
          type = if extended then "pagefind_extended" else "pagefind";

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
    }) {};

    # https://nixos.wiki/wiki/Development_environment_with_nix-shell
    rayPeatRodeoBuildTool = pkgs.callPackage ({ buildGoModule, pkgs }: buildGoModule {
      pname = "ray-peat-rodeo-build-tool";
      version = "unstable";
      src = ./.;
      vendorHash = "sha256-kaPEqLnKSwNfyQ4quDbqoccaq3y2INezxAZOY/BFj30=";
    }) {};

    rayPeatRodeo = pkgs.callPackage ({ stdenv, pkgs }: stdenv.mkDerivation {
      pname = "ray-peat-rodeo";
      version = "unstable";
      src = ./.;

      nativeBuildInputs = with pkgs; [
        go
        modd
        devd
        rayPeatRodeoBuildTool
        pagefind
      ];

      buildPhase = ''
        ${rayPeatRodeoBuildTool}/bin/ray-peat-rodeo build
        ${pagefind}/bin/pagefind --source ./build
      '';

      installPhase = ''
        mkdir $out;
        cp -r ./build/* $out
      '';
    }) {};


  in {
    packages = rec {
      ray-peat-rodeo = rayPeatRodeo;
      default = ray-peat-rodeo;
    };
  });
}
