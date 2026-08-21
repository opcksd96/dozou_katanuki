<!-- frontend/src/components/admin/audit/AuditReportOrphans.vue (100行以下) -->
<script setup lang="ts">
import { computed } from 'vue';
import { useAdminAudit } from '../../../composables/admin/useAdminAudit';

const adminAudit = useAdminAudit();

const props = withDefaults(defineProps<{
  report: any | null;
  purgingFiles?: boolean;
  purgingDb?: boolean;
  purgingDB?: boolean;
}>(), {
  report: null, purgingFiles: false, purgingDb: false, purgingDB: false,
});
const emit = defineEmits<{
  (e: 'purgeFiles', paths?: string[]): void;
  (e: 'purge-files', paths?: string[]): void;
  (e: 'purgeDB', mediaIDs?: string[]): void;
  (e: 'purge-db', mediaIDs?: string[]): void;
}>();

const isPurgingF = computed(() => props.purgingFiles || adminAudit.isPurgingFiles.value);
const isPurgingD = computed(() => props.purgingDb || props.purgingDB || adminAudit.isPurgingDB.value);

const totalOrphanBytes = computed(() => {
  return (props.report?.orphan_files || []).reduce((acc: number, f: any) => acc + (f.file_size || f.FileSize || 0), 0);
});

