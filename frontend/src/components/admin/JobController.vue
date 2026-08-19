<script setup lang="ts">
import { ref, watch, nextTick } from 'vue';
import { models } from '../../../wailsjs/go/models';

const props = defineProps<{
  activeJob: models.JobProgress | null;
  jobList: models.JobProgress[];
  logs: string[];
  salvageForm: { platform: string; account: string; limit: number };
  importForm: { warcPath: string; offline: boolean };
  actionLoading: boolean;
  loadingJobs: boolean;
}>();

const emit = defineEmits<{
  (e: 'startSalvage'): void;
  (e: 'startImport'): void;
  (e: 'cancelJob', jobId: string): void;
  (e: 'fetchJobs'): void;
  (e: 'clearLogs'): void;
}>();

const showManualImport = ref(false);
const autoScroll = ref(true);
const terminalRef = ref<HTMLElement | null>(null);

// ログ追加時に自動スクロール
watch(
  () => props.logs.length,
  () => {
    if (autoScroll.value && terminalRef.value) {
      nextTick(() => {
        if (terminalRef.value) {
          terminalRef.value.scrollTop = terminalRef.value.scrollHeight;
        }
      });
    }
  }
);

const copyLogs = () => {
  const text = props.logs.join('\n');
  navigator.clipboard.writeText(text).then(() => {
    alert('ログをクリップボードにコピーしました');
  });
};

const formatDate = (dateStr?: any) => {
  if (!dateStr) return '-';
  try {
    const d = new Date(dateStr);
    if (isNaN(d.getTime())) return String(dateStr);
    return d.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit', second: '2-digit' });
  } catch {
    return String(dateStr);
  }
};
</script>

