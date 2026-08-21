<!-- frontend/src/components/media/StashPlayer.vue (100行以下) -->
<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue';

const props = withDefaults(
  defineProps<{ src: string; poster?: string; autoplay?: boolean; controls?: boolean; stashSceneId?: string }>(),
  { autoplay: false, controls: true }
);

const videoRef = ref<HTMLVideoElement | null>(null);
const isError = ref(false), retryCount = ref(0);

const stashUrl = computed(() => {
  if (props.stashSceneId) return `http://127.0.0.1:9999/scenes/${props.stashSceneId}`;
  if (props.src) {
    const match = props.src.match(/\/stash-proxy\/scene\/([^/]+)/);
    if (match?.[1]) return `http://127.0.0.1:9999/scenes/${match[1]}`;
  }
  return null;
});

const handleRetry = () => {
  isError.value = false; retryCount.value++;
  if (videoRef.value) videoRef.value.load();
};

const handleError = () => {
  isError.value = true;
  if (retryCount.value === 0) setTimeout(handleRetry, 2000);
};

watch(() => props.src, () => { retryCount.value = 0; isError.value = !props.src; });
onMounted(() => { isError.value = !props.src; });
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

    <!-- Stash スマート別窓導線ボタン -->
    <a
      v-if="stashUrl && !isError"
      :href="stashUrl"
      target="_blank"
      rel="noopener noreferrer"
      @click.stop
      class="absolute top-2.5 right-2.5 z-20 opacity-0 group-hover/player:opacity-100 transition-opacity bg-black/75 hover:bg-purple-600/90 text-white text-[11px] font-mono px-2.5 py-1 rounded-lg border border-white/20 flex items-center gap-1.5"
    >
      <span>📦</span><span>Stash で開く</span><span class="text-[9px]">↗</span>
    </a>

    <!-- エラー時フォールバック -->
    <div v-else-if="isError" class="flex flex-col items-center justify-center p-6 text-slate-400 gap-2">
      <span class="text-2xl">⚠️</span>
      <span class="text-xs font-mono text-slate-400">Media Stream Offline</span>
      <div class="flex items-center gap-2 mt-1">
        <button @click="handleRetry" class="px-2.5 py-1 text-xs bg-slate-800 hover:bg-slate-700 text-cyan-400 rounded font-mono border border-slate-700">
          🔄 再試行
        </button>
        <a v-if="stashUrl" :href="stashUrl" target="_blank" class="px-2.5 py-1 text-xs bg-purple-950/60 hover:bg-purple-900/80 text-purple-300 rounded font-mono border border-purple-700/50">
          📦 Stash確認
        </a>
      </div>
    </div>
  </div>
</template>
