'use client'

import React, { useState, useEffect } from 'react'
import { Save, Loader2 } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Checkbox } from '@/components/ui/checkbox'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'

interface GlobalCustomField {
  id: string
  name: string
  label: string
  type: 'text' | 'textarea' | 'number' | 'email' | 'phone' | 'date' | 'boolean' | 'select' | 'multiselect'
  required: boolean
  options?: string[]
  validation?: {
    min_length?: number
    max_length?: number
    min?: number
    max?: number
    pattern?: string
  }
  order: number
}

interface CustomFieldWithValue extends GlobalCustomField {
  value?: string
}

interface UserProfileCustomFieldsProps {
  userId: string
  onFieldsChange?: (fields: CustomFieldWithValue[]) => void
}

export default function UserProfileCustomFields({ userId, onFieldsChange }: UserProfileCustomFieldsProps) {
  const [fields, setFields] = useState<CustomFieldWithValue[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [isSaving, setIsSaving] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [hasChanges, setHasChanges] = useState(false)

  useEffect(() => {
    loadUserFields()
  }, [userId])

  const loadUserFields = async () => {
    try {
      setIsLoading(true)
      const response = await fetch(`/api/users/${userId}/custom-field-values`, {
        credentials: 'include',
        cache: 'no-store'
      })

      if (!response.ok) {
        throw new Error('Failed to load custom fields')
      }

      const data: CustomFieldWithValue[] = await response.json()
      setFields(data)
      onFieldsChange?.(data)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load fields')
    } finally {
      setIsLoading(false)
    }
  }

  const updateFieldValue = (fieldId: string, value: any) => {
    setFields(prevFields =>
      prevFields.map(field =>
        field.id === fieldId ? { ...field, value: value } : field
      )
    )
    setHasChanges(true)
  }

  const handleMultiselectChange = (fieldId: string, option: string, checked: boolean) => {
    setFields(prevFields =>
      prevFields.map(field => {
        if (field.id === fieldId) {
          let currentValues: string[] = []
          if (field.value) {
            try {
              currentValues = JSON.parse(field.value)
            } catch {
              currentValues = field.value.split(',').map(v => v.trim())
            }
          }

          let newValues: string[]
          if (checked) {
            newValues = [...currentValues, option]
          } else {
            newValues = currentValues.filter(v => v !== option)
          }

          return { ...field, value: JSON.stringify(newValues) }
        }
        return field
      })
    )
    setHasChanges(true)
  }

  const saveChanges = async () => {
    try {
      setIsSaving(true)
      setError(null)

      // Prepare field values for API
      const fieldValues: { [key: string]: any } = {}
      fields.forEach(field => {
        if (field.value !== undefined && field.value !== '') {
          fieldValues[field.id] = field.value
        }
      })

      const response = await fetch(`/api/users/${userId}/custom-field-values`, {
        method: 'PUT',
        credentials: 'include',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({ field_values: fieldValues })
      })

      if (!response.ok) {
        const errorData = await response.json()
        throw new Error(errorData.message || 'Failed to save changes')
      }

      setHasChanges(false)
      await loadUserFields() // Reload to get fresh data
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to save changes')
    } finally {
      setIsSaving(false)
    }
  }

  const renderFieldInput = (field: CustomFieldWithValue) => {
    const commonProps = {
      id: field.id,
      className: "w-full"
    }

    switch (field.type) {
      case 'text':
      case 'email':
      case 'phone':
        return (
          <Input
            {...commonProps}
            type={field.type === 'email' ? 'email' : field.type === 'phone' ? 'tel' : 'text'}
            value={field.value || ''}
            onChange={(e) => updateFieldValue(field.id, e.target.value)}
            placeholder={`Enter ${field.label.toLowerCase()}`}
          />
        )

      case 'textarea':
        return (
          <Textarea
            {...commonProps}
            value={field.value || ''}
            onChange={(e) => updateFieldValue(field.id, e.target.value)}
            placeholder={`Enter ${field.label.toLowerCase()}`}
            rows={3}
          />
        )

      case 'number':
        return (
          <Input
            {...commonProps}
            type="number"
            value={field.value || ''}
            onChange={(e) => updateFieldValue(field.id, e.target.value)}
            min={field.validation?.min}
            max={field.validation?.max}
            placeholder={`Enter ${field.label.toLowerCase()}`}
          />
        )

      case 'date':
        return (
          <Input
            {...commonProps}
            type="date"
            value={field.value || ''}
            onChange={(e) => updateFieldValue(field.id, e.target.value)}
          />
        )

      case 'boolean':
        return (
          <div className="flex items-center space-x-2">
            <Checkbox
              id={field.id}
              checked={field.value === 'true' || field.value === '1'}
              onCheckedChange={(checked) => updateFieldValue(field.id, checked ? 'true' : 'false')}
            />
            <Label htmlFor={field.id} className="text-sm">
              {field.label}
            </Label>
          </div>
        )

      case 'select':
        return (
          <Select
            value={field.value || ''}
            onValueChange={(value) => updateFieldValue(field.id, value)}
          >
            <SelectTrigger>
              <SelectValue placeholder={`Select ${field.label.toLowerCase()}`} />
            </SelectTrigger>
            <SelectContent>
              {field.options?.map((option) => (
                <SelectItem key={option} value={option}>
                  {option}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        )

      case 'multiselect':
        let selectedValues: string[] = []
        if (field.value) {
          try {
            selectedValues = JSON.parse(field.value)
          } catch {
            selectedValues = field.value.split(',').map(v => v.trim())
          }
        }

        return (
          <div className="space-y-2">
            {field.options?.map((option) => (
              <div key={option} className="flex items-center space-x-2">
                <Checkbox
                  id={`${field.id}-${option}`}
                  checked={selectedValues.includes(option)}
                  onCheckedChange={(checked) => handleMultiselectChange(field.id, option, !!checked)}
                />
                <Label htmlFor={`${field.id}-${option}`} className="text-sm">
                  {option}
                </Label>
              </div>
            ))}
          </div>
        )

      default:
        return (
          <Input
            {...commonProps}
            value={field.value || ''}
            onChange={(e) => updateFieldValue(field.id, e.target.value)}
            placeholder={`Enter ${field.label.toLowerCase()}`}
          />
        )
    }
  }

  if (isLoading) {
    return (
      <div className="flex items-center justify-center p-8">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
      </div>
    )
  }

  if (fields.length === 0) {
    return (
      <Card>
        <CardContent className="flex flex-col items-center justify-center py-8">
          <div className="text-gray-400 mb-4">
            <div className="w-12 h-12 rounded-full bg-gray-100 flex items-center justify-center">
              📝
            </div>
          </div>
          <h3 className="text-lg font-medium text-gray-900 dark:text-white mb-2">No Custom Fields</h3>
          <p className="text-gray-600 dark:text-gray-400 text-center">
            There are no custom fields configured for your profile yet.
          </p>
        </CardContent>
      </Card>
    )
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-bold text-gray-900 dark:text-white">Profile Information</h2>
          <p className="text-gray-600 dark:text-gray-400">Manage your custom profile fields</p>
        </div>
        {hasChanges && (
          <Button
            onClick={saveChanges}
            disabled={isSaving}
            className="flex items-center space-x-2"
          >
            {isSaving ? (
              <Loader2 className="w-4 h-4 animate-spin" />
            ) : (
              <Save className="w-4 h-4" />
            )}
            <span>{isSaving ? 'Saving...' : 'Save Changes'}</span>
          </Button>
        )}
      </div>

      {error && (
        <Alert variant="destructive">
          <AlertDescription>{error}</AlertDescription>
        </Alert>
      )}

      <div className="grid gap-6">
        {fields.map((field) => (
          <Card key={field.id}>
            <CardHeader>
              <CardTitle className="flex items-center space-x-2">
                <span>{field.label}</span>
                {field.required && <span className="text-red-500">*</span>}
              </CardTitle>
            </CardHeader>
            <CardContent>
              {renderFieldInput(field)}
              {field.validation && (
                <div className="mt-2 text-sm text-gray-500">
                  {field.validation.min_length && `Minimum ${field.validation.min_length} characters`}
                  {field.validation.max_length && `Maximum ${field.validation.max_length} characters`}
                  {field.validation.min !== undefined && `Minimum value: ${field.validation.min}`}
                  {field.validation.max !== undefined && `Maximum value: ${field.validation.max}`}
                </div>
              )}
            </CardContent>
          </Card>
        ))}
      </div>

      {hasChanges && (
        <div className="sticky bottom-0 bg-white dark:bg-gray-900 border-t p-4 -mx-6 -mb-6">
          <div className="flex justify-end">
            <Button
              onClick={saveChanges}
              disabled={isSaving}
              className="flex items-center space-x-2"
            >
              {isSaving ? (
                <Loader2 className="w-4 h-4 animate-spin" />
              ) : (
                <Save className="w-4 h-4" />
              )}
              <span>{isSaving ? 'Saving...' : 'Save Changes'}</span>
            </Button>
          </div>
        </div>
      )}
    </div>
  )
}