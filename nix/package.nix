{ lib, buildGoModule, fetchFromGitHub }:

buildGoModule rec {
  pname = "justssh";
  version = "1.0.2";

  src = fetchFromGitHub {
    owner = "luynrs";
    repo = "justssh";
    rev = "v${version}";
    hash = "sha256-ZsRILOdd1DfjxNuFmGOz4lId+pggqZw1edZaUblABGg=";
  };

  vendorHash = "sha256-aJllcMJduoi8VBWMJWsxm8swXtNonYZzX8etmNZePzc=";

  subPackages = [ "cmd/jssh" ];

  meta = {
    description = "Minimal SSH launcher for the terminal";
    homepage = "https://github.com/luynrs/justssh";
    license = lib.licenses.mit;
    mainProgram = "jssh";
  };
}
