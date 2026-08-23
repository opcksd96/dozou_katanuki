<!-- frontend/src/App.vue (100行以下 - SPEC-FRONTEND-001) -->
<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';
import { useTimeline } from './composables/useTimeline';
import { useMediaOverlay } from './composables/useMediaOverlay';
import { useArticleDetail } from './composables/useArticleDetail';
import { useSkin } from './composables/useSkin';
import { useKeyboardNavigation } from './composables/useKeyboardNavigation';
import { useKeyboardReload } from './composables/useKeyboardReload';
import GlobalAppBar from './components/layout/GlobalAppBar.vue';
import AccountScopeSelector from './components/timeline/AccountScopeSelector.vue';
import AccountHeroHeader from './components/timeline/AccountHeroHeader.vue';
import TimelineContainer from './components/timeline/TimelineContainer.vue';
import ArticleDetailView from './components/article/ArticleDetailView.vue';
import MediaOverlay from './components/media/MediaOverlay.vue';
import AdminModal from './components/admin/AdminModal.vue';
import KeyboardShortcutModal from './components/layout/KeyboardShortcutModal.vue';
import ToastContainer from './components/layout/ToastContainer.vue';

const isAdminOpen = ref(false), activeArticleId = ref<string | null>(null);
const {
  articles, accounts, selectedAccount, currentFilter, searchQuery, systemLang,
  loading, hasMore, renderKey, isStashReady, selectAccount, setFilter, setSearchQuery, clearSearchQuery,
  toggleLike, retryMedia, loadMore, reloadAll,
} = useTimeline();

const { detail, loading: detailLoading, fetchDetail, clearDetail } = useArticleDetail();
const { activeMedia, activeArticle, hasNext, hasPrev, openMedia, closeMedia, nextMedia, prevMedia } = useMediaOverlay();
const { loadSkin } = useSkin();
useKeyboardReload();

const currentAccountObj = computed(() => accounts.value.find((a) => a.numeric_id === selectedAccount.value) || null);
const currentNavItems = computed(() => (activeArticleId.value && detail.value) ? [detail.value.article, ...(detail.value.thread || [])] : articles.value);
const openDetail = (id: string) => { activeArticleId.value = id; fetchDetail(id); window.scrollTo({ top: 0, behavior: 'smooth' }); };
const closeDetail = () => { activeArticleId.value = null; clearDetail(); };

const { focusedIndex, isHelpOpen } = useKeyboardNavigation({
  getItems: () => currentNavItems.value,
  onSelectArticle: openDetail, onToggleLike: toggleLike,
  onOpenMedia: (art) => { if (art.media?.length) openMedia(art.media[0], art.media, art); },
  onBack: closeDetail, isDetailView: () => !!activeArticleId.value,
  isOverlayOpen: () => !!activeMedia.value, isAdminOpen: () => isAdminOpen.value, openAdmin: () => { isAdminOpen.value = true; },
});

const focusedId = computed(() => (focusedIndex.value >= 0 && focusedIndex.value < currentNavItems.value.length) ? currentNavItems.value[focusedIndex.value].id : null);
onMounted(() => {
  loadSkin('twitter');
  try { if ((window as any)?.runtime?.EventsOnMultiple) (window as any).runtime.EventsOnMultiple('open:admin', () => { isAdminOpen.value = true; }, -1); } catch {}
});
</script>

<template>
  <div class="min-h-screen bg-slate-950 text-slate-100 flex flex-col items-center">
    <GlobalAppBar :active-article-id="activeArticleId" :active-article-handle="detail?.article?.author.handle || ''" :is-stash-online="isStashReady" @open-admin="isAdminOpen = true" @back-to-timeline="closeDetail" />

    <div class="w-full max-w-2xl px-4 py-4">
      <template v-if="!activeArticleId">
        <AccountHeroHeader v-if="currentAccountObj" :account="currentAccountObj" :total-articles="articles.length" @back-to-all="selectAccount('all')" @refresh="reloadAll" />
        <AccountScopeSelector v-else :accounts="accounts" :selected-id="selectedAccount" @select="selectAccount" />
      </template>

      <main :key="renderKey" class="w-full border border-slate-800 rounded-2xl bg-slate-950 overflow-hidden shadow-xl min-h-[600px]">
        <!-- Stash 待機画面 -->
        <div v-if="!isStashReady" class="p-16 flex flex-col items-center justify-center min-h-[500px] text-center space-y-6 animate-pulse">
          <div class="w-20 h-20 rounded-2xl bg-blue-600/20 border border-blue-500/40 flex items-center justify-center text-4xl shadow-xl">📦</div>
          <div class="space-y-2 max-w-sm"><h2 class="text-base font-bold text-slate-100">Stash メディアサーバー接続確認中...</h2><p class="text-xs text-slate-400 font-mono">ポート9999疎通プロービング中</p></div>
        </div>
        <!-- タイムライン / 詳細 -->
        <div v-else-if="activeArticleId">
          <div v-if="detailLoading && !detail" class="p-8 text-center text-slate-500 font-mono">Loading detail...</div>
          <ArticleDetailView v-else-if="detail" :article="detail.article" :thread="detail.thread" :target-lang="systemLang" :loading="detailLoading" :focused-article-id="focusedId || undefined" @back="closeDetail" @select-article="openDetail" @toggle-like="toggleLike" @retry-media="retryMedia" @click-tag="(t) => { closeDetail(); setSearchQuery('#' + t); }" @click-mention="(m) => { closeDetail(); setSearchQuery('@' + m); }" @click-media="(m, l, a) => openMedia(m, l, a)" />
        </div>
        <TimelineContainer v-else :articles="articles" :current-filter="currentFilter" :search-query="searchQuery" :system-lang="systemLang" :loading="loading" :has-more="hasMore" :focused-article-id="focusedId" @filter="setFilter" @clear-search="clearSearchQuery" @load-more="loadMore" @open-detail="openDetail" @toggle-like="toggleLike" @retry-media="retryMedia" @open-media="(m, l, a) => openMedia(m, l, a)" @click-tag="(t) => setSearchQuery('#' + t)" @click-mention="(m) => setSearchQuery('@' + m)" />
      </main>
    </div>

    <KeyboardShortcutModal :is-open="isHelpOpen" @close="isHelpOpen = false" />
    <AdminModal :is-open="isAdminOpen" @close="isAdminOpen = false" @whitelist-updated="reloadAll" />
    <MediaOverlay :media="activeMedia" :article="activeArticle" :target-lang="systemLang" :has-next="hasNext" :has-prev="hasPrev" @close="closeMedia" @next="nextMedia" @prev="prevMedia" @toggle-like="toggleLike" />
    <ToastContainer />
  </div>
</template>
