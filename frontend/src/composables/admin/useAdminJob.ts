// frontend/src/composables/admin/useAdminJob.ts (100行以下 - SPEC-PRINCIPLE-001)
import { ref, onMounted, onUnmounted } from 'vue';
import { EventsOn, EventsOff } from '../../../wailsjs/runtime/runtime';

const getApp = () => (window as any)?.go?.app?.App || (window as any)?.go?.main?.App;
const isScraperJob = (t?: string) => !t || ['salvage', 'import', 'restore', 'translate', 'job_salvage', 'job_import', 'job_restore'].some(k => t.toLowerCase().includes(k));
const isRunning = (s?: string) => (s || '').toLowerCase() === 'running';

const activeJob = ref<any>(null), jobList = ref<any[]>([]), jobLogs = ref<string[]>([]), isJobRunning = ref(false), backgroundJob = ref<any>(null);
let pollTimer: any = null;

export function useAdminJob() {
  const stopPolling = () => { if (pollTimer) { clearInterval(pollTimer); pollTimer = null; } };
  const startPolling = () => { stopPolling(); pollTimer = setInterval(fetchActiveJob, 1000); };

  const fetchActiveJob = async () => {
    const app = getApp();
    const job = app?.GetActiveJob ? await app.GetActiveJob() : await (await fetch('/api/jobs/status')).json().then((r: any) => Array.isArray(r) ? r.find((j: any) => isRunning(j.status)) || r[0] : r).catch(() => null);
    if (job && isScraperJob(job.type || job.id)) {
      activeJob.value = job; isJobRunning.value = isRunning(job.status);
      if (job.logs?.length) jobLogs.value = [...job.logs];
      if (!isJobRunning.value) stopPolling();
    } else if (!job) { isJobRunning.value = false; stopPolling(); }
  };

  const fetchJobList = async () => { const app = getApp(); if (app?.ListJobs) jobList.value = await app.ListJobs(); };

  const startSalvage = async (platform: string, account: string, limit: number, source: string = 'all') => {
    try {
      const app = getApp();
      activeJob.value = app?.StartSalvageJob ? await app.StartSalvageJob(platform, account, source, limit) : await (await fetch('/api/jobs/salvage', { method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify({ Platform: platform, Account: account, Source: source, Limit: limit }) })).json();
      isJobRunning.value = true; jobLogs.value = [`[Job] 🕷️ スクレイパー採取を開始: ${account} (source: ${source}, limit: ${limit})`];
      startPolling();
    } catch (e: any) { jobLogs.value = [`[Job] ❌ スクレイパー採取の起動に失敗しました: ${e?.message || e}`]; }
  };

  const startManualImport = async (warcPath: string, offline: boolean) => {
    try {
      const app = getApp();
      activeJob.value = app?.StartManualImportJob ? await app.StartManualImportJob(warcPath, offline) : await (await fetch('/api/jobs/import-manual', { method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify({ warc_path: warcPath, offline }) })).json();
      isJobRunning.value = true; jobLogs.value = [`[Job] 📦 WARCインポートを開始: ${warcPath}`]; startPolling();
    } catch (e: any) { jobLogs.value = [`[Job] ❌ WARCインポート起動失敗: ${e?.message || e}`]; }
  };

  const triggerRestore = async (dumpsDir: string, resetDB: boolean) => {
    try {
      const app = getApp();
      activeJob.value = app?.TriggerRestore ? await app.TriggerRestore(dumpsDir, resetDB) : await (await fetch('/api/jobs/restore', { method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify({ dumps_dir: dumpsDir }) })).json();
      isJobRunning.value = true; jobLogs.value = [`[Job] 🔄 原本ダンプからDB再構築を開始: ${dumpsDir || './backups/dumps'}`]; startPolling();
    } catch (e: any) { jobLogs.value = [`[Job] ❌ DB再構築起動失敗: ${e?.message || e}`]; }
  };

  const cancelJob = async (jobId: string) => {
    try { const app = getApp(); if (app?.CancelJob) await app.CancelJob(jobId); else await fetch('/api/jobs/cancel', { method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify({ id: jobId }) }); await fetchActiveJob(); } catch (_) {}
  };

  const setupEventListeners = () => {
    try {
      EventsOn('job:progress', (data: any) => {
        if (isScraperJob(data?.type || data?.id)) { activeJob.value = data; isJobRunning.value = isRunning(data?.status); if (data?.logs?.length) jobLogs.value = [...data.logs]; if (isJobRunning.value && !pollTimer) startPolling(); else if (!isJobRunning.value) stopPolling(); }
        else backgroundJob.value = data;
      });
      EventsOn('job:log', (data: { id: string; line: string }) => { if (isScraperJob(data?.id)) { jobLogs.value.push(data.line); if (jobLogs.value.length > 500) jobLogs.value.shift(); } });
      EventsOn('job:finished', (data: any) => {
        if (isScraperJob(data?.type || data?.id)) { activeJob.value = data; isJobRunning.value = false; if (data?.logs?.length) jobLogs.value = [...data.logs]; stopPolling(); }
        fetchJobList();
      });
    } catch (_) {}
  };

  onMounted(() => { fetchActiveJob(); fetchJobList(); setupEventListeners(); });
  onUnmounted(() => { stopPolling(); try { EventsOff('job:progress'); EventsOff('job:log'); EventsOff('job:finished'); } catch (_) {} });

  return { activeJob, backgroundJob, jobList, jobLogs, isJobRunning, fetchActiveJob, fetchJobList, startSalvage, startManualImport, triggerRestore, cancelJob, clearLogs: () => { jobLogs.value = []; } };
}

