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
      rustVersion = "1.68.2";
      rustChannel = "stable";
      packageFun = import ./engine/Cargo.nix;
    };
  in {
    packages.engine = (rustPkgs.workspace.engine {}).bin;
    packages.default = inputs.self.packages.${system}.engine;

    devShell = rustPkgs.workspaceShell {
      name = "ray-peat-rodeo";
      packages = with pkgs; [
        cargo-watch
        devd
        tmux
        (pkgs.writeScriptBin "serve" ''
          tmux new-session -d \
            cargo watch \
              --workdir engine \
              --ignore build \
              --exec "run -- ../build --clean" \
          \; \
          split-window -h \
          devd \
            --open \
            --livewatch \
            ./build \
          \; attach
        '')
      ];
    };
  });
}
