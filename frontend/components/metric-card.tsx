import { Card, CardContent } from './ui/card'

interface MetricCardProps {
  icon: React.ReactNode
  value: string | number
  title: string
  color: 'blue' | 'green' | 'orange' | 'purple'
  className?: string
}

const colorClasses = {
  blue: 'text-blue-400',
  green: 'text-green-400',
  orange: 'text-orange-400',
  purple: 'text-purple-400',
}

export default function MetricCard({ icon, value, title, color, className }: MetricCardProps) {
  return (
    <Card className={`bg-gray-800 border border-gray-700 hover:border-gray-600 rounded-xl shadow-sm hover:shadow-md transition-all duration-200 ${className}`}>
      <CardContent className="p-6">
        <div className="flex items-center justify-between mb-4">
          <div className="p-3 bg-gray-700 rounded-lg">
            <div className="text-gray-400">
              {icon}
            </div>
          </div>
          <div className="text-right">
            <div className={`text-3xl font-bold ${colorClasses[color]} mb-1`}>
              {value}
            </div>
            <div className="text-sm text-gray-400 font-medium">
              {title}
            </div>
          </div>
        </div>
        <div className="w-full bg-gray-700 rounded-full h-2">
          <div className={`h-2 rounded-full bg-gradient-to-r ${color === 'blue' ? 'from-blue-500 to-blue-400' :
            color === 'green' ? 'from-green-500 to-green-400' :
            color === 'orange' ? 'from-orange-500 to-orange-400' :
            'from-purple-500 to-purple-400'} transition-all duration-300`}
            style={{ width: '75%' }}>
          </div>
        </div>
      </CardContent>
    </Card>
  )
}