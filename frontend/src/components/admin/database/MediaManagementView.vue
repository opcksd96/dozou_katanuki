<!-- frontend/src/components/admin/database/MediaManagementView.vue (100行以下) -->
<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { BrowserOpenURL } from '../../../../wailsjs/runtime/runtime';
import MediaCard from './MediaCard.vue';
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
}>();

const emit = defineEmits<{
  (e: 'fetch'): void;
  (e: 'update:accountFilter', v: string): void;
  (e: 'update:statusFilter', v: string): void;
  (e: 'update:typeFilter', v: 'all' | 'image' | 'video'): void;
  (e: 'update:page', p: number): void;
  (e: 'retryMedia', mediaId: string): void;
  (e: 'purgeMedia', mediaId: string): void;
  (e: 'purgeByStatus', status: string): void;
  (e: 'viewPost', articleId: string): void;
  (e: 'saveMetadata', payload: any): void;
}>();

const selectedIndex = ref<number | null>(null);
const jumpPageInput = ref<number>(props.page);

const totalPages = computed(() => Math.max(1, Math.ceil(props.total / props.limit)));
const startRange = computed(() => props.total === 0 ? 0 : (props.page - 1) * props.limit + 1);
const endRange = computed(() => Math.min(props.page * props.limit, props.total));
const selectedMedia = computed(() => selectedIndex.value !== null ? props.mediaItems[selectedIndex.value] : null);

const currentHostname = computed(() => (typeof window !== 'undefined' && window.location?.hostname) ? window.location.hostname : '127.0.0.1');
const stashPort = computed(() => props.config?.network?.stash_port || 9999);
const openStashWebUI = () => {
  const url = `http://${currentHostname.value}:${stashPort.value}`;
  try { BrowserOpenURL(url); } catch { window.open(url, '_blank', 'noopener,noreferrer'); }
};

const onTypeChange = (val: 'all' | 'image' | 'video') => { emit('update:typeFilter', val); emit('update:page', 1); emit('fetch'); };
const onAccountChange = (val: string) => { emit('update:accountFilter', val); emit('update:page', 1); emit('fetch'); };
const onStatusChange = (val: string) => { emit('update:statusFilter', val); emit('update:page', 1); emit('fetch'); };
const handleJumpPage = () => {
  const p = Math.max(1, Math.min(Number(jumpPageInput.value) || 1, totalPages.value));
  emit('update:page', p); emit('fetch');
};

const handleConfirmPurgeStatus = (st: string, label: string) => {
  if (confirm(`本当に [${label}] のメディアレコードをDBから一括パージしますか？\n（※物理削除されます）`)) {
    emit('purgeByStatus', st);
  }
};

onMounted(() => { emit('fetch'); });
</script>

