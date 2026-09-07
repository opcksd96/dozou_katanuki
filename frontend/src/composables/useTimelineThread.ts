// frontend/src/composables/useTimelineThread.ts (100行以下)
import { ref } from 'vue';
import { GetArticleDetail } from '../../wailsjs/go/app/App';
import type { RenderTree } from '../models/RenderTree';

export function useTimelineThread(platform = 'twitter') {
  const threadsMap = ref<Record<string, RenderTree[]>>({});
  const loadingMap = ref<Record<string, boolean>>({});
  const expandedMap = ref<Record<string, boolean>>({});

  const toggleExpand = async (articleId: string) => {
    if (expandedMap.value[articleId]) {
      expandedMap.value[articleId] = false;
      return;
    }
    expandedMap.value[articleId] = true;
    if (!threadsMap.value[articleId]) {
      loadingMap.value[articleId] = true;
      try {
        const getApp = (window as any)?.go?.app?.App || (window as any)?.go?.main?.App;
        let res: any = null;
        if (getApp?.GetArticleDetail) {
          res = await getApp.GetArticleDetail(platform, articleId);
        } else if (typeof GetArticleDetail === 'function') {
          res = await GetArticleDetail(platform, articleId);
        } else {
          res = await (await fetch(`/api/article?platform=${encodeURIComponent(platform)}&id=${encodeURIComponent(articleId)}`)).json();
        }
        if (res && res.thread) {
          threadsMap.value[articleId] = res.thread;
        } else if (Array.isArray(res) && res.length > 1) {
          threadsMap.value[articleId] = res.slice(1);
        }
      } catch (err) {
        console.error('Failed to load thread:', err);
      } finally {
        loadingMap.value[articleId] = false;
      }
    }
  };

  const getParentArticles = (article: RenderTree): RenderTree[] => {
    const thread = threadsMap.value[article.id];
    if (!thread) return [];
    const mainTime = new Date(article.created_at).getTime();
    return thread
      .filter((t) => t.id !== article.id && new Date(t.created_at).getTime() < mainTime)
      .sort((a, b) => new Date(a.created_at).getTime() - new Date(b.created_at).getTime());
  };

  return {
    threadsMap,
    loadingMap,
    expandedMap,
    toggleExpand,
    getParentArticles,
  };
}
