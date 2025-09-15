'use client'

import React from 'react'

interface DashboardClientProps {
  children: React.ReactNode
}

/**
 * DashboardClient
 * - Previously wrapped the legacy DashboardLayout (which had its own Sidebar).
 * - Now acts as a thin client boundary to enable any client-only context/hooks
 *   while letting DashboardChrome control Sidebar/Topbar layout.
 */
export default function DashboardClient({ children }: DashboardClientProps) {
  return <>{children}</>
}