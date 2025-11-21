'use client'

import { useEffect, useState } from 'react'
import { useRouter } from 'next/navigation'
import { Button } from '@/components/ui/button'
import { api, type AnalyticsSummary, type CacheEntry, type Token, isAuthenticated } from '@/lib/api'
import { formatBytes, formatDuration, formatDate } from '@/lib/utils'

export default function DashboardPage() {
  const router = useRouter()
  const [analytics, setAnalytics] = useState<AnalyticsSummary | null>(null)
  const [cacheEntries, setCacheEntries] = useState<CacheEntry[]>([])
  const [tokens, setTokens] = useState<Token[]>([])
  const [newTokenName, setNewTokenName] = useState('')
  const [newToken, setNewToken] = useState<string | null>(null)
  const [loading, setLoading] = useState(true)
  const [activeTab, setActiveTab] = useState<'analytics' | 'cache' | 'tokens'>('analytics')

  useEffect(() => {
    if (!isAuthenticated()) {
      router.push('/login')
      return
    }

    loadData()
  }, [router])

  const loadData = async () => {
    try {
      const [analyticsData, cacheData, tokensData] = await Promise.all([
        api.getAnalytics(new Date(Date.now() - 30 * 24 * 60 * 60 * 1000)), // Last 30 days
        api.listCache(50, 0),
        api.listTokens(),
      ])

      setAnalytics(analyticsData)
      setCacheEntries(cacheData.entries)
      setTokens(tokensData.tokens)
    } catch (error) {
      console.error('Failed to load data:', error)
    } finally {
      setLoading(false)
    }
  }

  const handleLogout = async () => {
    await api.logout()
    router.push('/')
  }

  const handleCreateToken = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!newTokenName) return

    try {
      const result = await api.createToken(newTokenName)
      setNewToken(result.token || '')
      setNewTokenName('')
      await loadData()
    } catch (error: any) {
      alert('Failed to create token: ' + error.message)
    }
  }

  const handleRevokeToken = async (id: string) => {
    if (!confirm('Are you sure you want to revoke this token?')) return

    try {
      await api.revokeToken(id)
      await loadData()
    } catch (error: any) {
      alert('Failed to revoke token: ' + error.message)
    }
  }

  const handleDeleteCache = async (hash: string) => {
    if (!confirm('Are you sure you want to delete this cache entry?')) return

    try {
      await api.deleteCache(hash)
      await loadData()
    } catch (error: any) {
      alert('Failed to delete cache: ' + error.message)
    }
  }

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-lg">Loading...</div>
      </div>
    )
  }

  const totalHits = analytics?.hit?.count || 0
  const totalMisses = analytics?.miss?.count || 0
  const hitRate = totalHits + totalMisses > 0 ? (totalHits / (totalHits + totalMisses)) * 100 : 0
  const timeSaved = analytics?.time_saved_ms || 0

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <header className="bg-white border-b">
        <div className="container mx-auto px-4 py-4 flex items-center justify-between">
          <div className="flex items-center gap-2">
            <div className="text-2xl font-bold">Kasoku</div>
            <span className="text-sm text-muted-foreground">Dashboard</span>
          </div>
          <Button variant="ghost" onClick={handleLogout}>
            Logout
          </Button>
        </div>
      </header>

      <div className="container mx-auto px-4 py-8">
        {/* Tabs */}
        <div className="flex gap-4 mb-8 border-b">
          <button
            className={`pb-2 px-4 ${activeTab === 'analytics' ? 'border-b-2 border-primary font-semibold' : 'text-muted-foreground'}`}
            onClick={() => setActiveTab('analytics')}
          >
            Analytics
          </button>
          <button
            className={`pb-2 px-4 ${activeTab === 'cache' ? 'border-b-2 border-primary font-semibold' : 'text-muted-foreground'}`}
            onClick={() => setActiveTab('cache')}
          >
            Cache
          </button>
          <button
            className={`pb-2 px-4 ${activeTab === 'tokens' ? 'border-b-2 border-primary font-semibold' : 'text-muted-foreground'}`}
            onClick={() => setActiveTab('tokens')}
          >
            API Tokens
          </button>
        </div>

        {/* Analytics Tab */}
        {activeTab === 'analytics' && (
          <div className="space-y-6">
            <h2 className="text-2xl font-bold">Analytics (Last 30 Days)</h2>

            <div className="grid md:grid-cols-4 gap-4">
              <div className="bg-white p-6 rounded-lg border">
                <div className="text-sm text-muted-foreground mb-1">Cache Hits</div>
                <div className="text-3xl font-bold">{totalHits}</div>
              </div>
              <div className="bg-white p-6 rounded-lg border">
                <div className="text-sm text-muted-foreground mb-1">Cache Misses</div>
                <div className="text-3xl font-bold">{totalMisses}</div>
              </div>
              <div className="bg-white p-6 rounded-lg border">
                <div className="text-sm text-muted-foreground mb-1">Hit Rate</div>
                <div className="text-3xl font-bold">{hitRate.toFixed(1)}%</div>
              </div>
              <div className="bg-white p-6 rounded-lg border">
                <div className="text-sm text-muted-foreground mb-1">Time Saved</div>
                <div className="text-3xl font-bold">{formatDuration(timeSaved)}</div>
              </div>
            </div>

            <div className="bg-white p-6 rounded-lg border">
              <h3 className="text-lg font-semibold mb-4">Usage Overview</h3>
              <div className="space-y-3">
                <div className="flex justify-between">
                  <span className="text-muted-foreground">Total Operations</span>
                  <span className="font-semibold">{totalHits + totalMisses}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-muted-foreground">Cache Size</span>
                  <span className="font-semibold">{formatBytes(analytics?.hit?.total_size || 0)}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-muted-foreground">Storage Used</span>
                  <span className="font-semibold">{formatBytes((analytics?.hit?.total_size || 0) + (analytics?.put?.total_size || 0))}</span>
                </div>
              </div>
            </div>
          </div>
        )}

        {/* Cache Tab */}
        {activeTab === 'cache' && (
          <div className="space-y-6">
            <div className="flex justify-between items-center">
              <h2 className="text-2xl font-bold">Cache Entries</h2>
              <div className="text-sm text-muted-foreground">{cacheEntries.length} entries</div>
            </div>

            {cacheEntries.length === 0 ? (
              <div className="bg-white p-12 rounded-lg border text-center text-muted-foreground">
                No cache entries yet. Run some commands with kasoku to populate the cache.
              </div>
            ) : (
              <div className="bg-white rounded-lg border overflow-hidden">
                <table className="w-full">
                  <thead className="bg-gray-50 border-b">
                    <tr>
                      <th className="text-left p-4 text-sm font-semibold">Command</th>
                      <th className="text-left p-4 text-sm font-semibold">Hash</th>
                      <th className="text-left p-4 text-sm font-semibold">Size</th>
                      <th className="text-left p-4 text-sm font-semibold">Created</th>
                      <th className="text-left p-4 text-sm font-semibold">Actions</th>
                    </tr>
                  </thead>
                  <tbody>
                    {cacheEntries.map((entry) => (
                      <tr key={entry.id} className="border-b last:border-0 hover:bg-gray-50">
                        <td className="p-4 text-sm font-mono">{entry.command || '-'}</td>
                        <td className="p-4 text-sm font-mono">{entry.hash.substring(0, 12)}...</td>
                        <td className="p-4 text-sm">{formatBytes(entry.size)}</td>
                        <td className="p-4 text-sm">{formatDate(entry.created_at)}</td>
                        <td className="p-4">
                          <Button
                            variant="destructive"
                            size="sm"
                            onClick={() => handleDeleteCache(entry.hash)}
                          >
                            Delete
                          </Button>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            )}
          </div>
        )}

        {/* Tokens Tab */}
        {activeTab === 'tokens' && (
          <div className="space-y-6">
            <h2 className="text-2xl font-bold">API Tokens</h2>

            <div className="bg-white p-6 rounded-lg border">
              <h3 className="text-lg font-semibold mb-4">Create New Token</h3>
              <form onSubmit={handleCreateToken} className="flex gap-4">
                <input
                  type="text"
                  placeholder="Token name (e.g., CI/CD, Local Dev)"
                  value={newTokenName}
                  onChange={(e) => setNewTokenName(e.target.value)}
                  className="flex-1 px-3 py-2 border border-input rounded-md focus:outline-none focus:ring-2 focus:ring-ring"
                />
                <Button type="submit">Create Token</Button>
              </form>

              {newToken && (
                <div className="mt-4 p-4 bg-blue-50 border border-blue-200 rounded">
                  <div className="text-sm font-semibold mb-2">Your new token (copy now, it won&apos;t be shown again):</div>
                  <div className="font-mono text-sm bg-white p-2 rounded border break-all">{newToken}</div>
                  <Button
                    size="sm"
                    variant="outline"
                    className="mt-2"
                    onClick={() => {
                      navigator.clipboard.writeText(newToken)
                      alert('Token copied to clipboard!')
                    }}
                  >
                    Copy to Clipboard
                  </Button>
                </div>
              )}
            </div>

            {tokens.length === 0 ? (
              <div className="bg-white p-12 rounded-lg border text-center text-muted-foreground">
                No API tokens yet. Create one to use with the CLI or CI/CD.
              </div>
            ) : (
              <div className="bg-white rounded-lg border overflow-hidden">
                <table className="w-full">
                  <thead className="bg-gray-50 border-b">
                    <tr>
                      <th className="text-left p-4 text-sm font-semibold">Name</th>
                      <th className="text-left p-4 text-sm font-semibold">Created</th>
                      <th className="text-left p-4 text-sm font-semibold">Last Used</th>
                      <th className="text-left p-4 text-sm font-semibold">Actions</th>
                    </tr>
                  </thead>
                  <tbody>
                    {tokens.map((token) => (
                      <tr key={token.id} className="border-b last:border-0 hover:bg-gray-50">
                        <td className="p-4 text-sm font-semibold">{token.name}</td>
                        <td className="p-4 text-sm">{formatDate(token.created_at)}</td>
                        <td className="p-4 text-sm">
                          {token.last_used ? formatDate(token.last_used) : 'Never'}
                        </td>
                        <td className="p-4">
                          <Button
                            variant="destructive"
                            size="sm"
                            onClick={() => handleRevokeToken(token.id)}
                          >
                            Revoke
                          </Button>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            )}
          </div>
        )}
      </div>
    </div>
  )
}
