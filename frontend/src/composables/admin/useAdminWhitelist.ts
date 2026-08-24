// frontend/src/composables/admin/useAdminWhitelist.ts (100行以下)
import { ref } from 'vue';

const getApp = () => (window as any)?.go?.app?.App || (window as any)?.go?.main?.App;

export function useAdminWhitelist() {
  const whitelists = ref<any[]>([]);
  const isWhitelistLoading = ref(false);

  const fetchWhitelists = async () => {
    isWhitelistLoading.value = true;
    const app = getApp();
    if (app?.GetWhitelists) whitelists.value = await app.GetWhitelists() || [];
    isWhitelistLoading.value = false;
  };

  const addWhitelist = async (type: string, value: string, groupName = '', aliasOf = '') => {
    const app = getApp();
    if (app?.AddWhitelist) {
      await app.AddWhitelist(type, value, groupName, aliasOf);
      await fetchWhitelists();
    }
  };

  const toggleWhitelist = async (id: number) => {
    const app = getApp();
    if (app?.ToggleWhitelist) {
      await app.ToggleWhitelist(id);
      await fetchWhitelists();
    }
  };

  const deleteWhitelist = async (id: number) => {
    const app = getApp();
    if (app?.DeleteWhitelist) {
      await app.DeleteWhitelist(id);
      await fetchWhitelists();
    }
  };

  const updateWhitelist = async (id: number, type: string, value: string, groupName = '', aliasOf = '', isActive = true) => {
    const app = getApp();
    if (app?.UpdateWhitelist) {
      await app.UpdateWhitelist(id, type, value, groupName, aliasOf, isActive);
      await fetchWhitelists();
    }
  };

  return {
    whitelists, isWhitelistLoading,
    fetchWhitelists, addWhitelist, toggleWhitelist, deleteWhitelist, updateWhitelist,
  };
}

