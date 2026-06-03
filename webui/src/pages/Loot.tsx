import { useQuery } from '@tanstack/react-query'
import { Archive } from 'lucide-react'
import { api } from '../api/client'

export default function Loot() {
  const { data, isLoading } = useQuery({ queryKey: ['loot'], queryFn: api.getLootFiles })

  return (
    <div className="p-8 space-y-6">
      <div className="flex items-center gap-3">
        <Archive className="text-leonidas-400" size={20} />
        <h1 className="text-2xl font-bold text-leonidas-300">Loot</h1>
      </div>

      <div className="bg-surface-800 border border-surface-600 rounded-lg overflow-hidden">
        {isLoading ? (
          <div className="p-8 text-center text-gray-600">Loading loot...</div>
        ) : (
          <table className="w-full text-sm">
            <thead>
              <tr className="border-b border-surface-600 text-gray-500 text-xs uppercase tracking-wider">
                <th className="text-left px-4 py-3">Name</th>
                <th className="text-left px-4 py-3">Type</th>
                <th className="text-left px-4 py-3">Created</th>
              </tr>
            </thead>
            <tbody>
              {data?.LootFiles?.map((item) => (
                <tr key={item.ID} className="border-b border-surface-700 hover:bg-surface-700">
                  <td className="px-4 py-2.5 text-leonidas-300">{item.Name}</td>
                  <td className="px-4 py-2.5 text-gray-400">{item.Type}</td>
                  <td className="px-4 py-2.5 text-gray-500 text-xs">{item.CreatedAt}</td>
                </tr>
              ))}
              {(!data?.LootFiles || data.LootFiles.length === 0) && (
                <tr>
                  <td colSpan={3} className="px-4 py-12 text-center text-gray-600">
                    No loot collected yet
                  </td>
                </tr>
              )}
            </tbody>
          </table>
        )}
      </div>
    </div>
  )
}
