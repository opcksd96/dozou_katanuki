// frontend/src/composables/useTimeline.ts (100行以下)
import { ref, onMounted } from 'vue';
import { GetTimeline, GetAccounts } from '../../wailsjs/go/main/App';
import type { RenderTree, RenderAuthor, RenderMedia } from '../models/RenderTree';

export type LanguageCode = 'original' | 'ja' | 'en' | 'zh';
export type FilterType = 'all' | 'media' | 'reposts' | 'bookmarks';

export function useTimeline(platform: string = 'twitter') {
  const articles = ref<RenderTree[]>([]);
  const accounts = ref<RenderAuthor[]>([]);
  const selectedAccount = ref<string>('all');
  const currentFilter = ref<FilterType>('all');
  const currentLang = ref<LanguageCode>('ja');
  const loading = ref(false);
  const hasMore = ref(true);
  const activeMedia = ref<RenderMedia | null>(null);

  const fetchAccounts = async () => {
    try {
      accounts.value = await GetAccounts(platform) || [];
    } catch (e) {
      console.error('[useTimeline] Failed to fetch accounts:', e);
    }
  };

  const fetchTimeline = async (reset: boolean = false) => {
    if (loading.value || (!reset && !hasMore.value)) return;
    loading.value = true;
    const offset = reset ? 0 : articles.value.length;
    try {
      const res = await GetTimeline(platform, selectedAccount.value, currentFilter.value, 50, offset);
      if (reset) articles.value = res || [];
      else articles.value.push(...(res || []));
      hasMore.value = (res || []).length === 50;
    } catch (e) {
      console.error('[useTimeline] Fetch timeline failed:', e);
    } finally {
      loading.value = false;
    }
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

  const setLanguage = (lang: LanguageCode) => { currentLang.value = lang; };

  const toggleLike = (id: string) => {
    const target = articles.value.find((item) => item.id === id);
    if (target) target.is_liked = !target.is_liked;
  };

  const openLightbox = (m: RenderMedia) => { activeMedia.value = m; };
  const closeLightbox = () => { activeMedia.value = null; };

  onMounted(() => {
    fetchAccounts();
    fetchTimeline(true);
  });

  return {
    articles, accounts, selectedAccount, currentFilter, currentLang,
    loading, hasMore, activeMedia, selectAccount, setFilter, setLanguage,
    toggleLike, openLightbox, closeLightbox, loadMore: () => fetchTimeline(false),
  };
}
