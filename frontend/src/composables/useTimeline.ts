// frontend/src/composables/useTimeline.ts (100行以下)
import { ref, onMounted, onUnmounted } from 'vue';
import { GetTimeline, GetAccounts, GetSystemLanguage, SearchArticles, RetryMediaDownload } from '../../wailsjs/go/app/App';
import { EventsOn } from '../../wailsjs/runtime/runtime';
import type { RenderTree, RenderAuthor } from '../models/RenderTree';

export type LanguageCode = 'original' | 'ja' | 'en' | 'zh';
export type FilterType = 'all' | 'media' | 'reposts' | 'bookmarks';

export function useTimeline(platform = 'twitter') {
  const articles = ref<RenderTree[]>([]), accounts = ref<RenderAuthor[]>([]);
  const selectedAccount = ref('all'), currentFilter = ref<FilterType>('all'), searchQuery = ref('');
  const systemLang = ref<LanguageCode>('ja'), loading = ref(false), hasMore = ref(true);
  const renderKey = ref(0), isStashReady = ref(false);

  const fetchSystemLang = async () => {
    try { if (typeof GetSystemLanguage === 'function') { const l = await GetSystemLanguage(); if (l) systemLang.value = l as LanguageCode; } } catch {}
  };

  const fetchAccounts = async (retry = 2): Promise<void> => {
    try {
      const getApp = (window as any)?.go?.app?.App || (window as any)?.go?.main?.App;
      if (getApp?.GetAccounts) accounts.value = (await getApp.GetAccounts(platform)) || [];
      else { const res = await fetch(`/api/accounts?platform=${encodeURIComponent(platform)}`); accounts.value = (await res.json()) || []; }
    } catch { if (retry > 0) { await new Promise((r) => setTimeout(r, 300)); return fetchAccounts(retry - 1); } }
  };

  const fetchTimeline = async (reset = false, retry = 2): Promise<void> => {
    if (!reset && (loading.value || !hasMore.value)) return;
    loading.value = true;
    const offset = reset ? 0 : articles.value.length;
    try {
      let items: RenderTree[] = [];
      const getApp = (window as any)?.go?.app?.App || (window as any)?.go?.main?.App;
      if (getApp) {
        if (searchQuery.value.trim()) { const res = await SearchArticles(searchQuery.value.trim(), selectedAccount.value, currentFilter.value, 50, offset); items = res?.items || []; }
        else items = (await GetTimeline(platform, selectedAccount.value, currentFilter.value, 50, offset)) || [];
      } else {
        const url = searchQuery.value.trim() ? `/api/search?q=${encodeURIComponent(searchQuery.value.trim())}&account_id=${encodeURIComponent(selectedAccount.value)}&filter=${encodeURIComponent(currentFilter.value)}&limit=50&offset=${offset}` : `/api/timeline?platform=${encodeURIComponent(platform)}&account_id=${encodeURIComponent(selectedAccount.value)}&filter=${encodeURIComponent(currentFilter.value)}&limit=50&offset=${offset}`;
        const res = await fetch(url); const data = await res.json(); items = Array.isArray(data) ? data : data?.items || [];
      }
      if (reset) articles.value = items; else articles.value.push(...items);
      hasMore.value = items.length === 50;
    } catch { if (retry > 0) { loading.value = false; await new Promise((r) => setTimeout(r, 300)); return fetchTimeline(reset, retry - 1); } }
    finally { loading.value = false; }
  };

  const reloadAll = async () => { loading.value = false; hasMore.value = true; renderKey.value++; await Promise.all([fetchAccounts(), fetchTimeline(true)]); };
  const checkInitialStashState = async () => {
    try {
      const getApp = (window as any)?.go?.app?.App || (window as any)?.go?.main?.App;
      if (getApp?.IsStashReady && (await getApp.IsStashReady())) { isStashReady.value = true; return; }
      const res = await fetch('/stash-proxy/', { method: 'HEAD' });
      if (res.ok || res.status === 401 || res.status === 404) isStashReady.value = true;
    } catch {}
  };

  const unoffs: (() => void)[] = [];
  onMounted(async () => {
    fetchSystemLang(); await checkInitialStashState();
    try {
      if ((window as any)?.runtime?.EventsOnMultiple) {
        unoffs.push(EventsOn('app:ready', () => { fetchSystemLang(); reloadAll(); }));
        unoffs.push(EventsOn('stash:ready', (ready: boolean) => { isStashReady.value = !!ready; if (ready) { renderKey.value++; reloadAll(); } }));
      }
    } catch {}
    reloadAll();
  });
  onUnmounted(() => { unoffs.forEach((fn) => { try { fn(); } catch {} }); });

  return {
    articles, accounts, selectedAccount, currentFilter, searchQuery, systemLang, loading, hasMore, renderKey, isStashReady,
    selectAccount: (id: string) => { selectedAccount.value = id; hasMore.value = true; fetchTimeline(true); },
    setFilter: (f: FilterType) => { currentFilter.value = f; hasMore.value = true; fetchTimeline(true); },
    setSearchQuery: (q: string) => { searchQuery.value = q; hasMore.value = true; fetchTimeline(true); },
    clearSearchQuery: () => { searchQuery.value = ''; hasMore.value = true; fetchTimeline(true); },
    reloadAll, loadMore: () => fetchTimeline(false),
    toggleLike: (id: string) => { const t = articles.value.find((i) => i.id === id); if (t) t.is_liked = !t.is_liked; },
    retryMedia: async (mId: string) => {
      for (const art of articles.value) { const m = art.media.find((i) => i.id === mId); if (m) { m.download_status = 'QUEUED'; break; } }
      try { if (typeof RetryMediaDownload === 'function') await RetryMediaDownload(mId); } catch {}
    },
  };
}
