import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { Antenna, Plus, Trash2 } from 'lucide-react'
import { api } from '../api/client'

export default function Listeners() {
  const qc = useQueryClient()
  const { data, isLoading } = useQuery({ queryKey: ['jobs'], queryFn: api.getJobs })

  const [form, setForm] = useState({ type: 'mtls', host: '0.0.0.0', port: '31337' })
  const [error, setError] = useState('')

  const startListener = useMutation({
    mutationFn: async () => {
      setError('')
      if (form.type === 'mtls') {
        return api.startMTLSListener(form.host, parseInt(form.port))
      }
      if (form.type === 'icmp') {
        return api.startICMPListener(form.host)
      }
      throw new Error(`Unsupported listener type: ${form.type}`)
    },
    onSuccess: () => qc.invalidateQueries({ queryKey: ['jobs'] }),
    onError: (e: Error) => setError(e.message),
  })

  const killJob = useMutation({
    mutationFn: (id: number) => api.killJob(id),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['jobs'] }),
  })

  return (
    <div className="p-8 space-y-8">
      <div className="flex items-center gap-3">
        <Antenna className="text-leonidas-400" size={20} />
        <h1 className="text-2xl font-bold text-leonidas-300">Listeners</h1>
      </div>

      {/* Start listener form */}
      <section className="bg-surface-800 border border-surface-600 rounded-lg p-6 space-y-4">
        <h2 className="text-sm font-semibold text-gray-400 uppercase tracking-wider">
          Start New Listener
        </h2>

        <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
          <div>
            <label className="block text-xs text-gray-500 mb-1">Type</label>
            <select
              value={form.type}
              onChange={(e) => setForm({ ...form, type: e.target.value })}
              className="w-full bg-surface-700 border border-surface-500 rounded px-3 py-2 text-sm text-gray-200 focus:outline-none focus:border-leonidas-400"
            >
              <option value="mtls">mTLS</option>
              <option value="wg">WireGuard</option>
              <option value="https">HTTPS</option>
              <option value="dns">DNS</option>
              <option value="icmp">ICMP (root required)</option>
            </select>
          </div>
          <div>
            <label className="block text-xs text-gray-500 mb-1">Bind Address</label>
            <input
              value={form.host}
              onChange={(e) => setForm({ ...form, host: e.target.value })}
              placeholder="0.0.0.0"
              className="w-full bg-surface-700 border border-surface-500 rounded px-3 py-2 text-sm text-gray-200 focus:outline-none focus:border-leonidas-400"
            />
          </div>
          {form.type !== 'icmp' && (
            <div>
              <label className="block text-xs text-gray-500 mb-1">Port</label>
              <input
                value={form.port}
                onChange={(e) => setForm({ ...form, port: e.target.value })}
                placeholder="31337"
                type="number"
                className="w-full bg-surface-700 border border-surface-500 rounded px-3 py-2 text-sm text-gray-200 focus:outline-none focus:border-leonidas-400"
              />
            </div>
          )}
        </div>

        {error && <p className="text-red-400 text-xs">{error}</p>}

        <button
          onClick={() => startListener.mutate()}
          disabled={startListener.isPending}
          className="flex items-center gap-2 px-4 py-2 bg-leonidas-500 hover:bg-leonidas-600 text-white rounded text-sm transition-colors disabled:opacity-50"
        >
          <Plus size={14} />
          {startListener.isPending ? 'Starting...' : 'Start Listener'}
        </button>
      </section>

      {/* Active jobs */}
      <section>
        <h2 className="text-sm font-semibold text-gray-400 uppercase tracking-wider mb-3">
          Active Jobs
        </h2>
        <div className="bg-surface-800 border border-surface-600 rounded-lg overflow-hidden">
          {isLoading ? (
            <div className="p-8 text-center text-gray-600">Loading jobs...</div>
          ) : (
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-surface-600 text-gray-500 text-xs uppercase tracking-wider">
                  <th className="text-left px-4 py-3">ID</th>
                  <th className="text-left px-4 py-3">Name</th>
                  <th className="text-left px-4 py-3">Description</th>
                  <th className="text-left px-4 py-3">Protocol</th>
                  <th className="text-left px-4 py-3">Port</th>
                  <th className="text-left px-4 py-3"></th>
                </tr>
              </thead>
              <tbody>
                {data?.Active?.map((j) => (
                  <tr key={j.ID} className="border-b border-surface-700 hover:bg-surface-700">
                    <td className="px-4 py-2.5 text-gray-600">{j.ID}</td>
                    <td className="px-4 py-2.5 text-leonidas-300 font-semibold uppercase">{j.Name}</td>
                    <td className="px-4 py-2.5 text-gray-400">{j.Description}</td>
                    <td className="px-4 py-2.5 text-gray-500">{j.Protocol}</td>
                    <td className="px-4 py-2.5 text-gray-500">{j.Port || '—'}</td>
                    <td className="px-4 py-2.5 text-right">
                      <button
                        onClick={() => killJob.mutate(j.ID)}
                        className="text-red-500 hover:text-red-400 transition-colors"
                        title="Kill job"
                      >
                        <Trash2 size={14} />
                      </button>
                    </td>
                  </tr>
                ))}
                {(!data?.Active || data.Active.length === 0) && (
                  <tr>
                    <td colSpan={6} className="px-4 py-12 text-center text-gray-600">
                      No active listener jobs
                    </td>
                  </tr>
                )}
              </tbody>
            </table>
          )}
        </div>
      </section>
    </div>
  )
}
