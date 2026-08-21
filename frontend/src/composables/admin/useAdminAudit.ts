// frontend/src/composables/admin/useAdminAudit.ts (100行以下)
import { ref } from 'vue';

const getApp = () => (window as any)?.go?.main?.App;

export function useAdminAudit() {
  const auditReport = ref<any>(null);
  const isAuditing = ref(false);

  const runAudit = async (purgeFiles = false, purgeDB = false) => {
    isAuditing.value = true;
    const app = getApp();
    if (app?.RunAudit) {
      auditReport.value = await app.RunAudit(purgeFiles, purgeDB);
    }
    isAuditing.value = false;
  };

  const purgeOrphanFiles = async (paths: string[]) => {
    const app = getApp();
    if (app?.PurgeOrphanFiles) {
      await app.PurgeOrphanFiles(paths);
      await runAudit();
    }
  };

  const purgeOrphanDBMedia = async (mediaIds: string[]) => {
    const app = getApp();
    if (app?.PurgeOrphanDBMedia) {
      await app.PurgeOrphanDBMedia(mediaIds);
      await runAudit();
    }
  };

  return { auditReport, isAuditing, runAudit, purgeOrphanFiles, purgeOrphanDBMedia };
}
