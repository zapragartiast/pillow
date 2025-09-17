import React from 'react'
import DataExplorerClient from '@/components/dashboard/DataExplorer/DataExplorerClient'

export const metadata = {
  title: 'Table Editor',
}

export default function TableEditorPage() {
  return (
    <div className="max-w-7xl mx-auto">
      <div className="mb-6">
        <h1 className="text-2xl font-semibold text-gray-900 dark:text-white">Data Explorer</h1>
        {/* <p className="text-sm text-gray-500 dark:text-gray-400">
          Browse, filter, sort, resize columns, and edit inline with server-side pagination.
        </p> */}
      </div>

      <DataExplorerClient className="mt-2" />
    </div>
  )
}