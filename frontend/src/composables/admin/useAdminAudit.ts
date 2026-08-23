// frontend/src/composables/admin/useAdminAudit.ts (100行以下)
import { ref } from 'vue';

const getApp = () => (window as any)?.go?.app?.App || (window as any)?.go?.main?.App;

const auditReport = ref<any>(null), isAuditing = ref(false), isPurgingFiles = ref(false);
const isPurgingDB = ref(false), isRollingBack = ref(false), canRollback = ref(false);
const auditStatusMessage = ref<{ success: boolean; message: string } | null>(null);

export function useAdminAudit() {
  const checkCanRollback = async () => {
    try { const app = getApp(); if (app?.CanRollback) canRollback.value = await app.CanRollback(); } catch {}
  };

  const runAudit = async (purgeFiles = false, purgeDB = false) => {
    isAuditing.value = true; auditStatusMessage.value = null;
    try {
      const app = getApp();
      if (app?.RunAudit) {
        const report = await app.RunAudit(purgeFiles, purgeDB);
        auditReport.value = report;
        auditStatusMessage.value = { success: Boolean(report?.integrity_ok && report?.foreign_key_ok), message: report?.summary || '整合性監査完了' };
      }
    } catch (e: any) { auditStatusMessage.value = { success: false, message: `監査失敗: ${e?.message || e}` }; }
    finally { isAuditing.value = false; await checkCanRollback(); }
  };

  const purgeOrphanFiles = async (paths?: string[]) => {
    const raw = paths?.length ? paths : (auditReport.value?.orphan_files || []).map((f: any) => f?.path || f?.Path).filter(Boolean);
    if (!raw?.length) return;
    isPurgingFiles.value = true;
    try {
      const count = await getApp()?.PurgeOrphanFiles?.(raw);
      auditStatusMessage.value = { success: true, message: `${count} 件の孤立ファイルを退避しました` };
      await runAudit();
    } catch (e: any) { auditStatusMessage.value = { success: false, message: `ファイル退避失敗: ${e?.message || e}` }; }
    finally { isPurgingFiles.value = false; }
  };

  const purgeOrphanDBMedia = async (ids?: string[]) => {
    const raw = ids?.length ? ids : (auditReport.value?.orphan_db_media || []).map((m: any) => m?.media_id || m?.MediaID || m?.mediaId).filter(Boolean);
    if (!raw?.length) return;
    isPurgingDB.value = true;
    try {
      const count = await getApp()?.PurgeOrphanDBMedia?.(raw);
      auditStatusMessage.value = { success: true, message: `${count} 件の孤立DBレコードを削除しました` };
      await runAudit();
    } catch (e: any) { auditStatusMessage.value = { success: false, message: `DBパージ失敗: ${e?.message || e}` }; }
    finally { isPurgingDB.value = false; }
  };

  const rollbackLastPurge = async () => {
    isRollingBack.value = true;
    try {
      const count = await getApp()?.RollbackLastPurge?.();
      auditStatusMessage.value = { success: true, message: `↩️ ${count} 件のDBレコードを復元しました` };
      await runAudit();
    } catch (e: any) { auditStatusMessage.value = { success: false, message: `ロールバック失敗: ${e?.message || e}` }; }
    finally { isRollingBack.value = false; }
  };

  return {
    auditReport, isAuditing, isPurgingFiles, isPurgingDB, isRollingBack, canRollback,
    auditStatusMessage, runAudit, purgeOrphanFiles, purgeOrphanDBMedia, rollbackLastPurge, checkCanRollback,
  };
}
