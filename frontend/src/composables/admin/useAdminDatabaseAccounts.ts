// frontend/src/composables/admin/useAdminDatabaseAccounts.ts (100行以下 - SPEC-PRINCIPLE-001)
import { ref } from 'vue';

const getApp = () => (window as any)?.go?.app?.App || (window as any)?.go?.main?.App;

export function useAdminDatabaseAccounts() {
  const accountsList = ref<any[]>([]);
  const selectedAccountDetail = ref<any>(null);
  const isAccountLoading = ref(false);
  const errorMessage = ref<string | null>(null);

  const fetchAccounts = async () => {
    isAccountLoading.value = true;
    errorMessage.value = null;
    try {
      const app = getApp();
      if (app?.ListAllAccounts) accountsList.value = (await app.ListAllAccounts()) || [];
    } catch (e: any) { errorMessage.value = `アカウント一覧の取得に失敗: ${e?.message || e}`; }
    finally { isAccountLoading.value = false; }
  };

  const selectAccount = async (numericIdOrRow: any) => {
    const rawId = typeof numericIdOrRow === 'object' ? (numericIdOrRow?.numeric_id || numericIdOrRow?.id) : numericIdOrRow;
    const numericId = String(rawId || '').trim();
    if (!numericId) return;

    isAccountLoading.value = true;
    errorMessage.value = null;
    try {
      const app = getApp();
      if (app?.GetAccountDetail) selectedAccountDetail.value = await app.GetAccountDetail(numericId);
    } catch (e: any) { errorMessage.value = `アカウント詳細の取得に失敗: ${e?.message || e}`; }
    finally { isAccountLoading.value = false; }
  };

  const toggleAccountWhitelist = async (numericId: string, isWhitelist: boolean) => {
    errorMessage.value = null;
    try {
      const app = getApp();
      if (app?.ToggleAccountWhitelist) {
        await app.ToggleAccountWhitelist(numericId, isWhitelist);
        const target = accountsList.value.find(a => a.numeric_id === numericId);
        if (target) target.is_whitelist = isWhitelist;
        if (selectedAccountDetail.value?.account?.numeric_id === numericId) {
          selectedAccountDetail.value.account.is_whitelist = isWhitelist;
        }
        return true;
      }
      return false;
    } catch (e: any) {
      errorMessage.value = `Whitelist更新に失敗: ${e?.message || e}`;
      return false;
    }
  };

  const updateAccount = async (numericId: string, displayName: string, username: string, avatarUrl: string, description: string, aliasOf: string, groupName: string) => {
    isAccountLoading.value = true;
    errorMessage.value = null;
    try {
      const app = getApp();
      if (app?.UpdateAccount) {
        await app.UpdateAccount(numericId, displayName, username, avatarUrl, description, aliasOf, groupName);
        await fetchAccounts();
        await selectAccount(numericId);
        return true;
      }
      return false;
    } catch (e: any) {
      errorMessage.value = `アカウント情報の更新に失敗: ${e?.message || e}`;
      return false;
    } finally { isAccountLoading.value = false; }
  };

  const mergeAccounts = async (sourceId: string, targetId: string) => {
    isAccountLoading.value = true;
    errorMessage.value = null;
    try {
      const app = getApp();
      if (app?.MergeAccounts) {
        await app.MergeAccounts(sourceId, targetId);
        await fetchAccounts();
        await selectAccount(targetId);
        return true;
      }
      return false;
    } catch (e: any) {
      errorMessage.value = `アカウント統合に失敗: ${e?.message || e}`;
      return false;
    } finally { isAccountLoading.value = false; }
  };

  const saveAvatarImage = async (platform: string, virtualKey: string, base64Data: string) => {
    try {
      const app = getApp();
      if (app?.SaveCustomAvatar) {
        await app.SaveCustomAvatar(platform, virtualKey, base64Data);
        return true;
      }
      return false;
    } catch (e: any) {
      errorMessage.value = `アバター画像の保存に失敗: ${e?.message || e}`;
      return false;
    }
  };

  const fetchAvailableAvatars = async (platform = 'twitter') => {
    try { const app = getApp(); return (await app?.ListCustomAvatars?.(platform)) || []; } catch { return []; }
  };

  return {
    accountsList, selectedAccountDetail, isAccountLoading, errorMessage,
    fetchAccounts, selectAccount, toggleAccountWhitelist, updateAccount, mergeAccounts, saveAvatarImage, fetchAvailableAvatars,
  };
}
