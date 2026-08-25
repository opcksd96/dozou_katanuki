<!-- frontend/src/components/admin/database/MediaManagementView.vue (100行以下) -->
<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { BrowserOpenURL } from '../../../../wailsjs/runtime/runtime';
import MediaCockpitHeader from './MediaCockpitHeader.vue';
import MediaToolbar from './MediaToolbar.vue';
import MediaPaginationBar from './MediaPaginationBar.vue';
import MediaCard from './MediaCard.vue';
import MediaPreviewModal from './MediaPreviewModal.vue';

const props = defineProps<{
  mediaItems: any[];
  total: number;
  stats?: { total_count: number; image_count: number; video_count: number };
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
const selectedIndex = ref<number | null>(null);

const filteredItems = computed(() => {
  let list = props.mediaItems || [];
  if (onlyBookmarked.value) list = list.filter(m => m.is_bookmarked);
  if (searchQuery.value.trim()) {
    const q = searchQuery.value.trim().toLowerCase();
    list = list.filter(m => (m.media_id || m.id || '').toLowerCase().includes(q) || (m.username || '').toLowerCase().includes(q) || (m.full_text || '').toLowerCase().includes(q));
  }
  return list;
});

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
    <div class="flex-1 overflow-y-auto min-h-0">
      <div v-if="loading" class="py-12 text-center text-xs text-slate-500">メディアを読み込み中...</div>
      <div v-else-if="!filteredItems.length" class="py-12 text-center text-xs text-slate-500">該当するメディアが見つかりません</div>
      <div v-else class="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-6 gap-3">
        <MediaCard
          v-for="(item, idx) in filteredItems" :key="item.media_id || item.id" :media="item"
          @click="selectedIndex = idx" @retry="emit('retryMedia', item.media_id || item.id)"
          @purge="emit('purgeMedia', item.media_id || item.id)" @toggle-bookmark="emit('toggleBookmark', item.media_id || item.id)"
        />
      </div>
    </div>
    <MediaPaginationBar :page="page" :limit="limit" :total="total" @update:page="emit('update:page', $event)" @update:limit="emit('update:limit', $event)" />
    <MediaPreviewModal v-if="selectedIndex !== null" :media="filteredItems[selectedIndex]" :has-prev="selectedIndex > 0" :has-next="selectedIndex < filteredItems.length - 1" @close="selectedIndex = null" @prev="selectedIndex--" @next="selectedIndex++" />
  </div>
</template>

