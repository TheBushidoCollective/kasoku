const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'

export interface User {
  id: string
  email: string
  name: string
  plan: string
  created_at: string
}

export interface CacheEntry {
  id: string
  hash: string
  command: string
  size: number
  created_at: string
  accessed_at: string
}

export interface Token {
  id: string
  name: string
  token?: string
  created_at: string
  expires_at?: string
  last_used?: string
}

export interface AnalyticsSummary {
  hit?: { count: number; total_size: number }
  miss?: { count: number; total_size: number }
  put?: { count: number; total_size: number }
  time_saved_ms: number
}

export class ApiError extends Error {
  constructor(public status: number, message: string) {
    super(message)
    this.name = 'ApiError'
  }
}

async function fetchApi(endpoint: string, options: RequestInit = {}) {
  const token = typeof window !== 'undefined' ? localStorage.getItem('kasoku_token') : null

  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    ...(typeof options.headers === 'object' && options.headers !== null 
      ? (options.headers as Record<string, string>)
      : {}),
  }

  if (token) {
    headers['Authorization'] = `Bearer ${token}`
  }

  const response = await fetch(`${API_URL}${endpoint}`, {
    ...options,
    headers,
  })

  if (!response.ok) {
    const error = await response.json().catch(() => ({ error: response.statusText }))
    throw new ApiError(response.status, error.error || 'An error occurred')
  }

  return response.json()
}

export function isAuthenticated(): boolean {
  if (typeof window === 'undefined') return false
  return !!localStorage.getItem('kasoku_token')
}

export const api = {
  // Auth
  async register(email: string, name: string, password: string) {
    return fetchApi('/auth/register', {
      method: 'POST',
      body: JSON.stringify({ email, name, password }),
    })
  },

  async login(email: string, password: string): Promise<{ token: string; user: User }> {
    const data = await fetchApi('/auth/login', {
      method: 'POST',
      body: JSON.stringify({ email, password }),
    })
    localStorage.setItem('kasoku_token', data.token)
    return data
  },

  async logout() {
    localStorage.removeItem('kasoku_token')
  },

  // Cache
  async listCache(limit = 50, offset = 0): Promise<{ entries: CacheEntry[] }> {
    return fetchApi(`/cache?limit=${limit}&offset=${offset}`)
  },

  async deleteCache(hash: string) {
    return fetchApi(`/cache/${hash}`, { method: 'DELETE' })
  },

  // Analytics
  async getAnalytics(since?: Date): Promise<AnalyticsSummary> {
    const params = since ? `?since=${since.toISOString()}` : ''
    return fetchApi(`/analytics${params}`)
  },

  // Tokens
  async createToken(name: string): Promise<Token> {
    return fetchApi('/auth/tokens', {
      method: 'POST',
      body: JSON.stringify({ name }),
    })
  },

  async listTokens(): Promise<{ tokens: Token[] }> {
    return fetchApi('/auth/tokens')
  },

  async revokeToken(id: string) {
    return fetchApi(`/auth/tokens/${id}`, { method: 'DELETE' })
  },

  // User
  async getMe(): Promise<User> {
    return fetchApi('/auth/me')
  },
}
