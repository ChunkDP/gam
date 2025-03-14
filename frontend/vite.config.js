import { defineConfig } from 'vite';
import vue from '@vitejs/plugin-vue';
import path from 'path';
import tailwindcss from '@tailwindcss/vite'
export default defineConfig({
  plugins: [vue(),tailwindcss()],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, 'src'),
      '~': path.resolve(__dirname, './')
    },
  },
  server: {
    port: 3000,
    open: true,
    proxy: {
      
      '/ws/notifications': {
        target: 'http://localhost:8080',
        changeOrigin: true,
        ws: true  // 启用 websocket 代理
        //rewrite: (path) => path.replace(/^\/api/, '')
      },
      '/gam': {
        target: 'http://localhost:8080',
        changeOrigin: true,
        ws: true  // 启用 websocket 代理
        //rewrite: (path) => path.replace(/^\/api/, '')
      },
      '/uploads': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      }
    }
  },
  build: {
    outDir: 'dist',
    assetsDir: 'assets',
    rollupOptions: {
      input: 'src/main.js',
    },
  },
  optimizeDeps: {
    include: ['element-plus'],
  },
  css: {
    preprocessorOptions: {
      scss: {
        additionalData: `@use "@/assets/style/variables.scss";`
      }
    }
  },
});