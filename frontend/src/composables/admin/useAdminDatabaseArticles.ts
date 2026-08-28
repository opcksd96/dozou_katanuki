// frontend/src/composables/admin/useAdminDatabaseArticles.ts (100行以下 - SPEC-PRINCIPLE-001)
import { ref, computed } from 'vue';
import { useArticleHistory } from './useArticleHistory';

const getApp = () => (window as any)?.go?.app?.App || (window as any)?.go?.main?.App;

export function useAdminDatabaseArticles() {
  const searchResults = ref<any[]>([]), totalCount = ref(0), isSearchLoading = ref(false), isTranslating = ref(false);
  const searchQuery = ref(''), searchAccount = ref('all'), searchFilter = ref('all'), includeTrash = ref(false);
  const limit = ref(50), errorMessage = ref<string | null>(null), history = useArticleHistory();

  const getEffectiveFilter = () => (includeTrash.value ? 'all_with_trash' : searchFilter.value);

  const searchArticles = async () => {
    isSearchLoading.value = true; errorMessage.value = null;
    try {
      const app = getApp(), filt = getEffectiveFilter();
      if (app?.SearchArticles) {
        const res = await app.SearchArticles(searchQuery.value.trim(), searchAccount.value, filt, limit.value, 0);
        searchResults.value = res?.items || res?.Items || []; totalCount.value = res?.total || res?.Total || 0;
      } else {
        const url = `/api/search?q=${encodeURIComponent(searchQuery.value.trim())}&account_id=${encodeURIComponent(searchAccount.value)}&filter=${encodeURIComponent(filt)}&limit=${limit.value}&offset=0`;
        const data = await (await fetch(url)).json();
        searchResults.value = data?.items || (Array.isArray(data) ? data : []); totalCount.value = data?.total ?? searchResults.value.length;
      }
    } catch (e: any) { errorMessage.value = `投稿の検索に失敗: ${e?.message || e}`; }
    finally { isSearchLoading.value = false; }
  };

  const fetchMoreArticles = async () => {
    if (isSearchLoading.value || searchResults.value.length >= totalCount.value) return;
    isSearchLoading.value = true;
    try {
      const app = getApp(), offset = searchResults.value.length, filt = getEffectiveFilter();
      if (app?.SearchArticles) {
        const res = await app.SearchArticles(searchQuery.value.trim(), searchAccount.value, filt, limit.value, offset);
        searchResults.value = [...searchResults.value, ...(res?.items || res?.Items || [])];
      }
    } catch (e: any) { errorMessage.value = `追加取得に失敗: ${e?.message || e}`; }
    finally { isSearchLoading.value = false; }
  };

  const setAccount = (acc: string) => { searchAccount.value = acc; return searchArticles(); };
  const toggleIncludeTrash = (val?: boolean) => { includeTrash.value = val !== undefined ? val : !includeTrash.value; return searchArticles(); };
  const saveTranslation = async (id: string, ja: string, en: string, zh: string) => {
    try { const app = getApp(); if (app?.UpdateArticleTranslations) await app.UpdateArticleTranslations(id, ja, en, zh); await searchArticles(); }
    catch (e: any) { errorMessage.value = `翻訳の保存に失敗: ${e?.message || e}`; }
  };

  const trashArticle = async (id: string, trashedBy = 'admin', reason = '') => {
    try {
      if (getApp()?.TrashArticle) await getApp().TrashArticle(id, trashedBy, reason);
      history.pushAction({ type: 'TRASH', ids: [id], trashedBy, reason });
      await searchArticles(); return true;
    } catch (e: any) { errorMessage.value = `削除に失敗: ${e?.message || e}`; return false; }
  };

  const batchTrashArticles = async (ids: string[], trashedBy = 'admin', reason = '') => {
    if (ids.length === 0) return false;
    try {
      if (getApp()?.BatchTrashArticles) await getApp().BatchTrashArticles(ids, trashedBy, reason);
      history.pushAction({ type: 'TRASH', ids, trashedBy, reason });
      await searchArticles(); return true;
    } catch (e: any) { errorMessage.value = `一括削除に失敗: ${e?.message || e}`; return false; }
  };

  const batchRestoreArticles = async (ids: string[]) => {
    if (ids.length === 0) return false;
    try {
      if (getApp()?.BatchRestoreArticles) await getApp().BatchRestoreArticles(ids);
      history.pushAction({ type: 'RESTORE', ids });
      await searchArticles(); return true;
    } catch (e: any) { errorMessage.value = `一括復元に失敗: ${e?.message || e}`; return false; }
  };

  const batchResetTranslations = async (ids: string[]) => {
    if (ids.length === 0) return false;
    try {
      const targetItems = searchResults.value.filter(a => ids.includes(a.id)).map(a => ({ id: a.id, ja: a.content?.ja || '', en: a.content?.en || '', zh: a.content?.zh || '' }));
      if (getApp()?.BatchResetTranslations) await getApp().BatchResetTranslations(ids);
      history.pushAction({ type: 'RESET_TRANSLATIONS', items: targetItems });
      await searchArticles(); return true;
    } catch (e: any) { errorMessage.value = `翻訳一括リセットに失敗: ${e?.message || e}`; return false; }
  };

  return {
    searchResults, totalCount, isSearchLoading, isTranslating, searchQuery, searchAccount, searchFilter, includeTrash, limit, errorMessage,
    hasMore: computed(() => searchResults.value.length < totalCount.value), canUndo: history.canUndo, canRedo: history.canRedo,
    searchArticles, fetchMoreArticles, setAccount, toggleIncludeTrash, saveTranslation, trashArticle, batchTrashArticles, batchRestoreArticles, batchResetTranslations,
    undo: () => history.undo({ batchTrashArticles, batchRestoreArticles, batchResetTranslations, saveTranslation }, searchArticles),
    redo: () => history.redo({ batchTrashArticles, batchRestoreArticles, batchResetTranslations, saveTranslation }, searchArticles),
    autoTranslate: async (id: string) => { isTranslating.value = true; try { return await getApp()?.AutoTranslateArticle?.(id); } finally { isTranslating.value = false; } },
    startBatchTranslate: (acc = 'all', ovw = false) => getApp()?.StartTranslateJob?.(acc, ovw),
  };
}
