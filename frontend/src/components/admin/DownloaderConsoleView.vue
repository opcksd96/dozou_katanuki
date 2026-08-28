<!-- frontend/src/components/admin/DownloaderConsoleView.vue (100行以下 - SPEC-PRINCIPLE-001) -->
<script setup lang="ts">
import { useDownloaderConsole } from '../../composables/admin/useDownloaderConsole';

defineProps<{ admin?: any }>();
const { status, loading, controlMotrix, launchThunder, fetchStatus } = useDownloaderConsole();

const formatSpeed = (b: number) => {
  if (!b || b === 0) return '0 B/s';
  return b < 1024 * 1024 ? `${(b / 1024).toFixed(1)} KB/s` : `${(b / (1024 * 1024)).toFixed(2)} MB/s`;
};
</script>

<template>
  <div class="h-full flex flex-col p-3 sm:p-4 space-y-4 bg-slate-950 text-slate-100 overflow-y-auto font-sans max-w-5xl mx-auto">
    <div class="flex items-center justify-between border-b border-slate-800 pb-2">
      <div class="flex items-center gap-2"><span class="text-xl">⚡</span><h2 class="text-sm font-bold">ダウンローダー遠隔管理 (Motrix & Thunder)</h2></div>
      <button @click="fetchStatus" class="px-2.5 py-1 bg-slate-800 hover:bg-slate-700 text-slate-300 rounded text-xs active:scale-95 cursor-pointer">🔄 更新</button>
    </div>

    <!-- 上段: 2大ダウンローダーサマリー -->
    <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
      <div class="p-3.5 bg-slate-900/90 border border-slate-800 rounded-2xl space-y-2.5 shadow-lg">
        <div class="flex items-center justify-between">
          <div class="flex items-center gap-2"><span class="text-base">🚀</span><span class="font-bold text-xs">Motrix Next (Aria2)</span></div>
          <span :class="status?.motrix?.is_online ? 'bg-emerald-950 text-emerald-300 border-emerald-700/60' : 'bg-rose-950 text-rose-300 border-rose-700/60'" class="px-2 py-0.5 rounded text-[10px] font-mono font-bold border">
            {{ status?.motrix?.is_online ? '🟢 ONLINE' : '🔴 OFFLINE' }}
          </span>
        </div>
        <div class="grid grid-cols-4 gap-1.5 text-center">
          <div class="p-1.5 bg-slate-950 rounded-xl border border-slate-800"><div class="text-[10px] text-slate-500">アクティブ</div><div class="text-xs font-mono font-bold text-blue-400">{{ status?.motrix?.num_active ?? 0 }}</div></div>
          <div class="p-1.5 bg-slate-950 rounded-xl border border-slate-800"><div class="text-[10px] text-slate-500">待機</div><div class="text-xs font-mono font-bold text-amber-400">{{ status?.motrix?.num_waiting ?? 0 }}</div></div>
          <div class="p-1.5 bg-slate-950 rounded-xl border border-slate-800"><div class="text-[10px] text-slate-500">停止</div><div class="text-xs font-mono font-bold text-emerald-400">{{ status?.motrix?.num_stopped ?? 0 }}</div></div>
          <div class="p-1.5 bg-slate-950 rounded-xl border border-slate-800"><div class="text-[10px] text-slate-500">速度</div><div class="text-xs font-mono font-bold text-purple-400">{{ formatSpeed(status?.motrix?.download_speed) }}</div></div>
        </div>
        <div class="flex flex-wrap gap-1 pt-1">
          <button @click="controlMotrix('safe_limits')" :disabled="loading" class="px-2 py-0.5 bg-indigo-950 hover:bg-indigo-900 border border-indigo-700/60 text-indigo-200 rounded text-[11px] font-semibold active:scale-95 cursor-pointer">🛡️ 安全2並列</button>
          <button @click="controlMotrix('pause_all')" :disabled="loading" class="px-2 py-0.5 bg-amber-950 hover:bg-amber-900 border border-amber-700/60 text-amber-200 rounded text-[11px] font-semibold active:scale-95 cursor-pointer">⏸️ 一時停止</button>
          <button @click="controlMotrix('unpause_all')" :disabled="loading" class="px-2 py-0.5 bg-blue-950 hover:bg-blue-900 border border-blue-700/60 text-blue-200 rounded text-[11px] font-semibold active:scale-95 cursor-pointer">▶️ 再開</button>
          <button @click="controlMotrix('purge_all')" :disabled="loading" class="px-2 py-0.5 bg-rose-950 hover:bg-rose-900 border border-rose-700/60 text-rose-200 rounded text-[11px] font-semibold active:scale-95 cursor-pointer">🧹 履歴削除</button>
        </div>
      </div>

      <div class="p-3.5 bg-slate-900/90 border border-slate-800 rounded-2xl space-y-2.5 shadow-lg flex flex-col justify-between">
        <div class="space-y-2">
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-2"><span class="text-base">⚡</span><span class="font-bold text-xs">Thunder (迅雷 P2SP)</span></div>
            <span :class="status?.thunder?.is_installed ? 'bg-amber-950 text-amber-300 border-amber-700/60' : 'bg-slate-800 text-slate-400 border-slate-700'" class="px-2 py-0.5 rounded text-[10px] font-mono font-bold border">
              {{ status?.thunder?.is_installed ? 'INSTALLED' : 'NOT FOUND' }}
            </span>
          </div>
          <div class="p-2 bg-slate-950 rounded-xl border border-slate-800 flex items-center justify-between">
            <div class="text-xs text-slate-400">保留（RETAINED）メディア</div>
            <div class="text-sm font-mono font-bold text-amber-400">{{ status?.thunder?.retained_count ?? 0 }} 件</div>
          </div>
        </div>
        <div class="flex gap-2 pt-1">
          <button @click="launchThunder" :disabled="!status?.thunder?.is_installed" class="flex-1 py-1 bg-gradient-to-r from-amber-600 to-orange-600 hover:from-amber-500 text-white rounded text-xs font-bold shadow active:scale-95 cursor-pointer disabled:opacity-50">🚀 Thunder 起動</button>
          <button @click="admin?.startThunderEscalate?.()" :disabled="!status?.thunder?.is_installed || (status?.thunder?.retained_count || 0) === 0" class="flex-1 py-1 bg-purple-600 hover:bg-purple-500 text-white rounded text-xs font-bold shadow active:scale-95 cursor-pointer disabled:opacity-50">⚡ RETAINED 一括投入</button>
        </div>
      </div>
    </div>

    <!-- 下段: Motrix アクティブタスク一覧 (TASK10本格実装に向けたフル領域) -->
    <div class="flex-1 min-h-[220px] bg-slate-900/90 border border-slate-800 rounded-2xl p-3 flex flex-col space-y-2 shadow-lg overflow-hidden">
      <span class="text-xs font-bold text-slate-300">アクティブダウンロード一覧 ({{ status?.motrix?.active_tasks?.length ?? 0 }})</span>
      <div class="flex-1 overflow-y-auto space-y-2 pr-1">
        <div v-if="!status?.motrix?.active_tasks || status?.motrix?.active_tasks.length === 0" class="text-center py-10 text-slate-500 text-xs font-mono">アクティブなダウンロードタスクはありません</div>
        <div v-for="t in status?.motrix?.active_tasks" :key="t.gid" class="p-2.5 bg-slate-950 rounded-xl border border-slate-800/80 space-y-1.5">
          <div class="flex items-center justify-between text-xs"><span class="font-mono font-bold text-slate-200 truncate max-w-md">{{ t.file_name || t.gid }}</span><span class="font-mono text-purple-400 text-[11px]">{{ formatSpeed(t.download_speed) }}</span></div>
          <div class="w-full bg-slate-900 rounded-full h-1.5 overflow-hidden"><div class="bg-gradient-to-r from-blue-500 to-purple-500 h-full rounded-full" :style="{ width: `${t.progress}%` }"></div></div>
          <div class="flex items-center justify-between text-[10px] text-slate-500 font-mono"><span>GID: {{ t.gid }}</span><span>{{ t.progress.toFixed(1) }}%</span></div>
        </div>
      </div>
    </div>
  </div>
</template>
