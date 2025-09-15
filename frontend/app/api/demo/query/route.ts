import { NextResponse } from 'next/server'

type ResultSet = {
  columns: { key: string; label: string }[]
  rows: Record<string, any>[]
  rowCount: number
  executionTimeMs: number
  warnings?: string[]
}

function mockRun(sql: string): ResultSet {
  const start = Date.now()
  const lower = sql.trim().toLowerCase()

  // Very small DSL to return different sample datasets based on SQL keywords
  if (lower.includes('from users') || lower.includes('select * users')) {
    const rows = Array.from({ length: 10 }).map((_, i) => ({
      id: i + 1,
      name: `User ${i + 1}`,
      email: `user${i + 1}@example.com`,
      role: ['admin', 'editor', 'viewer'][i % 3],
      created_at: new Date(Date.now() - i * 3600_000).toISOString(),
    }))
    return {
      columns: [
        { key: 'id', label: 'id' },
        { key: 'name', label: 'name' },
        { key: 'email', label: 'email' },
        { key: 'role', label: 'role' },
        { key: 'created_at', label: 'created_at' },
      ],
      rows,
      rowCount: rows.length,
      executionTimeMs: Date.now() - start,
    }
  }

  if (lower.includes('count(') || lower.startsWith('select count')) {
    const rows = [{ count: 12345 }]
    return {
      columns: [{ key: 'count', label: 'count' }],
      rows,
      rowCount: rows.length,
      executionTimeMs: Date.now() - start,
    }
  }

  // Default: generic dataset
  const rows = Array.from({ length: 8 }).map((_, i) => ({
    id: i + 1,
    value: Math.round(Math.random() * 1000),
    ts: new Date(Date.now() - i * 600_000).toISOString(),
  }))
  return {
    columns: [
      { key: 'id', label: 'id' },
      { key: 'value', label: 'value' },
      { key: 'ts', label: 'timestamp' },
    ],
    rows,
    rowCount: rows.length,
    executionTimeMs: Date.now() - start,
    warnings: ['Using demo execution engine (mocked).'],
  }
}

export async function POST(req: Request) {
  try {
    const body = await req.json() as { sql: string; params?: any }
    if (!body?.sql || typeof body.sql !== 'string') {
      return NextResponse.json({ error: 'Missing SQL string' }, { status: 400 })
    }
    const result = mockRun(body.sql)
    return NextResponse.json(result)
  } catch (e) {
    return NextResponse.json({ error: 'Malformed request' }, { status: 400 })
  }
}