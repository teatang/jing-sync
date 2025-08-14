import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";
import path from "path";
import fs from "fs";

// https://vite.dev/config/
export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
    },
    extensions: [".ts", ".tsx", ".js", ".jsx", ".json", ".vue"],
  },
  build: {
    target: "esnext",
    chunkSizeWarningLimit: 500, // 单位 KB
    cssCodeSplit: true,
    rollupOptions: {
      output: {
        manualChunks(id) {
          if (id.includes("node_modules")) {
            const lib = id.split("node_modules/")[1].split("/")[0];
            return lib.endsWith("-heavy") ? `lib-${lib}` : "vendors";
          }
          try {
            return fs.statSync(id).size > 400 * 1024 ? "large-modules" : null;
          } catch {
            return null;
          }
        },
        entryFileNames: "assets/[name]-[hash].js",
        assetFileNames: "assets/[name]-[hash].[ext]",
      },
    },
  },
  server: {
    port: 8080,
    host: true,
    open: true,
    proxy: {
      "/api": {
        target: "htttp://localhost:8888",
        changeOrigin: true, // 是否跨域
        // rewrite: (path) => path.replace(/^\/api/, '') // 重写路径，去掉 /api
      },
    },
  },
});
