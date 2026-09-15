<!-- frontend/src/components/admin/pipeline/PipelineLogViewer.vue (100行以下 - SPEC-PRINCIPLE-001) -->
<script setup lang="ts">
import { computed } from 'vue';

const props = defineProps<{ logs: any[]; selectedStage: string }>();
const emit = defineEmits<{ (e: 'selectStage', stage: string): void }>();

const tabs = [
  { id: 'all', label: '混合ログ (All)' },
  { id: 'scraper', label: '🕷️ スクレイパー' },
  { id: 'requests', label: '🌐 Requests' },
  { id: 'motrix', label: '🚀 Motrix' },
  { id: 'thunder', label: '⚡ Thunder' },
  { id: 'stash', label: '🎬 Stash' },
];

interface NormalizedLog {
  timestamp: string;
  stage: string;
  level: string;
  message: string;
}

const normalize = (l: any): NormalizedLog => {
  if (!l) return { timestamp: '', stage: 'INFO', level: 'INFO', message: '(空エントリ)' };
  if (typeof l === 'string') return { timestamp: '', stage: 'LOG', level: 'INFO', message: l };
  const ts = l.timestamp || l.Timestamp || l.time || '';
  const stage = (l.stage || l.Stage || 'SYSTEM').toUpperCase();
  const level = (l.level || l.Level || 'INFO').toUpperCase();
  const rawMsg = l.message ?? l.Message ?? l.msg ?? l.text ?? '';
  const message = typeof rawMsg === 'string' ? rawMsg : JSON.stringify(rawMsg);
  return { timestamp: ts, stage, level, message: message || '(メッセージなし)' };
};

const filteredLogs = computed(() => {
  if (!props.logs || props.logs.length === 0) return [];
  const normalized = props.logs.map(normalize);
  if (props.selectedStage === 'all') return normalized;
  const target = props.selectedStage.toUpperCase();
  return normalized.filter((l) => l.stage === target || (target === 'REQUESTS' && l.stage === 'SYSTEM'));
});

const getLvlBadge = (lvl: string) => {
  if (lvl === 'SUCCESS') return 'text-emerald-400 bg-emerald-950/80 border-emerald-800';
  if (lvl === 'WARN') return 'text-amber-400 bg-amber-950/80 border-amber-800';
  if (lvl === 'ERROR') return 'text-rose-400 bg-rose-950/80 border-rose-800';
  return 'text-blue-400 bg-blue-950/80 border-blue-800';
};
</script>

<template>
  <div class="bg-slate-900/90 border border-slate-800 rounded-2xl p-3.5 flex flex-col space-y-2.5 shadow-lg w-full">
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

    <div class="space-y-1 font-mono text-[11px] pt-1 select-text w-full">
      <div v-if="filteredLogs.length === 0" class="text-center py-12 text-slate-500 text-xs font-mono">
        {{ selectedStage.toUpperCase() }} のログはまだありません
      </div>
      <div v-for="(l, i) in filteredLogs" :key="i" class="p-1.5 bg-slate-950/60 rounded-lg border border-slate-800/40 flex items-start gap-2">
        <span v-if="l.timestamp" class="text-slate-500 text-[10px] shrink-0 font-mono">{{ l.timestamp }}</span>
        <span :class="['px-1.5 py-0.5 rounded text-[9px] font-bold border shrink-0', getLvlBadge(l.level)]">{{ l.stage }}</span>
        <span class="text-slate-300 break-all flex-1 leading-relaxed">{{ l.message }}</span>
      </div>
    </div>
  </div>
</template>
