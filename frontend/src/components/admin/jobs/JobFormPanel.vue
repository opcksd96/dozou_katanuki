<!-- frontend/src/components/admin/jobs/JobFormPanel.vue (100行以下) -->
<script setup lang="ts">
import { ref } from 'vue';
const props = withDefaults(
  defineProps<{
    salvageForm: { platform: string; account: string; source?: string; limit: number };
    importForm: { warcPath: string; offline: boolean };
    actionLoading?: boolean;
    isJobRunning?: boolean;
  }>(),
  { actionLoading: false, isJobRunning: false }
);
const emit = defineEmits<{ (e: 'startSalvage'): void; (e: 'startImport'): void }>();
const showManualImport = ref(false);
</script>

<template>
  <div class="bg-slate-900/80 border border-slate-800 rounded-xl p-4 space-y-4">
    <div class="flex items-center justify-between">
      <div class="flex items-center gap-2">
        <h3 class="text-sm font-bold text-slate-200 flex items-center gap-2"><span>🚀</span> ジョブ実行</h3>
        <span class="text-[10px] px-2 py-0.5 rounded-full bg-blue-500/20 text-blue-400 font-mono font-semibold border border-blue-500/30">
          Target: {{ salvageForm.platform.toUpperCase() }}
        </span>
      </div>
      <button @click="showManualImport = !showManualImport" class="text-xs text-slate-400 hover:text-slate-200 underline">
        {{ showManualImport ? '通常サルベージへ' : 'WARC手動インポート' }}
      </button>
    </div>
    <form v-if="!showManualImport" @submit.prevent="emit('startSalvage')" class="grid grid-cols-1 sm:grid-cols-4 gap-3">
      <div>
        <label class="block text-xs text-slate-400 mb-1">サルベージソース</label>
        <select v-model="salvageForm.source" class="w-full bg-slate-950 border border-slate-700 rounded-lg px-3 py-1.5 text-xs text-slate-200">
          <option value="all">⚡ 全ソース自動 (推奨)</option>
          <option value="sotwe">🧬 Sotwe (スレッド/メディア)</option>
          <option value="twistalker">📜 Twistalker</option>
          <option value="nitter">🌐 Nitter</option>
          <option value="wayback">🏛️ Wayback 魚拓</option>
          <option value="official">🐦 X Syndication</option>
        </select>
      </div>
      <div class="sm:col-span-2">
        <label class="block text-xs text-slate-400 mb-1">アカウント名 / ID</label>
        <input v-model="salvageForm.account" type="text" placeholder="target_user" required class="w-full bg-slate-950 border border-slate-700 rounded-lg px-3 py-1.5 text-xs text-slate-200 font-mono" />
      </div>
      <div>
        <label class="block text-xs text-slate-400 mb-1">上限件数 (0 = 全件)</label>
        <div class="flex gap-2">
          <input v-model.number="salvageForm.limit" type="number" min="0" placeholder="0" class="w-full bg-slate-950 border border-slate-700 rounded-lg px-3 py-1.5 text-xs text-slate-200 font-mono" />
          <button type="submit" :disabled="actionLoading || isJobRunning" class="bg-blue-600 hover:bg-blue-500 disabled:bg-slate-700 text-white text-xs font-semibold px-4 py-1.5 rounded-lg whitespace-nowrap">
            {{ actionLoading ? '⏳' : '実行' }}
          </button>
        </div>
      </div>
    </form>
    <form v-else @submit.prevent="emit('startImport')" class="grid grid-cols-1 sm:grid-cols-4 gap-3">
      <div class="sm:col-span-3">
        <label class="block text-xs text-slate-400 mb-1">WARC パス</label>
        <input v-model="importForm.warcPath" type="text" placeholder="D:/archives/sample.warc.gz" required class="w-full bg-slate-950 border border-slate-700 rounded-lg px-3 py-1.5 text-xs text-slate-200 font-mono" />
      </div>
      <div class="flex items-end gap-2">
        <label class="flex items-center gap-1.5 text-xs text-slate-300 mb-2 cursor-pointer">
          <input type="checkbox" v-model="importForm.offline" class="rounded bg-slate-950 border-slate-700" />
          <span>オフライン</span>
        </label>
        <button type="submit" :disabled="actionLoading || isJobRunning" class="bg-emerald-600 hover:bg-emerald-500 disabled:bg-slate-700 text-white text-xs font-semibold px-4 py-1.5 rounded-lg whitespace-nowrap mb-0.5">
          インポート
        </button>
      </div>
    </form>
  </div>
</template>
