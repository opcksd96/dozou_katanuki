// frontend/src/composables/admin/useAdminDatabaseAccounts.ts (100行以下)
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
      if (app?.ListAllAccounts) {
        accountsList.value = (await app.ListAllAccounts()) || [];
      }
    } catch (e: any) {
      errorMessage.value = `アカウント一覧の取得に失敗: ${e?.message || e}`;
    } finally {
      isAccountLoading.value = false;
    }
  };

  const selectAccount = async (numericId: string) => {
    isAccountLoading.value = true;
    errorMessage.value = null;
    try {
      const app = getApp();
      if (app?.GetAccountDetail) {
        selectedAccountDetail.value = await app.GetAccountDetail(numericId);
      }
    } catch (e: any) {
      errorMessage.value = `アカウント詳細の取得に失敗: ${e?.message || e}`;
    } finally {
      isAccountLoading.value = false;
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
    } finally {
      isAccountLoading.value = false;
    }
  };

  const mergeAccounts = async (sourceNumericId: string, targetNumericId: string): Promise<boolean> => {
    errorMessage.value = null;
    try {
      const app = getApp();
      if (app?.MergeAccounts) {
        await app.MergeAccounts(sourceNumericId, targetNumericId);
        await fetchAccounts();
        if (selectedAccountDetail.value?.account?.numeric_id === sourceNumericId) await selectAccount(sourceNumericId);
        return true;
      }
      return false;
    } catch (e: any) {
      errorMessage.value = `アカウントの合併に失敗: ${e?.message || e}`;
      return false;
    }
  };
  const saveAvatarImage = async (platform: string, virtualKey: string, base64Data: string) => {
    errorMessage.value = null;
    try {
      const app = getApp();
      if (app?.SaveAvatarImage) {
        return await app.SaveAvatarImage(platform, virtualKey, base64Data);
      }
      return null;
    } catch (e: any) {
      errorMessage.value = `アバター画像の保存に失敗: ${e?.message || e}`;
      return null;
    }
  };
  const fetchAvailableAvatars = async (platform = 'twitter'): Promise<string[]> => {
    try {
      const app = getApp();
      return app?.ListAvailableAvatars ? (await app.ListAvailableAvatars(platform)) || [] : [];
    } catch {
      return [];
    }
  };

  return {
    accountsList, selectedAccountDetail, isAccountLoading, errorMessage,
    fetchAccounts, selectAccount, updateAccount, mergeAccounts, saveAvatarImage, fetchAvailableAvatars,
  };
}
