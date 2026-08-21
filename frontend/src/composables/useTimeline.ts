// frontend/src/composables/useTimeline.ts (100行以下)
import { ref, onMounted, onUnmounted } from 'vue';
import { GetTimeline, GetAccounts, GetSystemLanguage, SearchArticles, RetryMediaDownload } from '../../wailsjs/go/main/App';
import { EventsOn } from '../../wailsjs/runtime/runtime';
import type { RenderTree, RenderAuthor } from '../models/RenderTree';
import { useKeyboardReload } from './useKeyboardReload';

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
    try { accounts.value = await GetAccounts(platform) || []; } catch (e) {
      if (retry > 0) { await new Promise((r) => setTimeout(r, 400)); return fetchAccounts(retry - 1); }
    }
  };

  const fetchTimeline = async (reset = false, retry = 2): Promise<void> => {
    if (loading.value || (!reset && !hasMore.value)) return;
    loading.value = true;
    const offset = reset ? 0 : articles.value.length;
    try {
      let items: RenderTree[] = [];
      if (searchQuery.value.trim() !== '') {
        const res = await SearchArticles(searchQuery.value.trim(), selectedAccount.value, currentFilter.value, 50, offset);
        items = (res && res.items) || [];
      } else { items = await GetTimeline(platform, selectedAccount.value, currentFilter.value, 50, offset) || []; }
      if (reset) articles.value = items; else articles.value.push(...items);
      hasMore.value = items.length === 50;
    } catch (e) {
      if (retry > 0) { loading.value = false; await new Promise((r) => setTimeout(r, 400)); return fetchTimeline(reset, retry - 1); }
    } finally { loading.value = false; }
  };

  const reloadAll = async () => { hasMore.value = true; await Promise.all([fetchAccounts(), fetchTimeline(true)]); };
  useKeyboardReload(reloadAll);

  let unoffReady: (() => void) | null = null;
  onMounted(() => {
    fetchSystemLang();
    try { unoffReady = EventsOn('app:ready', () => { fetchSystemLang(); reloadAll(); }); } catch (_) {}
    reloadAll();
  });
  onUnmounted(() => { if (unoffReady) unoffReady(); });

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
