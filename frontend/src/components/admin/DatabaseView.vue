<!-- frontend/src/components/admin/DatabaseView.vue (100行以下 - SPEC-PRINCIPLE-001) -->
<script setup lang="ts">
import { ref, watch, onMounted, onUnmounted } from 'vue';
import { models } from '../../../wailsjs/go/models';
import { useToast } from '../../composables/useToast';
import AccountManagementView from './database/AccountManagementView.vue';
import PostManagementView from './database/PostManagementView.vue';
import MediaManagementView from './database/MediaManagementView.vue';

const props = defineProps<{ admin: any; view: 'accounts' | 'posts' | 'media' }>();
const emit = defineEmits<{ (e: 'navigate', tab: 'accounts' | 'posts' | 'media'): void; (e: 'jumpToTimelinePost', articleId: string): void; }>();
const { addToast } = useToast(), localSelected = ref<models.RenderTree | null>(null), availableAvatars = ref<string[]>([]);

const loadAvatars = async () => { if (props.admin?.fetchAvailableAvatars) availableAvatars.value = (await props.admin.fetchAvailableAvatars('twitter')) || []; };
const refreshView = (v: 'accounts' | 'posts' | 'media') => {
  props.admin?.clearError?.(); props.admin?.fetchAccounts?.();
  if (v === 'posts') props.admin?.searchArticles?.(); else if (v === 'accounts') loadAvatars(); else if (v === 'media') props.admin?.fetchMedia?.();
};

watch(() => props.view, (v) => refreshView(v), { immediate: true });
const handleKey = (e: KeyboardEvent) => {
  if (props.view !== 'posts' || (!e.ctrlKey && !e.metaKey)) return;
  if (e.key.toLowerCase() === 'z') { e.preventDefault(); if (e.shiftKey) props.admin?.redo?.(); else props.admin?.undo?.(); }
  else if (e.key.toLowerCase() === 'y') { e.preventDefault(); props.admin?.redo?.(); }
};
onMounted(() => { refreshView(props.view); window.addEventListener('keydown', handleKey); });
onUnmounted(() => window.removeEventListener('keydown', handleKey));

const handleAutoTranslate = async (autoSave = false) => {
  if (!localSelected.value) return;
  const aid = localSelected.value.id, res = await props.admin?.autoTranslate?.(aid);
  if (res) { localSelected.value = res; if (autoSave) await handleSaveTranslation(res.content?.ja || '', res.content?.en || '', res.content?.zh || ''); else addToast(`ℹ️ 記事 [${aid}] の翻訳下書きを展開しました`, 'info', 3000); }
};
const handleSaveTranslation = async (ja: string, en: string, zh: string) => {
  if (!localSelected.value) return;
  await props.admin?.saveTranslation?.(localSelected.value.id, ja, en, zh);
  if (localSelected.value.content) { localSelected.value.content.ja = ja; localSelected.value.content.en = en; localSelected.value.content.zh = zh; }
  addToast(`💾 記事 [${localSelected.value.id}] の翻訳データを保存しました`, 'success', 3000);
};
const handleTrashArticle = async (id: string, reason: string) => {
  if (await props.admin?.trashArticle?.(id, 'admin', reason)) { if (localSelected.value?.id === id) localSelected.value = null; addToast(`🗑️ 記事 [${id}] をゴミ箱へ移動しました`, 'info', 3000); }
};
const handleBatchTrash = async (ids: string[], reason: string) => {
  if (await props.admin?.batchTrashArticles?.(ids, 'admin', reason)) { if (ids.includes(localSelected.value?.id || '')) localSelected.value = null; addToast(`🗑️ ${ids.length} 件の記事をゴミ箱へ移動しました`, 'info', 4000); }
};
const handleBatchReset = async (ids: string[]) => {
  if (await props.admin?.batchResetTranslations?.(ids)) addToast(`🔄 ${ids.length} 件の記事の翻訳をリセットしました`, 'info', 4000);
};
const handleSaveAccount = async (p: any) => {
  if (await props.admin?.updateAccount?.(p.numericId, p.displayName, p.username, p.avatarUrl, p.description, p.aliasOf, p.groupName)) { addToast(`💾 アカウント [@${p.username || p.numericId}] を更新しました`, 'success', 3000); await loadAvatars(); }
};
const handleUploadAvatar = async (p: { virtualKey: string; base64Data: string }) => {
  if (await props.admin?.saveAvatarImage?.('twitter', p.virtualKey, p.base64Data)) { addToast(`📥 アバター [${p.virtualKey}.jpg] を登録しました！`, 'success', 4000); await loadAvatars(); await props.admin?.fetchAccounts?.(); }
};
const onNavigatePost = async (id: string) => {
  if (!id) return;
  if (props.admin?.searchQuery?.value !== undefined) props.admin.searchQuery.value = id; else if (props.admin) props.admin.searchQuery = id;
  if (props.admin?.searchAccount?.value !== undefined) props.admin.searchAccount.value = 'all'; else if (props.admin) props.admin.searchAccount = 'all';
  await props.admin?.searchArticles?.(); emit('navigate', 'posts');
};
const onNavigateAccountPosts = (id: string) => { try { props.admin?.showAccountPosts?.(id); } catch (_) {} emit('navigate', 'posts'); };
const onNavigateAccountMedia = (id: string) => { try { props.admin?.showAccountMedia?.(id); } catch (_) {} emit('navigate', 'media'); };
</script>

