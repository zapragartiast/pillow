import { registerAction } from './actions'

interface RegisterPageProps {
  searchParams: Promise<{ [key: string]: string | string[] | undefined }>
}

export default async function RegisterPage({ searchParams }: RegisterPageProps) {
  const params = await searchParams
  const error = typeof params.error === 'string' ? params.error : null

  return (
    <div className="max-w-md mx-auto mt-24 bg-white border rounded-lg shadow p-6">
      <h1 className="text-2xl font-semibold mb-4">Create account</h1>

      {error && (
        <div className="mb-4 p-3 bg-red-50 border border-red-200 rounded-md">
          <p className="text-sm text-red-600">{decodeURIComponent(error)}</p>
        </div>
      )}

      <form action={registerAction} className="space-y-4">
        <div>
          <label htmlFor="username" className="block text-sm font-medium text-slate-700">
            Username
          </label>
          <input
            id="username"
            name="username"
            type="text"
            required
            className="mt-1 block w-full rounded-md border border-slate-300 px-3 py-2 focus:border-slate-500 focus:outline-none focus:ring-1 focus:ring-slate-500"
            placeholder="your username"
          />
        </div>

        <div>
          <label htmlFor="email" className="block text-sm font-medium text-slate-700">
            Email
          </label>
          <input
            id="email"
            name="email"
            type="email"
            required
            className="mt-1 block w-full rounded-md border border-slate-300 px-3 py-2 focus:border-slate-500 focus:outline-none focus:ring-1 focus:ring-slate-500"
            placeholder="you@example.com"
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
            minLength={6}
            className="mt-1 block w-full rounded-md border border-slate-300 px-3 py-2 focus:border-slate-500 focus:outline-none focus:ring-1 focus:ring-slate-500"
            placeholder="••••••••"
          />
        </div>

        <div className="flex items-center justify-between">
          <button
            type="submit"
            className="inline-flex items-center px-4 py-2 bg-slate-800 text-white rounded-md hover:bg-slate-700 focus:outline-none focus:ring-2 focus:ring-slate-500 focus:ring-offset-2 disabled:opacity-50"
          >
            Create account
          </button>

          <a
            href="/login"
            className="text-sm text-slate-600 hover:underline focus:outline-none focus:ring-2 focus:ring-slate-500 focus:ring-offset-2 rounded"
          >
            Already have an account?
          </a>
        </div>
      </form>
    </div>
  )
}