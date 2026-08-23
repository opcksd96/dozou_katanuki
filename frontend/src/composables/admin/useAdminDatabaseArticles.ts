// frontend/src/composables/admin/useAdminDatabaseArticles.ts (100行以下)
import { ref } from 'vue';

const getApp = () => (window as any)?.go?.app?.App || (window as any)?.go?.main?.App;

export function useAdminDatabaseArticles() {
  const searchResults = ref<any[]>([]);
  const totalCount = ref(0);
  const isSearchLoading = ref(false);
  const isTranslating = ref(false);
  const searchQuery = ref('');
  const searchAccount = ref('all');
  const searchFilter = ref('all');
  const page = ref(1);
  const limit = ref(20);
  const errorMessage = ref<string | null>(null);

  const searchArticles = async (targetPage?: number) => {
    if (typeof targetPage === 'number') page.value = Math.max(1, targetPage);
    isSearchLoading.value = true;
    errorMessage.value = null;
    try {
      const app = getApp();
      if (app?.SearchArticles) {
        const offset = Math.max(0, (page.value - 1) * limit.value);
        const res = await app.SearchArticles(searchQuery.value.trim(), searchAccount.value, searchFilter.value, limit.value, offset);
        searchResults.value = res?.items || res?.Items || [];
        totalCount.value = res?.total || res?.Total || 0;
      }
    } catch (e: any) {
      errorMessage.value = `投稿の検索に失敗: ${e?.message || e}`;
    } finally {
      isSearchLoading.value = false;
    }
  };

  const setPage = (p: number) => {
    page.value = Math.max(1, p);
    return searchArticles();
  };

  const setAccount = (acc: string) => {
    searchAccount.value = acc;
    page.value = 1;
    return searchArticles();
  };

  const saveTranslation = async (id: string, ja: string, en: string, zh: string) => {
    try {
      const app = getApp();
      if (app?.UpdateArticleTranslations) {
        await app.UpdateArticleTranslations(id, ja, en, zh);
        await searchArticles();
      }
    } catch (e: any) {
      errorMessage.value = `翻訳の保存に失敗: ${e?.message || e}`;
    }
  };

  const autoTranslate = async (id: string) => {
    const app = getApp();
    if (!app?.AutoTranslateArticle) return null;
    isTranslating.value = true;
    errorMessage.value = null;
    try {
      return await app.AutoTranslateArticle(id);
    } catch (e: any) {
      errorMessage.value = `自動翻訳に失敗: ${e?.message || e}`;
      return null;
    } finally {
      isTranslating.value = false;
    }
  };

  const startBatchTranslate = async (acc = 'all', ovw = false) => {
    try { return await getApp()?.StartTranslateJob?.(acc, ovw); }
    catch (e: any) { errorMessage.value = `一括翻訳の開始に失敗: ${e?.message || e}`; }
  };

  return {
    searchResults, totalCount, isSearchLoading, isTranslating,
    searchQuery, searchAccount, searchFilter, page, limit, errorMessage,
    searchArticles, setPage, setAccount, saveTranslation, autoTranslate, startBatchTranslate,
  };
}

