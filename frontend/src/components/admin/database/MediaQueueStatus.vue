<!-- frontend/src/components/admin/database/MediaQueueStatus.vue (100行以下) -->
<script setup lang="ts">
import { computed } from 'vue';

interface DownloadQueueStats {
  queued?: number; completed?: number; dead_404?: number; outsourced?: number; retained?: number; failed?: number; total?: number;
}

const props = defineProps<{ stats?: DownloadQueueStats | any; activeJob?: any; statusFilter?: string; }>();
const emit = defineEmits<{ (e: 'update:statusFilter', v: string): void; }>();

const segments = computed(() => [
  { label: 'QUEUED', value: props.stats?.queued ?? props.stats?.Queued ?? 0, color: 'bg-slate-500', status: 'QUEUED' },
  { label: 'COMPLETED', value: props.stats?.completed ?? props.stats?.Completed ?? 0, color: 'bg-emerald-500', status: 'COMPLETED' },
  { label: 'OUTSOURCED', value: props.stats?.outsourced ?? props.stats?.Outsourced ?? 0, color: 'bg-purple-500', status: 'OUTSOURCED' },
  { label: 'RETAINED', value: props.stats?.retained ?? props.stats?.Retained ?? 0, color: 'bg-amber-500', status: 'RETAINED' },
  { label: 'DEAD_404', value: props.stats?.dead_404 ?? props.stats?.Dead404 ?? props.stats?.dead404 ?? 0, color: 'bg-rose-500', status: 'DEAD_404' },
  { label: 'FAILED', value: props.stats?.failed ?? props.stats?.Failed ?? 0, color: 'bg-red-600', status: 'FAILED' },
]);

const total = computed(() => (props.stats?.total ?? props.stats?.Total ?? segments.value.reduce((s, x) => s + x.value, 0)) || 0);
const pct = (v: number) => (total.value > 0 ? (v / total.value) * 100 : 0);
const isJobActive = computed(() => props.activeJob && props.activeJob.status === 'RUNNING');
</script>

<template>
  <div class="p-2.5 bg-slate-900 border border-slate-800 rounded-xl flex items-center justify-between gap-3">
    <div class="flex-1 min-w-0">
      <div class="flex h-3 w-full rounded-full overflow-hidden bg-slate-800">
        <div v-for="s in segments" :key="s.status" class="h-full transition-all duration-300" :class="s.color" :style="{ width: `${pct(s.value)}%` }" :title="`${s.label}: ${s.value}`" />
      </div>
      <div class="mt-1.5 flex flex-wrap gap-2">
        <button v-for="s in segments" :key="s.status" @click="emit('update:statusFilter', s.status)"
          :class="statusFilter === s.status ? 'ring-2 ring-slate-300' : 'opacity-80 hover:opacity-100'"
          class="flex items-center gap-1 px-1.5 py-0.5 rounded text-[11px] font-mono bg-slate-800 text-slate-200 cursor-pointer">
          <span class="w-2 h-2 rounded-full" :class="s.color" />
          <span class="font-bold">{{ s.value }}</span>
          <span class="text-slate-400">{{ s.label }}</span>
        </button>
      </div>
    </div>
    <div v-if="isJobActive" class="shrink-0 flex items-center gap-2 px-2 py-1 rounded-lg bg-blue-900/60 border border-blue-600/60 text-xs font-mono animate-pulse">
      <span class="text-blue-400 font-bold">● {{ activeJob.type }}</span>
      <span class="text-blue-300">{{ activeJob.percentage }}%</span>
    </div>
  </div>
</template>
