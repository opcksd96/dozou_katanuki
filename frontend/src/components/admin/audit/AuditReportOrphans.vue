<!-- frontend/src/components/admin/audit/AuditReportOrphans.vue (100行以下) -->
<script setup lang="ts">
const props = defineProps<{
  report: any | null;
  purgingFiles: boolean;
  purgingDB: boolean;
}>();
const emit = defineEmits<{
  (e: 'purgeFiles'): void;
  (e: 'purgeDB'): void;
}>();

const formatBytes = (bytes: number) => {
  if (bytes === 0) return '0 B';
  const k = 1024, sizes = ['B', 'KB', 'MB', 'GB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
};
</script>

<template>
  <div v-if="report" class="grid grid-cols-1 md:grid-cols-2 gap-4">
    <!-- 孤立ファイル一覧 -->
    <div class="bg-slate-900/60 p-4 border border-slate-800 rounded-xl space-y-3">
      <div class="flex items-center justify-between">
        <h4 class="text-xs font-bold text-slate-300">孤立ストレージファイル ({{ report.orphan_files?.length || 0 }}件)</h4>
        <button v-if="report.orphan_files?.length > 0" @click="$emit('purgeFiles')" :disabled="purgingFiles" class="px-2.5 py-1 bg-amber-600 hover:bg-amber-500 text-white text-[11px] font-bold rounded">
          🗑️ ゴミ箱へ退避
        </button>
      </div>
      <div class="max-h-48 overflow-y-auto font-mono text-[11px] divide-y divide-slate-800/60 border border-slate-800/80 rounded-lg bg-slate-950/60">
        <div v-if="!report.orphan_files?.length" class="p-4 text-center text-slate-500 italic">孤立ファイルはありません</div>
        <div v-for="(f, i) in report.orphan_files" :key="i" class="p-2 flex justify-between items-center text-slate-300">
          <span class="truncate max-w-[220px]" :title="f.path">{{ f.file_name }}</span>
          <span class="text-slate-500 text-[10px]">{{ formatBytes(f.file_size) }}</span>
        </div>
      </div>
    </div>

    <!-- 孤立DBレコード一覧 -->
    <div class="bg-slate-900/60 p-4 border border-slate-800 rounded-xl space-y-3">
      <div class="flex items-center justify-between">
        <h4 class="text-xs font-bold text-slate-300">孤立 DB レコード ({{ report.orphan_db_media?.length || 0 }}件)</h4>
        <button v-if="report.orphan_db_media?.length > 0" @click="$emit('purgeDB')" :disabled="purgingDB" class="px-2.5 py-1 bg-rose-600 hover:bg-rose-500 text-white text-[11px] font-bold rounded">
          🗑️ DBレコード削除
        </button>
      </div>
      <div class="max-h-48 overflow-y-auto font-mono text-[11px] divide-y divide-slate-800/60 border border-slate-800/80 rounded-lg bg-slate-950/60">
        <div v-if="!report.orphan_db_media?.length" class="p-4 text-center text-slate-500 italic">孤立レコードはありません</div>
        <div v-for="(m, i) in report.orphan_db_media" :key="i" class="p-2 text-slate-300 space-y-0.5">
          <div class="font-semibold text-rose-300">{{ m.media_id }}</div>
          <div class="text-[10px] text-slate-500">{{ m.reason }}</div>
        </div>
      </div>
    </div>
  </div>
</template>
