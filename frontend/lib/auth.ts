'use client'

import { api } from './api'
import Cookies from 'js-cookie'

const TOKEN_KEY = 'pillow_token'
const USER_KEY = 'pillow_user'

export interface User {
  id: string
  username: string
  email: string
  is_active?: boolean
  created_at?: string
}

export function setAuth(token: string, user: User) {
  Cookies.set(TOKEN_KEY, token, {
    expires: 1, // 1 day
    secure: process.env.NODE_ENV === 'production',
    sameSite: 'strict'
  })
  Cookies.set(USER_KEY, JSON.stringify(user), {
    expires: 1,
    secure: process.env.NODE_ENV === 'production',
    sameSite: 'strict'
  })
  api.defaults.headers.common['Authorization'] = `Bearer ${token}`
}

export function clearAuth() {
  Cookies.remove(TOKEN_KEY)
  Cookies.remove(USER_KEY)
  delete api.defaults.headers.common['Authorization']
}

export function getToken(): string | undefined {
  return Cookies.get(TOKEN_KEY)
}

export function getUser(): User | null {
  const raw = Cookies.get(USER_KEY)
  if (!raw) return null
  try {
    return JSON.parse(raw) as User
  } catch {
    return null
  }
}

export function isAuthenticated(): boolean {
  return !!getToken()
}

// Initialize axios header if token exists (client runtime)
if (typeof window !== 'undefined') {
  const token = getToken()
  if (token) {
    api.defaults.headers.common['Authorization'] = `Bearer ${token}`
  }
}