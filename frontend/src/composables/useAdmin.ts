import { ref, reactive, onMounted, onUnmounted } from 'vue';
import { models } from '../../wailsjs/go/models';
import * as WailsApp from '../../wailsjs/go/main/App';
import { EventsOn, EventsOff } from '../../wailsjs/runtime/runtime';

// Wails RPC の呼び出しラッパー（キャッシュ差異やモジュール解決の安全化）
const getRpc = () => (window as any)?.go?.main?.App || WailsApp;

const callGetConfig = async (): Promise<any> => {
  const rpc = getRpc();
  if (typeof rpc.GetConfig === 'function') return rpc.GetConfig();
  throw new Error('GetConfig RPC is not available yet');
};

const callSaveConfig = async (cfg: any): Promise<any> => {
  const rpc = getRpc();
  if (typeof rpc.SaveConfig === 'function') return rpc.SaveConfig(cfg);
  throw new Error('SaveConfig RPC is not available yet');
};

const callGetActiveJob = async (): Promise<any> => {
  const rpc = getRpc();
  if (typeof rpc.GetActiveJob === 'function') return rpc.GetActiveJob();
  return null;
};

const callListJobs = async (): Promise<any> => {
  const rpc = getRpc();
  if (typeof rpc.ListJobs === 'function') return rpc.ListJobs();
  return [];
};

const callStartSalvageJob = async (platform: string, account: string, limit: number): Promise<any> => {
  const rpc = getRpc();
  if (typeof rpc.StartSalvageJob === 'function') return rpc.StartSalvageJob(platform, account, limit);
  throw new Error('StartSalvageJob RPC is not available');
};

const callStartManualImportJob = async (warcPath: string, offline: boolean): Promise<any> => {
  const rpc = getRpc();
  if (typeof rpc.StartManualImportJob === 'function') return rpc.StartManualImportJob(warcPath, offline);
  throw new Error('StartManualImportJob RPC is not available');
};

const callCancelJob = async (jobId: string): Promise<any> => {
  const rpc = getRpc();
  if (typeof rpc.CancelJob === 'function') return rpc.CancelJob(jobId);
  throw new Error('CancelJob RPC is not available');
};

