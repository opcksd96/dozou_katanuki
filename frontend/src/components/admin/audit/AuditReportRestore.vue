<!-- frontend/src/components/admin/audit/AuditReportRestore.vue (100行以下) -->
<script setup lang="ts">
import { ref } from 'vue';
defineProps<{
  restoring?: boolean;
  restoreStatus?: { success: boolean; message: string } | null;
}>();
const emit = defineEmits<{ (e: 'triggerRestore', resetDB: boolean): void }>();
const resetDB = ref(false);

const handleRestore = () => {
  if (confirm('Layer 2 (dumps/) から全オフラインデータを再構築・リストアしますか？')) {
    emit('triggerRestore', resetDB.value);
  }
};
</script>

<template>
  <div class="bg-slate-900/60 p-4 border border-slate-800 rounded-xl space-y-3">
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-2">
      <div>
        <h4 class="text-xs font-bold text-slate-200 flex items-center gap-1.5">
          <span>🔄</span> ディザスタリカバリ (Layer 2 自動復元)
          <span class="text-[10px] font-mono bg-purple-950/80 text-purple-400 border border-purple-700/50 px-2 py-0.5 rounded">SPEC-RECOVERY-001</span>
        </h4>
        <p class="text-[11px] text-slate-400 mt-0.5">万が一のDB破損時、保存済み JSONダンプから完全オフラインでSQLite3を再構築します。</p>
      </div>
      <div class="flex items-center gap-3">
        <label class="flex items-center gap-1.5 text-xs text-slate-300 cursor-pointer select-none">
          <input type="checkbox" v-model="resetDB" class="rounded bg-slate-950 border-slate-700 text-rose-600" />
          <span>DB初期化</span>
        </label>
        <button @click="handleRestore" :disabled="restoring" class="px-4 py-1.5 bg-purple-600 hover:bg-purple-500 text-white text-xs font-bold rounded-lg disabled:opacity-50 flex items-center gap-1">
          <span>{{ restoring ? '⏳ リストア中...' : '🚀 リストア実行' }}</span>
        </button>
      </div>
    </div>
    <div v-if="restoreStatus" class="p-2 rounded text-xs font-semibold" :class="restoreStatus.success ? 'bg-emerald-950/70 border border-emerald-500/40 text-emerald-300' : 'bg-rose-950/70 border border-rose-500/40 text-rose-300'">
      {{ restoreStatus.success ? '✅' : '⚠️' }} {{ restoreStatus.message }}
    </div>
  </div>
</template>
