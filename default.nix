let
  pkgs = import <nixpkgs> {};
in
pkgs.buildGoModule rec {
  name = "webshot";
  version = "0.0.1";
  vendorSha256 = "0sh9dcfsan2bqqiwbzm99l10sik326lzv6rra68kvqrayb551qjn";
  src = ./.;
  meta = with pkgs.lib; {
    description = "Render your html files to images";
    homepage = "https://github.com/lucasew/webshot";
    platforms = platforms.linux;
  };
}
