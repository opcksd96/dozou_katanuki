<!-- frontend/src/components/admin/AdminModal.vue (100行以下 - SPEC-ADMINBOARD-001) -->
<script setup lang="ts">
import { ref, reactive, watch, onMounted, onUnmounted } from 'vue';
import { useAdmin } from '../../composables/useAdmin';
import AdminNavSidebar, { AdminTabId } from './AdminNavSidebar.vue';
import JobController from './JobController.vue';
import ConfigPortal from './ConfigPortal.vue';
import StashStatusView from './StashStatusView.vue';
import WhitelistView from './WhitelistView.vue';
import DatabaseView from './DatabaseView.vue';
import AuditReportView from './AuditReportView.vue';
import SkinFontEditor from './SkinFontEditor.vue';

const props = defineProps<{ isOpen: boolean }>();
const emit = defineEmits<{ (e: 'close'): void; (e: 'whitelistUpdated'): void }>();

const activeTab = ref<AdminTabId>('jobs');
const salvageForm = reactive({ platform: 'twitter', account: '', limit: 0 }), importForm = reactive({ warcPath: '', offline: true }), selectedPlatform = ref('twitter');
const fontPresets = {
  ja: [{ label: '標準ゴシック', value: 'Hiragino Sans, Meiryo, sans-serif' }],
  en: [{ label: '標準サンセリフ', value: 'Nunito, sans-serif' }],
  zh: [{ label: '標準簡体字', value: 'Microsoft YaHei, SimHei, sans-serif' }],
};

const admin = useAdmin();
const close = () => emit('close');
const handleKey = (e: KeyboardEvent) => { if (e.key === 'Escape' && props.isOpen) close(); };

const refreshCurrentTab = (tab: AdminTabId) => {
  try {
    if (tab === 'jobs') { admin.fetchActiveJob?.(); admin.fetchJobList?.(); }
    else if (tab === 'config') { admin.fetchConfig?.(); }
    else if (tab === 'skin') { admin.fetchSkinCSS?.(selectedPlatform.value || 'twitter'); }
    else if (tab === 'whitelist') { admin.fetchWhitelists?.(); }
    else if (tab === 'accounts') { admin.fetchAccounts?.(); }
    else if (tab === 'posts') { admin.searchArticles?.(); }
    else if (tab === 'media') { admin.fetchMedia?.(); }
  } catch (_) {}
};

watch(() => props.isOpen, (open) => {
  if (open) { refreshCurrentTab(activeTab.value); document.body.style.overflow = 'hidden'; }
  else { document.body.style.overflow = ''; }
}, { immediate: true });
watch(activeTab, (tab) => { if (props.isOpen) refreshCurrentTab(tab); });
onMounted(() => window.addEventListener('keydown', handleKey));
onUnmounted(() => { window.removeEventListener('keydown', handleKey); document.body.style.overflow = ''; });
</script>

