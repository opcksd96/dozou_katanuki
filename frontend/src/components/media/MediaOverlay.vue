<script setup lang="ts">
import { computed } from 'vue';
import type { RenderMedia } from '../../models/RenderTree';
import StashPlayer from './StashPlayer.vue';

const props = defineProps<{
  media: RenderMedia | null;
  hasNext?: boolean;
  hasPrev?: boolean;
}>();

const emit = defineEmits<{
  (e: 'close'): void;
  (e: 'next'): void;
  (e: 'prev'): void;
}>();

// Stash への直接導線 URL を算出
const stashDirectUrl = computed(() => {
  if (!props.media) return null;

  // 1. Scene ID がある場合
  if (props.media.stash_scene_id) {
    return `http://127.0.0.1:9999/scenes/${props.media.stash_scene_id}`;
  }
  // 2. Image ID がある場合
  if (props.media.stash_image_id) {
    return `http://127.0.0.1:9999/images/${props.media.stash_image_id}`;
  }
  // 3. URLs から抽出
  if (props.media.urls?.stream) {
    const match = props.media.urls.stream.match(/\/stash-proxy\/scene\/([^/]+)/);
    if (match && match[1]) {
      return `http://127.0.0.1:9999/scenes/${match[1]}`;
    }
  }
  if (props.media.urls?.image) {
    const match = props.media.urls.image.match(/\/stash-proxy\/image\/([^/]+)/);
    if (match && match[1]) {
      return `http://127.0.0.1:9999/images/${match[1]}`;
    }
  }
  // 4. Stash が有効ならルート
  return 'http://127.0.0.1:9999';
});
</script>

<template>
  <Transition name="fade">
    <div
      v-if="media"
      class="fixed inset-0 z-50 bg-black/95 backdrop-blur-md flex items-center justify-center p-4 select-none"
      @click="emit('close')"
    >
      <!-- 上部ツールバー -->
      <div class="absolute top-4 right-4 flex items-center gap-2 z-50" @click.stop>
        <!-- Stash スマート別窓導線 -->
        <a
          v-if="stashDirectUrl"
          :href="stashDirectUrl"
          target="_blank"
          rel="noopener noreferrer"
          title="Stash で開く (別窓)"
          class="flex items-center gap-1.5 px-3 py-1.5 rounded-full bg-purple-950/70 hover:bg-purple-900/90 text-purple-200 border border-purple-600/50 text-xs font-semibold backdrop-blur-md shadow-lg transition-colors cursor-pointer"
        >
          <span>📦</span>
          <span>Stash で開く</span>
          <span class="text-[10px] opacity-70">↗</span>
        </a>

        <!-- 原本リンク導線 (Wayback/CDN) -->
        <a
          v-if="media.urls?.original"
          :href="media.urls.original"
          target="_blank"
          rel="noopener noreferrer"
          title="オリジナル原本リンクを開く"
          class="flex items-center gap-1.5 px-3 py-1.5 rounded-full bg-slate-900/80 hover:bg-slate-800 text-slate-300 border border-slate-700 text-xs font-semibold backdrop-blur-md shadow-lg transition-colors cursor-pointer"
        >
          <span>🌐</span>
          <span>原本 (Source)</span>
          <span class="text-[10px] opacity-70">↗</span>
        </a>

        <!-- 閉じるボタン -->
        <button
          @click.stop="emit('close')"
          title="閉じる (Esc)"
          class="text-white/70 hover:text-white bg-white/10 hover:bg-white/20 p-2 rounded-full transition-colors"
        >
          <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>

      <!-- 前へ送りボタン -->
      <button
        v-if="hasPrev"
        @click.stop="emit('prev')"
        title="前のメディア (←)"
        class="absolute left-4 text-white/70 hover:text-white bg-white/10 hover:bg-white/20 p-3 rounded-full transition-colors z-50"
      >
        <svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
        </svg>
      </button>

      <!-- 次へ送りボタン -->
      <button
        v-if="hasNext"
        @click.stop="emit('next')"
        title="次のメディア (→)"
        class="absolute right-4 text-white/70 hover:text-white bg-white/10 hover:bg-white/20 p-3 rounded-full transition-colors z-50"
      >
        <svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
        </svg>
      </button>

      <!-- メディア本体コンテナ -->
      <div class="relative max-w-5xl max-h-[85vh] w-full flex items-center justify-center" @click.stop>
        <img
          v-if="media.type === 'image' || media.type === 'gif'"
          :src="media.urls.image || media.urls.original"
          :alt="media.id"
          class="max-w-full max-h-[85vh] object-contain rounded shadow-2xl transition-all"
        />
        <div v-else-if="media.type === 'video'" class="w-full max-w-4xl max-h-[80vh] flex items-center justify-center">
          <StashPlayer
            :src="media.urls.stream || media.urls.original"
            :poster="media.urls.thumbnail"
            :stashSceneId="media.stash_scene_id"
            :autoplay="true"
            :controls="true"
          />
        </div>
      </div>
    </div>
  </Transition>
</template>

<style scoped>
.fade-enter-active, .fade-leave-active { transition: opacity 0.2s ease; }
.fade-enter-from, .fade-leave-to { opacity: 0; }
</style>

