<!-- frontend/src/components/admin/database/MediaInspectorSqlite.vue (100行以下 - SPEC-PRINCIPLE-001) -->
<script setup lang="ts">
import { computed } from 'vue';

const props = defineProps<{ media: any; editStatus: string; editReason: string }>();

const emit = defineEmits<{
  (e: 'update:editStatus', val: string): void; (e: 'update:editReason', val: string): void;
  (e: 'save'): void; (e: 'retry'): void; (e: 'purge'): void;
  (e: 'viewPost'): void; (e: 'viewPostTimeline'): void; (e: 'openExplorer'): void; (e: 'openDefault'): void;
}>();

const isWails = computed(() => typeof window !== 'undefined' && !!((window as any)?.go?.app?.App || (window as any)?.go?.main?.App) && !(window as any)?._isWailsPolyfill);
const localPath = computed(() => props.media.file_path || props.media.local_path || '');
const tweetUrls = computed(() => {
  const v = props.media.tweet_urls;
  if (Array.isArray(v)) return v;
  if (typeof v === 'string' && v.startsWith('[')) {
    try { return JSON.parse(v); } catch { return []; }
  }
  return [];
});
</script>

<template>
  <div class="p-3 bg-slate-950/80 rounded-xl border border-slate-800 space-y-2.5">
    <div class="text-[10px] text-slate-400 font-bold uppercase tracking-wider">🗄️ SQLite レコード編集</div>
    
    <div v-if="localPath" class="space-y-1">
      <label class="text-[10px] text-slate-400">実体ファイルパス</label>
      <div class="p-1.5 bg-slate-900 border border-slate-700/80 rounded font-mono text-[11px] text-slate-300 select-all break-all leading-tight">{{ localPath }}</div>
      <div v-if="isWails" class="flex gap-2 pt-1">
        <button @click="emit('openExplorer')" class="flex-1 py-1 bg-slate-800 hover:bg-slate-700 text-slate-200 rounded text-[11px] font-bold active:scale-95 cursor-pointer">📁 Show in Folder</button>
        <button @click="emit('openDefault')" class="flex-1 py-1 bg-slate-800 hover:bg-slate-700 text-slate-200 rounded text-[11px] font-bold active:scale-95 cursor-pointer">▶️ Open</button>
      </div>
    </div>

    <div v-if="tweetUrls.length > 0" class="space-y-1">
      <label class="text-[10px] text-slate-400">参照ツイートURL ({{ tweetUrls.length }}件)</label>
      <div class="space-y-1 max-h-24 overflow-y-auto">
        <a v-for="u in tweetUrls" :key="u" :href="u" target="_blank" class="block p-1 bg-slate-900 border border-slate-800 rounded font-mono text-[10px] text-indigo-400 hover:underline truncate">{{ u }}</a>
      </div>
    </div>

    <div class="space-y-1">
      <label class="text-[10px] text-slate-400">ステータス</label>
      <select :value="editStatus" @change="emit('update:editStatus', ($event.target as HTMLSelectElement).value)" class="w-full bg-slate-900 border border-slate-700 rounded px-2 py-1 text-slate-200 text-xs cursor-pointer">
        <option value="COMPLETED">COMPLETED</option><option value="QUEUED">QUEUED</option><option value="OUTSOURCED">OUTSOURCED</option>
        <option value="RETAINED">RETAINED</option><option value="DEAD_404">DEAD_404</option><option value="FAILED">FAILED</option><option value="EXCLUDED">EXCLUDED</option>
      </select>
    </div>

    <div class="space-y-1">
      <label class="text-[10px] text-slate-400">失敗/除外理由</label>
      <input :value="editReason" @input="emit('update:editReason', ($event.target as HTMLInputElement).value)" type="text" placeholder="理由を入力..." class="w-full bg-slate-900 border border-slate-700 rounded px-2 py-1 text-slate-200 text-[11px]" />
    </div>

    <div class="flex gap-2 pt-1">
      <button @click="emit('save')" class="flex-1 py-1.5 bg-blue-600 hover:bg-blue-500 text-white font-bold rounded text-xs active:scale-95 cursor-pointer">💾 SQLite 更新</button>
      <button @click="emit('retry')" class="px-2.5 py-1.5 bg-amber-600 hover:bg-amber-500 text-white font-bold rounded text-xs active:scale-95 cursor-pointer">🔄 リトライ</button>
      <button @click="emit('purge')" class="px-2.5 py-1.5 bg-red-800 hover:bg-red-700 text-white font-bold rounded text-xs active:scale-95 cursor-pointer">🗑️ 削除</button>
    </div>

    <div class="flex gap-1.5 pt-1">
      <button @click="emit('viewPostTimeline')" class="flex-1 py-1.5 bg-gradient-to-r from-blue-600 to-indigo-600 hover:from-blue-500 hover:to-indigo-500 text-white font-bold rounded text-xs shadow-md active:scale-95 cursor-pointer" title="メインタイムラインで詳細展開">📱 タイムラインで見る</button>
      <button @click="emit('viewPost')" class="px-2.5 py-1.5 bg-slate-800 hover:bg-slate-700 text-slate-300 rounded text-[11px] active:scale-95 cursor-pointer" title="管理画面の投稿タブで開く">📄 投稿タブ</button>
    </div>
  </div>
</template>