<template>
  <div class="space-y-2.5 flex flex-col h-full">
    <!-- ツールバー -->
    <div class="flex flex-wrap items-center justify-between gap-2 text-xs">
      <div class="flex flex-wrap items-center gap-1.5">
        <select :value="accountFilter" @change="onAccountChange(($event.target as HTMLSelectElement).value)" class="bg-slate-950 border border-slate-700 rounded-lg px-2 py-1 text-slate-200 font-mono text-[11px]">
          <option value="all">全アカウント (ALL)</option>
          <option v-for="acc in accounts" :key="acc.numeric_id" :value="acc.numeric_id">@{{ acc.username }}</option>
        </select>
        <div class="flex rounded-lg border border-slate-800 p-0.5 bg-slate-950 text-[11px] font-mono">
          <button @click="onTypeChange('all')" class="px-2 py-0.5 rounded transition-colors" :class="typeFilter === 'all' ? 'bg-blue-600 text-white font-bold' : 'text-slate-400 hover:text-slate-200'">
            全種別 ({{ stats?.total_count || total }})
          </button>
          <button @click="onTypeChange('image')" class="px-2 py-0.5 rounded transition-colors" :class="typeFilter === 'image' ? 'bg-blue-600 text-white font-bold' : 'text-slate-400 hover:text-slate-200'">
            🖼️ 静止画 ({{ stats?.image_count ?? 0 }})
          </button>
          <button @click="onTypeChange('video')" class="px-2 py-0.5 rounded transition-colors" :class="typeFilter === 'video' ? 'bg-blue-600 text-white font-bold' : 'text-slate-400 hover:text-slate-200'">
            🎬 動画 ({{ stats?.video_count ?? 0 }})
          </button>
        </div>
        <select :value="statusFilter" @change="onStatusChange(($event.target as HTMLSelectElement).value)" class="bg-slate-950 border border-slate-700 rounded-lg px-2 py-1 text-slate-200 font-mono text-[11px]">
          <option value="all">全ステータス</option>
          <option value="COMPLETED">COMPLETED (確保済)</option>
          <option value="EXCLUDED">EXCLUDED (対象外)</option>
          <option value="QUEUED">QUEUED (待機中)</option>
          <option value="DEAD_404">DEAD_404 (消失)</option>
        </select>
        <!-- 一括パージドロップダウン/ボタン -->
        <button @click="handleConfirmPurgeStatus('EXCLUDED', 'EXCLUDED (対象外)')" class="px-2 py-1 bg-slate-900 hover:bg-slate-800 border border-slate-700 text-slate-300 rounded-lg text-[10px] font-mono transition-colors" title="Whitelist外のEXCLUDEDメディアをDBからパージ">
          🗑️ 対象外一括パージ
        </button>
        <!-- Stash WebUI クイックアクセス -->
        <button @click="openStashWebUI" class="px-2 py-1 bg-purple-950/70 hover:bg-purple-900/90 border border-purple-700/60 text-purple-300 rounded-lg text-[10px] font-mono transition-colors flex items-center gap-1" title="Stash WebUI (ポート9999) をブラウザで開く">
          🎛️ Stash WebUI ↗
        </button>
      </div>

      <!-- ページネーション -->
      <div class="text-[11px] text-slate-300 font-mono flex items-center gap-2">
        <span><span class="text-blue-400 font-bold">{{ startRange }}-{{ endRange }}</span> / 全 {{ total }} 件 ({{ page }}/{{ totalPages }}P)</span>
        <button @click="if (page > 1) { emit('update:page', page - 1); emit('fetch'); }" :disabled="page <= 1" class="px-2 py-0.5 bg-slate-800 hover:bg-slate-700 rounded disabled:opacity-40">◀</button>
        <button @click="if (page < totalPages) { emit('update:page', page + 1); emit('fetch'); }" :disabled="page >= totalPages" class="px-2 py-0.5 bg-slate-800 hover:bg-slate-700 rounded disabled:opacity-40">▶</button>
        <div class="flex items-center gap-1">
          <input type="number" min="1" :max="totalPages" v-model.number="jumpPageInput" @keyup.enter="handleJumpPage" class="w-12 px-1 py-0.5 bg-slate-950 border border-slate-700 rounded text-center text-slate-200 text-xs" />
          <button @click="handleJumpPage" class="px-1.5 py-0.5 bg-slate-800 hover:bg-slate-700 text-[10px] rounded text-slate-300">移動</button>
        </div>
      </div>
    </div>

    <!-- メディアカードグリッド -->
    <div class="flex-1 overflow-y-auto max-h-[440px] pr-1">
      <div v-if="loading" class="p-12 text-center text-slate-400 text-xs">読み込み中...</div>
      <div v-else-if="!mediaItems || mediaItems.length === 0" class="p-12 text-center text-slate-500 text-xs">メディアデータがありません</div>
      <div v-else class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-3">
        <MediaCard
          v-for="(m, idx) in mediaItems"
          :key="m.media_id || m.id"
          :media="m"
          @click="selectedIndex = idx"
          @retry="(id) => emit('retryMedia', id)"
          @purge="(id) => emit('purgeMedia', id)"
          @view-post="(artId) => emit('viewPost', artId)"
        />
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
