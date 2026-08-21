// frontend/src/composables/admin/useAdminJob.ts (100行以下)
import { ref, onMounted, onUnmounted } from 'vue';
import { EventsOn, EventsOff } from '../../../wailsjs/runtime/runtime';

const getApp = () => (window as any)?.go?.main?.App;

export function useAdminJob() {
  const activeJob = ref<any>(null), jobList = ref<any[]>([]), jobLogs = ref<string[]>([]), isJobRunning = ref(false);

  const fetchActiveJob = async () => {
    const app = getApp();
    if (app?.GetActiveJob) {
      activeJob.value = await app.GetActiveJob();
      isJobRunning.value = activeJob.value?.status === 'RUNNING';
      if (activeJob.value?.logs) jobLogs.value = [...activeJob.value.logs];
    }
  };

  const fetchJobList = async () => {
    const app = getApp();
    if (app?.ListJobs) jobList.value = await app.ListJobs();
  };

  const startSalvage = async (platform: string, account: string, limit: number) => {
    const app = getApp();
    if (app?.StartSalvageJob) {
      activeJob.value = await app.StartSalvageJob(platform, account, limit);
      isJobRunning.value = true;
      jobLogs.value = [`[Job] Enqueued salvage: ${account} (limit: ${limit})`];
    }
  };

  const startManualImport = async (warcPath: string, offline: boolean) => {
    const app = getApp();
    if (app?.StartManualImportJob) {
      activeJob.value = await app.StartManualImportJob(warcPath, offline);
      isJobRunning.value = true;
      jobLogs.value = [`[Job] Enqueued WARC import: ${warcPath}`];
    }
  };

  const triggerRestore = async (dumpsDir: string, resetDB: boolean) => {
    const app = getApp();
    if (app?.TriggerRestore) {
      activeJob.value = await app.TriggerRestore(dumpsDir, resetDB);
      isJobRunning.value = true;
      jobLogs.value = [`[Job] Enqueued Restore from ${dumpsDir || './backups/dumps'}`];
    }
  };

  const cancelJob = async (jobId: string) => {
    const app = getApp();
    if (app?.CancelJob) { await app.CancelJob(jobId); await fetchActiveJob(); }
  };

  const setupEventListeners = () => {
    try {
      EventsOn('job:progress', (data: any) => {
        activeJob.value = data;
        isJobRunning.value = data?.status === 'RUNNING';
        if (data?.logs) jobLogs.value = [...data.logs];
      });
      EventsOn('job:finished', (data: any) => {
        activeJob.value = data; isJobRunning.value = false; fetchJobList();
      });
    } catch (_) {}
  };

  const cleanupEventListeners = () => {
    try { EventsOff('job:progress'); EventsOff('job:finished'); } catch (_) {}
  };

  onMounted(() => { fetchActiveJob(); fetchJobList(); setupEventListeners(); });
  onUnmounted(cleanupEventListeners);

  return {
    activeJob, jobList, jobLogs, isJobRunning,
    fetchActiveJob, fetchJobList, startSalvage, startManualImport, triggerRestore, cancelJob,
    clearLogs: () => { jobLogs.value = []; }, setupEventListeners, cleanupEventListeners,
  };
}
