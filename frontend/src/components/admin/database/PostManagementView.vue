<!-- frontend/src/components/admin/database/PostManagementView.vue (100行以下) -->
<script setup lang="ts">
import { computed } from 'vue';
import { models } from '../../../../wailsjs/go/models';
import DatabaseSpreadsheet from './DatabaseSpreadsheet.vue';
import DatabaseTranslationEditor from './DatabaseTranslationEditor.vue';

const props = defineProps<{
  articles: models.RenderTree[]; total: number; selectedArticle: models.RenderTree | null;
  accounts: any[]; searchAccount: string; searchQuery: string; page: number; limit: number;
  loading: boolean; saving: boolean; translating?: boolean; activeJob?: any;
}>();

const emit = defineEmits<{
  (e: 'update:searchAccount', v: string): void; (e: 'update:searchQuery', v: string): void;
  (e: 'update:page', p: number): void; (e: 'changePage', p: number): void;
  (e: 'selectAccount', id: string): void;
  (e: 'search'): void; (e: 'select', art: models.RenderTree): void;
  (e: 'save', ja: string, en: string, zh: string): void; (e: 'autoTranslate', autoSave: boolean): void;
  (e: 'batchTranslate'): void; (e: 'cancelJob', id: string): void;
}>();

const selectAcc = (id: string) => {
  emit('selectAccount', id);
};

const handlePageNav = (newPage: number) => {
  emit('changePage', newPage);
};

const handleSearchSubmit = () => {
  emit('changePage', 1);
};

const cols = [
  { key: 'id', label: 'ID', width: '110px' }, { key: 'author_handle', label: 'Author', width: '100px' },
  { key: 'created_at', label: 'Created At', width: '130px' }, { key: 'orig_text', label: 'Original', width: '220px' },
  { key: 'ja_text', label: 'JA (日本語)', width: '200px' }, { key: 'en_text', label: 'EN', width: '200px' }, { key: 'zh_text', label: 'ZH', width: '200px' },
];

const spreadsheetRows = computed(() => (props.articles || []).filter(Boolean).map((a: any) => ({
  id: a?.id || '',
  author_handle: `@${a?.author?.handle || a?.author?.username || ''}`,
  created_at: String(a?.created_at || a?.createdAt || '').slice(0, 19).replace('T', ' '),
  orig_text: a?.content?.original || '',
  ja_text: a?.content?.ja || '',
  en_text: a?.content?.en || '',
  zh_text: a?.content?.zh || '',
  _raw: a,
})));

</script>

<template>
  <div class="space-y-3 flex flex-col h-full text-xs">
    <!-- アカウント切替バー -->
    <div class="flex flex-wrap items-center gap-1.5 bg-slate-950/40 p-2 rounded-xl border border-slate-800">
      <span class="text-[11px] font-bold text-slate-400 mr-1">アカウント:</span>
      <button @click="selectAcc('all')" :class="searchAccount === 'all' ? 'bg-blue-600 text-white' : 'bg-slate-900 text-slate-400'" class="px-2.5 py-1 rounded-lg border border-slate-700 font-bold cursor-pointer">ALL</button>
      <button v-for="acc in accounts" :key="acc.numeric_id" @click="selectAcc(acc.numeric_id)" :class="searchAccount === acc.numeric_id ? 'bg-blue-600 text-white' : 'bg-slate-900 text-slate-400'" class="px-2.5 py-1 rounded-lg border border-slate-700 font-mono cursor-pointer">@{{ acc.username }}</button>
    </div>

    <!-- 検索バー ＆ ページネーション -->
    <div class="flex gap-2 items-center">
      <input :value="searchQuery" @input="emit('update:searchQuery', ($event.target as HTMLInputElement).value)" @keyup.enter="handleSearchSubmit" type="text" placeholder="ID・本文を検索..." class="flex-1 bg-slate-950 border border-slate-700 rounded px-2.5 py-1.5 text-slate-200" />
      <button @click="handleSearchSubmit" class="px-3.5 py-1.5 bg-blue-600 hover:bg-blue-500 text-white font-bold rounded cursor-pointer">{{ loading ? '⏳' : '検索' }}</button>
      <div class="text-[11px] text-slate-400 font-mono flex items-center gap-2">
        <span>全 {{ total }} 件</span>
        <button @click="handlePageNav(page - 1)" :disabled="page <= 1" class="px-2 py-1 bg-slate-800 rounded disabled:opacity-40 hover:bg-slate-700 cursor-pointer">◀</button>
        <span>{{ page }} / {{ Math.max(1, Math.ceil(total / limit)) }}</span>
        <button @click="handlePageNav(page + 1)" :disabled="page * limit >= total" class="px-2 py-1 bg-slate-800 rounded disabled:opacity-40 hover:bg-slate-700 cursor-pointer">▶</button>
      </div>
    </div>

    <!-- メインエリア -->
    <div class="grid grid-cols-1 lg:grid-cols-12 gap-4 flex-1 min-h-0 overflow-hidden">
      <div class="lg:col-span-7 h-full min-h-0 flex flex-col">
        <DatabaseSpreadsheet :columns="cols" :rows="spreadsheetRows" :selected-row-id="selectedArticle?.id" @select-row="(r) => emit('select', r._raw)" />
      </div>
      <div class="lg:col-span-5 h-full flex flex-col gap-3 min-h-0 overflow-y-auto">
        <div class="bg-indigo-950/60 border border-indigo-500/30 rounded-xl p-3 space-y-2 shadow">
          <div class="flex items-center justify-between font-bold text-indigo-200">
            <span>🌐 一括自動翻訳</span>
            <span class="text-[10px] px-2 py-0.5 rounded bg-indigo-900 text-indigo-300 font-mono">対象: {{ searchAccount === 'all' ? '全件' : '@' + searchAccount }}</span>
          </div>
          <div v-if="activeJob?.type === 'translate' && activeJob.status === 'RUNNING'" class="space-y-1 bg-slate-950 p-2 rounded">
            <div class="flex justify-between font-bold text-indigo-200"><span>⏳ 処理中...</span><span>{{ Math.round(activeJob.percentage || 0) }}%</span></div>
            <div class="w-full bg-slate-800 rounded-full h-1.5"><div class="bg-indigo-500 h-1.5 rounded-full" :style="{ width: `${activeJob.percentage || 0}%` }"></div></div>
          </div>
          <button v-else @click="emit('batchTranslate')" class="w-full py-1.5 bg-indigo-600 hover:bg-indigo-500 text-white font-bold rounded shadow cursor-pointer">🚀 未翻訳を一括自動翻訳 (最大500件)</button>
        </div>
        <div class="flex-1 min-h-[300px]">
          <DatabaseTranslationEditor :article="selectedArticle" :saving="saving" :translating="translating" @auto-translate="(autoSave) => emit('autoTranslate', autoSave)" @save="(j, e, z) => emit('save', j, e, z)" />
        </div>
      </div>
    </div>
  </div>
</template>

