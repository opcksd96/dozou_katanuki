<!-- frontend/src/components/admin/PipelineConsoleView.vue (100行以下 - SPEC-PRINCIPLE-001) -->
<script setup lang="ts">
import { ref } from 'vue';
import { usePipelineConsole } from '../../composables/admin/usePipelineConsole';
import { ResetAllToQueuedAndBootstrap, IgnitePipeline } from '../../../wailsjs/go/app/App';
import PipelineCheckpointRibbon from './pipeline/PipelineCheckpointRibbon.vue';
import PipelineLogViewer from './pipeline/PipelineLogViewer.vue';

const { overview, logs, selectedLogStage, loading, isAutoEngineRunning, toggleAutoEngine, refreshAll, syncAndReconcile, setLogStage } = usePipelineConsole();
const logViewerRef = ref<HTMLElement | null>(null);
const scrollToLogs = () => { logViewerRef.value?.scrollIntoView({ behavior: 'smooth' }); };

const handleManualIgnite = async () => {
  try {
    const res = await IgnitePipeline();
    alert(`🔥 点火完了: QUEUED ${res.queued_count} 件 / 迅雷 ${res.escalated_count} 件 を再発火しました！`);
    await refreshAll();
  } catch (e: any) { alert(`エラー: ${e?.message || e}`); }
};

const handleResetAll = async () => {
  if (!confirm('全タスクを QUEUED に初期化して最上流から流し直しますか？')) return;
  try {
    const count = await ResetAllToQueuedAndBootstrap();
    alert(`初期化完了: ${count} 件を QUEUED に差し戻しました。完全自動運転により最上流から順次処理されます。`);
    await refreshAll();
  } catch (e: any) { alert(`初期化エラー: ${e?.message || e}`); }
};
</script>

<template>
  <div class="h-full w-full flex flex-col p-3 sm:p-4 space-y-3 bg-slate-950 text-slate-100 overflow-y-auto font-sans">
    <div class="flex items-center justify-between border-b border-slate-800 pb-2">
      <div class="flex items-center gap-2">
        <span class="text-xl">🚀</span>
        <div><h2 class="text-sm font-bold text-slate-100">統合ダウンロードパイプライン</h2><p class="text-[10px] text-slate-400 font-mono">Requests ➔ Motrix Next ➔ 迅雷 ➔ Stash 完全自律ストリーム</p></div>
      </div>
      <div class="flex items-center gap-1.5">
        <button @click="scrollToLogs" class="px-2.5 py-1 bg-slate-800 hover:bg-slate-700 text-slate-300 rounded text-xs font-semibold cursor-pointer active:scale-95">📜 ログ</button>
        <button @click="handleResetAll" class="px-2.5 py-1 bg-slate-900 hover:bg-rose-950/80 border border-slate-800 hover:border-rose-800 text-slate-400 hover:text-rose-300 rounded text-xs font-mono cursor-pointer active:scale-95" title="全タスクをQUEUEDに初期化して最上流から流し直す">🧹 QUEUED初期化</button>
        <button @click="refreshAll" :disabled="loading" class="px-2.5 py-1 bg-slate-800 hover:bg-slate-700 text-slate-300 rounded text-xs cursor-pointer">🔄 更新</button>
      </div>
    </div>

    <PipelineCheckpointRibbon :checkpoints="overview?.checkpoints || []" />

    <!-- 🌟 パイプライン完全自動運転コントロールバー -->
    <div class="p-3 bg-slate-900/90 border border-indigo-900/60 rounded-2xl space-y-2.5 shadow-lg">
      <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-2 border-b border-slate-800/80 pb-2">
        <div class="flex items-center gap-2">
          <span class="text-xs font-bold text-indigo-300">🎮 自律パイプライン運転:</span>
          <span v-if="isAutoEngineRunning" class="px-2.5 py-0.5 rounded text-[11px] font-mono font-bold bg-emerald-950 text-emerald-300 border border-emerald-700 flex items-center gap-1 animate-pulse">
            <span class="w-1.5 h-1.5 rounded-full bg-emerald-400"></span> 🟢 完全自動運転中 (常駐ループ)
          </span>
          <span v-else class="px-2.5 py-0.5 rounded text-[11px] font-mono font-bold bg-amber-950 text-amber-300 border border-amber-700">
            ⏸️ 自動運転一時停止中
          </span>
        </div>
        <div class="flex items-center gap-1.5">
          <button @click="toggleAutoEngine" :class="['px-3 py-1 font-bold rounded-lg text-xs cursor-pointer active:scale-95 shadow transition', isAutoEngineRunning ? 'bg-amber-600 hover:bg-amber-500 text-white' : 'bg-emerald-600 hover:bg-emerald-500 text-white']">
            {{ isAutoEngineRunning ? '⏸️ 自動運転を一時停止' : '▶️ 自動運転を開始' }}
          </button>
          <button @click="handleManualIgnite" class="px-2.5 py-1 bg-orange-600 hover:bg-orange-500 text-white font-bold rounded-lg text-xs cursor-pointer active:scale-95 shadow" title="今すぐ全ステージの未処理タスクを強制点火">
            🔥 即時点火
          </button>
          <button @click="syncAndReconcile" class="px-2.5 py-1 bg-slate-800 hover:bg-slate-700 text-slate-300 font-bold rounded-lg text-xs cursor-pointer active:scale-95" title="Stash同期を手動キック">
            🎬 Stash手動同期
          </button>
        </div>
      </div>

      <div class="space-y-1.5">
        <div class="flex items-center justify-between text-xs font-mono"><span class="text-slate-400">パイプライン全体救出進捗 (全 {{ overview?.total_media ?? 0 }} 件)</span><span class="font-bold text-purple-400">{{ (overview?.overall_progress ?? 0).toFixed(1) }}%</span></div>
        <div class="w-full bg-slate-950 rounded-full h-2 overflow-hidden border border-slate-800"><div class="bg-gradient-to-r from-blue-500 via-sky-500 via-purple-500 to-emerald-500 h-full rounded-full transition-all duration-300" :style="{ width: `${overview?.overall_progress ?? 0}%` }"></div></div>
        <div class="grid grid-cols-4 gap-2 text-center pt-1">
          <div class="p-1.5 bg-slate-950 rounded-xl border border-slate-800"><div class="text-[10px] text-slate-400 font-medium">救出完了</div><div class="text-xs font-mono font-bold text-emerald-400">{{ overview?.completed ?? 0 }} 件</div></div>
          <div class="p-1.5 bg-slate-950 rounded-xl border border-slate-800"><div class="text-[10px] text-slate-400 font-medium">迅雷探索中</div><div class="text-xs font-mono font-bold text-purple-400">{{ overview?.escalated ?? 0 }} 件</div></div>
          <div class="p-1.5 bg-slate-950 rounded-xl border border-slate-800"><div class="text-[10px] text-slate-400 font-medium">Motrix待機</div><div class="text-xs font-mono font-bold text-blue-400">{{ overview?.outsourced ?? 0 }} 件</div></div>
          <div class="p-1.5 bg-slate-950 rounded-xl border border-slate-800"><div class="text-[10px] text-slate-400 font-medium">司令塔リキュー待ち</div><div class="text-xs font-mono font-bold text-rose-400">{{ overview?.retained ?? 0 }} 件</div></div>
        </div>
      </div>
    </div>

    <div ref="logViewerRef"><PipelineLogViewer :logs="logs" :selected-stage="selectedLogStage" @select-stage="(s) => setLogStage(s)" /></div>
  </div>
</template>
