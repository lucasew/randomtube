with import <nixpkgs> {};
stdenv.mkDerivation {
  name = "environment";
  buildInputs = [
    go
    gopls
    ffmpeg
  ];
}
