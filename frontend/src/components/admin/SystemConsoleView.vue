<!-- frontend/src/components/admin/SystemConsoleView.vue (100行以下 - SPEC-PRINCIPLE-001) -->
<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { Terminal, RefreshCw, Zap, Copy, Trash2 } from 'lucide-vue-next';
import { useToast } from '../../composables/useToast';

const getApp = () => (window as any)?.go?.app?.App || (window as any)?.go?.main?.App;
const { addToast } = useToast();
const journals = ref<any[]>([]), selectedEntry = ref<any>(null);
const filterComp = ref('all'), filterLevel = ref('all'), isRestarting = ref(false);

const fetchJournals = async () => {
  try {
    const app = getApp();
    if (app?.GetSystemJournals) journals.value = (await app.GetSystemJournals(200)) || [];
  } catch (_) {}
};

const restartBackend = async () => {
  if (!window.confirm('⚡ バックエンドの内部サービスを再初期化しますか？')) return;
  isRestarting.value = true;
  try {
    const ok = await getApp()?.RestartBackendServices?.();
    if (ok) { addToast('⚡ バックエンドサービスを再初期化しました！', 'success', 3000); await fetchJournals(); }
  } catch (e: any) { addToast(`再初期化失敗: ${e?.message || e}`, 'error', 4000); }
  finally { isRestarting.value = false; }
};

onMounted(fetchJournals);
</script>

<template>
  <div class="flex-1 flex flex-col h-full space-y-3 font-sans min-h-0 bg-slate-950 p-2 sm:p-3">
    <!-- ヘッダー / コントロールバー -->
    <div class="flex flex-wrap items-center justify-between gap-2 bg-slate-900/90 border border-slate-800 rounded-xl p-2.5 shrink-0 shadow-lg">
      <div class="flex items-center gap-2">
        <Terminal class="w-4 h-4 text-emerald-400" />
        <h2 class="text-xs font-bold text-slate-100">📜 システムコンソール＆ジャーナル</h2>
        <span class="text-[10px] font-mono text-slate-400 bg-slate-950 px-2 py-0.5 rounded border border-slate-800">{{ journals.length }} records</span>
      </div>
      <div class="flex items-center gap-1.5">
        <select v-model="filterComp" class="bg-slate-950 border border-slate-700 text-xs text-slate-200 rounded-lg px-2.5 py-1">
          <option value="all">全コンポーネント</option>
          <option value="stash">Stash</option>
          <option value="downloader">Downloader</option>
          <option value="crawler">Crawler</option>
          <option value="audit">Audit</option>
          <option value="system">System</option>
        </select>
        <select v-model="filterLevel" class="bg-slate-950 border border-slate-700 text-xs text-slate-200 rounded-lg px-2.5 py-1">
          <option value="all">全レベル</option>
          <option value="INFO">INFO</option>
          <option value="WARN">WARN</option>
          <option value="ERROR">ERROR</option>
        </select>
        <button @click="fetchJournals" class="p-1.5 bg-slate-800 hover:bg-slate-700 text-slate-300 rounded-lg text-xs cursor-pointer active:scale-95" title="更新"><RefreshCw class="w-3.5 h-3.5" /></button>
        <button @click="restartBackend" :disabled="isRestarting" class="px-2.5 py-1 bg-gradient-to-r from-amber-600 to-rose-600 hover:from-amber-500 text-white rounded-lg text-xs font-bold shadow-md cursor-pointer active:scale-95 flex items-center gap-1.5 disabled:opacity-50" title="バックエンドサービスを再初期化">
          <Zap class="w-3.5 h-3.5 text-yellow-200" /><span>⚡ 再初期化</span>
        </button>
      </div>
    </div>

    <!-- 巨大コンソール領域: 左右2ペイン (モバイルは上下) -->
    <div class="flex-1 grid grid-cols-1 md:grid-cols-2 gap-3 min-h-0">
      <!-- 左ペイン: ジャーナル一覧 -->
      <div class="bg-slate-900/90 border border-slate-800 rounded-2xl p-2.5 flex flex-col min-h-0 overflow-hidden shadow-md">
        <div class="flex-1 overflow-y-auto space-y-1.5 pr-1 font-mono text-xs">
          <div
            v-for="j in journals.filter(x => (filterComp === 'all' || x.component === filterComp) && (filterLevel === 'all' || x.level === filterLevel))"
            :key="j.id"
            @click="selectedEntry = j"
            :class="[
              'p-2 rounded-xl border transition-all cursor-pointer flex items-center justify-between text-[11px]',
              selectedEntry?.id === j.id ? 'bg-blue-950/70 border-blue-500 text-blue-200 shadow-md' : 'bg-slate-950/80 border-slate-800/80 text-slate-400 hover:text-slate-200 hover:bg-slate-900'
            ]"
          >
            <div class="flex items-center gap-2 truncate">
              <span class="text-[9px] px-1.5 py-0.5 rounded font-bold font-mono" :class="j.level === 'ERROR' ? 'bg-rose-950 text-rose-300' : j.level === 'WARN' ? 'bg-amber-950 text-amber-300' : 'bg-slate-800 text-slate-300'">{{ j.component }}</span>
              <span class="font-bold text-slate-200 truncate">{{ j.event }}</span>
              <span class="text-[10px] text-slate-500 truncate hidden sm:inline">{{ j.message }}</span>
            </div>
            <span class="text-[10px] text-slate-500 shrink-0 font-mono">{{ new Date(j.timestamp).toLocaleTimeString() }}</span>
          </div>
          <div v-if="journals.length === 0" class="text-center py-16 text-slate-600">ジャーナル記録はありません</div>
        </div>
      </div>

      <!-- 右ペイン: 巨大 JSON インスペクター -->
      <div class="bg-slate-900/90 border border-slate-800 rounded-2xl p-3 flex flex-col min-h-0 overflow-hidden shadow-md font-mono text-xs">
        <div v-if="selectedEntry" class="flex-1 flex flex-col min-h-0 space-y-2">
          <div class="flex items-center justify-between border-b border-slate-800 pb-2">
            <span class="font-bold text-blue-400 text-xs">Entry: {{ selectedEntry.id }} ({{ selectedEntry.event }})</span>
            <span class="text-[10px] text-slate-500">{{ new Date(selectedEntry.timestamp).toLocaleString() }}</span>
          </div>
          <div class="text-xs text-slate-300">{{ selectedEntry.message }}</div>
          <div class="flex-1 min-h-0 bg-slate-950 border border-slate-800/90 rounded-xl p-3 overflow-y-auto">
            <pre class="text-emerald-400 text-[11px] whitespace-pre-wrap select-text">{{ JSON.stringify(selectedEntry.payload || {}, null, 2) }}</pre>
          </div>
        </div>
        <div v-else class="flex-1 flex items-center justify-center text-slate-600 text-xs">左側のジャーナルをクリックして JSON ペイロードを展開</div>
      </div>
    </div>
  </div>
</template>
