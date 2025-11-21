import Link from 'next/link'
import { Button } from '@/components/ui/button'

export default function DocsPage() {
  return (
    <div className="min-h-screen bg-white">
      {/* Header */}
      <header className="border-b bg-white sticky top-0 z-50">
        <div className="container mx-auto px-4 py-4 flex items-center justify-between">
          <Link href="/" className="flex items-center gap-2">
            <div className="text-xl font-bold">Kasoku</div>
            <span className="text-sm text-muted-foreground">加速</span>
          </Link>
          <Link href="/">
            <Button variant="ghost">Back to Home</Button>
          </Link>
        </div>
      </header>

      <div className="container mx-auto px-4 py-12">
        <div className="grid lg:grid-cols-[250px_1fr] gap-12">
          {/* Sidebar Navigation */}
          <aside className="lg:sticky lg:top-24 h-fit">
            <nav className="space-y-1">
              <div className="font-semibold mb-3 text-sm text-muted-foreground">GETTING STARTED</div>
              <a href="#installation" className="block px-3 py-2 rounded hover:bg-gray-100 text-sm">Installation</a>
              <a href="#quick-start" className="block px-3 py-2 rounded hover:bg-gray-100 text-sm">Quick Start</a>

              <div className="font-semibold mb-3 mt-6 text-sm text-muted-foreground">CLI USAGE</div>
              <a href="#basic-usage" className="block px-3 py-2 rounded hover:bg-gray-100 text-sm">Basic Usage</a>
              <a href="#patterns" className="block px-3 py-2 rounded hover:bg-gray-100 text-sm">Cache Patterns</a>
              <a href="#commands" className="block px-3 py-2 rounded hover:bg-gray-100 text-sm">Commands</a>

              <div className="font-semibold mb-3 mt-6 text-sm text-muted-foreground">REMOTE CACHING</div>
              <a href="#authentication" className="block px-3 py-2 rounded hover:bg-gray-100 text-sm">Authentication</a>
              <a href="#remote-setup" className="block px-3 py-2 rounded hover:bg-gray-100 text-sm">Setup</a>
              <a href="#teams" className="block px-3 py-2 rounded hover:bg-gray-100 text-sm">Team Caching</a>

              <div className="font-semibold mb-3 mt-6 text-sm text-muted-foreground">CI/CD</div>
              <a href="#github-actions" className="block px-3 py-2 rounded hover:bg-gray-100 text-sm">GitHub Actions</a>
              <a href="#gitlab-ci" className="block px-3 py-2 rounded hover:bg-gray-100 text-sm">GitLab CI</a>
              <a href="#other-ci" className="block px-3 py-2 rounded hover:bg-gray-100 text-sm">Other Platforms</a>

              <div className="font-semibold mb-3 mt-6 text-sm text-muted-foreground">SELF-HOSTING</div>
              <a href="#docker-compose" className="block px-3 py-2 rounded hover:bg-gray-100 text-sm">Docker Compose</a>
              <a href="#kubernetes" className="block px-3 py-2 rounded hover:bg-gray-100 text-sm">Kubernetes</a>

              <div className="font-semibold mb-3 mt-6 text-sm text-muted-foreground">REFERENCE</div>
              <a href="#configuration" className="block px-3 py-2 rounded hover:bg-gray-100 text-sm">Configuration</a>
              <a href="#api" className="block px-3 py-2 rounded hover:bg-gray-100 text-sm">API Reference</a>
            </nav>
          </aside>

          {/* Main Content */}
          <main className="prose prose-gray max-w-none">
            <h1 className="text-4xl font-bold mb-2">Documentation</h1>
            <p className="text-xl text-muted-foreground mb-8">
              Everything you need to accelerate your build times with Kasoku.
            </p>

            {/* Installation */}
            <section id="installation" className="mb-12">
              <h2 className="text-3xl font-bold mb-4 border-b pb-2">Installation</h2>

              <h3 className="text-xl font-semibold mb-3">macOS / Linux</h3>
              <div className="bg-gray-900 text-white rounded-lg p-4 font-mono text-sm mb-4">
                <code>curl -fsSL https://kasoku.dev/install.sh | sh</code>
              </div>

              <h3 className="text-xl font-semibold mb-3">Using Go</h3>
              <div className="bg-gray-900 text-white rounded-lg p-4 font-mono text-sm mb-4">
                <code>go install github.com/thebushidocollective/kasoku/cmd/kasoku@latest</code>
              </div>

              <h3 className="text-xl font-semibold mb-3">Using Homebrew</h3>
              <div className="bg-gray-900 text-white rounded-lg p-4 font-mono text-sm mb-4">
                <div className="text-gray-400"># Tap the Bushido Collective repository</div>
                <div>brew tap thebushidocollective/kasoku</div>
                <div className="mt-2 text-gray-400"># Install Kasoku</div>
                <div>brew install kasoku</div>
              </div>
              <p className="text-sm text-muted-foreground mb-4">Or in one command:</p>
              <div className="bg-gray-900 text-white rounded-lg p-4 font-mono text-sm mb-4">
                <code>brew install thebushidocollective/kasoku/kasoku</code>
              </div>

              <h3 className="text-xl font-semibold mb-3">Verify Installation</h3>
              <div className="bg-gray-900 text-white rounded-lg p-4 font-mono text-sm">
                <div>$ kasoku --version</div>
                <div className="text-green-400">kasoku version 0.1.0</div>
              </div>
            </section>

            {/* Quick Start */}
            <section id="quick-start" className="mb-12">
              <h2 className="text-3xl font-bold mb-4 border-b pb-2">Quick Start</h2>

              <p className="mb-4">
                Get started with Kasoku in under a minute. No configuration required!
              </p>

              <h3 className="text-xl font-semibold mb-3">1. Run Your First Cached Command</h3>
              <div className="bg-gray-900 text-white rounded-lg p-4 font-mono text-sm mb-4">
                <div className="text-gray-400"># First run executes normally</div>
                <div>$ kasoku exec go build -o myapp</div>
                <div className="text-yellow-400 mt-2">⏱  Building... (took 30.2s)</div>
                <div className="text-green-400">✓ Cache entry created: abc123def456</div>
              </div>

              <h3 className="text-xl font-semibold mb-3">2. Re-run Without Changes</h3>
              <div className="bg-gray-900 text-white rounded-lg p-4 font-mono text-sm mb-4">
                <div className="text-gray-400"># Subsequent runs are instant!</div>
                <div>$ kasoku exec go build -o myapp</div>
                <div className="text-green-400 mt-2">✨ Cache hit! Restored in 43ms</div>
              </div>

              <h3 className="text-xl font-semibold mb-3">3. Make a Change</h3>
              <div className="bg-gray-900 text-white rounded-lg p-4 font-mono text-sm mb-4">
                <div className="text-gray-400"># Edit a source file</div>
                <div>$ echo &apos;// comment&apos; &gt;&gt; main.go</div>
                <div className="mt-2">$ kasoku exec go build -o myapp</div>
                <div className="text-yellow-400 mt-2">⏱  Cache miss, rebuilding... (took 5.1s)</div>
                <div className="text-green-400">✓ Cache entry created: def456abc789</div>
              </div>
            </section>

            {/* Basic Usage */}
            <section id="basic-usage" className="mb-12">
              <h2 className="text-3xl font-bold mb-4 border-b pb-2">Basic Usage</h2>

              <h3 className="text-xl font-semibold mb-3">The <code>exec</code> Command</h3>
              <p className="mb-4">
                Wrap any command with <code>kasoku exec</code> to enable caching:
              </p>
              <div className="bg-gray-900 text-white rounded-lg p-4 font-mono text-sm mb-4">
                <div>kasoku exec [flags] &lt;command&gt;</div>
              </div>

              <h3 className="text-xl font-semibold mb-3">Common Examples</h3>
              <div className="bg-gray-900 text-white rounded-lg p-4 font-mono text-sm mb-4">
                <div className="text-gray-400"># Go builds</div>
                <div>kasoku exec go build ./...</div>
                <div className="text-gray-400 mt-3"># npm/Node.js</div>
                <div>kasoku exec npm run build</div>
                <div className="text-gray-400 mt-3"># Rust</div>
                <div>kasoku exec cargo build --release</div>
                <div className="text-gray-400 mt-3"># Python</div>
                <div>kasoku exec python -m pytest</div>
                <div className="text-gray-400 mt-3"># Maven</div>
                <div>kasoku exec mvn clean package</div>
                <div className="text-gray-400 mt-3"># Make</div>
                <div>kasoku exec make all</div>
              </div>
            </section>

            {/* Cache Patterns */}
            <section id="patterns" className="mb-12">
              <h2 className="text-3xl font-bold mb-4 border-b pb-2">Cache Patterns</h2>

              <p className="mb-4">
                Kasoku automatically detects which files to track, but you can customize this behavior:
              </p>

              <h3 className="text-xl font-semibold mb-3">Using <code>--pattern</code></h3>
              <div className="bg-gray-900 text-white rounded-lg p-4 font-mono text-sm mb-4">
                <div className="text-gray-400"># Only track Go source files</div>
                <div>kasoku exec --pattern &apos;**/*.go&apos; go build</div>
                <div className="text-gray-400 mt-3"># Track multiple patterns</div>
                <div>kasoku exec --pattern &apos;src/**/*.ts&apos; --pattern &apos;package.json&apos; npm run build</div>
              </div>

              <h3 className="text-xl font-semibold mb-3">Environment Variables</h3>
              <div className="bg-gray-900 text-white rounded-lg p-4 font-mono text-sm mb-4">
                <div className="text-gray-400"># Include environment variables in cache key</div>
                <div>kasoku exec --env NODE_ENV --env API_URL npm run build</div>
              </div>

              <h3 className="text-xl font-semibold mb-3">Excluding Files</h3>
              <div className="bg-gray-900 text-white rounded-lg p-4 font-mono text-sm mb-4">
                <div className="text-gray-400"># Exclude test files from cache key</div>
                <div>kasoku exec --exclude &apos;**/*_test.go&apos; go build</div>
              </div>
            </section>

            {/* Commands */}
            <section id="commands" className="mb-12">
              <h2 className="text-3xl font-bold mb-4 border-b pb-2">Commands</h2>

              <h3 className="text-xl font-semibold mb-3"><code>kasoku exec</code></h3>
              <p className="mb-2">Execute a command with caching.</p>
              <div className="bg-gray-50 rounded p-3 mb-4 text-sm">
                <strong>Usage:</strong> <code>kasoku exec [flags] &lt;command&gt;</code>
              </div>

              <h3 className="text-xl font-semibold mb-3"><code>kasoku list</code></h3>
              <p className="mb-2">List all cached entries.</p>
              <div className="bg-gray-900 text-white rounded-lg p-4 font-mono text-sm mb-4">
                <div>$ kasoku list</div>
                <div className="text-gray-400 mt-2">abc123  go build        2.1 MB  2 hours ago</div>
                <div className="text-gray-400">def456  npm run build   15.3 MB 1 day ago</div>
              </div>

              <h3 className="text-xl font-semibold mb-3"><code>kasoku clean</code></h3>
              <p className="mb-2">Clean local cache.</p>
              <div className="bg-gray-900 text-white rounded-lg p-4 font-mono text-sm mb-4">
                <div className="text-gray-400"># Clean all</div>
                <div>$ kasoku clean</div>
                <div className="text-gray-400 mt-3"># Clean specific entry</div>
                <div>$ kasoku clean abc123</div>
              </div>

              <h3 className="text-xl font-semibold mb-3"><code>kasoku stats</code></h3>
              <p className="mb-2">View cache statistics.</p>
              <div className="bg-gray-900 text-white rounded-lg p-4 font-mono text-sm mb-4">
                <div>$ kasoku stats</div>
                <div className="text-green-400 mt-2">Cache hits: 142 (87%)</div>
                <div className="text-yellow-400">Cache misses: 21 (13%)</div>
                <div className="text-blue-400">Time saved: 2.3 hours</div>
              </div>
            </section>

            {/* Authentication */}
            <section id="authentication" className="mb-12">
              <h2 className="text-3xl font-bold mb-4 border-b pb-2">Authentication</h2>

              <h3 className="text-xl font-semibold mb-3">Login to Kasoku Cloud</h3>
              <div className="bg-gray-900 text-white rounded-lg p-4 font-mono text-sm mb-4">
                <div>$ kasoku login</div>
                <div className="text-blue-400 mt-2">→ Opening browser for authentication...</div>
                <div className="text-green-400 mt-1">✓ Successfully logged in as user@example.com</div>
              </div>

              <h3 className="text-xl font-semibold mb-3">Using API Tokens</h3>
              <p className="mb-2">For CI/CD, use API tokens instead of interactive login:</p>
              <div className="bg-gray-900 text-white rounded-lg p-4 font-mono text-sm mb-4">
                <div className="text-gray-400"># Set token via environment variable</div>
                <div>export KASOKU_TOKEN=ksk_abc123...</div>
                <div className="text-gray-400 mt-3"># Or use --token flag</div>
                <div>kasoku exec --token ksk_abc123... go build</div>
              </div>

              <h3 className="text-xl font-semibold mb-3">Self-Hosted Server</h3>
              <div className="bg-gray-900 text-white rounded-lg p-4 font-mono text-sm mb-4">
                <div className="text-gray-400"># Configure custom server URL</div>
                <div>export KASOKU_SERVER=https://kasoku.yourcompany.com</div>
                <div className="mt-2">$ kasoku login</div>
              </div>
            </section>

            {/* Remote Setup */}
            <section id="remote-setup" className="mb-12">
              <h2 className="text-3xl font-bold mb-4 border-b pb-2">Remote Caching Setup</h2>

              <p className="mb-4">
                Enable remote caching to share cache across your team and CI/CD pipelines.
              </p>

              <h3 className="text-xl font-semibold mb-3">1. Enable Remote Caching</h3>
              <div className="bg-gray-900 text-white rounded-lg p-4 font-mono text-sm mb-4">
                <div>$ kasoku exec --remote go build</div>
              </div>

              <h3 className="text-xl font-semibold mb-3">2. Set as Default</h3>
              <div className="bg-gray-900 text-white rounded-lg p-4 font-mono text-sm mb-4">
                <div className="text-gray-400"># Always use remote caching</div>
                <div>export KASOKU_REMOTE=true</div>
              </div>

              <h3 className="text-xl font-semibold mb-3">Cache Priority</h3>
              <p className="mb-2">Kasoku checks caches in this order:</p>
              <ol className="list-decimal ml-6 mb-4">
                <li>Local cache (fastest)</li>
                <li>Remote cache (team-shared)</li>
                <li>Execute command (cache miss)</li>
              </ol>
            </section>

            {/* Team Caching */}
            <section id="teams" className="mb-12">
              <h2 className="text-3xl font-bold mb-4 border-b pb-2">Team Caching</h2>

              <p className="mb-4">
                Share cache entries across your entire team for maximum efficiency.
              </p>

              <h3 className="text-xl font-semibold mb-3">Create a Team</h3>
              <div className="bg-gray-50 border rounded p-4 mb-4">
                <p className="mb-2">1. Go to <Link href="/dashboard" className="text-blue-600 hover:underline">Dashboard</Link></p>
                <p className="mb-2">2. Click &quot;Create Team&quot;</p>
                <p>3. Invite team members via email</p>
              </div>

              <h3 className="text-xl font-semibold mb-3">Team Cache Benefits</h3>
              <ul className="list-disc ml-6 mb-4">
                <li>One team member builds, everyone benefits</li>
                <li>CI/CD cache shared with local development</li>
                <li>Parallel CI jobs avoid duplicate work</li>
                <li>50GB shared storage on Team plan</li>
              </ul>
            </section>

            {/* GitHub Actions */}
            <section id="github-actions" className="mb-12">
              <h2 className="text-3xl font-bold mb-4 border-b pb-2">GitHub Actions</h2>

              <p className="mb-4">
                Integrate Kasoku into your GitHub Actions workflows for faster CI builds.
              </p>

              <div className="bg-gray-900 text-white rounded-lg p-4 font-mono text-sm mb-4 overflow-x-auto">
                <pre className="text-xs">{`name: Build

on: [push, pull_request]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Install Kasoku
        run: |
          curl -fsSL https://kasoku.dev/install.sh | sh
          echo "$HOME/.kasoku/bin" >> $GITHUB_PATH

      - name: Build with Kasoku
        env:
          KASOKU_TOKEN: \${{ secrets.KASOKU_TOKEN }}
          KASOKU_REMOTE: true
        run: kasoku exec go build ./...`}</pre>
              </div>

              <div className="bg-blue-50 border-l-4 border-blue-500 p-4 mb-4">
                <p className="font-semibold">Pro Tip:</p>
                <p className="text-sm mt-1">
                  Store your KASOKU_TOKEN in GitHub Secrets for secure access.
                </p>
              </div>
            </section>

            {/* GitLab CI */}
            <section id="gitlab-ci" className="mb-12">
              <h2 className="text-3xl font-bold mb-4 border-b pb-2">GitLab CI</h2>

              <div className="bg-gray-900 text-white rounded-lg p-4 font-mono text-sm mb-4 overflow-x-auto">
                <pre className="text-xs">{`build:
  image: golang:1.21
  before_script:
    - curl -fsSL https://kasoku.dev/install.sh | sh
    - export PATH="$HOME/.kasoku/bin:$PATH"
  script:
    - kasoku exec go build ./...
  variables:
    KASOKU_TOKEN: $KASOKU_TOKEN
    KASOKU_REMOTE: "true"`}</pre>
              </div>

              <p className="mb-4">
                Add KASOKU_TOKEN to GitLab CI/CD Variables for secure access.
              </p>
            </section>

            {/* Other CI */}
            <section id="other-ci" className="mb-12">
              <h2 className="text-3xl font-bold mb-4 border-b pb-2">Other CI/CD Platforms</h2>

              <p className="mb-4">
                Kasoku works with any CI/CD platform. General setup pattern:
              </p>

              <ol className="list-decimal ml-6 mb-4">
                <li className="mb-2">Install Kasoku CLI during CI setup</li>
                <li className="mb-2">Set KASOKU_TOKEN environment variable (from CI secrets)</li>
                <li className="mb-2">Set KASOKU_REMOTE=true to enable remote caching</li>
                <li className="mb-2">Wrap build commands with <code>kasoku exec</code></li>
              </ol>

              <h3 className="text-xl font-semibold mb-3">CircleCI Example</h3>
              <div className="bg-gray-900 text-white rounded-lg p-4 font-mono text-sm mb-4 overflow-x-auto">
                <pre className="text-xs">{`version: 2.1
jobs:
  build:
    docker:
      - image: cimg/go:1.21
    steps:
      - checkout
      - run: curl -fsSL https://kasoku.dev/install.sh | sh
      - run: echo 'export PATH="$HOME/.kasoku/bin:$PATH"' >> $BASH_ENV
      - run: kasoku exec go build ./...
    environment:
      KASOKU_REMOTE: true`}</pre>
              </div>

              <h3 className="text-xl font-semibold mb-3">Jenkins Example</h3>
              <div className="bg-gray-900 text-white rounded-lg p-4 font-mono text-sm mb-4 overflow-x-auto">
                <pre className="text-xs">{`pipeline {
  agent any
  environment {
    KASOKU_TOKEN = credentials('kasoku-token')
    KASOKU_REMOTE = 'true'
  }
  stages {
    stage('Build') {
      steps {
        sh 'curl -fsSL https://kasoku.dev/install.sh | sh'
        sh 'export PATH="$HOME/.kasoku/bin:$PATH"'
        sh 'kasoku exec go build ./...'
      }
    }
  }
}`}</pre>
              </div>
            </section>

            {/* Docker Compose */}
            <section id="docker-compose" className="mb-12">
              <h2 className="text-3xl font-bold mb-4 border-b pb-2">Self-Hosting with Docker Compose</h2>

              <p className="mb-4">
                Deploy Kasoku on your own infrastructure for full control and privacy.
              </p>

              <h3 className="text-xl font-semibold mb-3">Quick Start</h3>
              <div className="bg-gray-900 text-white rounded-lg p-4 font-mono text-sm mb-4">
                <div className="text-gray-400"># Clone repository</div>
                <div>git clone https://github.com/thebushidocollective/brisk.git</div>
                <div>cd brisk/deployment/docker-compose</div>
                <div className="mt-3 text-gray-400"># Run setup script</div>
                <div>./scripts/setup.sh</div>
                <div className="mt-3 text-gray-400"># Your server is now running!</div>
                <div className="text-green-400">✓ Kasoku server running at https://kasoku.example.com</div>
              </div>

              <h3 className="text-xl font-semibold mb-3">What&apos;s Included</h3>
              <ul className="list-disc ml-6 mb-4">
                <li>Kasoku server with authentication</li>
                <li>PostgreSQL database</li>
                <li>Caddy reverse proxy with auto-SSL</li>
                <li>Web dashboard</li>
                <li>Automated backups</li>
              </ul>

              <div className="bg-yellow-50 border-l-4 border-yellow-500 p-4 mb-4">
                <p className="font-semibold">Requirements:</p>
                <ul className="text-sm mt-1 list-disc ml-4">
                  <li>Docker & Docker Compose</li>
                  <li>Domain name pointing to your server</li>
                  <li>Ports 80 and 443 open</li>
                </ul>
              </div>
            </section>

            {/* Kubernetes */}
            <section id="kubernetes" className="mb-12">
              <h2 className="text-3xl font-bold mb-4 border-b pb-2">Kubernetes Deployment</h2>

              <p className="mb-4">
                Production-ready Helm chart for Kubernetes deployments.
              </p>

              <h3 className="text-xl font-semibold mb-3">Install with Helm</h3>
              <div className="bg-gray-900 text-white rounded-lg p-4 font-mono text-sm mb-4">
                <div className="text-gray-400"># Add Helm repo</div>
                <div>helm repo add kasoku https://charts.kasoku.dev</div>
                <div>helm repo update</div>
                <div className="mt-3 text-gray-400"># Install</div>
                <div>helm install kasoku kasoku/kasoku \</div>
                <div>  --set global.domain=kasoku.example.com \</div>
                <div>  --set secrets.jwtSecret=$(openssl rand -hex 32)</div>
              </div>

              <h3 className="text-xl font-semibold mb-3">Features</h3>
              <ul className="list-disc ml-6 mb-4">
                <li>High availability with multiple replicas</li>
                <li>Horizontal pod autoscaling</li>
                <li>Built-in or external PostgreSQL</li>
                <li>S3-compatible storage support</li>
                <li>Ingress with SSL termination</li>
                <li>Prometheus metrics</li>
              </ul>

              <h3 className="text-xl font-semibold mb-3">Production Configuration</h3>
              <div className="bg-gray-900 text-white rounded-lg p-4 font-mono text-sm mb-4 overflow-x-auto">
                <pre className="text-xs">{`helm install kasoku kasoku/kasoku \\
  --values values-production.yaml \\
  --set global.domain=kasoku.example.com \\
  --set server.replicaCount=3 \\
  --set storage.type=s3 \\
  --set storage.s3.bucket=kasoku-cache`}</pre>
              </div>
            </section>

            {/* Configuration */}
            <section id="configuration" className="mb-12">
              <h2 className="text-3xl font-bold mb-4 border-b pb-2">Configuration</h2>

              <h3 className="text-xl font-semibold mb-3">Environment Variables</h3>
              <div className="overflow-x-auto">
                <table className="min-w-full border text-sm">
                  <thead className="bg-gray-50">
                    <tr>
                      <th className="border px-4 py-2 text-left">Variable</th>
                      <th className="border px-4 py-2 text-left">Description</th>
                      <th className="border px-4 py-2 text-left">Default</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr>
                      <td className="border px-4 py-2"><code>KASOKU_SERVER</code></td>
                      <td className="border px-4 py-2">Server URL</td>
                      <td className="border px-4 py-2">https://api.kasoku.dev</td>
                    </tr>
                    <tr>
                      <td className="border px-4 py-2"><code>KASOKU_TOKEN</code></td>
                      <td className="border px-4 py-2">API token</td>
                      <td className="border px-4 py-2">-</td>
                    </tr>
                    <tr>
                      <td className="border px-4 py-2"><code>KASOKU_REMOTE</code></td>
                      <td className="border px-4 py-2">Enable remote caching</td>
                      <td className="border px-4 py-2">false</td>
                    </tr>
                    <tr>
                      <td className="border px-4 py-2"><code>KASOKU_CACHE_DIR</code></td>
                      <td className="border px-4 py-2">Local cache directory</td>
                      <td className="border px-4 py-2">~/.kasoku/cache</td>
                    </tr>
                  </tbody>
                </table>
              </div>
            </section>

            {/* API Reference */}
            <section id="api" className="mb-12">
              <h2 className="text-3xl font-bold mb-4 border-b pb-2">API Reference</h2>

              <p className="mb-4">
                Kasoku server provides a REST API for programmatic access.
              </p>

              <h3 className="text-xl font-semibold mb-3">Authentication</h3>
              <div className="bg-gray-900 text-white rounded-lg p-4 font-mono text-sm mb-4">
                <div>Authorization: Bearer &lt;token&gt;</div>
              </div>

              <h3 className="text-xl font-semibold mb-3">Cache Endpoints</h3>

              <div className="mb-6">
                <h4 className="font-semibold mb-2">PUT /cache/:hash</h4>
                <p className="text-sm mb-2">Upload a cache entry</p>
                <div className="bg-gray-50 rounded p-3 text-sm">
                  <div><strong>Headers:</strong> Authorization, Content-Length</div>
                  <div><strong>Body:</strong> gzipped tarball</div>
                  <div><strong>Response:</strong> 201 Created</div>
                </div>
              </div>

              <div className="mb-6">
                <h4 className="font-semibold mb-2">GET /cache/:hash</h4>
                <p className="text-sm mb-2">Download a cache entry</p>
                <div className="bg-gray-50 rounded p-3 text-sm">
                  <div><strong>Response:</strong> gzipped tarball or 404 Not Found</div>
                </div>
              </div>

              <div className="mb-6">
                <h4 className="font-semibold mb-2">GET /cache</h4>
                <p className="text-sm mb-2">List cache entries</p>
                <div className="bg-gray-50 rounded p-3 text-sm">
                  <div><strong>Query params:</strong> limit, offset</div>
                  <div><strong>Response:</strong> Array of cache entries</div>
                </div>
              </div>

              <div className="mb-6">
                <h4 className="font-semibold mb-2">GET /analytics</h4>
                <p className="text-sm mb-2">Get cache analytics</p>
                <div className="bg-gray-50 rounded p-3 text-sm">
                  <div><strong>Response:</strong> Stats including hits, misses, time saved</div>
                </div>
              </div>

              <h3 className="text-xl font-semibold mb-3 mt-6">Job Coordination</h3>

              <div className="mb-6">
                <h4 className="font-semibold mb-2">POST /jobs/register/:hash</h4>
                <p className="text-sm mb-2">Register a new job</p>
                <div className="bg-gray-50 rounded p-3 text-sm">
                  <div><strong>Response:</strong> job_id, is_primary</div>
                </div>
              </div>

              <div className="mb-6">
                <h4 className="font-semibold mb-2">GET /jobs/wait/:hash</h4>
                <p className="text-sm mb-2">Wait for job completion (long-poll)</p>
                <div className="bg-gray-50 rounded p-3 text-sm">
                  <div><strong>Response:</strong> Job status (complete/failed)</div>
                </div>
              </div>
            </section>

            {/* Footer */}
            <div className="mt-16 pt-8 border-t text-center text-sm text-muted-foreground">
              <p className="mb-2">
                Need help? <Link href="/support" className="text-blue-600 hover:underline">Contact Support</Link>
              </p>
              <p>
                Kasoku © The Bushido Collective
              </p>
            </div>
          </main>
        </div>
      </div>
    </div>
  )
}
