'use client'

import * as React from "react"
import { X } from "lucide-react"

export interface DialogProps {
  open?: boolean
  onOpenChange?: (open: boolean) => void
  children?: React.ReactNode
}

export interface DialogContentProps {
  children?: React.ReactNode
  className?: string
}

export interface DialogHeaderProps {
  children?: React.ReactNode
}

export interface DialogTitleProps {
  children?: React.ReactNode
}

export interface DialogTriggerProps {
  children?: React.ReactNode
  asChild?: boolean
}

export function Dialog({ open, onOpenChange, children }: DialogProps) {
  const [internalOpen, setInternalOpen] = React.useState(false)

  const isOpen = open !== undefined ? open : internalOpen
  const setIsOpen = onOpenChange || setInternalOpen

  React.useEffect(() => {
    if (isOpen) {
      document.body.style.overflow = 'hidden'
    } else {
      document.body.style.overflow = 'unset'
    }

    return () => {
      document.body.style.overflow = 'unset'
    }
  }, [isOpen])

  if (!isOpen) return null

  return (
    <>
      <div
        className="fixed inset-0 z-50 bg-black/50"
        onClick={() => setIsOpen(false)}
      />
      <div className="fixed left-1/2 top-1/2 z-50 -translate-x-1/2 -translate-y-1/2">
        {children}
      </div>
    </>
  )
}

export function DialogContent({ children, className = "" }: DialogContentProps) {
  return (
    <div className={`relative grid w-full max-w-lg gap-4 border border-gray-200 bg-white p-6 shadow-lg duration-200 sm:rounded-lg dark:border-gray-700 dark:bg-gray-800 ${className}`}>
      {children}
    </div>
  )
}

export function DialogHeader({ children }: DialogHeaderProps) {
  return (
    <div className="flex flex-col space-y-1.5 text-center sm:text-left">
      {children}
    </div>
  )
}

export function DialogTitle({ children }: DialogTitleProps) {
  return (
    <h3 className="text-lg font-semibold leading-none tracking-tight">
      {children}
    </h3>
  )
}

export function DialogTrigger({ children, asChild }: DialogTriggerProps) {
  if (asChild && React.isValidElement(children)) {
    return React.cloneElement(children as React.ReactElement<any>, {
      onClick: (e: React.MouseEvent) => {
        e.preventDefault()
        // This would normally trigger the dialog, but we'll handle it in the parent
      }
    })
  }

  return <>{children}</>
}