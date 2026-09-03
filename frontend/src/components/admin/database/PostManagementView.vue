<!-- frontend/src/components/admin/database/PostManagementView.vue (100行以下 - SPEC-PRINCIPLE-001) -->
<script setup lang="ts">
import { ref, computed } from 'vue';
import { models } from '../../../../wailsjs/go/models';
import DatabaseSpreadsheet from './DatabaseSpreadsheet.vue';
import DatabaseTranslationEditor from './DatabaseTranslationEditor.vue';
import PostBatchToolbar from './PostBatchToolbar.vue';
import PostTrashModal from './PostTrashModal.vue';
import { useArticleBatchOps } from '../../../composables/admin/useArticleBatchOps';

const props = defineProps<{
  articles: models.RenderTree[]; total: number; selectedArticle: models.RenderTree | null;
  accounts: any[]; searchAccount: string; searchQuery: string; loading: boolean; saving: boolean;
  translating?: boolean; activeJob?: any; canUndo?: boolean; canRedo?: boolean; includeTrash?: boolean;
}>();

const emit = defineEmits<{
  (e: 'update:searchAccount', v: string): void; (e: 'update:searchQuery', v: string): void;
  (e: 'selectAccount', id: string): void; (e: 'search'): void; (e: 'fetchMore'): void;
  (e: 'select', art: models.RenderTree): void; (e: 'save', ja: string, en: string, zh: string): void;
  (e: 'autoTranslate', autoSave: boolean): void; (e: 'batchTranslate'): void; (e: 'cancelJob', id: string): void;
  (e: 'trash', id: string, reason: string): void; (e: 'batchTrash', ids: string[], reason: string): void;
  (e: 'batchResetTranslations', ids: string[]): void; (e: 'undo'): void; (e: 'redo'): void;
  (e: 'update:includeTrash', val: boolean): void;
}>();

const { selectedIds, selectedCount, toggleSelect, selectAll, clearSelection } = useArticleBatchOps();
const showBatchTrashModal = ref(false);

const cols = [
  { key: 'id', label: 'ID', width: '105px' }, { key: 'author_handle', label: 'Author', width: '90px' },
  { key: 'source_domain', label: 'Domain', width: '85px' }, { key: 'original_url', label: 'Original URL', width: '160px' },
  { key: 'wayback_url', label: 'Wayback URL', width: '160px' }, { key: 'sotwe_url', label: 'Sotwe URL', width: '140px' },
  { key: 'nitter_url', label: 'Nitter URL', width: '140px' }, { key: 'twistalker_url', label: 'Twistalker URL', width: '140px' },
  { key: 'created_at', label: 'Created At', width: '130px' }, { key: 'ja_text', label: 'JA (日本語)', width: '180px' },
  { key: 'orig_text', label: 'Original', width: '180px' },
];

const spreadsheetRows = computed(() => (props.articles || []).filter(Boolean).map((a: any) => ({
  id: a?.id || '', author_handle: `@${a?.author?.handle || a?.author?.username || ''}`,
  source_domain: a?.source_domain || (a?.source_url?.includes('x.com') ? 'x.com' : 'twitter.com'),
  original_url: a?.original_url || '', wayback_url: a?.wayback_url || a?.source_url || '',
  sotwe_url: a?.sotwe_url || '', nitter_url: a?.nitter_url || '', twistalker_url: a?.twistalker_url || '',
  created_at: String(a?.created_at || a?.createdAt || '').slice(0, 19).replace('T', ' '),
  orig_text: a?.content?.original || '', ja_text: a?.content?.ja || '', en_text: a?.content?.en || '', zh_text: a?.content?.zh || '',
  is_trash: !!a?.is_trash, _raw: a,
})));

const handleBatchTrashConfirm = (reason: string) => { emit('batchTrash', Array.from(selectedIds.value), reason); clearSelection(); };
</script>

