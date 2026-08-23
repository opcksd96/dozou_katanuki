import { ref } from 'vue';
import { GetArticleDetail } from '../../wailsjs/go/app/App';
import type { RenderTree } from '../models/RenderTree';

export interface ArticleDetailResult {
  article: RenderTree;
  thread: RenderTree[];
}

export function useArticleDetail(platform: string = 'twitter') {
  const detail = ref<ArticleDetailResult | null>(null);
  const loading = ref(false);
  const error = ref<string | null>(null);

  const fetchDetail = async (id: string) => {
    loading.value = true;
    error.value = null;
    try {
      if (typeof GetArticleDetail === 'function') {
        const res = await GetArticleDetail(platform, id);
        detail.value = res;
      } else {
        error.value = 'GetArticleDetail is not available';
      }
    } catch (e: any) {
      error.value = e?.message || 'Failed to fetch article detail';
    } finally {
      loading.value = false;
    }
  };

  const clearDetail = () => {
    detail.value = null;
    error.value = null;
  };

  return {
    detail,
    loading,
    error,
    fetchDetail,
    clearDetail,
  };
}
