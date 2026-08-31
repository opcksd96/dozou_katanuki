// frontend/src/composables/useExternalAppsHealth.ts (100行以下 - SPEC-PRINCIPLE-001)
import { onMounted } from 'vue';
import { useToast } from './useToast';

let hasCheckedHealth = false;

export function useExternalAppsHealth() {
  const { addToast } = useToast();

  const checkExternalApps = async () => {
    if (hasCheckedHealth) return;
    hasCheckedHealth = true;

    try {
      const getApp = () => (window as any)?.go?.app?.App;
      let motrixOk = false;
      let thunderOk = false;

      // 1. Motrix 疎通確認
      if (getApp()?.GetMotrixStatus) {
        const mStatus = await getApp().GetMotrixStatus();
        motrixOk = mStatus?.is_online || mStatus?.connected || false;
      }

      // 2. Thunder CDP 疎通確認
      if (getApp()?.GetThunderCDPStatus) {
        const tStatus = await getApp().GetThunderCDPStatus();
        thunderOk = tStatus?.is_connected || false;
      }

      const motrixLabel = motrixOk ? '🟢 Motrix [稼働中]' : '⚪ Motrix [待機中]';
      const thunderLabel = thunderOk ? '🟢 迅雷CDP [接続中]' : '⚪ 迅雷 [未接続]';

      addToast(`🔌 外部アプリ連携: ${motrixLabel} | ${thunderLabel}`, 'info', 4500);
    } catch (e) {
      console.error('External apps check failed', e);
    }
  };

  onMounted(() => {
    // 起動直後のチラつき防止のため少しディレイを挟んで1回だけ実行
    setTimeout(checkExternalApps, 1200);
  });

  return { checkExternalApps };
}
