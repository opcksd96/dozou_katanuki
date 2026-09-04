import { ref, computed } from 'vue';

// シングルトンとしてポート番号をキャッシュする
const cachedPort = ref<number | null>(null);
const isFetching = ref(false);

const getApp = () => (window as any)?.go?.dozou_katanuki?.app?.App || (window as any)?.go?.main?.App || (window as any)?.go?.app?.App;

export function useStashResolver() {
  const fetchStashPort = async () => {
    if (cachedPort.value !== null || isFetching.value) return;
    isFetching.value = true;
    try {
      const app = getApp();
      if (app?.GetConfig) {
        const cfg = await app.GetConfig();
        if (cfg?.network?.stash_port) {
          cachedPort.value = cfg.network.stash_port;
        }
      } else {
        // Web UI (外部ブラウザ) 向け: API経由で設定を取得
        try {
          const res = await fetch('/api/admin/config/get');
          if (res.ok) {
            const data = await res.json();
            if (data?.network?.stash_port) {
              cachedPort.value = data.network.stash_port;
            }
          }
        } catch (e) {
          // fallback
        }
      }
    } finally {
      if (cachedPort.value === null) {
        cachedPort.value = 9999;
      }
      isFetching.value = false;
    }
  };

  fetchStashPort();

  const stashPort = computed(() => cachedPort.value || 9999);

  const stashHost = computed(() => {
    if (typeof window !== 'undefined' && window.location?.hostname && window.location.hostname !== 'wails.localhost') {
      return window.location.hostname;
    }
    return '127.0.0.1';
  });

  const stashBaseUrl = computed(() => {
    return `http://${stashHost.value}:${stashPort.value}`;
  });

  const getStashSceneUrl = (sceneId: string | number) => {
    if (!sceneId) return '';
    return `${stashBaseUrl.value}/scenes/${sceneId}`;
  };

  const getStashImageUrl = (imageId: string | number) => {
    if (!imageId) return '';
    return `${stashBaseUrl.value}/images/${imageId}`;
  };

  const openStashWebUI = () => {
    const url = `${stashBaseUrl.value}/`;
    try {
      const app = getApp();
      if (app?.BrowserOpenURL) {
         app.BrowserOpenURL(url);
         return;
      }
    } catch(e) {}
    window.open(url, '_blank');
  };

  return {
    stashPort,
    stashHost,
    stashBaseUrl,
    getStashSceneUrl,
    getStashImageUrl,
    openStashWebUI,
    fetchStashPort
  };
}
