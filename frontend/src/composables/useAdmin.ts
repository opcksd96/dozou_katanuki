import { ref, reactive, onMounted, onUnmounted } from 'vue';
import { models } from '../../wailsjs/go/models';
import * as WailsApp from '../../wailsjs/go/main/App';
import { EventsOn, EventsOff } from '../../wailsjs/runtime/runtime';

// Wails RPC の安全な呼び出しヘルパー
const getApp = () => (window as any)?.go?.main?.App;

// デフォルト設定モック（ブラウザ単体テスト用）
const defaultMockConfig = {
  system: { language: 'ja', default_framework: 'twitter', env: 'production' },
  storage: { db_path: './archive.db', local_media_dir: './blobs', stash_dir: './stash', dumps_dir: './dumps', stash_enabled: true },
  network: { stash_port: 9999, frontend_port: 5173, internal_bind_address: '127.0.0.1' },
  scheduler: { poll_interval_sec: 60, backup_interval_hours: 24, max_backup_generations: 7 },
  appearance: { font_family_ja: 'Hiragino Sans, Meiryo, sans-serif', font_family_en: 'Nunito, sans-serif', font_family_zh: 'Microsoft YaHei, SimHei, sans-serif' },
};

let localMockConfig = JSON.parse(JSON.stringify(defaultMockConfig));

const callGetConfig = async (): Promise<any> => {
  const app = getApp();
  if (app && typeof app.GetConfig === 'function') {
    return await app.GetConfig();
  }
  return localMockConfig;
};

const callSaveConfig = async (cfg: any): Promise<any> => {
  const app = getApp();
  if (app && typeof app.SaveConfig === 'function') {
    return await app.SaveConfig(cfg);
  }
  localMockConfig = JSON.parse(JSON.stringify(cfg));
  return true;
};

const callGetActiveJob = async (): Promise<any> => {
  const app = getApp();
  if (app && typeof app.GetActiveJob === 'function') {
    return await app.GetActiveJob();
  }
  return null;
};

const callListJobs = async (): Promise<any> => {
  const app = getApp();
  if (app && typeof app.ListJobs === 'function') {
    return await app.ListJobs();
  }
  return [];
};

const callStartSalvageJob = async (platform: string, account: string, limit: number): Promise<any> => {
  const app = getApp();
  if (app && typeof app.StartSalvageJob === 'function') {
    return await app.StartSalvageJob(platform, account, limit);
  }
  return {
    id: `mock-salvage-${Date.now()}`,
    type: 'salvage',
    status: 'running',
    current: 10,
    total: limit,
    percentage: (10 / limit) * 100,
    started_at: new Date().toISOString(),
    logs: ['[Mock] Starting salvage job for ' + account, '[Mock] Fetching tweets...'],
  };
};

const callStartManualImportJob = async (warcPath: string, offline: boolean): Promise<any> => {
  const app = getApp();
  if (app && typeof app.StartManualImportJob === 'function') {
    return await app.StartManualImportJob(warcPath, offline);
  }
  return {
    id: `mock-import-${Date.now()}`,
    type: 'import',
    status: 'running',
    current: 5,
    total: 20,
    percentage: 25.0,
    started_at: new Date().toISOString(),
    logs: ['[Mock] Starting WARC import: ' + warcPath],
  };
};

const callCancelJob = async (jobId: string): Promise<any> => {
  const app = getApp();
  if (app && typeof app.CancelJob === 'function') {
    return await app.CancelJob(jobId);
  }
  return true;
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
    try {
      if (typeof EventsOn === 'function') {
        EventsOn('job:progress', handleJobProgress);
        EventsOn('job:started', handleJobProgress);
        EventsOn('job:finished', handleJobFinished);
        EventsOn('job:log', handleJobLog);
      }
    } catch (_) {}
  };

  const cleanupEventListeners = () => {
    isMounted = false;
    try {
      if (typeof EventsOff === 'function') {
        EventsOff('job:progress');
        EventsOff('job:started');
        EventsOff('job:finished');
        EventsOff('job:log');
      }
    } catch (_) {}
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
