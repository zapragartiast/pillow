import { cookies } from 'next/headers'
import { redirect } from 'next/navigation'
import { logoutAction } from './actions'

async function getUserProfile() {
  const cookieStore = await cookies()
  const token = cookieStore.get('pillow_token')?.value

  if (!token) {
    redirect('/login')
  }

  try {
    const apiUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:3000'

    const response = await fetch(`${apiUrl}/api/auth/profile`, {
      method: 'GET',
      headers: {
        'Cookie': `pillow_token=${token}`, // Pass token via cookie since it's server-side
      },
    })

    if (!response.ok) {
      if (response.status === 401) {
        // Token expired or invalid, clear cookies and redirect
        cookieStore.delete('pillow_token')
        cookieStore.delete('pillow_user')
        redirect('/login')
      }
      throw new Error('Failed to fetch user profile')
    }

    return await response.json()
  } catch (error) {
    console.error('Error fetching user profile:', error)
    // Clear cookies and redirect on any error
    cookieStore.delete('pillow_token')
    cookieStore.delete('pillow_user')
    redirect('/login')
  }
}

export default async function DashboardPage() {
  const user = await getUserProfile()

  return (
    <div className="mt-8">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-semibold">Dashboard</h1>
        <form action={logoutAction}>
          <button
            type="submit"
            className="px-3 py-2 bg-red-600 text-white rounded-md hover:bg-red-500 focus:outline-none focus:ring-2 focus:ring-red-500 focus:ring-offset-2"
          >
            Sign out
          </button>
        </form>
      </div>

      <div className="mt-6 bg-white border rounded-lg shadow p-6">
        <h2 className="text-lg font-medium mb-2">Profile</h2>
        {user ? (
          <div className="space-y-1 text-sm text-slate-700">
            <div><strong>ID:</strong> {user.id}</div>
            <div><strong>Username:</strong> {user.username}</div>
            <div><strong>Email:</strong> {user.email}</div>
            <div><strong>Active:</strong> {String(user.is_active ?? user.isActive ?? true)}</div>
          </div>
        ) : (
          <div>No profile available</div>
        )}
      </div>
    </div>
  )
}