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
        // Map the new array format ([]dto.RenderTree) back to expected format
        if (Array.isArray(res) && res.length > 0) {
          detail.value = {
            article: res[0],
            thread: res.slice(1)
          } as unknown as ArticleDetailResult;
        } else {
          detail.value = null;
        }
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
