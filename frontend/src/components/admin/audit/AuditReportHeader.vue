<!-- frontend/src/components/admin/audit/AuditReportHeader.vue (100行以下) -->
<script setup lang="ts">
defineProps<{
  report: any | null;
  loading: boolean;
  purgingFiles: boolean;
  statusMessage: { success: boolean; message: string } | null;
}>();
defineEmits<{ (e: 'runAudit', purgeFiles: boolean, purgeDB: boolean): void }>();
</script>

<template>
  <div class="bg-slate-900/60 p-4 border border-slate-800 rounded-xl space-y-3">
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-3">
      <div>
        <h3 class="text-base font-bold text-slate-100 flex items-center gap-2">
          <span>🩺</span> SQLite3 整合性監査 ＆ 孤立ファイルパージ
          <span class="text-[10px] font-mono bg-emerald-950/80 text-emerald-400 border border-emerald-700/50 px-2 py-0.5 rounded">SPEC-AUDIT-001</span>
        </h3>
        <p class="text-xs text-slate-400 mt-0.5">B-Tree・インデックス破損、外部キー整合性、孤立ファイルを完全監査します。</p>
      </div>
      <div class="flex items-center gap-2">
        <button @click="$emit('runAudit', false, false)" :disabled="loading" class="px-4 py-2 bg-blue-600 hover:bg-blue-500 text-white text-xs font-bold rounded-lg disabled:opacity-50">
          {{ loading ? '⏳ 監査中...' : '🔍 整合性監査' }}
        </button>
        <button @click="$emit('runAudit', true, false)" :disabled="loading || purgingFiles" class="px-3.5 py-2 bg-amber-600 hover:bg-amber-500 text-white text-xs font-bold rounded-lg disabled:opacity-50">
          🧹 自動パージ付監査
        </button>
      </div>
    </div>
    <div v-if="statusMessage" class="p-2.5 rounded-lg text-xs font-semibold" :class="statusMessage.success ? 'bg-emerald-950/70 border border-emerald-500/40 text-emerald-300' : 'bg-rose-950/70 border border-rose-500/40 text-rose-300'">
      {{ statusMessage.success ? '✅' : '⚠️' }} {{ statusMessage.message }}
    </div>
    <div v-if="report" class="grid grid-cols-2 md:grid-cols-4 gap-2 pt-1 font-mono text-xs">
      <div class="bg-slate-950/80 p-2.5 rounded-lg border border-slate-800">
        <div class="text-[10px] text-slate-500">B-Tree / インデックス</div>
        <div class="font-bold mt-0.5" :class="report.integrity_ok ? 'text-emerald-400' : 'text-rose-400'">{{ report.integrity_ok ? '🛡️ 健全 (OK)' : '❌ 破損検知' }}</div>
      </div>
      <div class="bg-slate-950/80 p-2.5 rounded-lg border border-slate-800">
        <div class="text-[10px] text-slate-500">外部キー (FK)</div>
        <div class="font-bold mt-0.5" :class="report.foreign_key_ok ? 'text-emerald-400' : 'text-amber-400'">{{ report.foreign_key_ok ? '🛡️ 正常 (0件)' : `⚠️ 違反 ${report.foreign_key_errors?.length || 0}件` }}</div>
      </div>
      <div class="bg-slate-950/80 p-2.5 rounded-lg border border-slate-800">
        <div class="text-[10px] text-slate-500">未紐付け孤立ファイル</div>
        <div class="font-bold mt-0.5" :class="report.orphan_files?.length > 0 ? 'text-amber-400' : 'text-slate-300'">{{ report.orphan_files?.length || 0 }} 件</div>
      </div>
      <div class="bg-slate-950/80 p-2.5 rounded-lg border border-slate-800">
        <div class="text-[10px] text-slate-500">孤立 DB レコード</div>
        <div class="font-bold mt-0.5" :class="report.orphan_db_media?.length > 0 ? 'text-amber-400' : 'text-slate-300'">{{ report.orphan_db_media?.length || 0 }} 件</div>
      </div>
    </div>
  </div>
</template>
