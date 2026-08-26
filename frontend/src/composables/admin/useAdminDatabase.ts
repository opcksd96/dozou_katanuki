import { ref, computed } from 'vue';
import { useAdminDatabaseAccounts } from './useAdminDatabaseAccounts';
import { useAdminDatabaseArticles } from './useAdminDatabaseArticles';
import { useAdminDatabaseMedia } from './useAdminDatabaseMedia';

const getApp = () => (window as any)?.go?.app?.App || (window as any)?.go?.main?.App;

export function useAdminDatabase() {
  const activeSubTab = ref<'accounts' | 'posts' | 'media' | 'whitelist'>('accounts');
  const tableData = ref<{ columns: string[]; rows: any[]; total: number }>({ columns: [], rows: [], total: 0 });
  const accounts = useAdminDatabaseAccounts();
  const articles = useAdminDatabaseArticles();
  const media = useAdminDatabaseMedia();

  const errorMessage = computed(() => accounts.errorMessage.value || articles.errorMessage.value || media.errorMessage.value || null);
  const clearError = () => {
    accounts.errorMessage.value = null;
    articles.errorMessage.value = null;
    media.errorMessage.value = null;
  };

  const showAccountPosts = (id: string, onNav?: (tab: 'posts') => void) => {
    articles.searchAccount.value = id;
    articles.page.value = 1;
    activeSubTab.value = 'posts';
    articles.searchArticles();
    if (onNav) onNav('posts');
  };

  const showAccountMedia = (id: string, onNav?: (tab: 'media') => void) => {
    media.mediaAccount.value = id;
    media.mediaPage.value = 1;
    activeSubTab.value = 'media';
    media.fetchMedia();
    if (onNav) onNav('media');
  };

  const startMediaDownload = async (mId = '') => { try { return await getApp()?.StartMediaDownloadJob?.('twitter', mId); } catch {} };
  const startMediaPoll = async () => { try { return await getApp()?.StartMediaPollJob?.('twitter'); } catch {} };
  const startMediaEscalate = async () => { try { return await getApp()?.StartMediaEscalateJob?.('twitter'); } catch {} };
  const requeueMedia = async (status = 'DEAD_404') => { try { const count = await getApp()?.RequeueMediaByStatus?.(status, media.mediaAccount.value); await media.fetchMedia(); return count; } catch { return 0; } };
  const mergeDuplicates = async () => { try { const count = await getApp()?.MergeDuplicateMedia?.(); await media.fetchMedia(); return count; } catch { return 0; } };
  const purgeLowResDuplicates = async () => { try { const count = await getApp()?.PurgeLowerResolutionDuplicates?.(); await media.fetchMedia(); return count; } catch { return 0; } };
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
    searchFilter: articles.searchFilter, page: articles.page, limit: articles.limit, searchArticles: articles.searchArticles,
    setPage: articles.setPage, setAccount: articles.setAccount,
    saveTranslation: articles.saveTranslation, autoTranslate: articles.autoTranslate, startBatchTranslate: articles.startBatchTranslate,
    // Accounts
    accountsList: accounts.accountsList, selectedAccountDetail: accounts.selectedAccountDetail, isAccountLoading: accounts.isAccountLoading,
    fetchAccounts: accounts.fetchAccounts, selectAccount: accounts.selectAccount, updateAccount: accounts.updateAccount, mergeAccounts: accounts.mergeAccounts,
    saveAvatarImage: accounts.saveAvatarImage, fetchAvailableAvatars: accounts.fetchAvailableAvatars,
    // Media
    mediaResults: media.mediaResults, mediaTotal: media.mediaTotal, isMediaLoading: media.isMediaLoading,
    mediaAccount: media.mediaAccount, mediaPage: media.mediaPage, mediaLimit: media.mediaLimit,
    mediaStatusFilter: media.mediaStatusFilter, mediaTypeFilter: media.mediaTypeFilter, mediaStats: media.mediaStats,
    fetchMedia: media.fetchMedia, setMediaPage: media.setMediaPage, setMediaLimit: media.setMediaLimit,
    setMediaAccount: media.setMediaAccount, setMediaStatusFilter: media.setMediaStatusFilter, setMediaTypeFilter: media.setMediaTypeFilter,
    updateMediaMetadata: media.updateMediaMetadata, purgeMedia: media.purgeMedia,
    purgeMediaByStatus: (status: string) => media.purgeMediaByStatus(status, media.mediaAccount.value),
    toggleBookmark: media.toggleBookmark, openInExplorer: media.openInExplorer, openWithDefaultApp: media.openWithDefaultApp,
    retryMedia: media.retryMedia, showAccountPosts, showAccountMedia,
    startMediaDownload, startMediaPoll, startMediaEscalate, requeueMedia, mergeDuplicates, purgeLowResDuplicates, reconcileStashMedia,
    fetchStashMetadata, updateStashMetadata,
  };
}

