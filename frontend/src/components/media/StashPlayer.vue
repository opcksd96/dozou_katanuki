<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch } from 'vue';

const props = withDefaults(
  defineProps<{
    src: string;
    poster?: string;
    autoplay?: boolean;
    controls?: boolean;
  }>(),
  {
    autoplay: false,
    controls: true,
  }
);

const videoRef = ref<HTMLVideoElement | null>(null);
const isError = ref(false);

const checkSource = () => {
  isError.value = false;
  if (!props.src) {
    isError.value = true;
  }
};

watch(() => props.src, checkSource);
onMounted(checkSource);
</script>

<template>
  <div class="relative w-full h-full flex items-center justify-center bg-black rounded overflow-hidden select-none">
    <video
      v-if="!isError"
      ref="videoRef"
      :src="src"
      :poster="poster"
      :controls="controls"
      :autoplay="autoplay"
      playsinline
      preload="metadata"
      class="max-w-full max-h-full object-contain rounded"
      @error="isError = true"
    >
      <source :src="src" type="video/mp4" />
      動画の再生に対応していません。
    </video>

    <!-- エラー・未接続時のフォールバック -->
    <div v-else class="flex flex-col items-center justify-center p-6 text-slate-400 gap-2">
      <svg class="w-12 h-12 text-red-400/80" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
      </svg>
      <span class="text-xs font-mono">Stream Load Error (404/502)</span>
    </div>
  </div>
</template>
