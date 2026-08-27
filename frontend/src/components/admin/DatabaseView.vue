<!-- frontend/src/components/admin/DatabaseView.vue (100行以下) -->
<script setup lang="ts">
import { ref, watch, onMounted } from 'vue';
import { models } from '../../../wailsjs/go/models';
import { useToast } from '../../composables/useToast';
import AccountManagementView from './database/AccountManagementView.vue';
import PostManagementView from './database/PostManagementView.vue';
import MediaManagementView from './database/MediaManagementView.vue';

const props = defineProps<{ admin: any; view: 'accounts' | 'posts' | 'media' }>();
const emit = defineEmits<{ (e: 'navigate', tab: 'accounts' | 'posts' | 'media'): void }>();
const { addToast } = useToast();
const localSelected = ref<models.RenderTree | null>(null), availableAvatars = ref<string[]>([]);

const loadAvatars = async () => { if (props.admin?.fetchAvailableAvatars) availableAvatars.value = (await props.admin.fetchAvailableAvatars('twitter')) || []; };
const refreshView = (v: 'accounts' | 'posts' | 'media') => {
  props.admin?.clearError?.();
  if (v === 'posts') props.admin?.searchArticles?.();
  else if (v === 'accounts') { props.admin?.fetchAccounts?.(); loadAvatars(); }
  else if (v === 'media') props.admin?.fetchMedia?.();
};

watch(() => props.view, (v) => refreshView(v), { immediate: true });
onMounted(() => refreshView(props.view));

const handleAutoTranslate = async (autoSave = false) => {
  if (!localSelected.value) return;
  const aid = localSelected.value.id, res = await props.admin?.autoTranslate?.(aid);
  if (res) {
    localSelected.value = res;
    if (autoSave) await handleSaveTranslation(res.content?.ja || '', res.content?.en || '', res.content?.zh || '');
    else addToast(`ℹ️ 記事 [${aid}] の翻訳下書きを展開しました（未保存）`, 'info', 3000);
  }
};
const handleSaveTranslation = async (ja: string, en: string, zh: string) => {
  if (!localSelected.value) return;
  const aid = localSelected.value.id;
  await props.admin?.saveTranslation?.(aid, ja, en, zh);
  if (localSelected.value.content) { localSelected.value.content.ja = ja; localSelected.value.content.en = en; localSelected.value.content.zh = zh; }
  addToast(`💾 記事 [${aid}] の翻訳データを保存しました`, 'success', 3000);
};
const handleSaveAccount = async (p: { numericId: string; displayName: string; username: string; avatarUrl: string; description: string; aliasOf: string; groupName: string }) => {
  if (await props.admin?.updateAccount?.(p.numericId, p.displayName, p.username, p.avatarUrl, p.description, p.aliasOf, p.groupName)) {
    addToast(`💾 アカウント [@${p.username || p.numericId}] の情報を更新しました`, 'success', 3000);
    await loadAvatars();
  }
};
const handleMergeAccounts = async (sourceId: string, targetId: string) => {
  if (await props.admin?.mergeAccounts?.(sourceId, targetId)) {
    addToast(`🔗 アカウント [${sourceId}] を [${targetId}] へ合併しました`, 'success', 3000);
  } else { addToast(`❌ アカウントの合併に失敗しました`, 'error', 3000); }
};
const handleUploadAvatar = async (p: { virtualKey: string; base64Data: string }) => {
  if (await props.admin?.saveAvatarImage?.('twitter', p.virtualKey, p.base64Data)) {
    addToast(`📥 アバター画像 [${p.virtualKey}.jpg] を assets に登録しました！`, 'success', 4000);
    await loadAvatars(); await props.admin?.fetchAccounts?.();
    const curAcc = props.admin?.selectedAccountDetail?.value?.account?.numeric_id ?? props.admin?.selectedAccountDetail?.account?.numeric_id;
    if (curAcc) await props.admin?.selectAccount?.(curAcc);
  }
};
const onNavigatePost = (id: string) => {
  if (props.admin?.searchQuery?.value !== undefined) props.admin.searchQuery.value = id; else if (props.admin) props.admin.searchQuery = id;
  if (props.admin?.searchAccount?.value !== undefined) props.admin.searchAccount.value = 'all'; else if (props.admin) props.admin.searchAccount = 'all';
  emit('navigate', 'posts');
};

const handleSmartRecovery = async () => {
  addToast('⚡ スマート一括回収ジョブを開始しました（直接 ➔ Motrix ➔ Stash）', 'info', 4000);
  await props.admin?.startSmartRecovery?.();
};
const handleThunderEscalate = async () => {
  addToast('⚡ 迅雷 (Thunder.exe) へメディア直リンクを投入しました', 'info', 4000);
  await props.admin?.startThunderEscalate?.();
};
const handleReconcileStash = async () => {
  addToast('📦 Stash との照合同期を開始しました...', 'info', 3000);
  const count = await props.admin?.reconcileStashMedia?.();
  addToast(`✅ Stash同期完了: ${count || 0} 件を照合・バインドしました`, 'success', 3000);
};
</script>

