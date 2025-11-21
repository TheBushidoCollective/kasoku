'use client'

import { useEffect } from 'react'
import { useRouter } from 'next/navigation'
import Link from 'next/link'
import { Button } from '@/components/ui/button'
import { StatCounter } from '@/components/stats-counter'
import { LiveCounter } from '@/components/live-counter'
import { ShellTabs } from '@/components/shell-tabs'

export default function LandingPage() {
  const router = useRouter()

  useEffect(() => {
    // If self-hosted, redirect to login/dashboard
    const isSelfHosted = process.env.NEXT_PUBLIC_SELF_HOSTED === 'true'
    if (isSelfHosted) {
      // Check if user is logged in
      const token = localStorage.getItem('kasoku_token')
      if (token) {
        router.push('/dashboard')
      } else {
        router.push('/login')
      }
    }
  }, [router])

  return (
    <div className="min-h-screen bg-gradient-to-b from-gray-50 to-white">
      {/* Header */}
      <header className="container mx-auto px-4 py-6 flex items-center justify-between">
        <div className="flex items-center gap-2">
          <div className="text-2xl font-bold">Kasoku</div>
          <span className="text-sm text-muted-foreground">加速</span>
        </div>
        <nav className="flex items-center gap-4">
          <Link href="/login">
            <Button variant="ghost">Login</Button>
          </Link>
          <Link href="/signup">
            <Button>Sign Up</Button>
          </Link>
        </nav>
      </header>

      {/* Hero */}
      <section className="container mx-auto px-4 py-20 text-center">
        <h1 className="text-5xl md:text-6xl font-bold mb-6">
          Accelerate Your Build Times,<br />Save Time, Save on Compute
        </h1>
        <p className="text-xl text-muted-foreground mb-8 max-w-2xl mx-auto">
          Kasoku makes rebuilding, retesting, and rerunning commands instant with intelligent caching.
          Works with any build tool, language, or command.
        </p>
        <div className="flex gap-4 justify-center">
          <Link href="/signup">
            <Button size="lg">Get Started Free</Button>
          </Link>
          <Link href="/docs">
            <Button size="lg" variant="outline">View Docs</Button>
          </Link>
        </div>

        {/* Quick example */}
        <div className="mt-12 max-w-3xl mx-auto">
          <div className="bg-gray-900 text-white rounded-lg p-6 text-left font-mono text-sm">
            <div className="text-gray-400"># First run: ~30s</div>
            <div className="text-green-400">$ kasoku exec go build -o myapp</div>
            <div className="mt-2 text-gray-400"># Subsequent runs with no changes: instant!</div>
            <div className="text-green-400">$ kasoku exec go build -o myapp</div>
            <div className="text-blue-400">✨ Cache hit! Restored in 43ms</div>
          </div>
        </div>
      </section>

      {/* Shell Integration */}
      <section className="container mx-auto px-4 py-20 bg-white">
        <h2 className="text-3xl font-bold text-center mb-4">Transparent Shell Integration</h2>
        <p className="text-xl text-muted-foreground text-center mb-12 max-w-2xl mx-auto">
          One line of setup. Every build command automatically cached. No changes to your workflow.
        </p>
        <ShellTabs />
      </section>

      {/* Statistics */}
      <section className="container mx-auto px-4 py-16">
        <div className="max-w-6xl mx-auto bg-primary text-primary-foreground rounded-2xl p-12">
          <div className="grid md:grid-cols-4 gap-8 text-center">
            <div>
              <div className="text-4xl md:text-5xl font-bold mb-2">
                <LiveCounter initialValue={142857} incrementPerSecond={10} />
              </div>
              <div className="text-sm md:text-base opacity-90">Developer Hours Saved</div>
            </div>
            <div>
              <div className="text-4xl md:text-5xl font-bold mb-2">
                <StatCounter end={2.4} decimals={1} suffix="M" />
              </div>
              <div className="text-sm md:text-base opacity-90">Builds Cached</div>
            </div>
            <div>
              <div className="text-4xl md:text-5xl font-bold mb-2">
                <StatCounter end={89} suffix="%" />
              </div>
              <div className="text-sm md:text-base opacity-90">Average Time Reduction</div>
            </div>
            <div>
              <div className="text-4xl md:text-5xl font-bold mb-2">
                <LiveCounter initialValue={14285} incrementPerSecond={1} prefix="$" />
              </div>
              <div className="text-sm md:text-base opacity-90">CI Compute Dollars Saved</div>
            </div>
          </div>
          <div className="mt-8 text-center text-sm opacity-75">
            Real-time statistics from Kasoku users worldwide · CI costs calculated at $0.10/hour
          </div>
        </div>
      </section>

      {/* Features */}
      <section className="container mx-auto px-4 py-20">
        <h2 className="text-3xl font-bold text-center mb-12">Key Features</h2>
        <div className="grid md:grid-cols-3 gap-8 max-w-6xl mx-auto">
          <div className="text-center p-6">
            <div className="text-5xl mb-4">⚡</div>
            <h3 className="text-xl font-semibold mb-3">Lightning Fast</h3>
            <p className="text-muted-foreground">
              Cache hits restore outputs in milliseconds. Pattern-based matching for automatic caching.
            </p>
          </div>
          <div className="text-center p-6">
            <div className="text-5xl mb-4">🔍</div>
            <h3 className="text-xl font-semibold mb-3">Smart Detection</h3>
            <p className="text-muted-foreground">
              Only re-executes when inputs actually change. Tracks files, globs, and environment variables.
            </p>
          </div>
          <div className="text-center p-6">
            <div className="text-5xl mb-4">🌐</div>
            <h3 className="text-xl font-semibold mb-3">Team Sharing</h3>
            <p className="text-muted-foreground">
              Share cache across teams with remote caching. Self-hosted or cloud options available.
            </p>
          </div>
          <div className="text-center p-6">
            <div className="text-5xl mb-4">🚀</div>
            <h3 className="text-xl font-semibold mb-3">CI/CD Integration</h3>
            <p className="text-muted-foreground">
              Seamlessly works with GitHub Actions, GitLab CI, Jenkins, CircleCI, and any CI/CD platform.
            </p>
          </div>
          <div className="text-center p-6">
            <div className="text-5xl mb-4">🛠️</div>
            <h3 className="text-xl font-semibold mb-3">Language Agnostic</h3>
            <p className="text-muted-foreground">
              Works with any build tool or language: Go, Node, Rust, Python, Java, Maven, Gradle, and more.
            </p>
          </div>
          <div className="text-center p-6">
            <div className="text-5xl mb-4">🔄</div>
            <h3 className="text-xl font-semibold mb-3">Distributed Coordination</h3>
            <p className="text-muted-foreground">
              Multiple jobs wait for the first to complete. No duplicate work across parallel CI builds.
            </p>
          </div>
        </div>
      </section>

      {/* Pricing */}
      <section className="container mx-auto px-4 py-20 bg-gray-50">
        <h2 className="text-3xl font-bold text-center mb-12">Simple Pricing</h2>
        <div className="grid md:grid-cols-3 gap-8 max-w-5xl mx-auto">
          {/* Free */}
          <div className="bg-white rounded-lg p-6 border-2">
            <h3 className="text-xl font-bold mb-2">Free</h3>
            <div className="text-3xl font-bold mb-4">$0</div>
            <ul className="space-y-2 mb-6 text-sm">
              <li>✓ Local caching only</li>
              <li>✓ Unlimited storage</li>
              <li>✓ No signup required</li>
              <li>✓ All core features</li>
            </ul>
            <Button variant="outline" className="w-full">Download CLI</Button>
          </div>

          {/* Individual */}
          <div className="bg-white rounded-lg p-6 border-2 border-primary">
            <div className="text-xs font-semibold text-primary mb-2">MOST POPULAR</div>
            <h3 className="text-xl font-bold mb-2">Individual</h3>
            <div className="text-3xl font-bold mb-4">$5<span className="text-sm text-muted-foreground">/month</span></div>
            <ul className="space-y-2 mb-6 text-sm">
              <li>✓ 2GB remote cache</li>
              <li>✓ FIFO eviction</li>
              <li>✓ CI/CD integration</li>
              <li>✓ Analytics dashboard</li>
              <li>✓ Single user</li>
            </ul>
            <Link href="/signup">
              <Button className="w-full">Start Free Trial</Button>
            </Link>
          </div>

          {/* Team */}
          <div className="bg-white rounded-lg p-6 border-2">
            <h3 className="text-xl font-bold mb-2">Team</h3>
            <div className="text-3xl font-bold mb-4">$20<span className="text-sm text-muted-foreground">/month</span></div>
            <ul className="space-y-2 mb-6 text-sm">
              <li>✓ 50GB remote cache</li>
              <li>✓ Up to 10 users</li>
              <li>✓ Shared team cache</li>
              <li>✓ Priority support</li>
              <li>✓ Advanced analytics</li>
            </ul>
            <Link href="/signup?plan=team">
              <Button variant="outline" className="w-full">Start Free Trial</Button>
            </Link>
          </div>
        </div>

        <div className="text-center mt-8">
          <p className="text-sm text-muted-foreground">
            Self-hosted option available for free. Deploy on your infrastructure.
          </p>
        </div>
      </section>

      {/* Footer */}
      <footer className="container mx-auto px-4 py-8 border-t">
        <div className="flex justify-between items-center">
          <div className="text-sm text-muted-foreground">
            Kasoku © The Bushido Collective
          </div>
          <div className="flex gap-6 text-sm">
            <Link href="/docs" className="hover:underline">Docs</Link>
            <Link href="/blog" className="hover:underline">Blog</Link>
            <Link href="/github" className="hover:underline">GitHub</Link>
          </div>
        </div>
      </footer>
    </div>
  )
}
