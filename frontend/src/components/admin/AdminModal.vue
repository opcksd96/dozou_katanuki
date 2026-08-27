<!-- frontend/src/components/admin/AdminModal.vue (100行以下 - SPEC-ADMINBOARD-001) -->
<script setup lang="ts">
import { ref, reactive, watch, onMounted, onUnmounted } from 'vue';
import { useAdmin } from '../../composables/useAdmin';
import { Minus, Square, X, ArrowLeft, Settings2, Menu } from 'lucide-vue-next';
import { WindowMinimise, WindowToggleMaximise, Quit } from '../../../wailsjs/runtime/runtime';
import { AdminTabId } from '../../models/adminTabs';
import AdminNavSidebar from './AdminNavSidebar.vue';
import AdminMobileMenu from './AdminMobileMenu.vue';
import PluginHubView from './PluginHubView.vue';
import RelationExplorerView from './RelationExplorerView.vue';
import ConfigPortal from './ConfigPortal.vue';
import StashStatusView from './StashStatusView.vue';
import DownloaderConsoleView from './DownloaderConsoleView.vue';
import DatabaseView from './DatabaseView.vue';
import AuditReportView from './AuditReportView.vue';

const props = defineProps<{ isOpen: boolean }>();
const emit = defineEmits<{ (e: 'close'): void; (e: 'whitelistUpdated'): void; (e: 'jumpToTimelinePost', articleId: string): void; }>();

const activeTab = ref<AdminTabId>('accounts'), isMobileNavOpen = ref(false);
const salvageForm = reactive({ platform: 'twitter', account: '', source: 'all', limit: 0 }), importForm = reactive({ warcPath: '', offline: true }), selectedPlatform = ref('twitter');
const admin = useAdmin();
const close = () => emit('close');
const handleKey = (e: KeyboardEvent) => { if (e.key === 'Escape' && props.isOpen) close(); };

const refreshCurrentTab = (tab: AdminTabId) => {
  try {
    if (tab === 'plugins') { admin.fetchActiveJob?.(); admin.fetchJobList?.(); admin.fetchSkinCSS?.(selectedPlatform.value || 'twitter'); }
    else if (tab === 'config') admin.fetchConfig?.();
    else if (tab === 'accounts') admin.fetchAccounts?.();
    else if (tab === 'posts') { admin.fetchAccounts?.(); admin.searchArticles?.(); }
    else if (tab === 'media') { admin.fetchAccounts?.(); admin.fetchMedia?.(); }
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
  <div v-if="isOpen" class="fixed inset-0 z-50 w-screen h-screen max-h-[100dvh] bg-slate-950 flex flex-col overflow-hidden select-none">
    <header class="px-3 py-2 border-b border-white/10 flex items-center justify-between bg-slate-900/90 backdrop-blur-xl shrink-0 wails-drag">
      <div class="flex items-center gap-2 wails-no-drag truncate">
        <div class="w-6 h-6 rounded-lg bg-blue-600/20 border border-blue-500/30 flex items-center justify-center text-blue-400 shrink-0"><Settings2 class="w-3.5 h-3.5" /></div>
        <h2 class="text-xs font-bold text-slate-100 truncate">Admin Governance Portal</h2>
      </div>

      <div class="flex items-center gap-1.5 wails-no-drag shrink-0">
        <!-- スマホ時: ハンバーガーメニューボタン -->
        <button @click="isMobileNavOpen = true" class="md:hidden p-2 rounded-lg bg-slate-800 border border-slate-700 text-slate-300 hover:text-white active:scale-95 cursor-pointer flex items-center gap-1" aria-label="Open tab menu">
          <Menu class="w-4 h-4" />
          <span class="text-[11px] font-bold text-blue-400">タブ</span>
        </button>

        <button @click="close" class="hidden sm:flex items-center gap-1.5 px-2.5 py-1.5 rounded-lg bg-slate-800 hover:bg-slate-700 text-slate-300 hover:text-white text-xs font-semibold border border-slate-700 cursor-pointer active:scale-95">
          <ArrowLeft class="w-3.5 h-3.5" /><span>タイムラインへ</span>
        </button>
        <div class="hidden sm:flex items-center ml-1 border-l border-slate-700 pl-1.5">
          <button @click="() => { try { WindowMinimise(); } catch {} }" class="p-1 rounded text-slate-400 hover:text-white"><Minus class="w-3 h-3" /></button>
          <button @click="() => { try { WindowToggleMaximise(); } catch {} }" class="p-1 rounded text-slate-400 hover:text-white"><Square class="w-3 h-3" /></button>
          <button @click="() => { try { Quit(); } catch {} }" class="p-1 rounded text-slate-400 hover:text-rose-400"><X class="w-3 h-3" /></button>
        </div>
      </div>
    </header>

    <div class="flex-1 min-h-0 flex overflow-hidden">
      <AdminNavSidebar :active-tab="activeTab" @select="(t) => activeTab = t" />
      <main class="flex-1 min-h-0 overflow-y-auto p-2 sm:p-3 bg-slate-950 flex flex-col">
        <PluginHubView v-if="activeTab === 'plugins'" :admin="admin" :salvage-form="salvageForm" :import-form="importForm" v-model:selected-platform="selectedPlatform" @start-salvage="admin.startSalvage(salvageForm.platform, salvageForm.account, salvageForm.limit, salvageForm.source)" @start-import="admin.startManualImport(importForm.warcPath, importForm.offline)" />
        <RelationExplorerView v-else-if="activeTab === 'explorer'" class="overflow-y-auto flex-1" />
        <AuditReportView v-else-if="activeTab === 'audit'" class="overflow-y-auto flex-1" :restoring="admin.isJobRunning?.value ?? admin.isJobRunning" @trigger-restore="(resetDB) => { admin.triggerRestore('', resetDB); activeTab = 'plugins'; }" />
        <StashStatusView v-else-if="activeTab === 'stash'" class="overflow-y-auto flex-1" :config="admin.configForm" />
        <DownloaderConsoleView v-else-if="activeTab === 'downloaders'" class="overflow-y-auto flex-1" :admin="admin" />
        <DatabaseView v-else-if="activeTab === 'posts' || activeTab === 'media' || activeTab === 'accounts'" class="flex-1 min-h-0" :admin="admin" :view="activeTab" @navigate="(t) => activeTab = t" @jump-to-timeline-post="(artId) => emit('jumpToTimelinePost', artId)" />
        <ConfigPortal v-else-if="activeTab === 'config'" class="overflow-y-auto flex-1" :config="admin.configForm" :loading-config="admin.isConfigLoading?.value ?? admin.isConfigLoading" :saving-config="admin.isConfigLoading?.value ?? admin.isConfigLoading" :save-status="admin.configSaved?.value ?? admin.configSaved ? { success: true, message: '設定を保存しました' } : null" @save-config="admin.saveConfig" @load-config="admin.fetchConfig" />
      </main>
    </div>

    <!-- モバイル用ハンバーガーメニュー・ドロワー -->
    <AdminMobileMenu :is-open="isMobileNavOpen" :active-tab="activeTab" @close="isMobileNavOpen = false" @select="(t) => activeTab = t" @back-to-timeline="close" />
  </div>
</template>


