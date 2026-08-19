<script setup lang="ts">
import { ref, onMounted, watch } from 'vue';

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
const retryCount = ref(0);

const checkSource = () => {
  isError.value = false;
  if (!props.src) {
    isError.value = true;
  }
};

const handleRetry = () => {
  isError.value = false;
  retryCount.value++;
  if (videoRef.value) {
    videoRef.value.load();
  }
};

const handleError = () => {
  isError.value = true;
  // 初回エラー時は2秒後に1回自動リトライ
  if (retryCount.value === 0) {
    setTimeout(() => {
      handleRetry();
    }, 2000);
  }
};

watch(() => props.src, () => {
  retryCount.value = 0;
  checkSource();
});
onMounted(checkSource);
</script>

<template>
  <div class="relative w-full h-full flex items-center justify-center bg-black rounded overflow-hidden select-none">
    <video
      v-if="!isError"
      ref="videoRef"
      :key="`${src}-${retryCount}`"
      :src="src"
      :poster="poster"
      :controls="controls"
      :autoplay="autoplay"
      playsinline
      preload="metadata"
      class="max-w-full max-h-full object-contain rounded"
      @error="handleError"
    >
      動画の再生に対応していません。
    </video>

    <!-- エラー・未接続時のフォールバック -->
    <div v-else class="flex flex-col items-center justify-center p-6 text-slate-400 gap-2">
      <svg class="w-10 h-10 text-red-400/80" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
      </svg>
      <span class="text-xs font-mono text-slate-400">Media Stream Offline (404/502)</span>
      <button
        @click="handleRetry"
        class="mt-1 px-2.5 py-1 text-xs bg-slate-800 hover:bg-slate-700 text-cyan-400 rounded font-mono border border-slate-700 transition-colors flex items-center gap-1 cursor-pointer"
      >
        <svg class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
        </svg>
        再試行 (Retry)
      </button>
    </div>
  </div>
</template>
