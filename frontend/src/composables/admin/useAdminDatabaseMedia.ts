// frontend/src/composables/admin/useAdminDatabaseMedia.ts (100行以下)
import { ref } from 'vue';

const getApp = () => (window as any)?.go?.app?.App || (window as any)?.go?.main?.App;

export function useAdminDatabaseMedia() {
  const mediaResults = ref<any[]>([]), mediaTotal = ref(0), isMediaLoading = ref(false), errorMessage = ref<string | null>(null);
  const mediaAccount = ref('all'), mediaPage = ref(1), mediaLimit = ref(24), mediaStatusFilter = ref('all'), mediaTypeFilter = ref<'all' | 'image' | 'video'>('all');
  const mediaStats = ref({ total_count: 0, image_count: 0, video_count: 0 });
  const downloadStatusStats = ref({ queued: 0, completed: 0, dead_404: 0, outsourced: 0, failed: 0, total: 0 });

  const fetchMedia = async (opts?: { account?: string; page?: number; limit?: number; status?: string; type?: 'all' | 'image' | 'video' }) => {
    if (opts?.account !== undefined) mediaAccount.value = opts.account;
    if (opts?.page !== undefined) mediaPage.value = Math.max(1, opts.page);
    if (opts?.limit !== undefined) mediaLimit.value = opts.limit;
    if (opts?.status !== undefined) mediaStatusFilter.value = opts.status;
    if (opts?.type !== undefined) mediaTypeFilter.value = opts.type;
    isMediaLoading.value = true; errorMessage.value = null;
    try {
      const app = getApp();
      if (app?.GetMediaList) {
        const offset = Math.max(0, (mediaPage.value - 1) * mediaLimit.value);
        const res = await app.GetMediaList(mediaAccount.value, mediaStatusFilter.value, mediaTypeFilter.value, mediaLimit.value, offset);
        mediaResults.value = res?.items || res?.Items || [];
        mediaTotal.value = res?.total || res?.Total || 0;
        if (res?.stats) mediaStats.value = res.stats;
      }
      if (app?.GetMediaDownloadStatusStats) {
        const ds = await app.GetMediaDownloadStatusStats(mediaAccount.value);
        if (ds) downloadStatusStats.value = ds;
      }
    } catch (e: any) { errorMessage.value = `メディア一覧の取得に失敗: ${e?.message || e}`; }
    finally { isMediaLoading.value = false; }
  };

  const setMediaPage = (p: number) => { mediaPage.value = Math.max(1, p); return fetchMedia(); };
  const setMediaLimit = (l: number) => { mediaLimit.value = l; mediaPage.value = 1; return fetchMedia(); };
  const setMediaAccount = (acc: string) => { mediaAccount.value = acc; mediaPage.value = 1; return fetchMedia(); };
  const setMediaStatusFilter = (st: string) => { mediaStatusFilter.value = st; mediaPage.value = 1; return fetchMedia(); };
  const setMediaTypeFilter = (t: 'all' | 'image' | 'video') => { mediaTypeFilter.value = t; mediaPage.value = 1; return fetchMedia(); };

  const updateMediaMetadata = async (mediaId: string, downloadStatus: string, stashSceneId: string, stashImageId: string, failedReason: string) => {
    try { const app = getApp(); if (app?.UpdateMediaMetadata) { await app.UpdateMediaMetadata(mediaId, downloadStatus, stashSceneId, stashImageId, failedReason); return true; } }
    catch (e: any) { errorMessage.value = `メタデータ更新に失敗: ${e?.message || e}`; }
    return false;
  };

  const purgeMedia = async (mId: string) => {
    try { const app = getApp(); if (app?.PurgeMedia) { await app.PurgeMedia(mId); return true; } }
    catch (e: any) { errorMessage.value = `メディアパージ失敗: ${e?.message || e}`; }
    return false;
  };

  const purgeMediaByStatus = async (status: string, account = 'all') => {
    try { const app = getApp(); if (app?.PurgeMediaByStatus) return await app.PurgeMediaByStatus(status, account); }
    catch (e: any) { errorMessage.value = `一括パージ失敗: ${e?.message || e}`; }
    return 0;
  };

  const toggleBookmark = async (mediaId: string) => {
    try {
      const app = getApp();
      if (app?.ToggleMediaBookmark) {
        const res = await app.ToggleMediaBookmark(mediaId);
        const item = mediaResults.value.find((m: any) => (m.media_id || m.id) === mediaId);
        if (item) item.is_bookmarked = res;
        return res;
      }
    } catch {}
  };

  const openInExplorer = async (mId: string) => { try { await getApp()?.OpenInExplorer?.(mId); } catch {} };
  const openWithDefaultApp = async (mId: string) => { try { await getApp()?.OpenWithDefaultApp?.(mId); } catch {} };
  const retryMedia = async (mId: string) => { try { await getApp()?.RetryMediaDownload?.(mId); } catch {} };

  return {
    mediaResults, mediaTotal, isMediaLoading, mediaAccount, mediaPage, mediaLimit, mediaStatusFilter, mediaTypeFilter, mediaStats, downloadStatusStats, errorMessage,
    fetchMedia, setMediaPage, setMediaLimit, setMediaAccount, setMediaStatusFilter, setMediaTypeFilter,
    updateMediaMetadata, purgeMedia, purgeMediaByStatus, toggleBookmark, openInExplorer, openWithDefaultApp, retryMedia,
  };
}
