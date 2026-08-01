{ lib, buildGoModule, fetchFromGitHub }:

buildGoModule rec {
  pname = "justssh";
  version = "1.0.1";

  src = fetchFromGitHub {
    owner = "luynrs";
    repo = "justssh";
    rev = "v${version}";
    hash = "sha256-BtcvrbaGQa07HZEJN31fF53HNdPwCLfOMi3NY3hA5Ec=";
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
