import { RegisterForm } from '@/components/register-form'

interface RegisterPageProps {
  searchParams: Promise<{ [key: string]: string | string[] | undefined }>
}

export default async function RegisterPage({ searchParams }: RegisterPageProps) {
  const params = await searchParams
  const error = typeof params.error === 'string' ? params.error : null

  return <RegisterForm error={error} />
}