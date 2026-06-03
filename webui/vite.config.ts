import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
  build: {
    outDir: '../server/webui/dist',
    emptyOutDir: true,
  },
  server: {
    proxy: {
      '/rpcpb.LeonidasRPC': {
        target: 'https://localhost:8443',
        changeOrigin: true,
        secure: false,
      },
    },
  },
})
