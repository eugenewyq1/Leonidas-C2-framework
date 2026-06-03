import { useQuery } from '@tanstack/react-query'
import { Monitor } from 'lucide-react'
import { api } from '../api/client'

export default function Sessions() {
  const { data, isLoading } = useQuery({ queryKey: ['sessions'], queryFn: api.getSessions })

  return (
    <div className="p-8 space-y-6">
      <div className="flex items-center gap-3">
        <Monitor className="text-leonidas-400" size={20} />
        <h1 className="text-2xl font-bold text-leonidas-300">Sessions</h1>
        <span className="ml-auto text-xs text-gray-500">
          {data?.Sessions?.length ?? 0} active
        </span>
      </div>

      <div className="bg-surface-800 border border-surface-600 rounded-lg overflow-hidden">
        {isLoading ? (
          <div className="p-8 text-center text-gray-600">Loading sessions...</div>
        ) : (
          <table className="w-full text-sm">
            <thead>
              <tr className="border-b border-surface-600 text-gray-500 text-xs uppercase tracking-wider">
                <th className="text-left px-4 py-3">ID</th>
                <th className="text-left px-4 py-3">Name</th>
                <th className="text-left px-4 py-3">Hostname</th>
                <th className="text-left px-4 py-3">Address</th>
                <th className="text-left px-4 py-3">OS / Arch</th>
                <th className="text-left px-4 py-3">PID</th>
                <th className="text-left px-4 py-3">Transport</th>
                <th className="text-left px-4 py-3">Last Seen</th>
              </tr>
            </thead>
            <tbody>
              {data?.Sessions?.map((s) => (
                <tr key={s.ID} className="border-b border-surface-700 hover:bg-surface-700 transition-colors cursor-pointer">
                  <td className="px-4 py-2.5 text-gray-600 text-xs font-mono">{s.ID.slice(0, 8)}</td>
                  <td className="px-4 py-2.5 text-leonidas-300 font-semibold">{s.Name}</td>
                  <td className="px-4 py-2.5 text-gray-300">{s.Hostname}</td>
                  <td className="px-4 py-2.5 text-gray-400 font-mono">{s.RemoteAddress}</td>
                  <td className="px-4 py-2.5 text-gray-400">{s.OS}/{s.Arch}</td>
                  <td className="px-4 py-2.5 text-gray-500">{s.PID}</td>
                  <td className="px-4 py-2.5">
                    <span className="px-2 py-0.5 rounded text-xs bg-leonidas-500/20 text-leonidas-300">
                      {s.ActiveC2}
                    </span>
                  </td>
                  <td className="px-4 py-2.5 text-gray-500 text-xs">{s.LastCheckin}</td>
                </tr>
              ))}
              {(!data?.Sessions || data.Sessions.length === 0) && (
                <tr>
                  <td colSpan={8} className="px-4 py-12 text-center text-gray-600">
                    No active sessions. Generate an implant and deploy it to get started.
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
