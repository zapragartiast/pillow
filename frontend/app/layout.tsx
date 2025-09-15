import './globals.css'
import React from 'react'
import { ThemeProvider } from '@/lib/theme-context'

export const metadata = {
  title: 'Pillow Dashboard',
  description: 'Admin dashboard',
}

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en" suppressHydrationWarning>
      <body className="min-h-screen bg-white text-gray-900 dark:lg:bg-zinc-900 dark:text-white transition-colors" suppressHydrationWarning={true}>
        <ThemeProvider>
          <div className="mx-auto">{children}</div>
        </ThemeProvider>
      </body>
    </html>
  )
}