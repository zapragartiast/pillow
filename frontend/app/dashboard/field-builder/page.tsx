'use client'

import { useState, useEffect } from 'react'
import { useRouter } from 'next/navigation'
import SimpleFieldBuilder from '@/components/dashboard/SimpleFieldBuilder'

export default function FieldBuilderPage() {
  const [userId, setUserId] = useState<string | null>(null)
  const router = useRouter()

  useEffect(() => {
    // Get user ID from localStorage or API
    const storedUser = localStorage.getItem('pillow_user')
    if (storedUser) {
      try {
        const user = JSON.parse(storedUser)
        setUserId(user.id)
      } catch (error) {
        console.error('Failed to parse user data:', error)
        router.push('/login')
      }
    } else {
      router.push('/login')
    }
  }, [router])

  if (!userId) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
      </div>
    )
  }

  return (
    <div className="max-w-4xl mx-auto p-6">
      <SimpleFieldBuilder userId={userId} />
    </div>
  )
}