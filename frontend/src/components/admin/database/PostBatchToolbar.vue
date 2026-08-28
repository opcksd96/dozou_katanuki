<!-- frontend/src/components/admin/database/PostBatchToolbar.vue (100行以下 - SPEC-PRINCIPLE-001) -->
<script setup lang="ts">
defineProps<{
  selectedCount: number; canUndo: boolean; canRedo: boolean; includeTrash: boolean; total: number;
}>();

const emit = defineEmits<{
  (e: 'batchTrash'): void; (e: 'batchResetTranslations'): void;
  (e: 'undo'): void; (e: 'redo'): void; (e: 'update:includeTrash', val: boolean): void;
  (e: 'clearSelection'): void;
}>();
</script>

<template>
  <div class="flex flex-wrap items-center justify-between gap-2 bg-slate-900/70 border border-slate-800 p-2 rounded-xl text-xs select-none shadow-sm">
    <!-- 左側：選択ステータス ＆ バッチ操作 -->
    <div class="flex items-center gap-2 flex-wrap">
      <span class="text-slate-400 font-mono text-[11px] flex items-center gap-1.5">
        <span>全 <strong class="text-slate-200">{{ total }}</strong> 件</span>
        <span v-if="selectedCount > 0" class="text-blue-400 font-bold bg-blue-950/60 border border-blue-800/60 px-2 py-0.5 rounded-md">
          ✓ 選択: {{ selectedCount }} 件
        </span>
      </span>

      <template v-if="selectedCount > 0">
        <button
          @click="emit('batchTrash')"
          class="px-2.5 py-1 bg-rose-950/80 hover:bg-rose-900 border border-rose-700/60 text-rose-300 hover:text-white font-bold rounded-lg flex items-center gap-1 cursor-pointer transition-colors shadow"
          title="選択した記事を一括でゴミ箱へ移動"
        >
          <span>🗑️</span> 一括削除 ({{ selectedCount }})
        </button>
        <button
          @click="emit('batchResetTranslations')"
          class="px-2.5 py-1 bg-amber-950/80 hover:bg-amber-900 border border-amber-700/60 text-amber-300 hover:text-white font-bold rounded-lg flex items-center gap-1 cursor-pointer transition-colors shadow"
          title="選択した記事の翻訳データを初期化"
        >
          <span>🔄</span> 翻訳リセット ({{ selectedCount }})
        </button>
        <button
          @click="emit('clearSelection')"
          class="px-2 py-1 text-slate-400 hover:text-slate-200 text-[11px] cursor-pointer"
        >
          選択解除
        </button>
      </template>
    </div>

    <!-- 右側：Undo / Redo & ゴミ箱表示トグル -->
    <div class="flex items-center gap-2 flex-wrap">
      <div class="flex items-center gap-1 bg-slate-950 p-0.5 rounded-lg border border-slate-800">
        <button
          @click="emit('undo')"
          :disabled="!canUndo"
          class="px-2 py-1 rounded text-slate-300 hover:bg-slate-800 disabled:opacity-30 disabled:cursor-not-allowed cursor-pointer flex items-center gap-1 text-[11px]"
          title="直前の操作を取り消す (Ctrl+Z)"
        >
          <span>↩</span> アンドゥ
        </button>
        <button
          @click="emit('redo')"
          :disabled="!canRedo"
          class="px-2 py-1 rounded text-slate-300 hover:bg-slate-800 disabled:opacity-30 disabled:cursor-not-allowed cursor-pointer flex items-center gap-1 text-[11px]"
          title="取り消した操作をやり直す (Ctrl+Y)"
        >
          <span>↪</span> リドゥ
        </button>
      </div>

      <label class="flex items-center gap-1.5 px-2.5 py-1 bg-slate-950/80 border border-slate-800 rounded-lg text-slate-400 hover:text-slate-200 cursor-pointer select-none text-[11px]">
        <input
          type="checkbox"
          :checked="includeTrash"
          @change="emit('update:includeTrash', ($event.target as HTMLInputElement).checked)"
          class="rounded border-slate-700 bg-slate-900 text-rose-500 focus:ring-0 w-3.5 h-3.5"
        />
        <span :class="{ 'text-rose-300 font-bold': includeTrash }">🗑️ 削除済みも含める</span>
      </label>
    </div>
  </div>
</template>
