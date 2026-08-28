<!-- frontend/src/components/admin/ThunderOrchestratorView.vue (100行以下 - SPEC-PRINCIPLE-001) -->
<script setup lang="ts">
import { ref } from 'vue';
import { useThunderOrchestrator } from '../../composables/admin/useThunderOrchestrator';
import ThunderSlotCard from './ThunderSlotCard.vue';
import ThunderProgressPanel from './ThunderProgressPanel.vue';
import ThunderConfirmModal from './ThunderConfirmModal.vue';

defineProps<{ admin?: any }>();
const { status, loading, intervalSec, tempDir, startOrchestrator, pauseOrchestrator, resumeOrchestrator, stopOrchestrator, launchThunder, syncDownloads, fetchStatus } = useThunderOrchestrator();

const showConfirmModal = ref(false);
const handleStartConfirmed = (sec: number) => {
  intervalSec.value = sec;
  startOrchestrator();
};
</script>

<template>
  <div class="h-full flex flex-col p-3 sm:p-4 space-y-3 bg-slate-950 text-slate-100 overflow-y-auto font-sans max-w-6xl mx-auto">
    <div class="flex items-center justify-between border-b border-slate-800 pb-2">
      <div class="flex items-center gap-2"><span class="text-xl">⚡</span><div><h2 class="text-sm font-bold text-slate-100">迅雷 (Thunder) COM オーケストレーター</h2><p class="text-[10px] text-slate-400 font-mono">*.xltd 生起監視・ESCALATED ➔ RETAINED ➔ COMPLETED 10ステップ自律同期</p></div></div>
      <div class="flex items-center gap-2">
        <button @click="launchThunder" class="px-2.5 py-1 bg-amber-600 hover:bg-amber-500 text-white font-bold rounded text-xs active:scale-95 cursor-pointer shadow">🚀 Thunder 起動</button>
        <button @click="fetchStatus" class="px-2.5 py-1 bg-slate-800 hover:bg-slate-700 text-slate-300 rounded text-xs active:scale-95 cursor-pointer">🔄 更新</button>
      </div>
    </div>

    <!-- オーケストレーション制御パネル -->
    <div class="p-3 bg-slate-900/90 border border-slate-800 rounded-2xl space-y-2.5 shadow-lg">
      <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-2">
        <div class="flex items-center gap-2">
          <span class="text-xs font-bold text-slate-200">稼働:</span>
          <span v-if="status?.is_running && !status?.is_paused" class="px-2 py-0.5 rounded text-[10px] font-mono font-bold bg-purple-950 text-purple-300 border border-purple-700 animate-pulse">⚡ 実行中</span>
          <span v-else-if="status?.is_paused" class="px-2 py-0.5 rounded text-[10px] font-mono font-bold bg-amber-950 text-amber-300 border border-amber-700">⏸️ 一時停止</span>
          <span v-else class="px-2 py-0.5 rounded text-[10px] font-mono font-bold bg-slate-800 text-slate-400 border border-slate-700">⏹️ 停止中</span>
        </div>
        <div class="flex items-center gap-1.5">
          <button v-if="!status?.is_running" @click="showConfirmModal = true" :disabled="loading" class="px-3 py-1 bg-gradient-to-r from-purple-600 to-indigo-600 hover:from-purple-500 text-white font-bold rounded-lg text-xs shadow cursor-pointer disabled:opacity-50">🚀 エスカレーション開始 ({{ status?.total_jobs || 303 }}ジョブ)</button>
          <template v-else>
            <button v-if="!status?.is_paused" @click="pauseOrchestrator" class="px-2.5 py-1 bg-amber-600 hover:bg-amber-500 text-white font-bold rounded-lg text-xs cursor-pointer">⏸️ 一時停止</button>
            <button v-else @click="resumeOrchestrator" class="px-2.5 py-1 bg-emerald-600 hover:bg-emerald-500 text-white font-bold rounded-lg text-xs cursor-pointer">▶️ 再開</button>
            <button @click="stopOrchestrator" class="px-2.5 py-1 bg-rose-600 hover:bg-rose-500 text-white font-bold rounded-lg text-xs cursor-pointer">🛑 中断</button>
          </template>
        </div>
      </div>
      <ThunderProgressPanel :status="status" />

      <!-- テンポラリフォルダからの手動取り込み・移動バー -->
      <div class="flex items-center justify-between gap-2 p-2 bg-slate-950/80 rounded-xl border border-slate-800">
        <div class="flex items-center gap-2 flex-1 min-w-0">
          <span class="text-xs text-amber-400 font-bold shrink-0">📂 迅雷保存先:</span>
          <input v-model="tempDir" class="bg-slate-900 border border-slate-700 rounded px-2 py-0.5 text-xs font-mono text-slate-200 flex-1 min-w-0" placeholder="D:\迅雷下载" />
        </div>
        <button @click="syncDownloads" :disabled="loading" class="px-3 py-1 bg-emerald-700 hover:bg-emerald-600 text-white rounded text-xs font-bold shrink-0 cursor-pointer active:scale-95 disabled:opacity-50">
          📦 完了ファイルをアカウントフォルダへ移動＆同期
        </button>
      </div>
    </div>

    <!-- 12スロット待機制御グリッド -->
    <div class="bg-slate-900/90 border border-slate-800 rounded-2xl p-3 space-y-2 shadow-lg">
      <div class="flex items-center justify-between"><span class="text-xs font-bold text-slate-200">迅雷 COM 投入スロット (最大 12 並列待機枠)</span><span class="text-[10px] font-mono text-slate-400">*.xltd 監視 / 投入間隔: {{ intervalSec }}秒</span></div>
      <div class="grid grid-cols-3 sm:grid-cols-4 md:grid-cols-6 gap-2">
        <ThunderSlotCard v-for="s in status?.slots || Array.from({ length: 12 }, (_, i) => ({ index: i, is_occupied: false }))" :key="s.index" :slot="s" />
      </div>
    </div>

    <!-- 最近のディスパッチログ -->
    <div class="flex-1 bg-slate-900/90 border border-slate-800 rounded-2xl p-3 flex flex-col space-y-2 shadow-lg overflow-hidden min-h-[160px]">
      <span class="text-xs font-bold text-slate-200">ディスパッチ実行ログ</span>
      <div class="flex-1 overflow-y-auto space-y-1.5 pr-1">
        <div v-if="!status?.recent_tasks || status.recent_tasks.length === 0" class="text-center py-6 text-slate-500 text-xs font-mono">ログはありません</div>
        <div v-for="t in status?.recent_tasks" :key="t.id" class="p-2 bg-slate-950 rounded-lg border border-slate-800/80 flex items-center justify-between text-xs">
          <div class="space-y-0.5 truncate flex-1 pr-2">
            <div class="flex items-center gap-2"><span class="font-mono font-bold text-slate-200">{{ t.media_id }}</span><span class="px-1.5 py-0.2 rounded text-[9px] font-mono bg-indigo-950 text-indigo-300 border border-indigo-800">{{ t.resolution_type }}</span></div>
            <div class="text-[10px] text-slate-400 font-mono truncate">{{ t.url }}</div>
          </div>
          <span class="px-2 py-0.5 rounded text-[10px] font-mono font-bold bg-purple-950 text-purple-300 border border-purple-800 shrink-0">DISPATCHED</span>
        </div>
      </div>
    </div>

    <!-- 免責・注意事項モーダル -->
    <ThunderConfirmModal :is-open="showConfirmModal" :job-count="status?.total_jobs || 303" @close="showConfirmModal = false" @confirm="handleStartConfirmed" />
  </div>
</template>
