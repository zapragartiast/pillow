interface UserAvatarProps {
  name: string
  email: string
  size?: 'sm' | 'md' | 'lg'
  className?: string
}

export default function UserAvatar({ name, email, size = 'md', className }: UserAvatarProps) {
  const initials = name
    .split(' ')
    .map(word => word.charAt(0).toUpperCase())
    .join('')
    .slice(0, 2)

  const sizeClasses = {
    sm: 'w-8 h-8 text-sm',
    md: 'w-10 h-10 text-base',
    lg: 'w-12 h-12 text-lg'
  }

  return (
    <div className={`flex items-center gap-3 ${className}`}>
      <div className={`${sizeClasses[size]} bg-gray-700 text-white border border-gray-600 rounded-full flex items-center justify-center font-medium`}>
        {initials}
      </div>
      <div className="hidden md:block">
        <div className="text-sm font-medium text-white">{name}</div>
        <div className="text-xs text-gray-400">{email}</div>
      </div>
    </div>
  )
}