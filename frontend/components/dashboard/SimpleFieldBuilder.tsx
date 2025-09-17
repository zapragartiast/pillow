'use client'

import React, { useState, useEffect } from 'react'
import { Plus, Edit, Trash2, Save, X } from 'lucide-react'

interface CustomField {
  id: string
  name: string
  label: string
  type: 'text' | 'textarea' | 'number' | 'email' | 'phone' | 'date' | 'boolean' | 'select' | 'multiselect'
  required: boolean
  value?: any
  options?: string[]
  order: number
}

interface CustomFieldsData {
  fields: CustomField[]
  metadata: {
    version: string
    last_updated: string
  }
}

interface SimpleFieldBuilderProps {
  userId: string
}

export default function SimpleFieldBuilder({ userId }: SimpleFieldBuilderProps) {
  const [fields, setFields] = useState<CustomField[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [showCreateForm, setShowCreateForm] = useState(false)
  const [editingField, setEditingField] = useState<CustomField | null>(null)

  useEffect(() => {
    loadFields()
  }, [userId])

  const loadFields = async () => {
    try {
      setIsLoading(true)
      const response = await fetch(`/api/users/${userId}/custom-fields`)

      if (!response.ok) {
        throw new Error('Failed to load custom fields')
      }

      const data: CustomFieldsData = await response.json()
      setFields(data.fields || [])
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load fields')
    } finally {
      setIsLoading(false)
    }
  }

  const createField = async (fieldData: Omit<CustomField, 'id' | 'order'>) => {
    try {
      const response = await fetch(`/api/users/${userId}/custom-fields`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(fieldData)
      })

      if (!response.ok) {
        const errorData = await response.json()
        throw new Error(errorData.message || 'Failed to create field')
      }

      await loadFields()
      setShowCreateForm(false)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to create field')
    }
  }

  const deleteField = async (fieldId: string) => {
    if (!confirm('Are you sure you want to delete this field?')) return

    try {
      const response = await fetch(`/api/users/${userId}/custom-fields/${fieldId}`, {
        method: 'DELETE'
      })

      if (!response.ok) {
        const errorData = await response.json()
        throw new Error(errorData.message || 'Failed to delete field')
      }

      await loadFields()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to delete field')
    }
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
        <button
          onClick={() => setShowCreateForm(true)}
          className="bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700 flex items-center space-x-2"
        >
          <Plus className="w-4 h-4" />
          <span>Add Field</span>
        </button>
      </div>

      {error && (
        <div className="bg-red-50 border border-red-200 text-red-800 px-4 py-3 rounded-lg">
          {error}
        </div>
      )}

      <div className="space-y-4">
        {fields.length === 0 ? (
          <div className="text-center py-12 bg-gray-50 dark:bg-gray-800 rounded-lg">
            <div className="text-gray-400 mb-4">
              <Plus className="w-12 h-12 mx-auto" />
            </div>
            <h3 className="text-lg font-medium text-gray-900 dark:text-white mb-2">No custom fields yet</h3>
            <p className="text-gray-600 dark:text-gray-400 mb-4">
              Create your first custom field to extend user profiles with additional information.
            </p>
            <button
              onClick={() => setShowCreateForm(true)}
              className="bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700"
            >
              Create Your First Field
            </button>
          </div>
        ) : (
          fields.map((field) => (
            <div
              key={field.id}
              className="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg p-4"
            >
              <div className="flex items-center justify-between">
                <div>
                  <h3 className="font-medium text-gray-900 dark:text-white">{field.label}</h3>
                  <p className="text-sm text-gray-600 dark:text-gray-400">
                    {field.type} • {field.name}
                    {field.required && <span className="text-red-500 ml-1">*</span>}
                  </p>
                </div>
                <div className="flex items-center space-x-2">
                  <button
                    onClick={() => setEditingField(field)}
                    className="text-blue-600 hover:text-blue-700 p-2"
                  >
                    <Edit className="w-4 h-4" />
                  </button>
                  <button
                    onClick={() => deleteField(field.id)}
                    className="text-red-600 hover:text-red-700 p-2"
                  >
                    <Trash2 className="w-4 h-4" />
                  </button>
                </div>
              </div>
            </div>
          ))
        )}
      </div>

      {showCreateForm && (
        <FieldForm
          onSave={createField}
          onCancel={() => setShowCreateForm(false)}
        />
      )}

      {editingField && (
        <FieldForm
          field={editingField}
          onSave={(fieldData) => {
            // Update field logic would go here
            console.log('Update field:', fieldData)
            setEditingField(null)
          }}
          onCancel={() => setEditingField(null)}
        />
      )}
    </div>
  )
}

interface FieldFormProps {
  field?: CustomField
  onSave: (fieldData: any) => void
  onCancel: () => void
}

function FieldForm({ field, onSave, onCancel }: FieldFormProps) {
  const [formData, setFormData] = useState({
    name: field?.name || '',
    label: field?.label || '',
    type: field?.type || 'text',
    required: field?.required || false,
    options: field?.options || []
  })

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    onSave(formData)
  }

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
      <div className="bg-white dark:bg-gray-800 rounded-lg p-6 w-full max-w-md">
        <h3 className="text-lg font-semibold mb-4">
          {field ? 'Edit Field' : 'Create New Field'}
        </h3>

        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="block text-sm font-medium mb-1">Field Name</label>
            <input
              type="text"
              value={formData.name}
              onChange={(e) => setFormData(prev => ({ ...prev, name: e.target.value }))}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500"
              placeholder="e.g., full_name"
              required
            />
          </div>

          <div>
            <label className="block text-sm font-medium mb-1">Field Label</label>
            <input
              type="text"
              value={formData.label}
              onChange={(e) => setFormData(prev => ({ ...prev, label: e.target.value }))}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500"
              placeholder="e.g., Full Name"
              required
            />
          </div>

          <div>
            <label className="block text-sm font-medium mb-1">Field Type</label>
            <select
              value={formData.type}
              onChange={(e) => setFormData(prev => ({ ...prev, type: e.target.value as any }))}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500"
            >
              <option value="text">Text</option>
              <option value="textarea">Textarea</option>
              <option value="number">Number</option>
              <option value="email">Email</option>
              <option value="phone">Phone</option>
              <option value="date">Date</option>
              <option value="boolean">Boolean</option>
              <option value="select">Select</option>
              <option value="multiselect">Multi-select</option>
            </select>
          </div>

          <div className="flex items-center space-x-2">
            <input
              type="checkbox"
              id="required"
              checked={formData.required}
              onChange={(e) => setFormData(prev => ({ ...prev, required: e.target.checked }))}
              className="rounded"
            />
            <label htmlFor="required" className="text-sm">Required field</label>
          </div>

          <div className="flex justify-end space-x-3">
            <button
              type="button"
              onClick={onCancel}
              className="px-4 py-2 text-gray-600 hover:text-gray-800"
            >
              Cancel
            </button>
            <button
              type="submit"
              className="bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700 flex items-center space-x-2"
            >
              <Save className="w-4 h-4" />
              <span>{field ? 'Update' : 'Create'}</span>
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}