<!-- frontend/src/components/admin/pipeline/PipelineCheckpointRibbon.vue (100行以下 - SPEC-PRINCIPLE-001) -->
<script setup lang="ts">
import { useThunderOrchestrator } from '../../../composables/admin/useThunderOrchestrator';
import { useDownloaderConsole } from '../../../composables/admin/useDownloaderConsole';
import { useStashResolver } from '../../../composables/useStashResolver';

defineProps<{ checkpoints: any[] }>();
const { launchThunder } = useThunderOrchestrator();
const { launchMotrix } = useDownloaderConsole();
const { openStashWebUI } = useStashResolver();

const kickApp = (key: string) => {
  if (key === 'thunder') launchThunder();
  else if (key === 'stash') openStashWebUI();
  else if (key === 'motrix') launchMotrix();
};
</script>

<template>
  <div class="flex items-stretch gap-2.5 overflow-x-auto pb-1 snap-x select-none">
    <div
      v-for="(cp, idx) in checkpoints || []"
      :key="cp.key"
      class="p-3 bg-slate-900/90 border border-slate-800 rounded-2xl flex flex-col justify-between shadow-md relative overflow-hidden group hover:border-slate-700 transition flex-1 min-w-[210px] shrink-0 snap-start"
    >
      <div class="flex items-center justify-between gap-2">
        <div class="flex items-center gap-1.5">
          <span class="text-xs font-bold text-slate-200 truncate">{{ cp.name }}</span>
        </div>

        <button
          @click="kickApp(cp.key)"
          :title="`${cp.name} を起動/接続する`"
          :class="[
            'px-2 py-0.5 rounded-lg text-[10px] font-mono font-bold border transition-all cursor-pointer flex items-center gap-1 active:scale-90 shadow-sm',
            cp.status_text?.includes('PAUSED')
              ? 'bg-amber-950/90 text-amber-300 border-amber-700/80 hover:bg-amber-900 shadow-amber-950'
              : cp.is_online
                ? 'bg-emerald-950/90 text-emerald-300 border-emerald-700/80 hover:bg-emerald-900 shadow-emerald-950'
                : 'bg-rose-950/90 text-rose-300 border-rose-700/80 hover:bg-rose-900 animate-pulse shadow-rose-950'
          ]"
        >
          <span :class="['w-1.5 h-1.5 rounded-full', cp.status_text?.includes('PAUSED') ? 'bg-amber-400 shadow-[0_0_6px_rgba(251,191,36,0.8)]' : cp.is_online ? 'bg-emerald-400 shadow-[0_0_6px_rgba(52,211,153,0.8)]' : 'bg-rose-500']"></span>
          <span>{{ cp.status_text }}</span>
        </button>
      </div>

      <div class="pt-2 mt-1 border-t border-slate-800/60 space-y-1 font-mono text-[11px]">
        <div class="flex items-baseline justify-between">
          <span class="text-slate-400 text-[10px]">通信タスク稼働数:</span>
          <span class="font-bold text-slate-100 text-xs">{{ cp.active_count ?? 0 }} タスク</span>
        </div>
        <div v-if="cp.summary_text" class="text-[9.5px] text-slate-400 bg-slate-950/80 px-2 py-0.5 rounded-md border border-slate-800/80 truncate leading-tight" :title="cp.summary_text">
          📊 {{ cp.summary_text }}
        </div>
      </div>
    </div>
  </div>
</template>
