{
  description = "JustSSH - minimal SSH launcher";

  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";

  outputs = { self, nixpkgs }:
    let
      system = "x86_64-linux";
      pkgs = nixpkgs.legacyPackages.${system};
    in {
      packages.${system}.default = pkgs.buildGoModule {
        pname = "justssh";
        version = "1.0.2";
        src = ./.;
        vendorHash = "sha256-aJllcMJduoi8VBWMJWsxm8swXtNonYZzX8etmNZePzc=";
        subPackages = [ "cmd/jssh" ];

        meta = {
          description = "Minimal SSH launcher for the terminal";
          mainProgram = "jssh";
        };
      };

      devShells.${system}.default = pkgs.mkShell {
        packages = [ pkgs.go pkgs.gofumpt pkgs.golangci-lint ];
      };
    };
}
