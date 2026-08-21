<!-- frontend/src/components/media/MediaOverlay.vue (100行以下) -->
<script setup lang="ts">
import { onMounted, onUnmounted } from 'vue';
import type { RenderMedia, RenderTree } from '../../models/RenderTree';
import type { LanguageCode } from '../../composables/useTimeline';
import MediaOverlayViewer from './MediaOverlayViewer.vue';
import MediaOverlaySidebar from './MediaOverlaySidebar.vue';

const props = defineProps<{
  media: RenderMedia | null;
  article?: RenderTree | null;
  targetLang?: LanguageCode;
  hasNext?: boolean;
  hasPrev?: boolean;
}>();

const emit = defineEmits<{
  (e: 'close'): void;
  (e: 'next'): void;
  (e: 'prev'): void;
  (e: 'toggleLike', id: string): void;
}>();

const handleKeyDown = (e: KeyboardEvent) => {
  if (!props.media) return;
  if (e.key === 'ArrowRight' || e.key === 'd' || e.key === 'D') { if (props.hasNext) { e.preventDefault(); emit('next'); } }
  else if (e.key === 'ArrowLeft' || e.key === 'a' || e.key === 'A') { if (props.hasPrev) { e.preventDefault(); emit('prev'); } }
  else if (e.key === 'Escape') { e.preventDefault(); emit('close'); }
};

onMounted(() => window.addEventListener('keydown', handleKeyDown));
onUnmounted(() => window.removeEventListener('keydown', handleKeyDown));
</script>

<template>
  <Transition name="fade">
    <div v-if="media" @click.self="$emit('close')" class="fixed inset-0 z-50 bg-black/95 backdrop-blur-md flex flex-col md:flex-row overflow-hidden">
      <!-- 閉じるボタン -->
      <button @click="$emit('close')" class="absolute top-4 right-4 z-30 w-10 h-10 rounded-full bg-black/60 hover:bg-black/90 text-white flex items-center justify-center text-lg transition-colors">✕</button>

      <MediaOverlayViewer
        :media="media"
        :has-next="hasNext"
        :has-prev="hasPrev"
        @next="$emit('next')"
        @prev="$emit('prev')"
        @close="$emit('close')"
      />

      <MediaOverlaySidebar
        v-if="article"
        :media="media"
        :article="article"
        :target-lang="targetLang"
        @toggle-like="(id) => $emit('toggleLike', id)"
      />
    </div>
  </Transition>
</template>

<style scoped>
.fade-enter-active, .fade-leave-active { transition: opacity 0.2s ease; }
.fade-enter-from, .fade-leave-to { opacity: 0; }
</style>
