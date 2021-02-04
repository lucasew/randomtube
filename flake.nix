{
  description = "memetube";
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-20.09";
  };
  outputs = { self, nixpkgs, ... }: 
  let
    pkgs = import nixpkgs {
      system = "x86_64-linux";
    };
  in rec {
    devShell.x86_64-linux = import ./shell.nix {inherit pkgs;};
    package = pkgs.callPackage ./package.nix {};
    app = pkgs.writeShellScriptBin "app" ''
      PATH=$PATH:${pkgs.ffmpeg}/bin
      exec ${package}/bin/randomtube "$@"
    '';
  };
}
