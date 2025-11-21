class KasokuServer < Formula
  desc "Kasoku cache server for team and remote caching"
  homepage "https://kasoku.dev"
  url "https://github.com/thebushidocollective/brisk/archive/refs/tags/v0.1.0.tar.gz"
  sha256 "" # Will be calculated when we create a release
  license "MIT"
  head "https://github.com/thebushidocollective/brisk.git", branch: "main"

  depends_on "go" => :build

  def install
    # Build the server
    cd "server/cmd/kasoku-server" do
      system "go", "build", *std_go_args(ldflags: "-s -w -X main.version=#{version}"), "."
    end

    # Install example configuration
    (etc/"kasoku").install "kasoku.yaml" => "kasoku.example.yaml" if File.exist?("kasoku.yaml")
  end

  service do
    run [opt_bin/"kasoku-server"]
    keep_alive true
    working_dir var/"kasoku"
    log_path var/"log/kasoku-server.log"
    error_log_path var/"log/kasoku-server.log"
    environment_variables PATH: std_service_path_env
  end

  def post_install
    (var/"kasoku").mkpath
    (var/"log").mkpath
  end

  test do
    # Test server starts
    system bin/"kasoku-server", "--help"
  end
end
