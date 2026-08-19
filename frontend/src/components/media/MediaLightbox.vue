<script setup lang="ts">
import type { RenderMedia } from '../../models/RenderTree';

defineProps<{
  media: RenderMedia | null;
}>();

const emit = defineEmits<{
  (e: 'close'): void;
}>();
</script>

<template>
  <Transition name="fade">
    <div
      v-if="media"
      class="fixed inset-0 z-50 bg-black/90 backdrop-blur-md flex items-center justify-center p-4"
      @click="emit('close')"
    >
      <button
        @click.stop="emit('close')"
        class="absolute top-4 right-4 text-white/70 hover:text-white bg-white/10 hover:bg-white/20 p-2 rounded-full transition-colors z-50"
      >
        <svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
        </svg>
      </button>

      <div class="relative max-w-5xl max-h-[90vh] flex items-center justify-center" @click.stop>
        <img
          v-if="media.type === 'image' || media.type === 'gif'"
          :src="media.urls.image || media.urls.original"
          :alt="media.id"
          class="max-w-full max-h-[85vh] object-contain rounded shadow-2xl"
        />
        <video
          v-else-if="media.type === 'video'"
          :src="media.urls.stream || media.urls.original"
          controls
          autoplay
          class="max-w-full max-h-[85vh] rounded shadow-2xl"
        />
      </div>
    </div>
  </Transition>
</template>

<style scoped>
.fade-enter-active, .fade-leave-active { transition: opacity 0.2s ease; }
.fade-enter-from, .fade-leave-to { opacity: 0; }
</style>
