<script setup lang="ts">
import { onMounted, onUnmounted, ref } from 'vue';
import { useTimeline } from './composables/useTimeline';
import ArticleCard from './components/article/ArticleCard.vue';
import AccountSelector from './components/timeline/AccountSelector.vue';
import TimelineFilter from './components/timeline/TimelineFilter.vue';
import MediaLightbox from './components/media/MediaLightbox.vue';

const {
  articles, accounts, selectedAccount, currentFilter, currentLang,
  loading, hasMore, activeMedia, selectAccount, setFilter, setLanguage,
  toggleLike, openLightbox, closeLightbox, loadMore,
} = useTimeline();

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
      <div class="flex items-center justify-between mb-3">
        <div>
          <h1 class="text-xl font-bold tracking-tight text-white">dozou_katanuki</h1>
          <p class="text-xs text-slate-400 font-mono">Dynamic Archival System (SPEC-FRONTEND-001)</p>
        </div>
        <div class="flex bg-slate-900 border border-slate-800 rounded-lg p-1 text-xs">
          <button
            v-for="lang in (['original', 'ja', 'en', 'zh'] as const)"
            :key="lang"
            @click="setLanguage(lang)"
            :class="[
              'px-2.5 py-1 rounded transition-colors',
              currentLang === lang ? 'bg-blue-600 text-white font-bold' : 'text-slate-400 hover:text-slate-200'
            ]"
          >
            {{ lang.toUpperCase() }}
          </button>
        </div>
      </div>
      <AccountSelector :accounts="accounts" :selectedId="selectedAccount" @select="selectAccount" />
    </header>

    <main class="w-full max-w-2xl">
      <TimelineFilter :currentFilter="currentFilter" @filter="setFilter" />

      <ArticleCard
        v-for="article in articles"
        :key="article.id"
        :article="article"
        :currentLang="currentLang"
        @toggleLike="toggleLike"
        @clickMedia="openLightbox"
      />

      <div ref="observerTarget" class="py-8 flex justify-center text-xs text-slate-500 font-mono">
        <span v-if="loading">Loading more archives...</span>
        <span v-else-if="!hasMore && articles.length > 0">End of local archive. ({{ articles.length }} items)</span>
      </div>
    </main>

    <MediaLightbox :media="activeMedia" @close="closeLightbox" />
  </div>
</template>
