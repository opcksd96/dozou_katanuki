// frontend/src/composables/admin/useAdminJob.ts (100行以下 - SPEC-PRINCIPLE-001)
import { ref, onMounted, onUnmounted } from 'vue';
import { EventsOn, EventsOff } from '../../../wailsjs/runtime/runtime';

const getApp = () => (window as any)?.go?.app?.App || (window as any)?.go?.main?.App;
const isScraperJob = (type?: string) => !type || ['salvage', 'import', 'restore', 'translate', 'job_salvage', 'job_import', 'job_restore'].some(k => type.toLowerCase().includes(k));

export function useAdminJob() {
  const activeJob = ref<any>(null), jobList = ref<any[]>([]), jobLogs = ref<string[]>([]), isJobRunning = ref(false);
  const backgroundJob = ref<any>(null);

  const fetchActiveJob = async () => {
    const app = getApp();
    if (app?.GetActiveJob) {
      const job = await app.GetActiveJob();
      if (job && isScraperJob(job.type || job.id)) {
        activeJob.value = job; isJobRunning.value = job.status === 'RUNNING';
        if (job.logs && job.logs.length > 0) jobLogs.value = [...job.logs];
      }
    }
  };

  const fetchJobList = async () => { const app = getApp(); if (app?.ListJobs) jobList.value = await app.ListJobs(); };

  const startSalvage = async (platform: string, account: string, limit: number, source: string = 'all') => {
    const app = getApp();
    if (app?.StartSalvageJob) {
      activeJob.value = await app.StartSalvageJob(platform, account, source, limit);
      isJobRunning.value = true;
      jobLogs.value = [`[Job] 🕷️ スクレイパー採取を開始: ${account} (source: ${source}, limit: ${limit})`];
    }
  };

  const startManualImport = async (warcPath: string, offline: boolean) => {
    const app = getApp();
    if (app?.StartManualImportJob) {
      activeJob.value = await app.StartManualImportJob(warcPath, offline);
      isJobRunning.value = true;
      jobLogs.value = [`[Job] 📦 WARCインポートを開始: ${warcPath}`];
    }
  };

  const triggerRestore = async (dumpsDir: string, resetDB: boolean) => {
    const app = getApp();
    if (app?.TriggerRestore) {
      activeJob.value = await app.TriggerRestore(dumpsDir, resetDB);
      isJobRunning.value = true;
      jobLogs.value = [`[Job] 🔄 原本ダンプからDB再構築を開始: ${dumpsDir || './backups/dumps'}`];
    }
  };

  const cancelJob = async (jobId: string) => {
    const app = getApp();
    if (app?.CancelJob) { await app.CancelJob(jobId); await fetchActiveJob(); }
  };

  const setupEventListeners = () => {
    try {
      EventsOn('job:progress', (data: any) => {
        if (isScraperJob(data?.type || data?.id)) {
          activeJob.value = data; isJobRunning.value = data?.status === 'RUNNING';
          if (data?.logs) jobLogs.value = [...data.logs];
        } else {
          backgroundJob.value = data;
        }
      });
      EventsOn('job:log', (data: { id: string; line: string }) => {
        if (isScraperJob(data?.id)) {
          jobLogs.value.push(data.line);
          if (jobLogs.value.length > 500) jobLogs.value.shift();
        }
      });
      EventsOn('job:finished', (data: any) => {
        if (isScraperJob(data?.type || data?.id)) {
          activeJob.value = data; isJobRunning.value = false;
          if (data?.logs) jobLogs.value = [...data.logs];
        }
        fetchJobList();
      });
    } catch (_) {}
  };

  const cleanupEventListeners = () => {
    try { EventsOff('job:progress'); EventsOff('job:log'); EventsOff('job:finished'); } catch (_) {}
  };

  onMounted(() => { fetchActiveJob(); fetchJobList(); setupEventListeners(); });
  onUnmounted(cleanupEventListeners);

  return {
    activeJob, backgroundJob, jobList, jobLogs, isJobRunning,
    fetchActiveJob, fetchJobList, startSalvage, startManualImport, triggerRestore, cancelJob,
    clearLogs: () => { jobLogs.value = []; }, setupEventListeners, cleanupEventListeners,
  };
}
