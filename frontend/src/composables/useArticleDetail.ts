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
      const getApp = (window as any)?.go?.app?.App || (window as any)?.go?.main?.App;
      let res: any = null;
      if (getApp?.GetArticleDetail) {
        res = await getApp.GetArticleDetail(platform, id);
      } else if (typeof GetArticleDetail === 'function') {
        res = await GetArticleDetail(platform, id);
      } else {
        const resp = await fetch(`/api/article?platform=${encodeURIComponent(platform)}&id=${encodeURIComponent(id)}`);
        res = await resp.json();
      }

      if (res && res.article) {
        detail.value = {
          article: res.article,
          thread: Array.isArray(res.thread) ? res.thread : [],
        };
      } else if (Array.isArray(res) && res.length > 0) {
        detail.value = {
          article: res[0],
          thread: res.slice(1),
        };
      } else {
        detail.value = null;
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
