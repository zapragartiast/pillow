import { NextRequest, NextResponse } from 'next/server'
import { cookies } from 'next/headers'

export async function GET(
  request: NextRequest,
  { params }: { params: Promise<{ userId: string }> }
) {
  console.log("🔥 FRONTEND API ROUTE: GET custom-field-values started")
  
  try {
    const { userId } = await params
    console.log("🔥 UserID:", userId)
    
    const cookieStore = await cookies()
    const token = cookieStore.get('pillow_token')?.value
    
    console.log("🔥 Token from cookies:", token ? `Found (${token.substring(0, 20)}...)` : "❌ NOT FOUND")

    if (!token) {
      console.log("🔥 No token - returning 401 from FRONTEND")
      return NextResponse.json(
        { error: 'Unauthorized' },
        { status: 401 }
      )
    }

    const backendUrl = process.env.BACKEND_URL || 'http://localhost:8080'
    console.log("🔥 Calling backend:", `${backendUrl}/api/users/${userId}/custom-field-values`)

    const response = await fetch(`${backendUrl}/api/users/${userId}/custom-field-values`, {
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
        { error: data.message || data.error || 'Failed to fetch custom field values' },
        { status: response.status }
      )
    }

    return NextResponse.json(data)
  } catch (error) {
    console.error('Custom field values API error:', error)
    return NextResponse.json(
      { error: 'Internal server error' },
      { status: 500 }
    )
  }
}

export async function PUT(
  request: NextRequest,
  { params }: { params: Promise<{ userId: string }> }
) {
  try {
    const { userId } = await params
    const cookieStore = await cookies()
    const token = cookieStore.get('pillow_token')?.value

    if (!token) {
      return NextResponse.json(
        { error: 'Unauthorized' },
        { status: 401 }
      )
    }

    const body = await request.text()
    const backendUrl = process.env.BACKEND_URL || 'http://localhost:8080'

    const response = await fetch(`${backendUrl}/api/users/${userId}/custom-field-values`, {
      method: 'PUT',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json',
      },
      body: body,
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
        { error: data.message || data.error || 'Failed to update custom field values' },
        { status: response.status }
      )
    }

    return NextResponse.json(data)
  } catch (error) {
    console.error('Custom field values PUT API error:', error)
    return NextResponse.json(
      { error: 'Internal server error' },
      { status: 500 }
    )
  }
}
