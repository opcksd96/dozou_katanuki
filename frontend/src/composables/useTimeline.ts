// frontend/src/composables/useTimeline.ts (100行以下)
import { ref, onMounted, onUnmounted } from 'vue';
import { GetTimeline, GetAccounts, GetSystemLanguage, SearchArticles, RetryMediaDownload } from '../../wailsjs/go/main/App';
import { EventsOn } from '../../wailsjs/runtime/runtime';
import type { RenderTree, RenderAuthor } from '../models/RenderTree';

export type LanguageCode = 'original' | 'ja' | 'en' | 'zh';
export type FilterType = 'all' | 'media' | 'reposts' | 'bookmarks';

export function useTimeline(platform: string = 'twitter') {
  const articles = ref<RenderTree[]>([]), accounts = ref<RenderAuthor[]>([]);
  const selectedAccount = ref('all'), currentFilter = ref<FilterType>('all'), searchQuery = ref('');
  const systemLang = ref<LanguageCode>('ja'), loading = ref(false), hasMore = ref(true);

  const fetchSystemLang = async () => {
    try {
      if (typeof GetSystemLanguage === 'function') {
        const lang = await GetSystemLanguage();
        if (lang) systemLang.value = lang as LanguageCode;
      }
    } catch (_) {}
  };

  const fetchAccounts = async (retry = 2): Promise<void> => {
    try {
      if (typeof window !== 'undefined' && (window as any).go?.main?.App?.GetAccounts) {
        accounts.value = (await GetAccounts(platform)) || [];
      } else {
        const res = await fetch(`/api/accounts?platform=${encodeURIComponent(platform)}`);
        accounts.value = (await res.json()) || [];
      }
    } catch (e) {
      if (retry > 0) { await new Promise((r) => setTimeout(r, 300)); return fetchAccounts(retry - 1); }
    }
  };

  const fetchTimeline = async (reset = false, retry = 2): Promise<void> => {
    if (!reset && (loading.value || !hasMore.value)) return;
    loading.value = true;
    const offset = reset ? 0 : articles.value.length;
    try {
      let items: RenderTree[] = [];
      if (typeof window !== 'undefined' && (window as any).go?.main?.App) {
        if (searchQuery.value.trim() !== '') {
          const res = await SearchArticles(searchQuery.value.trim(), selectedAccount.value, currentFilter.value, 50, offset);
          items = (res && res.items) || [];
        } else {
          items = (await GetTimeline(platform, selectedAccount.value, currentFilter.value, 50, offset)) || [];
        }
      } else {
        if (searchQuery.value.trim() !== '') {
          const res = await fetch(`/api/search?q=${encodeURIComponent(searchQuery.value.trim())}&account_id=${encodeURIComponent(selectedAccount.value)}&filter=${encodeURIComponent(currentFilter.value)}&limit=50&offset=${offset}`);
          const data = await res.json();
          items = (data && data.items) || [];
        } else {
          const res = await fetch(`/api/timeline?platform=${encodeURIComponent(platform)}&account_id=${encodeURIComponent(selectedAccount.value)}&filter=${encodeURIComponent(currentFilter.value)}&limit=50&offset=${offset}`);
          items = (await res.json()) || [];
        }
      }
      if (reset) articles.value = items;
      else articles.value.push(...items);
      hasMore.value = items.length === 50;
    } catch (e) {
      if (retry > 0) { loading.value = false; await new Promise((r) => setTimeout(r, 300)); return fetchTimeline(reset, retry - 1); }
    } finally { loading.value = false; }
  };

  const reloadAll = async () => {
    loading.value = false; hasMore.value = true;
    await Promise.all([fetchAccounts(), fetchTimeline(true)]);
  };

  const unoffs: (() => void)[] = [];
  onMounted(() => {
    fetchSystemLang();
    try {
      if ((window as any)?.runtime?.EventsOnMultiple) {
        unoffs.push(EventsOn('app:ready', () => { fetchSystemLang(); reloadAll(); }));
        unoffs.push(EventsOn('stash:ready', () => { reloadAll(); }));
      }
    } catch (_) {}
    reloadAll();
  });
  onUnmounted(() => { unoffs.forEach((fn) => { try { fn(); } catch (_) {} }); });
  const setSearchQuery = (q: string) => { searchQuery.value = q; hasMore.value = true; fetchTimeline(true); };
  const clearSearchQuery = () => setSearchQuery('');

  return {
    articles, accounts, selectedAccount, currentFilter, searchQuery, systemLang, loading, hasMore,
    selectAccount: (id: string) => { selectedAccount.value = id; hasMore.value = true; fetchTimeline(true); },
    setFilter: (f: FilterType) => { currentFilter.value = f; hasMore.value = true; fetchTimeline(true); },
    setSearchQuery, clearSearchQuery, reloadAll, loadMore: () => fetchTimeline(false),
    toggleLike: (id: string) => {
      const t = articles.value.find((i) => i.id === id); if (t) t.is_liked = !t.is_liked;
    },
    retryMedia: async (mediaId: string) => {
      for (const art of articles.value) {
        const m = art.media.find((i) => i.id === mediaId);
        if (m) { m.download_status = 'QUEUED'; m.failed_reason = undefined; break; }
      }
      try { if (typeof RetryMediaDownload === 'function') await RetryMediaDownload(mediaId); } catch (_) {}
    },
  };
}
