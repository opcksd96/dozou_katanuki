<script setup lang="ts">
import { onMounted, onUnmounted, ref } from 'vue';
import { useTimeline } from './composables/useTimeline';
import { useMediaOverlay } from './composables/useMediaOverlay';
import ArticleCard from './components/article/ArticleCard.vue';
import AccountSelector from './components/timeline/AccountSelector.vue';
import TimelineFilter from './components/timeline/TimelineFilter.vue';
import MediaOverlay from './components/media/MediaOverlay.vue';
import AdminModal from './components/admin/AdminModal.vue';

const isAdminOpen = ref(false);

const {
  articles, accounts, selectedAccount, currentFilter, systemLang,
  loading, hasMore, selectAccount, setFilter,
  toggleLike, loadMore, reloadAll,
} = useTimeline();

const {
  activeMedia, hasNext, hasPrev, openMedia, closeMedia, nextMedia, prevMedia,
} = useMediaOverlay();

const observerTarget = ref<HTMLElement | null>(null);
let observer: IntersectionObserver | null = null;

onMounted(() => {
  observer = new IntersectionObserver((entries) => {
    if (entries[0].isIntersecting && hasMore.value && !loading.value) {
      loadMore();
    }
  }, { rootMargin: '200px' });
  if (observerTarget.value) observer.observe(observerTarget.value);
});

onUnmounted(() => { observer?.disconnect(); });
</script>

<template>
  <div class="min-h-screen bg-slate-950 text-slate-100 flex flex-col items-center py-6 px-4">
    <header class="w-full max-w-2xl pb-4 border-b border-slate-800 mb-4">
      <div class="flex items-center justify-between mb-3 gap-2">
        <div>
          <h1 class="text-xl font-bold tracking-tight text-white flex items-center gap-2">
            dozou_katanuki
            <button
              @click="reloadAll"
              title="データを再読み込み (Ctrl+R / F5)"
              class="text-xs text-slate-400 hover:text-blue-400 transition-colors p-1 rounded hover:bg-slate-800"
            >
              🔄
            </button>
          </h1>
          <p class="text-xs text-slate-400 font-mono">Dynamic Archival System (SPEC-FRONTEND-001)</p>
        </div>
        <button
          @click="isAdminOpen = true"
          title="管理ダッシュボード ＆ 設定を開く"
          class="flex items-center gap-1.5 px-3 py-1.5 bg-slate-900 hover:bg-slate-800 text-slate-300 hover:text-white border border-slate-700/80 rounded-lg text-xs font-semibold transition-all shadow-sm"
        >
          <span>⚙️</span>
          <span>設定・ジョブ管理</span>
        </button>
      </div>
      <AccountSelector :accounts="accounts" :selectedId="selectedAccount" @select="selectAccount" />
    </header>

    <main class="w-full max-w-2xl">
      <TimelineFilter :currentFilter="currentFilter" @filter="setFilter" />

      <!-- 初期ロード時のスケルトン表示 -->
      <div v-if="articles.length === 0 && loading" class="space-y-4">
        <div v-for="i in 3" :key="i" class="bg-slate-900 border border-slate-800 rounded-xl p-4 animate-pulse">
          <div class="flex items-center space-x-3 mb-3">
            <div class="w-10 h-10 bg-slate-800 rounded-full"></div>
            <div class="space-y-1.5 flex-1">
              <div class="h-3.5 bg-slate-800 rounded w-1/3"></div>
              <div class="h-2.5 bg-slate-800 rounded w-1/4"></div>
            </div>
          </div>
          <div class="space-y-2 mb-3">
            <div class="h-3 bg-slate-800 rounded w-full"></div>
            <div class="h-3 bg-slate-800 rounded w-5/6"></div>
          </div>
          <div class="h-44 bg-slate-800/60 rounded-lg"></div>
        </div>
      </div>

      <ArticleCard
        v-for="article in articles"
        :key="article.id"
        :article="article"
        :targetLang="systemLang"
        @toggleLike="toggleLike"
        @clickMedia="(media, list) => openMedia(media, list)"
      />

      <div ref="observerTarget" class="py-8 flex flex-col items-center justify-center text-xs text-slate-500 font-mono gap-2">
        <span v-if="loading && articles.length > 0">Loading more archives...</span>
        <span v-else-if="!hasMore && articles.length > 0">End of local archive. ({{ articles.length }} items)</span>
        <template v-else-if="!loading && articles.length === 0">
          <span>No archive items found.</span>
          <button
            @click="reloadAll"
            class="px-3 py-1.5 bg-slate-800 hover:bg-slate-700 text-slate-200 rounded border border-slate-700 text-xs transition-colors flex items-center gap-1.5"
          >
            <span>🔄</span> 最新の情報に更新 (Ctrl+R)
          </button>
        </template>
      </div>
    </main>

    <MediaOverlay
      :media="activeMedia"
      :hasNext="hasNext"
      :hasPrev="hasPrev"
      @close="closeMedia"
      @next="nextMedia"
      @prev="prevMedia"
    />

    <AdminModal
      :isOpen="isAdminOpen"
      @close="isAdminOpen = false"
    />
  </div>
</template>
