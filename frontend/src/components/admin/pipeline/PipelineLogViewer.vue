<!-- frontend/src/components/admin/pipeline/PipelineLogViewer.vue (100行以下 - SPEC-PRINCIPLE-001) -->
<script setup lang="ts">
import { computed } from 'vue';

const props = defineProps<{ logs: any[]; selectedStage: string }>();
const emit = defineEmits<{ (e: 'selectStage', stage: string): void }>();

const tabs = [
  { id: 'all', label: '混合ログ (All)' },
  { id: 'requests', label: '🌐 Requests' },
  { id: 'motrix', label: '🚀 Motrix' },
  { id: 'thunder', label: '⚡ Thunder' },
  { id: 'stash', label: '🎬 Stash' },
];

const filteredLogs = computed(() => {
  if (!props.logs || props.logs.length === 0) return [];
  if (props.selectedStage === 'all') return props.logs;
  const target = props.selectedStage.toUpperCase();
  return props.logs.filter((l) => l.stage === target || (target === 'REQUESTS' && l.stage === 'SYSTEM'));
});

const getLvlBadge = (lvl: string) => {
  if (lvl === 'SUCCESS') return 'text-emerald-400 bg-emerald-950/80 border-emerald-800';
  if (lvl === 'WARN') return 'text-amber-400 bg-amber-950/80 border-amber-800';
  if (lvl === 'ERROR') return 'text-rose-400 bg-rose-950/80 border-rose-800';
  return 'text-blue-400 bg-blue-950/80 border-blue-800';
};
</script>

<template>
  <div class="bg-slate-900/90 border border-slate-800 rounded-2xl p-3 flex flex-col space-y-2 shadow-lg overflow-hidden flex-1 h-full min-h-[220px]">
    <div class="flex items-center justify-between border-b border-slate-800 pb-2">
      <div class="flex items-center gap-1 bg-slate-950 p-1 rounded-xl border border-slate-800/80 overflow-x-auto">
        <button
          v-for="t in tabs"
          :key="t.id"
          @click="emit('selectStage', t.id)"
          :class="['px-2.5 py-0.5 rounded-lg text-xs font-bold cursor-pointer transition whitespace-nowrap', selectedStage === t.id ? 'bg-indigo-600 text-white shadow' : 'text-slate-400 hover:text-slate-200']"
        >
          {{ t.label }}
        </button>
      </div>
      <span class="text-[10px] font-mono text-slate-500">表示 {{ filteredLogs.length }} 件 / 全 {{ logs.length }} 件</span>
    </div>

    <div class="flex-1 overflow-y-auto space-y-1 font-mono text-[11px] pr-1 select-text">
      <div v-if="filteredLogs.length === 0" class="text-center py-12 text-slate-500 text-xs font-mono">
        {{ selectedStage.toUpperCase() }} のログはまだありません
      </div>
      <div v-for="(l, i) in filteredLogs" :key="i" class="p-1.5 bg-slate-950/60 rounded-lg border border-slate-800/40 flex items-start gap-2">
        <span class="text-slate-500 text-[10px] shrink-0 font-mono">{{ l.timestamp }}</span>
        <span :class="['px-1.5 py-0.2 rounded text-[9px] font-bold border shrink-0', getLvlBadge(l.level)]">{{ l.stage }}</span>
        <span class="text-slate-300 break-all flex-1">{{ l.message }}</span>
      </div>
    </div>
  </div>
</template>
