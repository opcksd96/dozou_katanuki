<!-- frontend/src/components/admin/ThunderProgressPanel.vue (100行以下 - SPEC-PRINCIPLE-001) -->
<script setup lang="ts">
const props = defineProps<{
  status: any;
}>();

const progressPercent = () => {
  if (!props.status?.total_jobs || props.status.total_jobs === 0) return 0;
  const done = (props.status.success_jobs || 0) + (props.status.failed_jobs || 0);
  return (done / props.status.total_jobs) * 100;
};
</script>

<template>
  <div class="space-y-1 bg-slate-950 p-2.5 rounded-xl border border-slate-800">
    <div class="flex items-center justify-between text-xs font-mono">
      <span class="text-slate-400">全体進捗 (厳選3解像度 × {{ status?.total_media_count ?? 101 }} 件 = {{ status?.total_jobs ?? 303 }} ジョブ)</span>
      <span class="font-bold text-purple-400">{{ progressPercent().toFixed(1) }}%</span>
    </div>
    <div class="w-full bg-slate-900 rounded-full h-2 overflow-hidden">
      <div class="bg-gradient-to-r from-purple-500 via-indigo-500 to-blue-500 h-full rounded-full transition-all duration-300" :style="{ width: `${progressPercent()}%` }"></div>
    </div>
    <div class="grid grid-cols-4 gap-1.5 text-center text-[10px] font-mono pt-1">
      <div class="p-1 bg-slate-900 rounded border border-slate-800/80"><span class="text-slate-400">待機:</span> <span class="text-amber-400 font-bold">{{ status?.pending_jobs ?? 0 }}</span></div>
      <div class="p-1 bg-slate-900 rounded border border-slate-800/80"><span class="text-slate-400">実行中:</span> <span class="text-purple-400 font-bold">{{ status?.running_jobs ?? 0 }}</span></div>
      <div class="p-1 bg-slate-900 rounded border border-slate-800/80"><span class="text-slate-400">投入済:</span> <span class="text-emerald-400 font-bold">{{ status?.success_jobs ?? 0 }}</span></div>
      <div class="p-1 bg-slate-900 rounded border border-slate-800/80"><span class="text-slate-400">失敗:</span> <span class="text-rose-400 font-bold">{{ status?.failed_jobs ?? 0 }}</span></div>
    </div>
  </div>
</template>
