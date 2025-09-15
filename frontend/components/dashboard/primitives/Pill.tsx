'use client'

import React from 'react'
import clsx from 'clsx'

type PillVariant = 'muted' | 'amber' | 'emerald' | 'gray' | 'outline'
type PillProps = {
  children: React.ReactNode
  className?: string
  variant?: PillVariant
  size?: 'xs' | 'sm'
}

const base =
  'inline-flex items-center gap-1.5 rounded-full border px-1.5 py-[1px] align-middle'

const variants: Record<PillVariant, string> = {
  muted:
    'border-neutral-200 bg-neutral-50 text-neutral-700 dark:border-neutral-800 dark:bg-neutral-900/30 dark:text-neutral-300',
  amber:
    'border-amber-200 bg-amber-100 text-amber-800 dark:border-amber-900/40 dark:bg-amber-900/20 dark:text-amber-300',
  emerald:
    'border-emerald-200 bg-emerald-50 text-emerald-700 dark:border-emerald-900/40 dark:bg-emerald-900/20 dark:text-emerald-300',
  gray:
    'border-neutral-200 bg-white text-neutral-700 shadow-sm dark:bg-neutral-900 dark:border-neutral-800 dark:text-neutral-200',
  outline:
    'border-neutral-300 text-neutral-700 dark:border-neutral-700 dark:text-neutral-200',
}

const sizes = {
  xs: 'text-[10px]',
  sm: 'text-[12px] px-2 py-[3px]',
}

export default function Pill({ children, className, variant = 'muted', size = 'xs' }: PillProps) {
  return <span className={clsx(base, variants[variant], sizes[size], className)}>{children}</span>
}