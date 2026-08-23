<!-- frontend/src/components/admin/database/MediaManagementView.vue (100行以下 - SPEC-PRINCIPLE-001) -->
<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue';
import { BrowserOpenURL } from '../../../../wailsjs/runtime/runtime';
import MediaCockpitHeader from './MediaCockpitHeader.vue';
import MediaCard from './MediaCard.vue';
import MediaTableView from './MediaTableView.vue';
import MediaPreviewModal from './MediaPreviewModal.vue';

const props = defineProps<{
  mediaItems: any[];
  total: number;
  stats?: { total_count: number; image_count: number; video_count: number };
  accounts: any[];
  accountFilter: string;
  statusFilter: string;
  typeFilter: 'all' | 'image' | 'video';
  loading: boolean;
  page: number;
  limit: number;
  config?: any;
  activeJob?: any;
}>();

const emit = defineEmits<{
  (e: 'fetch'): void;
  (e: 'update:accountFilter', v: string): void;
  (e: 'update:statusFilter', v: string): void;
  (e: 'update:typeFilter', v: 'all' | 'image' | 'video'): void;
  (e: 'update:page', p: number): void;
  (e: 'update:limit', l: number): void;
  (e: 'retryMedia', mediaId: string): void;
  (e: 'purgeMedia', mediaId: string): void;
  (e: 'purgeByStatus', status: string): void;
  (e: 'viewPost', articleId: string): void;
  (e: 'saveMetadata', payload: any): void;
  (e: 'startDownload'): void;
  (e: 'startPoll'): void;
  (e: 'requeueFailed'): void;
  (e: 'reconcileStash'): void;
  (e: 'mergeDuplicates'): void;
  (e: 'purgeLowRes'): void;
  (e: 'openExplorer', id: string): void;
  (e: 'openDefault', id: string): void;
  (e: 'toggleBookmark', id: string): void;
  (e: 'cancelJob', id: string): void;
}>();

const viewMode = ref<'large' | 'compact' | 'table'>('large');
const searchQuery = ref('');
const onlyBookmarked = ref(false);
const showMoreActions = ref(false);
const selectedIndex = ref<number | null>(null);
const jumpPageTop = ref<number>(props.page);
const jumpPageBottom = ref<number>(props.page);

watch(() => props.page, (newP) => { jumpPageTop.value = newP; jumpPageBottom.value = newP; });

const filteredItems = computed(() => {
  let list = props.mediaItems || [];
  if (onlyBookmarked.value) list = list.filter(m => m.is_bookmarked);
  if (searchQuery.value.trim()) {
    const q = searchQuery.value.trim().toLowerCase();
    list = list.filter(m => (m.media_id || m.id || '').toLowerCase().includes(q) || (m.username || '').toLowerCase().includes(q) || (m.download_url || '').toLowerCase().includes(q) || (m.full_text || '').toLowerCase().includes(q));
  }
  return list;
});

const totalPages = computed(() => Math.max(1, Math.ceil(props.total / props.limit)));
const startRange = computed(() => props.total === 0 ? 0 : (props.page - 1) * props.limit + 1);
const endRange = computed(() => Math.min(props.page * props.limit, props.total));
const selectedMedia = computed(() => selectedIndex.value !== null ? filteredItems.value[selectedIndex.value] : null);

const currentHostname = computed(() => (typeof window !== 'undefined' && window.location?.hostname) ? window.location.hostname : '127.0.0.1');
const stashPort = computed(() => props.config?.network?.stash_port || 9999);
const openStashWebUI = () => {
  const url = `http://${currentHostname.value}:${stashPort.value}`;
  try { BrowserOpenURL(url); } catch { window.open(url, '_blank', 'noopener,noreferrer'); }
};

const handleCopyMedia = async (media: any) => {
  const text = media.download_url || media.media_id || media.id;
  try { await navigator.clipboard.writeText(text); alert(`📋 クリップボードにコピーしました:\n${text}`); }
  catch { prompt('コピーしてください:', text); }
};

const onTypeChange = (val: 'all' | 'image' | 'video') => { emit('update:typeFilter', val); emit('update:page', 1); emit('fetch'); };
const onAccountChange = (val: string) => { emit('update:accountFilter', val); emit('update:page', 1); emit('fetch'); };
const onStatusChange = (val: string) => { emit('update:statusFilter', val); emit('update:page', 1); emit('fetch'); };
const onLimitChange = (val: number) => { emit('update:limit', val); emit('update:page', 1); emit('fetch'); };
const goToPage = (p: number) => { emit('update:page', Math.max(1, Math.min(p, totalPages.value))); emit('fetch'); };

