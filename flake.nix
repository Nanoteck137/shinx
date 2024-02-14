{
  description = "Simple auth server for self-hosted services";

  inputs = {
    nixpkgs.url      = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url  = "github:numtide/flake-utils";
    templ.url        = "github:a-h/templ";
  };

  outputs = { self, nixpkgs, flake-utils, ... }@inputs:
    flake-utils.lib.eachDefaultSystem (system:
      let
        overlays = [];
        pkgs = import nixpkgs {
          inherit system overlays;
        };

        templ = inputs.templ.packages.${system}.templ;
      in
      {
        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
            air
            templ
          ];
        };
      }
    );
}
