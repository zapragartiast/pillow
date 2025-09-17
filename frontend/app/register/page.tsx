import { cookies } from 'next/headers'
import { redirect } from 'next/navigation'
import { RegisterForm } from '@/components/register-form'

interface RegisterPageProps {
  searchParams?: { [key: string]: string | string[] | undefined }
}

export default async function RegisterPage({ searchParams }: RegisterPageProps) {
  const cookieStore = await cookies()
  const token = cookieStore.get('pillow_token')?.value

  if (token) {
    redirect('/dashboard')
  }

  const errorParam = typeof searchParams?.error === 'string' ? searchParams?.error : null

  return <RegisterForm error={errorParam} />
}
