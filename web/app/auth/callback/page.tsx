'use client'

import { useEffect } from 'react'
import { useRouter, useSearchParams } from 'next/navigation'

export default function OAuthCallback() {
  const router = useRouter()
  const searchParams = useSearchParams()

  useEffect(() => {
    const token = searchParams.get('token')
    const userId = searchParams.get('user_id')

    if (token && userId) {
      // Store token in localStorage
      localStorage.setItem('kasoku_token', token)
      localStorage.setItem('kasoku_user_id', userId)

      // Redirect to dashboard
      router.push('/dashboard')
    } else {
      // If no token, redirect to login with error
      router.push('/login?error=oauth_failed')
    }
  }, [searchParams, router])

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-b from-gray-50 to-white">
      <div className="text-center">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary mx-auto mb-4"></div>
        <p className="text-muted-foreground">Completing authentication...</p>
      </div>
    </div>
  )
}
