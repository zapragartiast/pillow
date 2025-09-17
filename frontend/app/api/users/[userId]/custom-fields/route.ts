import { NextRequest, NextResponse } from 'next/server'

export async function GET(
  request: NextRequest,
  { params }: { params: { userId: string } }
) {
  try {
    // Forward the request to the backend API
    const backendUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:3000'
    const backendResponse = await fetch(`${backendUrl}/api/users/${params.userId}/custom-fields`, {
      headers: {
        'Authorization': request.headers.get('Authorization') || '',
        'Content-Type': 'application/json',
      },
    })

    if (!backendResponse.ok) {
      return NextResponse.json(
        { error: 'Failed to fetch custom fields' },
        { status: backendResponse.status }
      )
    }

    const data = await backendResponse.json()
    return NextResponse.json(data)
  } catch (error) {
    console.error('API Error:', error)
    return NextResponse.json(
      { error: 'Internal server error' },
      { status: 500 }
    )
  }
}

export async function POST(
  request: NextRequest,
  { params }: { params: { userId: string } }
) {
  try {
    const body = await request.json()

    // Forward the request to the backend API
    const backendUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:3000'
    const backendResponse = await fetch(`${backendUrl}/api/users/${params.userId}/custom-fields`, {
      method: 'POST',
      headers: {
        'Authorization': request.headers.get('Authorization') || '',
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(body),
    })

    if (!backendResponse.ok) {
      const errorData = await backendResponse.json()
      return NextResponse.json(
        { error: errorData.message || 'Failed to create custom field' },
        { status: backendResponse.status }
      )
    }

    const data = await backendResponse.json()
    return NextResponse.json(data)
  } catch (error) {
    console.error('API Error:', error)
    return NextResponse.json(
      { error: 'Internal server error' },
      { status: 500 }
    )
  }
}