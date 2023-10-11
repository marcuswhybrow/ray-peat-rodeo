{
  description = "The engine that builds Ray Peat Rodeo from markdown to HTML";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    rust-overlay.url = "github:oxalica/rust-overlay";
    flake-utils.url = "github:numtide/flake-utils";
    cargo2nix = {
      url = "github:marcuswhybrow/cargo2nix";
      inputs.nixpkgs.follows = "nixpkgs";
      inputs.rust-overlay.follows = "rust-overlay";
    };
  };

  outputs = inputs: with inputs; flake-utils.lib.eachDefaultSystem (system: let
    pkgs = import inputs.nixpkgs {
      inherit system;
      overlays = [ cargo2nix.overlays.default ];
    };

    # https://github.com/cargo2nix/cargo2nix#arguments-to-makepackageset
    rustPkgs = pkgs.rustBuilder.makePackageSet {
      rustVersion = "1.72.0";
      rustChannel = "stable";
      packageFun = import ./Cargo.nix;
    };
  in {
    packages = rec {
      ray-peat-rodeo = (rustPkgs.workspace.ray-peat-rodeo {});

      rpr-with-search = pkgs.stdenv.mkDerivation {
        pname = "rpr-with-search";
        version = "unstable";
        src = ./.;

        buildPhase = let 
          buildScript = pkgs.writeScript "build" ''
            ${ray-peat-rodeo}/bin/ray-peat-rodeo
            ${pagefind}/bin/pagefind --site ./build
          '';
        in ''
          mkdir --parents $out/bin
          cp ${buildScript} $out/bin/rpr-with-search
        '';
      };

      default = inputs.self.packages.${system}.rpr-with-search;

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

    devShell = rustPkgs.workspaceShell {
      name = "ray-peat-rodeo";
      packages = with pkgs; [
        inputs.self.packages.${system}.pagefind
        openssl
        pkg-config
        cargo-watch
        devd
        tmux
        (pkgs.writeScriptBin "watch" ''
          RUST_BACKTRACE=full cargo watch \
            --watch src \
            --ignore stash.yml \
            --watch content \
            --exec "run --bin ray-peat-rodeo -- --update-stash" \
            --shell "pagefind --site ./build"
        '')
        (pkgs.writeScriptBin "serve" ''devd --open --livewatch ./build '')
        (pkgs.writeScriptBin "watch-and-serve" ''tmux new-session -d watch \; split-window serve \; attach'')
        (pkgs.writeScriptBin "deps" ''
          cargo generate-lockfile
          nix run github:cargo2nix/cargo2nix
        '')
      ];
    };
  });
}
