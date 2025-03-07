import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import path from 'path';

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: [{ find: '@', replacement: path.resolve(__dirname, 'src') }],
  },
  server: {
    proxy: {
      '^/log/': {
        target: 'http://localhost:9000',
        changeOrigin: true,
        headers: {
          // cookie: env.VITE_COOKIES,
          // Authorization: `Bearer ${env.VITE_TOKEN}`
        },
      },
      '^/api/': {
        target: 'http://localhost:9000',
        changeOrigin: true,
        headers: {
          // cookie: env.VITE_COOKIES,
          // Authorization: `Bearer ${env.VITE_TOKEN}`
        },
      },
      '^/api/v1/misc/sockjs*': {
        target: 'ws://localhost:9000',
        changeOrigin: false,
        secure: false,
        ws: true,
      },
    }
  }
})
