<!-- frontend/src/components/admin/DownloaderConsoleView.vue (100行以下 - SPEC-PRINCIPLE-001) -->
<script setup lang="ts">
import { computed } from 'vue';
import { useMotrixQueue } from '../../composables/admin/useMotrixQueue';
import { useDownloaderBatchOps } from '../../composables/admin/useDownloaderBatchOps';
import MotrixTaskCard from './MotrixTaskCard.vue';
import MotrixReserveCard from './MotrixReserveCard.vue';

defineProps<{ admin?: any }>();
const { activeTasks, waitingTasks, stoppedTasks, reserves, globalStat, loading, activeTab, fetchQueue, controlMotrix, restoreReserve } = useMotrixQueue();
const { selectedCount, isOperating, isSelected, toggleSelect, selectAll, selectOnlyErrors, clearSelection, batchControl, batchSafePurge, batchEscalateToThunder } = useDownloaderBatchOps(fetchQueue);

const currentList = computed(() => {
  if (activeTab.value === 'active') return activeTasks.value;
  if (activeTab.value === 'waiting') return waitingTasks.value;
  if (activeTab.value === 'stopped') return stoppedTasks.value;
  return [];
});
const formatSpeed = (b: number) => (!b || b === 0 ? '0 B/s' : b < 1048576 ? `${(b / 1024).toFixed(1)} KB/s` : `${(b / 1048576).toFixed(2)} MB/s`);
</script>

