'use client'

import React, { useCallback } from 'react'
import TableView, { Column, FetchParams, FetchResult } from './TableView'

type DataExplorerClientProps = {
  className?: string
}

/**
 * DataExplorerClient
 * - Wires TableView to demo API: /api/demo/tables
 * - Supports server-side pagination, sorting, filtering, and inline edits via PATCH
 */
export default function DataExplorerClient({ className = '' }: DataExplorerClientProps) {
  const columns: Column[] = [
    { key: 'id', header: 'ID', width: 90 },
    { key: 'name', header: 'Name', editable: true, width: 220 },
    { key: 'email', header: 'Email', editable: true, width: 260 },
    { key: 'role', header: 'Role', editable: true, width: 140 },
    { key: 'status', header: 'Status', editable: true, width: 140 },
    { key: 'createdAt', header: 'Created At', width: 220, accessor: (r) => new Date(r.createdAt).toLocaleString() },
  ]

  const fetchPage = useCallback(async (params: FetchParams): Promise<FetchResult> => {
    const url = new URL('/api/demo/tables', window.location.origin)
    url.searchParams.set('page', String(params.page))
    url.searchParams.set('pageSize', String(params.pageSize))
    if (params.sortKey) url.searchParams.set('sortKey', params.sortKey)
    if (params.sortDir) url.searchParams.set('sortDir', params.sortDir)
    if (params.filters?.q) url.searchParams.set('q', params.filters.q)

    const res = await fetch(url.toString(), { cache: 'no-store' })
    if (!res.ok) {
      throw new Error('Failed to load data')
    }
    const json = await res.json()
    return { rows: json.rows, total: json.total }
  }, [])

  const onEditCell = useCallback(async (row: any, key: string, value: string | number | null) => {
    // Persist change via demo PATCH API
    const res = await fetch('/api/demo/tables', {
      method: 'PATCH',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ id: row.id, key, value }),
    })
    if (!res.ok) {
      const data = await res.json().catch(() => ({}))
      throw new Error(data?.error || 'Failed to save')
    }
  }, [])

  return (
    <div className={className}>
      <TableView
        ariaLabel="Data Explorer table"
        columns={columns}
        fetchPage={fetchPage}
        onEditCell={onEditCell}
        initialPageSize={20}
      />
    </div>
  )
}