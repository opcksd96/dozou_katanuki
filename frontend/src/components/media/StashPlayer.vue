<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue';

const props = withDefaults(
  defineProps<{
    src: string;
    poster?: string;
    autoplay?: boolean;
    controls?: boolean;
    stashSceneId?: string;
  }>(),
  {
    autoplay: false,
    controls: true,
  }
);

const videoRef = ref<HTMLVideoElement | null>(null);
const isError = ref(false);
const retryCount = ref(0);

const stashUrl = computed(() => {
  if (props.stashSceneId) {
    return `http://127.0.0.1:9999/scenes/${props.stashSceneId}`;
  }
  if (props.src) {
    const match = props.src.match(/\/stash-proxy\/scene\/([^/]+)/);
    if (match && match[1]) {
      return `http://127.0.0.1:9999/scenes/${match[1]}`;
    }
  }
  return null;
});

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
  <div class="relative w-full h-full flex items-center justify-center bg-black rounded overflow-hidden select-none group/player">
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

    <!-- Stash スマート別窓導線ボタン (ホバー時表示) -->
    <a
      v-if="stashUrl && !isError"
      :href="stashUrl"
      target="_blank"
      rel="noopener noreferrer"
      title="Stash でこのシーンを開く"
      @click.stop
      class="absolute top-2.5 right-2.5 z-20 opacity-0 group-hover/player:opacity-100 transition-opacity duration-200 bg-black/75 hover:bg-purple-600/90 text-white text-[11px] font-mono font-semibold px-2.5 py-1 rounded-lg backdrop-blur-md border border-white/20 shadow-lg flex items-center gap-1.5 cursor-pointer"
    >
      <span>📦</span>
      <span>Stash で開く</span>
      <span class="text-[9px] opacity-70">↗</span>
    </a>

    <!-- エラー・未接続時のフォールバック -->
    <div v-else-if="isError" class="flex flex-col items-center justify-center p-6 text-slate-400 gap-2">
      <svg class="w-10 h-10 text-red-400/80" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
      </svg>
      <span class="text-xs font-mono text-slate-400">Media Stream Offline (404/502)</span>
      
      <div class="flex items-center gap-2 mt-1">
        <button
          @click="handleRetry"
          class="px-2.5 py-1 text-xs bg-slate-800 hover:bg-slate-700 text-cyan-400 rounded font-mono border border-slate-700 transition-colors flex items-center gap-1 cursor-pointer"
        >
          <svg class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
          </svg>
          再試行 (Retry)
        </button>

        <a
          v-if="stashUrl"
          :href="stashUrl"
          target="_blank"
          rel="noopener noreferrer"
          class="px-2.5 py-1 text-xs bg-purple-950/60 hover:bg-purple-900/80 text-purple-300 rounded font-mono border border-purple-700/50 transition-colors flex items-center gap-1 cursor-pointer"
        >
          <span>📦</span>
          <span>Stash を確認</span>
        </a>
      </div>
    </div>
  </div>
</template>

