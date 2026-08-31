// frontend/src/composables/admin/useThunderOrchestrator.ts (100行以下 - SPEC-PRINCIPLE-001)
import { ref, onMounted, onUnmounted } from 'vue';
import { useToast } from '../useToast';

export function useThunderOrchestrator() {
  const status = ref<any>(null), loading = ref(false), maxSlots = ref(12), intervalSec = ref(3);
  const tempDir = ref<string>('D:\\迅雷下载');
  const { addToast } = useToast();
  let timer: any = null;
  const getApp = () => (window as any)?.go?.app?.App;

  const fetchStatus = async () => {
    try {
      if (getApp()?.GetThunderOrchestratorStatus) status.value = await getApp().GetThunderOrchestratorStatus();
      if (getApp()?.GetConfig) {
        const cfg = await getApp().GetConfig();
        if (cfg?.storage?.thunder_download_dir) tempDir.value = cfg.storage.thunder_download_dir;
      }
    } catch {}
  };

  const startOrchestrator = async () => {
    loading.value = true;
    try {
      const res = await getApp()?.StartThunderOrchestrator?.(maxSlots.value, intervalSec.value);
      if (res) { status.value = res; addToast(`⚡ 迅雷ディスパッチ開始 (${res.total_jobs} ジョブ)`, 'success', 4000); }
    } catch { addToast('❌ 起動失敗', 'error', 3000); }
    finally { loading.value = false; }
  };

  const pauseOrchestrator = async () => {
    if (await getApp()?.PauseThunderOrchestrator?.()) { addToast('⏸️ 一時停止しました', 'info', 3000); await fetchStatus(); }
  };

  const resumeOrchestrator = async () => {
    if (await getApp()?.ResumeThunderOrchestrator?.()) { addToast('▶️ 再開しました', 'success', 3000); await fetchStatus(); }
  };

  const stopOrchestrator = async () => {
    if (await getApp()?.StopThunderOrchestrator?.()) { addToast('🛑 停止しました', 'warning', 3000); await fetchStatus(); }
  };

  const launchThunder = async () => {
    if (await getApp()?.LaunchThunder?.()) addToast('⚡ Thunder.exe を起動しました', 'success', 3000);
    else addToast('❌ 起動失敗', 'error', 3000);
  };

  const syncDownloads = async () => {
    loading.value = true;
    try {
      const count = await getApp()?.SyncThunderDownloads?.(tempDir.value);
      addToast(count > 0 ? `📦 ${count} 件の完了ファイルを同期・重複解消！` : 'ℹ️ 新規完了ファイルなし', count > 0 ? 'success' : 'info', 3000);
      await fetchStatus();
    } catch { addToast('❌ 同期エラー', 'error', 3000); }
    finally { loading.value = false; }
  };

  onMounted(() => { fetchStatus(); timer = setInterval(fetchStatus, 2000); });
  onUnmounted(() => { if (timer) clearInterval(timer); });

  return { status, loading, maxSlots, intervalSec, tempDir, fetchStatus, startOrchestrator, pauseOrchestrator, resumeOrchestrator, stopOrchestrator, launchThunder, syncDownloads };
}