<template>
  <div class="h-full flex flex-col min-h-0 bg-slate-950 overflow-hidden">
    <div class="flex-1 min-h-0 overflow-hidden flex flex-col">
      <AccountManagementView v-if="view === 'accounts'" :accounts="admin?.accountsList?.value ?? admin?.accountsList ?? []" :selected-detail="admin?.selectedAccountDetail?.value ?? admin?.selectedAccountDetail ?? null" :loading="admin?.isAccountLoading?.value ?? admin?.isAccountLoading ?? false" :available-avatars="availableAvatars" @select-account="admin?.selectAccount" @save-account="handleSaveAccount" @merge-accounts="admin?.mergeAccounts" @upload-avatar="handleUploadAvatar" @toggle-whitelist="(id, w) => admin?.toggleAccountWhitelist?.(id, w)" @view-posts="onNavigateAccountPosts" @view-media="onNavigateAccountMedia" @refresh="() => { admin?.fetchAccounts?.(); loadAvatars(); }" />
      <PostManagementView v-else-if="view === 'posts'" :articles="admin?.searchResults?.value ?? admin?.searchResults ?? []" :total="admin?.totalCount?.value ?? admin?.totalCount ?? 0" :selected-article="localSelected" :accounts="admin?.accountsList?.value ?? admin?.accountsList ?? []" :search-account="admin?.searchAccount?.value ?? admin?.searchAccount ?? 'all'" :search-query="admin?.searchQuery?.value ?? admin?.searchQuery ?? ''" :loading="admin?.isSearchLoading?.value ?? admin?.isSearchLoading ?? false" :saving="false" :translating="admin?.isTranslating?.value ?? admin?.isTranslating ?? false" :active-job="admin?.activeJob?.value ?? admin?.activeJob" :can-undo="admin?.canUndo?.value ?? admin?.canUndo" :can-redo="admin?.canRedo?.value ?? admin?.canRedo" :include-trash="admin?.includeTrash?.value ?? admin?.includeTrash" @update:search-account="(v) => { if (admin?.searchAccount?.value !== undefined) admin.searchAccount.value = v; else if (admin) admin.searchAccount = v; }" @update:search-query="(v) => { if (admin?.searchQuery?.value !== undefined) admin.searchQuery.value = v; else if (admin) admin.searchQuery = v; }" @select-account="(acc) => (admin?.setAccount ? admin.setAccount(acc) : admin?.searchArticles?.())" @search="() => admin?.searchArticles?.()" @fetch-more="() => admin?.fetchMoreArticles?.()" @select="localSelected = $event" @save="handleSaveTranslation" @auto-translate="handleAutoTranslate" @batch-translate="() => admin?.startBatchTranslate?.(admin?.searchAccount?.value ?? admin?.searchAccount ?? 'all', false)" @cancel-job="admin?.cancelJob" @trash="handleTrashArticle" @batch-trash="handleBatchTrash" @batch-reset-translations="handleBatchReset" @undo="() => admin?.undo?.()" @redo="() => admin?.redo?.()" @update:include-trash="(v) => admin?.toggleIncludeTrash?.(v)" />
      <MediaManagementView v-else-if="view === 'media'" :media-items="admin?.mediaResults?.value ?? admin?.mediaResults ?? []" :total="admin?.mediaTotal?.value ?? admin?.mediaTotal ?? 0" :stats="admin?.mediaStats?.value ?? admin?.mediaStats" :queue-stats="admin?.downloadStatusStats?.value ?? admin?.downloadStatusStats" :accounts="admin?.accountsList?.value ?? admin?.accountsList ?? []" :account-filter="admin?.mediaAccount?.value ?? admin?.mediaAccount ?? 'all'" :status-filter="admin?.mediaStatusFilter?.value ?? admin?.mediaStatusFilter ?? 'all'" :type-filter="admin?.mediaTypeFilter?.value ?? admin?.mediaTypeFilter ?? 'all'" :page="admin?.mediaPage?.value ?? admin?.mediaPage ?? 1" :limit="admin?.mediaLimit?.value ?? admin?.mediaLimit ?? 24" :loading="admin?.isMediaLoading?.value ?? admin?.isMediaLoading ?? false" :config="admin?.configForm" :active-job="admin?.activeJob?.value ?? admin?.activeJob" @fetch="admin?.fetchMedia" @update:account-filter="(v) => admin?.setMediaAccount?.(v)" @update:status-filter="(v) => admin?.setMediaStatusFilter?.(v)" @update:type-filter="(v) => admin?.setMediaTypeFilter?.(v)" @update:page="(p) => admin?.setMediaPage?.(p)" @update:limit="(l) => admin?.setMediaLimit?.(l)" @retry-media="admin?.retryMedia" @purge-media="admin?.purgeMedia" @purge-by-status="admin?.purgeMediaByStatus" @save-metadata="(p) => admin?.updateMediaMetadata(p.mediaId, p.downloadStatus, p.stashSceneId, p.stashImageId, p.failedReason)" @view-post="onNavigatePost" @view-post-timeline="emit('jumpToTimelinePost', $event)" @start-smart-recovery="() => admin?.startSmartRecovery?.()" @start-thunder="() => admin?.startThunderEscalate?.()" @reconcile-stash="() => admin?.reconcileStashMedia?.()" @open-explorer="admin?.openInExplorer" @open-default="admin?.openWithDefaultApp" @toggle-bookmark="admin?.toggleBookmark" @cancel-job="admin?.cancelJob" />
    </div>
  </div>
</template>
