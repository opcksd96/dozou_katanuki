<!-- frontend/src/components/media/StashPlayer.vue (100行以下) -->
<script setup lang="ts">
import { ref, shallowRef, computed, onMounted, onBeforeUnmount, watch } from 'vue';
import Plyr from 'plyr';
import Hls from 'hls.js';
import { EventsOn } from '../../../wailsjs/runtime/runtime';

const props = withDefaults(defineProps<{ src: string; poster?: string; autoplay?: boolean; controls?: boolean; stashSceneId?: string; showExpandButton?: boolean }>(), { autoplay: false, controls: true, showExpandButton: true });
const emit = defineEmits<{ (e: 'expand'): void; (e: 'fullscreenChange', active: boolean): void }>();

const videoRef = ref<HTMLVideoElement | null>(null), playerRef = shallowRef<Plyr | null>(null), hlsRef = shallowRef<Hls | null>(null);
const isError = ref(false), retryCount = ref(0);
let unoffStashReady: (() => void) | null = null;
const host = computed(() => typeof window !== 'undefined' && window.location?.hostname ? window.location.hostname : '127.0.0.1');
const stashUrl = computed(() => {
  if (props.stashSceneId) return `http://${host.value}:9999/scenes/${props.stashSceneId}`;
  const m = props.src?.match(/\/stash-proxy\/scene\/([^/]+)/);
  return m?.[1] ? `http://${host.value}:9999/scenes/${m[1]}` : null;
});

const destroyPlayer = () => { if (hlsRef.value) { hlsRef.value.destroy(); hlsRef.value = null; } if (playerRef.value) { playerRef.value.destroy(); playerRef.value = null; } };
const initPlayer = () => {
  destroyPlayer();
  if (!videoRef.value || !props.src) { isError.value = !props.src; return; }
  isError.value = false;
  if ((props.src.includes('.m3u8') || props.src.includes('/m3u8')) && Hls.isSupported()) {
    const hls = new Hls({ enableWorker: true, lowLatencyMode: true });
    hls.loadSource(props.src); hls.attachMedia(videoRef.value);
    hls.on(Hls.Events.ERROR, (_, data) => { if (data.fatal) handleError(); });
    hlsRef.value = hls;
  }
  playerRef.value = new Plyr(videoRef.value, {
    controls: ['play-large', 'play', 'rewind', 'fast-forward', 'progress', 'current-time', 'duration', 'mute', 'volume', 'settings', 'pip', 'fullscreen'],
    seekTime: 10, settings: ['speed'], keyboard: { focused: false, global: false }, autoplay: props.autoplay, muted: props.autoplay, clickToPlay: true
  });
  playerRef.value.on('enterfullscreen', () => emit('fullscreenChange', true));
  playerRef.value.on('exitfullscreen', () => emit('fullscreenChange', false));
};

const handleRetry = () => { isError.value = false; retryCount.value++; initPlayer(); };
const handleError = () => { isError.value = true; if (retryCount.value < 2) setTimeout(handleRetry, 2000 * (retryCount.value + 1)); };
watch(() => props.src, () => { retryCount.value = 0; initPlayer(); });
onMounted(() => {
  initPlayer();
  try { if ((window as any)?.runtime?.EventsOnMultiple) unoffStashReady = EventsOn('stash:ready', () => { if (isError.value) handleRetry(); }); } catch {}
});
onBeforeUnmount(() => { if (unoffStashReady) try { unoffStashReady(); } catch {} destroyPlayer(); });
</script>

<template>
  <div @click.stop class="relative w-full h-full flex items-center justify-center bg-black rounded overflow-hidden select-none group/player">
    <video v-show="!isError" ref="videoRef" :src="src" :poster="poster" playsinline preload="metadata" class="max-w-full max-h-full object-contain rounded" @error="handleError">動画の再生に対応していません。</video>
    <div class="absolute top-2.5 right-2.5 z-30 opacity-0 group-hover/player:opacity-100 transition-opacity flex items-center gap-1.5">
      <button v-if="showExpandButton && !isError" @click.stop="emit('expand')" class="bg-black/75 hover:bg-blue-600 text-white text-[11px] font-mono px-2.5 py-1 rounded-lg border border-white/20 flex items-center gap-1 cursor-pointer">
        <span>全画面詳細</span>
      </button>
      <a v-if="stashUrl && !isError" :href="stashUrl" target="_blank" rel="noopener noreferrer" @click.stop class="bg-black/75 hover:bg-purple-600 text-white text-[11px] font-mono px-2.5 py-1 rounded-lg border border-white/20 flex items-center gap-1">
        <span>📦 Stash ↗</span>
      </a>
    </div>
    <div v-if="isError" class="flex flex-col items-center justify-center p-6 text-slate-400 gap-2">
      <span class="text-2xl">⚠️</span>
      <div class="flex gap-2 mt-1">
        <button @click.stop="handleRetry" class="px-2.5 py-1 text-xs bg-slate-800 hover:bg-slate-700 text-cyan-400 rounded font-mono border border-slate-700">🔄 再試行</button>
        <a v-if="stashUrl" :href="stashUrl" target="_blank" @click.stop class="px-2.5 py-1 text-xs bg-purple-950 hover:bg-purple-900 text-purple-300 rounded font-mono border border-purple-700">📦 Stash確認</a>
      </div>
    </div>
  </div>
</template>
