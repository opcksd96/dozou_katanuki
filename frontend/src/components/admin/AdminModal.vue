<!-- frontend/src/components/admin/AdminModal.vue (100行以下 - SPEC-ADMINBOARD-001) -->
<script setup lang="ts">
import { ref, reactive, watch, onMounted, onUnmounted } from 'vue';
import { useAdmin } from '../../composables/useAdmin';
import { WindowMinimise, WindowToggleMaximise, Quit } from '../../../wailsjs/runtime/runtime';
import { Minus, Square, X, ArrowLeft, Settings2 } from 'lucide-vue-next';
import AdminNavSidebar, { AdminTabId } from './AdminNavSidebar.vue';
import PluginHubView from './PluginHubView.vue';
import ConfigPortal from './ConfigPortal.vue';
import StashStatusView from './StashStatusView.vue';
import WhitelistView from './WhitelistView.vue';
import DatabaseView from './DatabaseView.vue';
import AuditReportView from './AuditReportView.vue';

const props = defineProps<{ isOpen: boolean }>();
const emit = defineEmits<{ (e: 'close'): void; (e: 'whitelistUpdated'): void }>();

const activeTab = ref<AdminTabId>('plugins');
const salvageForm = reactive({ platform: 'twitter', account: '', source: 'all', limit: 0 }), importForm = reactive({ warcPath: '', offline: true }), selectedPlatform = ref('twitter');
const admin = useAdmin();
const close = () => emit('close');
const handleKey = (e: KeyboardEvent) => { if (e.key === 'Escape' && props.isOpen) close(); };

const refreshCurrentTab = (tab: AdminTabId) => {
  try {
    if (tab === 'plugins') { admin.fetchActiveJob?.(); admin.fetchJobList?.(); admin.fetchSkinCSS?.(selectedPlatform.value || 'twitter'); }
    else if (tab === 'config') { admin.fetchConfig?.(); }
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
watch(selectedPlatform, (p) => { salvageForm.platform = p; if (activeTab.value === 'plugins') refreshCurrentTab('plugins'); });
onMounted(() => window.addEventListener('keydown', handleKey));
onUnmounted(() => { window.removeEventListener('keydown', handleKey); document.body.style.overflow = ''; });
</script>

<template>
  <div v-if="isOpen" class="fixed inset-0 z-50 w-screen h-screen bg-slate-950 flex flex-col overflow-hidden select-none overscroll-contain">
    <!-- モーダルヘッダー -->
    <header class="px-3 py-2 border-b border-white/10 flex items-center justify-between bg-slate-900/90 backdrop-blur-xl shrink-0 wails-drag select-none">
      <div class="flex items-center gap-2 wails-no-drag">
        <div class="w-6 h-6 rounded-lg bg-blue-600/20 border border-blue-500/30 flex items-center justify-center text-blue-400"><Settings2 class="w-3.5 h-3.5" /></div>
        <h2 class="text-xs font-bold text-slate-100">土蔵・型抜き Admin Governance Portal</h2>
        <span class="text-[9px] px-1.5 py-0.5 rounded-full bg-blue-500/20 text-blue-400 font-mono font-semibold border border-blue-500/30">Admin Board</span>
      </div>

      <div class="flex items-center gap-2 wails-no-drag">
        <button @click="close" class="px-2.5 py-1.5 rounded-lg bg-slate-800/90 hover:bg-slate-700 text-slate-300 hover:text-white flex items-center gap-1.5 text-xs font-semibold transition-all border border-slate-700/80 shadow-sm cursor-pointer" title="タイムラインに戻る (Esc)">
          <ArrowLeft class="w-3.5 h-3.5" /><span>タイムラインに戻る</span><span class="text-[10px] opacity-60 font-mono">(Esc)</span>
        </button>
        <div class="flex items-center ml-2 border-l border-slate-700/80 pl-2">
          <button @click="() => { try { WindowMinimise(); } catch {} }" class="p-1.5 rounded-md hover:bg-slate-800 text-slate-400 hover:text-slate-200 transition-colors" title="最小化"><Minus class="w-3.5 h-3.5" /></button>
          <button @click="() => { try { WindowToggleMaximise(); } catch {} }" class="p-1.5 rounded-md hover:bg-slate-800 text-slate-400 hover:text-slate-200 transition-colors" title="最大化 / 復元"><Square class="w-3.5 h-3.5" /></button>
          <button @click="() => { try { Quit(); } catch {} }" class="p-1.5 rounded-md hover:bg-red-600/80 text-slate-400 hover:text-white transition-colors" title="閉じる"><X class="w-3.5 h-3.5" /></button>
        </div>
      </div>
    </header>

    <!-- 2ペイン本体 (左: サイドバー, 右: コンテンツ) -->
    <div class="flex-1 min-h-0 flex overflow-hidden">
      <AdminNavSidebar :active-tab="activeTab" @select="(t) => activeTab = t" />
      <main class="flex-1 min-h-0 overflow-hidden p-3 bg-slate-950 flex flex-col">
        <PluginHubView v-if="activeTab === 'plugins'" :admin="admin" :salvage-form="salvageForm" :import-form="importForm" v-model:selected-platform="selectedPlatform" @start-salvage="admin.startSalvage(salvageForm.platform, salvageForm.account, salvageForm.limit, salvageForm.source)" @start-import="admin.startManualImport(importForm.warcPath, importForm.offline)" />
        <AuditReportView v-else-if="activeTab === 'audit'" class="overflow-y-auto flex-1" :restoring="admin.isJobRunning?.value ?? admin.isJobRunning" @trigger-restore="(resetDB) => { admin.triggerRestore('', resetDB); activeTab = 'plugins'; }" />
        <StashStatusView v-else-if="activeTab === 'stash'" class="overflow-y-auto flex-1" :config="admin.configForm" />
        <DatabaseView v-else-if="activeTab === 'posts' || activeTab === 'media' || activeTab === 'accounts'" class="flex-1 min-h-0" :admin="admin" :view="activeTab" @navigate="(t) => activeTab = t" />
        <WhitelistView v-else-if="activeTab === 'whitelist'" class="overflow-y-auto flex-1" :whitelist-list="admin.whitelists?.value ?? admin.whitelists" :loading="admin.isWhitelistLoading?.value ?? admin.isWhitelistLoading" :status-message="null" @fetch="admin.fetchWhitelists" @add="async (t, v) => { await admin.addWhitelist(t, v); emit('whitelistUpdated'); }" @update="async (id, t, v, act) => { await admin.updateWhitelist(id, t, v, act); emit('whitelistUpdated'); }" @delete="async (id) => { await admin.deleteWhitelist(id); emit('whitelistUpdated'); }" @toggle="async (id) => { await admin.toggleWhitelist(id); emit('whitelistUpdated'); }" />
        <ConfigPortal v-else-if="activeTab === 'config'" class="overflow-y-auto flex-1" :config="admin.configForm" :loading-config="admin.isConfigLoading?.value ?? admin.isConfigLoading" :saving-config="admin.isConfigLoading?.value ?? admin.isConfigLoading" :save-status="admin.configSaved?.value ?? admin.configSaved ? { success: true, message: '設定を保存しました' } : null" @save-config="admin.saveConfig" @load-config="admin.fetchConfig" />
      </main>
    </div>
  </div>
</template>
