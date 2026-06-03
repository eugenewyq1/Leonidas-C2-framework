import { useQuery } from '@tanstack/react-query'
import { Radio } from 'lucide-react'
import { api } from '../api/client'

function formatInterval(ms: number) {
  const s = Math.floor(ms / 1000)
  if (s < 60) return `${s}s`
  return `${Math.floor(s / 60)}m ${s % 60}s`
}

export default function Beacons() {
  const { data, isLoading } = useQuery({ queryKey: ['beacons'], queryFn: api.getBeacons })

  return (
    <div className="p-8 space-y-6">
      <div className="flex items-center gap-3">
        <Radio className="text-leonidas-400" size={20} />
        <h1 className="text-2xl font-bold text-leonidas-300">Beacons</h1>
        <span className="ml-auto text-xs text-gray-500">
          {data?.Beacons?.length ?? 0} registered
        </span>
      </div>

      <div className="bg-surface-800 border border-surface-600 rounded-lg overflow-hidden">
        {isLoading ? (
          <div className="p-8 text-center text-gray-600">Loading beacons...</div>
        ) : (
          <table className="w-full text-sm">
            <thead>
              <tr className="border-b border-surface-600 text-gray-500 text-xs uppercase tracking-wider">
                <th className="text-left px-4 py-3">Name</th>
                <th className="text-left px-4 py-3">Address</th>
                <th className="text-left px-4 py-3">OS</th>
                <th className="text-left px-4 py-3">Interval</th>
                <th className="text-left px-4 py-3">Jitter</th>
                <th className="text-left px-4 py-3">Next Check-in</th>
                <th className="text-left px-4 py-3">Transport</th>
              </tr>
            </thead>
            <tbody>
              {data?.Beacons?.map((b) => (
                <tr key={b.ID} className="border-b border-surface-700 hover:bg-surface-700 transition-colors">
                  <td className="px-4 py-2.5 text-leonidas-300 font-semibold">{b.Name}</td>
                  <td className="px-4 py-2.5 text-gray-400 font-mono">{b.RemoteAddress}</td>
                  <td className="px-4 py-2.5 text-gray-400">{b.OS}</td>
                  <td className="px-4 py-2.5 text-gray-400">{formatInterval(b.Interval)}</td>
                  <td className="px-4 py-2.5 text-gray-400">{formatInterval(b.Jitter)}</td>
                  <td className="px-4 py-2.5 text-gray-500 text-xs">{b.NextCheckin}</td>
                  <td className="px-4 py-2.5">
                    <span className="px-2 py-0.5 rounded text-xs bg-leonidas-500/20 text-leonidas-300">
                      {b.ActiveC2}
                    </span>
                  </td>
                </tr>
              ))}
              {(!data?.Beacons || data.Beacons.length === 0) && (
                <tr>
                  <td colSpan={7} className="px-4 py-12 text-center text-gray-600">
                    No beacons registered
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
