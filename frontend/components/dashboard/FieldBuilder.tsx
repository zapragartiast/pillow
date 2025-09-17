'use client'

import React, { useState, useEffect } from 'react'
import { Plus, Edit, Trash2, GripVertical, Save, X } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Checkbox } from '@/components/ui/checkbox'
import { Textarea } from '@/components/ui/textarea'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from '@/components/ui/dialog'
import { Alert, AlertDescription } from '@/components/ui/alert'

export interface CustomField {
  id: string
  name: string
  label: string
  type: 'text' | 'textarea' | 'number' | 'email' | 'phone' | 'date' | 'boolean' | 'select' | 'multiselect'
  required: boolean
  value?: any
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

export interface CustomFieldsData {
  fields: CustomField[]
  metadata: {
    version: string
    last_updated: string
  }
}

interface FieldBuilderProps {
  userId: string
  onFieldsChange?: (fields: CustomField[]) => void
}

export default function FieldBuilder({ userId, onFieldsChange }: FieldBuilderProps) {
  const [fields, setFields] = useState<CustomField[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [isCreateDialogOpen, setIsCreateDialogOpen] = useState(false)
  const [editingField, setEditingField] = useState<CustomField | null>(null)
  const [draggedField, setDraggedField] = useState<string | null>(null)

  useEffect(() => {
    loadFields()
  }, [userId])

  const loadFields = async () => {
    try {
      setIsLoading(true)
      const response = await fetch(`/api/users/${userId}/custom-field-values`, {
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        }
      })

      if (!response.ok) {
        throw new Error('Failed to load custom field values')
      }

      const data = await response.json()
      // Transform global fields with values to component format
      const transformedFields = data.map((field: any) => ({
        id: field.id,
        name: field.name,
        label: field.label,
        type: field.type,
        required: field.required,
        value: field.Value ? JSON.parse(field.Value) : undefined,
        options: field.Options || [],
        validation: field.Validation || {},
        order: field.Order || 0
      }))

      setFields(transformedFields)
      onFieldsChange?.(transformedFields)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load field values')
    } finally {
      setIsLoading(false)
    }
  }

  const createField = async (fieldData: Omit<CustomField, 'id' | 'order'>) => {
    // Global custom fields don't allow users to create new field definitions
    // Only admins can create field definitions through /api/global-custom-fields
    setError('You cannot create new custom fields. Please contact an administrator.')
    setIsCreateDialogOpen(false)
  }

  const updateField = async (fieldId: string, fieldData: Partial<CustomField>) => {
    try {
      // Users can only update their values for global fields, not the field definitions
      const fieldValues: Record<string, any> = {}
      if (fieldData.value !== undefined) {
        fieldValues[fieldId] = fieldData.value
      }

      const response = await fetch(`/api/users/${userId}/custom-field-values`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        },
        body: JSON.stringify({ field_values: fieldValues })
      })

      if (!response.ok) {
        const errorData = await response.json()
        throw new Error(errorData.message || 'Failed to update field value')
      }

      await loadFields()
      setEditingField(null)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to update field value')
    }
  }

  const clearFieldValue = async (fieldId: string) => {
    if (!confirm('Are you sure you want to clear this field value?')) return

    try {
      // Clear the user's value by setting it to null/undefined
      const fieldValues: Record<string, any> = {}
      fieldValues[fieldId] = null

      const response = await fetch(`/api/users/${userId}/custom-field-values`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        },
        body: JSON.stringify({ field_values: fieldValues })
      })

      if (!response.ok) {
        const errorData = await response.json()
        throw new Error(errorData.message || 'Failed to clear field value')
      }

      await loadFields()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to clear field value')
    }
  }

  // Reordering is handled automatically by the backend based on field.order

  const handleDragStart = (fieldId: string) => {
    setDraggedField(fieldId)
  }

  const handleDragOver = (e: React.DragEvent) => {
    e.preventDefault()
  }

  const handleDrop = (e: React.DragEvent, targetFieldId: string) => {
    e.preventDefault()

    if (!draggedField || draggedField === targetFieldId) return

    const draggedIndex = fields.findIndex(f => f.id === draggedField)
    const targetIndex = fields.findIndex(f => f.id === targetFieldId)

    if (draggedIndex === -1 || targetIndex === -1) return

    const newFields = [...fields]
    const [draggedFieldData] = newFields.splice(draggedIndex, 1)
    newFields.splice(targetIndex, 0, draggedFieldData)

    // Update order values
    newFields.forEach((field, index) => {
      field.order = index
    })

    setFields(newFields)
    setDraggedField(null)

    // Note: Field reordering is handled server-side by the global custom fields system
    // Users cannot reorder fields, only admins can through /api/global-custom-fields
  }

  if (isLoading) {
    return (
      <div className="flex items-center justify-center p-8">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-bold text-gray-900 dark:text-white">Custom Fields</h2>
          <p className="text-gray-600 dark:text-gray-400">Create and manage your custom user fields</p>
        </div>
        <Dialog open={isCreateDialogOpen} onOpenChange={setIsCreateDialogOpen}>
          <DialogTrigger asChild>
            <Button>
              <Plus className="w-4 h-4 mr-2" />
              Add Field
            </Button>
          </DialogTrigger>
          <DialogContent className="max-w-2xl">
            <DialogHeader>
              <DialogTitle>Create New Field</DialogTitle>
            </DialogHeader>
            <FieldEditor
              onSave={(fieldData) => createField(fieldData)}
              onCancel={() => setIsCreateDialogOpen(false)}
            />
          </DialogContent>
        </Dialog>
      </div>

      {error && (
        <Alert variant="destructive">
          <AlertDescription>{error}</AlertDescription>
        </Alert>
      )}

      <div className="space-y-4">
        {fields.length === 0 ? (
          <Card>
            <CardContent className="flex flex-col items-center justify-center py-12">
              <div className="text-gray-400 mb-4">
                <Plus className="w-12 h-12" />
              </div>
              <h3 className="text-lg font-medium text-gray-900 dark:text-white mb-2">No custom fields yet</h3>
              <p className="text-gray-600 dark:text-gray-400 text-center mb-4">
                Create your first custom field to extend user profiles with additional information.
              </p>
              <Button onClick={() => setIsCreateDialogOpen(true)}>
                <Plus className="w-4 h-4 mr-2" />
                Create Your First Field
              </Button>
            </CardContent>
          </Card>
        ) : (
          fields.map((field) => (
            <Card
              key={field.id}
              draggable
              onDragStart={() => handleDragStart(field.id)}
              onDragOver={handleDragOver}
              onDrop={(e) => handleDrop(e, field.id)}
              className="cursor-move"
            >
              <CardContent className="p-4">
                <div className="flex items-center justify-between">
                  <div className="flex items-center space-x-3">
                    <GripVertical className="w-5 h-5 text-gray-400" />
                    <div>
                      <h3 className="font-medium text-gray-900 dark:text-white">{field.label}</h3>
                      <p className="text-sm text-gray-600 dark:text-gray-400">
                        {field.type} • {field.name}
                        {field.required && <span className="text-red-500 ml-1">*</span>}
                      </p>
                    </div>
                  </div>
                  <div className="flex items-center space-x-2">
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => setEditingField(field)}
                    >
                      <Edit className="w-4 h-4" />
                    </Button>
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => clearFieldValue(field.id)}
                      className="text-red-600 hover:text-red-700"
                    >
                      <Trash2 className="w-4 h-4" />
                    </Button>
                  </div>
                </div>
              </CardContent>
            </Card>
          ))
        )}
      </div>

      {editingField && (
        <Dialog open={!!editingField} onOpenChange={() => setEditingField(null)}>
          <DialogContent className="max-w-2xl">
            <DialogHeader>
              <DialogTitle>Edit Field</DialogTitle>
            </DialogHeader>
            <FieldEditor
              field={editingField}
              onSave={(fieldData) => updateField(editingField.id, fieldData)}
              onCancel={() => setEditingField(null)}
            />
          </DialogContent>
        </Dialog>
      )}
    </div>
  )
}

