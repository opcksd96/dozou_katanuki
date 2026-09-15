<!-- frontend/src/components/admin/pipeline/PipelineEngineRibbon.vue (100行以下 - SPEC-PRINCIPLE-001) -->
<script setup lang="ts">
import { useThunderOrchestrator } from '../../../composables/admin/useThunderOrchestrator';
import { useDownloaderConsole } from '../../../composables/admin/useDownloaderConsole';
import { useStashResolver } from '../../../composables/useStashResolver';
import type { CheckpointStatus } from '../../../types/pipeline';

defineProps<{ checkpoints: CheckpointStatus[] }>();
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
  <div class="space-y-1.5 pt-1 border-t border-slate-800/80 font-sans select-none">
    <div class="text-[10px] font-bold text-slate-400 font-mono flex items-center gap-1.5">
      <span>通信インフラ・ダウンローダー稼働状況</span>
    </div>
    <div class="grid grid-cols-2 sm:grid-cols-4 gap-2">
      <div
        v-for="cp in checkpoints || []"
        :key="cp.key"
        class="p-2 bg-slate-950/70 border border-slate-800/90 rounded-xl flex flex-col justify-between hover:border-slate-700 transition space-y-1"
      >
        <div class="flex items-center justify-between gap-1">
          <span class="text-[11px] font-bold text-slate-200 truncate" :title="cp.name">{{ cp.name }}</span>
          <button
            @click="kickApp(cp.key)"
            :title="`${cp.name} 起動 / 接続`"
            :class="[
              'px-1.5 py-0.5 rounded text-[9.5px] font-mono font-bold border transition flex items-center gap-1 cursor-pointer active:scale-95',
              cp.status_text?.includes('PAUSED')
                ? 'bg-amber-950 text-amber-300 border-amber-800'
                : cp.is_online
                  ? 'bg-emerald-950 text-emerald-300 border-emerald-800'
                  : 'bg-rose-950 text-rose-300 border-rose-800 animate-pulse'
            ]"
          >
            <span :class="['w-1.5 h-1.5 rounded-full', cp.status_text?.includes('PAUSED') ? 'bg-amber-400' : cp.is_online ? 'bg-emerald-400 shadow-[0_0_6px_#34d399]' : 'bg-rose-500']"></span>
            <span>{{ cp.status_text }}</span>
          </button>
        </div>
        <div class="flex items-center justify-between text-[10px] font-mono text-slate-400">
          <span>稼働: {{ cp.active_count ?? 0 }}</span>
          <span v-if="cp.summary_text" class="text-[9px] text-slate-500 truncate max-w-[110px]" :title="cp.summary_text">{{ cp.summary_text }}</span>
        </div>
      </div>
    </div>
  </div>
</template>
