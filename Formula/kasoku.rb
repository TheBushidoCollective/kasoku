class Kasoku < Formula
  desc "Accelerate your build times with intelligent command caching"
  homepage "https://kasoku.dev"
  url "https://github.com/thebushidocollective/brisk/archive/refs/tags/v0.1.0.tar.gz"
  sha256 "" # Will be calculated when we create a release
  license "MIT"
  head "https://github.com/thebushidocollective/brisk.git", branch: "main"

  depends_on "go" => :build

  def install
    # Build the CLI
    cd "cmd/kasoku" do
      system "go", "build", *std_go_args(ldflags: "-s -w -X main.version=#{version}"), "."
    end

    # Install shell completion scripts
    generate_completions_from_executable(bin/"kasoku", "completion")
  end

  test do
    # Test version command
    assert_match "kasoku version #{version}", shell_output("#{bin}/kasoku version")

    # Test basic help
    assert_match "Accelerate your builds", shell_output("#{bin}/kasoku --help")
  end
end
