'use client'

import React, { useEffect, useMemo, useState } from 'react'
import { usePathname, useRouter } from 'next/navigation'
import {
  Gauge,
  Database,
  KeyRound,
  HardDrive,
  Bolt,
  FileCode,
  ListTree,
  Settings,
  CreditCard,
  Users,
  ChevronLeft,
  Menu,
} from 'lucide-react'

type Project = { id: string; name: string }

type NavItem = {
  label: string
  href: string
  icon: React.ComponentType<{ className?: string }>
  badge?: string
  soon?: boolean
  section?: 'primary' | 'secondary'
}

const STORAGE_KEY = 'dashboard.sidebar.collapsed:v1'

interface SidebarProps {
  projects: Project[]
  activeProjectId: string
  onProjectChange: (id: string) => void
  mobileOpen: boolean
  onCloseMobile: () => void
}

/**
 * Supabase-like Sidebar
 * - Desktop: sticky, collapsible rail with smooth transitions and persisted state
 * - Mobile (<lg): drawer from left with overlay
 * - Light/Dark aware via Tailwind tokens and .dark class
 */
export default function Sidebar({
  projects,
  activeProjectId,
  onProjectChange,
  mobileOpen,
  onCloseMobile,
}: SidebarProps) {
  const router = useRouter()
  const pathname = usePathname()
  const [collapsed, setCollapsed] = useState(false)

  // Load persisted collapsed state
  useEffect(() => {
    try {
      const raw = localStorage.getItem(STORAGE_KEY)
      if (raw != null) setCollapsed(raw === '1')
    } catch {}
  }, [])
  // Persist on change
  useEffect(() => {
    try {
      localStorage.setItem(STORAGE_KEY, collapsed ? '1' : '0')
    } catch {}
  }, [collapsed])

  const widthClass = collapsed ? 'w-[72px]' : 'w-[260px]'
  const railHideClass = collapsed ? 'opacity-0 pointer-events-none' : 'opacity-100'

  const nav: NavItem[] = useMemo(
    () => [
      { label: 'Dashboard', href: '/dashboard', icon: Gauge, section: 'primary' },
      { label: 'Table Editor', href: '/dashboard/table-editor', icon: ListTree, section: 'primary' },
      { label: 'SQL Editor', href: '/dashboard/sql-editor', icon: FileCode, section: 'primary', soon: true },
      { label: 'Database', href: '/dashboard/database', icon: Database, section: 'primary' },
      { label: 'Authentication', href: '/dashboard/auth', icon: KeyRound, section: 'primary' },
      { label: 'Storage', href: '/dashboard/storage', icon: HardDrive, section: 'primary' },
      { label: 'Edge Functions', href: '/dashboard/edge-functions', icon: Bolt, section: 'primary' },
      { label: 'Logs', href: '/dashboard/logs', icon: ListTree, section: 'primary' },

      { label: 'Settings', href: '/dashboard/settings', icon: Settings, section: 'secondary' },
      { label: 'Billing', href: '/dashboard/billing', icon: CreditCard, section: 'secondary' },
      { label: 'Team & Access', href: '/dashboard/team', icon: Users, section: 'secondary' },
    ],
    []
  )

  const renderItem = (item: NavItem) => {
    const isActive = pathname === item.href
    return (
      <button
        key={item.href}
        onClick={() => {
          router.push(item.href)
          onCloseMobile()
        }}
        className={[
          'group relative w-full flex items-center gap-3 rounded-md px-3 py-2 text-[13px] transition-colors',
          'text-neutral-700 hover:bg-neutral-100 dark:text-neutral-300 dark:hover:bg-neutral-800',
          isActive
            ? 'bg-emerald-50 text-emerald-800 dark:bg-emerald-900/20 dark:text-emerald-300'
            : '',
        ].join(' ')}
        aria-current={isActive ? 'page' : undefined}
      >
        {/* Left active indicator */}
        <span
          className={[
            'absolute left-0 top-1/2 -translate-y-1/2 h-5 w-[3px] rounded-r-full',
            isActive ? 'bg-emerald-600' : 'bg-transparent',
          ].join(' ')}
          aria-hidden
        />
        {/* Icon */}
        <item.icon
          className={[
            'h-4.5 w-4.5 flex-none',
            isActive
              ? 'text-emerald-600 dark:text-emerald-400'
              : 'text-neutral-500 group-hover:text-neutral-700 dark:text-neutral-400 dark:group-hover:text-neutral-200',
          ].join(' ')}
        />
        {/* Label */}
        <span className={collapsed ? 'sr-only' : 'truncate'}>{item.label}</span>
        {/* Badge / Soon */}
        {!collapsed && (item.badge || item.soon) ? (
          <span
            className={[
              'ml-auto inline-flex items-center rounded-full border px-1.5 py-[1px] text-[10px]',
              item.soon
                ? 'border-amber-200 bg-amber-50 text-amber-800 dark:border-amber-900/40 dark:bg-amber-900/20 dark:text-amber-300'
                : 'border-emerald-200 bg-emerald-50 text-emerald-700 dark:border-emerald-900/40 dark:bg-emerald-900/20 dark:text-emerald-300',
            ].join(' ')}
          >
            {item.soon ? 'Soon' : item.badge}
          </span>
        ) : null}
        {/* Tooltip when collapsed */}
        {collapsed && (
          <span className="pointer-events-none absolute left-[72px] z-20 ml-2 hidden rounded-md bg-neutral-900 px-2 py-1 text-xs text-white shadow-lg group-hover:block">
            {item.label}
          </span>
        )}
      </button>
    )
  }

  const ProjectSwitcher = () => (
    <div className="px-3 py-3">
      <label htmlFor="sidebar-project" className="sr-only">
        Project
      </label>
      <div className={['relative', collapsed ? 'px-0' : ''].join(' ')}>
        <select
          id="sidebar-project"
          value={activeProjectId}
          onChange={(e) => onProjectChange(e.target.value)}
          aria-label="Project"
          className={[
            'w-full appearance-none rounded-md border text-[13px] transition-colors',
            'bg-white text-neutral-900 border-neutral-200 hover:bg-neutral-50',
            'dark:bg-neutral-900 dark:text-neutral-100 dark:border-neutral-800 dark:hover:bg-neutral-800',
            collapsed ? 'sr-only' : 'px-2.5 py-2 pr-6',
          ].join(' ')}
        >
          {projects.map((p) => (
            <option key={p.id} value={p.id}>
              {p.name}
            </option>
          ))}
        </select>
        {!collapsed && (
          <span className="pointer-events-none absolute right-2 top-1/2 -translate-y-1/2 text-neutral-400">
            ▾
          </span>
        )}
      </div>
    </div>
  )

  // Desktop sidebar
  const Desktop = (
    <aside
      className={[
        'hidden lg:flex lg:flex-col shrink-0 border-r',
        'bg-white/90 backdrop-blur supports-[backdrop-filter]:bg-white/80',
        'dark:bg-neutral-900/80 dark:border-neutral-800 border-neutral-200',
        'transition-all duration-300 ease-in-out',
        widthClass,
      ].join(' ')}
      aria-label="Primary Navigation"
    >
      {/* Header / Toggle */}
      {/* <div className="flex items-center justify-between px-3 py-2 border-b border-neutral-200 dark:border-neutral-800">
        <div className={['flex items-center gap-2', railHideClass].join(' ')}>
          <div className="h-7 w-7 rounded-md bg-emerald-500 text-white grid place-items-center">DB</div>
          <span className="text-[13px] font-medium text-neutral-800 dark:text-neutral-200">Project</span>
        </div>
        <button
          className="inline-flex h-8 w-8 items-center justify-center rounded-md text-neutral-600 hover:bg-neutral-100 focus:outline-none focus:ring-2 focus:ring-emerald-400 dark:text-neutral-300 dark:hover:bg-neutral-800"
          onClick={() => setCollapsed((v) => !v)}
          aria-label={collapsed ? 'Expand sidebar' : 'Collapse sidebar'}
          aria-pressed={collapsed}
        >
          <ChevronLeft className={['h-4 w-4 transition-transform', collapsed ? 'rotate-180' : ''].join(' ')} />
        </button>
      </div> */}

      {/* Project Switcher */}
      {/* <div className={railHideClass}>
        <ProjectSwitcher />
      </div> */}

      {/* Primary nav */}
      <nav className="px-2 py-2" aria-label="Main">
        <div className="space-y-1">
          {nav.filter((n) => n.section !== 'secondary').map(renderItem)}
        </div>
      </nav>

      {/* Divider */}
      <div className="my-2 border-t border-neutral-200 dark:border-neutral-800" />

      {/* Secondary nav */}
      <nav className="px-2 pb-3" aria-label="Secondary">
        <div className="space-y-1">{nav.filter((n) => n.section === 'secondary').map(renderItem)}</div>
      </nav>

      {/* Footer */}
      <div className="mt-auto px-3 py-2 text-[11px] text-neutral-500 dark:text-neutral-400">
        Pillow Console • v0.1
      </div>
    </aside>
  )

  // Mobile drawer
  const Mobile = (
    <>
      {/* Toggle button is in Topbar; here we render the drawer */}
      <div
        className={[
          'fixed inset-0 z-40 lg:hidden transition-opacity',
          mobileOpen ? 'opacity-100' : 'pointer-events-none opacity-0',
        ].join(' ')}
        aria-hidden={!mobileOpen}
        onClick={onCloseMobile}
      >
        <div className="absolute inset-0 bg-black/40" />
      </div>
      <div
        className={[
          'fixed inset-y-0 left-0 z-50 w-[260px] lg:hidden',
          'bg-white dark:bg-neutral-900 border-r border-neutral-200 dark:border-neutral-800 shadow-xl',
          'transition-transform duration-300 ease-in-out',
          mobileOpen ? 'translate-x-0' : '-translate-x-full',
        ].join(' ')}
        role="dialog"
        aria-modal="true"
        aria-label="Sidebar"
      >
        {/* Mobile header */}
        <div className="flex items-center justify-between px-3 py-2 border-b border-neutral-200 dark:border-neutral-800">
          <div className="flex items-center gap-2">
            <div className="h-7 w-7 rounded-md bg-emerald-500 text-white grid place-items-center">DB</div>
            <span className="text-[13px] font-medium text-neutral-800 dark:text-neutral-200">Navigation</span>
          </div>
          <button
            className="inline-flex h-8 w-8 items-center justify-center rounded-md text-neutral-600 hover:bg-neutral-100 focus:outline-none focus:ring-2 focus:ring-emerald-400 dark:text-neutral-300 dark:hover:bg-neutral-800"
            onClick={onCloseMobile}
            aria-label="Close sidebar"
          >
            <Menu className="h-4 w-4" />
          </button>
        </div>
        <div className="overflow-y-auto h-[calc(100vh-44px)] pb-3">
          <ProjectSwitcher />
          <nav className="px-2 py-2" aria-label="Main">
            <div className="space-y-1">
              {nav.filter((n) => n.section !== 'secondary').map(renderItem)}
            </div>
          </nav>
          <div className="my-2 border-t border-neutral-200 dark:border-neutral-800" />
          <nav className="px-2 pb-3" aria-label="Secondary">
            <div className="space-y-1">{nav.filter((n) => n.section === 'secondary').map(renderItem)}</div>
          </nav>
        </div>
      </div>
    </>
  )

  return (
    <>
      {Desktop}
      {Mobile}
    </>
  )
}