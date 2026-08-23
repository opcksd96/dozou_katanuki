<!-- frontend/src/components/admin/DatabaseView.vue (100行以下) -->
<script setup lang="ts">
import { ref, onMounted, watch } from 'vue';
import { models } from '../../../wailsjs/go/models';
import { useToast } from '../../composables/useToast';
import AccountManagementView from './database/AccountManagementView.vue';
import PostManagementView from './database/PostManagementView.vue';
import MediaManagementView from './database/MediaManagementView.vue';
import WhitelistManagementView from './database/WhitelistManagementView.vue';

const props = defineProps<{ admin: any }>();
const { addToast } = useToast();
const localSelected = ref<models.RenderTree | null>(null);
const availableAvatars = ref<string[]>([]);

const loadAvailableAvatars = async () => {
  if (props.admin?.fetchAvailableAvatars) {
    availableAvatars.value = await props.admin.fetchAvailableAvatars('twitter');
  }
};

const onTabChange = (tab: 'accounts' | 'posts' | 'media' | 'whitelist') => {
  props.admin.activeSubTab.value = tab; props.admin.clearError?.();
  if (tab === 'posts') props.admin.searchArticles?.();
  else if (tab === 'accounts') { props.admin.fetchAccounts?.(); loadAvailableAvatars(); }
  else if (tab === 'media') props.admin.fetchMedia?.();
  else if (tab === 'whitelist') props.admin.fetchWhitelists?.();
};

const handleAutoTranslate = async (autoSave = false) => {
  if (!localSelected.value) return;
  const aid = localSelected.value.id, res = await props.admin.autoTranslate(aid);
  if (res) {
    localSelected.value = res;
    if (autoSave) await handleSaveTranslation(res.content?.ja || '', res.content?.en || '', res.content?.zh || '');
    else addToast(`ℹ️ 記事 [${aid}] の翻訳下書きを展開しました（未保存）`, 'info', 3000);
  }
};

const handleSaveTranslation = async (ja: string, en: string, zh: string) => {
  if (!localSelected.value) return;
  const aid = localSelected.value.id;
  await props.admin.saveTranslation(aid, ja, en, zh);
  if (localSelected.value.content) {
    localSelected.value.content.ja = ja; localSelected.value.content.en = en; localSelected.value.content.zh = zh;
  }
  addToast(`💾 記事 [${aid}] の翻訳データを保存しました`, 'success', 3000);
};

const handleSaveAccount = async (payload: { numericId: string; displayName: string; username: string; avatarUrl: string; description: string }) => {
  const success = await props.admin.updateAccount(payload.numericId, payload.displayName, payload.username, payload.avatarUrl, payload.description);
  if (success) {
    addToast(`💾 アカウント [@${payload.username || payload.numericId}] の情報を更新しました`, 'success', 3000);
    await loadAvailableAvatars();
  }
};

const handleUploadAvatar = async (payload: { virtualKey: string; base64Data: string }) => {
  const savedPath = await props.admin.saveAvatarImage('twitter', payload.virtualKey, payload.base64Data);
  if (savedPath) {
    addToast(`📥 アバター画像 [${payload.virtualKey}.jpg] を assets に登録しました！`, 'success', 4000);
    await loadAvailableAvatars();
    await props.admin.fetchAccounts?.();
    if (props.admin.selectedAccountDetail?.value?.account?.numeric_id) {
      await props.admin.selectAccount(props.admin.selectedAccountDetail.value.account.numeric_id);
    }
  }
};

const handleBatchTranslate = async () => {
  const acc = props.admin.searchAccount.value || 'all';
  addToast(`🚀 一括自動翻訳を開始しました (対象: ${acc === 'all' ? '全件' : '@' + acc})`, 'info', 3000);
  await props.admin.startBatchTranslate(acc, false);
};

const handleSaveMediaMetadata = async (payload: { mediaId: string; downloadStatus: string; stashSceneId: string; stashImageId: string; failedReason: string }) => {
  const ok = await props.admin.updateMediaMetadata?.(payload.mediaId, payload.downloadStatus, payload.stashSceneId, payload.stashImageId, payload.failedReason);
  if (ok) {
    addToast(`💾 メディア [${payload.mediaId}] のメタデータを更新しました`, 'success', 3000);
  }
};

