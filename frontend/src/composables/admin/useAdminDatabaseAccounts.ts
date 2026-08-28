// frontend/src/composables/admin/useAdminDatabaseAccounts.ts (100行以下 - SPEC-PRINCIPLE-001)
import { ref } from 'vue';
import { useToast } from '../useToast';

const getApp = () => (window as any)?.go?.app?.App || (window as any)?.go?.main?.App;

export function useAdminDatabaseAccounts() {
  const { addToast } = useToast();
  const accountsList = ref<any[]>([]), selectedAccountDetail = ref<any>(null), isAccountLoading = ref(false), errorMessage = ref<string | null>(null), showTrash = ref(false);

  const fetchAccounts = async () => {
    isAccountLoading.value = true; errorMessage.value = null;
    try {
      const app = getApp(), raw = showTrash.value ? ((await app?.ListTrashedAccounts?.()) || []) : (app?.ListAllAccounts ? await app.ListAllAccounts() : await (await fetch('/api/accounts')).json());
      accountsList.value = (raw || []).filter((a: any) => showTrash.value ? !!a?.is_trash : !a?.is_trash).map((a: any) => ({ ...a, username: a.username || a.handle || a.numeric_id || '', display_name: a.display_name || a.username || a.handle || '' }));
    } catch (e: any) { errorMessage.value = `取得失敗: ${e?.message || e}`; }
    finally { isAccountLoading.value = false; }
  };

  const selectAccount = async (numericIdOrRow: any) => {
    const rawId = typeof numericIdOrRow === 'object' ? (numericIdOrRow?.numeric_id || numericIdOrRow?.id) : numericIdOrRow, numericId = String(rawId || '').trim();
    if (!numericId) return;
    isAccountLoading.value = true;
    try {
      const app = getApp();
      if (app?.GetAccountDetail) selectedAccountDetail.value = await app.GetAccountDetail(numericId);
      else { const res = await fetch(`/api/account/detail?id=${encodeURIComponent(numericId)}`); selectedAccountDetail.value = res.ok ? await res.json() : null; }
    } catch (e: any) { errorMessage.value = `詳細取得失敗: ${e?.message || e}`; }
    finally { isAccountLoading.value = false; }
  };

  const toggleAccountWhitelist = async (numericId: string, isWhitelist: boolean) => {
    try {
      if (getApp()?.ToggleAccountWhitelist) {
        await getApp().ToggleAccountWhitelist(numericId, isWhitelist);
        const t = accountsList.value.find(a => a.numeric_id === numericId); if (t) t.is_whitelist = isWhitelist;
        if (selectedAccountDetail.value?.account?.numeric_id === numericId) selectedAccountDetail.value.account.is_whitelist = isWhitelist;
        addToast(isWhitelist ? '🛡️ 巡回対象に設定' : '⚪ 巡回対象から除外', 'info', 2000); return true;
      }
      return false;
    } catch (e: any) { addToast(`更新失敗: ${e?.message || e}`, 'error', 3500); return false; }
  };

  const updateAccount = async (numericId: string, displayName: string, username: string, avatarUrl: string, description: string, aliasOf: string, groupName: string) => {
    isAccountLoading.value = true;
    try {
      if (getApp()?.UpdateAccount) {
        await getApp().UpdateAccount(numericId, displayName, username, avatarUrl, description, aliasOf, groupName);
        addToast('💾 アカウント情報を更新しました', 'success', 2500); await fetchAccounts(); await selectAccount(numericId); return true;
      }
      return false;
    } catch (e: any) { addToast(`更新失敗: ${e?.message || e}`, 'error', 3500); return false; }
    finally { isAccountLoading.value = false; }
  };

  const mergeAccounts = async (sourceId: string, targetId: string) => {
    isAccountLoading.value = true;
    try {
      if (getApp()?.MergeAccounts) {
        await getApp().MergeAccounts(sourceId, targetId);
        addToast('✅ アカウントの名寄せ・統合が完了しました！', 'success', 3500); await fetchAccounts(); await selectAccount(targetId); return true;
      }
      return false;
    } catch (e: any) { addToast(`統合失敗: ${e?.message || e}`, 'error', 4000); return false; }
    finally { isAccountLoading.value = false; }
  };

  const trashAccount = async (idOrPayload: any, maybeReason = '手動整理') => {
    const rawId = typeof idOrPayload === 'object' ? (idOrPayload?.numericId || idOrPayload?.numeric_id || idOrPayload?.id) : idOrPayload;
    const reason = String((typeof idOrPayload === 'object' ? idOrPayload?.reason : maybeReason) || '手動整理').trim(), numericId = String(rawId || '').trim();
    if (!numericId) return false;
    try {
      if (getApp()?.TrashAccount) {
        await getApp().TrashAccount(numericId, reason);
        addToast(`🗑️ アカウントをゴミ箱へ移動しました (${reason})`, 'info', 3000); selectedAccountDetail.value = null; await fetchAccounts(); return true;
      }
      return false;
    } catch (e: any) { addToast(`ゴミ箱移動失敗: ${e?.message || e}`, 'error', 3500); return false; }
  };

  const restoreAccount = async (idOrRow: any) => {
    const rawId = typeof idOrRow === 'object' ? (idOrRow?.numericId || idOrRow?.numeric_id || idOrRow?.id) : idOrRow, numericId = String(rawId || '').trim();
    if (!numericId) return false;
    try {
      if (getApp()?.RestoreAccount) {
        await getApp().RestoreAccount(numericId);
        addToast('♻️ アカウントを復元しました', 'success', 2500); await fetchAccounts(); await selectAccount(numericId); return true;
      }
      return false;
    } catch (e: any) { addToast(`復元失敗: ${e?.message || e}`, 'error', 3500); return false; }
  };

  return {
    accountsList, selectedAccountDetail, isAccountLoading, errorMessage, showTrash,
    toggleShowTrash: async () => { showTrash.value = !showTrash.value; selectedAccountDetail.value = null; await fetchAccounts(); },
    fetchAccounts, selectAccount, toggleAccountWhitelist, updateAccount, mergeAccounts, trashAccount, restoreAccount,
    saveAvatarImage: async (p: string, k: string, b: string) => (await getApp()?.SaveCustomAvatar?.(p, k, b)) ?? false,
    fetchAvailableAvatars: async (p = 'twitter') => (await getApp()?.ListCustomAvatars?.(p)) || [],
  };
}