export function useAdmin() {
  const activeJob = ref<models.JobProgress | null>(null);
  const jobList = ref<models.JobProgress[]>([]);
  const logs = ref<string[]>([]);
  const config = ref<models.AppConfig | null>(null);

  const loadingJobs = ref(false);
  const loadingConfig = ref(false);
  const savingConfig = ref(false);
  const actionLoading = ref(false);
  const saveStatus = ref<{ success: boolean; message: string } | null>(null);

  const salvageForm = reactive({
    platform: 'twitter',
    account: '',
    limit: 50,
  });

  const importForm = reactive({
    warcPath: '',
    offline: false,
  });

  // ジョブ一覧および実行中ジョブの取得
  const fetchJobs = async () => {
    loadingJobs.value = true;
    try {
      const [active, list] = await Promise.all([
        callGetActiveJob().catch(() => null),
        callListJobs().catch(() => []),
      ]);
      activeJob.value = active && active.id ? active : null;
      jobList.value = (list || []).sort((a, b) => {
        const timeA = a.started_at ? new Date(a.started_at).getTime() : 0;
        const timeB = b.started_at ? new Date(b.started_at).getTime() : 0;
        return timeB - timeA;
      });

      if (activeJob.value && activeJob.value.logs) {
        logs.value = [...activeJob.value.logs];
      }
    } catch (err: any) {
      console.error('[useAdmin] Failed to fetch jobs:', err);
    } finally {
      loadingJobs.value = false;
    }
  };

  // 設定のロード
  const loadConfig = async () => {
    loadingConfig.value = true;
    saveStatus.value = null;
    try {
      const res = await callGetConfig();
      config.value = models.AppConfig.createFrom(res);
    } catch (err: any) {
      console.error('[useAdmin] Failed to load config:', err);
      saveStatus.value = {
        success: false,
        message: `設定の読み込みに失敗しました: ${err?.message || err}`,
      };
    } finally {
      loadingConfig.value = false;
    }
  };

  // 設定の保存
  const saveConfig = async () => {
    if (!config.value) return;
    savingConfig.value = true;
    saveStatus.value = null;
    try {
      await callSaveConfig(config.value);
      saveStatus.value = {
        success: true,
        message: '設定を正常に保存しました (config.json)',
      };
      setTimeout(() => {
        if (saveStatus.value?.success) {
          saveStatus.value = null;
        }
      }, 4000);
    } catch (err: any) {
      console.error('[useAdmin] Failed to save config:', err);
      saveStatus.value = {
        success: false,
        message: `保存に失敗しました: ${err?.message || err}`,
      };
    } finally {
      savingConfig.value = false;
    }
  };

  // サルベージジョブ開始
  const startSalvage = async () => {
    if (!salvageForm.account.trim()) {
      alert('アカウント名/ID を入力してください');
      return;
    }
    actionLoading.value = true;
    logs.value = [];
    try {
      const job = await callStartSalvageJob(
        salvageForm.platform,
        salvageForm.account.trim(),
        Number(salvageForm.limit) || 50
      );
      activeJob.value = job;
      if (job && job.logs) {
        logs.value = [...job.logs];
      }
      await fetchJobs();
    } catch (err: any) {
      alert(`サルベージジョブの開始に失敗しました: ${err?.message || err}`);
    } finally {
      actionLoading.value = false;
    }
  };

  // 手動インポート開始
  const startImport = async () => {
    if (!importForm.warcPath.trim()) {
      alert('WARC ファイルパスを入力してください');
      return;
    }
    actionLoading.value = true;
    logs.value = [];
    try {
      const job = await callStartManualImportJob(
        importForm.warcPath.trim(),
        importForm.offline
      );
      activeJob.value = job;
      await fetchJobs();
    } catch (err: any) {
      alert(`インポートジョブの開始に失敗しました: ${err?.message || err}`);
    } finally {
      actionLoading.value = false;
    }
  };

  // ジョブキャンセル
  const cancelJob = async (jobId: string) => {
    if (!confirm(`ジョブ [${jobId}] を中止しますか？`)) return;
    try {
      await callCancelJob(jobId);
      await fetchJobs();
    } catch (err: any) {
      alert(`ジョブのキャンセルに失敗しました: ${err?.message || err}`);
    }
  };

  // ログクリア
  const clearLogs = () => {
    logs.value = [];
  };

  // イベント購読ハンドラ
  const handleJobProgress = (progress: models.JobProgress) => {
    if (activeJob.value && activeJob.value.id === progress.id) {
      activeJob.value = progress;
    } else if (progress.status === 'running') {
      activeJob.value = progress;
    }
    if (progress.logs && progress.logs.length > logs.value.length) {
      logs.value = [...progress.logs];
    }
    // 履歴リストの該当アイテムも更新
    const idx = jobList.value.findIndex((j) => j.id === progress.id);
    if (idx !== -1) {
      jobList.value[idx] = progress;
    } else {
      jobList.value.unshift(progress);
    }
  };

  const handleJobFinished = (progress: models.JobProgress) => {
    handleJobProgress(progress);
    if (activeJob.value && activeJob.value.id === progress.id) {
      activeJob.value = null;
    }
    fetchJobs();
  };

  const handleJobLog = (logLine: string) => {
    logs.value.push(logLine);
    if (logs.value.length > 1000) {
      logs.value.shift();
    }
  };

  // ライフサイクル管理
  let isMounted = false;
  const setupEventListeners = () => {
    if (isMounted) return;
    isMounted = true;
    EventsOn('job:progress', handleJobProgress);
    EventsOn('job:started', handleJobProgress);
    EventsOn('job:finished', handleJobFinished);
    EventsOn('job:log', handleJobLog);
  };

  const cleanupEventListeners = () => {
    isMounted = false;
    EventsOff('job:progress');
    EventsOff('job:started');
    EventsOff('job:finished');
    EventsOff('job:log');
  };

  return {
    activeJob,
    jobList,
    logs,
    config,
    loadingJobs,
    loadingConfig,
    savingConfig,
    actionLoading,
    saveStatus,
    salvageForm,
    importForm,
    fetchJobs,
    loadConfig,
    saveConfig,
    startSalvage,
    startImport,
    cancelJob,
    clearLogs,
    setupEventListeners,
    cleanupEventListeners,
  };
}
