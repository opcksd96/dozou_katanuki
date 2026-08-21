<!-- frontend/src/components/media/MediaOverlay.vue -->
<script setup lang="ts">
import { onMounted, onUnmounted } from 'vue';
import type { RenderMedia, RenderTree } from '../../models/RenderTree';
import type { LanguageCode } from '../../composables/useTimeline';
import MediaOverlayViewer from './MediaOverlayViewer.vue';
import MediaOverlayBottomCard from './MediaOverlayBottomCard.vue';

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
    <div v-if="media" @click.self="$emit('close')" class="fixed inset-0 z-50 bg-black flex flex-col items-center justify-center overflow-hidden select-none">
      <!-- メディア全画面ビューワー -->
      <MediaOverlayViewer
        :media="media"
        :has-next="hasNext"
        :has-prev="hasPrev"
        @next="$emit('next')"
        @prev="$emit('prev')"
        @close="$emit('close')"
      />

      <!-- 下部自然フェードグラデーション (動画上の文字可読性向上用・ボックスなし) -->
      <div class="absolute inset-x-0 bottom-0 h-64 bg-gradient-to-t from-black/90 via-black/40 to-transparent pointer-events-none z-20"></div>

      <!-- 左上/右上 ナビゲーションバー (戻る/閉じる) -->
      <div class="absolute top-4 inset-x-4 z-40 flex items-center justify-between pointer-events-none">
        <button
          @click="$emit('close')"
          title="戻る / 閉じる (Esc)"
          class="w-10 h-10 rounded-full bg-black/50 hover:bg-black/80 text-white flex items-center justify-center text-lg transition-colors border border-white/20 shadow-lg backdrop-blur-md cursor-pointer pointer-events-auto active:scale-95"
        >
          ✕
        </button>
      </div>

      <!-- 下部文字情報 ＆ アクション列 (ダイレクトオーバーレイ) -->
      <div v-if="article" class="absolute bottom-0 inset-x-0 z-30 pointer-events-none">
        <MediaOverlayBottomCard
          :media="media"
          :article="article"
          :target-lang="targetLang"
          @toggle-like="(id) => $emit('toggleLike', id)"
        />
      </div>
    </div>
  </Transition>
</template>

<style scoped>
.fade-enter-active, .fade-leave-active { transition: opacity 0.2s ease; }
.fade-enter-from, .fade-leave-to { opacity: 0; }
</style>
