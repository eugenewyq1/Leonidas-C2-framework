import { Outlet, NavLink } from 'react-router-dom'
import {
  LayoutDashboard,
  Monitor,
  Radio,
  Antenna,
  Cpu,
  Archive,
  Users,
  Shield,
} from 'lucide-react'
import { clsx } from 'clsx'

const navItems = [
  { to: '/dashboard',  label: 'Dashboard',  Icon: LayoutDashboard },
  { to: '/sessions',   label: 'Sessions',   Icon: Monitor },
  { to: '/beacons',    label: 'Beacons',    Icon: Radio },
  { to: '/listeners',  label: 'Listeners',  Icon: Antenna },
  { to: '/generate',   label: 'Generate',   Icon: Cpu },
  { to: '/loot',       label: 'Loot',       Icon: Archive },
  { to: '/operators',  label: 'Operators',  Icon: Users },
]

export default function Layout() {
  return (
    <div className="flex h-screen overflow-hidden">
      {/* Sidebar */}
      <aside className="w-56 shrink-0 bg-surface-800 border-r border-surface-600 flex flex-col">
        {/* Logo */}
        <div className="flex items-center gap-3 px-4 py-5 border-b border-surface-600">
          <Shield className="text-leonidas-400" size={22} />
          <span className="text-leonidas-300 font-semibold tracking-widest text-sm uppercase">
            Leonidas
          </span>
        </div>

        {/* Nav */}
        <nav className="flex-1 py-3 space-y-0.5 overflow-y-auto">
          {navItems.map(({ to, label, Icon }) => (
            <NavLink
              key={to}
              to={to}
              className={({ isActive }) =>
                clsx(
                  'flex items-center gap-3 px-4 py-2.5 text-sm transition-colors',
                  isActive
                    ? 'bg-leonidas-500/20 text-leonidas-300 border-l-2 border-leonidas-400'
                    : 'text-gray-400 hover:bg-surface-700 hover:text-gray-200'
                )
              }
            >
              <Icon size={16} />
              {label}
            </NavLink>
          ))}
        </nav>

        {/* Footer */}
        <div className="px-4 py-3 border-t border-surface-600 text-xs text-gray-600">
          Leonidas C2 Framework
        </div>
      </aside>

      {/* Main content */}
      <main className="flex-1 overflow-y-auto bg-surface-900">
        <Outlet />
      </main>
    </div>
  )
}
