<!-- frontend/src/components/admin/ThunderProgressPanel.vue (100行以下 - SPEC-PRINCIPLE-001) -->
<script setup lang="ts">
const props = defineProps<{ status: any }>();
const progressPercent = () => {
  if (!props.status?.total_jobs || props.status.total_jobs === 0) return 0;
  const done = (props.status.success_jobs || 0) + (props.status.failed_jobs || 0);
  return (done / props.status.total_jobs) * 100;
};
</script>

<template>
  <div class="space-y-2 pt-1">
    <!-- 4グリッド統計カード (Motrix と完全同一のレイアウト・トーン＆マナー) -->
    <div class="grid grid-cols-4 gap-2 text-center">
      <div class="p-1.5 bg-slate-950 rounded-xl border border-slate-800">
        <div class="text-[10px] text-slate-400 font-medium">稼働 (Running)</div>
        <div class="text-xs font-mono font-bold text-purple-400">{{ status?.running_jobs ?? 0 }} 件</div>
      </div>
      <div class="p-1.5 bg-slate-950 rounded-xl border border-slate-800">
        <div class="text-[10px] text-slate-400 font-medium">待機 (Pending)</div>
        <div class="text-xs font-mono font-bold text-amber-400">{{ status?.pending_jobs ?? 0 }} 件</div>
      </div>
      <div class="p-1.5 bg-slate-950 rounded-xl border border-slate-800">
        <div class="text-[10px] text-slate-400 font-medium">救出完了 (Done)</div>
        <div class="text-xs font-mono font-bold text-emerald-400">{{ status?.success_jobs ?? 0 }} 件</div>
      </div>
      <div class="p-1.5 bg-slate-950 rounded-xl border border-slate-800">
        <div class="text-[10px] text-slate-400 font-medium">枯渇退避 (Depleted)</div>
        <div class="text-xs font-mono font-bold text-rose-400">{{ status?.failed_jobs ?? 0 }} 件</div>
      </div>
    </div>

    <!-- 全体プログレスバー -->
    <div class="bg-slate-950 p-2 rounded-xl border border-slate-800 space-y-1">
      <div class="flex items-center justify-between text-[11px] font-mono">
        <span class="text-slate-400">救出エスカレーション進捗 (全 {{ status?.total_jobs ?? 303 }} ジョブ)</span>
        <span class="font-bold text-purple-400">{{ progressPercent().toFixed(1) }}%</span>
      </div>
      <div class="w-full bg-slate-900 rounded-full h-1.5 overflow-hidden">
        <div class="bg-gradient-to-r from-purple-500 via-indigo-500 to-emerald-500 h-full rounded-full transition-all duration-300" :style="{ width: `${progressPercent()}%` }"></div>
      </div>
    </div>
  </div>
</template>
