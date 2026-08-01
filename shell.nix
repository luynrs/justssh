{ pkgs ? import <nixpkgs> { } }:

pkgs.mkShell {
  packages = with pkgs; [
    go
    gofumpt
    golangci-lint
  ];
}
