<!-- frontend/src/components/admin/AdminModal.vue (100行以下 - SPEC-ADMINBOARD-001) -->
<script setup lang="ts">
import { ref, reactive, watch } from 'vue';
import { useAdmin } from '../../composables/useAdmin';
import JobController from './JobController.vue';
import ConfigPortal from './ConfigPortal.vue';
import StashStatusView from './StashStatusView.vue';
import WhitelistView from './WhitelistView.vue';
import DatabaseView from './DatabaseView.vue';
import AuditReportView from './AuditReportView.vue';
import SkinFontEditor from './SkinFontEditor.vue';

const props = defineProps<{ isOpen: boolean }>();
const emit = defineEmits<{
  (e: 'close'): void; (e: 'whitelistUpdated'): void; (e: 'whitelist-updated'): void;
}>();

const activeTab = ref<'jobs' | 'config' | 'skin' | 'stash' | 'whitelist' | 'db' | 'audit'>('jobs');
const salvageForm = reactive({ platform: 'twitter', account: '', limit: 50 });
const importForm = reactive({ warcPath: '', offline: true });
const selectedPlatform = ref('twitter');
const fontPresets = {
  ja: [{ label: '標準ゴシック', value: 'Hiragino Sans, Meiryo, sans-serif' }],
  en: [{ label: '標準サンセリフ', value: 'Nunito, sans-serif' }],
  zh: [{ label: '標準簡体字', value: 'Microsoft YaHei, SimHei, sans-serif' }],
};

const admin = useAdmin();
watch(() => props.isOpen, (open) => {
  if (open) {
    admin.fetchConfig(); admin.fetchBroadcastStatus();
    document.body.style.overflow = 'hidden';
  } else {
    document.body.style.overflow = '';
  }
});
const onWhitelistChanged = () => { emit('whitelistUpdated'); };
</script>

<template>
  <div v-show="isOpen" class="fixed inset-0 z-50 w-screen h-screen bg-slate-950 flex flex-col overflow-hidden select-none overscroll-contain">
    <!-- モーダルヘッダー -->
    <div class="px-6 py-3 border-b border-slate-800 flex items-center justify-between bg-slate-900/90 shrink-0">
      <div class="flex items-center gap-3">
        <span class="text-xl">⚙️</span>
        <div>
          <h2 class="text-sm font-bold text-slate-100">土蔵・型抜き Admin Governance Portal</h2>
          <p class="text-[11px] text-slate-400">システム設定・ジョブ監視・データベース・翻訳・監査管理</p>
        </div>
      </div>
      <button @click="$emit('close')" class="px-3 py-1.5 rounded-lg bg-slate-800 hover:bg-slate-700 text-slate-300 hover:text-white flex items-center gap-1.5 text-xs font-bold transition-colors border border-slate-700 shadow">
        <span>←</span> タイムラインに戻る
      </button>
    </div>

    <!-- タブバー -->
    <div class="px-6 bg-slate-900/40 border-b border-slate-800 flex gap-2 overflow-x-auto text-xs font-semibold shrink-0">
      <button v-for="t in [
        { id: 'jobs', label: '🚀 ジョブ制御' }, { id: 'config', label: '🌐 システム設定' },
        { id: 'skin', label: '🎨 スキン・フォント' }, { id: 'whitelist', label: '📋 Whitelist' },
        { id: 'db', label: '🗄️ DB閲覧・翻訳' }, { id: 'audit', label: '🩺 整合性監査' },
        { id: 'stash', label: '🎛️ Stash状態' }
      ]" :key="t.id" @click="activeTab = t.id as any" class="py-3 px-3 border-b-2 transition-colors whitespace-nowrap" :class="activeTab === t.id ? 'border-blue-500 text-blue-400 bg-blue-500/10' : 'border-transparent text-slate-400 hover:text-slate-200'">
        {{ t.label }}
      </button>
    </div>

    <!-- コンテンツ領域（画面いっぱいに広がるメインワークスペース） -->
    <div class="flex-1 overflow-y-auto p-6 bg-slate-950">
      <JobController v-if="activeTab === 'jobs'" :active-job="admin.activeJob.value" :job-list="admin.jobList.value" :logs="admin.jobLogs.value" :salvage-form="salvageForm" :import-form="importForm" :action-loading="admin.isJobRunning.value" :loading-jobs="false" @start-salvage="admin.startSalvage(salvageForm.platform, salvageForm.account, salvageForm.limit)" @start-import="admin.startManualImport(importForm.warcPath, importForm.offline)" @cancel-job="(id) => admin.cancelJob(id)" @fetch-jobs="admin.fetchJobList" @clear-logs="admin.clearLogs" />
      <ConfigPortal v-else-if="activeTab === 'config'" :config="admin.configForm" :loading-config="admin.isConfigLoading.value" :saving-config="admin.isConfigLoading.value" :save-status="admin.configSaved.value ? { success: true, message: '設定を保存しました' } : null" @save-config="admin.saveConfig" @load-config="admin.fetchConfig" />
      <SkinFontEditor v-else-if="activeTab === 'skin'" :skin-c-s-s="admin.skinCSS.value" :loading-skin="admin.isSkinLoading.value" :saving-skin="admin.isSkinLoading.value" :skin-status="admin.isSkinSaved.value ? { success: true, message: 'スキンを保存しました' } : null" :selected-platform="selectedPlatform" :font-presets="fontPresets" :config="admin.configForm" :saving-config="admin.isConfigLoading.value" :save-status="null" @fetch-skin="(p) => admin.fetchSkinCSS(p)" @save-skin="(p, css) => admin.saveSkinCSS(p, css)" @apply-dynamic-skin="() => {}" @save-config="admin.saveConfig" />
      <WhitelistView v-else-if="activeTab === 'whitelist'" :whitelist-list="admin.whitelists.value" :loading="admin.isWhitelistLoading.value" :status-message="null" @fetch="admin.fetchWhitelists" @add="async (t, v) => { await admin.addWhitelist(t, v); onWhitelistChanged(); }" @update="async (id, t, v, act) => { await admin.updateWhitelist(id, t, v, act); onWhitelistChanged(); }" @delete="async (id) => { await admin.deleteWhitelist(id); onWhitelistChanged(); }" @toggle="async (id) => { await admin.toggleWhitelist(id); onWhitelistChanged(); }" />
      <DatabaseView v-else-if="activeTab === 'db'" :admin="admin" />
      <AuditReportView v-else-if="activeTab === 'audit'" :report="admin.auditReport.value" :loading="admin.isAuditing.value" :purging-files="admin.isPurgingFiles.value" :purging-db="admin.isPurgingDB.value" :status-message="admin.auditStatusMessage.value" :restoring="admin.isJobRunning.value" @run-audit="(pf, pdb) => admin.runAudit(pf, pdb)" @purge-orphan-files="(paths) => admin.purgeOrphanFiles(paths)" @purge-orphan-db-media="(ids) => admin.purgeOrphanDBMedia(ids)" @purgeOrphanDBMedia="(ids) => admin.purgeOrphanDBMedia(ids)" @trigger-restore="(resetDB) => { admin.triggerRestore('', resetDB); activeTab = 'jobs'; }" />
      <StashStatusView v-else-if="activeTab === 'stash'" :config="admin.configForm" />
    </div>
  </div>
</template>
