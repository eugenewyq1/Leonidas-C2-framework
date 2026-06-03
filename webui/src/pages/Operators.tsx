import { useQuery } from '@tanstack/react-query'
import { Users, Circle } from 'lucide-react'
import { api } from '../api/client'

export default function Operators() {
  const { data, isLoading } = useQuery({ queryKey: ['operators'], queryFn: api.getOperators })

  return (
    <div className="p-8 space-y-6">
      <div className="flex items-center gap-3">
        <Users className="text-leonidas-400" size={20} />
        <h1 className="text-2xl font-bold text-leonidas-300">Operators</h1>
      </div>

      <div className="bg-surface-800 border border-surface-600 rounded-lg overflow-hidden">
        {isLoading ? (
          <div className="p-8 text-center text-gray-600">Loading operators...</div>
        ) : (
          <table className="w-full text-sm">
            <thead>
              <tr className="border-b border-surface-600 text-gray-500 text-xs uppercase tracking-wider">
                <th className="text-left px-4 py-3">Status</th>
                <th className="text-left px-4 py-3">Name</th>
              </tr>
            </thead>
            <tbody>
              {data?.Operators?.map((op) => (
                <tr key={op.Name} className="border-b border-surface-700 hover:bg-surface-700">
                  <td className="px-4 py-2.5">
                    <Circle
                      size={8}
                      fill={op.Online ? '#22c55e' : '#6b7280'}
                      className={op.Online ? 'text-green-500' : 'text-gray-500'}
                    />
                  </td>
                  <td className="px-4 py-2.5 text-gray-200">{op.Name}</td>
                </tr>
              ))}
              {(!data?.Operators || data.Operators.length === 0) && (
                <tr>
                  <td colSpan={2} className="px-4 py-12 text-center text-gray-600">
                    No operators registered
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
