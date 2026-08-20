<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from 'vue';
import { useTimeline } from './composables/useTimeline';
import { useMediaOverlay } from './composables/useMediaOverlay';
import { useArticleDetail } from './composables/useArticleDetail';
import { useSkin } from './composables/useSkin';
import { useKeyboardNavigation } from './composables/useKeyboardNavigation';
import GlobalAppBar from './components/layout/GlobalAppBar.vue';
import ArticleCard from './components/article/ArticleCard.vue';
import ArticleDetailView from './components/article/ArticleDetailView.vue';
import AccountScopeSelector from './components/timeline/AccountScopeSelector.vue';
import AccountHeroHeader from './components/timeline/AccountHeroHeader.vue';
import TimelineFilter from './components/timeline/TimelineFilter.vue';
import MediaOverlay from './components/media/MediaOverlay.vue';
import AdminModal from './components/admin/AdminModal.vue';
import KeyboardShortcutModal from './components/layout/KeyboardShortcutModal.vue';
import { EventsOn, EventsOff } from '../wailsjs/runtime/runtime';

const isAdminOpen = ref(false);
const activeArticleId = ref<string | null>(null);

const {
  articles, accounts, selectedAccount, currentFilter, searchQuery, systemLang,
  loading, hasMore, selectAccount, setFilter, setSearchQuery, clearSearchQuery,
  toggleLike, retryMedia, loadMore, reloadAll,
} = useTimeline();

const {
  detail, loading: detailLoading, fetchDetail, clearDetail,
} = useArticleDetail();

const {
  activeMedia, activeArticle, hasNext, hasPrev, openMedia, closeMedia, nextMedia, prevMedia,
} = useMediaOverlay();

const { loadSkin, applyCSS } = useSkin();

const currentAccountObj = computed(() => {
  if (selectedAccount.value === 'all') return null;
  return accounts.value.find((acc) => acc.numeric_id === selectedAccount.value) || null;
});

const activeArticleHandle = computed(() => {
  if (detail.value?.article) return detail.value.article.author.handle;
  return '';
});

// 現在アクティブ/選択可能なアイテムリスト
const currentNavItems = computed(() => {
  if (activeArticleId.value && detail.value) {
    const list = [detail.value.article];
    if (detail.value.thread) {
      list.push(...detail.value.thread.filter((t) => t.id !== detail.value?.article.id));
    }
    return list;
  }
  return articles.value;
});

const openArticleDetail = (id: string) => {
  activeArticleId.value = id;
  fetchDetail(id);
  window.scrollTo({ top: 0, behavior: 'smooth' });
};

const closeArticleDetail = () => {
  activeArticleId.value = null;
  clearDetail();
};

const handleTagClick = (tag: string) => {
  closeArticleDetail();
  setSearchQuery('#' + tag);
};

const handleMentionClick = (handle: string) => {
  closeArticleDetail();
  const found = accounts.value.find(
    (a) => a.handle.toLowerCase() === handle.toLowerCase() || a.numeric_id === handle
  );
  if (found) {
    selectAccount(found.numeric_id);
  } else {
    setSearchQuery('@' + handle);
  }
};

// キーボードナビゲーションの統合
const { focusedIndex, isHelpOpen } = useKeyboardNavigation({
  getItems: () => currentNavItems.value,
  onSelectArticle: (id) => openArticleDetail(id),
  onToggleLike: (id) => toggleLike(id),
  onOpenMedia: (art) => {
    if (art.media && art.media.length > 0) {
      openMedia(art.media[0], art.media, art);
    }
  },
  onBack: () => closeArticleDetail(),
  isDetailView: () => !!activeArticleId.value,
  isOverlayOpen: () => !!activeMedia.value,
  isAdminOpen: () => isAdminOpen.value,
  openAdmin: () => { isAdminOpen.value = true; },
});

const focusedArticleId = computed(() => {
  if (focusedIndex.value >= 0 && focusedIndex.value < currentNavItems.value.length) {
    return currentNavItems.value[focusedIndex.value].id;
  }
  return null;
});

const observerTarget = ref<HTMLElement | null>(null);
let observer: IntersectionObserver | null = null;

const applyInitialFontsAndSkin = async () => {
  try {
    const app = (window as any)?.go?.main?.App;
    if (app && typeof app.GetConfig === 'function') {
      const cfg = await app.GetConfig();
      if (cfg?.appearance) {
        const root = document.documentElement;
        if (cfg.appearance.font_family_ja) root.style.setProperty('--font-family-ja', cfg.appearance.font_family_ja);
        if (cfg.appearance.font_family_en) root.style.setProperty('--font-family-en', cfg.appearance.font_family_en);
        if (cfg.appearance.font_family_zh) root.style.setProperty('--font-family-zh', cfg.appearance.font_family_zh);
      }
    }
    await loadSkin('twitter', {
      platform: 'twitter',
      router: { push: (path) => console.log('Skin navigate:', path) },
      api: { toggleLike },
    });
  } catch (e) {
    console.warn('[App] Failed to apply initial fonts/skin:', e);
  }
};

