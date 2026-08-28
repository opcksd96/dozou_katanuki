// frontend/src/composables/admin/useAdminDatabaseMedia.ts (100行以下 - SPEC-PRINCIPLE-001)
import { ref, computed, onMounted, onUnmounted } from 'vue';
import { EventsOn, EventsOff } from '../../../wailsjs/runtime/runtime';
import { useToast } from '../useToast';

const getApp = () => (window as any)?.go?.app?.App || (window as any)?.go?.main?.App;

export function useAdminDatabaseMedia() {
  const { addToast } = useToast();
  const mediaResults = ref<any[]>([]), mediaTotal = ref(0), isMediaLoading = ref(false), errorMessage = ref<string | null>(null);
  const mediaAccount = ref('all'), mediaPage = ref(1), mediaLimit = ref(24), mediaStatusFilter = ref('all'), mediaTypeFilter = ref<'all' | 'image' | 'video'>('all');
  const mediaStats = ref({ total_count: 0, image_count: 0, video_count: 0 }), downloadStatusStats = ref({ queued: 0, completed: 0, dead_404: 0, outsourced: 0, failed: 0, total: 0 });
  const isWails = computed(() => typeof window !== 'undefined' && !!((window as any)?.go?.app?.App || (window as any)?.go?.main?.App) && !(window as any)?._isWailsPolyfill);

  const fetchMedia = async (opts?: any) => {
    if (opts?.account !== undefined) mediaAccount.value = opts.account; if (opts?.page !== undefined) mediaPage.value = Math.max(1, opts.page);
    if (opts?.limit !== undefined) mediaLimit.value = opts.limit; if (opts?.status !== undefined) mediaStatusFilter.value = opts.status; if (opts?.type !== undefined) mediaTypeFilter.value = opts.type;
    isMediaLoading.value = true; errorMessage.value = null;
    try {
      const app = getApp(), offset = Math.max(0, (mediaPage.value - 1) * mediaLimit.value);
      const res = app?.GetMediaList ? await app.GetMediaList(mediaAccount.value, mediaStatusFilter.value, mediaTypeFilter.value, mediaLimit.value, offset) : await (await fetch(`/api/media?account_id=${encodeURIComponent(mediaAccount.value)}&status=${encodeURIComponent(mediaStatusFilter.value)}&type=${encodeURIComponent(mediaTypeFilter.value)}&limit=${mediaLimit.value}&offset=${offset}`)).json();
      const raw = res?.items || res?.Items || [];
      mediaResults.value = raw.filter((m: any) => mediaStatusFilter.value === 'TRASH' ? !!m?.is_trash : !m?.is_trash);
      mediaTotal.value = res?.total || res?.Total || 0; if (res?.stats) mediaStats.value = res.stats;
      if (app?.GetMediaDownloadStatusStats) { const ds = await app.GetMediaDownloadStatusStats(mediaAccount.value); if (ds) downloadStatusStats.value = ds; }
    } catch (e: any) { errorMessage.value = `メディア取得失敗: ${e?.message || e}`; }
    finally { isMediaLoading.value = false; }
  };

  const trashMedia = async (idOrPayload: any, maybeReason = '手動整理') => {
    const mId = String(typeof idOrPayload === 'object' ? (idOrPayload?.mediaId || idOrPayload?.media_id || idOrPayload?.id) : idOrPayload).trim();
    const reason = String((typeof idOrPayload === 'object' ? idOrPayload?.reason : maybeReason) || '手動整理').trim();
    if (!mId) return false;
    try {
      if (getApp()?.TrashMedia) await getApp().TrashMedia(mId, reason);
      addToast(`🗑️ メディアをゴミ箱へ移動しました (${reason})`, 'info', 2500); await fetchMedia(); return true;
    } catch (e: any) { addToast(`退避失敗: ${e?.message || e}`, 'error', 3500); return false; }
  };

  const restoreMedia = async (idOrRow: any) => {
    const mId = String(typeof idOrRow === 'object' ? (idOrRow?.mediaId || idOrRow?.media_id || idOrRow?.id) : idOrRow).trim();
    if (!mId) return false;
    try {
      if (getApp()?.RestoreMedia) await getApp().RestoreMedia(mId);
      addToast('♻️ メディアを復元しました', 'success', 2500); await fetchMedia(); return true;
    } catch (e: any) { addToast(`復元失敗: ${e?.message || e}`, 'error', 3500); return false; }
  };

  const updateMediaMetadata = async (mediaId: string, downloadStatus: string, stashSceneId: string, stashImageId: string, failedReason: string) => {
    try {
      if (getApp()?.UpdateMediaMetadata) await getApp().UpdateMediaMetadata(mediaId, downloadStatus, stashSceneId, stashImageId, failedReason);
      addToast(`💾 メタデータを更新しました [${downloadStatus}]`, 'success', 2500); await fetchMedia(); return true;
    } catch (e: any) { addToast(`更新失敗: ${e?.message || e}`, 'error', 4000); return false; }
  };

  const toggleBookmark = async (mId: string) => {
    try {
      const res = (await getApp()?.ToggleMediaBookmark?.(mId)) ?? false;
      const item = mediaResults.value.find((m: any) => (m.media_id || m.id) === mId); if (item) item.is_bookmarked = res;
      addToast(res ? '⭐ お気に入り追加' : '☆ お気に入り解除', 'info', 2000); return res;
    } catch {}
  };

  onMounted(() => { try { EventsOn('media:retried', () => fetchMedia()); } catch (_) {} });
  onUnmounted(() => { try { EventsOff('media:retried'); } catch (_) {} });

  return {
    mediaResults, mediaTotal, isMediaLoading, mediaAccount, mediaPage, mediaLimit, mediaStatusFilter, mediaTypeFilter, mediaStats, downloadStatusStats, errorMessage, isWails,
    fetchMedia, setMediaPage: (p: number) => { mediaPage.value = Math.max(1, p); return fetchMedia(); },
    setMediaLimit: (l: number) => { mediaLimit.value = l; mediaPage.value = 1; return fetchMedia(); },
    setMediaAccount: (acc: string) => { mediaAccount.value = acc; mediaPage.value = 1; return fetchMedia(); },
    setMediaStatusFilter: (st: string) => { mediaStatusFilter.value = st; mediaPage.value = 1; return fetchMedia(); },
    setMediaTypeFilter: (t: 'all' | 'image' | 'video') => { mediaTypeFilter.value = t; mediaPage.value = 1; return fetchMedia(); },
    trashMedia, restoreMedia, updateMediaMetadata, toggleBookmark,
    escalateMediaToThunder: async (m: any) => {
      const mId = String(typeof m === 'object' ? (m.media_id || m.id) : m).trim();
      const url = String(typeof m === 'object' ? (m.download_url || m.url || '') : '').trim();
      try {
        const ok = await getApp()?.EscalateToThunder?.(mId, url);
        if (ok) { addToast(`⚡ Thunder (迅雷) へエスカレーション投入しました`, 'success', 3000); await fetchMedia(); }
        else { addToast('❌ Thunder への投入に失敗しました', 'error', 3000); }
      } catch (e: any) { addToast(`エスカレーション失敗: ${e?.message || e}`, 'error', 3500); }
    },
    retryMedia: async (mId: string) => { try { await getApp()?.RetryMediaDownload?.(mId); addToast(`🔄 再ダウンロードを発行`, 'info', 3000); await fetchMedia(); } catch (e: any) { addToast(`リトライ失敗: ${e?.message || e}`, 'error', 3500); } },
    openInExplorer: async (mId: string) => { try { await getApp()?.OpenInExplorer?.(mId); addToast('📂 フォルダ表示', 'info', 2000); } catch { addToast('フォルダ表示失敗', 'warning', 3000); } },
    openWithDefaultApp: async (mId: string) => { try { await getApp()?.OpenWithDefaultApp?.(mId); addToast('▶️ 既定アプリ起動', 'info', 2000); } catch { addToast('起動失敗', 'warning', 3000); } },
  };
}
