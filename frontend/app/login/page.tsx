import { loginAction } from './actions'

interface LoginPageProps {
  searchParams: Promise<{ [key: string]: string | string[] | undefined }>
}

export default async function LoginPage({ searchParams }: LoginPageProps) {
  const params = await searchParams
  const error = typeof params.error === 'string' ? params.error : null

  return (
    <div className="max-w-md mx-auto mt-24 bg-white border rounded-lg shadow p-6">
      <h1 className="text-2xl font-semibold mb-4">Sign in</h1>

      {error && (
        <div className="mb-4 p-3 bg-red-50 border border-red-200 rounded-md">
          <p className="text-sm text-red-600">{decodeURIComponent(error)}</p>
        </div>
      )}

      <form action={loginAction} className="space-y-4">
        <div>
          <label htmlFor="identifier" className="block text-sm font-medium text-slate-700">
            Username or Email
          </label>
          <input
            id="identifier"
            name="identifier"
            type="text"
            required
            className="mt-1 block w-full rounded-md border border-slate-300 px-3 py-2 focus:border-slate-500 focus:outline-none focus:ring-1 focus:ring-slate-500"
            placeholder="you@example.com or username"
          />
        </div>

        <div>
          <label htmlFor="password" className="block text-sm font-medium text-slate-700">
            Password
          </label>
          <input
            id="password"
            name="password"
            type="password"
            required
            className="mt-1 block w-full rounded-md border border-slate-300 px-3 py-2 focus:border-slate-500 focus:outline-none focus:ring-1 focus:ring-slate-500"
            placeholder="••••••••"
          />
        </div>

        <div className="flex items-center justify-between">
          <button
            type="submit"
            className="inline-flex items-center px-4 py-2 bg-slate-800 text-white rounded-md hover:bg-slate-700 focus:outline-none focus:ring-2 focus:ring-slate-500 focus:ring-offset-2 disabled:opacity-50"
          >
            Sign in
          </button>

          <a
            href="/register"
            className="text-sm text-slate-600 hover:underline focus:outline-none focus:ring-2 focus:ring-slate-500 focus:ring-offset-2 rounded"
          >
            Create account
          </a>
        </div>
      </form>
    </div>
  )
}