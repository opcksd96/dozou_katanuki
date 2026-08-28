// frontend/src/composables/admin/useDownloaderConsole.ts (100行以下)
import { ref, onMounted, onUnmounted } from 'vue';
import { useToast } from '../useToast';

export function useDownloaderConsole() {
  const status = ref<any>(null);
  const loading = ref(false);
  const { addToast } = useToast();
  let timer: any = null;

  const fetchStatus = async () => {
    try {
      const getApp = () => (window as any)?.go?.app?.App;
      if (getApp()?.GetDownloaderDashboardStatus) {
        status.value = await getApp().GetDownloaderDashboardStatus();
      }
    } catch (e) { console.error('fetchStatus failed', e); }
  };

  const controlMotrix = async (action: 'pause_all' | 'unpause_all' | 'purge_all' | 'safe_limits') => {
    loading.value = true;
    try {
      const getApp = () => (window as any)?.go?.app?.App;
      const ok = await getApp()?.ControlMotrix?.(action);
      if (ok) {
        const msgs = { pause_all: 'Motrixの全タスクを一時停止しました', unpause_all: 'Motrixのタスクを再開しました', purge_all: 'Motrixの完了・失敗履歴をパージしました', safe_limits: 'Motrixに安全2並列・スロットル制限を適用しました' };
        addToast(`✅ ${msgs[action]}`, 'success', 3000);
      } else { addToast('❌ Motrixとの通信に失敗しました (RPC未起動)', 'error', 3000); }
      await fetchStatus();
    } catch { addToast('❌ コマンド実行に失敗しました', 'error', 3000); }
    finally { loading.value = false; }
  };

  const launchThunder = async () => {
    const getApp = () => (window as any)?.go?.app?.App;
    if (await getApp()?.LaunchThunder?.()) addToast('⚡ Thunder.exe を起動しました', 'success', 3000);
    else addToast('❌ Thunder.exe の起動に失敗しました', 'error', 3000);
  };

  const escalateToThunder = async (mediaID: string, downloadURL: string) => {
    try {
      const getApp = () => (window as any)?.go?.app?.App;
      const ok = await getApp()?.EscalateToThunder?.(mediaID, downloadURL);
      if (ok) {
        addToast('⚡ Thunder (迅雷) へエスカレーション投入しました', 'success', 3000);
        await fetchStatus();
      } else { addToast('❌ Thunder への投入に失敗しました', 'error', 3000); }
    } catch { addToast('❌ エスカレーション処理でエラーが発生しました', 'error', 3000); }
  };

  const giveUpRetained = async (mediaID: string) => {
    try {
      const getApp = () => (window as any)?.go?.app?.App;
      const ok = await getApp()?.GiveUpRetainedMedia?.(mediaID);
      if (ok) {
        addToast('🛑 タスクの探索を諦め、ステータスを確定しました', 'info', 3000);
        await fetchStatus();
      } else { addToast('❌ ステータス更新に失敗しました', 'error', 3000); }
    } catch { addToast('❌ 処理でエラーが発生しました', 'error', 3000); }
  };

  const startPolling = (ms = 2500) => {
    fetchStatus();
    timer = setInterval(fetchStatus, ms);
  };
  const stopPolling = () => { if (timer) { clearInterval(timer); timer = null; } };

  onMounted(() => startPolling());
  onUnmounted(() => stopPolling());

  return { status, loading, fetchStatus, controlMotrix, launchThunder, escalateToThunder, giveUpRetained, startPolling, stopPolling };
}
