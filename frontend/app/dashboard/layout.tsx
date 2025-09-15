import React from 'react'
import DashboardClient from '@/components/dashboard-client'
import DashboardChrome from '@/components/dashboard/DashboardChrome'

export const metadata = {
  title: 'Dashboard',
}

type Project = { id: string; name: string }

export default function DashboardRouteLayout({ children }: { children: React.ReactNode }) {
  // Server component wrapper for the dashboard route.
  // Only pass serializable props to client components.
  const projects: Project[] = [
    { id: 'pillow', name: 'Pillow' },
    { id: 'staging', name: 'Staging' },
    { id: 'prod', name: 'Production' },
  ]
  const activeProjectId: string = 'pillow'
  const breadcrumbs = [{ label: 'Projects', href: '/dashboard' }, { label: 'Overview' }]

  return (
    <DashboardClient>
      <DashboardChrome
        breadcrumbs={breadcrumbs}
        projects={projects}
        activeProjectId={activeProjectId}
      >
        {children}
      </DashboardChrome>
    </DashboardClient>
  )
}