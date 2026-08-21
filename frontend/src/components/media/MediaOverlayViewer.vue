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
  <div @click.self="$emit('close')" class="flex-1 relative flex items-center justify-center p-4 min-h-0 cursor-pointer">
    <!-- 前へボタン -->
    <button v-if="hasPrev" @click.stop="$emit('prev')" class="absolute left-4 z-20 w-12 h-12 rounded-full bg-black/60 hover:bg-black/90 text-white flex items-center justify-center text-xl transition-all shadow-lg">
      ‹
    </button>

    <!-- メディア表示本体 -->
    <div @click.self="$emit('close')" class="max-w-full max-h-full flex items-center justify-center cursor-default">
      <img
        v-if="media.type === 'image' && media.urls?.image"
        :src="media.urls.image"
        @click.stop
        class="max-w-full max-h-[85vh] object-contain rounded-lg shadow-2xl select-none"
        alt="Enlarged media"
      />
      <StashPlayer
        v-else-if="media.urls?.stream"
        :src="media.urls.stream"
        :poster="media.urls.thumbnail"
        :stashSceneId="media.stash_scene_id"
        :autoplay="true"
        class="max-w-full max-h-[85vh] shadow-2xl rounded-lg overflow-hidden"
      />
      <div v-else class="text-slate-500 font-mono text-sm p-8 bg-slate-900 rounded-xl border border-slate-800">
        メディアを読み込めません
      </div>
    </div>

    <!-- 次へボタン -->
    <button v-if="hasNext" @click.stop="$emit('next')" class="absolute right-4 z-20 w-12 h-12 rounded-full bg-black/60 hover:bg-black/90 text-white flex items-center justify-center text-xl transition-all shadow-lg">
      ›
    </button>
  </div>
</template>

