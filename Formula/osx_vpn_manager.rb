class OsxVpnManager < Formula
  desc "Simple command-line vpn manager for OSX"
  homepage "https://github.com/Sjeanpierre/osx_vpn_manager"
  url "https://github.com/Sjeanpierre/osx_vpn_manager/releases/download/0.0.2/osx_vpn_manager-0.0.2.tar.gz"
  sha256 "e0819ee9c0b946b0cd716bf34392c9915674b272250114941e371d746c571728"

  depends_on "macosvpn"

  def install
    bin.install "vpn"
  end

  test do
    assert_equal "0.0.2", `#{bin}/vpn --version`
  end
end
