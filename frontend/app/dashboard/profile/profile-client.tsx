'use client'

import UserProfileCustomFields from '@/components/dashboard/UserProfileCustomFields'

interface ProfileClientProps {
  user: {
    id: string
    username: string
    email: string
  }
}

export default function ProfileClient({ user }: ProfileClientProps) {
  return (
    <div className="max-w-4xl mx-auto space-y-8">
      {/* Profile Header */}
      <div className="bg-white dark:bg-gray-800 rounded-lg shadow p-6">
        <div className="flex items-center space-x-4">
          <div className="w-16 h-16 bg-blue-600 rounded-full flex items-center justify-center text-white text-2xl font-bold">
            {user.username?.charAt(0).toUpperCase() || 'U'}
          </div>
          <div>
            <h1 className="text-2xl font-bold text-gray-900 dark:text-white">
              {user.username}
            </h1>
            <p className="text-gray-600 dark:text-gray-400">
              {user.email}
            </p>
          </div>
        </div>
      </div>

      {/* Custom Fields Section */}
      <div className="bg-white dark:bg-gray-800 rounded-lg shadow">
        <div className="p-6">
          <UserProfileCustomFields
            userId={user.id}
            onFieldsChange={(fields) => {
              console.log('Fields changed:', fields)
            }}
          />
        </div>
      </div>
    </div>
  )
}
