import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import path from 'path'

// https://vite.dev/config/
export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
    },
    extensions: ['.ts', '.tsx', '.js', '.jsx', '.json', '.vue']
  },
  server: {
    port: 8080,
    host: true,
    open: true,
    proxy: {
      '/api': {
        target: 'htttp://localhost:8888',
        changeOrigin: true, // 是否跨域
        // rewrite: (path) => path.replace(/^\/api/, '') // 重写路径，去掉 /api
      } 
    } 
  }
})
