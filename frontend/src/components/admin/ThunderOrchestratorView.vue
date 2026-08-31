<!-- frontend/src/components/admin/ThunderOrchestratorView.vue (100行以下 - SPEC-PRINCIPLE-001) -->
<script setup lang="ts">
import { ref } from 'vue';
import { useThunderOrchestrator } from '../../composables/admin/useThunderOrchestrator';
import ThunderProgressPanel from './ThunderProgressPanel.vue';
import ThunderConfirmModal from './ThunderConfirmModal.vue';
import ThunderCDPTaskList from './ThunderCDPTaskList.vue';

defineProps<{ admin?: any }>();
const { status, loading, intervalSec, tempDir, startOrchestrator, pauseOrchestrator, resumeOrchestrator, stopOrchestrator, launchThunder, syncDownloads, fetchStatus } = useThunderOrchestrator();
const showConfirmModal = ref(false), activeTab = ref<'running' | 'all'>('running');
const handleStartConfirmed = (sec: number) => { intervalSec.value = sec; startOrchestrator(); };
</script>

<template>
  <div class="h-full w-full flex flex-col p-3 sm:p-4 space-y-3 bg-slate-950 text-slate-100 overflow-y-auto font-sans">
    <!-- 統合ヘッダー (Motrix と完全同一のレイアウト・サブタイトル形式) -->
    <div class="flex items-center justify-between border-b border-slate-800 pb-2">
      <div class="flex items-center gap-2">
        <span class="text-xl">⚡</span>
        <div>
          <h2 class="text-sm font-bold text-slate-100">迅雷 (Thunder) キュー統合管理</h2>
          <p class="text-[10px] text-slate-400 font-mono">P2SP 救出エスカレーション・重複解消・Stash 自動連携</p>
        </div>
      </div>
      <div class="flex items-center gap-1.5">
        <button @click="launchThunder" class="px-2.5 py-1 bg-amber-600 hover:bg-amber-500 text-white font-bold rounded text-xs transition cursor-pointer flex items-center gap-1 active:scale-95 shadow">
          🚀 Thunder 起動
        </button>
        <button @click="fetchStatus" class="px-2.5 py-1 bg-slate-800 hover:bg-slate-700 text-slate-300 rounded text-xs cursor-pointer">
          🔄 更新
        </button>
      </div>
    </div>

    <!-- グローバルステータス & 制御パネル -->
    <div class="p-3 bg-slate-900/90 border border-slate-800 rounded-2xl space-y-2 shadow-lg">
      <div class="flex items-center justify-between">
        <div class="flex items-center gap-2">
          <span class="text-xs font-bold text-slate-300">Thunder Orchestrator</span>
          <span v-if="status?.is_running && !status?.is_paused" class="px-2 py-0.5 rounded text-[10px] font-mono font-bold bg-purple-950 text-purple-300 border border-purple-700 animate-pulse">⚡ RUNNING</span>
          <span v-else-if="status?.is_paused" class="px-2 py-0.5 rounded text-[10px] font-mono font-bold bg-amber-950 text-amber-300 border border-amber-700">⏸️ PAUSED</span>
          <span v-else class="px-2 py-0.5 rounded text-[10px] font-mono font-bold bg-slate-800 text-slate-400 border border-slate-700">⏹️ STOPPED</span>
        </div>
        <div class="flex items-center gap-1">
          <button v-if="!status?.is_running" @click="showConfirmModal = true" :disabled="loading" class="px-2.5 py-0.5 bg-gradient-to-r from-purple-600 to-indigo-600 hover:from-purple-500 text-white rounded text-[11px] font-semibold cursor-pointer disabled:opacity-50">
            🚀 救出開始 ({{ status?.total_jobs || 303 }})
          </button>
          <template v-else>
            <button v-if="!status?.is_paused" @click="pauseOrchestrator" class="px-2 py-0.5 bg-amber-950 hover:bg-amber-900 text-amber-200 rounded text-[11px] font-semibold cursor-pointer">⏸️ 一時停止</button>
            <button v-else @click="resumeOrchestrator" class="px-2 py-0.5 bg-emerald-950 hover:bg-emerald-900 text-emerald-200 rounded text-[11px] font-semibold cursor-pointer">▶️ 再開</button>
            <button @click="stopOrchestrator" class="px-2 py-0.5 bg-rose-950 hover:bg-rose-900 text-rose-200 rounded text-[11px] font-semibold cursor-pointer">🛑 中断</button>
          </template>
        </div>
      </div>

      <ThunderProgressPanel :status="status" />

      <!-- 保存先フォルダ & 同期バー -->
      <div class="flex items-center justify-between gap-2 p-2 bg-slate-950 rounded-xl border border-slate-800">
        <div class="flex items-center gap-2 flex-1 min-w-0">
          <span class="text-[11px] text-amber-400 font-bold shrink-0">📂 保存先:</span>
          <input v-model="tempDir" class="bg-slate-900 border border-slate-700 rounded px-2 py-0.5 text-xs font-mono text-slate-200 flex-1 min-w-0" placeholder="D:\迅雷下载" />
        </div>
        <button @click="syncDownloads" :disabled="loading" class="px-3 py-1 bg-emerald-700 hover:bg-emerald-600 text-white rounded text-xs font-bold shrink-0 cursor-pointer active:scale-95 disabled:opacity-50">
          📦 完了回収 & Stash連携
        </button>
      </div>
    </div>

    <!-- タブ切り替え & リスト領域 (Motrix と完全同一の構造) -->
    <div class="flex items-center justify-between border-b border-slate-800/80 pb-2">
      <div class="flex items-center gap-1.5 bg-slate-900/80 p-1 rounded-xl border border-slate-800">
        <button @click="activeTab = 'running'" :class="activeTab === 'running' ? 'bg-purple-600 text-white font-bold' : 'text-slate-400 hover:text-slate-200'" class="px-2.5 py-1 rounded-lg text-xs cursor-pointer">
          投入中 ({{ status?.recent_tasks?.length || 0 }})
        </button>
        <button @click="activeTab = 'all'" :class="activeTab === 'all' ? 'bg-indigo-600 text-white font-bold' : 'text-slate-400 hover:text-slate-200'" class="px-2.5 py-1 rounded-lg text-xs cursor-pointer">
          全マスター ({{ status?.total_jobs || 303 }})
        </button>
      </div>
      <span class="text-[11px] font-mono text-slate-400">
        スロット: <strong class="text-purple-400">{{ status?.occupied_slots || 0 }}</strong> / {{ status?.total_slots || 12 }}
      </span>
    </div>

    <ThunderCDPTaskList :status="status" :loading="loading" :tab="activeTab" />
    <ThunderConfirmModal :is-open="showConfirmModal" :job-count="status?.total_jobs || 303" @close="showConfirmModal = false" @confirm="handleStartConfirmed" />
  </div>
</template>
