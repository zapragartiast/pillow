"use client"

import { useState, useTransition } from "react"
import Link from "next/link"
import { Eye, EyeOff } from "lucide-react"
import { loginAction } from "@/app/login/actions"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Checkbox } from "@/components/ui/checkbox"

interface LoginFormProps {
  error?: string | null
}

export function LoginForm({ error }: LoginFormProps) {
  const [isPending, startTransition] = useTransition()
  const [showPassword, setShowPassword] = useState(false)

  return (
    <div className="min-h-screen flex">
      {/* Left Section */}
      <div className="hidden lg:flex lg:w-1/2 bg-gradient-to-br from-purple-600 via-blue-600 to-indigo-700 relative overflow-hidden">
        <div className="absolute inset-0 bg-black/20"></div>
        <div className="relative z-10 flex flex-col justify-center px-12 text-white">
          <div className="mb-8">
            <p className="text-sm font-medium tracking-wider uppercase opacity-80">
              A WISE QUOTE
            </p>
          </div>
          <div className="max-w-md">
            <h1 className="text-4xl lg:text-5xl font-serif font-bold mb-6 leading-tight">
              Get Everything You Want
            </h1>
            <p className="text-lg opacity-90 leading-relaxed">
              You can get everything you want if you work hard, trust the process, and stick to the plan.
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
            <div className="mx-auto h-12 w-12 bg-gradient-to-r from-blue-600 to-indigo-600 rounded-xl flex items-center justify-center mb-4">
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
                  d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"
                />
              </svg>
            </div>
          </div>

          {/* Welcome Text */}
          <div className="text-center">
            <h2 className="text-3xl font-bold text-gray-900 mb-2">Welcome Back</h2>
            <p className="text-gray-600">Enter your email and password to access your account</p>
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

          {/* Login Form */}
          <form action={loginAction} className="space-y-6">
            <div className="space-y-2">
              <Label htmlFor="identifier" className="text-sm font-medium text-gray-700">
                Email
              </Label>
              <Input
                id="identifier"
                name="identifier"
                type="email"
                placeholder="Enter your email"
                required
                disabled={isPending}
                className="h-11 mt-[5px] text-gray-700"
              />
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
                  placeholder="Enter your password"
                  required
                  disabled={isPending}
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
            </div>

            {/* Remember me and Forgot Password */}
            <div className="flex items-center justify-between">
              <div className="flex items-center space-x-2">
                <Checkbox id="remember" name="remember" />
                <Label htmlFor="remember" className="text-sm text-gray-600 cursor-pointer">
                  Remember me
                </Label>
              </div>
              <Link
                href="/forgot-password"
                className="text-sm text-blue-600 hover:text-blue-500"
              >
                Forgot Password?
              </Link>
            </div>

            {/* Sign In Button */}
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
                  Signing in...
                </div>
              ) : (
                "Sign In"
              )}
            </Button>

            {/* Divider */}
            <div className="relative">
              <div className="absolute inset-0 flex items-center">
                <div className="w-full border-t border-gray-300"></div>
              </div>
              <div className="relative flex justify-center text-sm">
                <span className="px-2 bg-white text-gray-500">Or</span>
              </div>
            </div>

            {/* Google Sign In Button */}
            <Button
              type="button"
              variant="outline"
              className="w-full h-11 border border-gray-300 hover:bg-gray-50 text-gray-800"
            >
              <svg className="w-4 h-4 mr-2" viewBox="0 0 24 24">
                <path fill="#4285F4" d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"/>
                <path fill="#34A853" d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"/>
                <path fill="#FBBC05" d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"/>
                <path fill="#EA4335" d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"/>
              </svg>
              Sign in with Google
            </Button>

            {/* Sign Up Link */}
            <div className="text-center">
              <p className="text-sm text-gray-600">
                Don't have an account?{" "}
                <Link
                  href="/register"
                  className="font-medium text-blue-600 hover:text-blue-500"
                >
                  Sign Up
                </Link>
              </p>
            </div>
          </form>
        </div>
      </div>
    </div>
  )
}