"use client"

import { useState, useTransition } from "react"
import Link from "next/link"
import { Eye, EyeOff } from "lucide-react"
import { registerAction } from "@/app/register/actions"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Checkbox } from "@/components/ui/checkbox"

interface RegisterFormProps {
  error?: string | null
}

export function RegisterForm({ error }: RegisterFormProps) {
  const [isPending, startTransition] = useTransition()
  const [showPassword, setShowPassword] = useState(false)
  const [showConfirmPassword, setShowConfirmPassword] = useState(false)
  const [formData, setFormData] = useState({
    username: '',
    email: '',
    password: '',
    confirmPassword: ''
  })
  const [fieldErrors, setFieldErrors] = useState<Record<string, string>>({})

  const validateForm = () => {
    const errors: Record<string, string> = {}

    if (!formData.username.trim()) {
      errors.username = 'Username is required'
    } else if (formData.username.length < 3) {
      errors.username = 'Username must be at least 3 characters'
    }

    if (!formData.email.trim()) {
      errors.email = 'Email is required'
    } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(formData.email)) {
      errors.email = 'Please enter a valid email address'
    }

    if (!formData.password) {
      errors.password = 'Password is required'
    } else if (formData.password.length < 6) {
      errors.password = 'Password must be at least 6 characters'
    }

    if (!formData.confirmPassword) {
      errors.confirmPassword = 'Please confirm your password'
    } else if (formData.password !== formData.confirmPassword) {
      errors.confirmPassword = 'Passwords do not match'
    }

    setFieldErrors(errors)
    return Object.keys(errors).length === 0
  }

  const handleInputChange = (field: string, value: string) => {
    setFormData(prev => ({ ...prev, [field]: value }))
    // Clear field error when user starts typing
    if (fieldErrors[field]) {
      setFieldErrors(prev => ({ ...prev, [field]: '' }))
    }
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()

    if (!validateForm()) {
      return
    }

    startTransition(async () => {
      try {
        const formDataToSubmit = new FormData()
        formDataToSubmit.append('username', formData.username)
        formDataToSubmit.append('email', formData.email)
        formDataToSubmit.append('password', formData.password)

        await registerAction(formDataToSubmit)
      } catch (error) {
        // Error handling is done in the server action
      }
    })
  }

  return (
    <div className="min-h-screen flex">
      {/* Left Section */}
      <div className="hidden lg:flex lg:w-1/2 bg-gradient-to-br from-emerald-600 via-teal-600 to-cyan-700 relative overflow-hidden">
        <div className="absolute inset-0 bg-black/20"></div>
        <div className="relative z-10 flex flex-col justify-center px-12 text-white">
          <div className="mb-8">
            <p className="text-sm font-medium tracking-wider uppercase opacity-80">
              START YOUR JOURNEY
            </p>
          </div>
          <div className="max-w-md">
            <h1 className="text-4xl lg:text-5xl font-serif font-bold mb-6 leading-tight">
              Join Our Community
            </h1>
            <p className="text-lg opacity-90 leading-relaxed">
              Create your account and unlock access to powerful tools and insights that will transform the way you work and collaborate.
            </p>
          </div>
        </div>
        {/* Abstract background pattern */}
        <div className="absolute inset-0 opacity-10">
          <div className="absolute top-20 left-20 w-32 h-32 bg-white rounded-full"></div>
          <div className="absolute bottom-32 right-16 w-24 h-24 bg-white rounded-full"></div>
          <div className="absolute top-1/2 left-1/3 w-16 h-16 bg-white rounded-full"></div>
        </div>
      </div>

      {/* Right Section */}
      <div className="w-full lg:w-1/2 flex items-center justify-center p-8 bg-white">
        <div className="w-full max-w-md space-y-8">
          {/* Logo and App Name */}
          <div className="text-center">
            <div className="mx-auto h-12 w-12 bg-gradient-to-r from-emerald-600 to-teal-600 rounded-xl flex items-center justify-center mb-4">
              <svg
                className="h-6 w-6 text-white"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M18 9v3m0 0v3m0-3h3m-3 0h-3m-2-5a4 4 0 11-8 0 4 4 0 018 0zM3 20a6 6 0 0112 0v1H3v-1z"
                />
              </svg>
            </div>
          </div>

          {/* Welcome Text */}
          <div className="text-center">
            <h2 className="text-3xl font-bold text-gray-900 mb-2">Create Account</h2>
            <p className="text-gray-600">Join us today and get started</p>
          </div>

          {/* Error Message */}
          {error && (
            <div className="p-4 text-sm text-red-700 bg-red-50 border border-red-200 rounded-lg">
              <div className="flex items-center">
                <svg className="w-4 h-4 mr-2 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20">
                  <path fillRule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z" clipRule="evenodd" />
                </svg>
                {decodeURIComponent(error)}
              </div>
            </div>
          )}

          {/* Registration Form */}
          <form onSubmit={handleSubmit} className="space-y-6">
            <div className="space-y-2">
              <Label htmlFor="username" className="text-sm font-medium text-gray-700">
                Username
              </Label>
              <Input
                id="username"
                name="username"
                type="text"
                placeholder="Choose a username"
                required
                disabled={isPending}
                value={formData.username}
                onChange={(e) => handleInputChange('username', e.target.value)}
                className="h-11 mt-[5px] text-gray-700"
              />
              {fieldErrors.username && (
                <p className="text-sm text-red-600">{fieldErrors.username}</p>
              )}
            </div>

            <div className="space-y-2">
              <Label htmlFor="email" className="text-sm font-medium text-gray-700">
                Email
              </Label>
              <Input
                id="email"
                name="email"
                type="email"
                placeholder="Enter your email"
                required
                disabled={isPending}
                value={formData.email}
                onChange={(e) => handleInputChange('email', e.target.value)}
                className="h-11 mt-[5px] text-gray-700"
              />
              {fieldErrors.email && (
                <p className="text-sm text-red-600">{fieldErrors.email}</p>
              )}
            </div>

            <div className="space-y-2">
              <Label htmlFor="password" className="text-sm font-medium text-gray-700">
                Password
              </Label>
              <div className="relative">
                <Input
                  id="password"
                  name="password"
                  type={showPassword ? "text" : "password"}
                  placeholder="Create a password"
                  required
                  disabled={isPending}
                  value={formData.password}
                  onChange={(e) => handleInputChange('password', e.target.value)}
                  className="h-11 pr-10 mt-[5px] text-gray-700"
                />
                <button
                  type="button"
                  className="absolute inset-y-0 right-0 pr-3 flex items-center"
                  onClick={() => setShowPassword(!showPassword)}
                >
                  {showPassword ? (
                    <EyeOff className="h-4 w-4 text-gray-400" />
                  ) : (
                    <Eye className="h-4 w-4 text-gray-400" />
                  )}
                </button>
              </div>
              {fieldErrors.password && (
                <p className="text-sm text-red-600">{fieldErrors.password}</p>
              )}
            </div>

            <div className="space-y-2">
              <Label htmlFor="confirmPassword" className="text-sm font-medium text-gray-700">
                Confirm Password
              </Label>
              <div className="relative">
                <Input
                  id="confirmPassword"
                  name="confirmPassword"
                  type={showConfirmPassword ? "text" : "password"}
                  placeholder="Confirm your password"
                  required
                  disabled={isPending}
                  value={formData.confirmPassword}
                  onChange={(e) => handleInputChange('confirmPassword', e.target.value)}
                  className="h-11 pr-10 mt-[5px] text-gray-700"
                />
                <button
                  type="button"
                  className="absolute inset-y-0 right-0 pr-3 flex items-center"
                  onClick={() => setShowConfirmPassword(!showConfirmPassword)}
                >
                  {showConfirmPassword ? (
                    <EyeOff className="h-4 w-4 text-gray-400" />
                  ) : (
                    <Eye className="h-4 w-4 text-gray-400" />
                  )}
                </button>
              </div>
              {fieldErrors.confirmPassword && (
                <p className="text-sm text-red-600">{fieldErrors.confirmPassword}</p>
              )}
            </div>

            {/* Terms and Conditions */}
            <div className="flex items-start space-x-2">
              <Checkbox
                id="terms"
                name="terms"
                required
                className="mt-1"
              />
              <div className="text-sm">
                <Label htmlFor="terms" className="text-gray-600 cursor-pointer">
                  I agree to the{" "}
                  <Link href="/terms" className="text-blue-600 hover:text-blue-500 underline">
                    Terms of Service
                  </Link>{" "}
                  and{" "}
                  <Link href="/privacy" className="text-blue-600 hover:text-blue-500 underline">
                    Privacy Policy
                  </Link>
                </Label>
              </div>
            </div>

            {/* Sign Up Button */}
            <Button
              type="submit"
              className="w-full h-11 bg-black hover:bg-gray-800 text-white font-medium"
              disabled={isPending}
            >
              {isPending ? (
                <div className="flex items-center justify-center">
                  <svg className="animate-spin -ml-1 mr-2 h-4 w-4 text-white" fill="none" viewBox="0 0 24 24">
                    <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                    <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                  </svg>
                  Creating account...
                </div>
              ) : (
                "Create Account"
              )}
            </Button>

            {/* Sign In Link */}
            <div className="text-center">
              <p className="text-sm text-gray-600">
                Already have an account?{" "}
                <Link
                  href="/login"
                  className="font-medium text-blue-600 hover:text-blue-500"
                >
                  Sign In
                </Link>
              </p>
            </div>
          </form>
        </div>
      </div>
    </div>
  )
}