<template>
  <div class="h-full w-full flex flex-col p-3 sm:p-4 space-y-3 bg-slate-950 text-slate-100 overflow-y-auto font-sans">
    <div class="flex items-center justify-between border-b border-slate-800 pb-2">
      <div class="flex items-center gap-2"><span class="text-xl">🚀</span><div><h2 class="text-sm font-bold text-slate-100">Motrix / Aria2 キュー統合管理</h2><p class="text-[10px] text-slate-400 font-mono">OUTSOURCED キューの遠隔制御・安全退避パージ</p></div></div>
      <button @click="fetchQueue" class="px-2.5 py-1 bg-slate-800 hover:bg-slate-700 text-slate-300 rounded text-xs cursor-pointer">🔄 更新</button>
    </div>

    <!-- グローバルステータス & 制御パネル -->
    <div class="p-3 bg-slate-900/90 border border-slate-800 rounded-2xl space-y-2 shadow-lg">
      <div class="flex items-center justify-between">
        <div class="flex items-center gap-2"><span class="text-xs font-bold text-slate-300">Motrix Next</span><span :class="globalStat?.is_online ? 'bg-emerald-950 text-emerald-300 border-emerald-700/60' : 'bg-rose-950 text-rose-300 border-rose-700/60'" class="px-2 py-0.5 rounded text-[10px] font-mono font-bold border">{{ globalStat?.is_online ? '🟢 ONLINE' : '🔴 OFFLINE' }}</span></div>
        <div class="flex gap-1">
          <button @click="controlMotrix('safe_limits')" :disabled="loading" class="px-2 py-0.5 bg-indigo-950 hover:bg-indigo-900 text-indigo-200 rounded text-[11px] font-semibold cursor-pointer">🛡️ 安全2並列</button>
          <button @click="controlMotrix('pause_all')" :disabled="loading" class="px-2 py-0.5 bg-amber-950 hover:bg-amber-900 text-amber-200 rounded text-[11px] font-semibold cursor-pointer">⏸️ 全停止</button>
          <button @click="controlMotrix('unpause_all')" :disabled="loading" class="px-2 py-0.5 bg-blue-950 hover:bg-blue-900 text-blue-200 rounded text-[11px] font-semibold cursor-pointer">▶️ 全再開</button>
          <button @click="controlMotrix('purge_all')" :disabled="loading" class="px-2 py-0.5 bg-rose-950 hover:bg-rose-900 text-rose-200 rounded text-[11px] font-semibold cursor-pointer">🧹 履歴削除</button>
        </div>
      </div>
      <div class="grid grid-cols-4 gap-2 text-center pt-1">
        <div class="p-1.5 bg-slate-950 rounded-xl border border-slate-800"><div class="text-[10px] text-slate-400">稼働 (Active)</div><div class="text-xs font-mono font-bold text-blue-400">{{ globalStat?.num_active ?? 0 }} 件</div></div>
        <div class="p-1.5 bg-slate-950 rounded-xl border border-slate-800"><div class="text-[10px] text-slate-400">待機 (Waiting)</div><div class="text-xs font-mono font-bold text-amber-400">{{ globalStat?.num_waiting ?? 0 }} 件</div></div>
        <div class="p-1.5 bg-slate-950 rounded-xl border border-slate-800"><div class="text-[10px] text-slate-400">停止 (Stopped)</div><div class="text-xs font-mono font-bold text-slate-300">{{ globalStat?.num_stopped ?? 0 }} 件</div></div>
        <div class="p-1.5 bg-slate-950 rounded-xl border border-slate-800"><div class="text-[10px] text-slate-400">速度</div><div class="text-xs font-mono font-bold text-purple-400">{{ formatSpeed(globalStat?.download_speed) }}</div></div>
      </div>
    </div>

    <!-- タブ切り替え & 一括操作ツールバー -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-2 border-b border-slate-800/80 pb-2">
      <div class="flex items-center gap-1.5 bg-slate-900/80 p-1 rounded-xl border border-slate-800">
        <button @click="activeTab = 'active'" :class="activeTab === 'active' ? 'bg-blue-600 text-white font-bold' : 'text-slate-400 hover:text-slate-200'" class="px-2.5 py-1 rounded-lg text-xs cursor-pointer">Active ({{ activeTasks.length }})</button>
        <button @click="activeTab = 'waiting'" :class="activeTab === 'waiting' ? 'bg-amber-600 text-white font-bold' : 'text-slate-400 hover:text-slate-200'" class="px-2.5 py-1 rounded-lg text-xs cursor-pointer">Waiting ({{ waitingTasks.length }})</button>
        <button @click="activeTab = 'stopped'" :class="activeTab === 'stopped' ? 'bg-purple-600 text-white font-bold' : 'text-slate-400 hover:text-slate-200'" class="px-2.5 py-1 rounded-lg text-xs cursor-pointer">Stopped ({{ stoppedTasks.length }})</button>
        <button @click="activeTab = 'reserves'" :class="activeTab === 'reserves' ? 'bg-emerald-600 text-white font-bold' : 'text-slate-400 hover:text-slate-200'" class="px-2.5 py-1 rounded-lg text-xs cursor-pointer">🛡️ 退避保管 ({{ reserves.length }})</button>
      </div>

      <div v-if="activeTab !== 'reserves'" class="flex flex-wrap items-center gap-1.5">
        <span class="text-[11px] font-mono text-slate-400">選択: {{ selectedCount }}</span>
        <button @click="selectAll(currentList)" class="px-2 py-0.5 bg-slate-800 hover:bg-slate-700 text-slate-300 rounded text-[10px] cursor-pointer">全選択</button>
        <button @click="selectOnlyErrors(currentList)" class="px-2 py-0.5 bg-rose-950 hover:bg-rose-900 text-rose-300 rounded text-[10px] cursor-pointer">エラー選択</button>
        <button @click="clearSelection" class="px-2 py-0.5 bg-slate-800 hover:bg-slate-700 text-slate-400 rounded text-[10px] cursor-pointer">クリア</button>
        <div class="h-3 w-px bg-slate-700 mx-0.5"></div>
        <button @click="batchControl('pause')" :disabled="selectedCount === 0 || isOperating" class="px-2 py-0.5 bg-amber-900/80 hover:bg-amber-800 text-amber-200 rounded text-[10px] font-bold cursor-pointer disabled:opacity-30">⏸️ 停止</button>
        <button @click="batchControl('unpause')" :disabled="selectedCount === 0 || isOperating" class="px-2 py-0.5 bg-blue-900/80 hover:bg-blue-800 text-blue-200 rounded text-[10px] font-bold cursor-pointer disabled:opacity-30">▶️ 再開</button>
        <button @click="batchSafePurge" :disabled="selectedCount === 0 || isOperating" class="px-2 py-0.5 bg-emerald-900/80 hover:bg-emerald-800 text-emerald-200 rounded text-[10px] font-bold cursor-pointer disabled:opacity-30">🛡️ 退避削除</button>
        <button @click="batchEscalateToThunder(currentList)" :disabled="selectedCount === 0 || isOperating" class="px-2 py-0.5 bg-purple-900/80 hover:bg-purple-800 text-purple-200 rounded text-[10px] font-bold cursor-pointer disabled:opacity-30">⚡ 迅雷転送</button>
      </div>
    </div>

    <!-- リスト表示領域 -->
    <div class="flex-1 bg-slate-900/90 border border-slate-800 rounded-2xl p-3 flex flex-col space-y-2 shadow-lg overflow-hidden min-h-[300px]">
      <div v-if="activeTab !== 'reserves'" class="flex-1 overflow-y-auto space-y-2 pr-1">
        <div v-if="currentList.length === 0" class="text-center py-12 text-slate-500 text-xs font-mono">該当するタスクはありません</div>
        <MotrixTaskCard v-for="t in currentList" :key="t.gid" :task="t" :is-selected="isSelected(t.gid)" @toggle="toggleSelect" />
      </div>
      <div v-else class="flex-1 overflow-y-auto space-y-2 pr-1">
        <div v-if="reserves.length === 0" class="text-center py-12 text-slate-500 text-xs font-mono">退避されたタスクはありません</div>
        <MotrixReserveCard v-for="r in reserves" :key="r.id" :reserve="r" @restore="restoreReserve" />
      </div>
    </div>
  </div>
</template>
