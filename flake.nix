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

      customPnpm = pkgs.pnpm.override { nodejs = pkgs.nodejs_24; };

      p-shortcut = pkgs.writeShellScriptBin "p" ''
        exec ${customPnpm}/bin/pnpm "$@"
      '';

      no-npm = pkgs.writeShellScriptBin "npm" ''
        echo -e "Por favor, usa \033[0;32mpnpm\033[0m para mantener la consistencia del monorepo."
        exit 1
      '';

      # UFO RPC
      urpc = pkgs.stdenv.mkDerivation {
        pname = "urpc";
        version = "0.3.9";

        src = pkgs.fetchurl {
          url = "https://github.com/uforg/uforpc/releases/download/v0.3.9/urpc-linux-amd64";
          sha256 = "ee7364dc17bcafa4fc983649a0c8200cce0b0b96ae7b78f8a10beb0df6403193";
        };

        # Since it's a single binary, we skip the unpacking and building phases
        phases = [ "installPhase" ];

        installPhase = ''
          mkdir -p $out/bin
          cp $src $out/bin/urpc
          chmod +x $out/bin/urpc
        '';
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
          p-shortcut
          no-npm
          customPnpm
          nodejs_24
          urpc # RPC Framework
          librsvg # To convert SVG Logos to PNG
          imagemagick # To convert PNG Logos to ICO
          uutils-coreutils-noprefix # For sha256sum
          xdotool # For simulating keyboard input (paste from clipboard)
          pkg-config # Required for oto audio library build
          alsa-lib.dev # ALSA headers for audio playback
          libnotify # For desktop notifications
          stdenv.cc.cc.lib # libstdc++ for ONNX Runtime
        ];

        shellHook = ''
          export LD_LIBRARY_PATH="${pkgs.stdenv.cc.cc.lib}/lib:$LD_LIBRARY_PATH"
        '';
      };
    };
}
