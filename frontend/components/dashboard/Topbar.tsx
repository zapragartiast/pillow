'use client'

import React from 'react'
import { useTheme } from '@/lib/theme-context'
import {
  ChevronRight,
  Database,
  Plus,
  Settings,
  UserPlus,
  Search,
  Moon,
  Sun,
  Play,
  Zap,
  User,
  GitBranch,
  ShieldCheck,
  Link2,
  HelpCircle,
  Bell,
  ChevronDown,
} from 'lucide-react'

type Breadcrumb = { label: string; href?: string }
type Project = { id: string; name: string }

interface TopbarProps {
  breadcrumbs: Breadcrumb[]
  projects: Project[]
  activeProjectId: string
  onProjectChange: (id: string) => void
  onNewQuery: () => void
  onInvite: () => void
  onQuickAction: (action: 'new-table' | 'new-query' | 'invite') => void
  onOpenSidebar?: () => void
}

export default function Topbar({
  breadcrumbs,
  projects,
  activeProjectId,
  onProjectChange,
  onNewQuery,
  onInvite,
  onQuickAction,
  onOpenSidebar,
}: TopbarProps) {
  const { theme, toggleTheme } = useTheme()
  const activeProject = projects.find(p => p.id === activeProjectId) ?? projects[0]

  return (
    <header
      className="sticky top-0 z-30 bg-white/90 dark:bg-neutral-900/80 backdrop-blur supports-[backdrop-filter]:bg-white/80 border-b dark:border-neutral-800 border-neutral-200 dark:border-gray-800"
      aria-label="Top navigation"
    >
      <div className="mx-auto px-3 h-12 sm:h-[52px] flex items-center gap-2 text-[13px]">
        {/* Mobile menu button */}
        <button
          type="button"
          className="inline-flex lg:hidden items-center justify-center h-9 w-9 rounded-md text-gray-600 hover:bg-gray-100 focus:outline-none focus:ring-2 focus:ring-emerald-400 dark:text-gray-300 dark:hover:bg-gray-800"
          aria-label="Open sidebar"
          onClick={() => onOpenSidebar?.()}
        >
          {/* Simple hamburger */}
          <span className="block w-4 h-0.5 bg-current mb-1.5 rounded"></span>
          <span className="block w-4 h-0.5 bg-current mb-1.5 rounded"></span>
          <span className="block w-4 h-0.5 bg-current rounded"></span>
        </button>
        {/* Brand bolt */}
        <div className="hidden sm:grid place-items-center h-7 w-7 rounded-md bg-emerald-500 text-white shadow-sm" aria-hidden="true">
          <Zap className="h-4 w-4" />
        </div>

        {/* Org / crumbs row (Supabase-like chips) */}
        <nav aria-label="Breadcrumb" className="flex items-center gap-2 text-[13px] min-w-0">
          {/* Org chip */}
          <span className="inline-flex items-center gap-1.5 px-2 py-[5px] rounded-full border border-gray-200 text-gray-700 bg-white shadow-sm dark:lg:bg-zinc-900 dark:border-gray-800 dark:text-gray-200">
            <User className="h-3.5 w-3.5" />
            <span className="truncate max-w-[9rem]">hosting.com Ticket</span>
            <span className="ml-1 rounded-full bg-gray-100 text-gray-700 px-1.5 text-[11px] border border-gray-200 hidden sm:inline">Super Premium</span>
          </span>

          <span className="text-gray-300 dark:text-gray-600 select-none">/</span>

          {/* Project selector chip */}
          <label htmlFor="project-selector" className="sr-only">Select project</label>
          <div className="relative">
            <select
              id="project-selector"
              className="appearance-none pr-6 inline-flex items-center gap-1.5 px-2 py-[5px] rounded-full border border-gray-200 text-gray-800 bg-white shadow-sm text-[13px] dark:lg:bg-zinc-900 dark:text-gray-100 dark:border-gray-800 focus:outline-none focus:ring-2 focus:ring-emerald-400"
              value={activeProjectId}
              onChange={(e) => onProjectChange(e.target.value)}
              aria-label="Project selector"
            >
              {projects.map((p) => (
                <option key={p.id} value={p.id}>{p.name}</option>
              ))}
            </select>
            <ChevronDown className="pointer-events-none absolute right-1.5 top-1/2 -translate-y-1/2 h-3.5 w-3.5 text-gray-400" />
          </div>

          <span className="text-gray-300 dark:text-gray-600 select-none">/</span>

          {/* Branch chip */}
          {/* <span className="inline-flex items-center gap-1.5 px-2 py-[5px] rounded-full border border-gray-200 text-gray-700 bg-white shadow-sm dark:lg:bg-zinc-900 dark:border-gray-800 dark:text-gray-200">
            <GitBranch className="h-3.5 w-3.5" />
            <span className="truncate max-w-[6rem]">main</span>
          </span> */}

          {/* Environment pill */}
          <span className="inline-flex items-center gap-1.5 px-2 py-[5px] rounded-full bg-amber-100 text-amber-800 border border-amber-200 text-[12px]">
            <ShieldCheck className="h-3.5 w-3.5" />
            Production
          </span>

          {/* Connect button (ghost chip) */}
          {/* <button
            type="button"
            className="hidden md:inline-flex items-center gap-1.5 px-2.5 py-[5px] rounded-full border border-gray-200 text-gray-700 bg-white shadow-sm hover:bg-gray-50 dark:lg:bg-zinc-900 dark:border-gray-800 dark:text-gray-200 dark:hover:bg-gray-800"
          >
            <Link2 className="h-3.5 w-3.5" />
            <span className="text-sm">Connect</span>
          </button> */}
        </nav>

        {/* Spacer */}
        <div className="flex-1" />

        {/* Right controls */}
        <div className="flex items-center gap-1.5">
          {/* Search (placeholder) */}
          <button
            type="button"
            className="hidden sm:inline-flex items-center gap-2 px-2.5 py-1.5 rounded-md text-[13px] text-gray-700 hover:bg-gray-100 focus:outline-none focus:ring-2 focus:ring-emerald-400 dark:text-gray-300 dark:hover:bg-gray-800"
            aria-label="Search"
          >
            <Search className="h-4 w-4" />
            <span className="hidden lg:inline">Search</span>
          </button>

          {/* Quick actions */}
          {/* <button
            type="button"
            onClick={() => onQuickAction('new-table')}
            className="hidden md:inline-flex items-center gap-1.5 px-2.5 py-1.5 rounded-md text-[13px] text-gray-700 hover:bg-gray-100 border border-gray-200 bg-white shadow-sm focus:outline-none focus:ring-2 focus:ring-emerald-400 dark:text-gray-300 dark:lg:bg-zinc-900 dark:border-gray-800 dark:hover:bg-gray-800"
          >
            <Plus className="h-4 w-4" /> New table
          </button> */}
          {/* <button
            type="button"
            onClick={() => { onNewQuery(); onQuickAction('new-query') }}
            className="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-md text-[13px] bg-emerald-500 hover:bg-emerald-600 text-white shadow focus:outline-none focus:ring-2 focus:ring-emerald-400"
          >
            <Play className="h-4 w-4" /> New query
          </button> */}
          <button
            type="button"
            onClick={() => { onInvite(); onQuickAction('invite') }}
            className="hidden md:inline-flex items-center gap-1.5 px-2.5 py-1.5 rounded-md text-[13px] text-gray-700 hover:bg-gray-100 border border-gray-200 bg-white shadow-sm focus:outline-none focus:ring-2 focus:ring-emerald-400 dark:text-gray-900 dark:bg-gray-100 dark:border-gray-800 dark:hover:bg-gray-100"
          >
            <UserPlus className="h-4 w-4" /> Invite
          </button>
          {/* <button
            type="button"
            className="inline-flex h-8 w-8 sm:h-9 sm:w-9 items-center justify-center rounded-md text-gray-600 hover:bg-gray-100 focus:outline-none focus:ring-2 focus:ring-emerald-400 dark:text-gray-300 dark:hover:bg-gray-800"
            aria-label="Feedback"
            title="Feedback"
          >
            <span className="text-xs font-medium">F</span>
          </button> */}
          <button
            type="button"
            className="inline-flex h-8 w-8 sm:h-9 sm:w-9 items-center justify-center rounded-md text-gray-600 hover:bg-gray-100 focus:outline-none focus:ring-2 focus:ring-emerald-400 dark:text-gray-300 dark:hover:bg-gray-800"
            aria-label="Help"
            title="Help"
          >
            <HelpCircle className="h-5 w-5" />
          </button>
          <button
            type="button"
            className="inline-flex h-8 w-8 sm:h-9 sm:w-9 items-center justify-center rounded-md text-gray-600 hover:bg-gray-100 focus:outline-none focus:ring-2 focus:ring-emerald-400 dark:text-gray-300 dark:hover:bg-gray-800"
            aria-label="Notifications"
            title="Notifications"
          >
            <Bell className="h-5 w-5" />
          </button>

          {/* <button
            type="button"
            onClick={toggleTheme}
            className="inline-flex h-8 w-8 sm:h-9 sm:w-9 items-center justify-center rounded-md text-gray-600 hover:bg-gray-100 focus:outline-none focus:ring-2 focus:ring-emerald-400 dark:text-gray-300 dark:hover:bg-gray-800"
            aria-label={`Switch to ${theme === 'dark' ? 'light' : 'dark'} theme`}
            title="Toggle theme"
          >
            {theme === 'dark' ? <Sun className="h-5 w-5" /> : <Moon className="h-5 w-5" />}
          </button> */}

          {/* Avatar placeholder */}
          <button
            type="button"
            className="inline-flex items-center justify-center h-8 w-8 sm:h-9 sm:w-9 rounded-full bg-zinc-200 text-gray-700 border border-gray-300 hover:ring-2 hover:ring-emerald-400 focus:outline-none dark:bg-zinc-800 dark:text-gray-200 dark:border-gray-700"
            aria-label="Account menu"
            title="Account"
          >
            <Settings className="h-4 w-4" />
          </button>
        </div>
      </div>
    </header>
  )
}