{
  description = "Tribar Voice DevShell";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";
  };

  outputs =
    { nixpkgs, ... }:
    let
      system = "x86_64-linux";
      pkgs = import nixpkgs {
        inherit system;
      };
    in
    {
      devShells.${system}.default = pkgs.mkShell {
        packages = with pkgs; [
          go
          gopls
          delve
          golangci-lint
          go-task
          git
          curl
          librsvg # To convert SVG Logos to PNG
          imagemagick # To convert PNG Logos to ICO
          uutils-coreutils-noprefix # For sha256sum
          xdotool # For simulating keyboard input (paste from clipboard)
        ];
      };
    };
}
