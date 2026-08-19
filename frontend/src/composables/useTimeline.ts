// frontend/src/composables/useTimeline.ts (100行以下)
import { ref, onMounted } from 'vue';
import { GetTimeline } from '../../wailsjs/go/main/App';
import type { RenderTree } from '../models/RenderTree';

export type LanguageCode = 'original' | 'ja' | 'en' | 'zh';

export function useTimeline(platform: string = 'twitter') {
  const articles = ref<RenderTree[]>([]);
  const loading = ref(false);
  const error = ref<string | null>(null);
  const currentLang = ref<LanguageCode>('ja');

  // Wails Goバックエンドからタイムライン取得
  const fetchTimeline = async (
    accountID: string = 'all',
    filter: string = 'all',
    limit: number = 50,
    offset: number = 0
  ) => {
    loading.value = true;
    error.value = null;
    try {
      const result = await GetTimeline(platform, accountID, filter, limit, offset);
      console.log('[useTimeline] 受領データ:', result);
      if (result) {
        articles.value = offset === 0 ? result : [...articles.value, ...result];
      }
    } catch (err: any) {
      console.error('[useTimeline] Fetch failed:', err);
      error.value = err?.message || 'Failed to fetch timeline';
    } finally {
      loading.value = false;
    }
  };

  // 言語切り替え
  const setLanguage = (lang: LanguageCode) => {
    currentLang.value = lang;
  };

  // いいね（ブックマーク）状態の即時反転
  const toggleLike = (id: string) => {
    const target = articles.value.find((item) => item.id === id);
    if (target) {
      target.is_liked = !target.is_liked;
    }
  };

  // メディアの再取得ハンドラ
  const retryDownload = (mediaId: string) => {
    console.log('[useTimeline] Retry download requested for media:', mediaId);
  };

  // マウント時に自動で初期データを取得
  onMounted(() => {
    fetchTimeline();
  });

  return {
    articles,
    loading,
    error,
    currentLang,
    setLanguage,
    toggleLike,
    retryDownload,
    fetchTimeline,
  };
}
