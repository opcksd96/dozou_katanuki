// frontend/src/composables/admin/useDownloaderBatchOps.ts (100行以下 - SPEC-PRINCIPLE-001)
import { ref } from 'vue';
import { useToast } from '../useToast';
import { useDownloaderSelection } from './useDownloaderSelection';

export function useDownloaderBatchOps(onRefresh?: () => Promise<void>) {
  const selection = useDownloaderSelection();
  const isOperating = ref(false);
  const { addToast } = useToast();

  const batchControl = async (action: 'pause' | 'unpause' | 'remove') => {
    if (selection.selectedGids.value.size === 0) return;
    isOperating.value = true;
    try {
      const getApp = () => (window as any)?.go?.app?.App;
      const gids = Array.from(selection.selectedGids.value);
      const count = await getApp()?.BatchControlMotrix?.(action, gids);
      const labels = { pause: '一時停止', unpause: '再開', remove: 'キュー削除' };
      addToast(`✅ ${count} 件のタスクを${labels[action]}しました`, 'success', 3000);
      selection.clearSelection();
      if (onRefresh) await onRefresh();
    } catch {
      addToast('❌ 一括操作に失敗しました', 'error', 3000);
    } finally {
      isOperating.value = false;
    }
  };

  const batchSafePurge = async () => {
    if (selection.selectedGids.value.size === 0) return;
    isOperating.value = true;
    try {
      const getApp = () => (window as any)?.go?.app?.App;
      const gids = Array.from(selection.selectedGids.value);
      const count = await getApp()?.SafePurgeWithBackup?.(gids);
      addToast(`🛡️ ${count} 件を退避保存してキューから削除しました`, 'success', 4000);
      selection.clearSelection();
      if (onRefresh) await onRefresh();
    } catch {
      addToast('❌ 安全パージに失敗しました', 'error', 3000);
    } finally {
      isOperating.value = false;
    }
  };

  return { ...selection, isOperating, batchControl, batchSafePurge };
}
