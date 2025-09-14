'use server'

import { cookies } from 'next/headers'
import { redirect } from 'next/navigation'

export async function registerAction(formData: FormData) {
  const username = formData.get('username') as string
  const email = formData.get('email') as string
  const password = formData.get('password') as string

  if (!username || !email || !password) {
    redirect(`/register?error=${encodeURIComponent('All fields are required')}`)
  }

  if (password.length < 6) {
    redirect(`/register?error=${encodeURIComponent('Password must be at least 6 characters')}`)
  }

  try {
    // Call Next.js API route instead of backend directly
    const apiUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:3000'
    const registerResponse = await fetch(`${apiUrl}/api/auth/register`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ username, email, password }),
    })

    if (!registerResponse.ok) {
      const responseText = await registerResponse.text()
      let errorData
      try {
        errorData = JSON.parse(responseText)
      } catch {
        errorData = { error: responseText }
      }
      const errorMessage = errorData.error || errorData.message || 'Registration failed'
      redirect(`/register?error=${encodeURIComponent(errorMessage)}`)
    }

    // Auto-login after successful registration
    const loginResponse = await fetch(`${apiUrl}/api/auth/login`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ identifier: username, password }),
    })

    if (!loginResponse.ok) {
      redirect(`/register?error=${encodeURIComponent('Registration successful, but login failed. Please try logging in manually.')}`)
    }

    const loginData = await loginResponse.json()

    if (loginData.token) {
      // Set authentication cookie
      const cookieStore = await cookies()
      cookieStore.set('pillow_token', loginData.token, {
        httpOnly: true,
        secure: process.env.NODE_ENV === 'production',
        sameSite: 'strict',
        maxAge: 60 * 60 * 24, // 24 hours
      })

      // Set user data cookie
      cookieStore.set('pillow_user', JSON.stringify(loginData.user), {
        httpOnly: true,
        secure: process.env.NODE_ENV === 'production',
        sameSite: 'strict',
        maxAge: 60 * 60 * 24,
      })

      redirect('/dashboard')
    }

    redirect(`/register?error=${encodeURIComponent('Registration successful, but login failed. Please try logging in manually.')}`)
  } catch (error: any) {
    // Re-throw redirect errors to let Next.js handle them
    if (error?.digest?.startsWith('NEXT_REDIRECT')) {
      throw error
    }
    console.error('Registration error:', error)
    redirect(`/register?error=${encodeURIComponent('Network error. Please try again.')}`)
  }
}