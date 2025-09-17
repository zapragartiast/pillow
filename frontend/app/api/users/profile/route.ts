import { NextRequest, NextResponse } from 'next/server'
import { cookies } from 'next/headers'

export async function GET(request: NextRequest) {
  try {
    const cookieStore = await cookies()
    const token = cookieStore.get('pillow_token')?.value

    if (!token) {
      return NextResponse.json(
        { error: 'Unauthorized' },
        { status: 401 }
      )
    }

    const backendUrl = process.env.BACKEND_URL || 'http://localhost:8080'

    const response = await fetch(`${backendUrl}/api/users/profile`, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json',
      },
    })

    const responseText = await response.text()
    let data
    try {
      data = JSON.parse(responseText)
    } catch {
      data = { message: responseText }
    }

    if (!response.ok) {
      if (response.status === 401) {
        // Clear invalid token
        cookieStore.delete('pillow_token')
        cookieStore.delete('pillow_user')
      }
      return NextResponse.json(
        { error: data.message || data.error || 'Failed to fetch profile' },
        { status: response.status }
      )
    }

    return NextResponse.json(data)
  } catch (error) {
    console.error('Profile API error:', error)
    return NextResponse.json(
      { error: 'Internal server error' },
      { status: 500 }
    )
  }
}