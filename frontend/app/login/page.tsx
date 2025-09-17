import { cookies } from 'next/headers'
import { redirect } from 'next/navigation'
import { LoginForm } from '@/components/login-form'

interface LoginPageProps {
  searchParams?: { [key: string]: string | string[] | undefined }
}

export default async function LoginPage({ searchParams }: LoginPageProps) {
  const cookieStore = await cookies()
  const token = cookieStore.get('pillow_token')?.value

  if (token) {
    redirect('/dashboard')
  }

  const errorParam = typeof searchParams?.error === 'string' ? searchParams?.error : null

  return <LoginForm error={errorParam} />
}
