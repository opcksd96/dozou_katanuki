<!-- frontend/src/views/mobile/MobileConsole.vue -->
<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { useTimeline } from '../../composables/useTimeline';
import { useMediaOverlay } from '../../composables/useMediaOverlay';
import MediaOverlay from '../../components/media/MediaOverlay.vue';
import { ArrowLeft, Server, Settings2, Smartphone, Image as ImageIcon } from 'lucide-vue-next';

const router = useRouter();
const { articles, systemLang } = useTimeline();
const { activeMedia, activeArticle, hasNext, hasPrev, openMedia, closeMedia, nextMedia, prevMedia } = useMediaOverlay();

const isStashOnline = ref(false);

const checkStash = async () => {
  try {
    const r = await fetch('/stash-proxy/', { method: 'HEAD' }); 
    isStashOnline.value = r.ok || r.status === 401 || r.status === 404;
  } catch { isStashOnline.value = false; }
};

const backToWebUI = () => {
  router.push('/webui');
};

const openLatestMedia = () => {
  const mediaArticles = articles.value.filter(a => a.media && a.media.length > 0);
  if (mediaArticles.length > 0) {
    const art = mediaArticles[0];
    if (art.media) openMedia(art.media[0], art.media, art);
  }
};

onMounted(() => {
  checkStash();
});
</script>

<template>
  <div class="min-h-screen bg-slate-950 text-slate-100 flex flex-col font-sans">
    <!-- モバイル用ヘッダー -->
    <header class="px-4 py-3 border-b border-slate-800 bg-slate-900 sticky top-0 z-10 flex items-center justify-between">
      <div class="flex items-center gap-2">
        <Smartphone class="w-5 h-5 text-blue-400" />
        <h1 class="text-sm font-bold">Mobile Console</h1>
      </div>
      <button @click="backToWebUI" class="p-2 rounded-lg bg-slate-800 border border-slate-700 text-slate-300 active:scale-95">
        <ArrowLeft class="w-4 h-4" />
      </button>
    </header>

    <!-- ダッシュボードコンテンツ (AdSense風) -->
    <main class="flex-1 p-4 space-y-4">
      <div class="grid grid-cols-2 gap-3">
        <div class="bg-slate-900 border border-slate-800 rounded-xl p-4 flex flex-col items-center justify-center gap-2">
          <span class="text-xs text-slate-400 font-semibold">Total Posts</span>
          <span class="text-2xl font-bold text-slate-100">{{ articles.length }}</span>
        </div>
        <div class="bg-slate-900 border border-slate-800 rounded-xl p-4 flex flex-col items-center justify-center gap-2">
          <span class="text-xs text-slate-400 font-semibold">Stash Status</span>
          <div class="flex items-center gap-1.5">
            <span :class="['w-2 h-2 rounded-full', isStashOnline ? 'bg-emerald-400' : 'bg-amber-500']"></span>
            <span class="text-sm font-bold text-slate-100">{{ isStashOnline ? 'Online' : 'Offline' }}</span>
          </div>
        </div>
      </div>

      <div class="bg-slate-900 border border-slate-800 rounded-xl p-4 space-y-3">
        <h2 class="text-xs font-bold text-slate-400 uppercase tracking-wider flex items-center gap-2">
          <Settings2 class="w-4 h-4" /> Operations
        </h2>
        <div class="grid gap-2">
          <button @click="openLatestMedia" class="w-full py-3 px-4 rounded-lg bg-blue-600/20 hover:bg-blue-600/30 border border-blue-500/30 text-blue-400 text-sm font-bold flex items-center justify-center gap-2 active:scale-95 transition-all">
            <ImageIcon class="w-4 h-4" /> 最新メディアを開く
          </button>
        </div>
      </div>
    </main>

    <!-- 既存のMediaOverlayを再利用 -->
    <MediaOverlay 
      :media="activeMedia" 
      :article="activeArticle" 
      :target-lang="systemLang" 
      :has-next="hasNext" 
      :has-prev="hasPrev" 
      @close="closeMedia" 
      @next="nextMedia" 
      @prev="prevMedia" 
      @toggle-like="() => {}" 
    />
  </div>
</template>
