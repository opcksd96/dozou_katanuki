<!-- frontend/src/components/admin/PluginHubView.vue (100行以下 - SPEC-PLUGIN-001) -->
<script setup lang="ts">
import { ref } from 'vue';
import JobController from './JobController.vue';
import SkinFontEditor from './SkinFontEditor.vue';
import PluginSourcesView from './plugins/PluginSourcesView.vue';

export type PluginSubTab = 'jobs' | 'skin' | 'sources';

defineProps<{
  admin: any;
  salvageForm: { platform: string; account: string; source?: string; limit: number };
  importForm: { warcPath: string; offline: boolean };
  selectedPlatform: string;
}>();

defineEmits<{
  (e: 'update:selectedPlatform', platform: string): void;
  (e: 'startSalvage'): void;
  (e: 'startImport'): void;
}>();

const currentSubTab = ref<PluginSubTab>('jobs');
const subTabs: { id: PluginSubTab; icon: string; label: string }[] = [
  { id: 'jobs', icon: '🚀', label: 'サルベージ＆ジョブ' },
  { id: 'skin', icon: '🎨', label: 'スキン＆外観' },
  { id: 'sources', icon: '📡', label: 'ソースエンジン' },
];
</script>

<template>
  <div class="flex flex-col h-full space-y-4">
    <!-- プラグイン選択 ＆ サブタブヘッダー -->
    <div class="flex flex-wrap items-center justify-between gap-3 bg-slate-900/90 border border-slate-800 rounded-xl p-3 shrink-0">
      <div class="flex items-center gap-3">
        <div class="flex items-center gap-2">
          <span class="text-base">🧩</span>
          <span class="text-xs font-bold text-slate-300">対象プラグイン:</span>
        </div>
        <select
          :value="selectedPlatform"
          @change="(e) => $emit('update:selectedPlatform', (e.target as HTMLSelectElement).value)"
          class="bg-slate-950 border border-slate-700 text-xs font-bold text-slate-100 rounded-lg px-3 py-1.5 focus:outline-none focus:border-blue-500"
        >
          <option value="twitter">🐦 Twitter / X プラグイン (Active 🟢)</option>
          <option value="bsky">🦋 Bluesky プラグイン (Ready 🟡)</option>
        </select>
      </div>

      <!-- サブタブセレクタ -->
      <div class="flex items-center bg-slate-950 p-1 rounded-lg border border-slate-800 gap-1">
        <button
          v-for="st in subTabs"
          :key="st.id"
          @click="currentSubTab = st.id"
          :class="[
            'flex items-center gap-1.5 px-3 py-1 rounded-md text-xs font-semibold transition-all cursor-pointer',
            currentSubTab === st.id ? 'bg-blue-600 text-white shadow-sm' : 'text-slate-400 hover:text-slate-200 hover:bg-slate-900'
          ]"
        >
          <span>{{ st.icon }}</span><span>{{ st.label }}</span>
        </button>
      </div>
    </div>

    <!-- コンテンツ領域 -->
    <div class="flex-1 overflow-y-auto min-h-0 pr-1">
      <JobController
        v-if="currentSubTab === 'jobs'"
        :active-job="admin.activeJob?.value ?? admin.activeJob"
        :job-list="admin.jobList?.value ?? admin.jobList"
        :logs="admin.jobLogs?.value ?? admin.jobLogs"
        :salvage-form="salvageForm"
        :import-form="importForm"
        :action-loading="admin.isJobRunning?.value ?? admin.isJobRunning"
        :loading-jobs="false"
        @start-salvage="$emit('startSalvage')"
        @start-import="$emit('startImport')"
        @cancel-job="(id) => admin.cancelJob(id)"
        @fetch-jobs="admin.fetchJobList"
        @clear-logs="admin.clearLogs"
      />
      <SkinFontEditor
        v-else-if="currentSubTab === 'skin'"
        :skin-c-s-s="admin.skinCSS?.value ?? admin.skinCSS"
        :loading-skin="admin.isSkinLoading?.value ?? admin.isSkinLoading"
        :saving-skin="admin.isSkinLoading?.value ?? admin.isSkinLoading"
        :skin-status="admin.isSkinSaved?.value ?? admin.isSkinSaved ? { success: true, message: 'スキンを保存しました' } : null"
        :selected-platform="selectedPlatform"
        @fetch-skin="(p) => admin.fetchSkinCSS(p)"
        @save-skin="(p, css) => admin.saveSkinCSS(p, css)"
        @apply-dynamic-skin="() => {}"
      />
      <PluginSourcesView v-else-if="currentSubTab === 'sources'" :platform="selectedPlatform" />
    </div>
  </div>
</template>
