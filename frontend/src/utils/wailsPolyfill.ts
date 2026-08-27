// frontend/src/utils/wailsPolyfill.ts (100行以下)
export function initWailsPolyfill() {
  if (typeof window === 'undefined') return;
  if (!window.runtime) {
    (window as any).runtime = {
      EventsOn: () => () => {}, EventsOnMultiple: () => () => {}, EventsOff: () => {}, EventsOffAll: () => {},
      EventsOnce: () => () => {}, EventsEmit: () => {},
      BrowserOpenURL: (url: string) => window.open(url, '_blank', 'noopener,noreferrer'),
      LogPrint: console.log, LogTrace: console.trace, LogDebug: console.debug, LogInfo: console.info, LogWarning: console.warn, LogError: console.error,
    };
  }

  if (!window.go || !window.go.app) {
    const mockApp = {
      GetSystemLanguage: async () => 'ja',
      GetAccounts: async (p: string) => { try { return await (await fetch(`/api/accounts?platform=${encodeURIComponent(p || 'twitter')}`)).json(); } catch { return []; } },
      GetTimeline: async (p: string, acc: string, f: string, l: number, o: number) => {
        try { return await (await fetch(`/api/timeline?platform=${encodeURIComponent(p || 'twitter')}&account_id=${encodeURIComponent(acc || 'all')}&filter=${encodeURIComponent(f || 'all')}&limit=${l || 50}&offset=${o || 0}`)).json(); } catch { return []; }
      },
      SearchArticles: async (q: string, acc: string, f: string, l: number, o: number) => {
        try { return await (await fetch(`/api/search?q=${encodeURIComponent(q || '')}&account_id=${encodeURIComponent(acc || 'all')}&filter=${encodeURIComponent(f || 'all')}&limit=${l || 50}&offset=${o || 0}`)).json(); } catch { return { items: [], total: 0 }; }
      },
      GetArticleDetail: async (p: string, id: string) => { try { return await (await fetch(`/api/article?platform=${encodeURIComponent(p || 'twitter')}&id=${encodeURIComponent(id)}`)).json(); } catch { return null; } },
      GetBroadcastStatus: async () => { try { return await (await fetch('/api/broadcast/status')).json(); } catch { return { enabled: true, use_tls: true, port: 5175, active_clients: 0, networks: [] }; } },
      GetConfig: async () => ({
        system: { language: 'ja', default_framework: 'twitter', env: 'production' },
        storage: { db_path: './archive.db', local_media_dir: 'G:/Media_Storage/Influencers', stash_dir: './stash', dumps_dir: './dumps', stash_enabled: true },
        network: { stash_port: 9999, frontend_port: 5173, internal_bind_address: '127.0.0.1', middleware_port: 5175, public_bind_address: '0.0.0.0' },
        scheduler: { poll_interval_sec: 300, backup_interval_hours: 24, max_backup_generations: 7 },
        broadcast: { enabled: true, allowed_networks: ['192.168.10.0/24', '192.168.3.0/24', '127.0.0.1/32'] },
        appearance: { font_family_ja: 'Hiragino Sans, Meiryo, sans-serif', font_family_en: 'Nunito, sans-serif', font_family_zh: 'Microsoft YaHei, SimHei, sans-serif' },
      }),
      GetMediaList: async () => ({ items: [], total: 0, stats: { total_count: 0, image_count: 0, video_count: 0 } }),
      GetMediaDownloadStatusStats: async () => ({ queued: 0, completed: 1103, dead_404: 0, outsourced: 933, retained: 818, failed: 0, total: 2854 }),
      StartMediaEscalateJob: async () => ({ id: 'esc_job', status: 'RUNNING', percentage: 0 }),
      StartSmartRecoveryJob: async () => ({ id: 'smart_job', status: 'RUNNING', percentage: 0 }),
      StartThunderEscalateJob: async () => ({ id: 'thun_job', status: 'RUNNING', percentage: 0 }),
      GetSkinPackage: async (p: string) => {
        try {
          const [cssRes, yamlRes] = await Promise.all([fetch(`/plugins/${p}/skin/design.css`), fetch(`/plugins/${p}/skin/layout.yaml`)]);
          return { platform: p, design_css: cssRes.ok ? await cssRes.text() : '', layout_yaml: yamlRes.ok ? await yamlRes.text() : '', controller: '' };
        } catch { return { platform: p, design_css: '', layout_yaml: '', controller: '' }; }
      },
      GetSkinCSS: async (p: string) => { try { const r = await fetch(`/plugins/${p}/skin/design.css`); return r.ok ? await r.text() : ''; } catch { return ''; } },
      UpdateAccount: async (...args: any[]) => { console.log('[Polyfill] UpdateAccount:', args); },
      MergeAccounts: async (...args: any[]) => { console.log('[Polyfill] MergeAccounts:', args); return null; },
      SaveAvatarImage: async (p: string, k: string) => `/avatars/${p || 'twitter'}/${k}.jpg`,
      ListAvailableAvatars: async (p: string) => [`/avatars/${p || 'twitter'}/msluo14_avatar_001.jpg`, `/avatars/${p || 'twitter'}/default_avatar.jpg`],
    };
    (window as any).go = { app: { App: mockApp }, main: { App: mockApp } };
  }
}
