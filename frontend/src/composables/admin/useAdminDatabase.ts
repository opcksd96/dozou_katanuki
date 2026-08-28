// frontend/src/composables/admin/useAdminDatabase.ts (100行以下 - SPEC-PRINCIPLE-001)
import { ref, computed } from 'vue';
import { useAdminDatabaseAccounts } from './useAdminDatabaseAccounts';
import { useAdminDatabaseArticles } from './useAdminDatabaseArticles';
import { useAdminDatabaseMedia } from './useAdminDatabaseMedia';

const getApp = () => (window as any)?.go?.app?.App || (window as any)?.go?.main?.App;

export function useAdminDatabase() {
  const activeSubTab = ref<'accounts' | 'posts' | 'media'>('accounts');
  const tableData = ref<{ columns: string[]; rows: any[]; total: number }>({ columns: [], rows: [], total: 0 });
  const accounts = useAdminDatabaseAccounts();
  const articles = useAdminDatabaseArticles();
  const media = useAdminDatabaseMedia();

  const errorMessage = computed(() => accounts.errorMessage.value || articles.errorMessage.value || media.errorMessage.value || null);
  const clearError = () => { accounts.errorMessage.value = null; articles.errorMessage.value = null; media.errorMessage.value = null; };

  const showAccountPosts = (id: string, onNav?: (tab: 'posts') => void) => {
    articles.searchAccount.value = id; articles.searchQuery.value = ''; activeSubTab.value = 'posts';
    articles.searchArticles(); if (onNav) onNav('posts');
  };

  const showAccountMedia = (id: string, onNav?: (tab: 'media') => void) => {
    media.mediaAccount.value = id; media.mediaPage.value = 1; activeSubTab.value = 'media';
    media.fetchMedia(); if (onNav) onNav('media');
  };

  const startMediaDownload = async (mId = '') => { try { return await getApp()?.StartMediaDownloadJob?.('twitter', mId); } catch {} };
  const startMediaPoll = async () => { try { return await getApp()?.StartMediaPollJob?.('twitter'); } catch {} };
  const startMediaEscalate = async () => { try { return await getApp()?.StartMediaEscalateJob?.('twitter'); } catch {} };
  const startSmartRecovery = async () => { try { return await getApp()?.StartSmartRecoveryJob?.('twitter'); } catch {} };
  const startThunderEscalate = async () => { try { return await getApp()?.StartThunderEscalateJob?.('twitter'); } catch {} };
  const requeueMedia = async (status = 'DEAD_404') => { try { const count = await getApp()?.RequeueMediaByStatus?.(status, media.mediaAccount.value); await media.fetchMedia(); return count; } catch { return 0; } };
  const reconcileStashMedia = async () => { try { const count = await getApp()?.ReconcileStashMedia?.(); await media.fetchMedia(); return count; } catch { return 0; } };
  const fetchStashMetadata = async (sId: string, iId: string) => { try { return await getApp()?.GetStashMetadata?.(sId, iId); } catch { return null; } };
  const updateStashMetadata = async (isScene: boolean, id: string, title: string, details: string, rating100: number) => {
    try { return await getApp()?.UpdateStashMetadata?.(isScene, id, title, details, rating100); } catch { return null; }
  };

  return {
    activeSubTab, tableData, errorMessage, clearError,
    // Articles
    searchResults: articles.searchResults, totalCount: articles.totalCount, isSearchLoading: articles.isSearchLoading,
    isTranslating: articles.isTranslating, searchQuery: articles.searchQuery, searchAccount: articles.searchAccount,
    searchFilter: articles.searchFilter, includeTrash: articles.includeTrash, hasMore: articles.hasMore,
    canUndo: articles.canUndo, canRedo: articles.canRedo, searchArticles: articles.searchArticles, fetchMoreArticles: articles.fetchMoreArticles,
    setAccount: articles.setAccount, toggleIncludeTrash: articles.toggleIncludeTrash, saveTranslation: articles.saveTranslation,
    autoTranslate: articles.autoTranslate, startBatchTranslate: articles.startBatchTranslate,
    trashArticle: articles.trashArticle, batchTrashArticles: articles.batchTrashArticles, batchRestoreArticles: articles.batchRestoreArticles,
    batchResetTranslations: articles.batchResetTranslations, undo: articles.undo, redo: articles.redo,
    // Accounts
    accountsList: accounts.accountsList, selectedAccountDetail: accounts.selectedAccountDetail, isAccountLoading: accounts.isAccountLoading,
    showTrash: accounts.showTrash, toggleShowTrash: accounts.toggleShowTrash, trashAccount: accounts.trashAccount, restoreAccount: accounts.restoreAccount,
    fetchAccounts: accounts.fetchAccounts, selectAccount: accounts.selectAccount, toggleAccountWhitelist: accounts.toggleAccountWhitelist,
    updateAccount: accounts.updateAccount, mergeAccounts: accounts.mergeAccounts, saveAvatarImage: accounts.saveAvatarImage,
    fetchAvailableAvatars: accounts.fetchAvailableAvatars, showAccountPosts, showAccountMedia,
    // Media
    mediaResults: media.mediaResults, mediaTotal: media.mediaTotal, mediaStats: media.mediaStats,
    downloadStatusStats: media.downloadStatusStats, isMediaLoading: media.isMediaLoading,
    mediaAccount: media.mediaAccount, mediaStatusFilter: media.mediaStatusFilter, mediaTypeFilter: media.mediaTypeFilter,
    mediaPage: media.mediaPage, mediaLimit: media.mediaLimit, fetchMedia: media.fetchMedia, setMediaAccount: media.setMediaAccount,
    setMediaStatusFilter: media.setMediaStatusFilter, setMediaTypeFilter: media.setMediaTypeFilter, setMediaPage: media.setMediaPage,
    setMediaLimit: media.setMediaLimit, retryMedia: media.retryMedia, purgeMedia: media.purgeMedia, purgeMediaByStatus: media.purgeMediaByStatus,
    trashMedia: media.trashMedia, restoreMedia: media.restoreMedia,
    openInExplorer: media.openInExplorer, openWithDefaultApp: media.openWithDefaultApp, toggleBookmark: media.toggleBookmark,
    startMediaDownload, startMediaPoll, startMediaEscalate, startSmartRecovery, startThunderEscalate, requeueMedia, reconcileStashMedia, fetchStashMetadata, updateStashMetadata,
  };
}
