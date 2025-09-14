'use server'

import { cookies } from 'next/headers'
import { redirect } from 'next/navigation'

export async function logoutAction() {
  const cookieStore = await cookies()

  // Clear authentication cookies
  cookieStore.delete('pillow_token')
  cookieStore.delete('pillow_user')

  redirect('/login')
}