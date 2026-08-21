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
  (e: 'update:page', p: number): void; (e: 'search'): void; (e: 'select', art: models.RenderTree): void;
  (e: 'save', ja: string, en: string, zh: string): void; (e: 'autoTranslate', autoSave: boolean): void;
  (e: 'batchTranslate'): void; (e: 'cancelJob', id: string): void;
}>();

const selectAcc = (id: string) => { emit('update:searchAccount', id); emit('update:page', 1); emit('search'); };
const cols = [
  { key: 'id', label: 'ID', width: '110px' }, { key: 'author_handle', label: 'Author', width: '100px' },
  { key: 'created_at', label: 'Created At', width: '130px' }, { key: 'orig_text', label: 'Original (FullText)', width: '220px' },
  { key: 'ja_text', label: 'JA (日本語)', width: '200px' }, { key: 'en_text', label: 'EN (英語)', width: '200px' }, { key: 'zh_text', label: 'ZH (中国語)', width: '200px' },
];

const spreadsheetRows = computed(() => (props.articles || []).map((a: any) => ({
  id: a.id, author_handle: `@${a.author?.handle || a.author?.username || ''}`,
  created_at: (a.created_at || a.createdAt || '').slice(0, 19).replace('T', ' '),
  orig_text: a.content?.original || '', ja_text: a.content?.ja || '',
  en_text: a.content?.en || '', zh_text: a.content?.zh || '', _raw: a,
})));
</script>

<template>
  <div class="space-y-3 flex flex-col h-full">
    <!-- 1. アカウント切替バー（折り返し全表示） -->
    <div class="flex flex-wrap items-center gap-1.5 text-xs bg-slate-950/40 p-2 rounded-xl border border-slate-800/80">
      <span class="text-[11px] font-bold text-slate-400 mr-1">アカウント:</span>
      <button @click="selectAcc('all')" class="px-2.5 py-1 rounded-lg font-bold border transition-colors whitespace-nowrap text-xs" :class="searchAccount === 'all' ? 'bg-blue-600 border-blue-500 text-white shadow' : 'bg-slate-900 border-slate-700 text-slate-400 hover:text-slate-200'">ALL (全アカウント)</button>
      <button v-for="acc in accounts" :key="acc.numeric_id" @click="selectAcc(acc.numeric_id)" class="px-2.5 py-1 rounded-lg font-mono border transition-colors whitespace-nowrap text-xs" :class="searchAccount === acc.numeric_id ? 'bg-blue-600 border-blue-500 text-white shadow' : 'bg-slate-900 border-slate-700 text-slate-400 hover:text-slate-200'">@{{ acc.username }}</button>
    </div>

    <!-- 2. 検索バー ＆ ページネーション -->
    <div class="flex gap-2 items-center">
      <input :value="searchQuery" @input="emit('update:searchQuery', ($event.target as HTMLInputElement).value)" @keyup.enter="emit('update:page', 1); emit('search')" type="text" placeholder="ID・本文・翻訳文を検索..." class="flex-1 bg-slate-950 border border-slate-700 rounded-lg px-3 py-1.5 text-xs text-slate-200" />
      <button @click="emit('update:page', 1); emit('search')" class="px-4 py-1.5 bg-blue-600 hover:bg-blue-500 text-white text-xs font-bold rounded-lg">{{ loading ? '⏳' : '検索' }}</button>
      <div class="text-[11px] text-slate-400 font-mono flex items-center gap-2">
        <span>全 {{ total }} 件</span>
        <button @click="if (page > 1) { emit('update:page', page - 1); emit('search'); }" :disabled="page <= 1" class="px-2 py-1 bg-slate-800 rounded disabled:opacity-40">◀</button>
        <span>{{ page }}</span>
        <button @click="if (page * limit < total) { emit('update:page', page + 1); emit('search'); }" :disabled="page * limit >= total" class="px-2 py-1 bg-slate-800 rounded disabled:opacity-40">▶</button>
      </div>
    </div>

    <!-- 3. メインエリア: 左側スプレッドシート / 右側翻訳パネル -->
    <div class="grid grid-cols-1 lg:grid-cols-12 gap-4 flex-1 min-h-[400px]">
      <div class="lg:col-span-7 h-full min-h-[300px]">
        <DatabaseSpreadsheet :columns="cols" :rows="spreadsheetRows" :selected-row-id="selectedArticle?.id" @select-row="(r) => emit('select', r._raw)" />
      </div>

      <!-- 右側：翻訳カラム（上部に一括翻訳バー、下部に個別エディタ） -->
      <div class="lg:col-span-5 h-full flex flex-col gap-3">
        <!-- 一括自動翻訳アクションカード -->
        <div class="bg-gradient-to-r from-indigo-950/60 to-slate-900/80 border border-indigo-500/30 rounded-xl p-3 space-y-2 shadow">
          <div class="flex items-center justify-between">
            <span class="text-xs font-bold text-indigo-200 flex items-center gap-1.5"><span>🌐</span> 一括自動翻訳 (Batch)</span>
            <span class="text-[10px] px-2 py-0.5 rounded bg-indigo-900/60 text-indigo-300 font-mono border border-indigo-700/50">対象: {{ searchAccount === 'all' ? '全アカウント' : '選択アカウント' }}</span>
          </div>
          <!-- 実行中プログレスバー -->
          <div v-if="activeJob && activeJob.type === 'translate' && activeJob.status?.toUpperCase() === 'RUNNING'" class="space-y-1.5 bg-slate-950/70 p-2.5 rounded-lg border border-indigo-500/50">
            <div class="flex justify-between text-[11px] text-indigo-200 font-bold">
              <span class="animate-pulse flex items-center gap-1"><span>⏳</span> 翻訳処理中...</span>
              <span class="font-mono">{{ Math.round(activeJob.percentage || 0) }}% ({{ activeJob.current }}/{{ activeJob.total }})</span>
            </div>
            <div class="w-full bg-slate-800 rounded-full h-2 overflow-hidden">
              <div class="bg-indigo-500 h-2 rounded-full transition-all duration-300 shadow-lg" :style="{ width: `${activeJob.percentage || 0}%` }"></div>
            </div>
            <div class="flex justify-between items-center text-[10px] text-slate-400 pt-0.5">
              <span class="truncate max-w-[220px]">{{ activeJob.message }}</span>
              <button @click="emit('cancelJob', activeJob.id)" class="text-rose-400 hover:text-rose-300 px-2 py-0.5 rounded bg-rose-950/60 border border-rose-800/50 text-[10px]">✕ 中断</button>
            </div>
          </div>
          <!-- 通常時の説明＆実行ボタン -->
          <template v-else>
            <p class="text-[10px] text-slate-400 leading-tight">※ 未翻訳記事を最大 <strong class="text-indigo-300">500件</strong> ずつ自動取得しDBに順次保存します。</p>
            <button @click="emit('batchTranslate')" class="w-full py-1.5 bg-indigo-600 hover:bg-indigo-500 text-white text-xs font-bold rounded-lg shadow flex items-center justify-center gap-1.5 transition-colors">
              <span>🚀</span> 未翻訳を一括自動翻訳実行 (最大500件)
            </button>
          </template>
        </div>

        <!-- 個別記事翻訳エディタ -->
        <div class="flex-1 min-h-[300px]">
          <DatabaseTranslationEditor :article="selectedArticle" :saving="saving" :translating="translating" @auto-translate="(autoSave) => emit('autoTranslate', autoSave)" @save="(j, e, z) => emit('save', j, e, z)" />
        </div>
      </div>
    </div>
  </div>
</template>
