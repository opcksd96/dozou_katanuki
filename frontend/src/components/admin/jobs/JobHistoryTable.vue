<!-- frontend/src/components/admin/jobs/JobHistoryTable.vue (100行以下) -->
<script setup lang="ts">
defineProps<{ jobList: any[]; loadingJobs: boolean }>();
defineEmits<{ (e: 'fetchJobs'): void; (e: 'cancelJob', id: string): void }>();
const formatDate = (dateStr?: any) => {
  if (!dateStr) return '-';
  try {
    const d = new Date(dateStr);
    return isNaN(d.getTime()) ? String(dateStr) : d.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit', second: '2-digit' });
  } catch { return String(dateStr); }
};
</script>

<template>
  <div class="bg-slate-900/60 border border-slate-800 rounded-xl p-4 space-y-3">
    <div class="flex items-center justify-between">
      <h3 class="text-sm font-bold text-slate-200 flex items-center gap-2"><span>📜</span> ジョブ実行履歴</h3>
      <button @click="$emit('fetchJobs')" :disabled="loadingJobs" class="text-xs text-slate-400 hover:text-blue-400 flex items-center gap-1">
        <span :class="{ 'animate-spin': loadingJobs }">🔄</span> 更新
      </button>
    </div>
    <div class="overflow-x-auto max-h-48 overflow-y-auto">
      <table class="w-full text-left text-xs font-mono text-slate-300">
        <thead class="bg-slate-950/80 text-slate-400 sticky top-0 border-b border-slate-800 text-[11px]">
          <tr>
            <th class="py-2 px-3">Status</th><th class="py-2 px-3">Job ID</th><th class="py-2 px-3">Type</th>
            <th class="py-2 px-3">Progress</th><th class="py-2 px-3">Started</th><th class="py-2 px-3">Action</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-slate-800/60">
          <tr v-if="jobList.length === 0"><td colspan="6" class="py-4 text-center text-slate-500 italic">実行履歴はありません</td></tr>
          <tr v-for="job in jobList" :key="job.id" class="hover:bg-slate-800/40">
            <td class="py-2 px-3">
              <span class="inline-block px-1.5 py-0.5 rounded text-[10px] font-semibold"
                :class="{
                  'bg-blue-500/20 text-blue-400': job.status === 'running', 'bg-emerald-500/20 text-emerald-400': job.status === 'completed',
                  'bg-rose-500/20 text-rose-400': job.status === 'failed', 'bg-slate-700 text-slate-300': job.status === 'cancelled',
                }">{{ job.status }}</span>
            </td>
            <td class="py-2 px-3 text-slate-200 font-semibold truncate max-w-[140px]">{{ job.id }}</td>
            <td class="py-2 px-3 text-slate-400">{{ job.type }}</td>
            <td class="py-2 px-3 text-slate-300">{{ job.current }} / {{ job.total }}</td>
            <td class="py-2 px-3 text-slate-400 text-[11px]">{{ formatDate(job.started_at) }}</td>
            <td class="py-2 px-3">
              <button v-if="job.status === 'running' || job.status === 'pending'" @click="$emit('cancelJob', job.id)" class="text-rose-400 hover:text-rose-300 underline text-[11px]">中止</button>
              <span v-else class="text-slate-600">-</span>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>
