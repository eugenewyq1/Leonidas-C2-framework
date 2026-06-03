import { useState } from 'react'
import { Cpu } from 'lucide-react'

export default function Generate() {
  const [form, setForm] = useState({
    os: 'linux',
    arch: 'amd64',
    format: 'exe',
    transport: 'mtls',
    c2: '',
    debug: false,
  })

  const handleGenerate = () => {
    // In production, this calls StartGenerate RPC via the API client
    alert(`Generating ${form.format} implant for ${form.os}/${form.arch} via ${form.transport} → ${form.c2}`)
  }

  return (
    <div className="p-8 space-y-8">
      <div className="flex items-center gap-3">
        <Cpu className="text-leonidas-400" size={20} />
        <h1 className="text-2xl font-bold text-leonidas-300">Generate Implant</h1>
      </div>

      <div className="bg-surface-800 border border-surface-600 rounded-lg p-6 space-y-6 max-w-2xl">
        <div className="grid grid-cols-2 gap-4">
          <div>
            <label className="block text-xs text-gray-500 mb-1">Target OS</label>
            <select
              value={form.os}
              onChange={(e) => setForm({ ...form, os: e.target.value })}
              className="w-full bg-surface-700 border border-surface-500 rounded px-3 py-2 text-sm text-gray-200 focus:outline-none focus:border-leonidas-400"
            >
              <option value="linux">Linux</option>
              <option value="windows">Windows</option>
              <option value="darwin">macOS</option>
            </select>
          </div>
          <div>
            <label className="block text-xs text-gray-500 mb-1">Architecture</label>
            <select
              value={form.arch}
              onChange={(e) => setForm({ ...form, arch: e.target.value })}
              className="w-full bg-surface-700 border border-surface-500 rounded px-3 py-2 text-sm text-gray-200 focus:outline-none focus:border-leonidas-400"
            >
              <option value="amd64">amd64</option>
              <option value="arm64">arm64</option>
              <option value="386">386</option>
            </select>
          </div>
          <div>
            <label className="block text-xs text-gray-500 mb-1">Format</label>
            <select
              value={form.format}
              onChange={(e) => setForm({ ...form, format: e.target.value })}
              className="w-full bg-surface-700 border border-surface-500 rounded px-3 py-2 text-sm text-gray-200 focus:outline-none focus:border-leonidas-400"
            >
              <option value="exe">Executable</option>
              <option value="shared">Shared Library</option>
              <option value="shellcode">Shellcode</option>
            </select>
          </div>
          <div>
            <label className="block text-xs text-gray-500 mb-1">C2 Transport</label>
            <select
              value={form.transport}
              onChange={(e) => setForm({ ...form, transport: e.target.value })}
              className="w-full bg-surface-700 border border-surface-500 rounded px-3 py-2 text-sm text-gray-200 focus:outline-none focus:border-leonidas-400"
            >
              <option value="mtls">mTLS</option>
              <option value="https">HTTPS</option>
              <option value="dns">DNS</option>
              <option value="wg">WireGuard</option>
              <option value="icmp">ICMP</option>
            </select>
          </div>
        </div>

        <div>
          <label className="block text-xs text-gray-500 mb-1">C2 Address</label>
          <input
            value={form.c2}
            onChange={(e) => setForm({ ...form, c2: e.target.value })}
            placeholder="e.g. 10.0.0.1:31337"
            className="w-full bg-surface-700 border border-surface-500 rounded px-3 py-2 text-sm text-gray-200 focus:outline-none focus:border-leonidas-400"
          />
        </div>

        <label className="flex items-center gap-2 cursor-pointer">
          <input
            type="checkbox"
            checked={form.debug}
            onChange={(e) => setForm({ ...form, debug: e.target.checked })}
            className="accent-leonidas-400"
          />
          <span className="text-sm text-gray-400">Debug build (includes symbols)</span>
        </label>

        <button
          onClick={handleGenerate}
          className="px-6 py-2.5 bg-leonidas-500 hover:bg-leonidas-600 text-white rounded text-sm font-medium transition-colors"
        >
          Generate
        </button>
      </div>
    </div>
  )
}
