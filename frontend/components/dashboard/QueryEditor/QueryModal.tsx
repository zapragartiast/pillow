'use client'

import React, { useEffect, useRef, useState } from 'react'

type Column = { key: string; label: string }
type ResultSet = {
  columns: Column[]
  rows: Record<string, any>[]
  rowCount: number
  executionTimeMs: number
  warnings?: string[]
}

export type QueryModalProps = {
  open: boolean
  onClose: () => void
}

export default function QueryModal({ open, onClose }: QueryModalProps) {
  const [sql, setSql] = useState<string>('select * from users limit 10;')
  const [running, setRunning] = useState(false)
  const [result, setResult] = useState<ResultSet | null>(null)
  const [error, setError] = useState<string | null>(null)
  const dialogRef = useRef<HTMLDivElement | null>(null)

  useEffect(() => {
    function onKey(e: KeyboardEvent) {
      if (!open) return
      if (e.key === 'Escape') {
        e.preventDefault()
        onClose()
      }
      if ((e.metaKey || e.ctrlKey) && e.key.toLowerCase() === 'enter') {
        e.preventDefault()
        runQuery().catch(() => {})
      }
    }
    window.addEventListener('keydown', onKey)
    return () => window.removeEventListener('keydown', onKey)
  }, [open, sql])

  useEffect(() => {
    if (open) {
      // simple focus trap: focus the textarea when opened
      const ta = dialogRef.current?.querySelector('textarea') as HTMLTextAreaElement | null
      ta?.focus()
    }
  }, [open])

  if (!open) return null

  async function runQuery() {
    setRunning(true)
    setError(null)
    setResult(null)
    try {
      const res = await fetch('/api/demo/query', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ sql }),
      })
      if (!res.ok) {
        const data = await res.json().catch(() => ({}))
        throw new Error(data?.error || 'Failed to run query')
      }
      const json = await res.json()
      setResult(json)
    } catch (e: any) {
      setError(e?.message || 'Failed to run query')
    } finally {
      setRunning(false)
    }
  }

  return (
    <div
      className="fixed inset-0 z-50 flex items-center justify-center"
      role="dialog"
      aria-modal="true"
      aria-labelledby="query-modal-title"
    >
      {/* Backdrop */}
      <div className="absolute inset-0 bg-black/50" onClick={onClose} />

      {/* Modal */}
      <div
        ref={dialogRef}
        className="relative z-10 w-[95vw] max-w-5xl max-h-[85vh] bg-white dark:lg:bg-zinc-900 border border-gray-200 dark:border-gray-800 rounded-lg shadow-xl overflow-hidden"
      >
        {/* Header */}
        <div className="flex items-center justify-between px-4 py-3 border-b border-gray-200 dark:border-gray-800">
          <div>
            <h2 id="query-modal-title" className="text-base font-semibold text-gray-900 dark:text-white">
              SQL Editor
            </h2>
            <p className="text-xs text-gray-500 dark:text-gray-400">Cmd/Ctrl + Enter to run</p>
          </div>
          <div className="flex items-center gap-2">
            <button
              onClick={runQuery}
              disabled={running}
              className="inline-flex items-center gap-2 px-3 py-1.5 rounded-md text-sm bg-emerald-500 hover:bg-emerald-600 disabled:bg-emerald-600/60 text-white focus:outline-none focus:ring-2 focus:ring-emerald-400"
            >
              {running ? (
                <svg className="h-4 w-4 animate-spin" viewBox="0 0 24 24" fill="none">
                  <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" />
                  <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8H4z" />
                </svg>
              ) : (
                <svg className="h-4 w-4" viewBox="0 0 24 24" fill="currentColor">
                  <path d="M8 5v14l11-7-11-7z" />
                </svg>
              )}
              Run
            </button>
            <button
              onClick={onClose}
              className="inline-flex items-center px-2.5 py-1.5 rounded-md text-sm text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-800 focus:outline-none focus:ring-2 focus:ring-primary-500"
              aria-label="Close"
            >
              Close
            </button>
          </div>
        </div>

        {/* Body grid */}
        <div className="grid grid-cols-1 lg:grid-cols-[1fr_420px] gap-0 lg:gap-3 p-3">
          {/* Editor */}
          <div className="flex flex-col min-h-[40vh] lg:min-h-[60vh]">
            <label htmlFor="sql-input" className="sr-only">
              SQL
            </label>
            <textarea
              id="sql-input"
              value={sql}
              onChange={(e) => setSql(e.target.value)}
              className="flex-1 w-full rounded-md border border-gray-200 dark:border-gray-800 bg-white dark:lg:bg-zinc-900 text-sm text-gray-900 dark:text-gray-100 p-3 font-mono leading-6 focus:outline-none focus:ring-2 focus:ring-primary-500"
              placeholder="Write SQL here..."
              spellCheck={false}
            />
            {error && (
              <div className="mt-2 text-sm text-red-600 dark:text-red-400" role="alert">
                {error}
              </div>
            )}
            {result?.warnings?.length ? (
              <div className="mt-2 text-xs text-yellow-700 dark:text-yellow-300">
                {result.warnings.map((w, i) => (
                  <div key={i}>• {w}</div>
                ))}
              </div>
            ) : null}
          </div>

          {/* Result pane */}
          <div className="mt-3 lg:mt-0 bg-gray-50 dark:bg-gray-950/40 border border-gray-200 dark:border-gray-800 rounded-md overflow-hidden">
            <div className="flex items-center justify-between px-3 py-2 border-b border-gray-200 dark:border-gray-800">
              <div className="text-sm text-gray-700 dark:text-gray-300">
                {running ? 'Running…' : result ? `${result.rowCount} rows` : 'Results'}
              </div>
              {result && (
                <div className="text-xs text-gray-500 dark:text-gray-400">
                  {result.executionTimeMs} ms
                </div>
              )}
            </div>

            <div className="overflow-auto max-h-[50vh]">
              {!result ? (
                <div className="p-4 text-sm text-gray-500 dark:text-gray-400">
                  Run a query to see results.
                </div>
              ) : result.rows.length === 0 ? (
                <div className="p-4 text-sm text-gray-500 dark:text-gray-400">No rows returned.</div>
              ) : (
                <table className="min-w-full text-sm">
                  <thead>
                    <tr className="bg-gray-100 dark:lg:bg-zinc-900/60 sticky top-0">
                      {result.columns.map((c) => (
                        <th
                          key={c.key}
                          scope="col"
                          className="text-left px-3 py-2 font-semibold text-gray-700 dark:text-gray-200 border-b border-gray-200 dark:border-gray-800"
                        >
                          {c.label}
                        </th>
                      ))}
                    </tr>
                  </thead>
                  <tbody>
                    {result.rows.map((r, i) => (
                      <tr key={i} className="border-b border-gray-100 dark:border-gray-900">
                        {result.columns.map((c) => (
                          <td key={c.key} className="px-3 py-2 text-gray-900 dark:text-gray-100">
                            {formatCell(r[c.key])}
                          </td>
                        ))}
                      </tr>
                    ))}
                  </tbody>
                </table>
              )}
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}

function formatCell(v: any) {
  if (v == null) return <span className="text-gray-400">NULL</span>
  if (typeof v === 'string') {
    // simple timestamp highlight
    if (/^\d{4}-\d{2}-\d{2}T/.test(v)) {
      return <time className="font-mono">{v}</time>
    }
    return v
  }
  if (typeof v === 'number') return v
  if (typeof v === 'boolean') return String(v)
  return JSON.stringify(v)
}