const formatBytes = (bytes: number) => {
  if (!bytes || bytes <= 0) return '0 B';
  const k = 1024, sizes = ['B', 'KB', 'MB', 'GB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
};

const handlePurgeAllFiles = () => {
  const paths = (props.report?.orphan_files || []).map((f: any) => f.path || f.Path).filter(Boolean);
  emit('purgeFiles', paths); emit('purge-files', paths);
  adminAudit.purgeOrphanFiles(paths);
};

const handlePurgeSingleFile = (f: any) => {
  const p = f?.path || f?.Path;
  if (!p) return;
  emit('purgeFiles', [p]); emit('purge-files', [p]);
  adminAudit.purgeOrphanFiles([p]);
};

const handlePurgeAllDB = () => {
  const ids = (props.report?.orphan_db_media || []).map((m: any) => m.media_id || m.MediaID || m.mediaId).filter(Boolean);
  emit('purgeDB', ids); emit('purge-db', ids);
  adminAudit.purgeOrphanDBMedia(ids);
};

const handlePurgeSingleDB = (m: any) => {
  const id = m?.media_id || m?.MediaID || m?.mediaId;
  if (!id) return;
  emit('purgeDB', [id]); emit('purge-db', [id]);
  adminAudit.purgeOrphanDBMedia([id]);
};
</script>

<template>
  <div v-if="report" class="grid grid-cols-1 md:grid-cols-2 gap-4">
    <!-- 孤立ファイル一覧 -->
    <div class="bg-slate-900/60 p-4 border border-slate-800 rounded-xl space-y-3 flex flex-col">
      <div class="flex items-center justify-between">
        <div>
          <h4 class="text-xs font-bold text-slate-300">孤立ストレージファイル ({{ report.orphan_files?.length || 0 }}件)</h4>
          <p class="text-[10px] text-slate-500 font-mono">合計: {{ formatBytes(totalOrphanBytes) }}</p>
        </div>
        <button v-if="report.orphan_files?.length > 0" @click="handlePurgeAllFiles" :disabled="isPurgingF" class="px-2.5 py-1 bg-amber-600 hover:bg-amber-500 text-white text-[11px] font-bold rounded shadow disabled:opacity-50 flex items-center gap-1 cursor-pointer">
          <span>{{ isPurgingF ? '⏳' : '🗑️' }}</span>
          <span>{{ isPurgingF ? '退避中...' : '全件ゴミ箱へ退避' }}</span>
        </button>
      </div>
      <div class="max-h-56 overflow-y-auto font-mono text-[11px] divide-y divide-slate-800/60 border border-slate-800/80 rounded-lg bg-slate-950/60 flex-1">
        <div v-if="!report.orphan_files?.length" class="p-6 text-center text-slate-500 italic text-xs">✨ 孤立ファイルはありません</div>
        <div v-for="(f, i) in report.orphan_files" :key="i" class="p-2 flex justify-between items-center text-slate-300 hover:bg-slate-900/50">
          <div class="truncate max-w-[200px] flex items-center gap-1.5" :title="f.path || f.Path">
            <span class="text-[9px] px-1 py-0.5 rounded bg-slate-800 text-slate-400 font-sans uppercase">{{ f.category || f.Category || 'file' }}</span>
            <span class="truncate">{{ f.file_name || f.FileName }}</span>
          </div>
          <div class="flex items-center gap-2 shrink-0">
            <span class="text-slate-500 text-[10px]">{{ formatBytes(f.file_size || f.FileSize) }}</span>
            <button @click="handlePurgeSingleFile(f)" :disabled="isPurgingF" class="text-rose-400 hover:text-rose-300 text-[10px] px-2 py-1 rounded bg-rose-950/50 hover:bg-rose-900/60 border border-rose-800/50 cursor-pointer disabled:opacity-50 font-sans font-semibold">退避</button>
          </div>
        </div>
      </div>
    </div>

    <!-- 孤立DBレコード一覧 -->
    <div class="bg-slate-900/60 p-4 border border-slate-800 rounded-xl space-y-3 flex flex-col">
      <div class="flex items-center justify-between">
        <div>
          <h4 class="text-xs font-bold text-slate-300">孤立 DB レコード ({{ report.orphan_db_media?.length || 0 }}件)</h4>
          <p class="text-[10px] text-slate-500">実ファイル未存在 または 親記事未紐付け</p>
        </div>
        <div class="flex items-center gap-1.5">
          <button v-if="adminAudit.canRollback.value" @click="adminAudit.rollbackLastPurge" :disabled="adminAudit.isRollingBack.value || isPurgingD" class="px-2 py-1 bg-purple-700/80 hover:bg-purple-600 text-white text-[10px] font-bold rounded shadow disabled:opacity-50 flex items-center gap-1 cursor-pointer" title="直前に削除したレコードを退避JSONから復元">
            <span>{{ adminAudit.isRollingBack.value ? '⏳' : '↩️' }}</span>
            <span>{{ adminAudit.isRollingBack.value ? '復元中...' : 'Undo' }}</span>
          </button>
          <button v-if="report.orphan_db_media?.length > 0" @click="handlePurgeAllDB" :disabled="isPurgingD" class="px-2.5 py-1 bg-rose-600 hover:bg-rose-500 text-white text-[11px] font-bold rounded shadow disabled:opacity-50 flex items-center gap-1 cursor-pointer">
            <span>{{ isPurgingD ? '⏳' : '🗑️' }}</span>
            <span>{{ isPurgingD ? '削除中...' : '全件レコード削除' }}</span>
          </button>
        </div>
      </div>
      <div class="max-h-56 overflow-y-auto font-mono text-[11px] divide-y divide-slate-800/60 border border-slate-800/80 rounded-lg bg-slate-950/60 flex-1">
        <div v-if="!report.orphan_db_media?.length" class="p-6 text-center text-slate-500 italic text-xs">✨ 孤立レコードはありません</div>
        <div v-for="(m, i) in report.orphan_db_media" :key="i" class="p-2 text-slate-300 flex items-center justify-between hover:bg-slate-900/50">
          <div class="space-y-0.5 truncate max-w-[220px]">
            <div class="font-semibold text-rose-300 truncate" :title="m.media_id || m.MediaID || m.mediaId">{{ m.media_id || m.MediaID || m.mediaId }}</div>
            <div class="text-[10px] text-slate-500 truncate">{{ m.reason || m.Reason }} <span v-if="m.article_id || m.ArticleID">({{ m.article_id || m.ArticleID }})</span></div>
          </div>
          <button @click="handlePurgeSingleDB(m)" :disabled="isPurgingD" class="text-rose-400 hover:text-rose-300 text-[10px] px-2 py-1 rounded bg-rose-950/50 hover:bg-rose-900/60 border border-rose-800/50 shrink-0 cursor-pointer disabled:opacity-50 font-sans font-semibold">削除</button>
        </div>
      </div>
    </div>
  </div>
</template>
