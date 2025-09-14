import UserAvatar from './user-avatar'

interface DashboardHeaderProps {
  title: string
  subtitle?: string
  username?: string
  email?: string
  showAvatar?: boolean
  className?: string
}

export default function DashboardHeader({
  title,
  subtitle,
  username,
  email,
  showAvatar = true,
  className = ''
}: DashboardHeaderProps) {
  return (
    <div className={`space-y-6 ${className}`}>
      {/* Header Content */}
      <header className="hidden md:flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-semibold text-white mb-1">{title}</h1>
          {subtitle && (
            <p className="text-sm text-gray-400">{subtitle}</p>
          )}
        </div>
        {showAvatar && username && (
          <div className="flex items-center gap-4">
            <UserAvatar
              name={username}
              email={email || ''}
              size="md"
            />
          </div>
        )}
      </header>

      {/* Visual Separator */}
      <div className="border-b border-gray-700"></div>
    </div>
  )
}