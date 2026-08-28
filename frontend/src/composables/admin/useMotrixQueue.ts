// frontend/src/composables/admin/useMotrixQueue.ts (100行以下 - SPEC-PRINCIPLE-001)
import { ref, onMounted, onUnmounted } from 'vue';
import { useToast } from '../useToast';

export function useMotrixQueue() {
  const activeTasks = ref<any[]>([]), waitingTasks = ref<any[]>([]), stoppedTasks = ref<any[]>([]);
  const reserves = ref<any[]>([]), globalStat = ref<any>(null), loading = ref(false);
  const activeTab = ref<'active' | 'waiting' | 'stopped' | 'reserves'>('active');
  const { addToast } = useToast();
  let timer: any = null;

  const fetchQueue = async () => {
    try {
      const getApp = () => (window as any)?.go?.app?.App;
      if (getApp()?.FetchMotrixFullQueue) {
        const full = await getApp().FetchMotrixFullQueue();
        activeTasks.value = full?.active_tasks || [];
        waitingTasks.value = full?.waiting_tasks || [];
        stoppedTasks.value = full?.stopped_tasks || [];
      }
      if (getApp()?.GetDownloaderDashboardStatus) {
        globalStat.value = (await getApp().GetDownloaderDashboardStatus())?.motrix || null;
      }
      if (getApp()?.FetchDownloadReserves) {
        reserves.value = (await getApp().FetchDownloadReserves('all')) || [];
      }
    } catch (e) { console.error('fetchQueue failed', e); }
  };

  const controlMotrix = async (action: 'pause_all' | 'unpause_all' | 'purge_all' | 'safe_limits') => {
    loading.value = true;
    try {
      const getApp = () => (window as any)?.go?.app?.App;
      if (await getApp()?.ControlMotrix?.(action)) {
        const msgs = { pause_all: '全タスク一時停止', unpause_all: '全タスク再開', purge_all: '履歴パージ', safe_limits: '安全2並列適用' };
        addToast(`✅ ${msgs[action]}を実行しました`, 'success', 3000);
      } else { addToast('❌ Motrixとの通信に失敗しました', 'error', 3000); }
      await fetchQueue();
    } catch { addToast('❌ コマンド実行に失敗しました', 'error', 3000); }
    finally { loading.value = false; }
  };

  const restoreReserve = async (id: number) => {
    try {
      const getApp = () => (window as any)?.go?.app?.App;
      if (await getApp()?.RestoreReserveToMotrix?.(id)) {
        addToast('🚀 退避タスクを Motrix へ再投入しました', 'success', 3000);
        await fetchQueue();
      } else { addToast('❌ 再投入に失敗しました', 'error', 3000); }
    } catch { addToast('❌ 再投入エラー', 'error', 3000); }
  };

  const startPolling = (ms = 2500) => { fetchQueue(); timer = setInterval(fetchQueue, ms); };
  const stopPolling = () => { if (timer) { clearInterval(timer); timer = null; } };

  onMounted(() => startPolling());
  onUnmounted(() => stopPolling());

  return { activeTasks, waitingTasks, stoppedTasks, reserves, globalStat, loading, activeTab, fetchQueue, controlMotrix, restoreReserve, startPolling, stopPolling };
}
