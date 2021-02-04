{pkgs ? import <nixpkgs> {}, ...}:
pkgs.buildGoModule rec {
  name = "randomtube";
  version = "0.0.1";
  vendorSha256 = "sha256-jFzgIDqWhYiS1LfWCp1Oofm+5qNQV/rg7dBYUbT5jU4=";
  src = ./.;
  buildInputs = with pkgs; [
    ffmpeg
  ];
  meta = with pkgs.lib; {
    description = "Generate videos to youtube from a telegram group";
    homepage = "https://github.com/lucasew/randomtube";
    platforms = platforms.linux;
  };
}
