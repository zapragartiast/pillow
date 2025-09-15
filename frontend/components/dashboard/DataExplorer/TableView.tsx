'use client'

import React, { useCallback, useEffect, useMemo, useRef, useState } from 'react'

type SortDir = 'asc' | 'desc' | null

export type Column = {
  key: string
  header: string
  width?: number
  editable?: boolean
  accessor?: (row: any) => React.ReactNode
}

export type FetchParams = {
  page: number
  pageSize: number
  sortKey?: string
  sortDir?: Exclude<SortDir, null>
  filters?: Record<string, string>
}

export type FetchResult = {
  rows: any[]
  total: number
}

export type TableViewProps = {
  columns: Column[]
  fetchPage: (params: FetchParams) => Promise<FetchResult>
  initialPageSize?: number
  className?: string
  onEditCell?: (row: any, key: string, value: string | number | null) => Promise<void> | void
  ariaLabel?: string
}

type EditState = {
  rowIndex: number
  colKey: string
  initialValue: string
}

export default function TableView({
  columns,
  fetchPage,
  initialPageSize = 20,
  className = '',
  onEditCell,
  ariaLabel = 'Data table',
}: TableViewProps) {
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(initialPageSize)
  const [sortKey, setSortKey] = useState<string | undefined>(undefined)
  const [sortDir, setSortDir] = useState<SortDir>(null)
  const [filters, setFilters] = useState<Record<string, string>>({})
  const [rows, setRows] = useState<any[]>([])
  const [total, setTotal] = useState(0)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [editing, setEditing] = useState<EditState | null>(null)

  // Column widths and resizing
  const [colWidths, setColWidths] = useState<Record<string, number>>(() => {
    const init: Record<string, number> = {}
    columns.forEach((c) => (init[c.key] = c.width ?? 180))
    return init
  })
  const resizerInfo = useRef<{ key: string; startX: number; startW: number } | null>(null)

  const onMouseMove = (e: MouseEvent) => {
    if (!resizerInfo.current) return
    const delta = e.clientX - resizerInfo.current.startX
    const newW = Math.max(80, resizerInfo.current.startW + delta)
    setColWidths((prev) => ({ ...prev, [resizerInfo.current!.key]: newW }))
  }
  const onMouseUp = () => {
    window.removeEventListener('mousemove', onMouseMove)
    window.removeEventListener('mouseup', onMouseUp)
    resizerInfo.current = null
  }
  const startResize = (key: string, e: React.MouseEvent) => {
    resizerInfo.current = { key, startX: e.clientX, startW: colWidths[key] ?? 180 }
    window.addEventListener('mousemove', onMouseMove)
    window.addEventListener('mouseup', onMouseUp)
    e.preventDefault()
  }
  useEffect(() => {
    return () => {
      window.removeEventListener('mousemove', onMouseMove)
      window.removeEventListener('mouseup', onMouseUp)
    }
  }, [])

  const totalPages = useMemo(() => Math.max(1, Math.ceil(total / pageSize)), [total, pageSize])

  const load = useCallback(async () => {
    setLoading(true)
    setError(null)
    try {
      const result = await fetchPage({
        page,
        pageSize,
        sortKey,
        sortDir: sortDir ?? undefined,
        filters,
      })
      setRows(result.rows)
      setTotal(result.total)
    } catch (err: any) {
      setError(err?.message ?? 'Failed to load data')
    } finally {
      setLoading(false)
    }
  }, [fetchPage, page, pageSize, sortKey, sortDir, filters])

  useEffect(() => {
    load()
  }, [load])

  const toggleSort = (key: string) => {
    if (sortKey !== key) {
      setSortKey(key)
      setSortDir('asc')
    } else {
      setSortDir((prev) => (prev === 'asc' ? 'desc' : prev === 'desc' ? null : 'asc'))
      if (sortDir === null) setSortKey(undefined)
    }
    setPage(1)
  }

  const handleCellDoubleClick = (rowIndex: number, colKey: string, currentValue: any) => {
    const col = columns.find((c) => c.key === colKey)
    if (!col?.editable) return
    setEditing({
      rowIndex,
      colKey,
      initialValue: String(currentValue ?? ''),
    })
  }

  const saveEdit = async (rowIndex: number, colKey: string, newValue: string) => {
    if (!onEditCell) {
      setEditing(null)
      return
    }
    try {
      await onEditCell(rows[rowIndex], colKey, newValue)
      // optimistic update
      setRows((prev) => {
        const clone = [...prev]
        clone[rowIndex] = { ...clone[rowIndex], [colKey]: newValue }
        return clone
      })
    } finally {
      setEditing(null)
    }
  }

  const onKeyDownEditor = (e: React.KeyboardEvent<HTMLInputElement>, rowIndex: number, colKey: string) => {
    if (e.key === 'Enter') {
      saveEdit(rowIndex, colKey, (e.target as HTMLInputElement).value)
    } else if (e.key === 'Escape') {
      setEditing(null)
    }
  }

  const headerCell = (c: Column) => {
    const isSorted = sortKey === c.key && sortDir
    return (
      <th
        key={c.key}
        scope="col"
        style={{ width: colWidths[c.key] }}
        className="group sticky top-0 z-10 bg-white dark:lg:bg-zinc-900 border-b border-gray-200 dark:border-gray-800 text-left text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wide select-none"
        aria-sort={
          isSorted ? (sortDir === 'asc' ? 'ascending' : sortDir === 'desc' ? 'descending' : 'none') : 'none'
        }
      >
        <div className="relative flex items-center">
          <button
            type="button"
            onClick={() => toggleSort(c.key)}
            className="text-left w-full px-3 py-2 hover:text-gray-800 dark:hover:text-gray-200 focus:outline-none focus:ring-2 focus:ring-primary-500 rounded-sm"
          >
            <span>{c.header}</span>
            {isSorted && (
              <span className="ml-1 inline-block align-middle text-[10px]" aria-hidden>
                {sortDir === 'asc' ? '▲' : '▼'}
              </span>
            )}
          </button>
          <div
            className="absolute right-0 top-0 h-full w-2 cursor-col-resize opacity-0 group-hover:opacity-100 transition-opacity"
            onMouseDown={(e) => startResize(c.key, e)}
            aria-hidden
            title="Resize column"
          />
        </div>
      </th>
    )
  }

  return (
    <div
      className={`rounded-lg border border-gray-200 dark:border-gray-800 bg-white dark:lg:bg-zinc-900 ${className}`}
      role="region"
      aria-label={ariaLabel}
    >
      {/* Toolbar */}
      <div className="flex items-center justify-between px-3 py-2 border-b border-gray-200 dark:border-gray-800">
        <div className="flex items-center gap-2">
          <label htmlFor="global-filter" className="sr-only">
            Filter
          </label>
          <input
            id="global-filter"
            type="text"
            placeholder="Filter..."
            className="w-48 md:w-72 text-sm px-3 py-1.5 rounded-md bg-gray-100 dark:bg-gray-800 text-gray-900 dark:text-gray-100 border border-gray-200 dark:border-gray-700 focus:outline-none focus:ring-2 focus:ring-primary-500"
            value={filters.q ?? ''}
            onChange={(e) => {
              setFilters((f) => ({ ...f, q: e.target.value }))
              setPage(1)
            }}
          />
        </div>
        <div className="flex items-center gap-2 text-sm">
          <label htmlFor="page-size" className="text-gray-600 dark:text-gray-300">
            Rows / page
          </label>
          <select
            id="page-size"
            className="bg-gray-100 dark:bg-gray-800 text-gray-900 dark:text-gray-100 px-2 py-1 rounded-md border border-gray-200 dark:border-gray-700 focus:outline-none focus:ring-2 focus:ring-primary-500"
            value={pageSize}
            onChange={(e) => {
              setPageSize(Number(e.target.value))
              setPage(1)
            }}
          >
            {[10, 20, 50, 100].map((n) => (
              <option key={n} value={n}>
                {n}
              </option>
            ))}
          </select>
        </div>
      </div>

      {/* Table container */}
      <div className="overflow-auto" role="grid" aria-rowcount={rows.length} aria-colcount={columns.length}>
        <table className="min-w-full divide-y divide-gray-200 dark:divide-gray-800">
          <thead>
            <tr>{columns.map((c) => headerCell(c))}</tr>
          </thead>
          <tbody className="divide-y divide-gray-100 dark:divide-gray-800">
            {loading ? (
              // Skeleton rows
              Array.from({ length: 5 }).map((_, i) => (
                <tr key={`sk-${i}`} className="animate-pulse">
                  {columns.map((c) => (
                    <td key={c.key} style={{ width: colWidths[c.key] }} className="px-3 py-2">
                      <div className="h-4 bg-gray-100 dark:bg-gray-800 rounded" />
                    </td>
                  ))}
                </tr>
              ))
            ) : error ? (
              <tr>
                <td colSpan={columns.length} className="px-3 py-6 text-sm text-red-600 dark:text-red-400">
                  {error}
                </td>
              </tr>
            ) : rows.length === 0 ? (
              <tr>
                <td colSpan={columns.length} className="px-3 py-6 text-sm text-gray-500 dark:text-gray-400">
                  No data
                </td>
              </tr>
            ) : (
              rows.map((row, rIdx) => (
                <tr key={rIdx} className="hover:bg-gray-50 dark:hover:bg-gray-800">
                  {columns.map((c) => {
                    const display = c.accessor ? c.accessor(row) : row[c.key]
                    const isEditing = editing && editing.rowIndex === rIdx && editing.colKey === c.key
                    return (
                      <td
                        key={c.key}
                        style={{ width: colWidths[c.key] }}
                        className="px-3 py-2 text-sm text-gray-900 dark:text-gray-100 align-top"
                        onDoubleClick={() => handleCellDoubleClick(rIdx, c.key, display)}
                        tabIndex={0}
                        role="gridcell"
                        aria-colindex={columns.findIndex((x) => x.key === c.key) + 1}
                      >
                        {isEditing ? (
                          <input
                            autoFocus
                            defaultValue={editing.initialValue}
                            onKeyDown={(e) => onKeyDownEditor(e, rIdx, c.key)}
                            onBlur={(e) => saveEdit(rIdx, c.key, (e.target as HTMLInputElement).value)}
                            className="w-full text-sm px-2 py-1 rounded border border-gray-300 dark:border-gray-700 bg-white dark:lg:bg-zinc-900 focus:outline-none focus:ring-2 focus:ring-primary-500"
                            aria-label={`Edit ${c.header}`}
                          />
                        ) : (
                          <span>{display as any}</span>
                        )}
                      </td>
                    )
                  })}
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>

      {/* Pagination */}
      <div className="flex items-center justify-between px-3 py-2 border-t border-gray-200 dark:border-gray-800 text-sm">
        <div className="text-gray-600 dark:text-gray-300">
          Page <span className="font-medium">{page}</span> of <span className="font-medium">{totalPages}</span> • {total}{' '}
          rows
        </div>
        <div className="flex items-center gap-2">
          <button
            className="px-2.5 py-1.5 rounded-md border border-gray-300 dark:border-gray-700 text-gray-700 dark:text-gray-300 disabled:opacity-50 hover:bg-gray-100 dark:hover:bg-gray-800 focus:outline-none focus:ring-2 focus:ring-primary-500"
            onClick={() => setPage(1)}
            disabled={page === 1}
            aria-label="First page"
          >
            «
          </button>
          <button
            className="px-2.5 py-1.5 rounded-md border border-gray-300 dark:border-gray-700 text-gray-700 dark:text-gray-300 disabled:opacity-50 hover:bg-gray-100 dark:hover:bg-gray-800 focus:outline-none focus:ring-2 focus:ring-primary-500"
            onClick={() => setPage((p) => Math.max(1, p - 1))}
            disabled={page === 1}
            aria-label="Previous page"
          >
            ‹
          </button>
          <button
            className="px-2.5 py-1.5 rounded-md border border-gray-300 dark:border-gray-700 text-gray-700 dark:text-gray-300 disabled:opacity-50 hover:bg-gray-100 dark:hover:bg-gray-800 focus:outline-none focus:ring-2 focus:ring-primary-500"
            onClick={() => setPage((p) => Math.min(totalPages, p + 1))}
            disabled={page === totalPages}
            aria-label="Next page"
          >
            ›
          </button>
          <button
            className="px-2.5 py-1.5 rounded-md border border-gray-300 dark:border-gray-700 text-gray-700 dark:text-gray-300 disabled:opacity-50 hover:bg-gray-100 dark:hover:bg-gray-800 focus:outline-none focus:ring-2 focus:ring-primary-500"
            onClick={() => setPage(totalPages)}
            disabled={page === totalPages}
            aria-label="Last page"
          >
            »
          </button>
        </div>
      </div>
    </div>
  )
}