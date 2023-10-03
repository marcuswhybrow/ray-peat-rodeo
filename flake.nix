{
  description = "The engine that builds Ray Peat Rodeo from markdown to HTML";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    rust-overlay.url = "github:oxalica/rust-overlay";
    flake-utils.url = "github:numtide/flake-utils";
    cargo2nix = {
      url = "github:cargo2nix/cargo2nix";
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
    packages.ray-peat-rodeo = (rustPkgs.workspace.ray-peat-rodeo {}).bin;
    packages.default = inputs.self.packages.${system}.ray-peat-rodeo;

    devShell = rustPkgs.workspaceShell {
      name = "ray-peat-rodeo";
      packages = with pkgs; [
        openssl
        pkg-config
        cargo-watch
        devd
        tmux
        (pkgs.writeScriptBin "watch" ''
          RUST_BACKTRACE=full cargo watch \
            --watch src \
            --watch content \
            --exec "run --bin ray-peat-rodeo" \
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
