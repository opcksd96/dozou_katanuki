import { ref, computed } from 'vue';
import type { RenderTree } from '../models/RenderTree';
import mockData from '../mock/sample_render_trees.json';

export type SupportedLang = 'original' | 'ja' | 'en' | 'zh';

export function useTimeline() {
  const articles = ref<RenderTree[]>(mockData as RenderTree[]);
  const currentLang = ref<SupportedLang>('ja');
  const isLoading = ref<boolean>(false);

  const setLanguage = (lang: SupportedLang) => {
    currentLang.value = lang;
  };

  const toggleLike = (id: string) => {
    const target = articles.value.find((a) => a.id === id);
    if (target) {
      target.is_liked = !target.is_liked;
      target.metrics.likes += target.is_liked ? 1 : -1;
    }
  };

  const retryDownload = (mediaId: string) => {
    console.log(`[UDF Action] メディア再取得要求: ${mediaId}`);
  };

  return {
    articles,
    currentLang,
    isLoading,
    articleCount: computed(() => articles.value.length),
    setLanguage,
    toggleLike,
    retryDownload,
  };
}
