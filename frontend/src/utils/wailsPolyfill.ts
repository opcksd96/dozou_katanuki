// frontend/src/utils/wailsPolyfill.ts (100行以下 - SPEC-PRINCIPLE-001)

export function initWailsPolyfill() {
  if (typeof window === 'undefined') return;

  const isDev = window.location.port === '5173';
  const baseUrl = isDev ? 'http://127.0.0.1:5175' : '';

  let sseListeners: Record<string, Function[]> = {};

  if (!window.runtime) (window as any).runtime = {};
  
  window.runtime.EventsOn = (eventName: string, callback: Function) => {
    if (!sseListeners[eventName]) sseListeners[eventName] = [];
    sseListeners[eventName].push(callback);
    return () => { sseListeners[eventName] = sseListeners[eventName].filter(cb => cb !== callback); };
  };
  window.runtime.EventsOnMultiple = () => () => {};
  window.runtime.EventsOff = () => {};
  window.runtime.EventsOffAll = () => {};
  window.runtime.EventsOnce = () => () => {};
  window.runtime.EventsEmit = () => {};
  window.runtime.BrowserOpenURL = (url: string) => window.open(url, '_blank', 'noopener,noreferrer');
  window.runtime.LogPrint = console.log; window.runtime.LogTrace = console.trace; window.runtime.LogDebug = console.debug;
  window.runtime.LogInfo = console.info; window.runtime.LogWarning = console.warn; window.runtime.LogError = console.error;

  const sseUrl = `${baseUrl}/api/events`;
  const evtSource = new EventSource(sseUrl);
  evtSource.onmessage = (event) => {
    try {
      const payload = JSON.parse(event.data);
      if (payload && payload.name && sseListeners[payload.name]) {
         sseListeners[payload.name].forEach(cb => cb(...(payload.data || [])));
      }
    } catch(e) {}
  };

  const postJson = async (url: string, data: any) => (await fetch(baseUrl + url, { method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify(data) })).json();
  const getJson = async (url: string, fallback: any) => { try { return await (await fetch(baseUrl + url)).json(); } catch { return fallback; } };

  const mockApp = {
    GetSystemLanguage: async () => 'ja',
    GetAccounts: async (p: string) => getJson(`/api/accounts?platform=${encodeURIComponent(p || 'twitter')}`, []),
    GetTimeline: async (p: string, acc: string, f: string, l: number, o: number) => getJson(`/api/timeline?platform=${encodeURIComponent(p || 'twitter')}&account_id=${encodeURIComponent(acc || 'all')}&filter=${encodeURIComponent(f || 'all')}&limit=${l || 50}&offset=${o || 0}`, []),
    SearchArticles: async (q: string, acc: string, f: string, l: number, o: number) => getJson(`/api/search?q=${encodeURIComponent(q || '')}&account_id=${encodeURIComponent(acc || 'all')}&filter=${encodeURIComponent(f || 'all')}&limit=${l || 50}&offset=${o || 0}`, { items: [], total: 0 }),
    TrashArticle: async (id: string, trashedBy = 'admin', reason = '') => { try { await postJson('/api/article/trash', { id, trashed_by: trashedBy, reason }); return true; } catch { return false; } },
    RestoreArticle: async (id: string) => { try { await postJson('/api/article/restore', { id }); return true; } catch { return false; } },
    BatchTrashArticles: async (ids: string[], trashedBy = 'admin', reason = '') => { try { await postJson('/api/article/batch-trash', { ids, trashed_by: trashedBy, reason }); return true; } catch { return false; } },
    BatchRestoreArticles: async (ids: string[]) => { try { await postJson('/api/article/batch-restore', { ids }); return true; } catch { return false; } },
    BatchResetTranslations: async (ids: string[]) => { try { await postJson('/api/article/batch-reset-translations', { ids }); return true; } catch { return false; } },
    GetArticlesByIDs: async (ids: string[]) => { try { return (await Promise.all(ids.map(id => fetch(`${baseUrl}/api/article?id=${id}`).then(r => r.json())))).filter(Boolean); } catch { return []; } },
    GetArticleDetail: async (p: string, id: string) => getJson(`/api/article?platform=${encodeURIComponent(p || 'twitter')}&id=${encodeURIComponent(id)}`, null),
    GetBroadcastStatus: async () => getJson('/api/broadcast/status', { enabled: true, use_tls: false, port: 5175, active_clients: 0, networks: [] }),
    GetConfig: async () => getJson('/api/config', { system: { language: 'ja', default_framework: 'twitter', env: 'production' } }),
    GetPipelineOverview: async () => getJson('/api/admin/pipeline/overview', null),
    GetPipelineLogs: async (stage: string, limit: number) => getJson(`/api/admin/pipeline/logs?stage=${encodeURIComponent(stage || 'all')}&limit=${limit || 50}`, []),
    IsPipelineAutoEngineRunning: async () => { try { const o = await getJson('/api/admin/pipeline/overview', null); return o?.auto_engine_running || false; } catch { return false; } },
    TogglePipelineAutoEngine: async (enable: boolean) => { try { const r = await postJson('/api/admin/pipeline/toggle', { enable }); return !!r?.isRunning; } catch { return false; } },
    SyncThunderDownloads: async () => { try { await postJson('/api/admin/pipeline/sync-thunder', {}); } catch {} },
    GetMediaList: async (acc = 'all', st = 'all', t = 'all', l = 24, o = 0) => getJson(`/api/media?account_id=${encodeURIComponent(acc)}&status=${encodeURIComponent(st)}&type=${encodeURIComponent(t)}&limit=${l}&offset=${o}`, { items: [], total: 0, stats: { total_count: 0, image_count: 0, video_count: 0 } }),
    GetMediaDownloadStatusStats: async (acc = 'all') => getJson(`/api/media/stats?account_id=${encodeURIComponent(acc)}`, { queued: 0, completed: 0, dead_404: 0, outsourced: 0, retained: 0, failed: 0, total: 0 }),
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
        const [cssRes, yamlRes] = await Promise.all([fetch(`${baseUrl}/plugins/${p}/skin/design.css`), fetch(`${baseUrl}/plugins/${p}/skin/layout.yaml`)]);
        return { platform: p, design_css: cssRes.ok ? await cssRes.text() : '', layout_yaml: yamlRes.ok ? await yamlRes.text() : '', controller: '' };
      } catch { return { platform: p, design_css: '', layout_yaml: '', controller: '' }; }
    },
    GetSkinCSS: async (p: string) => { try { const r = await fetch(`${baseUrl}/plugins/${p}/skin/design.css`); return r.ok ? await r.text() : ''; } catch { return ''; } },
    UpdateAccount: async (...args: any[]) => { console.log('[Polyfill] UpdateAccount:', args); },
    MergeAccounts: async (...args: any[]) => { console.log('[Polyfill] MergeAccounts:', args); return null; },
    SaveAvatarImage: async (p: string, k: string) => `${baseUrl}/avatars/${p || 'twitter'}/${k}.jpg`,
    ListAvailableAvatars: async (p: string) => [`${baseUrl}/avatars/${p || 'twitter'}/msluo14_avatar_001.jpg`, `${baseUrl}/avatars/${p || 'twitter'}/default_avatar.jpg`],
  };
  (window as any)._isWailsPolyfill = true;
  (window as any).go = { app: { App: mockApp }, main: { App: mockApp } };
}
