<!-- frontend/src/components/admin/database/MediaManagementView.vue (100行以下 - SPEC-PRINCIPLE-001) -->
<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { BrowserOpenURL } from '../../../../wailsjs/runtime/runtime';
import MediaCockpitHeader from './MediaCockpitHeader.vue';
import MediaToolbar from './MediaToolbar.vue';
import MediaPaginationBar from './MediaPaginationBar.vue';
import MediaQueueStatus from './MediaQueueStatus.vue';
import MediaGrid from './MediaGrid.vue';
import MediaTableView from './MediaTableView.vue';
import MediaTrashModal from './MediaTrashModal.vue';

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
  (e: 'viewPostTimeline', articleId: string): void; (e: 'saveMetadata', payload: any): void; (e: 'startDownload'): void;
  (e: 'startPoll'): void; (e: 'startEscalate'): void; (e: 'startSmartRecovery'): void; (e: 'startThunder'): void; (e: 'escalateThunder', m: any): void;
  (e: 'requeueFailed'): void; (e: 'reconcileStash'): void; (e: 'openExplorer', id: string): void; (e: 'openDefault', id: string): void;
  (e: 'toggleBookmark', id: string): void; (e: 'cancelJob', id: string): void;
}>();

const viewMode = ref<'large' | 'compact' | 'table'>('large');
const searchQuery = ref(''), onlyBookmarked = ref(false), showTrashModal = ref(false), selectedMediaForTrash = ref<any>(null);

const openTrashModal = (m: any) => { selectedMediaForTrash.value = m; showTrashModal.value = true; };
const openStash = () => {
  const url = `http://127.0.0.1:${props.config?.network?.stash_port || 9999}`;
  try { BrowserOpenURL(url); } catch { window.open(url, '_blank', 'noopener,noreferrer'); }
};

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
      @open-stash="openStash" @start-smart-recovery="emit('startSmartRecovery')" @start-thunder="emit('startThunder')"
      @reconcile-stash="emit('reconcileStash')"
    />
    <MediaQueueStatus :stats="queueStats" :active-job="activeJob" :status-filter="statusFilter" @update:status-filter="emit('update:statusFilter', $event)" />
    <div class="flex-1 min-h-0 overflow-y-auto">
      <MediaTableView v-if="viewMode === 'table'" :items="mediaItems" @retry="emit('retryMedia', $event)" @purge="emit('purgeMedia', $event)" @open-explorer="emit('openExplorer', $event)" @open-default="emit('openDefault', $event)" @toggle-bookmark="emit('toggleBookmark', $event)" @escalate-thunder="emit('escalateThunder', $event)" />
      <MediaGrid v-else :media-items="mediaItems" :loading="loading" :search-query="searchQuery" :only-bookmarked="onlyBookmarked" @retry-media="emit('retryMedia', $event)" @trash-media="openTrashModal" @restore-media="emit('restoreMedia', $event)" @purge-media="emit('purgeMedia', $event)" @toggle-bookmark="emit('toggleBookmark', $event)" @open-explorer="emit('openExplorer', $event)" @open-default="emit('openDefault', $event)" @save-metadata="emit('saveMetadata', $event)" @escalate-thunder="emit('escalateThunder', $event)" @view-post="emit('viewPost', $event)" @view-post-timeline="emit('viewPostTimeline', $event)" />
    </div>
    <MediaPaginationBar :page="page" :limit="limit" :total="total" @update:page="emit('update:page', $event)" @update:limit="emit('update:limit', $event)" />
    <MediaTrashModal :show="showTrashModal" :media-item="selectedMediaForTrash" @close="showTrashModal = false" @confirm="(p) => emit('trashMedia', p)" />
  </div>
</template>
