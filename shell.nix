{pkgs ? import <nixpkgs> {}}:
with pkgs;
stdenv.mkDerivation {
  name = "environment";
  buildInputs = [
    go
    gopls
    ffmpeg
  ];
}
