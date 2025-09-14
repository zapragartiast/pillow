import { NextRequest, NextResponse } from 'next/server'
import { cookies } from 'next/headers'

export async function POST(request: NextRequest) {
  try {
    // Clear authentication cookies
    const cookieStore = await cookies()

    // Delete the authentication tokens
    cookieStore.delete('pillow_token')
    cookieStore.delete('pillow_user')

    // Return success response
    return NextResponse.json(
      { message: 'Logged out successfully' },
      { status: 200 }
    )
  } catch (error) {
    console.error('Logout error:', error)
    return NextResponse.json(
      { error: 'Internal server error' },
      { status: 500 }
    )
  }
}