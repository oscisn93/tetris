#
# NixOS game dev flake by jchw 
# https://discourse.nixos.org/t/how-to-build-a-golang-program-which-needs-system-libs/45682/14?u=uskeutia
#
# Use:
# if flake.nix is inside a Git repo, but not added:
# nix develop path://$PWD
#
# otherwise, just:
# nix develop

{
  inputs.flake-utils.url = "github:numtide/flake-utils";
  outputs =
    { nixpkgs, flake-utils, ... }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        devShells.default = pkgs.mkShell {
          nativeBuildInputs = with pkgs; [
            go
            libGL
            xorg.libX11
            xorg.libXrandr
            xorg.libXcursor
            xorg.libXinerama
            xorg.libXi
            xorg.libXxf86vm
          ];

          shellHook = ''
            export LD_LIBRARY_PATH=${pkgs.wayland}/lib:${pkgs.lib.getLib pkgs.libGL}/lib:${pkgs.lib.getLib pkgs.libGL}/lib:$LD_LIBRARY_PATH
          '';
        };
      }
    );
}
