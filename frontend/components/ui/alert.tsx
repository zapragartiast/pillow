import * as React from "react"
import { AlertTriangle, CheckCircle, Info, XCircle } from "lucide-react"

export interface AlertProps {
  children?: React.ReactNode
  variant?: "default" | "destructive" | "success" | "warning"
  className?: string
}

export interface AlertDescriptionProps {
  children?: React.ReactNode
  className?: string
}

export function Alert({ children, variant = "default", className = "" }: AlertProps) {
  const variants = {
    default: "border-blue-200 bg-blue-50 text-blue-800 dark:border-blue-800 dark:bg-blue-950 dark:text-blue-200",
    destructive: "border-red-200 bg-red-50 text-red-800 dark:border-red-800 dark:bg-red-950 dark:text-red-200",
    success: "border-green-200 bg-green-50 text-green-800 dark:border-green-800 dark:bg-green-950 dark:text-green-200",
    warning: "border-yellow-200 bg-yellow-50 text-yellow-800 dark:border-yellow-800 dark:bg-yellow-950 dark:text-yellow-200"
  }

  const icons = {
    default: <Info className="h-4 w-4" />,
    destructive: <XCircle className="h-4 w-4" />,
    success: <CheckCircle className="h-4 w-4" />,
    warning: <AlertTriangle className="h-4 w-4" />
  }

  return (
    <div className={`flex items-start space-x-3 rounded-lg border p-4 ${variants[variant]} ${className}`}>
      <div className="flex-shrink-0">
        {icons[variant]}
      </div>
      <div className="flex-1">
        {children}
      </div>
    </div>
  )
}

export function AlertDescription({ children, className = "" }: AlertDescriptionProps) {
  return (
    <div className={`text-sm ${className}`}>
      {children}
    </div>
  )
}