import { ref, onMounted, onUnmounted } from 'vue';
import { GetTimeline, GetAccounts, GetSystemLanguage, SearchArticles } from '../../wailsjs/go/main/App';
import { EventsOn } from '../../wailsjs/runtime/runtime';
import type { RenderTree, RenderAuthor } from '../models/RenderTree';
import { useKeyboardReload } from './useKeyboardReload';

export type LanguageCode = 'original' | 'ja' | 'en' | 'zh';
export type FilterType = 'all' | 'media' | 'reposts' | 'bookmarks';

export function useTimeline(platform: string = 'twitter') {
  const articles = ref<RenderTree[]>([]);
  const accounts = ref<RenderAuthor[]>([]);
  const selectedAccount = ref<string>('all');
  const currentFilter = ref<FilterType>('all');
  const searchQuery = ref<string>('');
  const systemLang = ref<LanguageCode>('ja');
  const loading = ref(false);
  const hasMore = ref(true);

  const fetchSystemLang = async () => {
    try {
      if (typeof GetSystemLanguage === 'function') {
        const lang = await GetSystemLanguage();
        if (lang) systemLang.value = lang as LanguageCode;
      }
    } catch (e) {
      console.warn('[useTimeline] Failed to fetch system language:', e);
    }
  };

  const fetchAccounts = async (retry = 3): Promise<void> => {
    try {
      const res = await GetAccounts(platform);
      accounts.value = res || [];
    } catch (e) {
      if (retry > 0) {
        await new Promise((r) => setTimeout(r, 500));
        return fetchAccounts(retry - 1);
      }
    }
  };

  const fetchTimeline = async (reset = false, retry = 3): Promise<void> => {
    if (loading.value || (!reset && !hasMore.value)) return;
    loading.value = true;
    const offset = reset ? 0 : articles.value.length;
    try {
      let items: RenderTree[] = [];
      if (searchQuery.value && searchQuery.value.trim() !== '') {
        const res = await SearchArticles(searchQuery.value.trim(), selectedAccount.value, currentFilter.value, 50, offset);
        items = (res && res.items) || [];
      } else {
        const res = await GetTimeline(platform, selectedAccount.value, currentFilter.value, 50, offset);
        items = res || [];
      }

      if (reset) articles.value = items;
      else articles.value.push(...items);
      hasMore.value = items.length === 50;
    } catch (e) {
      if (retry > 0) {
        loading.value = false;
        await new Promise((r) => setTimeout(r, 600));
        return fetchTimeline(reset, retry - 1);
      }
    } finally {
      loading.value = false;
    }
  };

  const reloadAll = async () => {
    hasMore.value = true;
    await Promise.all([fetchAccounts(), fetchTimeline(true)]);
  };

  // reloadAll 定義後に呼び出し
  useKeyboardReload(reloadAll);

  let unoffReady: (() => void) | null = null;
  onMounted(() => {
    fetchSystemLang();
    try {
      unoffReady = EventsOn('app:ready', () => {
        fetchSystemLang();
        reloadAll();
      });
    } catch (_) {}
    reloadAll();
  });

  onUnmounted(() => {
    if (unoffReady) unoffReady();
  });

  const setSearchQuery = (q: string) => {
    searchQuery.value = q;
    hasMore.value = true;
    fetchTimeline(true);
  };

  const clearSearchQuery = () => {
    searchQuery.value = '';
    hasMore.value = true;
    fetchTimeline(true);
  };

  return {
    articles, accounts, selectedAccount, currentFilter, searchQuery, systemLang,
    loading, hasMore,
    selectAccount: (id: string) => { selectedAccount.value = id; hasMore.value = true; fetchTimeline(true); },
    setFilter: (f: FilterType) => { currentFilter.value = f; hasMore.value = true; fetchTimeline(true); },
    setSearchQuery,
    clearSearchQuery,
    toggleLike: (id: string) => {
      const target = articles.value.find((item) => item.id === id);
      if (target) target.is_liked = !target.is_liked;
    },
    reloadAll, loadMore: () => fetchTimeline(false),
  };
}
