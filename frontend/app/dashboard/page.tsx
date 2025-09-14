import { cookies } from 'next/headers'
import { redirect } from 'next/navigation'
import DashboardClient from '@/components/dashboard-client'
import DashboardHeader from '@/components/dashboard-header'
import MetricCard from '@/components/metric-card'

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
    <DashboardClient>
      <div className="max-w-7xl mx-auto">
        <DashboardHeader
          title="Dashboard"
          subtitle={`Welcome back, ${user.username}`}
          username={user.username}
          email={user.email}
        />

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
          <MetricCard
            icon={
              <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
              </svg>
            }
            value="2,456"
            title="Active Users"
            color="blue"
          />
          <MetricCard
            icon={
              <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M4 7v10c0 2.21 3.582 4 8 4s8-1.79 8-4V7M4 7c0 2.21 3.582 4 8 4s8-1.79 8-4M4 7c0-2.21 3.582-4 8-4s8 1.79 8 4" />
              </svg>
            }
            value="1,234"
            title="Database Queries"
            color="green"
          />
          <MetricCard
            icon={
              <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M5 19a2 2 0 01-2-2V7a2 2 0 012-2h4l2 2h4a2 2 0 012 2v1M5 19h14a2 2 0 002-2v-5a2 2 0 00-2-2H9a2 2 0 00-2 2v5a2 2 0 01-2 2z" />
              </svg>
            }
            value="89.2 GB"
            title="Storage Used"
            color="orange"
          />
          <MetricCard
            icon={
              <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6" />
              </svg>
            }
            value="99.9%"
            title="Uptime"
            color="purple"
          />
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
          <div className="bg-gray-800 border border-gray-700 rounded-xl p-6 shadow-sm">
            <div className="flex items-center justify-between mb-6">
              <div>
                <h2 className="text-lg font-semibold text-white mb-1">Database Performance</h2>
                <p className="text-sm text-gray-400">Real-time metrics and insights</p>
              </div>
              <div className="p-3 bg-gray-700 rounded-lg">
                <svg className="w-6 h-6 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
                </svg>
              </div>
            </div>
            <div className="h-64 bg-gradient-to-br from-gray-700 to-gray-800 rounded-lg border border-gray-600 flex items-center justify-center">
              <div className="text-center">
                <div className="w-16 h-16 bg-gradient-to-br from-blue-500 to-blue-600 rounded-full flex items-center justify-center mx-auto mb-4">
                  <svg className="w-8 h-8 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
                  </svg>
                </div>
                <p className="text-gray-400 text-sm">Performance charts would appear here</p>
              </div>
            </div>
          </div>

          <div className="bg-gray-800 border border-gray-700 rounded-xl p-6 shadow-sm">
            <div className="flex items-center justify-between mb-6">
              <div>
                <h2 className="text-lg font-semibold text-white mb-1">Recent Activity</h2>
                <p className="text-sm text-gray-400">Latest system events and updates</p>
              </div>
              <div className="p-3 bg-gray-700 rounded-lg">
                <svg className="w-6 h-6 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
              </div>
            </div>
            <div className="space-y-4">
              <div className="flex items-start gap-4 p-4 bg-gray-700 border border-gray-600 rounded-lg hover:bg-gray-650 transition-colors">
                <div className="w-3 h-3 bg-green-400 rounded-full mt-1.5 flex-shrink-0"></div>
                <div className="flex-1 min-w-0">
                  <p className="text-sm text-white font-medium">Database backup completed</p>
                  <p className="text-xs text-gray-400 mt-1">Automated backup finished successfully • 2 minutes ago</p>
                </div>
              </div>
              <div className="flex items-start gap-4 p-4 bg-gray-700 border border-gray-600 rounded-lg hover:bg-gray-650 transition-colors">
                <div className="w-3 h-3 bg-blue-400 rounded-full mt-1.5 flex-shrink-0"></div>
                <div className="flex-1 min-w-0">
                  <p className="text-sm text-white font-medium">New user registered</p>
                  <p className="text-xs text-gray-400 mt-1">john.doe@example.com joined the platform • 15 minutes ago</p>
                </div>
              </div>
              <div className="flex items-start gap-4 p-4 bg-gray-700 border border-gray-600 rounded-lg hover:bg-gray-650 transition-colors">
                <div className="w-3 h-3 bg-orange-400 rounded-full mt-1.5 flex-shrink-0"></div>
                <div className="flex-1 min-w-0">
                  <p className="text-sm text-white font-medium">API rate limit warning</p>
                  <p className="text-xs text-gray-400 mt-1">Approaching 80% of monthly limit • 1 hour ago</p>
                </div>
              </div>
              <div className="flex items-start gap-4 p-4 bg-gray-700 border border-gray-600 rounded-lg hover:bg-gray-650 transition-colors">
                <div className="w-3 h-3 bg-purple-400 rounded-full mt-1.5 flex-shrink-0"></div>
                <div className="flex-1 min-w-0">
                  <p className="text-sm text-white font-medium">Edge function deployed</p>
                  <p className="text-xs text-gray-400 mt-1">user-authentication function updated • 2 hours ago</p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </DashboardClient>
  )
}