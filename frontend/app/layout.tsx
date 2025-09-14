import './globals.css'
import React from 'react'

export const metadata = {
  title: 'Pillow Dashboard',
  description: 'Admin dashboard',
}

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en">
      <body className="min-h-screen bg-slate-50 text-slate-900" suppressHydrationWarning={true}>
        <div className="mx-auto">{children}</div>
      </body>
    </html>
  )
}