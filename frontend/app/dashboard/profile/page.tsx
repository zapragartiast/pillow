'use client'

import { useEffect, useState } from 'react'
import { useRouter } from 'next/navigation'
import ProfileClient from './profile-client'

export default function ProfilePage() {
  const router = useRouter()
  const [user, setUser] = useState<any | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    let mounted = true
    const load = async () => {
      try {
        const res = await fetch('/api/auth/profile', {
          method: 'GET',
          credentials: 'include',
          cache: 'no-store',
        })
        if (!mounted) return
        if (res.status === 401) {
          router.replace('/login')
          return
        }
        if (!res.ok) {
          throw new Error('Failed to fetch profile')
        }
        const data = await res.json()
        setUser(data)
      } catch (e) {
        console.error('Profile load error:', e)
        router.replace('/login')
      } finally {
        if (mounted) setLoading(false)
      }
    }
    load()
    return () => {
      mounted = false
    }
  }, [router])

  if (loading || !user) {
    return (
      <div className="flex items-center justify-center min-h-[60vh]">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600" />
      </div>
    )
  }

  return <ProfileClient user={user} />
}
