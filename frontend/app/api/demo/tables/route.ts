import { NextResponse } from 'next/server'

type Row = {
  id: number
  name: string
  email: string
  role: 'admin' | 'editor' | 'viewer'
  createdAt: string
  status: 'active' | 'invited' | 'disabled'
}

function generateSeed(count = 250): Row[] {
  const roles: Row['role'][] = ['admin', 'editor', 'viewer']
  const statuses: Row['status'][] = ['active', 'invited', 'disabled']
  const rows: Row[] = []
  for (let i = 1; i <= count; i++) {
    const role = roles[i % roles.length]
    const status = statuses[i % statuses.length]
    const createdAt = new Date(Date.now() - i * 8640000).toISOString()
    rows.push({
      id: i,
      name: `User ${i}`,
      email: `user${i}@example.com`,
      role,
      createdAt,
      status,
    })
  }
  return rows
}

// In-memory cache for demo data (ephemeral per server instance)
let CACHE: Row[] | null = null
function getData(): Row[] {
  if (!CACHE) CACHE = generateSeed()
  return CACHE
}

export async function GET(request: Request) {
  const { searchParams } = new URL(request.url)

  const page = Number(searchParams.get('page') || '1')
  const pageSize = Number(searchParams.get('pageSize') || '20')
  const sortKey = searchParams.get('sortKey') as keyof Row | null
  const sortDir = searchParams.get('sortDir') as 'asc' | 'desc' | null
  const q = searchParams.get('q')?.toLowerCase() || ''

  let data = [...getData()]

  // Filter (global q across a few fields)
  if (q) {
    data = data.filter((r) =>
      String(r.id).includes(q) ||
      r.name.toLowerCase().includes(q) ||
      r.email.toLowerCase().includes(q) ||
      r.role.toLowerCase().includes(q) ||
      r.status.toLowerCase().includes(q)
    )
  }

  // Sorting
  if (sortKey && sortDir) {
    data.sort((a, b) => {
      const av = a[sortKey]
      const bv = b[sortKey]
      if (av === bv) return 0
      if (av > bv) return sortDir === 'asc' ? 1 : -1
      return sortDir === 'asc' ? -1 : 1
    })
  }

  const total = data.length
  const start = (page - 1) * pageSize
  const end = start + pageSize
  const rows = data.slice(start, end)

  return NextResponse.json({ rows, total })
}

export async function PATCH(request: Request) {
  // Simple inline-edit save demo:
  // body: { id, key, value }
  try {
    const payload = await request.json() as { id: number; key: keyof Row; value: any }
    if (!payload || !payload.id || !payload.key) {
      return NextResponse.json({ error: 'Invalid payload' }, { status: 400 })
    }
    const data = getData()
    const idx = data.findIndex((r) => r.id === payload.id)
    if (idx === -1) {
      return NextResponse.json({ error: 'Row not found' }, { status: 404 })
    }
    // Basic validation for constrained fields
    if (payload.key === 'role' && !['admin', 'editor', 'viewer'].includes(payload.value)) {
      return NextResponse.json({ error: 'Invalid role' }, { status: 400 })
    }
    if (payload.key === 'status' && !['active', 'invited', 'disabled'].includes(payload.value)) {
      return NextResponse.json({ error: 'Invalid status' }, { status: 400 })
    }

    // Persist in-memory (demo)
    ;(data[idx] as any)[payload.key] = payload.value
    return NextResponse.json({ success: true })
  } catch (e) {
    return NextResponse.json({ error: 'Bad Request' }, { status: 400 })
  }
}