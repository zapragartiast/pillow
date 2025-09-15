'use client'

import React from 'react'
import { useRouter, usePathname, useSearchParams } from 'next/navigation'
import Topbar from './Topbar'
import Sidebar from './Sidebar'
import QueryModal from './QueryEditor/QueryModal'

type Breadcrumb = { label: string; href?: string }
type Project = { id: string; name: string }

interface DashboardChromeProps {
  breadcrumbs: Breadcrumb[]
  projects: Project[]
  activeProjectId: string
  children: React.ReactNode
}

export default function DashboardChrome({
  breadcrumbs,
  projects,
  activeProjectId,
  children,
}: DashboardChromeProps) {
  const router = useRouter()
  const pathname = usePathname()
  const searchParams = useSearchParams()
  const [isQueryOpen, setIsQueryOpen] = React.useState(false)
  const [mobileSidebarOpen, setMobileSidebarOpen] = React.useState(false)

  const handleProjectChange = (id: string) => {
    const sp = new URLSearchParams(searchParams?.toString() ?? '')
    sp.set('project', id)
    router.push(`${pathname}?${sp.toString()}`)
  }

  return (
    <>
      <Topbar
        breadcrumbs={breadcrumbs}
        projects={projects}
        activeProjectId={activeProjectId}
        onProjectChange={handleProjectChange}
        onNewQuery={() => setIsQueryOpen(true)}
        onInvite={() => {
          // placeholder for invite modal; can be wired to Auth flow later
          alert('Open Invite modal (to be implemented)')
        }}
        onQuickAction={(action) => {
          if (action === 'new-query') setIsQueryOpen(true)
          if (action === 'new-table') router.push('/dashboard/table-editor')
        }}
        onOpenSidebar={() => setMobileSidebarOpen(true)}
      />
      <div className="relative">
        <div className="flex min-h-[calc(100vh-52px)]">
          {/* Sidebar: desktop inline + mobile drawer */}
          <Sidebar
            projects={projects}
            activeProjectId={activeProjectId}
            onProjectChange={handleProjectChange}
            mobileOpen={mobileSidebarOpen}
            onCloseMobile={() => setMobileSidebarOpen(false)}
          />
          {/* Content */}
          <section aria-labelledby="page-title" className="flex-1 max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4">
            {children}
          </section>
        </div>
      </div>
      <QueryModal open={isQueryOpen} onClose={() => setIsQueryOpen(false)} />
    </>
  )
}