onMounted(() => { emit('fetch'); });
</script>

<template>
  <div class="h-full flex flex-col min-h-0 space-y-1.5 overflow-hidden">
    <!-- 薄型アクティブジョブ進行状況バー (実行時のみ) -->
    <MediaCockpitHeader :active-job="activeJob" @cancel-job="(id) => emit('cancelJob', id)" />

    <!-- 🌟 スマート1段 統合コマンドバー (超薄型・全機能集約) -->
    <div class="shrink-0 flex flex-wrap items-center justify-between gap-1.5 text-xs bg-slate-900/95 px-3 py-1.5 rounded-xl border border-slate-800/80 shadow-md">
      <!-- 左ブロック: フィルタ ＆ 検索 -->
      <div class="flex flex-wrap items-center gap-1.5">
        <select :value="accountFilter" @change="onAccountChange(($event.target as HTMLSelectElement).value)" class="bg-slate-950 border border-slate-700 rounded-lg px-2 py-1 text-slate-200 font-mono text-xs">
          <option value="all">全アカウント (ALL)</option>
          <option v-for="acc in accounts" :key="acc.numeric_id" :value="acc.numeric_id">@{{ acc.username }}</option>
        </select>
        <div class="flex rounded-lg border border-slate-800 p-0.5 bg-slate-950 text-xs font-mono">
          <button @click="onTypeChange('all')" class="px-2 py-0.5 rounded transition-colors" :class="typeFilter === 'all' ? 'bg-blue-600 text-white font-bold' : 'text-slate-400 hover:text-slate-200'">全 ({{ stats?.total_count || total }})</button>
          <button @click="onTypeChange('image')" class="px-2 py-0.5 rounded transition-colors" :class="typeFilter === 'image' ? 'bg-blue-600 text-white font-bold' : 'text-slate-400 hover:text-slate-200'">🖼️ ({{ stats?.image_count ?? 0 }})</button>
          <button @click="onTypeChange('video')" class="px-2 py-0.5 rounded transition-colors" :class="typeFilter === 'video' ? 'bg-blue-600 text-white font-bold' : 'text-slate-400 hover:text-slate-200'">🎬 ({{ stats?.video_count ?? 0 }})</button>
        </div>
        <select :value="statusFilter" @change="onStatusChange(($event.target as HTMLSelectElement).value)" class="bg-slate-950 border border-slate-700 rounded-lg px-2 py-1 text-slate-200 font-mono text-xs">
          <option value="all">全ステータス</option>
          <option value="COMPLETED">COMPLETED (確保済)</option>
          <option value="QUEUED">QUEUED (待機中)</option>
          <option value="EXCLUDED">EXCLUDED (対象外)</option>
          <option value="DEAD_404">DEAD_404 (消失)</option>
        </select>
        <input v-model="searchQuery" type="text" placeholder="🔍 検索..." class="bg-slate-950 border border-slate-700 rounded-lg px-2 py-1 text-slate-200 font-mono text-xs w-28 sm:w-36" />
        <button @click="onlyBookmarked = !onlyBookmarked" class="px-2 py-1 rounded-lg border text-xs font-mono transition-colors font-bold" :class="onlyBookmarked ? 'bg-amber-500 text-slate-950 border-amber-400' : 'bg-slate-950 border-slate-700 text-slate-300'">⭐</button>
      </div>

      <!-- 右ブロック: コックピットアクション ＆ 表示切替 ＆ ページ送り -->
      <div class="flex flex-wrap items-center gap-1.5 font-mono text-xs">
        <!-- 指令ボタン群 -->
        <button @click="emit('startDownload')" class="px-2 py-1 bg-blue-600 hover:bg-blue-500 text-white rounded-lg font-bold flex items-center gap-1" title="QUEUEDメディアの一括ダウンロード開始">🚀 DL開始</button>
        <button @click="emit('reconcileStash')" class="px-2 py-1 bg-indigo-600 hover:bg-indigo-500 text-white rounded-lg flex items-center gap-1 font-bold" title="Stashと土蔵DBの逆引き同期">🎛️ 同期</button>
        <button @click="emit('startPoll')" class="px-2 py-1 bg-purple-600 hover:bg-purple-500 text-white rounded-lg flex items-center gap-1" title="Aria2委託完了メディアの回収">🔄 回収</button>
        <button @click="emit('requeueFailed')" class="px-2 py-1 bg-amber-600 hover:bg-amber-500 text-white rounded-lg flex items-center gap-1" title="失敗の一括リトライ">🔁 リトライ</button>
        
        <!-- その他メニュー (重複統合、低解像度パージ) -->
        <div class="relative">
          <button @click="showMoreActions = !showMoreActions" class="px-2 py-1 bg-slate-800 hover:bg-slate-700 text-slate-300 rounded-lg border border-slate-700">⋯ その他</button>
          <div v-if="showMoreActions" class="absolute right-0 top-full mt-1 bg-slate-900 border border-slate-700 rounded-xl p-1.5 shadow-2xl z-50 flex flex-col gap-1 w-36" @click="showMoreActions = false">
            <button @click="emit('mergeDuplicates')" class="px-2 py-1 bg-slate-800 hover:bg-slate-700 text-slate-200 rounded text-left text-xs">🧬 重複統合</button>
            <button @click="emit('purgeLowRes')" class="px-2 py-1 bg-slate-800 hover:bg-slate-700 text-amber-300 rounded text-left text-xs">🧹 低解像度パージ</button>
            <button @click="openStashWebUI" class="px-2 py-1 bg-purple-950 hover:bg-purple-900 text-purple-300 rounded text-left text-xs">🎛️ Stash WebUI ↗</button>
          </div>
        </div>

        <!-- 表示切替スイッチ -->
        <div class="flex rounded-lg border border-slate-800 p-0.5 bg-slate-950">
          <button @click="viewMode = 'large'" class="px-1.5 py-0.5 rounded text-xs" :class="viewMode === 'large' ? 'bg-slate-700 text-white font-bold' : 'text-slate-400'" title="大サムネイル">🖼️</button>
          <button @click="viewMode = 'compact'" class="px-1.5 py-0.5 rounded text-xs" :class="viewMode === 'compact' ? 'bg-slate-700 text-white font-bold' : 'text-slate-400'" title="中サムネイル">🔲</button>
          <button @click="viewMode = 'table'" class="px-1.5 py-0.5 rounded text-xs" :class="viewMode === 'table' ? 'bg-slate-700 text-white font-bold' : 'text-slate-400'" title="詳細リスト">📋</button>
        </div>

        <!-- ページ送り -->
        <div class="flex items-center gap-1 bg-slate-950 p-0.5 rounded-lg border border-slate-800 text-xs">
          <button @click="goToPage(1)" :disabled="page <= 1" class="px-1.5 py-0.5 bg-slate-800 hover:bg-slate-700 rounded disabled:opacity-30">«</button>
          <button @click="goToPage(page - 1)" :disabled="page <= 1" class="px-1.5 py-0.5 bg-slate-800 hover:bg-slate-700 rounded disabled:opacity-30">◀</button>
          <input type="number" min="1" :max="totalPages" v-model.number="jumpPageTop" @keyup.enter="goToPage(jumpPageTop)" class="w-12 px-1 py-0.5 bg-slate-900 border border-slate-700 rounded text-center text-slate-100 font-bold" />
          <span class="text-slate-500 text-xs">/{{ totalPages }}P</span>
          <button @click="goToPage(page + 1)" :disabled="page >= totalPages" class="px-1.5 py-0.5 bg-slate-800 hover:bg-slate-700 rounded disabled:opacity-30">▶</button>
          <button @click="goToPage(totalPages)" :disabled="page >= totalPages" class="px-1.5 py-0.5 bg-slate-800 hover:bg-slate-700 rounded disabled:opacity-30">»</button>
        </div>
      </div>
    </div>

    <!-- 🖼️ メインコンテンツエリア (画面の80%以上を占有する広大なスクロールビュー) -->
    <div class="flex-1 min-h-0 overflow-y-auto pr-1">
      <div v-if="loading" class="p-16 text-center text-slate-400 text-sm">読み込み中...</div>
      <div v-else-if="!filteredItems || filteredItems.length === 0" class="p-16 text-center text-slate-500 text-sm">条件に合致するメディアデータがありません</div>
      
      <!-- 1. 詳細表形式 (Table View) -->
      <MediaTableView
        v-else-if="viewMode === 'table'"
        :items="filteredItems"
        @select="(m) => { selectedIndex = filteredItems.indexOf(m); }"
        @retry="(id) => emit('retryMedia', id)"
        @purge="(id) => emit('purgeMedia', id)"
        @open-explorer="(id) => emit('openExplorer', id)"
        @open-default="(id) => emit('openDefault', id)"
        @copy="handleCopyMedia"
        @toggle-bookmark="(id) => emit('toggleBookmark', id)"
      />

      <!-- 2. グリッド形式 (Large / Compact) -->
      <div v-else :class="['grid gap-3.5', viewMode === 'compact' ? 'grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-6' : 'grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5']">
        <MediaCard
          v-for="(m, idx) in filteredItems"
          :key="m.media_id || m.id"
          :media="m"
          :compact="viewMode === 'compact'"
          @click="selectedIndex = idx"
          @retry="(id) => emit('retryMedia', id)"
          @purge="(id) => emit('purgeMedia', id)"
          @view-post="(artId) => emit('viewPost', artId)"
          @open-explorer="(id) => emit('openExplorer', id)"
          @open-default="(id) => emit('openDefault', id)"
          @copy="handleCopyMedia"
          @toggle-bookmark="(id) => emit('toggleBookmark', id)"
        />
      </div>
    </div>

    <!-- 下部ミニマル・フッターバー (超薄型) -->
    <div class="shrink-0 flex items-center justify-between gap-2 text-xs bg-slate-900/90 px-3 py-1.5 rounded-xl border border-slate-800/80 font-mono text-slate-300 shadow">
      <div class="flex items-center gap-2">
        <span class="text-slate-400 text-xs">表示件数:</span>
        <select :value="limit" @change="onLimitChange(Number(($event.target as HTMLSelectElement).value))" class="bg-slate-950 border border-slate-700 rounded px-2 py-0.5 text-slate-200 text-xs">
          <option :value="24">24件</option><option :value="48">48件</option><option :value="96">96件</option><option :value="200">200件</option>
        </select>
        <span><span class="text-blue-400 font-bold">{{ startRange }}-{{ endRange }}</span> / 全 {{ total }} 件</span>
      </div>

      <div class="flex items-center gap-1 bg-slate-950 p-0.5 rounded-lg border border-slate-800">
        <button @click="goToPage(1)" :disabled="page <= 1" class="px-2 py-0.5 bg-slate-800 hover:bg-slate-700 rounded disabled:opacity-30 text-xs font-bold">«</button>
        <button @click="goToPage(page - 1)" :disabled="page <= 1" class="px-2 py-0.5 bg-slate-800 hover:bg-slate-700 rounded disabled:opacity-30 text-xs font-bold">◀</button>
        <input type="number" min="1" :max="totalPages" v-model.number="jumpPageBottom" @keyup.enter="goToPage(jumpPageBottom)" class="w-12 px-1 py-0.5 bg-slate-900 border border-slate-700 rounded text-center text-slate-100 font-bold text-xs" />
        <span class="text-slate-500 text-xs px-0.5">/{{ totalPages }}P</span>
        <button @click="goToPage(jumpPageBottom)" class="px-2 py-0.5 bg-blue-600 hover:bg-blue-500 text-white rounded text-xs font-bold">移動</button>
        <button @click="goToPage(page + 1)" :disabled="page >= totalPages" class="px-2 py-0.5 bg-slate-800 hover:bg-slate-700 rounded disabled:opacity-30 text-xs font-bold">▶</button>
        <button @click="goToPage(totalPages)" :disabled="page >= totalPages" class="px-2 py-0.5 bg-slate-800 hover:bg-slate-700 rounded disabled:opacity-30 text-xs font-bold">»</button>
      </div>
    </div>

    <!-- 詳細インスペクタ＆プレビューモーダル -->
    <MediaPreviewModal
      v-if="selectedMedia"
      :media="selectedMedia"
      @close="selectedIndex = null"
      @save-metadata="(p) => emit('saveMetadata', p)"
      @retry="(id) => { emit('retryMedia', id); }"
      @purge="(id) => { emit('purgeMedia', id); selectedIndex = null; }"
      @view-post="(artId) => { emit('viewPost', artId); selectedIndex = null; }"
    />
  </div>
</template>