interface FieldEditorProps {
  field?: CustomField
  onSave: (fieldData: any) => void
  onCancel: () => void
}

function FieldEditor({ field, onSave, onCancel }: FieldEditorProps) {
  const [formData, setFormData] = useState({
    name: field?.name || '',
    label: field?.label || '',
    type: field?.type || 'text',
    required: field?.required || false,
    options: field?.options || [],
    validation: field?.validation || {}
  })

  const [newOption, setNewOption] = useState('')
  const [errors, setErrors] = useState<Record<string, string>>({})

  const validateForm = () => {
    const newErrors: Record<string, string> = {}

    if (!formData.name.trim()) {
      newErrors.name = 'Field name is required'
    } else if (!/^[a-zA-Z_][a-zA-Z0-9_]*$/.test(formData.name)) {
      newErrors.name = 'Field name must start with a letter and contain only letters, numbers, and underscores'
    }

    if (!formData.label.trim()) {
      newErrors.label = 'Field label is required'
    }

    if ((formData.type === 'select' || formData.type === 'multiselect') && formData.options.length === 0) {
      newErrors.options = 'At least one option is required for select fields'
    }

    setErrors(newErrors)
    return Object.keys(newErrors).length === 0
  }

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (!validateForm()) return

    onSave(formData)
  }

  const addOption = () => {
    if (newOption.trim() && !formData.options.includes(newOption.trim())) {
      setFormData(prev => ({
        ...prev,
        options: [...prev.options, newOption.trim()]
      }))
      setNewOption('')
    }
  }

  const removeOption = (index: number) => {
    setFormData(prev => ({
      ...prev,
      options: prev.options.filter((_, i) => i !== index)
    }))
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      <div className="grid grid-cols-2 gap-4">
        <div>
          <Label htmlFor="name">Field Name *</Label>
          <Input
            id="name"
            value={formData.name}
            onChange={(e) => setFormData(prev => ({ ...prev, name: e.target.value }))}
            placeholder="e.g., full_name"
          />
          {errors.name && <p className="text-sm text-red-600 mt-1">{errors.name}</p>}
        </div>
        <div>
          <Label htmlFor="label">Field Label *</Label>
          <Input
            id="label"
            value={formData.label}
            onChange={(e) => setFormData(prev => ({ ...prev, label: e.target.value }))}
            placeholder="e.g., Full Name"
          />
          {errors.label && <p className="text-sm text-red-600 mt-1">{errors.label}</p>}
        </div>
      </div>

      <div className="grid grid-cols-2 gap-4">
        <div>
          <Label htmlFor="type">Field Type</Label>
          <Select value={formData.type} onValueChange={(value) => setFormData(prev => ({ ...prev, type: value as any }))}>
            <SelectTrigger>
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="text">Text</SelectItem>
              <SelectItem value="textarea">Textarea</SelectItem>
              <SelectItem value="number">Number</SelectItem>
              <SelectItem value="email">Email</SelectItem>
              <SelectItem value="phone">Phone</SelectItem>
              <SelectItem value="date">Date</SelectItem>
              <SelectItem value="boolean">Boolean</SelectItem>
              <SelectItem value="select">Select</SelectItem>
              <SelectItem value="multiselect">Multi-select</SelectItem>
            </SelectContent>
          </Select>
        </div>
        <div className="flex items-center space-x-2">
          <Checkbox
            id="required"
            checked={formData.required}
            onChange={(e) => setFormData(prev => ({ ...prev, required: e.target.checked }))}
          />
          <Label htmlFor="required">Required field</Label>
        </div>
      </div>

      {(formData.type === 'select' || formData.type === 'multiselect') && (
        <div>
          <Label>Options</Label>
          <div className="flex space-x-2 mb-2">
            <Input
              value={newOption}
              onChange={(e) => setNewOption(e.target.value)}
              placeholder="Add option"
              onKeyPress={(e) => e.key === 'Enter' && (e.preventDefault(), addOption())}
            />
            <Button type="button" onClick={addOption} variant="outline">
              Add
            </Button>
          </div>
          <div className="space-y-1">
            {formData.options.map((option, index) => (
              <div key={index} className="flex items-center justify-between bg-gray-50 dark:bg-gray-800 p-2 rounded">
                <span className="text-sm">{option}</span>
                <Button
                  type="button"
                  variant="ghost"
                  size="sm"
                  onClick={() => removeOption(index)}
                  className="text-red-600 hover:text-red-700"
                >
                  <X className="w-4 h-4" />
                </Button>
              </div>
            ))}
          </div>
          {errors.options && <p className="text-sm text-red-600 mt-1">{errors.options}</p>}
        </div>
      )}

      <div className="flex justify-end space-x-3">
        <Button type="button" variant="outline" onClick={onCancel}>
          Cancel
        </Button>
        <Button type="submit">
          <Save className="w-4 h-4 mr-2" />
          {field ? 'Update' : 'Create'} Field
        </Button>
      </div>
    </form>
  )
}
