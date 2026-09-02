<!-- frontend/src/components/admin/database/MediaGrid.vue (100行以下 - SPEC-PRINCIPLE-001) -->
<script setup lang="ts">
import { ref, computed } from 'vue';
import MediaCard from './MediaCard.vue';
import MediaPreviewModal from './MediaPreviewModal.vue';

const props = defineProps<{
  mediaItems: any[]; loading: boolean; searchQuery: string; onlyBookmarked: boolean;
  selectedIds?: Set<string>;
}>();

const emit = defineEmits<{
  (e: 'retryMedia', mediaId: string): void; (e: 'trashMedia', m: any): void; (e: 'restoreMedia', id: string): void;
  (e: 'toggleBookmark', id: string): void; (e: 'close'): void; (e: 'saveMetadata', payload: any): void;
  (e: 'viewPost', articleId: string): void; (e: 'viewPostTimeline', articleId: string): void;
  (e: 'fullscreenChange', active: boolean): void; (e: 'openExplorer', mediaId: string): void; (e: 'openDefault', mediaId: string): void;
  (e: 'escalateThunder', m: any): void; (e: 'toggleSelect', id: string): void;
}>();

const selectedIndex = ref<number | null>(null);
const filteredItems = computed(() => {
  let list = props.mediaItems || [];
  if (props.onlyBookmarked) list = list.filter(m => m.is_bookmarked);
  if (props.searchQuery.trim()) {
    const q = props.searchQuery.trim().toLowerCase();
    list = list.filter(m => (m.media_id || m.id || '').toLowerCase().includes(q) || (m.username || '').toLowerCase().includes(q) || (m.full_text || '').toLowerCase().includes(q));
  }
  return list;
});
</script>

<template>
  <div class="flex-1 overflow-y-auto min-h-0">
    <div v-if="loading" class="py-12 text-center text-xs text-slate-500">メディアを読み込み中...</div>
    <div v-else-if="!filteredItems.length" class="py-12 text-center text-xs text-slate-500">該当するメディアが見つかりません</div>
    <div v-else class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-3 xl:grid-cols-4 gap-4 p-1">
      <MediaCard
        v-for="(item, idx) in filteredItems" :key="item.media_id || item.id" :media="item"
        :selected="selectedIds?.has(item.media_id || item.id)"
        @click="selectedIndex = idx" @retry="emit('retryMedia', item.media_id || item.id)"
        @trash="emit('trashMedia', item)" @restore="emit('restoreMedia', item.media_id || item.id)"
        @toggle-bookmark="emit('toggleBookmark', item.media_id || item.id)"
        @toggle-select="emit('toggleSelect', $event)"
        @open-explorer="emit('openExplorer', item.media_id || item.id)" @open-default="emit('openDefault', item.media_id || item.id)"
        @escalate-thunder="emit('escalateThunder', $event)"
        @view-post="emit('viewPost', $event)" @view-post-timeline="emit('viewPostTimeline', $event)"
      />
    </div>
    <MediaPreviewModal v-if="selectedIndex !== null" :media="filteredItems[selectedIndex]" :has-prev="selectedIndex > 0" :has-next="selectedIndex < filteredItems.length - 1" @close="selectedIndex = null" @prev="selectedIndex--" @next="selectedIndex++" @save-metadata="(p) => emit('saveMetadata', p)" @retry="(id) => emit('retryMedia', id)" @escalate-thunder="emit('escalateThunder', $event)" @view-post="(artId) => emit('viewPost', artId)" @view-post-timeline="(artId) => emit('viewPostTimeline', artId)" @fullscreen-change="emit('fullscreenChange', $event)" @open-explorer="emit('openExplorer', $event)" @open-default="emit('openDefault', $event)" />
  </div>
</template>
