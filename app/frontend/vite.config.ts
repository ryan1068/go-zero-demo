import { defineConfig, loadEnv } from 'vite'
import vue from '@vitejs/plugin-vue'
import { resolve } from 'path'

// https://vite.dev/config/
export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd())

  return {
    plugins: [vue()],
    resolve: {
      alias: {
        '@': resolve(__dirname, 'src'),
      },
    },
    server: {
      port: 5173,
      proxy: {
        // 代理到 user-api (8000)
        '/api/user': {
          target: env.VITE_USER_API_URL || 'http://localhost:8000',
          changeOrigin: true,
          rewrite: (path) => path.replace(/^\/api\/user/, ''),
        },
        // 代理到 order-api (8001)
        '/api/order': {
          target: env.VITE_ORDER_API_URL || 'http://localhost:8001',
          changeOrigin: true,
          rewrite: (path) => path.replace(/^\/api\/order/, ''),
        },
      },
    },
  }
})
