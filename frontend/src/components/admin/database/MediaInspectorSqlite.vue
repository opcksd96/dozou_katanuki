<!-- frontend/src/components/admin/database/MediaInspectorSqlite.vue (100行以下) -->
<script setup lang="ts">
const props = defineProps<{
  media: any;
  editStatus: string;
  editReason: string;
}>();

const emit = defineEmits<{
  (e: 'update:editStatus', val: string): void;
  (e: 'update:editReason', val: string): void;
  (e: 'save'): void;
  (e: 'retry'): void;
  (e: 'purge'): void;
  (e: 'viewPost'): void;
}>();
</script>

<template>
  <div class="p-3 bg-slate-950/80 rounded-xl border border-slate-800 space-y-2.5">
    <div class="text-[10px] text-slate-400 font-bold uppercase tracking-wider">🗄️ SQLite レコード編集</div>
    
    <div class="space-y-1">
      <label class="text-[10px] text-slate-400">ステータス</label>
      <select :value="editStatus" @change="emit('update:editStatus', ($event.target as HTMLSelectElement).value)" class="w-full bg-slate-900 border border-slate-700 rounded px-2 py-1 text-slate-200 text-xs">
        <option value="COMPLETED">COMPLETED</option>
        <option value="QUEUED">QUEUED</option>
        <option value="FAILED">FAILED</option>
        <option value="DEAD_404">DEAD_404</option>
        <option value="EXCLUDED">EXCLUDED</option>
      </select>
    </div>

    <div class="space-y-1">
      <label class="text-[10px] text-slate-400">失敗/除外理由</label>
      <input :value="editReason" @input="emit('update:editReason', ($event.target as HTMLInputElement).value)" type="text" placeholder="理由を入力..." class="w-full bg-slate-900 border border-slate-700 rounded px-2 py-1 text-slate-200 text-[11px]" />
    </div>

    <div class="flex gap-2 pt-1">
      <button @click="emit('save')" class="flex-1 py-1.5 bg-blue-600 hover:bg-blue-500 text-white font-bold rounded text-xs">
        💾 SQLite 更新
      </button>
      <button @click="emit('retry')" class="px-2.5 py-1.5 bg-amber-600/80 hover:bg-amber-500 text-white font-bold rounded text-xs">
        🔄 リトライ
      </button>
      <button @click="emit('purge')" class="px-2.5 py-1.5 bg-red-800/80 hover:bg-red-700 text-white font-bold rounded text-xs">
        🗑️ 削除
      </button>
    </div>

    <button @click="emit('viewPost')" class="w-full mt-1 py-1 bg-slate-800 hover:bg-slate-700 text-slate-300 rounded text-[11px]">
      📄 該当ツイートを見る
    </button>
  </div>
</template>
