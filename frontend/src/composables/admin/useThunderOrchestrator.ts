// frontend/src/composables/admin/useThunderOrchestrator.ts (100行以下 - SPEC-PRINCIPLE-001)
import { ref, onMounted, onUnmounted } from 'vue';
import { useToast } from '../useToast';

export function useThunderOrchestrator() {
  const status = ref<any>(null), loading = ref(false), maxSlots = ref(12), intervalSec = ref(5);
  const tempDir = ref<string>('D:\\迅雷下载');
  const { addToast } = useToast();
  let timer: any = null;

  const fetchStatus = async () => {
    try {
      const getApp = () => (window as any)?.go?.app?.App;
      if (getApp()?.GetThunderOrchestratorStatus) status.value = await getApp().GetThunderOrchestratorStatus();
      if (getApp()?.GetConfig) {
        const cfg = await getApp().GetConfig();
        if (cfg?.storage?.thunder_download_dir) tempDir.value = cfg.storage.thunder_download_dir;
      }
    } catch (e) { console.error('fetchStatus failed', e); }
  };

  const startOrchestrator = async () => {
    loading.value = true;
    try {
      const getApp = () => (window as any)?.go?.app?.App;
      const res = await getApp()?.StartThunderOrchestrator?.(maxSlots.value, intervalSec.value);
      if (res) { status.value = res; addToast(`⚡ 迅雷オーケストレーター開始 (${res.total_jobs} ジョブ)`, 'success', 4000); }
    } catch { addToast('❌ 起動失敗', 'error', 3000); }
    finally { loading.value = false; }
  };

  const pauseOrchestrator = async () => {
    try {
      const getApp = () => (window as any)?.go?.app?.App;
      if (await getApp()?.PauseThunderOrchestrator?.()) { addToast('⏸️ 一時停止しました', 'info', 3000); await fetchStatus(); }
    } catch { addToast('❌ 一時停止失敗', 'error', 3000); }
  };

  const resumeOrchestrator = async () => {
    try {
      const getApp = () => (window as any)?.go?.app?.App;
      if (await getApp()?.ResumeThunderOrchestrator?.()) { addToast('▶️ 再開しました', 'success', 3000); await fetchStatus(); }
    } catch { addToast('❌ 再開失敗', 'error', 3000); }
  };

  const stopOrchestrator = async () => {
    try {
      const getApp = () => (window as any)?.go?.app?.App;
      if (await getApp()?.StopThunderOrchestrator?.()) { addToast('🛑 停止しました', 'warning', 3000); await fetchStatus(); }
    } catch { addToast('❌ 停止処理失敗', 'error', 3000); }
  };

  const resetAndRebuildQueue = async (resetVideos = true) => {
    loading.value = true;
    try {
      const getApp = () => (window as any)?.go?.app?.App;
      const res = await getApp()?.ResetAndRebuildThunderQueue?.(resetVideos);
      if (res) { status.value = res; addToast(`🔄 動画ステータスを差し戻し、キューを再構築しました (${res.total_jobs} ジョブ)`, 'success', 4000); }
    } catch { addToast('❌ キュー再構築失敗', 'error', 3000); }
    finally { loading.value = false; }
  };

  const launchThunder = async () => {
    try {
      const getApp = () => (window as any)?.go?.app?.App;
      if (await getApp()?.LaunchThunder?.()) addToast('⚡ Thunder.exe を起動しました', 'success', 3000);
      else addToast('❌ Thunder.exe 起動失敗', 'error', 3000);
    } catch { addToast('❌ 起動エラー', 'error', 3000); }
  };

  const syncDownloads = async () => {
    loading.value = true;
    try {
      const getApp = () => (window as any)?.go?.app?.App;
      const count = await getApp()?.SyncThunderDownloads?.(tempDir.value);
      if (count > 0) addToast(`📦 ${count} 件の完了ファイルをアカウントフォルダへ移動・同期しました！`, 'success', 4000);
      else addToast('ℹ️ 新規に完了したファイルはありませんでした', 'info', 3000);
      await fetchStatus();
    } catch { addToast('❌ ファイル同期エラー', 'error', 3000); }
    finally { loading.value = false; }
  };

  const startPolling = (ms = 1500) => { fetchStatus(); timer = setInterval(fetchStatus, ms); };
  const stopPolling = () => { if (timer) { clearInterval(timer); timer = null; } };

  onMounted(() => startPolling());
  onUnmounted(() => stopPolling());

  return { status, loading, maxSlots, intervalSec, tempDir, fetchStatus, startOrchestrator, pauseOrchestrator, resumeOrchestrator, stopOrchestrator, resetAndRebuildQueue, launchThunder, syncDownloads };
}
