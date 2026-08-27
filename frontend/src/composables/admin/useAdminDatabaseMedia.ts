// frontend/src/composables/admin/useAdminDatabaseMedia.ts (100行以下 - SPEC-PRINCIPLE-001)
import { ref, computed } from 'vue';

const getApp = () => (window as any)?.go?.app?.App || (window as any)?.go?.main?.App;

export function useAdminDatabaseMedia() {
  const mediaResults = ref<any[]>([]), mediaTotal = ref(0), isMediaLoading = ref(false), errorMessage = ref<string | null>(null);
  const mediaAccount = ref('all'), mediaPage = ref(1), mediaLimit = ref(24), mediaStatusFilter = ref('all'), mediaTypeFilter = ref<'all' | 'image' | 'video'>('all');
  const mediaStats = ref({ total_count: 0, image_count: 0, video_count: 0 }), downloadStatusStats = ref({ queued: 0, completed: 0, dead_404: 0, outsourced: 0, failed: 0, total: 0 });
  const isWails = computed(() => typeof window !== 'undefined' && !!((window as any)?.go?.app?.App || (window as any)?.go?.main?.App) && !(window as any)?._isWailsPolyfill);

  const postApi = async (path: string, body: any) => {
    const res = await fetch(path, { method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify(body) });
    if (!res.ok) throw new Error(await res.text());
    return res.json();
  };

  const fetchMedia = async (opts?: { account?: string; page?: number; limit?: number; status?: string; type?: 'all' | 'image' | 'video' }) => {
    if (opts?.account !== undefined) mediaAccount.value = opts.account;
    if (opts?.page !== undefined) mediaPage.value = Math.max(1, opts.page);
    if (opts?.limit !== undefined) mediaLimit.value = opts.limit;
    if (opts?.status !== undefined) mediaStatusFilter.value = opts.status;
    if (opts?.type !== undefined) mediaTypeFilter.value = opts.type;
    isMediaLoading.value = true; errorMessage.value = null;
    try {
      const app = getApp(), offset = Math.max(0, (mediaPage.value - 1) * mediaLimit.value);
      if (app?.GetMediaList) {
        const res = await app.GetMediaList(mediaAccount.value, mediaStatusFilter.value, mediaTypeFilter.value, mediaLimit.value, offset);
        mediaResults.value = res?.items || res?.Items || []; mediaTotal.value = res?.total || res?.Total || 0;
        if (res?.stats) mediaStats.value = res.stats;
      } else {
        const res = await (await fetch(`/api/media?account_id=${encodeURIComponent(mediaAccount.value)}&status=${encodeURIComponent(mediaStatusFilter.value)}&type=${encodeURIComponent(mediaTypeFilter.value)}&limit=${mediaLimit.value}&offset=${offset}`)).json();
        mediaResults.value = res?.items || res?.Items || []; mediaTotal.value = res?.total || res?.Total || 0;
        if (res?.stats) mediaStats.value = res.stats;
      }
      if (app?.GetMediaDownloadStatusStats) {
        const ds = await app.GetMediaDownloadStatusStats(mediaAccount.value); if (ds) downloadStatusStats.value = ds;
      } else {
        const ds = await (await fetch(`/api/media/stats?account_id=${encodeURIComponent(mediaAccount.value)}`)).json(); if (ds) downloadStatusStats.value = ds;
      }
    } catch (e: any) { errorMessage.value = `メディア取得失敗: ${e?.message || e}`; }
    finally { isMediaLoading.value = false; }
  };

  const updateMediaMetadata = async (mediaId: string, downloadStatus: string, stashSceneId: string, stashImageId: string, failedReason: string) => {
    try {
      const app = getApp();
      if (app?.UpdateMediaMetadata) { await app.UpdateMediaMetadata(mediaId, downloadStatus, stashSceneId, stashImageId, failedReason); }
      else { await postApi('/api/media/update', { media_id: mediaId, download_status: downloadStatus, stash_scene_id: stashSceneId, stash_image_id: stashImageId, failed_reason: failedReason }); }
      return true;
    } catch (e: any) { errorMessage.value = `メタデータ更新失敗: ${e?.message || e}`; return false; }
  };

  const purgeMedia = async (mId: string) => {
    try {
      const app = getApp();
      if (app?.PurgeMedia) { await app.PurgeMedia(mId); } else { await postApi('/api/media/purge', { media_id: mId }); }
      mediaResults.value = mediaResults.value.filter((m: any) => (m.media_id || m.id) !== mId);
      return true;
    } catch (e: any) { errorMessage.value = `削除失敗: ${e?.message || e}`; return false; }
  };

  const purgeMediaByStatus = async (status: string, account = 'all') => {
    try {
      const app = getApp();
      if (app?.PurgeMediaByStatus) return await app.PurgeMediaByStatus(status, account);
      const res = await postApi('/api/media/purge-status', { status, account_id: account });
      return res?.purged_count || 0;
    } catch (e: any) { errorMessage.value = `一括削除失敗: ${e?.message || e}`; return 0; }
  };

  const toggleBookmark = async (mediaId: string) => {
    try {
      const app = getApp(); let res = false;
      if (app?.ToggleMediaBookmark) { res = await app.ToggleMediaBookmark(mediaId); }
      else { const data = await postApi('/api/media/bookmark', { media_id: mediaId }); res = !!data?.is_bookmarked; }
      const item = mediaResults.value.find((m: any) => (m.media_id || m.id) === mediaId);
      if (item) item.is_bookmarked = res;
      return res;
    } catch {}
  };

  const openInExplorer = async (mId: string) => { try { if (getApp()?.OpenInExplorer) await getApp().OpenInExplorer(mId); else await postApi('/api/media/open', { media_id: mId, action: 'explorer' }); } catch {} };
  const openWithDefaultApp = async (mId: string) => { try { if (getApp()?.OpenWithDefaultApp) await getApp().OpenWithDefaultApp(mId); else await postApi('/api/media/open', { media_id: mId, action: 'default' }); } catch {} };
  const retryMedia = async (mId: string) => { try { if (getApp()?.RetryMediaDownload) await getApp().RetryMediaDownload(mId); else await postApi('/api/media/requeue', { status: '', account_id: '' }); } catch {} };

  return {
    mediaResults, mediaTotal, isMediaLoading, mediaAccount, mediaPage, mediaLimit, mediaStatusFilter, mediaTypeFilter, mediaStats, downloadStatusStats, errorMessage, isWails,
    fetchMedia, setMediaPage: (p: number) => { mediaPage.value = Math.max(1, p); return fetchMedia(); },
    setMediaLimit: (l: number) => { mediaLimit.value = l; mediaPage.value = 1; return fetchMedia(); },
    setMediaAccount: (acc: string) => { mediaAccount.value = acc; mediaPage.value = 1; return fetchMedia(); },
    setMediaStatusFilter: (st: string) => { mediaStatusFilter.value = st; mediaPage.value = 1; return fetchMedia(); },
    setMediaTypeFilter: (t: 'all' | 'image' | 'video') => { mediaTypeFilter.value = t; mediaPage.value = 1; return fetchMedia(); },
    updateMediaMetadata, purgeMedia, purgeMediaByStatus, toggleBookmark, openInExplorer, openWithDefaultApp, retryMedia,
  };
}
