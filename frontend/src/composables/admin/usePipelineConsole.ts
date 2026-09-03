// frontend/src/composables/admin/usePipelineConsole.ts (100行以下 - SPEC-PRINCIPLE-001)
import { ref, onMounted, onUnmounted } from 'vue';
import { GetPipelineOverview, GetPipelineLogs, SyncThunderDownloads, TogglePipelineAutoEngine, IsPipelineAutoEngineRunning } from '../../../wailsjs/go/app/App';

export function usePipelineConsole() {
  const overview = ref<any>(null);
  const logs = ref<any[]>([]);
  const selectedLogStage = ref<string>('all');
  const loading = ref(false);
  const syncing = ref(false);
  const isAutoEngineRunning = ref(true);
  let pollTimer: any = null;

  const isWails = () => typeof (window as any).go !== 'undefined';

  const fetchOverview = async () => {
    try {
      if (isWails()) {
        const res = await GetPipelineOverview();
        if (res) overview.value = res;
        isAutoEngineRunning.value = await IsPipelineAutoEngineRunning();
      } else {
        const res = await fetch('/api/admin/pipeline/overview');
        if (res.ok) {
          const data = await res.json();
          if (data) overview.value = data;
        }
      }
    } catch (_) {}
  };

  const fetchLogs = async (stageOverride?: string) => {
    const stage = stageOverride ?? selectedLogStage.value;
    try {
      if (isWails()) {
        const res = await GetPipelineLogs(stage, 50);
        if (res) logs.value = res;
      } else {
        const res = await fetch(`/api/admin/pipeline/logs?stage=${stage}&limit=50`);
        if (res.ok) {
          const data = await res.json();
          if (data) logs.value = data;
        }
      }
    } catch (_) {}
  };

  const setLogStage = async (stage: string) => {
    selectedLogStage.value = stage;
    await fetchLogs(stage);
  };

  const toggleAutoEngine = async () => {
    try {
      const next = !isAutoEngineRunning.value;
      if (isWails()) {
        isAutoEngineRunning.value = await TogglePipelineAutoEngine(next);
      } else {
        const res = await fetch('/api/admin/pipeline/toggle', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ enable: next })
        });
        if (!res.ok) throw new Error('Failed to toggle pipeline engine');
        const data = await res.json();
        isAutoEngineRunning.value = data.isRunning;
      }
      await fetchOverview();
    } catch (_) {}
  };

  const refreshAll = async () => {
    loading.value = true;
    await Promise.all([fetchOverview(), fetchLogs()]);
    loading.value = false;
  };

  const syncAndReconcile = async () => {
    syncing.value = true;
    try {
      if (isWails()) {
        await SyncThunderDownloads('');
      } else {
        await fetch('/api/admin/pipeline/sync-thunder', { method: 'POST' });
      }
      await refreshAll();
    } catch (_) {}
    finally { syncing.value = false; }
  };

  onMounted(() => {
    refreshAll();
    // 起動時に完全自動運転エンジンを自動着火
    if (isWails()) TogglePipelineAutoEngine(true).then((r) => { isAutoEngineRunning.value = r; }).catch(() => {});
    pollTimer = setInterval(() => {
      fetchOverview();
      fetchLogs();
    }, 3000);
  });

  onUnmounted(() => {
    if (pollTimer) clearInterval(pollTimer);
  });

  return {
    overview,
    logs,
    selectedLogStage,
    loading,
    syncing,
    isAutoEngineRunning,
    toggleAutoEngine,
    fetchOverview,
    fetchLogs,
    setLogStage,
    refreshAll,
    syncAndReconcile,
  };
}
