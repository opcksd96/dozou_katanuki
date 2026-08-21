<!-- frontend/src/components/media/InlineVideoPlayer.vue -->
<script setup lang="ts">
import { ref, computed, onMounted, onBeforeUnmount, watch } from 'vue';
import Hls from 'hls.js';

const props = defineProps<{
  src: string;
  poster?: string;
  stashSceneId?: string;
}>();

const emit = defineEmits<{
  (e: 'expand'): void;
}>();

const videoRef = ref<HTMLVideoElement | null>(null);
const hlsRef = ref<Hls | null>(null);
const isError = ref(false);

const currentHostname = computed(() => {
  if (typeof window !== 'undefined' && window.location?.hostname) {
    return window.location.hostname;
  }
  return '127.0.0.1';
});

const stashUrl = computed(() => {
  if (props.stashSceneId) return `http://${currentHostname.value}:9999/scenes/${props.stashSceneId}`;
  const match = props.src?.match(/\/stash-proxy\/scene\/([^/]+)/);
  return match?.[1] ? `http://${currentHostname.value}:9999/scenes/${match[1]}` : null;
});

const cleanupHls = () => {
  if (hlsRef.value) {
    hlsRef.value.destroy();
    hlsRef.value = null;
  }
};

const setupVideo = () => {
  cleanupHls();
  isError.value = false;
  if (!videoRef.value || !props.src) return;

  const video = videoRef.value;
  if ((props.src.includes('.m3u8') || props.src.includes('/m3u8')) && Hls.isSupported()) {
    const hls = new Hls({ enableWorker: true, lowLatencyMode: true });
    hls.loadSource(props.src);
    hls.attachMedia(video);
    hls.on(Hls.Events.ERROR, (_, data) => {
      if (data.fatal) isError.value = true;
    });
    hlsRef.value = hls;
  } else {
    video.src = props.src;
  }
};

watch(() => props.src, setupVideo);
onMounted(setupVideo);
onBeforeUnmount(cleanupHls);
</script>

<template>
  <div @click.stop class="relative w-full h-full flex items-center justify-center bg-black rounded overflow-hidden select-none group/inline-player">
    <video
      v-show="!isError"
      ref="videoRef"
      :poster="poster"
      controls
      playsinline
      preload="metadata"
      class="w-full h-full max-h-[580px] object-contain rounded"
      @error="isError = true"
    >
      動画の再生に対応していません。
    </video>

    <!-- 右上オーバーレイツールバー (全画面詳細展開 & Stashリンク) -->
    <div class="absolute top-2.5 right-2.5 z-20 opacity-0 group-hover/inline-player:opacity-100 transition-opacity flex items-center gap-1.5 pointer-events-auto">
      <button
        v-if="!isError"
        @click.stop="emit('expand')"
        title="全画面・ツイート詳細オーバーレイで開く"
        class="bg-black/80 hover:bg-blue-600 text-white text-[11px] font-mono px-2.5 py-1 rounded-lg border border-white/20 flex items-center gap-1.5 cursor-pointer shadow-lg backdrop-blur-md transition-all active:scale-95"
      >
        <svg class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 8V4m0 0h4M4 4l5 5m11-1V4m0 0h-4m4 0l-5 5M4 16v4m0 0h4m-4 0l5-5m11 5l-5-5m5 5v-4m0 4h-4" />
        </svg>
        <span>全画面詳細</span>
      </button>

      <a
        v-if="stashUrl && !isError"
        :href="stashUrl"
        target="_blank"
        rel="noopener noreferrer"
        @click.stop
        title="Stash でこのシーンを開く"
        class="bg-black/80 hover:bg-purple-600 text-white text-[11px] font-mono px-2.5 py-1 rounded-lg border border-white/20 flex items-center gap-1.5 shadow-lg backdrop-blur-md transition-all"
      >
        <span>📦</span><span>Stash</span><span class="text-[9px] opacity-70">↗</span>
      </a>
    </div>

    <!-- エラー時表示 -->
    <div v-if="isError" class="flex flex-col items-center justify-center p-6 text-slate-400 gap-2">
      <span class="text-2xl">⚠️</span>
      <div class="flex items-center gap-2 mt-1">
        <button @click.stop="setupVideo" class="px-2.5 py-1 text-xs bg-slate-800 hover:bg-slate-700 text-cyan-400 rounded font-mono border border-slate-700 cursor-pointer">🔄 再試行</button>
        <a v-if="stashUrl" :href="stashUrl" target="_blank" @click.stop class="px-2.5 py-1 text-xs bg-purple-950/60 hover:bg-purple-900/80 text-purple-300 rounded font-mono border border-purple-700/50">📦 Stash確認</a>
      </div>
    </div>
  </div>
</template>