<template>
  <div v-if="isOpen" class="fixed inset-0 z-50 w-screen h-screen bg-slate-950 flex flex-col overflow-hidden select-none overscroll-contain">
    <!-- モーダルヘッダー -->
    <header class="px-4 py-2 border-b border-slate-800 flex items-center justify-between bg-slate-900/90 shrink-0">
      <div class="flex items-center gap-2">
        <span class="text-base">⚙️</span>
        <h2 class="text-xs font-bold text-slate-200">土蔵・型抜き Admin Governance Portal</h2>
        <span class="text-[10px] px-1.5 py-0.5 rounded bg-blue-500/20 text-blue-400 font-mono font-bold">1-Generation Flatten Nav</span>
      </div>
      <button @click="close" class="px-2.5 py-1 rounded-lg bg-slate-800 hover:bg-slate-700 text-slate-300 hover:text-white flex items-center gap-1 text-xs font-bold transition-colors border border-slate-700 shadow cursor-pointer">
        <span>←</span> タイムラインに戻る (Esc)
      </button>
    </header>

    <!-- 2ペイン本体 (左: 1世代サイドバー, 右: コンテンツ領域) -->
    <div class="flex-1 min-h-0 flex overflow-hidden">
      <AdminNavSidebar :active-tab="activeTab" @select="(t) => activeTab = t" />

      <main class="flex-1 min-h-0 overflow-hidden p-3 bg-slate-950 flex flex-col">
        <JobController v-if="activeTab === 'jobs'" class="overflow-y-auto flex-1" :active-job="admin.activeJob?.value ?? admin.activeJob" :job-list="admin.jobList?.value ?? admin.jobList" :logs="admin.jobLogs?.value ?? admin.jobLogs" :salvage-form="salvageForm" :import-form="importForm" :action-loading="admin.isJobRunning?.value ?? admin.isJobRunning" :loading-jobs="false" @start-salvage="admin.startSalvage(salvageForm.platform, salvageForm.account, salvageForm.limit)" @start-import="admin.startManualImport(importForm.warcPath, importForm.offline)" @cancel-job="(id) => admin.cancelJob(id)" @fetch-jobs="admin.fetchJobList" @clear-logs="admin.clearLogs" />
        <AuditReportView v-else-if="activeTab === 'audit'" class="overflow-y-auto flex-1" :restoring="admin.isJobRunning?.value ?? admin.isJobRunning" @trigger-restore="(resetDB) => { admin.triggerRestore('', resetDB); activeTab = 'jobs'; }" />
        <StashStatusView v-else-if="activeTab === 'stash'" class="overflow-y-auto flex-1" :config="admin.configForm" />
        <DatabaseView v-else-if="activeTab === 'posts' || activeTab === 'media' || activeTab === 'accounts'" class="flex-1 min-h-0" :admin="admin" :view="activeTab" @navigate="(t) => activeTab = t" />
        <WhitelistView v-else-if="activeTab === 'whitelist'" class="overflow-y-auto flex-1" :whitelist-list="admin.whitelists?.value ?? admin.whitelists" :loading="admin.isWhitelistLoading?.value ?? admin.isWhitelistLoading" :status-message="null" @fetch="admin.fetchWhitelists" @add="async (t, v) => { await admin.addWhitelist(t, v); emit('whitelistUpdated'); }" @update="async (id, t, v, act) => { await admin.updateWhitelist(id, t, v, act); emit('whitelistUpdated'); }" @delete="async (id) => { await admin.deleteWhitelist(id); emit('whitelistUpdated'); }" @toggle="async (id) => { await admin.toggleWhitelist(id); emit('whitelistUpdated'); }" />
        <ConfigPortal v-else-if="activeTab === 'config'" class="overflow-y-auto flex-1" :config="admin.configForm" :loading-config="admin.isConfigLoading?.value ?? admin.isConfigLoading" :saving-config="admin.isConfigLoading?.value ?? admin.isConfigLoading" :save-status="admin.configSaved?.value ?? admin.configSaved ? { success: true, message: '設定を保存しました' } : null" @save-config="admin.saveConfig" @load-config="admin.fetchConfig" />
        <SkinFontEditor v-else-if="activeTab === 'skin'" class="overflow-y-auto flex-1" :skin-c-s-s="admin.skinCSS?.value ?? admin.skinCSS" :loading-skin="admin.isSkinLoading?.value ?? admin.isSkinLoading" :saving-skin="admin.isSkinLoading?.value ?? admin.isSkinLoading" :skin-status="admin.isSkinSaved?.value ?? admin.isSkinSaved ? { success: true, message: 'スキンを保存しました' } : null" :selected-platform="selectedPlatform" :font-presets="fontPresets" :config="admin.configForm" :saving-config="admin.isConfigLoading?.value ?? admin.isConfigLoading" :save-status="null" @fetch-skin="(p) => admin.fetchSkinCSS(p)" @save-skin="(p, css) => admin.saveSkinCSS(p, css)" @apply-dynamic-skin="() => {}" @save-config="admin.saveConfig" />
      </main>
    </div>
  </div>
</template>
