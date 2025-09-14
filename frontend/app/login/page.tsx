import { LoginForm } from '@/components/login-form'

interface LoginPageProps {
  searchParams: Promise<{ [key: string]: string | string[] | undefined }>
}

export default async function LoginPage({ searchParams }: LoginPageProps) {
  const params = await searchParams
  const error = typeof params.error === 'string' ? params.error : null

  return <LoginForm error={error} />
}