import { Routes, Route, Navigate } from 'react-router-dom'
import Layout from './components/Layout'
import Dashboard from './pages/Dashboard'
import Sessions from './pages/Sessions'
import Beacons from './pages/Beacons'
import Listeners from './pages/Listeners'
import Generate from './pages/Generate'
import Loot from './pages/Loot'
import Operators from './pages/Operators'

export default function App() {
  return (
    <Routes>
      <Route path="/" element={<Layout />}>
        <Route index element={<Navigate to="/dashboard" replace />} />
        <Route path="dashboard" element={<Dashboard />} />
        <Route path="sessions" element={<Sessions />} />
        <Route path="beacons" element={<Beacons />} />
        <Route path="listeners" element={<Listeners />} />
        <Route path="generate" element={<Generate />} />
        <Route path="loot" element={<Loot />} />
        <Route path="operators" element={<Operators />} />
      </Route>
    </Routes>
  )
}
