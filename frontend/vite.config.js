// frontend/vite.config.js (100行以下)
import { defineConfig } from 'vite';
import vue from '@vitejs/plugin-vue';
import tailwindcss from '@tailwindcss/vite';
import path from 'path';
import fs from 'fs';

export default defineConfig({
  plugins: [
    vue(),
    tailwindcss(),
    {
      name: 'avatar-pure-provider',
      configureServer(server) {
        server.middlewares.use((req, res, next) => {
          if (req.url && (req.url.startsWith('/avatars/') || req.url.startsWith('/assets/'))) {
            const cleanPath = req.url.split('?')[0].replace(/^\/(avatars|assets)\//, '');
            // ルート直下の assets/ を解決
            const baseDir = path.resolve(__dirname, '..', 'assets');
            const targetPath = path.resolve(baseDir, cleanPath);

            if (fs.existsSync(targetPath) && !fs.statSync(targetPath).isDirectory()) {
              res.setHeader('Content-Type', 'image/jpeg');
              fs.createReadStream(targetPath).pipe(res);
              return;
            }
          }
          next();
        });
      },
    },
  ],
  server: {
    proxy: {
      '/stash-proxy': {
        target: 'http://127.0.0.1:9999',
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/stash-proxy/, ''),
      },
    },
  },
});
