// frontend/src/composables/useTimeline.ts (100行以下)
import { ref, onMounted, onUnmounted } from 'vue';
import { GetTimeline, GetAccounts, GetSystemLanguage } from '../../wailsjs/go/main/App';
import { EventsOn, WindowReload } from '../../wailsjs/runtime/runtime';
import type { RenderTree, RenderAuthor } from '../models/RenderTree';

export type LanguageCode = 'original' | 'ja' | 'en' | 'zh';
export type FilterType = 'all' | 'media' | 'reposts' | 'bookmarks';

export function useTimeline(platform: string = 'twitter') {
  const articles = ref<RenderTree[]>([]);
  const accounts = ref<RenderAuthor[]>([]);
  const selectedAccount = ref<string>('all');
  const currentFilter = ref<FilterType>('all');
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

  const fetchAccounts = async (retryCount = 3): Promise<void> => {
    try {
      const res = await GetAccounts(platform);
      accounts.value = res || [];
    } catch (e) {
      if (retryCount > 0) {
        await new Promise((r) => setTimeout(r, 500));
        return fetchAccounts(retryCount - 1);
      }
    }
  };

  const fetchTimeline = async (reset: boolean = false, retryCount = 3): Promise<void> => {
    if (loading.value || (!reset && !hasMore.value)) return;
    loading.value = true;
    const offset = reset ? 0 : articles.value.length;
    try {
      const res = await GetTimeline(platform, selectedAccount.value, currentFilter.value, 50, offset);
      if (reset) articles.value = res || [];
      else articles.value.push(...(res || []));
      hasMore.value = (res || []).length === 50;
    } catch (e) {
      if (retryCount > 0) {
        loading.value = false;
        await new Promise((r) => setTimeout(r, 600));
        return fetchTimeline(reset, retryCount - 1);
      }
    } finally {
      loading.value = false;
    }
  };

  const reloadAll = async () => {
    hasMore.value = true;
    await Promise.all([fetchAccounts(), fetchTimeline(true)]);
  };

  const selectAccount = (id: string) => {
    selectedAccount.value = id;
    hasMore.value = true;
    fetchTimeline(true);
  };

  const setFilter = (filter: FilterType) => {
    currentFilter.value = filter;
    hasMore.value = true;
    fetchTimeline(true);
  };

  const toggleLike = (id: string) => {
    const target = articles.value.find((item) => item.id === id);
    if (target) target.is_liked = !target.is_liked;
  };

  const handleKeyDown = (e: KeyboardEvent) => {
    const isCtrlR = (e.ctrlKey || e.metaKey) && (e.key === 'r' || e.key === 'R' || e.code === 'KeyR');
    if (isCtrlR || e.key === 'F5' || e.code === 'F5') {
      e.preventDefault();
      try {
        if (typeof WindowReload === 'function') WindowReload();
        else window.location.reload();
      } catch (_) {
        window.location.reload();
      }
    }
  };

  let unoffReady: (() => void) | null = null;
  onMounted(() => {
    window.addEventListener('keydown', handleKeyDown, true);
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
    window.removeEventListener('keydown', handleKeyDown, true);
    if (unoffReady) unoffReady();
  });

  return {
    articles, accounts, selectedAccount, currentFilter, systemLang,
    loading, hasMore, selectAccount, setFilter,
    toggleLike, reloadAll, loadMore: () => fetchTimeline(false),
  };
}
