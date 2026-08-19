<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from 'vue';
import { useTimeline } from './composables/useTimeline';
import { useMediaOverlay } from './composables/useMediaOverlay';
import { useArticleDetail } from './composables/useArticleDetail';
import GlobalAppBar from './components/layout/GlobalAppBar.vue';
import ArticleCard from './components/article/ArticleCard.vue';
import ArticleDetailView from './components/article/ArticleDetailView.vue';
import AccountScopeSelector from './components/timeline/AccountScopeSelector.vue';
import AccountHeroHeader from './components/timeline/AccountHeroHeader.vue';
import TimelineFilter from './components/timeline/TimelineFilter.vue';
import MediaOverlay from './components/media/MediaOverlay.vue';
import AdminModal from './components/admin/AdminModal.vue';
import { EventsOn, EventsOff } from '../wailsjs/runtime/runtime';

const isAdminOpen = ref(false);
const activeArticleId = ref<string | null>(null);

const {
  articles, accounts, selectedAccount, currentFilter, systemLang,
  loading, hasMore, selectAccount, setFilter,
  toggleLike, loadMore, reloadAll,
} = useTimeline();

const {
  detail, loading: detailLoading, fetchDetail, clearDetail,
} = useArticleDetail();

const {
  activeMedia, activeArticle, hasNext, hasPrev, openMedia, closeMedia, nextMedia, prevMedia,
} = useMediaOverlay();

const currentAccountObj = computed(() => {
  if (selectedAccount.value === 'all') return null;
  return accounts.value.find((acc) => acc.numeric_id === selectedAccount.value) || null;
});

const activeArticleHandle = computed(() => {
  if (detail.value?.article) return detail.value.article.author.handle;
  return '';
});

const observerTarget = ref<HTMLElement | null>(null);
let observer: IntersectionObserver | null = null;

const openArticleDetail = (id: string) => {
  activeArticleId.value = id;
  fetchDetail(id);
  window.scrollTo({ top: 0, behavior: 'smooth' });
};

const closeArticleDetail = () => {
  activeArticleId.value = null;
  clearDetail();
};

const handleGlobalKeyDown = (e: KeyboardEvent) => {
  if ((e.ctrlKey || e.metaKey) && e.key === ',') {
    e.preventDefault();
    isAdminOpen.value = true;
  }
};

onMounted(() => {
  try {
    EventsOn('open:admin', () => {
      isAdminOpen.value = true;
    });
  } catch (err) {
    console.warn('Wails EventsOn not available:', err);
  }

  window.addEventListener('keydown', handleGlobalKeyDown);

  observer = new IntersectionObserver((entries) => {
    if (entries[0].isIntersecting && hasMore.value && !loading.value && !activeArticleId.value) {
      loadMore();
    }
  }, { rootMargin: '200px' });
  if (observerTarget.value) observer.observe(observerTarget.value);
});

onUnmounted(() => {
  observer?.disconnect();
  window.removeEventListener('keydown', handleGlobalKeyDown);
  try {
    EventsOff('open:admin');
  } catch {}
});
</script>

