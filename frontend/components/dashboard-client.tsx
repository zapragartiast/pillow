'use client'

import { useTheme } from '@/lib/theme-context'
import DashboardLayout from './dashboard-layout'

interface DashboardClientProps {
  children: React.ReactNode
}

export default function DashboardClient({ children }: DashboardClientProps) {
  // This ensures the theme context is available on the client side
  const { theme } = useTheme()

  return (
    <DashboardLayout>
      {children}
    </DashboardLayout>
  )
}