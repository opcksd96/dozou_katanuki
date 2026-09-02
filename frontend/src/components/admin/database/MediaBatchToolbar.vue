<!-- frontend/src/components/admin/database/MediaBatchToolbar.vue (100行以下 - SPEC-PRINCIPLE-001) -->
<script setup lang="ts">
defineProps<{
  selectedCount: number;
  isTrashView: boolean;
  total: number;
}>();

const emit = defineEmits<{
  (e: 'batchRetry'): void;
  (e: 'batchRevertToQueued'): void;
  (e: 'batchTrash'): void;
  (e: 'batchRestore'): void;
  (e: 'selectAll'): void;
  (e: 'clearSelection'): void;
}>();
</script>

<template>
  <div v-if="selectedCount > 0" class="flex flex-wrap items-center justify-between gap-2 bg-blue-950/40 border border-blue-800/60 p-2 rounded-xl text-xs select-none shadow-sm animate-fadeIn">
    <!-- 左側：選択数 ＆ 操作ボタン群 -->
    <div class="flex items-center gap-2 flex-wrap">
      <span class="text-blue-300 font-bold bg-blue-900/60 border border-blue-700/60 px-2.5 py-1 rounded-lg flex items-center gap-1">
        <span>✓ 選択:</span> <strong class="text-white">{{ selectedCount }}</strong> 件
      </span>

      <template v-if="!isTrashView">
        <button
          @click="emit('batchRetry')"
          class="px-2.5 py-1 bg-blue-950/80 hover:bg-blue-900 border border-blue-700/60 text-blue-300 hover:text-white font-bold rounded-lg flex items-center gap-1 cursor-pointer transition-colors shadow active:scale-95"
          title="選択したメディアの再ダウンロードを一括実行"
        >
          <span>🔄</span> 一括リトライ ({{ selectedCount }})
        </button>
        <button
          @click="emit('batchRevertToQueued')"
          class="px-2.5 py-1 bg-amber-950/80 hover:bg-amber-900 border border-amber-600/60 text-amber-300 hover:text-white font-bold rounded-lg flex items-center gap-1 cursor-pointer transition-colors shadow active:scale-95"
          title="選択したメディアを QUEUED に差し戻し"
        >
          <span>🧹</span> QUEUEDへ差し戻し ({{ selectedCount }})
        </button>
        <button
          @click="emit('batchTrash')"
          class="px-2.5 py-1 bg-rose-950/80 hover:bg-rose-900 border border-rose-700/60 text-rose-300 hover:text-white font-bold rounded-lg flex items-center gap-1 cursor-pointer transition-colors shadow active:scale-95"
          title="選択したメディアを一括でゴミ箱へ移動"
        >
          <span>🗑️</span> 一括ゴミ箱 ({{ selectedCount }})
        </button>
      </template>

      <template v-else>
        <button
          @click="emit('batchRestore')"
          class="px-2.5 py-1 bg-emerald-950/80 hover:bg-emerald-900 border border-emerald-700/60 text-emerald-300 hover:text-white font-bold rounded-lg flex items-center gap-1 cursor-pointer transition-colors shadow active:scale-95"
          title="選択したメディアを一括復元"
        >
          <span>♻️</span> 一括復元 ({{ selectedCount }})
        </button>
      </template>
    </div>

    <!-- 右側：選択解除 / 全選択 -->
    <div class="flex items-center gap-2">
      <button
        @click="emit('selectAll')"
        class="px-2.5 py-1 text-slate-300 hover:text-white bg-slate-900 hover:bg-slate-800 border border-slate-700 rounded-lg text-[11px] cursor-pointer"
      >
        表示中全選択
      </button>
      <button
        @click="emit('clearSelection')"
        class="px-2.5 py-1 text-slate-400 hover:text-slate-200 text-[11px] cursor-pointer"
      >
        選択解除
      </button>
    </div>
  </div>
</template>
