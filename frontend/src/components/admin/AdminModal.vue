<!-- frontend/src/components/admin/AdminModal.vue (100行以下 - SPEC-ADMINBOARD-001) -->
<script setup lang="ts">
import { ref, reactive, watch, onMounted, onUnmounted } from 'vue';
import { useAdmin } from '../../composables/useAdmin';
import { Minus, Square, X, ArrowLeft, Settings2, RefreshCw, Server, ExternalLink } from 'lucide-vue-next';
import { WindowMinimise, WindowToggleMaximise, Quit, EventsOn, BrowserOpenURL } from '../../../wailsjs/runtime/runtime';
import { AdminTabId } from '../../models/adminTabs';
import AdminNavSidebar from './AdminNavSidebar.vue';
import AdminMobileMenu from './AdminMobileMenu.vue';
import PluginHubView from './PluginHubView.vue';
import RelationExplorerView from './RelationExplorerView.vue';
import ConfigPortal from './ConfigPortal.vue';
import PipelineConsoleView from './PipelineConsoleView.vue';
import AuditReportView from './AuditReportView.vue';
import SystemConsoleView from './SystemConsoleView.vue';
import DatabaseView from './DatabaseView.vue';

const props = defineProps<{ isOpen: boolean }>();
const emit = defineEmits<{ (e: 'close'): void; (e: 'whitelistUpdated'): void; (e: 'jumpToTimelinePost', articleId: string): void; }>();
const initialTab = (sessionStorage.getItem('admin_auto_open_tab') as AdminTabId) || 'plugins';
sessionStorage.removeItem('admin_auto_open_tab');
const activeTab = ref<AdminTabId>(initialTab), isMobileNavOpen = ref(false), isStashOnline = ref(false);
const salvageForm = reactive({ platform: 'twitter', account: '', source: 'all', limit: 0 }), importForm = reactive({ warcPath: '', offline: true }), selectedPlatform = ref('twitter');
const admin = useAdmin();
let unoffStash: (() => void) | null = null;
const close = () => emit('close');
const handleKey = (e: KeyboardEvent) => { if (e.key === 'Escape' && props.isOpen) close(); };
const checkStash = async () => {
  try {
    const getApp = (window as any)?.go?.app?.App || (window as any)?.go?.main?.App;
    if (getApp?.IsStashReady && await getApp.IsStashReady()) { isStashOnline.value = true; return; }
    const r = await fetch('/stash-proxy/', { method: 'HEAD' }); isStashOnline.value = r.ok || r.status === 401 || r.status === 404;
  } catch { isStashOnline.value = false; }
};
const openStashWeb = () => {
  try { BrowserOpenURL('http://127.0.0.1:9999/'); }
  catch { window.open('http://127.0.0.1:9999/', '_blank'); }
};
const handleGlobalHardReload = () => {
  sessionStorage.setItem('admin_auto_open_tab', activeTab.value);
  window.location.reload();
};
watch(() => props.isOpen, (open) => { if (open) { checkStash(); document.body.style.overflow = 'hidden'; } else document.body.style.overflow = ''; }, { immediate: true });
onMounted(() => {
  window.addEventListener('keydown', handleKey);
  try { if ((window as any)?.runtime?.EventsOnMultiple) unoffStash = EventsOn('stash:ready', (ready: boolean) => { isStashOnline.value = !!ready; }); } catch {}
  checkStash();
});
onUnmounted(() => { window.removeEventListener('keydown', handleKey); if (unoffStash) try { unoffStash(); } catch {}; document.body.style.overflow = ''; });
</script>

<template>
  <div v-if="isOpen" class="fixed inset-0 z-50 w-screen h-screen max-h-[100dvh] bg-slate-950 flex flex-col overflow-hidden select-none font-sans">
    <header class="px-3 py-2 border-b border-white/10 flex items-center justify-between bg-slate-900/90 backdrop-blur-xl shrink-0 wails-drag">
      <div class="flex items-center gap-2 wails-no-drag truncate">
        <div class="w-6 h-6 rounded-lg bg-blue-600/20 border border-blue-500/30 flex items-center justify-center text-blue-400 shrink-0"><Settings2 class="w-3.5 h-3.5" /></div>
        <h2 class="text-xs font-bold text-slate-100 truncate">Admin Governance Portal</h2>
      </div>
      <div class="flex items-center gap-1.5 wails-no-drag shrink-0">
        <button @click="openStashWeb" class="flex items-center gap-1.5 px-2 py-1 rounded-lg bg-slate-800/80 hover:bg-slate-700/80 border border-slate-700/80 text-[11px] font-mono transition-all cursor-pointer active:scale-95 group" title="Stash Web UI を開く">
          <span :class="['w-1.5 h-1.5 rounded-full', isStashOnline ? 'bg-emerald-400 shadow-[0_0_8px_rgba(52,211,153,0.8)]' : 'bg-amber-500']"></span>
          <span class="text-slate-400 group-hover:text-white flex items-center gap-1"><Server class="w-3 h-3" />Stash</span>
          <ExternalLink class="w-2.5 h-2.5 text-slate-500 group-hover:text-blue-400 ml-0.5" />
        </button>
        <button @click="handleGlobalHardReload" class="px-2.5 py-1.5 rounded-lg bg-slate-800 hover:bg-slate-700 text-slate-300 hover:text-white text-xs font-semibold border border-slate-700 cursor-pointer active:scale-95 flex items-center gap-1.5" title="現在のタブを保持して再読み込み">
          <RefreshCw class="w-3.5 h-3.5" /><span class="hidden sm:inline">再読込</span>
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
        <PipelineConsoleView v-else-if="activeTab === 'pipeline'" class="overflow-y-auto flex-1" @jump-to-media="(mId) => { activeTab = 'media'; admin.fetchMedia?.(); }" />
        <AuditReportView v-else-if="activeTab === 'audit'" class="overflow-y-auto flex-1" :restoring="admin.isJobRunning?.value ?? admin.isJobRunning" @trigger-restore="(resetDB) => { admin.triggerRestore('', resetDB); activeTab = 'plugins'; }" />
        <SystemConsoleView v-else-if="activeTab === 'console'" class="flex-1 min-h-0" />
        <RelationExplorerView v-else-if="activeTab === 'explorer'" class="overflow-y-auto flex-1" />
        <DatabaseView v-else-if="activeTab === 'posts' || activeTab === 'media' || activeTab === 'accounts'" class="flex-1 min-h-0" :admin="admin" :view="activeTab" @navigate="(t) => activeTab = t" @jump-to-timeline-post="(artId) => emit('jumpToTimelinePost', artId)" />
        <ConfigPortal v-else-if="activeTab === 'config'" class="overflow-y-auto flex-1" :config="admin.configForm" :loading-config="admin.isConfigLoading?.value ?? admin.isConfigLoading" :saving-config="admin.isConfigLoading?.value ?? admin.isConfigLoading" :save-status="admin.configSaved?.value ?? admin.configSaved ? { success: true, message: '設定を保存しました' } : null" @save-config="admin.saveConfig" @load-config="admin.fetchConfig" />
      </main>
    </div>
    <AdminMobileMenu :is-open="isMobileNavOpen" :active-tab="activeTab" @close="isMobileNavOpen = false" @select="(t) => activeTab = t" @back-to-timeline="close" />
  </div>
</template>
