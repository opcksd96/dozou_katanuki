<!-- frontend/src/components/media/MediaOverlayViewer.vue (100行以下) -->
<script setup lang="ts">
import type { RenderMedia } from '../../models/RenderTree';
import StashPlayer from './StashPlayer.vue';

defineProps<{
  media: RenderMedia;
  hasNext?: boolean;
  hasPrev?: boolean;
}>();
defineEmits<{ (e: 'next'): void; (e: 'prev'): void; (e: 'close'): void }>();
</script>

<template>
  <div @click.self="$emit('close')" class="flex-1 w-full h-full relative flex items-center justify-center p-2 md:p-4 min-h-0 cursor-pointer overflow-hidden">
    <!-- 前へボタン -->
    <button v-if="hasPrev" @click.stop="$emit('prev')" class="absolute left-4 z-30 w-12 h-12 rounded-full bg-black/60 hover:bg-black/90 text-white flex items-center justify-center text-xl transition-all shadow-lg backdrop-blur-sm cursor-pointer active:scale-95">
      ‹
    </button>

    <!-- メディア表示本体 (アスペクト比完全維持 & 画面内フィット) -->
    <div @click.self="$emit('close')" class="w-full h-full flex items-center justify-center cursor-default pointer-events-auto">
      <img
        v-if="(media.type === 'image' || media.type === 'gif') && (media.urls?.image || media.urls?.thumbnail)"
        :src="media.urls.image || media.urls.thumbnail"
        @click.stop
        class="max-w-[95vw] max-h-[88vh] w-auto h-auto object-contain rounded-lg shadow-2xl select-none mx-auto transition-transform"
        alt="Enlarged media"
      />
      <div
        v-else-if="media.urls?.stream"
        @click.stop
        class="w-full h-full max-w-[95vw] max-h-[88vh] flex items-center justify-center"
      >
        <StashPlayer
          :src="media.urls.stream"
          :poster="media.urls.thumbnail"
          :stashSceneId="media.stash_scene_id"
          :autoplay="true"
          :show-expand-button="false"
          class="w-full h-full max-w-full max-h-full flex items-center justify-center"
        />
      </div>
      <div v-else class="text-slate-500 font-mono text-sm p-8 bg-slate-900 rounded-xl border border-slate-800">
        メディアを読み込めません
      </div>
    </div>

    <!-- 次へボタン -->
    <button v-if="hasNext" @click.stop="$emit('next')" class="absolute right-4 z-30 w-12 h-12 rounded-full bg-black/60 hover:bg-black/90 text-white flex items-center justify-center text-xl transition-all shadow-lg backdrop-blur-sm cursor-pointer active:scale-95">
      ›
    </button>
  </div>
</template>

