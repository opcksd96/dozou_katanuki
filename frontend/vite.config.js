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

            // セキュリティ: パストラバーサル防止
            if (!targetPath.startsWith(baseDir)) {
              res.statusCode = 403;
              res.end('Forbidden');
              return;
            }

            // 画像拡張子のみ許可
            const ext = path.extname(targetPath).toLowerCase();
            const allowedExts = ['.jpg', '.jpeg', '.png', '.webp', '.svg', '.gif'];
            if (!allowedExts.includes(ext)) {
              res.statusCode = 404;
              res.end('Not Found');
              return;
            }

            // 実体ファイルが存在する場合はサーブ
            if (fs.existsSync(targetPath) && !fs.statSync(targetPath).isDirectory()) {
              const contentType = ext === '.png' ? 'image/png' : ext === '.svg' ? 'image/svg+xml' : ext === '.webp' ? 'image/webp' : ext === '.gif' ? 'image/gif' : 'image/jpeg';
              res.setHeader('Content-Type', contentType);
              fs.createReadStream(targetPath).pipe(res);
              return;
            }

            // 明示的なデフォルトアバター要求時のみデフォルトSVGを返却
            if (cleanPath.includes('default_avatar') || cleanPath.includes('default-avatar')) {
              res.setHeader('Content-Type', 'image/svg+xml');
              res.statusCode = 200;
              res.end(`<svg xmlns="http://www.w3.org/2000/svg" width="64" height="64"><rect width="100%" height="100%" fill="#334155" rx="32"/><circle cx="32" cy="24" r="12" fill="#94a3b8"/><path d="M16 54c0-8.837 7.163-16 16-16s16 7.163 16 16" fill="#94a3b8"/></svg>`);
              return;
            }

            // 個別ファイルが存在しない場合は隠蔽せず 404 を返却
            res.statusCode = 404;
            res.end('Avatar Not Found');
            return;
          }
          if (req.url && req.url.startsWith('/media/')) {
            const mID = path.basename(req.url.split('?')[0]);
            const bases = ['G:/Media_Storage/Influencers', path.resolve(__dirname, '..', 'blobs'), path.resolve(__dirname, '..', 'stash')];
            const exts = ['', '.jpg', '.jpeg', '.png', '.webp', '.mp4', '.gif'];
            for (const base of bases) {
              if (!fs.existsSync(base)) continue;
              for (const ext of exts) {
                const direct = path.join(base, mID + ext);
                if (fs.existsSync(direct) && !fs.statSync(direct).isDirectory()) {
                  const targetExt = path.extname(direct).toLowerCase();
                  res.setHeader('Content-Type', targetExt === '.png' ? 'image/png' : targetExt === '.webp' ? 'image/webp' : targetExt === '.mp4' ? 'video/mp4' : targetExt === '.gif' ? 'image/gif' : 'image/jpeg');
                  fs.createReadStream(direct).pipe(res);
                  return;
                }
              }
              try {
                const subdirs = fs.readdirSync(base);
                const subAssetDirs = ['X(Twitter)/_assets', 'Twitter/_assets', 'Base/_assets', ''];
                for (const sub of subdirs) {
                  for (const subAsset of subAssetDirs) {
                    for (const ext of exts) {
                      const cand = path.join(base, sub, subAsset, mID + ext);
                      if (fs.existsSync(cand) && !fs.statSync(cand).isDirectory()) {
                        const targetExt = path.extname(cand).toLowerCase();
                        res.setHeader('Content-Type', targetExt === '.png' ? 'image/png' : targetExt === '.webp' ? 'image/webp' : targetExt === '.mp4' ? 'video/mp4' : targetExt === '.gif' ? 'image/gif' : 'image/jpeg');
                        fs.createReadStream(cand).pipe(res);
                        return;
                      }
                    }
                  }
                }
              } catch (_) {}
            }
            res.setHeader('Content-Type', 'image/svg+xml');
            res.statusCode = 200;
            res.end(`<svg xmlns="http://www.w3.org/2000/svg" width="400" height="250"><rect width="100%" height="100%" fill="#1e293b" rx="8"/><text x="50%" y="50%" dominant-baseline="middle" text-anchor="middle" fill="#64748b" font-family="sans-serif">Attached Media Preview</text></svg>`);
            return;
          }
          next();
        });
      },
    },
  ],
  server: {
    host: '0.0.0.0',
    port: 5173,
    proxy: {
      '/stash-proxy': {
        target: 'http://127.0.0.1:9999',
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/stash-proxy/, ''),
      },
      '/media': {
        target: 'http://127.0.0.1:5175',
        changeOrigin: true,
      },
    },
  },
});
