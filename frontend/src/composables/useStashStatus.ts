// frontend/src/composables/useStashStatus.ts (100行以下 - SPEC-PRINCIPLE-001)
import { ref } from 'vue';

const isStashOnline = ref(false);
const stashLatency = ref<number | null>(null);

export function useStashStatus() {
  const isWails = () => typeof window !== 'undefined' && !!((window as any)?.go?.app?.App || (window as any)?.go?.main?.App);

  const checkStashHealth = async () => {
    const start = Date.now();
    try {
      if (isWails()) {
        const getApp = (window as any)?.go?.app?.App || (window as any)?.go?.main?.App;
        if (getApp?.IsStashReady && await getApp.IsStashReady()) {
          isStashOnline.value = true;
          stashLatency.value = Date.now() - start;
          return;
        }
      }
      const r = await fetch('/stash-proxy/', { method: 'HEAD' });
      stashLatency.value = Date.now() - start;
      isStashOnline.value = r.ok || r.status === 401 || r.status === 404;
    } catch {
      isStashOnline.value = false;
      stashLatency.value = null;
    }
  };

  return {
    isStashOnline,
    stashLatency,
    checkStashHealth,
  };
}
