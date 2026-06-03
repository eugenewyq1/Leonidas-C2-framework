import { useQuery } from '@tanstack/react-query'
import { Monitor, Radio, Antenna, Tag } from 'lucide-react'
import { api } from '../api/client'
import StatCard from '../components/StatCard'

export default function Dashboard() {
  const sessions  = useQuery({ queryKey: ['sessions'],  queryFn: api.getSessions })
  const beacons   = useQuery({ queryKey: ['beacons'],   queryFn: api.getBeacons })
  const jobs      = useQuery({ queryKey: ['jobs'],       queryFn: api.getJobs })
  const version   = useQuery({ queryKey: ['version'],   queryFn: api.getVersion })

  return (
    <div className="p-8 space-y-8">
      <div>
        <h1 className="text-2xl font-bold text-leonidas-300 tracking-wide">Dashboard</h1>
        {version.data && (
          <p className="text-xs text-gray-600 mt-1">
            Leonidas C2 v{version.data.Major}.{version.data.Minor}.{version.data.Patch}
            {version.data.Commit ? ` (${version.data.Commit.slice(0, 7)})` : ''}
          </p>
        )}
      </div>

      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        <StatCard
          label="Active Sessions"
          value={sessions.data?.Sessions?.length ?? '—'}
          Icon={Monitor}
          accent
        />
        <StatCard
          label="Beacons"
          value={beacons.data?.Beacons?.length ?? '—'}
          Icon={Radio}
          accent
        />
        <StatCard
          label="Listener Jobs"
          value={jobs.data?.Active?.length ?? '—'}
          Icon={Antenna}
        />
        <StatCard
          label="Framework"
          value="Leonidas"
          Icon={Tag}
        />
      </div>

      {/* Recent Sessions */}
      <section>
        <h2 className="text-sm font-semibold text-gray-400 uppercase tracking-wider mb-3">
          Recent Sessions
        </h2>
        <div className="bg-surface-800 border border-surface-600 rounded-lg overflow-hidden">
          <table className="w-full text-sm">
            <thead>
              <tr className="border-b border-surface-600 text-gray-500 text-xs uppercase tracking-wider">
                <th className="text-left px-4 py-3">Name</th>
                <th className="text-left px-4 py-3">Address</th>
                <th className="text-left px-4 py-3">OS</th>
                <th className="text-left px-4 py-3">Transport</th>
              </tr>
            </thead>
            <tbody>
              {sessions.data?.Sessions?.slice(0, 5).map((s) => (
                <tr key={s.ID} className="border-b border-surface-700 hover:bg-surface-700 transition-colors">
                  <td className="px-4 py-2.5 text-leonidas-300 font-medium">{s.Name}</td>
                  <td className="px-4 py-2.5 text-gray-400">{s.RemoteAddress}</td>
                  <td className="px-4 py-2.5 text-gray-400">{s.OS}/{s.Arch}</td>
                  <td className="px-4 py-2.5">
                    <span className="px-2 py-0.5 rounded text-xs bg-surface-600 text-gray-300">
                      {s.ActiveC2}
                    </span>
                  </td>
                </tr>
              ))}
              {(!sessions.data?.Sessions || sessions.data.Sessions.length === 0) && (
                <tr>
                  <td colSpan={4} className="px-4 py-8 text-center text-gray-600">
                    No active sessions
                  </td>
                </tr>
              )}
            </tbody>
          </table>
        </div>
      </section>
    </div>
  )
}
