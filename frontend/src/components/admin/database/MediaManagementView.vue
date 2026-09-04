<!-- frontend/src/components/admin/database/MediaManagementView.vue (100行以下 - SPEC-PRINCIPLE-001) -->
<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { BrowserOpenURL } from '../../../../wailsjs/runtime/runtime';
import { useStashResolver } from '../../../composables/useStashResolver';
import MediaCockpitHeader from './MediaCockpitHeader.vue';
import MediaToolbar from './MediaToolbar.vue';
import MediaBatchToolbar from './MediaBatchToolbar.vue';
import MediaPaginationBar from './MediaPaginationBar.vue';
import MediaQueueStatus from './MediaQueueStatus.vue';
import MediaGrid from './MediaGrid.vue';
import MediaTableView from './MediaTableView.vue';
import MediaTrashModal from './MediaTrashModal.vue';
import { useMediaBatchOps } from '../../../composables/admin/useMediaBatchOps';

const props = defineProps<{
  mediaItems: any[]; total: number; accounts: any[]; accountFilter: string; statusFilter: string;
  typeFilter: 'all' | 'image' | 'video'; loading: boolean; page: number; limit: number; config?: any; activeJob?: any;
  stats?: { total_count: number; image_count: number; video_count: number };
  queueStats?: { queued: number; completed: number; dead_404: number; outsourced: number; escalated?: number; retained?: number; failed?: number; total: number };
}>();

const emit = defineEmits<{
  (e: 'fetch'): void; (e: 'update:accountFilter', v: string): void; (e: 'update:statusFilter', v: string): void;
  (e: 'update:typeFilter', v: 'all' | 'image' | 'video'): void; (e: 'update:page', p: number): void; (e: 'update:limit', l: number): void;
  (e: 'retryMedia', mediaId: string): void; (e: 'trashMedia', payload: { mediaId: string; reason: string }): void; (e: 'restoreMedia', id: string): void;
  (e: 'purgeMedia', mediaId: string): void; (e: 'purgeByStatus', status: string): void; (e: 'viewPost', articleId: string): void;
  (e: 'startPoll'): void; (e: 'startEscalate'): void; (e: 'startSmartRecovery'): void;
  (e: 'requeueFailed'): void; (e: 'reconcileStash'): void; (e: 'openExplorer', id: string): void; (e: 'openDefault', id: string): void;
  (e: 'toggleBookmark', id: string): void; (e: 'cancelJob', id: string): void;
  (e: 'batchTrashMedia', mediaIds: string[], reason?: string): void; (e: 'batchRestoreMedia', mediaIds: string[]): void;
  (e: 'batchRevertToQueued', mediaIds: string[]): void;
}>();

const viewMode = ref<'large' | 'compact' | 'table'>('large');
const searchQuery = ref(''), onlyBookmarked = ref(false), showTrashModal = ref(false), selectedMediaForTrash = ref<any>(null);
const { selectedIds, selectedCount, toggleSelect, selectAll, clearSelection } = useMediaBatchOps();

const openTrashModal = (m: any) => { selectedMediaForTrash.value = m; showTrashModal.value = true; };
const { openStashWebUI } = useStashResolver();
const openStash = () => openStashWebUI();

const handleBatchRetry = () => { selectedIds.value.forEach(id => emit('retryMedia', id)); clearSelection(); };
const handleBatchRevertToQueued = () => {
  emit('batchRevertToQueued', Array.from(selectedIds.value));
  clearSelection();
};
const handleBatchTrash = () => { emit('batchTrashMedia', Array.from(selectedIds.value), '一括退避'); clearSelection(); };
const handleBatchRestore = () => { emit('batchRestoreMedia', Array.from(selectedIds.value)); clearSelection(); };

onMounted(() => emit('fetch'));
</script>

<template>
  <div class="h-full flex flex-col min-h-0 space-y-2 overflow-hidden bg-slate-950 p-3">
    <MediaCockpitHeader :active-job="activeJob" @cancel-job="emit('cancelJob', $event)" />
    <MediaToolbar
      v-model:search-query="searchQuery" :account-filter="accountFilter" :status-filter="statusFilter"
      :type-filter="typeFilter" v-model:view-mode="viewMode" v-model:only-bookmarked="onlyBookmarked"
      :accounts="accounts" :stats="stats" :active-job="activeJob"
      @update:account-filter="emit('update:accountFilter', $event)"
      @update:status-filter="emit('update:statusFilter', $event)"
      @update:type-filter="emit('update:typeFilter', $event)"
      @open-stash="openStash" @start-smart-recovery="emit('startSmartRecovery')"
      @reconcile-stash="emit('reconcileStash')"
    />
    <MediaBatchToolbar
      :selected-count="selectedCount" :is-trash-view="statusFilter === 'TRASH'" :total="total"
      @batch-retry="handleBatchRetry" @batch-revert-to-queued="handleBatchRevertToQueued" @batch-trash="handleBatchTrash"
      @batch-restore="handleBatchRestore" @select-all="selectAll(mediaItems)" @clear-selection="clearSelection"
    />
    <MediaQueueStatus :stats="queueStats" :active-job="activeJob" :status-filter="statusFilter" @update:status-filter="emit('update:statusFilter', $event)" />
    <div class="flex-1 min-h-0 overflow-y-auto">
      <MediaTableView v-if="viewMode === 'table'" :items="mediaItems" :selected-ids="selectedIds" @select="toggleSelect($event.media_id || $event.id)" @toggle-select="toggleSelect" @toggle-select-all="() => selectedIds.size === mediaItems.length ? clearSelection() : selectAll(mediaItems)" @retry="emit('retryMedia', $event)" @purge="emit('purgeMedia', $event)" @open-explorer="emit('openExplorer', $event)" @open-default="emit('openDefault', $event)" @toggle-bookmark="emit('toggleBookmark', $event)" />
      <MediaGrid v-else :media-items="mediaItems" :selected-ids="selectedIds" :loading="loading" :search-query="searchQuery" :only-bookmarked="onlyBookmarked" @toggle-select="toggleSelect" @retry-media="emit('retryMedia', $event)" @trash-media="openTrashModal" @restore-media="emit('restoreMedia', $event)" @purge-media="emit('purgeMedia', $event)" @toggle-bookmark="emit('toggleBookmark', $event)" @open-explorer="emit('openExplorer', $event)" @open-default="emit('openDefault', $event)" @save-metadata="emit('saveMetadata', $event)" @view-post="emit('viewPost', $event)" @view-post-timeline="emit('viewPostTimeline', $event)" />
    </div>
    <MediaPaginationBar :page="page" :limit="limit" :total="total" @update:page="emit('update:page', $event)" @update:limit="emit('update:limit', $event)" />
    <MediaTrashModal :show="showTrashModal" :media-item="selectedMediaForTrash" @close="showTrashModal = false" @confirm="(p) => emit('trashMedia', p)" />
  </div>
</template>