<template>
  <div class="min-h-screen bg-slate-950 text-slate-100 flex flex-col items-center">
    <!-- Level 0: グローバル・アプリバー (最上位 / システム全体コントロール) -->
    <GlobalAppBar
      :activeArticleId="activeArticleId"
      :activeArticleHandle="activeArticleHandle"
      @openAdmin="isAdminOpen = true"
      @reloadAll="reloadAll"
      @backToTimeline="closeArticleDetail"
    />

    <!-- メインコンテンツラッパー -->
    <div class="w-full max-w-2xl px-4 py-4">
      <!-- Level 1: スコープ選択領域 (タイムライン閲覧時のみ表示) -->
      <template v-if="!activeArticleId">
        <!-- 個別アカウント選択時のヒーローヘッダ -->
        <AccountHeroHeader
          v-if="currentAccountObj"
          :account="currentAccountObj"
          :totalArticles="articles.length"
          @backToAll="selectAccount('all')"
          @refresh="reloadAll"
        />

        <!-- アカウント切り替えスコープセレクター (全アカウント閲覧時) -->
        <AccountScopeSelector
          v-else
          :accounts="accounts"
          :selectedId="selectedAccount"
          @select="selectAccount"
        />
      </template>

      <!-- Level 2 & 3: タイムライン / 個別詳細コンテナ -->
      <main class="w-full border border-slate-800 rounded-2xl bg-slate-950 overflow-hidden shadow-xl min-h-[600px]">
        <!-- 1. 個別ツイート詳細ページ (Level 3 Detail Focus) -->
        <div v-if="activeArticleId">
          <div v-if="detailLoading && !detail" class="p-8 text-center text-slate-500 font-mono">
            Loading article detail and conversation thread...
          </div>
          <ArticleDetailView
            v-else-if="detail"
            :article="detail.article"
            :thread="detail.thread"
            :targetLang="systemLang"
            :loading="detailLoading"
            @back="closeArticleDetail"
            @selectArticle="openArticleDetail"
            @toggleLike="toggleLike"
            @clickMedia="(media, list, art) => openMedia(media, list, art)"
          />
        </div>

        <!-- 2. 通常タイムライン表示 (Level 2 Filter + Level 3 List) -->
        <div v-else>
          <!-- Level 2: コンテンツ種別フィルター -->
          <TimelineFilter :currentFilter="currentFilter" @filter="setFilter" />

          <!-- 初期ロード時のスケルトン表示 -->
          <div v-if="articles.length === 0 && loading" class="p-4 space-y-4">
            <div v-for="i in 3" :key="i" class="bg-slate-900/50 border border-slate-800/80 rounded-xl p-4 animate-pulse">
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

          <!-- 記事タイムライン一覧 -->
          <div class="divide-y divide-slate-800">
            <ArticleCard
              v-for="article in articles"
              :key="article.id"
              :article="article"
              :targetLang="systemLang"
              @clickArticle="openArticleDetail"
              @toggleLike="toggleLike"
              @clickMedia="(media, list, art) => openMedia(media, list, art)"
            />
          </div>

          <!-- スクロール終端 / 空状態表示 -->
          <div ref="observerTarget" class="py-8 flex flex-col items-center justify-center text-xs text-slate-500 font-mono gap-2">
            <span v-if="loading && articles.length > 0">Loading more archives...</span>
            <span v-else-if="!hasMore && articles.length > 0">End of local archive. ({{ articles.length }} items)</span>
            <template v-else-if="!loading && articles.length === 0">
              <span>No archive items found.</span>
              <button
                @click="reloadAll"
                class="px-3 py-1.5 bg-slate-800 hover:bg-slate-700 text-slate-200 rounded border border-slate-700 text-xs transition-colors flex items-center gap-1.5 cursor-pointer"
              >
                <span>🔄</span> 最新の情報に更新 (Ctrl+R)
              </button>
            </template>
          </div>
        </div>
      </main>
    </div>

    <!-- フローティングクイック設定ボタン (画面スクロール中も常時アクセス可能) -->
    <div class="fixed bottom-6 right-6 z-30">
      <button
        @click="isAdminOpen = true"
        title="設定・ジョブ管理を開く (Ctrl+,)"
        class="w-12 h-12 bg-blue-600 hover:bg-blue-500 text-white rounded-full shadow-lg shadow-blue-600/30 flex items-center justify-center text-xl transition-all transform hover:scale-105 active:scale-95 border border-blue-400/30 cursor-pointer"
      >
        ⚙️
      </button>
    </div>

    <!-- メディア Lightbox オーバーレイ -->
    <MediaOverlay
      :media="activeMedia"
      :article="activeArticle"
      :targetLang="systemLang"
      :hasNext="hasNext"
      :hasPrev="hasPrev"
      @close="closeMedia"
      @next="nextMedia"
      @prev="prevMedia"
      @toggleLike="toggleLike"
    />

    <!-- 管理・設定 AdminModal -->
    <AdminModal
      :isOpen="isAdminOpen"
      @close="isAdminOpen = false"
    />
  </div>
</template>
