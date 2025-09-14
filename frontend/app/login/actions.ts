'use server'

import { cookies } from 'next/headers'
import { redirect } from 'next/navigation'

export async function loginAction(formData: FormData) {
  const identifier = formData.get('identifier') as string
  const password = formData.get('password') as string

  if (!identifier || !password) {
    redirect(`/login?error=${encodeURIComponent('Username/email and password are required')}`)
  }

  try {
    // Call Next.js API route instead of backend directly
    const response = await fetch('http://localhost:3000/api/auth/login', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ identifier, password }),
    })

    if (!response.ok) {
      const errorData = await response.json().catch(() => ({}))
      const errorMessage = errorData.error || 'Login failed'
      redirect(`/login?error=${encodeURIComponent(errorMessage)}`)
    }

    const data = await response.json()

    if (data.token) {
      // Set authentication cookie
      const cookieStore = await cookies()
      cookieStore.set('pillow_token', data.token, {
        httpOnly: true,
        secure: process.env.NODE_ENV === 'production',
        sameSite: 'strict',
        maxAge: 60 * 60 * 24, // 24 hours
      })

      // Set user data cookie
      cookieStore.set('pillow_user', JSON.stringify(data.user), {
        httpOnly: true,
        secure: process.env.NODE_ENV === 'production',
        sameSite: 'strict',
        maxAge: 60 * 60 * 24,
      })

      redirect('/dashboard')
    }

    redirect(`/login?error=${encodeURIComponent('Invalid response from server')}`)
  } catch (error: any) {
    // Re-throw redirect errors to let Next.js handle them
    if (error?.digest?.startsWith('NEXT_REDIRECT')) {
      throw error
    }
    console.error('Login error:', error)
    redirect(`/login?error=${encodeURIComponent('Network error. Please try again.')}`)
  }
}