watch(() => props.admin.activeJob.value?.status, (newStatus, oldStatus) => {
  if (props.admin.activeJob.value?.type === 'translate') {
    if (oldStatus === 'RUNNING' && newStatus === 'COMPLETED') {
      addToast('✅ 未翻訳の一括自動翻訳が完了しました！', 'success', 4000);
      props.admin.searchArticles?.();
    } else if (oldStatus === 'RUNNING' && newStatus === 'FAILED') {
      addToast('❌ 一括自動翻訳ジョブが失敗しました', 'error', 4000);
    }
  }
});

watch(() => props.admin.searchResults.value, (newResults) => {
  if (localSelected.value && newResults) {
    const found = newResults.find((a: any) => a.id === localSelected.value?.id);
    if (found) localSelected.value = found;
  }
});

onMounted(() => {
  props.admin?.fetchAccounts?.();
  loadAvailableAvatars();
  if (props.admin?.activeSubTab?.value === 'posts') props.admin?.searchArticles?.();
});
</script>

<template>
  <div class="space-y-1.5 flex flex-col h-full min-h-0">
    <!-- エラー通知バナー -->
    <div v-if="admin.errorMessage.value" class="px-3 py-1.5 bg-rose-950/80 border border-rose-600/80 rounded-lg text-rose-200 text-xs flex items-center justify-between shadow">
      <div class="flex items-center gap-2"><span class="font-bold">⚠️ エラー:</span><span>{{ admin.errorMessage.value }}</span></div>
      <button @click="admin.clearError" class="text-rose-400 hover:text-rose-100 font-bold px-1.5 py-0.5 rounded">✕</button>
    </div>

    <!-- サブタブナビゲーション (薄型コンパクト) -->
    <div class="flex items-center justify-between border-b border-slate-800 pb-1.5 text-xs font-bold shrink-0">
      <div class="flex gap-1.5">
        <button v-for="tab in [
          { id: 'accounts', label: '👤 アカウント (accounts)' }, { id: 'posts', label: '📝 投稿・翻訳 (articles)' },
          { id: 'media', label: '🖼️ メディア (media / Stash)' }, { id: 'whitelist', label: '🛡️ Whitelist' }
        ]" :key="tab.id" @click="onTabChange(tab.id as any)" class="px-2.5 py-1 rounded-lg border transition-colors flex items-center gap-1.5 text-xs" :class="admin.activeSubTab.value === tab.id ? 'bg-blue-600 border-blue-500 text-white shadow' : 'bg-slate-900 border-slate-800 text-slate-400 hover:text-slate-200'">
          {{ tab.label }}
        </button>
      </div>
    </div>

    <!-- 各ドメイン別ビュー -->
    <div class="flex-1 min-h-0 overflow-hidden flex flex-col">
      <AccountManagementView
        v-if="admin.activeSubTab.value === 'accounts'"
        :accounts="admin.accountsList.value"
        :selected-detail="admin.selectedAccountDetail.value"
        :loading="admin.isAccountLoading.value"
        :available-avatars="availableAvatars"
        @select-account="(id) => admin.selectAccount(id)"
        @save-account="handleSaveAccount"
        @upload-avatar="handleUploadAvatar"
        @view-posts="(id) => admin.showAccountPosts(id)"
        @view-media="(id) => admin.showAccountMedia(id)"
        @refresh="() => { admin.fetchAccounts(); loadAvailableAvatars(); }"
      />
      <PostManagementView v-else-if="admin.activeSubTab.value === 'posts'" :articles="admin.searchResults.value" :total="admin.totalCount.value" :selected-article="localSelected" :accounts="admin.accountsList.value" :search-account="admin.searchAccount.value" :search-query="admin.searchQuery.value" :page="admin.page.value" :limit="admin.limit.value" :loading="admin.isSearchLoading.value" :saving="false" :translating="admin.isTranslating.value" :active-job="admin.activeJob.value" @update:search-account="(v) => admin.searchAccount.value = v" @update:search-query="(v) => admin.searchQuery.value = v" @update:page="(p) => admin.page.value = p" @search="admin.searchArticles" @select="(art) => localSelected = art" @save="handleSaveTranslation" @auto-translate="handleAutoTranslate" @batch-translate="handleBatchTranslate" @cancel-job="(id) => admin.cancelJob(id)" />
      <MediaManagementView
        v-else-if="admin.activeSubTab.value === 'media'"
        :media-items="admin.mediaResults?.value || []"
        :total="admin.mediaTotal?.value || 0"
        :stats="admin.mediaStats?.value || { total_count: 0, image_count: 0, video_count: 0 }"
        :accounts="admin.accountsList?.value || []"
        :account-filter="admin.searchAccount?.value || 'all'"
        :status-filter="admin.mediaStatusFilter?.value || 'all'"
        :type-filter="admin.mediaTypeFilter?.value || 'all'"
        :page="admin.page?.value || 1"
        :limit="admin.limit?.value || 24"
        :loading="admin.isMediaLoading?.value || false"
        :config="admin.configForm"
        :active-job="admin.activeJob?.value"
        @fetch="admin.fetchMedia"
        @save-metadata="handleSaveMediaMetadata"
        @update:account-filter="(v) => { if (admin.searchAccount) admin.searchAccount.value = v; }"
        @update:status-filter="(v) => { if (admin.mediaStatusFilter) admin.mediaStatusFilter.value = v; }"
        @update:type-filter="(v) => { if (admin.mediaTypeFilter) admin.mediaTypeFilter.value = v; }"
        @update:page="(p) => { if (admin.page) admin.page.value = p; }"
        @update:limit="(l) => { if (admin.limit) admin.limit.value = l; }"
        @retry-media="(id) => admin.retryMedia?.(id)"
        @purge-media="async (id) => { const ok = await admin.purgeMedia?.(id); if (ok) addToast(`🗑️ メディア [${id}] をDBからパージしました`, 'info', 3000); }"
        @purge-by-status="async (st) => { const cnt = await admin.purgeMediaByStatus?.(st); addToast(`🗑️ [${st}] のメディア ${cnt} 件をパージしました`, 'info', 4000); }"
        @view-post="(artId) => { if (admin.searchQuery) admin.searchQuery.value = artId; if (admin.page) admin.page.value = 1; admin.activeSubTab.value = 'posts'; admin.searchArticles?.(); }"
        @start-download="async () => { await admin.startMediaDownload?.(); addToast('🚀 メディアダウンロードジョブを開始しました！', 'info', 3000); }"
        @start-poll="async () => { await admin.startMediaPoll?.(); addToast('🔄 Aria2 委託回収ジョブを開始しました！', 'info', 3000); }"
        @reconcile-stash="async () => { const cnt = await admin.reconcileStashMedia?.(); addToast(`🎛️ Stash連携同期完了: ${cnt} 件のメディアを自動紐付けしました！`, 'success', 4000); }"
        @requeue-failed="async () => { const cnt = await admin.requeueMedia?.('DEAD_404'); addToast(`🔁 失敗メディア ${cnt} 件を QUEUED へ再キューイングしました`, 'success', 4000); }"
        @merge-duplicates="async () => { const cnt = await admin.mergeDuplicates?.(); addToast(`🧬 重複メディア ${cnt} 件を統合・クレンジングしました`, 'success', 4000); }"
        @purge-low-res="async () => { const cnt = await admin.purgeLowResDuplicates?.(); addToast(`🧹 低解像度メディア ${cnt} 件をごみ箱へ退避しました`, 'success', 4000); }"
        @open-explorer="(id) => admin.openInExplorer?.(id)"
        @open-default="(id) => admin.openWithDefaultApp?.(id)"
        @toggle-bookmark="(id) => admin.toggleBookmark?.(id)"
        @cancel-job="(id) => admin.cancelJob?.(id)"
      />
      <WhitelistManagementView v-else-if="admin.activeSubTab.value === 'whitelist'" :whitelist-list="admin.whitelists.value" :loading="admin.isWhitelistLoading.value" @fetch="admin.fetchWhitelists" @add="(t, v) => admin.addWhitelist(t, v)" @toggle="(id) => admin.toggleWhitelist(id)" @delete="(id) => admin.deleteWhitelist(id)" />
    </div>
  </div>
</template>
