{ pkgs ? import <nixpkgs> { } }:

pkgs.mkShell {
  name = "nymeria.go";

  nativeBuildInputs = with pkgs; [
    # Go itself for building and development.
    go_1_18

    # Go related tooling for development.
    gotools
    go-tools
    gopls
  ];

  shellHook = ''
    export GOPATH="$HOME/.go"
    export PATH="$PATH:$HOME/.go/bin"
  '';
}
