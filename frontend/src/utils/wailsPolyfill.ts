// frontend/src/utils/wailsPolyfill.ts (100行以下 - SPEC-PRINCIPLE-001)
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
    const postJson = async (url: string, data: any) => (await fetch(url, { method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify(data) })).json();
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
      GetBroadcastStatus: async () => { try { return await (await fetch('/api/broadcast/status')).json(); } catch { return { enabled: true, use_tls: false, port: 5175, active_clients: 0, networks: [] }; } },
      GetConfig: async () => { try { const r = await fetch('/api/config'); if (r.ok) return await r.json(); } catch {} return { system: { language: 'ja', default_framework: 'twitter', env: 'production' } }; },
      GetMediaList: async (acc = 'all', st = 'all', t = 'all', l = 24, o = 0) => {
        try { return await (await fetch(`/api/media?account_id=${encodeURIComponent(acc)}&status=${encodeURIComponent(st)}&type=${encodeURIComponent(t)}&limit=${l}&offset=${o}`)).json(); } catch { return { items: [], total: 0, stats: { total_count: 0, image_count: 0, video_count: 0 } }; }
      },
      GetMediaDownloadStatusStats: async (acc = 'all') => {
        try { return await (await fetch(`/api/media/stats?account_id=${encodeURIComponent(acc)}`)).json(); } catch { return { queued: 0, completed: 0, dead_404: 0, outsourced: 0, retained: 0, failed: 0, total: 0 }; }
      },
      ToggleMediaBookmark: async (id: string) => { try { const d = await postJson('/api/media/bookmark', { media_id: id }); return !!d?.is_bookmarked; } catch { return false; } },
      UpdateMediaMetadata: async (id: string, st: string, sId: string, iId: string, r: string) => { try { await postJson('/api/media/update', { media_id: id, download_status: st, stash_scene_id: sId, stash_image_id: iId, failed_reason: r }); return true; } catch { return false; } },
      PurgeMedia: async (id: string) => { try { await postJson('/api/media/purge', { media_id: id }); return true; } catch { return false; } },
      PurgeMediaByStatus: async (st: string, acc = 'all') => { try { const d = await postJson('/api/media/purge-status', { status: st, account_id: acc }); return d?.purged_count || 0; } catch { return 0; } },
      RequeueMediaByStatus: async (st: string, acc = 'all') => { try { const d = await postJson('/api/media/requeue', { status: st, account_id: acc }); return d?.requeued_count || 0; } catch { return 0; } },
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
    (window as any)._isWailsPolyfill = true;
    (window as any).go = { app: { App: mockApp }, main: { App: mockApp } };
  }
}
