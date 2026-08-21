<!-- frontend/src/components/admin/jobs/JobProgressBar.vue (100行以下) -->
<script setup lang="ts">
defineProps<{ activeJob: any }>();
defineEmits<{ (e: 'cancelJob', jobId: string): void }>();
</script>

<template>
  <div v-if="activeJob" class="bg-slate-900/90 border border-blue-500/40 rounded-xl p-4 space-y-3">
    <div class="flex items-center justify-between text-xs">
      <div class="flex items-center gap-2">
        <span class="inline-flex items-center px-2 py-0.5 rounded text-[10px] font-semibold font-mono"
          :class="{
            'bg-blue-500/20 text-blue-400 border border-blue-500/30 animate-pulse': activeJob.status === 'running',
            'bg-amber-500/20 text-amber-400 border border-amber-500/30': activeJob.status === 'pending',
            'bg-emerald-500/20 text-emerald-400 border border-emerald-500/30': activeJob.status === 'completed',
            'bg-rose-500/20 text-rose-400 border border-rose-500/30': activeJob.status === 'failed',
            'bg-slate-700 text-slate-300': activeJob.status === 'cancelled',
          }"
        >
          ● {{ activeJob.status?.toUpperCase() }}
        </span>
        <span class="font-mono text-slate-300 font-bold">[{{ activeJob.type }}] {{ activeJob.id }}</span>
      </div>
      <div class="flex items-center gap-3">
        <span class="text-slate-400 font-mono text-[11px]">
          {{ activeJob.current }} / {{ activeJob.total > 0 ? activeJob.total : '?' }} 件 ({{ (activeJob.percentage || 0).toFixed(1) }}%)
        </span>
        <button
          v-if="activeJob.status === 'running' || activeJob.status === 'pending'"
          @click="$emit('cancelJob', activeJob.id)"
          class="px-2.5 py-1 bg-rose-900/40 hover:bg-rose-800 text-rose-300 border border-rose-700/50 rounded text-xs"
        >
          ■ 中断
        </button>
      </div>
    </div>
    <div class="w-full bg-slate-950 h-3 rounded-full overflow-hidden p-0.5 border border-slate-800">
      <div class="h-full bg-gradient-to-r from-blue-600 via-indigo-500 to-emerald-400 rounded-full transition-all duration-300" :style="{ width: `${Math.min(100, Math.max(0, activeJob.percentage || 0))}%` }"></div>
    </div>
    <div class="text-xs text-slate-400 truncate font-mono">{{ activeJob.message || '処理中...' }}</div>
  </div>
</template>
