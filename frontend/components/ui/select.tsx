'use client'

import * as React from "react"
import { Check, ChevronDown } from "lucide-react"

export interface SelectProps {
  children?: React.ReactNode
  value?: string
  onValueChange?: (value: string) => void
}

export interface SelectTriggerProps {
  children?: React.ReactNode
  className?: string
}

export interface SelectContentProps {
  children?: React.ReactNode
}

export interface SelectItemProps {
  children?: React.ReactNode
  value: string
}

export interface SelectValueProps {
  placeholder?: string
}

const SelectContext = React.createContext<{
  value?: string
  onValueChange?: (value: string) => void
  open: boolean
  setOpen: (open: boolean) => void
}>({
  open: false,
  setOpen: () => {},
})

export function Select({ children, value, onValueChange }: SelectProps) {
  const [open, setOpen] = React.useState(false)

  return (
    <SelectContext.Provider value={{ value, onValueChange, open, setOpen }}>
      <div className="relative">
        {children}
      </div>
    </SelectContext.Provider>
  )
}

export function SelectTrigger({ children, className = "" }: SelectTriggerProps) {
  const { open, setOpen } = React.useContext(SelectContext)

  return (
    <button
      type="button"
      onClick={() => setOpen(!open)}
      className={`flex h-10 w-full items-center justify-between rounded-md border border-gray-300 bg-white px-3 py-2 text-sm ring-offset-white placeholder:text-gray-500 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50 dark:border-gray-600 dark:bg-gray-800 dark:ring-offset-gray-900 dark:placeholder:text-gray-400 ${className}`}
    >
      {children}
      <ChevronDown className="h-4 w-4 opacity-50" />
    </button>
  )
}

export function SelectValue({ placeholder }: SelectValueProps) {
  const { value } = React.useContext(SelectContext)

  return <span className="text-sm">{value || placeholder}</span>
}

export function SelectContent({ children }: SelectContentProps) {
  const { open, setOpen } = React.useContext(SelectContext)

  if (!open) return null

  return (
    <>
      <div
        className="fixed inset-0 z-50"
        onClick={() => setOpen(false)}
      />
      <div className="absolute top-full z-50 min-w-[8rem] overflow-hidden rounded-md border border-gray-200 bg-white shadow-md dark:border-gray-700 dark:bg-gray-800">
        <div className="p-1">
          {children}
        </div>
      </div>
    </>
  )
}

export function SelectItem({ children, value }: SelectItemProps) {
  const { value: selectedValue, onValueChange, setOpen } = React.useContext(SelectContext)

  const handleClick = () => {
    onValueChange?.(value)
    setOpen(false)
  }

  return (
    <div
      onClick={handleClick}
      className="relative flex w-full cursor-pointer select-none items-center rounded-sm py-1.5 pl-8 pr-2 text-sm outline-none hover:bg-gray-100 dark:hover:bg-gray-700"
    >
      {selectedValue === value && (
        <Check className="absolute left-2 h-4 w-4" />
      )}
      <span className="ml-2">{children}</span>
    </div>
  )
}