<template>
  <div class="h-full flex flex-col min-h-0 bg-slate-950 overflow-hidden">
    <div class="flex-1 min-h-0 overflow-hidden flex flex-col">
       <AccountManagementView v-if="view === 'accounts'" :accounts="admin?.accountsList?.value ?? admin?.accountsList ?? []" :selected-detail="admin?.selectedAccountDetail?.value ?? admin?.selectedAccountDetail ?? null" :loading="admin?.isAccountLoading?.value ?? admin?.isAccountLoading ?? false" :available-avatars="availableAvatars" @select-account="admin?.selectAccount" @save-account="handleSaveAccount" @merge-accounts="handleMergeAccounts" @upload-avatar="handleUploadAvatar" @view-posts="(id) => { admin?.showAccountPosts(id); emit('navigate', 'posts'); }" @view-media="(id) => { admin?.showAccountMedia(id); emit('navigate', 'media'); }" @refresh="() => { admin?.fetchAccounts?.(); loadAvatars(); }" />
      <PostManagementView v-else-if="view === 'posts'" :articles="admin?.searchResults?.value ?? admin?.searchResults ?? []" :total="admin?.totalCount?.value ?? admin?.totalCount ?? 0" :selected-article="localSelected" :accounts="admin?.accountsList?.value ?? admin?.accountsList ?? []" :search-account="admin?.searchAccount?.value ?? admin?.searchAccount ?? 'all'" :search-query="admin?.searchQuery?.value ?? admin?.searchQuery ?? ''" :page="admin?.page?.value ?? admin?.page ?? 1" :limit="admin?.limit?.value ?? admin?.limit ?? 20" :loading="admin?.isSearchLoading?.value ?? admin?.isSearchLoading ?? false" :saving="false" :translating="admin?.isTranslating?.value ?? admin?.isTranslating ?? false" :active-job="admin?.activeJob?.value ?? admin?.activeJob" @update:search-account="(v) => { if (admin?.searchAccount?.value !== undefined) admin.searchAccount.value = v; else if (admin) admin.searchAccount = v; }" @update:search-query="(v) => { if (admin?.searchQuery?.value !== undefined) admin.searchQuery.value = v; else if (admin) admin.searchQuery = v; }" @update:page="(p) => { if (admin?.page?.value !== undefined) admin.page.value = p; else if (admin) admin.page = p; }" @change-page="(p) => (admin?.setPage ? admin.setPage(p) : admin?.searchArticles?.(p))" @select-account="(acc) => (admin?.setAccount ? admin.setAccount(acc) : admin?.searchArticles?.(1))" @search="() => admin?.searchArticles?.()" @select="localSelected = $event" @save="handleSaveTranslation" @auto-translate="handleAutoTranslate" @batch-translate="() => admin?.startBatchTranslate?.(admin?.searchAccount?.value ?? admin?.searchAccount ?? 'all', false)" @cancel-job="admin?.cancelJob" />
      <MediaManagementView v-else-if="view === 'media'" :media-items="admin?.mediaResults?.value ?? admin?.mediaResults ?? []" :total="admin?.mediaTotal?.value ?? admin?.mediaTotal ?? 0" :stats="admin?.mediaStats?.value ?? admin?.mediaStats" :queue-stats="admin?.downloadStatusStats?.value ?? admin?.downloadStatusStats" :accounts="admin?.accountsList?.value ?? admin?.accountsList ?? []" :account-filter="admin?.mediaAccount?.value ?? admin?.mediaAccount ?? 'all'" :status-filter="admin?.mediaStatusFilter?.value ?? admin?.mediaStatusFilter ?? 'all'" :type-filter="admin?.mediaTypeFilter?.value ?? admin?.mediaTypeFilter ?? 'all'" :page="admin?.mediaPage?.value ?? admin?.mediaPage ?? 1" :limit="admin?.mediaLimit?.value ?? admin?.mediaLimit ?? 24" :loading="admin?.isMediaLoading?.value ?? admin?.isMediaLoading ?? false" :config="admin?.configForm" :active-job="admin?.activeJob?.value ?? admin?.activeJob" @fetch="admin?.fetchMedia" @update:account-filter="(v) => admin?.setMediaAccount ? admin.setMediaAccount(v) : (admin && admin.mediaAccount ? (admin.mediaAccount.value = v) : null, admin?.fetchMedia?.())" @update:status-filter="(v) => admin?.setMediaStatusFilter ? admin.setMediaStatusFilter(v) : (admin && admin.mediaStatusFilter ? (admin.mediaStatusFilter.value = v) : null, admin?.fetchMedia?.())" @update:type-filter="(v) => admin?.setMediaTypeFilter ? admin.setMediaTypeFilter(v) : (admin && admin.mediaTypeFilter ? (admin.mediaTypeFilter.value = v) : null, admin?.fetchMedia?.())" @update:page="(p) => admin?.setMediaPage ? admin.setMediaPage(p) : (admin && admin.mediaPage ? (admin.mediaPage.value = p) : null, admin?.fetchMedia?.())" @update:limit="(l) => admin?.setMediaLimit ? admin.setMediaLimit(l) : (admin && admin.mediaLimit ? (admin.mediaLimit.value = l) : null, admin?.fetchMedia?.())" @retry-media="admin?.retryMedia" @purge-media="admin?.purgeMedia" @purge-by-status="admin?.purgeMediaByStatus" @view-post="onNavigatePost" @start-smart-recovery="handleSmartRecovery" @start-thunder="handleThunderEscalate" @reconcile-stash="handleReconcileStash" @open-explorer="admin?.openInExplorer" @open-default="admin?.openWithDefaultApp" @toggle-bookmark="admin?.toggleBookmark" @cancel-job="admin?.cancelJob" />
    </div>
  </div>
</template>