<template>
  <div class="space-y-2 flex flex-col h-full text-xs">
    <div class="flex items-center gap-2 bg-slate-950/40 p-2 rounded-xl border border-slate-800 shrink-0">
      <span class="text-[11px] font-bold text-slate-400">アカウント:</span>
      <select :value="searchAccount" @change="emit('selectAccount', ($event.target as HTMLSelectElement).value)" class="flex-1 min-w-0 bg-slate-900 border border-slate-700 rounded-lg px-2 py-1.5 text-slate-200 text-xs font-mono focus:outline-none focus:border-blue-500 cursor-pointer">
        <option value="all">ALL (@全体)</option>
        <option v-for="acc in accounts" :key="acc.numeric_id" :value="acc.numeric_id">
          @{{ acc.username || acc.handle || acc.numeric_id }}
        </option>
      </select>
    </div>

    <!-- 検索バー -->
    <div class="flex gap-2 items-center">
      <input :value="searchQuery" @input="emit('update:searchQuery', ($event.target as HTMLInputElement).value)" @keyup.enter="emit('search')" type="text" placeholder="ID・本文を検索... (↑↓キーで行移動)" class="flex-1 bg-slate-950 border border-slate-700 rounded px-2.5 py-1.5 text-slate-200" />
      <button @click="emit('search')" class="px-3.5 py-1.5 bg-blue-600 hover:bg-blue-500 text-white font-bold rounded cursor-pointer">{{ loading ? '⏳' : '検索' }}</button>
    </div>

    <!-- バッチツールバー -->
    <PostBatchToolbar :selected-count="selectedCount" :can-undo="!!canUndo" :can-redo="!!canRedo" :include-trash="!!includeTrash" :total="total" @batch-trash="showBatchTrashModal = true" @batch-reset-translations="() => { emit('batchResetTranslations', Array.from(selectedIds)); clearSelection(); }" @undo="emit('undo')" @redo="emit('redo')" @update:include-trash="(v) => emit('update:includeTrash', v)" @clear-selection="clearSelection" />

    <!-- メインエリア -->
    <div class="grid grid-cols-1 lg:grid-cols-12 gap-3 flex-1 min-h-0 overflow-hidden">
      <div class="lg:col-span-7 h-full min-h-0 flex flex-col">
        <DatabaseSpreadsheet :columns="cols" :rows="spreadsheetRows" :selected-row-id="selectedArticle?.id" :selected-ids="selectedIds" @select-row="(r) => emit('select', r._raw)" @toggle-select="toggleSelect" @toggle-select-all="() => selectedIds.size === spreadsheetRows.length ? clearSelection() : selectAll(spreadsheetRows)" @scroll-bottom="emit('fetchMore')" />
      </div>
      <div class="lg:col-span-5 h-full flex flex-col gap-2.5 min-h-0 overflow-y-auto">
        <div class="bg-indigo-950/60 border border-indigo-500/30 rounded-xl p-2.5 space-y-1.5 shadow">
          <div class="flex items-center justify-between font-bold text-indigo-200 text-[11px]">
            <span>🌐 一括自動翻訳</span>
            <span class="text-[10px] px-2 py-0.5 rounded bg-indigo-900 text-indigo-300 font-mono">対象: {{ searchAccount === 'all' ? '全件' : '@' + searchAccount }}</span>
          </div>
          <button @click="emit('batchTranslate')" class="w-full py-1 bg-indigo-600 hover:bg-indigo-500 text-white font-bold rounded shadow cursor-pointer text-xs">🚀 未翻訳を一括自動翻訳 (最大500件)</button>
        </div>
        <div class="flex-1 min-h-[300px]">
          <DatabaseTranslationEditor :article="selectedArticle" :saving="saving" :translating="translating" @auto-translate="(autoSave) => emit('autoTranslate', autoSave)" @save="(j, e, z) => emit('save', j, e, z)" @trash="(id, r) => emit('trash', id, r)" />
        </div>
      </div>
    </div>

    <!-- 一括削除モーダル -->
    <PostTrashModal :show="showBatchTrashModal" :article-id="`${selectedCount} 件の記事`" @close="showBatchTrashModal = false" @confirm="handleBatchTrashConfirm" />
  </div>
</template>
