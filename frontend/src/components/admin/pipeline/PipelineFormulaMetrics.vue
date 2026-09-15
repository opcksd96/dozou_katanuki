<!-- frontend/src/components/admin/pipeline/PipelineFormulaMetrics.vue (100行以下 - SPEC-PRINCIPLE-001) -->
<script setup lang="ts">
import { computed } from 'vue';
import type { PipelineMetrics } from '../../../types/pipeline';

const props = defineProps<{ metrics?: PipelineMetrics }>();
const m = computed(() => props.metrics || {
  media_count: 0, total_media_count: 0, total_tasks: 0,
  registered_tasks: 0, active_tasks: 0, pending_tasks: 0, success_tasks: 0, failed_tasks: 0
});
</script>

<template>
  <div class="space-y-2.5 font-sans select-none">
    <!-- レベル1: 作品数 と タスク数の関係式 -->
    <div class="flex flex-wrap items-center gap-2 text-xs font-mono">
      <div class="flex items-center gap-1.5 px-3 py-1.5 bg-slate-950 rounded-xl border border-indigo-900/60 shadow-inner">
        <span class="text-slate-400 text-[11px]">作品数 (代表ID):</span>
        <span class="font-bold text-indigo-300 text-sm">{{ m.media_count }}</span>
        <span class="text-[10px] text-slate-500">(全 {{ m.total_media_count }} 作品)</span>
      </div>
      <div class="text-slate-500 font-bold text-[11px] flex items-center gap-1">
        <span>──</span><span class="px-1.5 py-0.5 rounded bg-slate-800 text-indigo-300 text-[10px] border border-slate-700">× 7 候補</span><span>──▶</span>
      </div>
      <div class="flex items-center gap-1.5 px-3 py-1.5 bg-slate-950 rounded-xl border border-purple-900/60 shadow-inner">
        <span class="text-slate-400 text-[11px]">タスク総数:</span>
        <span class="font-bold text-purple-300 text-sm">{{ m.total_tasks }}</span>
        <span class="text-[10px] text-slate-500">個別URL</span>
      </div>
    </div>

    <!-- レベル2: タスク内訳 (ダウンロード登録数 + 成功数 + 失敗数) -->
    <div class="grid grid-cols-1 sm:grid-cols-3 gap-2">
      <!-- 下載登録数 (稼働数 + 継続探査数) -->
      <div class="p-2.5 bg-slate-950/90 rounded-xl border border-sky-900/50 space-y-1.5">
        <div class="flex items-center justify-between">
          <span class="text-[11px] font-bold text-sky-400 font-mono">下載登録数</span>
          <span class="text-xs font-mono font-bold text-sky-200">{{ m.registered_tasks }} タスク</span>
        </div>
        <div class="space-y-1 pt-1 border-t border-slate-800/80 text-[10.5px] font-mono">
          <div class="flex justify-between items-center text-slate-300">
            <span class="flex items-center gap-1"><span class="w-1.5 h-1.5 rounded-full bg-emerald-400 shadow-[0_0_4px_#34d399]"></span>稼働数 (ONBOARDED):</span>
            <span class="font-bold text-emerald-300">{{ m.active_tasks }}</span>
          </div>
          <div class="flex justify-between items-center text-slate-300">
            <span class="flex items-center gap-1"><span class="w-1.5 h-1.5 rounded-full" :class="m.pending_tasks > 0 ? 'bg-amber-400 shadow-[0_0_4px_#fbbf24]' : 'bg-slate-600'"></span>継続探査数 (PENDING):</span>
            <span class="font-bold" :class="m.pending_tasks > 0 ? 'text-amber-300' : 'text-slate-400'">{{ m.pending_tasks }}</span>
          </div>
        </div>
        <div v-if="m.pending_tasks > 0" class="text-[9px] text-amber-400/90 bg-amber-950/40 px-1.5 py-0.5 rounded border border-amber-900/50 leading-tight">
          ⚠️ 差し戻し再登録 (要フロー確認)
        </div>
      </div>

      <!-- 成功数 (COMPLETED) -->
      <div class="p-2.5 bg-slate-950/90 rounded-xl border border-emerald-900/50 flex flex-col justify-between">
        <div class="flex items-center justify-between">
          <span class="text-[11px] font-bold text-emerald-400 font-mono">成功数</span>
          <span class="text-xs font-mono font-bold text-emerald-300">{{ m.success_tasks }} タスク</span>
        </div>
        <div class="text-[10px] text-slate-400 font-mono pt-2 border-t border-slate-800/80">
          <span>COMPLETED 救出完了タスク</span>
        </div>
      </div>

      <!-- 失敗数 (RETIRED -> RETAINED) -->
      <div class="p-2.5 bg-slate-950/90 rounded-xl border border-rose-900/50 flex flex-col justify-between">
        <div class="flex items-center justify-between">
          <span class="text-[11px] font-bold text-rose-400 font-mono">失敗数</span>
          <span class="text-xs font-mono font-bold text-rose-300">{{ m.failed_tasks }} タスク</span>
        </div>
        <div class="text-[10px] text-slate-400 font-mono pt-2 border-t border-slate-800/80">
          <span>RETIRED / RETAINED 枯渇メディア</span>
        </div>
      </div>
    </div>

    <!-- レベル3: 関係式整合フッター -->
    <div class="flex items-center justify-between px-2 py-1 bg-slate-950/60 rounded-lg border border-slate-800/80 text-[10px] font-mono text-slate-400">
      <span>📐 関係式: タスク数({{ m.total_tasks }}) ＝ 下載登録({{ m.registered_tasks }}) ＋ 成功({{ m.success_tasks }}) ＋ 失敗({{ m.failed_tasks }})</span>
      <span v-if="m.total_tasks === (m.registered_tasks + m.success_tasks + m.failed_tasks)" class="text-emerald-400 font-bold">● 整合</span>
      <span v-else class="text-amber-400">▲ 収束中 ({{ m.total_tasks - (m.registered_tasks + m.success_tasks + m.failed_tasks) }}差分)</span>
    </div>
  </div>
</template>
