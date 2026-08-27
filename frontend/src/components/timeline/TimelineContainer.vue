<!-- frontend/src/components/timeline/TimelineContainer.vue (100行以下) -->
<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue';
import type { RenderTree, RenderMedia } from '../../models/RenderTree';
import type { LanguageCode, FilterType } from '../../composables/useTimeline';
import { useTimelineThread } from '../../composables/useTimelineThread';
import TimelineFilter from './TimelineFilter.vue';
import ArticleCard from '../article/ArticleCard.vue';
import ScrollTopBottomButtons from './ScrollTopBottomButtons.vue';

const props = defineProps<{
  articles: RenderTree[];
  currentFilter: FilterType;
  searchQuery: string;
  systemLang: LanguageCode;
  loading: boolean;
  hasMore: boolean;
  focusedArticleId?: string | null;
}>();

const emit = defineEmits<{
  (e: 'filter', f: FilterType): void;
  (e: 'clearSearch'): void;
  (e: 'loadMore'): void;
  (e: 'openDetail', id: string): void;
  (e: 'toggleLike', id: string): void;
  (e: 'retryMedia', id: string): void;
  (e: 'openMedia', m: RenderMedia, list: RenderMedia[], a: RenderTree): void;
  (e: 'clickTag', tag: string): void;
  (e: 'clickMention', handle: string): void;
}>();

const { loadingMap, expandedMap, toggleExpand, getParentArticles } = useTimelineThread();
const targetRef = ref<HTMLElement | null>(null);
let observer: IntersectionObserver | null = null;

const isConnected = (a?: RenderTree, b?: RenderTree) => {
  if (!a || !b) return false;
  return a.id === b.parent_id || b.id === a.parent_id || (!!a.conversation_id && a.conversation_id === b.conversation_id && a.author.handle === b.author.handle);
};

onMounted(() => {
  observer = new IntersectionObserver((entries) => {
    if (entries[0].isIntersecting && props.hasMore && !props.loading) emit('loadMore');
  }, { rootMargin: '200px' });
  if (targetRef.value) observer.observe(targetRef.value);
});
onUnmounted(() => observer?.disconnect());
</script>

<template>
  <div>
    <TimelineFilter :current-filter="currentFilter" @filter="(f) => emit('filter', f)" />

    <div v-if="searchQuery" class="px-4 py-2.5 bg-blue-950/40 border-b border-blue-900/50 flex items-center justify-between gap-3 text-xs">
      <div class="flex items-center gap-2 truncate">
        <span class="text-blue-400 font-semibold">🔍 絞込:</span>
        <span class="px-2 py-0.5 rounded bg-blue-900/60 text-blue-200 font-mono truncate">{{ searchQuery }}</span>
      </div>
      <button @click="emit('clearSearch')" class="px-2 py-1 rounded bg-slate-800 text-slate-300 text-[11px]">✕ 解除</button>
    </div>

    <div v-if="articles.length === 0 && loading" class="p-4 space-y-4">
      <div v-for="i in 3" :key="i" class="bg-slate-900/50 border border-slate-800 rounded-xl p-4 animate-pulse h-32"></div>
    </div>
    <div v-else-if="articles.length === 0" class="p-12 text-center text-slate-500 font-mono text-sm">
      No articles found. サルベージジョブを実行してください。
    </div>
    <div v-else class="divide-y divide-slate-800">
      <ArticleCard
        v-for="(art, idx) in articles"
        :key="art.id"
        :article="art"
        :target-lang="systemLang"
        :is-focused="focusedArticleId === art.id"
        :has-parent-line="idx > 0 && isConnected(art, articles[idx - 1])"
        :has-child-line="idx < articles.length - 1 && isConnected(art, articles[idx + 1])"
        :is-expanded="!!expandedMap[art.id]"
        :parent-articles="getParentArticles(art)"
        :loading-thread="!!loadingMap[art.id]"
        @toggle-expand-thread="toggleExpand"
        @click-article="(id) => emit('openDetail', id)"
        @toggle-like="(id) => emit('toggleLike', id)"
        @retry-media="(id) => emit('retryMedia', id)"
        @click-media="(m, l, a) => emit('openMedia', m, l, a)"
        @click-tag="(t) => emit('clickTag', t)"
        @click-mention="(m) => emit('clickMention', m)"
      />
    </div>

    <div ref="targetRef" class="py-6 flex justify-center text-xs font-mono text-slate-500">
      <span v-if="loading && articles.length > 0">⏳ Loading...</span>
      <span v-else-if="!hasMore && articles.length > 0">— End of Timeline —</span>
    </div>

    <ScrollTopBottomButtons />
  </div>
</template>

