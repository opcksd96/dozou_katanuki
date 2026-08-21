// frontend/src/composables/admin/useAdminAudit.ts (100行以下)
import { ref } from 'vue';

const getApp = () => (window as any)?.go?.main?.App;

// シングルトンステート（コンポーネント間で監査結果・パージ・ロールバック状態を共有）
const auditReport = ref<any>(null);
const isAuditing = ref(false);
const isPurgingFiles = ref(false);
const isPurgingDB = ref(false);
const isRollingBack = ref(false);
const canRollback = ref(false);
const auditStatusMessage = ref<{ success: boolean; message: string } | null>(null);

export function useAdminAudit() {
  const checkCanRollback = async () => {
    try {
      const app = getApp();
      if (app?.CanRollback) canRollback.value = await app.CanRollback();
    } catch (_) {}
  };

  const runAudit = async (purgeFiles = false, purgeDB = false) => {
    isAuditing.value = true;
    auditStatusMessage.value = null;
    try {
      const app = getApp();
      if (app?.RunAudit) {
        const report = await app.RunAudit(purgeFiles, purgeDB);
        auditReport.value = report;
        auditStatusMessage.value = {
          success: Boolean(report?.integrity_ok && report?.foreign_key_ok),
          message: report?.summary || '整合性監査が完了しました',
        };
      }
    } catch (e: any) {
      auditStatusMessage.value = { success: false, message: `監査失敗: ${e?.message || e}` };
    } finally {
      isAuditing.value = false;
      await checkCanRollback();
    }
  };

  const purgeOrphanFiles = async (paths?: string[]) => {
    const rawPaths = paths && paths.length > 0
      ? paths : (auditReport.value?.orphan_files || []).map((f: any) => f?.path || f?.Path);
    const targetPaths = (rawPaths || []).filter(Boolean);
    if (targetPaths.length === 0) return;

    isPurgingFiles.value = true;
    try {
      const app = getApp();
      if (app?.PurgeOrphanFiles) {
        const count = await app.PurgeOrphanFiles(targetPaths);
        auditStatusMessage.value = { success: true, message: `${count} 件の孤立ファイルをゴミ箱へ退避しました` };
        await runAudit(false, false);
      }
    } catch (e: any) {
      auditStatusMessage.value = { success: false, message: `ファイルパージ失敗: ${e?.message || e}` };
    } finally {
      isPurgingFiles.value = false;
    }
  };

  const purgeOrphanDBMedia = async (mediaIds?: string[]) => {
    const rawIds = mediaIds && mediaIds.length > 0
      ? mediaIds : (auditReport.value?.orphan_db_media || []).map((m: any) => m?.media_id || m?.MediaID || m?.mediaId);
    const targetIds = (rawIds || []).filter(Boolean);
    if (targetIds.length === 0) return;

    isPurgingDB.value = true;
    try {
      const app = getApp();
      if (app?.PurgeOrphanDBMedia) {
        const count = await app.PurgeOrphanDBMedia(targetIds);
        auditStatusMessage.value = { success: true, message: `${count} 件の孤立DBレコードを削除しました（退避済み・Undo可能）` };
        await runAudit(false, false);
      }
    } catch (e: any) {
      auditStatusMessage.value = { success: false, message: `DBレコードパージ失敗: ${e?.message || e}` };
    } finally {
      isPurgingDB.value = false;
    }
  };

  const rollbackLastPurge = async () => {
    isRollingBack.value = true;
    try {
      const app = getApp();
      if (app?.RollbackLastPurge) {
        const count = await app.RollbackLastPurge();
        auditStatusMessage.value = { success: true, message: `↩️ ${count} 件のDBレコードを直前のスナップショットから復元しました` };
        await runAudit(false, false);
      }
    } catch (e: any) {
      auditStatusMessage.value = { success: false, message: `ロールバック失敗: ${e?.message || e}` };
    } finally {
      isRollingBack.value = false;
    }
  };

  return {
    auditReport, isAuditing, isPurgingFiles, isPurgingDB, isRollingBack, canRollback,
    auditStatusMessage, runAudit, purgeOrphanFiles, purgeOrphanDBMedia, rollbackLastPurge, checkCanRollback,
  };
}
