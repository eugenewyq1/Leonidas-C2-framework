import { type LucideIcon } from 'lucide-react'

interface StatCardProps {
  label: string
  value: number | string
  Icon: LucideIcon
  accent?: boolean
}

export default function StatCard({ label, value, Icon, accent }: StatCardProps) {
  return (
    <div className="bg-surface-800 border border-surface-600 rounded-lg p-4 flex items-center gap-4">
      <div className={`p-2 rounded-md ${accent ? 'bg-leonidas-500/20 text-leonidas-400' : 'bg-surface-700 text-gray-400'}`}>
        <Icon size={20} />
      </div>
      <div>
        <p className="text-2xl font-bold text-gray-100">{value}</p>
        <p className="text-xs text-gray-500 uppercase tracking-wider">{label}</p>
      </div>
    </div>
  )
}
