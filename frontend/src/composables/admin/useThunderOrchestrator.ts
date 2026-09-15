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

  const registerPendingTasks = async (limit = 3) => {
    try {
      const n = await getApp()?.RegisterThunderTasksFromDB?.(limit, tempDir.value);
      addToast(n > 0 ? `🚀 ${n} 件のDBタスクを迅雷へ投入しました` : 'ℹ️ 投入対象のPENDINGタスクなし', 'info', 3000);
      await fetchStatus();
    } catch (e: any) { addToast(`❌ 投入失敗: ${e}`, 'error', 3000); }
  };

  const reapEmptyFailedTasks = async () => {
    try {
      const n = await getApp()?.ReapFailedEmptyThunderTasks?.();
      addToast(n > 0 ? `🗑️ ${n} 件の0B枯渇タスクをゴミ箱へ移動しました` : 'ℹ️ 対象の0B失敗タスクなし', 'info', 3000);
      await fetchStatus();
    } catch (e: any) { addToast(`❌ ゴミ箱移動失敗: ${e}`, 'error', 3000); }
  };

  const reactivateTasks = async () => {
    try {
      await getApp()?.ReactivateThunderTasks?.(true);
      addToast('♻️ 迅雷・DBタスクを再活性化しました', 'success', 3000);
      await fetchStatus();
    } catch (e: any) { addToast(`❌ 再活性化失敗: ${e}`, 'error', 3000); }
  };

  const launchThunder = async () => {
    try {
      let ok = getApp()?.LaunchThunder ? await getApp().LaunchThunder() : false;
      addToast(ok ? '⚡ 迅雷 (Thunder) をキックしました' : '❌ 迅雷起動に失敗', ok ? 'success' : 'error', 3000);
    } catch { addToast('❌ 迅雷起動リクエスト失敗', 'error', 3000); }
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

  return { status, loading, maxSlots, intervalSec, tempDir, fetchStatus, startOrchestrator, pauseOrchestrator, resumeOrchestrator, stopOrchestrator, registerPendingTasks, reapEmptyFailedTasks, reactivateTasks, launchThunder, syncDownloads };
}
