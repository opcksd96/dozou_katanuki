<!-- frontend/src/components/admin/database/MediaManagementView.vue (100行以下) -->
<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { BrowserOpenURL } from '../../../../wailsjs/runtime/runtime';
import MediaCockpitHeader from './MediaCockpitHeader.vue';
import MediaToolbar from './MediaToolbar.vue';
import MediaPaginationBar from './MediaPaginationBar.vue';
import MediaQueueStatus from './MediaQueueStatus.vue';
import MediaGrid from './MediaGrid.vue';

const props = defineProps<{
  mediaItems: any[];
  total: number;
  stats?: { total_count: number; image_count: number; video_count: number };
  queueStats?: { queued: number; completed: number; dead_404: number; outsourced: number; failed: number; total: number };
  accounts: any[];
  accountFilter: string;
  statusFilter: string;
  typeFilter: 'all' | 'image' | 'video';
  loading: boolean;
  page: number;
  limit: number;
  config?: any;
  activeJob?: any;
}>();

const emit = defineEmits<{
  (e: 'fetch'): void;
  (e: 'update:accountFilter', v: string): void;
  (e: 'update:statusFilter', v: string): void;
  (e: 'update:typeFilter', v: 'all' | 'image' | 'video'): void;
  (e: 'update:page', p: number): void;
  (e: 'update:limit', l: number): void;
  (e: 'retryMedia', mediaId: string): void;
  (e: 'purgeMedia', mediaId: string): void;
  (e: 'purgeByStatus', status: string): void;
  (e: 'viewPost', articleId: string): void;
  (e: 'saveMetadata', payload: any): void;
  (e: 'startDownload'): void;
  (e: 'startPoll'): void;
  (e: 'startEscalate'): void;
  (e: 'requeueFailed'): void;
  (e: 'reconcileStash'): void;
  (e: 'openExplorer', id: string): void;
  (e: 'openDefault', id: string): void;
  (e: 'toggleBookmark', id: string): void;
  (e: 'cancelJob', id: string): void;
}>();

const viewMode = ref<'large' | 'compact' | 'table'>('large');
const searchQuery = ref('');
const onlyBookmarked = ref(false);

const openStash = () => {
  const url = `http://127.0.0.1:${props.config?.network?.stash_port || 9999}`;
  try { BrowserOpenURL(url); } catch { window.open(url, '_blank', 'noopener,noreferrer'); }
};

onMounted(() => { emit('fetch'); });
</script>

<template>
  <div class="h-full flex flex-col min-h-0 space-y-2 overflow-hidden bg-slate-950 p-3">
    <MediaCockpitHeader :active-job="activeJob" @cancel-job="emit('cancelJob', $event)" />
    <MediaToolbar
      v-model:search-query="searchQuery" :account-filter="accountFilter" :status-filter="statusFilter"
      :type-filter="typeFilter" v-model:view-mode="viewMode" v-model:only-bookmarked="onlyBookmarked"
      :accounts="accounts" :stats="stats"
      @update:account-filter="emit('update:accountFilter', $event)"
      @update:status-filter="emit('update:statusFilter', $event)"
      @update:type-filter="emit('update:typeFilter', $event)"
      @open-stash="openStash" @start-download="emit('startDownload')" @start-poll="emit('startPoll')"
      @start-escalate="emit('startEscalate')"
      @requeue-failed="emit('requeueFailed')" @reconcile-stash="emit('reconcileStash')"
    />
    <MediaQueueStatus :stats="queueStats" :active-job="activeJob" :status-filter="statusFilter" @update:status-filter="emit('update:statusFilter', $event)" />
    <MediaGrid :media-items="mediaItems" :loading="loading" :search-query="searchQuery" :only-bookmarked="onlyBookmarked" @retry-media="emit('retryMedia', $event)" @purge-media="emit('purgeMedia', $event)" @toggle-bookmark="emit('toggleBookmark', $event)" />
    <MediaPaginationBar :page="page" :limit="limit" :total="total" @update:page="emit('update:page', $event)" @update:limit="emit('update:limit', $event)" />
  </div>
</template>