<template>
  <div class="space-y-6">
    <!-- 1. ジョブ実行コントロールパネル -->
    <div class="bg-slate-900/80 border border-slate-800 rounded-xl p-4 space-y-4">
      <div class="flex items-center justify-between">
        <h3 class="text-sm font-bold text-slate-200 flex items-center gap-2">
          <span>🚀</span> サルベージジョブ キック (SPEC-ADMINBOARD-001)
        </h3>
        <button
          @click="showManualImport = !showManualImport"
          class="text-xs text-slate-400 hover:text-slate-200 transition-colors underline"
        >
          {{ showManualImport ? '通常サルベージへ' : 'WARC手動インポート' }}
        </button>
      </div>

      <!-- 通常サルベージフォーム -->
      <form v-if="!showManualImport" @submit.prevent="emit('startSalvage')" class="grid grid-cols-1 sm:grid-cols-4 gap-3">
        <div>
          <label class="block text-xs text-slate-400 mb-1">プラットフォーム</label>
          <select
            v-model="salvageForm.platform"
            class="w-full bg-slate-950 border border-slate-700 rounded-lg px-3 py-1.5 text-xs text-slate-200 focus:outline-none focus:border-blue-500"
          >
            <option value="twitter">Twitter / X</option>
            <option value="bsky">Bluesky</option>
          </select>
        </div>
        <div class="sm:col-span-2">
          <label class="block text-xs text-slate-400 mb-1">アカウント名 / ID</label>
          <input
            v-model="salvageForm.account"
            type="text"
            placeholder="e.g. target_user"
            required
            class="w-full bg-slate-950 border border-slate-700 rounded-lg px-3 py-1.5 text-xs text-slate-200 focus:outline-none focus:border-blue-500 font-mono"
          />
        </div>
        <div>
          <label class="block text-xs text-slate-400 mb-1">取得上限件数</label>
          <div class="flex gap-2">
            <input
              v-model.number="salvageForm.limit"
              type="number"
              min="1"
              max="5000"
              class="w-full bg-slate-950 border border-slate-700 rounded-lg px-3 py-1.5 text-xs text-slate-200 focus:outline-none focus:border-blue-500 font-mono"
            />
            <button
              type="submit"
              :disabled="actionLoading || (activeJob && activeJob.status === 'running')"
              class="bg-blue-600 hover:bg-blue-500 disabled:bg-slate-700 disabled:cursor-not-allowed text-white text-xs font-semibold px-4 py-1.5 rounded-lg transition-colors whitespace-nowrap flex items-center gap-1 shadow-sm"
            >
              <span v-if="actionLoading" class="animate-spin">⏳</span>
              <span v-else>実行</span>
            </button>
          </div>
        </div>
      </form>

      <!-- WARC 手動インポートフォーム -->
      <form v-else @submit.prevent="emit('startImport')" class="grid grid-cols-1 sm:grid-cols-4 gap-3">
        <div class="sm:col-span-3">
          <label class="block text-xs text-slate-400 mb-1">WARC ファイル絶対パス</label>
          <input
            v-model="importForm.warcPath"
            type="text"
            placeholder="D:/data/archives/sample.warc.gz"
            required
            class="w-full bg-slate-950 border border-slate-700 rounded-lg px-3 py-1.5 text-xs text-slate-200 focus:outline-none focus:border-blue-500 font-mono"
          />
        </div>
        <div class="flex items-end gap-2">
          <label class="flex items-center gap-1.5 text-xs text-slate-300 mb-2 cursor-pointer select-none">
            <input type="checkbox" v-model="importForm.offline" class="rounded bg-slate-950 border-slate-700 text-blue-600" />
            <span>オフライン</span>
          </label>
          <button
            type="submit"
            :disabled="actionLoading || (activeJob && activeJob.status === 'running')"
            class="bg-emerald-600 hover:bg-emerald-500 disabled:bg-slate-700 disabled:cursor-not-allowed text-white text-xs font-semibold px-4 py-1.5 rounded-lg transition-colors whitespace-nowrap mb-0.5"
          >
            インポート
          </button>
        </div>
      </form>
    </div>

    <!-- 2. アクティブジョブ進捗バー ＆ ステータス -->
    <div v-if="activeJob" class="bg-slate-900/90 border border-blue-500/40 rounded-xl p-4 space-y-3">
      <div class="flex items-center justify-between text-xs">
        <div class="flex items-center gap-2">
          <span class="inline-flex items-center px-2 py-0.5 rounded text-[10px] font-semibold font-mono"
            :class="{
              'bg-blue-500/20 text-blue-400 border border-blue-500/30 animate-pulse': activeJob.status === 'running',
              'bg-amber-500/20 text-amber-400 border border-amber-500/30': activeJob.status === 'pending',
              'bg-emerald-500/20 text-emerald-400 border border-emerald-500/30': activeJob.status === 'completed',
              'bg-rose-500/20 text-rose-400 border border-rose-500/30': activeJob.status === 'failed',
              'bg-slate-700 text-slate-300': activeJob.status === 'cancelled',
            }"
          >
            ● {{ activeJob.status.toUpperCase() }}
          </span>
          <span class="font-mono text-slate-300 font-bold">[{{ activeJob.type }}] {{ activeJob.id }}</span>
        </div>
        <div class="flex items-center gap-3">
          <span class="text-slate-400 font-mono text-[11px]">
            {{ activeJob.current }} / {{ activeJob.total > 0 ? activeJob.total : '?' }} 件 ({{ (activeJob.percentage || 0).toFixed(1) }}%)
          </span>
          <button
            v-if="activeJob.status === 'running' || activeJob.status === 'pending'"
            @click="emit('cancelJob', activeJob.id)"
            class="px-2.5 py-1 bg-rose-900/40 hover:bg-rose-800 text-rose-300 border border-rose-700/50 rounded text-xs transition-colors font-semibold"
          >
            ■ 中断 (Cancel)
          </button>
        </div>
      </div>

      <!-- プログレスバー -->
      <div class="w-full bg-slate-950 h-3 rounded-full overflow-hidden p-0.5 border border-slate-800">
        <div
          class="h-full bg-gradient-to-r from-blue-600 via-indigo-500 to-emerald-400 rounded-full transition-all duration-300 shadow-sm"
          :style="{ width: `${Math.min(100, Math.max(0, activeJob.percentage || 0))}%` }"
        ></div>
      </div>

      <div class="text-xs text-slate-400 truncate font-mono">
        {{ activeJob.message || '処理を実行中...' }}
      </div>
    </div>

    <!-- 3. Scraper View (StdoutPipe リアルタイム疑似端末) -->
    <div class="bg-black/90 border border-slate-800 rounded-xl overflow-hidden shadow-inner">
      <div class="bg-slate-900 px-4 py-2 border-b border-slate-800 flex items-center justify-between text-xs">
        <div class="flex items-center gap-2">
          <span class="w-2.5 h-2.5 rounded-full bg-red-500/80"></span>
          <span class="w-2.5 h-2.5 rounded-full bg-yellow-500/80"></span>
          <span class="w-2.5 h-2.5 rounded-full bg-green-500/80"></span>
          <span class="ml-2 font-mono text-slate-400 font-semibold">Scraper View Terminal (StdoutPipe)</span>
        </div>
        <div class="flex items-center gap-3 text-slate-400">
          <label class="flex items-center gap-1 cursor-pointer select-none text-[11px]">
            <input type="checkbox" v-model="autoScroll" class="rounded bg-slate-950 border-slate-700 text-blue-600" />
            <span>自動追従</span>
          </label>
          <button @click="copyLogs" class="hover:text-slate-200 transition-colors" title="ログをコピー">
            📋 コピー
          </button>
          <button @click="emit('clearLogs')" class="hover:text-slate-200 transition-colors" title="ログ消去">
            🧹 クリア
          </button>
        </div>
      </div>

      <div
        ref="terminalRef"
        class="p-3 h-52 overflow-y-auto font-mono text-xs text-emerald-400 space-y-0.5 select-text leading-relaxed bg-black/60"
      >
        <div v-if="logs.length === 0" class="text-slate-600 italic">
          --- 待機中 / ログはありません ---
        </div>
        <div v-for="(line, idx) in logs" :key="idx" class="whitespace-pre-wrap break-all hover:bg-white/5 px-1 rounded">
          <span class="text-slate-600 select-none mr-2 text-[10px]">{{ idx + 1 }}</span>
          <span>{{ line }}</span>
        </div>
      </div>
    </div>

    <!-- 4. ジョブ履歴一覧 -->
    <div class="bg-slate-900/60 border border-slate-800 rounded-xl p-4 space-y-3">
      <div class="flex items-center justify-between">
        <h3 class="text-sm font-bold text-slate-200 flex items-center gap-2">
          <span>📜</span> ジョブ実行履歴
        </h3>
        <button
          @click="emit('fetchJobs')"
          :disabled="loadingJobs"
          class="text-xs text-slate-400 hover:text-blue-400 transition-colors flex items-center gap-1"
        >
          <span :class="{ 'animate-spin': loadingJobs }">🔄</span> 更新
        </button>
      </div>

      <div class="overflow-x-auto max-h-48 overflow-y-auto">
        <table class="w-full text-left text-xs font-mono text-slate-300">
          <thead class="bg-slate-950/80 text-slate-400 sticky top-0 border-b border-slate-800 text-[11px]">
            <tr>
              <th class="py-2 px-3">Status</th>
              <th class="py-2 px-3">Job ID</th>
              <th class="py-2 px-3">Type</th>
              <th class="py-2 px-3">Progress</th>
              <th class="py-2 px-3">Started</th>
              <th class="py-2 px-3">Action</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-800/60">
            <tr v-if="jobList.length === 0">
              <td colspan="6" class="py-4 text-center text-slate-500 italic">
                実行履歴はありません
              </td>
            </tr>
            <tr v-for="job in jobList" :key="job.id" class="hover:bg-slate-800/40 transition-colors">
              <td class="py-2 px-3">
                <span
                  class="inline-block px-1.5 py-0.5 rounded text-[10px] font-semibold"
                  :class="{
                    'bg-blue-500/20 text-blue-400': job.status === 'running',
                    'bg-emerald-500/20 text-emerald-400': job.status === 'completed',
                    'bg-rose-500/20 text-rose-400': job.status === 'failed',
                    'bg-slate-700 text-slate-300': job.status === 'cancelled',
                    'bg-amber-500/20 text-amber-400': job.status === 'pending',
                  }"
                >
                  {{ job.status }}
                </span>
              </td>
              <td class="py-2 px-3 text-slate-200 font-semibold truncate max-w-[140px]" :title="job.id">
                {{ job.id }}
              </td>
              <td class="py-2 px-3 text-slate-400">{{ job.type }}</td>
              <td class="py-2 px-3 text-slate-300">{{ job.current }} / {{ job.total }}</td>
              <td class="py-2 px-3 text-slate-400 text-[11px]">{{ formatDate(job.started_at) }}</td>
              <td class="py-2 px-3">
                <button
                  v-if="job.status === 'running' || job.status === 'pending'"
                  @click="emit('cancelJob', job.id)"
                  class="text-rose-400 hover:text-rose-300 underline text-[11px]"
                >
                  中止
                </button>
                <span v-else class="text-slate-600">-</span>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>
