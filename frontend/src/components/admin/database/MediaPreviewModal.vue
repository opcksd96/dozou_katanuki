<!-- frontend/src/components/admin/database/MediaPreviewModal.vue (100行以下) -->
<script setup lang="ts">
import { computed, onMounted, onUnmounted } from 'vue';
import StashPlayer from '../../media/StashPlayer.vue';
import MediaInspectorPanel from './MediaInspectorPanel.vue';

const props = defineProps<{ media: any }>();
const emit = defineEmits<{
  (e: 'close'): void;
  (e: 'saveMetadata', payload: any): void;
  (e: 'retry', mediaId: string): void;
  (e: 'purge', mediaId: string): void;
  (e: 'viewPost', articleId: string): void;
  (e: 'fullscreenChange', active: boolean): void;
}>();

const isVideo = computed(() => {
  const t = props.media.type?.toLowerCase();
  return t === 'video' || t === 'gif' || t === 'animated_gif' || !!props.media.stash_scene_id;
});

const handleKeyDown = (e: KeyboardEvent) => {
  if (e.key === 'Escape') { e.preventDefault(); emit('close'); }
};

const handleFullscreenChange = (active: boolean) => {
  if (!active) emit('close');
};

onMounted(() => window.addEventListener('keydown', handleKeyDown));
onUnmounted(() => window.removeEventListener('keydown', handleKeyDown));
</script>

<template>
  <div class="fixed inset-0 z-[60] flex items-center justify-center bg-black/85 backdrop-blur-md p-3 md:p-5 select-none" @click.self="emit('close')">
    <!-- モーダルコンテナ (大画面・2ペイン分割) -->
    <div class="bg-slate-950/95 border border-slate-700/80 rounded-2xl w-full max-w-[96vw] h-[92vh] flex flex-col overflow-hidden shadow-2xl">
      <!-- ヘッダー -->
      <div class="px-4 py-2.5 border-b border-slate-800 flex items-center justify-between bg-slate-900/80 shrink-0">
        <div class="flex items-center gap-2 font-mono text-xs text-slate-200 min-w-0">
          <span class="px-2 py-0.5 rounded bg-blue-950 text-blue-300 font-bold border border-blue-800 shrink-0">{{ media.type?.toUpperCase() || 'MEDIA' }}</span>
          <span class="font-bold truncate max-w-md">{{ media.media_id || media.id }}</span>
          <span v-if="media.width && media.height" class="text-slate-400 text-[11px] shrink-0">({{ media.width }}x{{ media.height }})</span>
        </div>
        <button @click="emit('close')" title="閉じる (Esc)" class="w-7 h-7 rounded-lg bg-slate-800 hover:bg-slate-700 text-slate-400 hover:text-white flex items-center justify-center text-sm font-mono transition-colors">✕</button>
      </div>

      <!-- メインコンテンツ (左: 大画面メディアビューア, 右: 詳細インスペクタ＆エディタ) -->
      <div class="flex-1 flex flex-col md:flex-row min-h-0 overflow-hidden">
        <!-- 左ペイン: メディア本体 (アスペクト比維持 & 最大化表示) -->
        <div class="flex-1 h-full bg-black/90 flex items-center justify-center p-3 relative overflow-hidden" @click.self="emit('close')">
          <div v-if="isVideo && media.urls?.stream" class="w-full h-full flex items-center justify-center">
            <StashPlayer :src="media.urls.stream" :poster="media.urls.thumbnail" :stashSceneId="media.stash_scene_id" :autoplay="true" :show-expand-button="false" class="max-w-full max-h-full" @fullscreen-change="handleFullscreenChange" />
          </div>
          <img v-else-if="media.urls?.image || media.urls?.thumbnail" :src="media.urls.image || media.urls.thumbnail" :alt="media.media_id || media.id" class="max-w-full max-h-full object-contain rounded shadow-2xl select-none mx-auto" />
          <div v-else class="text-slate-500 font-mono text-xs">メディア実体を表示できません</div>
        </div>

        <!-- 右ペイン: 詳細インスペクタ ＆ メタデータエディタ -->
        <MediaInspectorPanel
          :media="media"
          @save-metadata="(p) => emit('saveMetadata', p)"
          @retry="(id) => emit('retry', id)"
          @purge="(id) => emit('purge', id)"
          @view-post="(artId) => emit('viewPost', artId)"
        />
      </div>
    </div>
  </div>
</template>
