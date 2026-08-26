// frontend/src/utils/wailsPolyfill.ts (100行以下)
export function initWailsPolyfill() {
  if (typeof window === 'undefined') return;
  if (!window.runtime) {
    (window as any).runtime = {
      EventsOn: () => () => {},
      EventsOnMultiple: () => () => {},
      EventsOff: () => {},
      EventsOffAll: () => {},
      EventsOnce: () => () => {},
      EventsEmit: () => {},
      BrowserOpenURL: (url: string) => window.open(url, '_blank', 'noopener,noreferrer'),
      LogPrint: console.log, LogTrace: console.trace, LogDebug: console.debug,
      LogInfo: console.info, LogWarning: console.warn, LogError: console.error,
    };
  }

  if (!window.go) {
    (window as any).go = {
      main: {
        App: {
          GetSystemLanguage: async () => 'ja',
          GetAccounts: async (platform: string) => {
            try { const res = await fetch(`/api/accounts?platform=${encodeURIComponent(platform || 'twitter')}`); return await res.json(); } catch { return []; }
          },
          GetTimeline: async (platform: string, accountId: string, filter: string, limit: number, offset: number) => {
            try { const res = await fetch(`/api/timeline?platform=${encodeURIComponent(platform || 'twitter')}&account_id=${encodeURIComponent(accountId || 'all')}&filter=${encodeURIComponent(filter || 'all')}&limit=${limit || 50}&offset=${offset || 0}`); return await res.json(); } catch { return []; }
          },
          SearchArticles: async (query: string, accountId: string, filter: string, limit: number, offset: number) => {
            try { const res = await fetch(`/api/search?q=${encodeURIComponent(query || '')}&account_id=${encodeURIComponent(accountId || 'all')}&filter=${encodeURIComponent(filter || 'all')}&limit=${limit || 50}&offset=${offset || 0}`); return await res.json(); } catch { return { items: [], total: 0 }; }
          },
          GetArticleDetail: async (platform: string, id: string) => {
            try { const res = await fetch(`/api/article?platform=${encodeURIComponent(platform || 'twitter')}&id=${encodeURIComponent(id)}`); return await res.json(); } catch { return null; }
          },
          GetBroadcastStatus: async () => {
            try { const res = await fetch('/api/broadcast/status'); return await res.json(); } catch { return { enabled: true, use_tls: true, port: 5175, active_clients: 0, networks: [] }; }
          },
          GetConfig: async () => {
            return {
              system: { language: 'ja', default_framework: 'twitter', env: 'production' },
              storage: { db_path: './archive.db', local_media_dir: './blobs', stash_dir: './stash', dumps_dir: './dumps', stash_enabled: true },
              network: { stash_port: 9999, frontend_port: 5173, internal_bind_address: '127.0.0.1', middleware_port: 5175, public_bind_address: '0.0.0.0' },
              scheduler: { poll_interval_sec: 300, backup_interval_hours: 24, max_backup_generations: 7 },
              broadcast: { enabled: true, allowed_networks: ['192.168.10.0/24', '192.168.3.0/24', '127.0.0.1/32'] },
              appearance: { font_family_ja: 'Hiragino Sans, Meiryo, sans-serif', font_family_en: 'Nunito, sans-serif', font_family_zh: 'Microsoft YaHei, SimHei, sans-serif' },
            };
          },
          GetSkinPackage: async (platform: string) => {
            try {
              const [cssRes, yamlRes] = await Promise.all([fetch(`/plugins/${platform}/skin/design.css`), fetch(`/plugins/${platform}/skin/layout.yaml`)]);
              return { platform, design_css: cssRes.ok ? await cssRes.text() : '', layout_yaml: yamlRes.ok ? await yamlRes.text() : '', controller: '' };
            } catch { return { platform, design_css: '', layout_yaml: '', controller: '' }; }
          },
          GetSkinCSS: async (platform: string) => {
            try { const res = await fetch(`/plugins/${platform}/skin/design.css`); return res.ok ? await res.text() : ''; } catch { return ''; }
          },
          RetryMediaDownload: async () => {},
          ToggleBroadcast: async () => {},
          SaveConfig: async () => {},
           UpdateAccount: async (numericId: string, displayName: string, username: string, avatarUrl: string, description: string, aliasOf: string, groupName: string) => {
             console.log('[Polyfill] UpdateAccount:', { numericId, displayName, username, avatarUrl, description, aliasOf, groupName });
           },
           MergeAccounts: async (sourceId: string, targetId: string) => {
             console.log('[Polyfill] MergeAccounts:', { sourceId, targetId });
             return null;
           },
          SaveAvatarImage: async (platform: string, virtualKey: string, base64Data: string) => {
            console.log('[Polyfill] SaveAvatarImage:', { platform, virtualKey, dataLen: base64Data?.length });
            return `/avatars/${platform || 'twitter'}/${virtualKey}.jpg`;
          },
          ListAvailableAvatars: async (platform: string) => {
            return [
              `/avatars/${platform || 'twitter'}/msluo14_avatar_001.jpg`,
              `/avatars/${platform || 'twitter'}/mash_kyrielight_avatar_001.jpg`,
              `/avatars/${platform || 'twitter'}/default_avatar.jpg`,
            ];
          },
        }
      }
    };
  }
}
