// frontend/src/composables/admin/useAdminConfig.ts (100行以下)
import { ref, reactive } from 'vue';

const getApp = () => (window as any)?.go?.main?.App;

export function useAdminConfig() {
  const configForm = reactive<any>({
    system: { language: 'ja', default_framework: 'twitter', env: 'production' },
    storage: { db_path: './archive.db', local_media_dir: './blobs', stash_dir: './stash', dumps_dir: './dumps', stash_enabled: true },
    network: { stash_port: 9999, frontend_port: 5173, internal_bind_address: '127.0.0.1', middleware_port: 5175, public_bind_address: '0.0.0.0' },
    scheduler: { poll_interval_sec: 300, backup_interval_hours: 24, max_backup_generations: 7 },
    broadcast: { enabled: false, allowed_networks: ['192.168.0.0/16', '10.0.0.0/8', '172.16.0.0/12', '127.0.0.1/32'] },
    appearance: { font_family_ja: 'Hiragino Sans, Meiryo, sans-serif', font_family_en: 'Nunito, sans-serif', font_family_zh: 'Microsoft YaHei, SimHei, sans-serif' },
  });

  const isConfigLoading = ref(false);
  const configSaved = ref(false);

  const fetchConfig = async () => {
    isConfigLoading.value = true;
    const app = getApp();
    if (app?.GetConfig) {
      const cfg = await app.GetConfig();
      if (cfg) Object.assign(configForm, cfg);
    }
    isConfigLoading.value = false;
  };

  const saveConfig = async () => {
    isConfigLoading.value = true;
    const app = getApp();
    if (app?.SaveConfig) {
      await app.SaveConfig(configForm);
      configSaved.value = true;
      setTimeout(() => { configSaved.value = false; }, 3000);
    }
    isConfigLoading.value = false;
  };

  return { configForm, isConfigLoading, configSaved, fetchConfig, saveConfig };
}