onMounted(() => {
  applyInitialFontsAndSkin();

  try {
    EventsOn('open:admin', () => {
      isAdminOpen.value = true;
    });
    EventsOn('skin:updated', (payload: any) => {
      if (payload?.css) applyCSS(payload.css);
    });
  } catch (err) {
    console.warn('Wails EventsOn not available:', err);
  }

  observer = new IntersectionObserver((entries) => {
    if (entries[0].isIntersecting && hasMore.value && !loading.value && !activeArticleId.value) {
      loadMore();
    }
  }, { rootMargin: '200px' });
  if (observerTarget.value) observer.observe(observerTarget.value);
});

onUnmounted(() => {
  observer?.disconnect();
  try {
    EventsOff('open:admin');
    EventsOff('skin:updated');
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
        <!-- 1. 個別ツイート詳細ページ (Level 3 Detail Focus & Thread Tree) -->
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
            :focusedArticleId="focusedArticleId || undefined"
            @back="closeArticleDetail"
            @selectArticle="openArticleDetail"
            @toggleLike="toggleLike"
            @retryMedia="retryMedia"
            @clickTag="handleTagClick"
            @clickMention="handleMentionClick"
            @clickMedia="(media, list, art) => openMedia(media, list, art)"
          />
        </div>

        <!-- 2. 通常タイムライン表示 (Level 2 Filter + Level 3 List) -->
        <div v-else>
          <!-- Level 2: コンテンツ種別フィルター -->
          <TimelineFilter :currentFilter="currentFilter" @filter="setFilter" />

          <!-- アクティブ検索・タグ絞り込み表示バー -->
          <div
            v-if="searchQuery"
            class="px-4 py-2.5 bg-blue-950/40 border-b border-blue-900/50 flex items-center justify-between gap-3 text-xs"
          >
            <div class="flex items-center gap-2 min-w-0">
              <span class="text-blue-400 font-semibold flex items-center gap-1">
                <span>🔍</span> 絞り込み中:
              </span>
              <span class="px-2 py-0.5 rounded bg-blue-900/60 text-blue-200 font-mono font-medium truncate border border-blue-700/60">
                {{ searchQuery }}
              </span>
            </div>
            <button
              @click="clearSearchQuery"
              class="px-2 py-1 rounded bg-slate-800 hover:bg-slate-700 text-slate-300 hover:text-white border border-slate-700 text-[11px] font-medium transition-colors flex items-center gap-1 cursor-pointer flex-shrink-0"
              title="絞り込みを解除"
            >
              <span>✕</span>
              <span>解除</span>
            </button>
          </div>

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

          <!-- 0件表示 -->
          <div v-else-if="articles.length === 0" class="p-12 text-center text-slate-500 font-mono text-sm space-y-2">
            <div>No articles found.</div>
            <div class="text-xs text-slate-600">サルベージジョブを実行してデータを取得してください。</div>
          </div>

          <!-- 記事リスト -->
          <div v-else class="divide-y divide-slate-800">
            <ArticleCard
              v-for="article in articles"
              :key="article.id"
              :article="article"
              :targetLang="systemLang"
              :isFocused="focusedArticleId === article.id"
              @clickArticle="openArticleDetail"
              @toggleLike="toggleLike"
              @retryMedia="retryMedia"
              @clickMedia="(media, list, art) => openMedia(media, list, art)"
              @clickTag="handleTagClick"
              @clickMention="handleMentionClick"
            />
          </div>

          <!-- 無限スクロール用監視ターゲット ＆ ロード中インジケータ -->
          <div ref="observerTarget" class="py-6 flex justify-center items-center text-xs font-mono text-slate-500">
            <div v-if="loading && articles.length > 0" class="flex items-center gap-2">
              <span class="animate-spin text-blue-500">⏳</span>
              <span>Loading more articles...</span>
            </div>
            <span v-else-if="!hasMore && articles.length > 0" class="text-slate-600">
              — End of Timeline —
            </span>
          </div>
        </div>
      </main>
    </div>

    <!-- キーボードショートカット一覧モーダル -->
    <KeyboardShortcutModal
      :isOpen="isHelpOpen"
      @close="isHelpOpen = false"
    />

    <!-- 管理者モーダル -->
    <AdminModal
      :isOpen="isAdminOpen"
      @close="isAdminOpen = false"
      @whitelistUpdated="reloadAll"
    />

    <!-- メディア拡大オーバーレイ (画像/動画フルスクリーンビューア) -->
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
  </div>
</template>
