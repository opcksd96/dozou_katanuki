<!-- frontend/src/components/admin/pipeline/PipelineCheckpointRibbon.vue (100行以下 - SPEC-PRINCIPLE-001) -->
<script setup lang="ts">
import { BrowserOpenURL } from '../../../../wailsjs/runtime/runtime';
import { useThunderOrchestrator } from '../../../composables/admin/useThunderOrchestrator';

defineProps<{ checkpoints: any[] }>();
const { launchThunder } = useThunderOrchestrator();

const kickApp = (key: string) => {
  if (key === 'thunder') launchThunder();
  else if (key === 'stash') {
    try { BrowserOpenURL('http://127.0.0.1:9999/'); } catch { window.open('http://127.0.0.1:9999/', '_blank'); }
  } else if (key === 'motrix') {
    try { BrowserOpenURL('motrix://'); } catch {}
  }
};
</script>

<template>
  <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-2.5">
    <div
      v-for="(cp, idx) in checkpoints || []"
      :key="cp.key"
      class="p-3 bg-slate-900/90 border border-slate-800 rounded-2xl flex flex-col justify-between shadow-md relative overflow-hidden group hover:border-slate-700 transition"
    >
      <div class="flex items-center justify-between gap-2">
        <div class="flex items-center gap-2">
          <span class="text-sm">
            {{ idx === 0 ? '🌐' : idx === 1 ? '🚀' : idx === 2 ? '⚡' : '🎬' }}
          </span>
          <span class="text-xs font-bold text-slate-200 truncate">{{ cp.name }}</span>
        </div>

        <!-- パイロットランプ 兼 キックボタン (クリックでアプリ起動/接続) -->
        <button
          @click="kickApp(cp.key)"
          :title="`${cp.name} を起動/接続する`"
          :class="[
            'px-2 py-0.5 rounded-lg text-[10px] font-mono font-bold border transition-all cursor-pointer flex items-center gap-1 active:scale-90 shadow-sm',
            cp.is_online
              ? 'bg-emerald-950/90 text-emerald-300 border-emerald-700/80 hover:bg-emerald-900 shadow-emerald-950'
              : 'bg-rose-950/90 text-rose-300 border-rose-700/80 hover:bg-rose-900 animate-pulse shadow-rose-950'
          ]"
        >
          <span :class="['w-1.5 h-1.5 rounded-full', cp.is_online ? 'bg-emerald-400 shadow-[0_0_6px_rgba(52,211,153,0.8)]' : 'bg-rose-500']"></span>
          <span>{{ cp.status_text }}</span>
        </button>
      </div>

      <div class="flex items-baseline justify-between pt-2 mt-1 border-t border-slate-800/60 text-[11px] font-mono">
        <span class="text-slate-400 text-[10px]">処理中/稼働数:</span>
        <span class="font-bold text-slate-100 text-xs">{{ cp.active_count ?? 0 }} 件</span>
      </div>
    </div>
  </div>
</template>
