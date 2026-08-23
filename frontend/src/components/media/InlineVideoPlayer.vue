<!-- frontend/src/components/media/InlineVideoPlayer.vue (100行以下) -->
<script setup lang="ts">
import { ref, computed, onMounted, onBeforeUnmount, watch } from 'vue';
import Hls from 'hls.js';

const props = defineProps<{ src: string; poster?: string; stashSceneId?: string }>();
const emit = defineEmits<{ (e: 'expand'): void }>();

const videoRef = ref<HTMLVideoElement | null>(null), hlsRef = ref<Hls | null>(null), isError = ref(false);
const host = computed(() => typeof window !== 'undefined' && window.location?.hostname ? window.location.hostname : '127.0.0.1');
const stashUrl = computed(() => {
  if (props.stashSceneId) return `http://${host.value}:9999/scenes/${props.stashSceneId}`;
  const m = props.src?.match(/\/stash-proxy\/scene\/([^/]+)/);
  return m?.[1] ? `http://${host.value}:9999/scenes/${m[1]}` : null;
});

const cleanupHls = () => { if (hlsRef.value) { hlsRef.value.destroy(); hlsRef.value = null; } };
const setupVideo = () => {
  cleanupHls(); isError.value = false;
  if (!videoRef.value || !props.src) return;
  const video = videoRef.value;
  if ((props.src.includes('.m3u8') || props.src.includes('/m3u8')) && Hls.isSupported()) {
    const hls = new Hls({ enableWorker: true, lowLatencyMode: true });
    hls.loadSource(props.src); hls.attachMedia(video);
    hls.on(Hls.Events.ERROR, (_, data) => { if (data.fatal) isError.value = true; });
    hlsRef.value = hls;
  } else video.src = props.src;
};

watch(() => props.src, setupVideo);
onMounted(setupVideo);
onBeforeUnmount(cleanupHls);
</script>

<template>
  <div @click.stop class="relative w-full h-full flex items-center justify-center bg-black rounded overflow-hidden select-none group/inline-player">
    <video v-show="!isError" ref="videoRef" :poster="poster" controls playsinline preload="metadata" class="w-full h-full max-h-[580px] object-contain rounded" @error="isError = true">動画の再生に対応していません。</video>
    <div class="absolute top-2.5 right-2.5 z-20 opacity-0 group-hover/inline-player:opacity-100 transition-opacity flex items-center gap-1.5 pointer-events-auto">
      <button v-if="!isError" @click.stop="emit('expand')" class="bg-black/80 hover:bg-blue-600 text-white text-[11px] font-mono px-2.5 py-1 rounded-lg border border-white/20 flex items-center gap-1"><span>全画面詳細</span></button>
      <a v-if="stashUrl && !isError" :href="stashUrl" target="_blank" rel="noopener noreferrer" @click.stop class="bg-black/80 hover:bg-purple-600 text-white text-[11px] font-mono px-2.5 py-1 rounded-lg border border-white/20 flex items-center gap-1"><span>📦 Stash ↗</span></a>
    </div>
    <div v-if="isError" class="flex flex-col items-center justify-center p-6 text-slate-400 gap-2">
      <span class="text-2xl">⚠️</span>
      <div class="flex items-center gap-2 mt-1">
        <button @click.stop="setupVideo" class="px-2.5 py-1 text-xs bg-slate-800 hover:bg-slate-700 text-cyan-400 rounded font-mono border border-slate-700">🔄 再試行</button>
        <a v-if="stashUrl" :href="stashUrl" target="_blank" @click.stop class="px-2.5 py-1 text-xs bg-purple-950 hover:bg-purple-900 text-purple-300 rounded font-mono border border-purple-700">📦 Stash確認</a>
      </div>
    </div>
  </div>
</template>
