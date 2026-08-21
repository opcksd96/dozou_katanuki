// frontend/src/composables/admin/useAdminBroadcast.ts (100行以下)
import { ref } from 'vue';

const getApp = () => (window as any)?.go?.main?.App;

export function useAdminBroadcast() {
  const broadcastStatus = ref<any>(null);
  const isBroadcastLoading = ref(false);

  const fetchBroadcastStatus = async () => {
    isBroadcastLoading.value = true;
    const app = getApp();
    if (app?.GetBroadcastStatus) {
      broadcastStatus.value = await app.GetBroadcastStatus();
    }
    isBroadcastLoading.value = false;
  };

  const toggleBroadcast = async (enabled: boolean) => {
    isBroadcastLoading.value = true;
    const app = getApp();
    if (app?.ToggleBroadcast) {
      broadcastStatus.value = await app.ToggleBroadcast(enabled);
    }
    isBroadcastLoading.value = false;
  };

  return { broadcastStatus, isBroadcastLoading, fetchBroadcastStatus, toggleBroadcast };
}
