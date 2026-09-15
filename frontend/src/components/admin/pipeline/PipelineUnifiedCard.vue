<!-- frontend/src/components/admin/pipeline/PipelineUnifiedCard.vue (100行以下 - SPEC-PRINCIPLE-001) -->
<script setup lang="ts">
import PipelineFormulaMetrics from './PipelineFormulaMetrics.vue';
import PipelineEngineRibbon from './PipelineEngineRibbon.vue';
import type { PipelineOverview } from '../../../types/pipeline';

defineProps<{
  overview?: PipelineOverview | null;
  loading?: boolean;
  isAutoEngineRunning: boolean;
}>();

const emit = defineEmits<{
  (e: 'toggleAutoEngine'): void;
  (e: 'executeCycleNow'): void;
}>();
</script>

<template>
  <div class="p-3.5 bg-slate-900/95 border border-indigo-900/60 rounded-2xl space-y-3 shadow-xl backdrop-blur select-none">
    <!-- ヘッダー: 自律運転制御 ＆ アクション -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-2 border-b border-slate-800/80 pb-2.5">
      <div class="flex items-center gap-2">
        <span class="text-xs font-bold text-indigo-300 font-mono">自律運転:</span>
        <span v-if="isAutoEngineRunning" class="px-2.5 py-0.5 rounded-md text-[11px] font-mono font-bold bg-emerald-950 text-emerald-300 border border-emerald-700/80 flex items-center gap-1.5 shadow-sm">
          <span class="w-2 h-2 rounded-full bg-emerald-400 shadow-[0_0_8px_rgba(52,211,153,0.9)] animate-pulse"></span>
          <span>完全自動運転中 (常駐ループ)</span>
        </span>
        <span v-else class="px-2.5 py-0.5 rounded-md text-[11px] font-mono font-bold bg-amber-950 text-amber-300 border border-amber-700/80 flex items-center gap-1.5">
          <span class="w-2 h-2 rounded-full bg-amber-400 shadow-[0_0_6px_rgba(251,191,36,0.8)]"></span>
          <span>一時停止中</span>
        </span>
      </div>
      <div class="flex items-center gap-1.5">
        <button
          @click="emit('toggleAutoEngine')"
          :class="['px-3 py-1 font-bold rounded-lg text-xs cursor-pointer active:scale-95 shadow transition', isAutoEngineRunning ? 'bg-amber-600 hover:bg-amber-500 text-white' : 'bg-emerald-600 hover:bg-emerald-500 text-white']"
        >
          {{ isAutoEngineRunning ? '自動運転を一時停止' : '自動運転を開始' }}
        </button>
        <button
          @click="emit('executeCycleNow')"
          :disabled="loading"
          class="px-3 py-1 bg-indigo-600 hover:bg-indigo-500 disabled:opacity-50 text-white font-bold rounded-lg text-xs cursor-pointer active:scale-95 shadow"
        >
          今すぐ実行
        </button>
      </div>
    </div>

    <!-- 進捗バー: メディア救出進捗 (全体) -->
    <div class="space-y-1">
      <div class="flex items-center justify-between text-xs font-mono">
        <span class="text-slate-400">メディア救出進捗 (全 {{ overview?.total_media ?? 0 }} 作品中 {{ overview?.completed ?? 0 }} 作品完了)</span>
        <span class="font-bold text-purple-400">{{ (overview?.overall_progress ?? 0).toFixed(1) }}%</span>
      </div>
      <div class="w-full bg-slate-950 rounded-full h-2 overflow-hidden border border-slate-800">
        <div class="bg-gradient-to-r from-blue-500 via-sky-500 via-purple-500 to-emerald-500 h-full rounded-full transition-all duration-300" :style="{ width: `${overview?.overall_progress ?? 0}%` }"></div>
      </div>
    </div>

    <!-- 中段: ユーザー指定の関係式モデルによるメトリクス -->
    <PipelineFormulaMetrics :metrics="overview?.metrics" />

    <!-- 下段: 各チェックポイント (ダウンローダー) の稼働状態 -->
    <PipelineEngineRibbon :checkpoints="overview?.checkpoints || []" />
  </div>
</template>
