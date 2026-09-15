<!-- frontend/src/components/admin/PipelineConsoleView.vue (100行以下 - SPEC-PRINCIPLE-001) -->
<script setup lang="ts">
import { ref } from 'vue';
import { usePipelineConsole } from '../../composables/admin/usePipelineConsole';
import { ResetAllToQueuedAndBootstrap } from '../../../wailsjs/go/app/App';
import PipelineUnifiedCard from './pipeline/PipelineUnifiedCard.vue';
import PipelineLogViewer from './pipeline/PipelineLogViewer.vue';
import PipelineMediaAccordionViewer from './pipeline/PipelineMediaAccordionViewer.vue';

const { overview, logs, selectedLogStage, loading, isAutoEngineRunning, toggleAutoEngine, refreshAll, executePipelineCycleNow, setLogStage } = usePipelineConsole();
const activeBottomTab = ref<'carousel' | 'raw'>('carousel');
const isWails = () => typeof (window as any).go !== 'undefined';

const handleExecuteCycleNow = async () => {
  try { await executePipelineCycleNow(); await refreshAll(); } catch (e: any) { alert(`実行エラー: ${e?.message || e}`); }
};

const handleResetAll = async () => {
  if (!confirm('全タスクを QUEUED に初期化して最上流から流し直しますか？')) return;
  try {
    const count = isWails() ? await ResetAllToQueuedAndBootstrap() : await (await fetch('/api/admin/pipeline/reset-all', { method: 'POST' })).json();
    alert(`初期化完了: ${count} 件を QUEUED に差し戻しました。`);
    await refreshAll();
  } catch (e: any) { alert(`初期化エラー: ${e?.message || e}`); }
};
</script>

<template>
  <div class="w-full flex flex-col p-3 sm:p-4 space-y-4 bg-slate-950 text-slate-100 font-sans min-h-full">
    <div class="flex items-center justify-between border-b border-slate-800 pb-2">
      <div>
        <h2 class="text-sm font-bold text-slate-100 tracking-wide">統合ダウンロードパイプライン</h2>
        <p class="text-[10px] text-slate-400 font-mono">作品救出進捗 ＆ 通信タスク自律オーケストレーション</p>
      </div>
      <div class="flex items-center gap-1.5">
        <button @click="handleResetAll" class="px-2.5 py-1 bg-slate-900 hover:bg-rose-950/80 border border-slate-800 text-slate-400 hover:text-rose-300 rounded text-xs font-mono cursor-pointer active:scale-95" title="全タスク初期化">QUEUED初期化</button>
        <button @click="refreshAll" :disabled="loading" class="px-2.5 py-1 bg-slate-800 hover:bg-slate-700 text-slate-300 rounded text-xs cursor-pointer">更新</button>
      </div>
    </div>

    <!-- 統合パイプライン管理カード (進捗・メトリクス・チェックポイントを1枚に集約) -->
    <PipelineUnifiedCard
      :overview="overview"
      :loading="loading"
      :is-auto-engine-running="isAutoEngineRunning"
      @toggle-auto-engine="toggleAutoEngine"
      @execute-cycle-now="handleExecuteCycleNow"
    />

    <!-- ペイン3: 詳細探査・ログ (全展開・全体スクロール) -->
    <div class="flex flex-col space-y-2 pt-1 w-full">
      <div class="flex items-center justify-between border-b border-slate-800/80 pb-1.5">
        <div class="flex items-center gap-1 bg-slate-900 p-1 rounded-xl border border-slate-800">
          <button @click="activeBottomTab = 'carousel'" :class="['px-3 py-1 rounded-lg text-xs font-bold transition cursor-pointer', activeBottomTab === 'carousel' ? 'bg-purple-600 text-white shadow' : 'text-slate-400 hover:text-slate-200']">作品別・候補URLコラプス探査</button>
          <button @click="activeBottomTab = 'raw'" :class="['px-3 py-1 rounded-lg text-xs font-bold transition cursor-pointer', activeBottomTab === 'raw' ? 'bg-indigo-600 text-white shadow' : 'text-slate-400 hover:text-slate-200']">生ログタイムライン</button>
        </div>
      </div>
      <div class="w-full">
        <PipelineMediaAccordionViewer v-if="activeBottomTab === 'carousel'" />
        <PipelineLogViewer v-else :logs="logs" :selected-stage="selectedLogStage" @select-stage="(s) => setLogStage(s)" />
      </div>
    </div>
  </div>
</template>
