<!-- frontend/src/components/admin/AuditReportView.vue (100行以下 - SPEC-PRINCIPLE-001) -->
<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useAdminAudit } from '../../composables/admin/useAdminAudit';
import { Server, Database, RefreshCw, ExternalLink } from 'lucide-vue-next';
import { BrowserOpenURL } from '../../../wailsjs/runtime/runtime';
import { useStashResolver } from '../../composables/useStashResolver';
import AuditReportHeader from './audit/AuditReportHeader.vue';
import AuditReportOrphans from './audit/AuditReportOrphans.vue';
import AuditReportRestore from './audit/AuditReportRestore.vue';

const props = defineProps<{ restoring?: boolean }>();
const emit = defineEmits<{ (e: 'triggerRestore', resetDB?: boolean): void }>();
const { auditReport, isAuditing, isPurgingFiles, isPurgingDB, auditStatusMessage, runAudit, purgeOrphanFiles, purgeOrphanDBMedia } = useAdminAudit();

const isStashOnline = ref(false), stashLatency = ref<number | null>(null);
const checkStashHealth = async () => {
  const start = Date.now();
  try {
    const r = await fetch('/stash-proxy/', { method: 'HEAD' });
    stashLatency.value = Date.now() - start;
    isStashOnline.value = r.ok || r.status === 401 || r.status === 404;
  } catch { isStashOnline.value = false; stashLatency.value = null; }
};

const { openStashWebUI } = useStashResolver();
const openStash = () => openStashWebUI();

onMounted(() => {
  checkStashHealth();
  if (!auditReport.value && !isAuditing.value) runAudit(false, false);
});
</script>

<template>
  <div class="space-y-4 font-sans max-w-4xl mx-auto py-1">
    <!-- 最上段: Stash & データベース統合ヘルスサマリーカード -->
    <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
      <div class="p-3.5 bg-slate-900/90 border border-slate-800 rounded-2xl flex items-center justify-between shadow-md">
        <div class="flex items-center gap-3">
          <div class="w-8 h-8 rounded-xl bg-blue-600/20 border border-blue-500/30 flex items-center justify-center text-blue-400"><Server class="w-4 h-4" /></div>
          <div>
            <div class="text-xs font-bold text-slate-100 flex items-center gap-1.5">
              <span>Stash Media Server</span>
              <button @click="openStash" class="text-[10px] text-blue-400 hover:text-blue-300 font-normal flex items-center gap-0.5 cursor-pointer underline" title="Stash Web UI を開く">
                <span>Stashを開く</span><ExternalLink class="w-2.5 h-2.5" />
              </button>
            </div>
            <div class="text-[10px] font-mono text-slate-400">{{ isStashOnline ? `応答時間: ${stashLatency}ms (GraphQL Proxy)` : 'オフライン (接続待機中)' }}</div>
          </div>
        </div>
        <div class="flex items-center gap-2">
          <span :class="isStashOnline ? 'bg-emerald-950 text-emerald-300 border-emerald-700/60' : 'bg-rose-950 text-rose-300 border-rose-700/60'" class="px-2 py-0.5 rounded-full text-[10px] font-mono font-bold border">
            {{ isStashOnline ? '🟢 ONLINE' : '🔴 OFFLINE' }}
          </span>
          <button @click="checkStashHealth" class="p-1 text-slate-400 hover:text-slate-200 cursor-pointer active:scale-95" title="ヘルスチェック"><RefreshCw class="w-3.5 h-3.5" /></button>
        </div>
      </div>

      <div class="p-3.5 bg-slate-900/90 border border-slate-800 rounded-2xl flex items-center justify-between shadow-md">
        <div class="flex items-center gap-3">
          <div class="w-8 h-8 rounded-xl bg-indigo-600/20 border border-indigo-500/30 flex items-center justify-center text-indigo-400"><Database class="w-4 h-4" /></div>
          <div>
            <div class="text-xs font-bold text-slate-100">SQLite SSOT (archive.db)</div>
            <div class="text-[10px] font-mono text-slate-400">WAL Mode / 整合性チェック Ready</div>
          </div>
        </div>
        <span class="px-2.5 py-0.5 rounded-full text-[10px] font-mono font-bold border bg-indigo-950 text-indigo-300 border-indigo-700/60">
          HEALTHY 🛡️
        </span>
      </div>
    </div>

    <AuditReportHeader :report="auditReport" :loading="isAuditing" :purging-files="isPurgingFiles" :status-message="auditStatusMessage" @run-audit="(pf, pdb) => runAudit(pf, pdb)" />
    <AuditReportOrphans :report="auditReport" :purging-files="isPurgingFiles" :purging-db="isPurgingDB" @purge-files="(paths) => purgeOrphanFiles(paths)" @purge-db="(ids) => purgeOrphanDBMedia(ids)" />
    <AuditReportRestore :restoring="restoring" :restore-status="null" @trigger-restore="(resetDB) => emit('triggerRestore', resetDB)" />
  </div>
</template>
