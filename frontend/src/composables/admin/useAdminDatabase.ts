// frontend/src/composables/admin/useAdminDatabase.ts (100行以下)
import { ref } from 'vue';
const getApp = () => (window as any)?.go?.main?.App;

export function useAdminDatabase() {
  const activeSubTab = ref<'accounts' | 'posts' | 'media' | 'whitelist'>('accounts');
  const searchResults = ref<any[]>([]), totalCount = ref(0), isSearchLoading = ref(false), isTranslating = ref(false);
  const searchQuery = ref(''), searchAccount = ref('all'), searchFilter = ref('all'), page = ref(1), limit = ref(20);
  const errorMessage = ref<string | null>(null);

  const accountsList = ref<any[]>([]), selectedAccountDetail = ref<any>(null), isAccountLoading = ref(false);
  const mediaResults = ref<any[]>([]), mediaTotal = ref(0), isMediaLoading = ref(false);
  const mediaStatusFilter = ref('all'), mediaTypeFilter = ref<'all' | 'image' | 'video'>('all');
  const mediaStats = ref<{ total_count: number; image_count: number; video_count: number }>({ total_count: 0, image_count: 0, video_count: 0 });
  const tableData = ref<{ columns: string[]; rows: any[]; total: number }>({ columns: [], rows: [], total: 0 });

  const clearError = () => { errorMessage.value = null; };

  const searchArticles = async () => {
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
      console.error('[AdminDB] SearchArticles error:', e);
      errorMessage.value = `投稿の検索に失敗しました: ${e?.message || e}`;
    } finally {
      isSearchLoading.value = false;
    }
  };

  const fetchAccounts = async () => {
    isAccountLoading.value = true;
    errorMessage.value = null;
    try {
      const app = getApp();
      if (app?.ListAllAccounts) {
        const res = await app.ListAllAccounts();
        accountsList.value = res || [];
      }
    } catch (e: any) {
      console.error('[AdminDB] ListAllAccounts error:', e);
      errorMessage.value = `アカウント一覧の取得に失敗しました: ${e?.message || e}`;
    } finally {
      isAccountLoading.value = false;
    }
  };

  const selectAccount = async (numericId: string) => {
    isAccountLoading.value = true;
    errorMessage.value = null;
    try {
      const app = getApp();
      if (app?.GetAccountDetail) {
        selectedAccountDetail.value = await app.GetAccountDetail(numericId);
      }
    } catch (e: any) {
      console.error('[AdminDB] GetAccountDetail error:', e);
      errorMessage.value = `アカウント詳細の取得に失敗しました: ${e?.message || e}`;
    } finally {
      isAccountLoading.value = false;
    }
  };

  const fetchMedia = async () => {
    isMediaLoading.value = true;
    errorMessage.value = null;
    try {
      const app = getApp();
      if (app?.GetMediaList) {
        const offset = Math.max(0, (page.value - 1) * limit.value);
        const res = await app.GetMediaList(searchAccount.value, mediaStatusFilter.value, mediaTypeFilter.value, limit.value, offset);
        mediaResults.value = res?.items || res?.Items || [];
        mediaTotal.value = res?.total || res?.Total || 0;
        if (res?.stats) {
          mediaStats.value = res.stats;
        }
      }
    } catch (e: any) {
      console.error('[AdminDB] GetMediaList error:', e);
      errorMessage.value = `メディア一覧の取得に失敗しました: ${e?.message || e}`;
    } finally {
      isMediaLoading.value = false;
    }
  };

  const saveTranslation = async (id: string, ja: string, en: string, zh: string) => {
    try {
      const app = getApp();
      if (app?.UpdateArticleTranslations) {
        await app.UpdateArticleTranslations(id, ja, en, zh);
        await searchArticles();
      }
    } catch (e: any) {
      console.error('[AdminDB] UpdateArticleTranslations error:', e);
      errorMessage.value = `翻訳の保存に失敗しました: ${e?.message || e}`;
    }
  };

  const autoTranslate = async (id: string) => {
    const app = getApp();
    if (!app?.AutoTranslateArticle) return null;
    isTranslating.value = true;
    errorMessage.value = null;
    try {
      const res = await app.AutoTranslateArticle(id);
      return res;
    } catch (e: any) {
      console.error('[AdminDB] AutoTranslateArticle error:', e);
      errorMessage.value = `自動翻訳に失敗しました: ${e?.message || e}`;
      return null;
    } finally {
      isTranslating.value = false;
    }
  };

  const showAccountPosts = (numericId: string) => {
    searchAccount.value = numericId;
    page.value = 1;
    activeSubTab.value = 'posts';
    searchArticles();
  };

  const showAccountMedia = (numericId: string) => {
    searchAccount.value = numericId;
    page.value = 1;
    activeSubTab.value = 'media';
    fetchMedia();
  };

  const startBatchTranslate = async (acc = 'all', ovw = false) => {
    try { return await getApp()?.StartTranslateJob?.(acc, ovw); }
    catch (e: any) { errorMessage.value = `一括翻訳の開始に失敗しました: ${e?.message || e}`; }
  };
  const retryMedia = async (mId: string) => {
    try { await getApp()?.RetryMediaDownload?.(mId); await fetchMedia(); }
    catch (e: any) { errorMessage.value = `メディア再取得に失敗しました: ${e?.message || e}`; }
  };

  const purgeMedia = async (mId: string) => {
    try {
      const app = getApp();
      if (app?.PurgeMedia) {
        await app.PurgeMedia(mId);
        await fetchMedia();
        return true;
      }
    } catch (e: any) {
      errorMessage.value = `メディアのパージに失敗しました: ${e?.message || e}`;
    }
    return false;
  };

  const purgeMediaByStatus = async (status: string) => {
    try {
      const app = getApp();
      if (app?.PurgeMediaByStatus) {
        const count = await app.PurgeMediaByStatus(status, searchAccount.value);
        await fetchMedia();
        return count;
      }
    } catch (e: any) {
      errorMessage.value = `一括パージに失敗しました: ${e?.message || e}`;
    }
    return 0;
  };

  const updateAccount = async (numericId: string, displayName: string, username: string, avatarUrl: string, description: string) => {
    isAccountLoading.value = true;
    errorMessage.value = null;
    try {
      const app = getApp();
      if (app?.UpdateAccount) {
        await app.UpdateAccount(numericId, displayName, username, avatarUrl, description);
        await fetchAccounts();
        await selectAccount(numericId);
        return true;
      }
      return false;
    } catch (e: any) {
      console.error('[AdminDB] UpdateAccount error:', e);
      errorMessage.value = `アカウント情報の更新に失敗しました: ${e?.message || e}`;
      return false;
    } finally {
      isAccountLoading.value = false;
    }
  };

  const saveAvatarImage = async (platform: string, virtualKey: string, base64Data: string) => {
    errorMessage.value = null;
    try {
      const app = getApp();
      if (app?.SaveAvatarImage) {
        const savedPath = await app.SaveAvatarImage(platform, virtualKey, base64Data);
        return savedPath;
      }
      return null;
    } catch (e: any) {
      console.error('[AdminDB] SaveAvatarImage error:', e);
      errorMessage.value = `アバター画像の保存に失敗しました: ${e?.message || e}`;
      return null;
    }
  };

  const fetchAvailableAvatars = async (platform = 'twitter'): Promise<string[]> => {
    try {
      const app = getApp();
      if (app?.ListAvailableAvatars) {
        return (await app.ListAvailableAvatars(platform)) || [];
      }
      return [];
    } catch (e: any) {
      console.error('[AdminDB] ListAvailableAvatars error:', e);
      return [];
    }
  };

  return {
    activeSubTab, searchResults, totalCount, isSearchLoading, isTranslating, searchQuery, searchAccount,
    searchFilter, page, limit, errorMessage, clearError, accountsList, selectedAccountDetail, isAccountLoading,
    mediaResults, mediaTotal, isMediaLoading, mediaStatusFilter, mediaTypeFilter, mediaStats, tableData, searchArticles, fetchAccounts,
    selectAccount, updateAccount, saveAvatarImage, fetchAvailableAvatars, fetchMedia, saveTranslation, autoTranslate, startBatchTranslate, retryMedia,
    purgeMedia, purgeMediaByStatus,
    showAccountPosts, showAccountMedia,
  };
}
