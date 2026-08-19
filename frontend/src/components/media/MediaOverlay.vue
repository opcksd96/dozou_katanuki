<script setup lang="ts">
import type { RenderMedia } from '../../models/RenderTree';
import StashPlayer from './StashPlayer.vue';

defineProps<{
  media: RenderMedia | null;
  hasNext?: boolean;
  hasPrev?: boolean;
}>();

const emit = defineEmits<{
  (e: 'close'): void;
  (e: 'next'): void;
  (e: 'prev'): void;
}>();
</script>

<template>
  <Transition name="fade">
    <div
      v-if="media"
      class="fixed inset-0 z-50 bg-black/90 backdrop-blur-md flex items-center justify-center p-4 select-none"
      @click="emit('close')"
    >
      <!-- 閉じるボタン -->
      <button
        @click.stop="emit('close')"
        title="閉じる (Esc)"
        class="absolute top-4 right-4 text-white/70 hover:text-white bg-white/10 hover:bg-white/20 p-2.5 rounded-full transition-colors z-50"
      >
        <svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
        </svg>
      </button>

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
