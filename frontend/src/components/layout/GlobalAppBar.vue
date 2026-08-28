<!-- frontend/src/components/layout/GlobalAppBar.vue (100行以下 - SPEC-PRINCIPLE-001) -->
<script setup lang="ts">
import { computed, ref, onMounted, onUnmounted } from 'vue';
import { reloadWindow } from '../../composables/useKeyboardReload';
import { EventsOn, WindowMinimise, WindowToggleMaximise, Quit, BrowserOpenURL } from '../../../wailsjs/runtime/runtime';
import { RefreshCw, Settings, Minus, Square, X, Server, Layers, Menu, ExternalLink } from 'lucide-vue-next';
import GlobalMobileMenu from './GlobalMobileMenu.vue';

const props = defineProps<{ activeArticleHandle?: string; activeArticleId?: string | null; isStashOnline?: boolean }>();
const emit = defineEmits<{ (e: 'openAdmin'): void; (e: 'backToTimeline'): void }>();

const localOnline = ref(false), isMobileMenuOpen = ref(false);
const isOnline = computed(() => props.isStashOnline ?? localOnline.value);
let unoff: (() => void) | null = null;

const checkStash = async () => {
  try {
    const res = await fetch('/stash-proxy/', { method: 'HEAD' });
    localOnline.value = res.ok || res.status === 401 || res.status === 404;
  } catch { localOnline.value = false; }
};

const openStashWeb = () => {
  try { BrowserOpenURL('http://127.0.0.1:9999/'); }
  catch { window.open('http://127.0.0.1:9999/', '_blank'); }
};

onMounted(() => {
  try { if ((window as any)?.runtime?.EventsOnMultiple) unoff = EventsOn('stash:ready', (ready: boolean) => { localOnline.value = !!ready; }); } catch {}
  checkStash();
});
onUnmounted(() => { if (unoff) try { unoff(); } catch {} });
</script>

<template>
  <nav class="w-full bg-slate-950/70 backdrop-blur-xl border-b border-white/10 sticky top-0 z-40 px-3 py-2 flex items-center justify-between select-none wails-drag transition-colors font-sans">
    <div class="flex items-center gap-2.5 wails-no-drag">
      <div @click="emit('backToTimeline')" class="flex items-center gap-2 cursor-pointer group" title="トップに戻る">
        <div class="w-7 h-7 rounded-lg bg-gradient-to-br from-blue-500 to-indigo-600 flex items-center justify-center text-white font-bold text-xs shadow-md group-hover:scale-105 transition-transform">
          <Layers class="w-4 h-4 text-white" />
        </div>
        <div class="flex items-center gap-1.5">
          <span class="text-xs font-bold text-slate-100 group-hover:text-blue-400">dozou_katanuki</span>
          <span class="text-[9px] px-1.5 py-0.2 rounded-full bg-slate-800/80 border border-slate-700/50 text-slate-400 font-mono">v4.0</span>
        </div>
      </div>

      <div v-if="activeArticleId" class="hidden sm:flex items-center gap-1 text-xs text-slate-400 font-mono pl-2 border-l border-slate-800/80 truncate max-w-[200px]">
        <button @click="emit('backToTimeline')" class="hover:text-blue-400 cursor-pointer">TL</button>
        <span>/</span><span class="text-slate-300 truncate">@{{ activeArticleHandle }}</span>
      </div>
    </div>

    <div class="flex items-center gap-1.5">
      <button @click="isMobileMenuOpen = true" class="sm:hidden p-2 rounded-lg bg-slate-900 border border-slate-800 text-slate-300 hover:text-white active:scale-95 cursor-pointer flex items-center gap-1.5" aria-label="Open menu">
        <Menu class="w-4 h-4" />
      </button>

      <div class="hidden sm:flex items-center gap-1.5 wails-no-drag">
        <!-- Stash パイロットランプ (クリックで Stash Web UI へダイレクトジャンプ) -->
        <button @click="openStashWeb" class="flex items-center gap-1.5 px-2 py-1 rounded-lg bg-slate-900/80 hover:bg-slate-800 border border-slate-800 hover:border-slate-700 text-[11px] font-mono transition-all cursor-pointer active:scale-95 group" title="Stash Web UI をブラウザで開く (http://127.0.0.1:9999/)">
          <span :class="['w-1.5 h-1.5 rounded-full', isOnline ? 'bg-emerald-400 shadow-[0_0_8px_rgba(52,211,153,0.8)]' : 'bg-amber-500']"></span>
          <span :class="isOnline ? 'text-slate-300 group-hover:text-white' : 'text-slate-500'" class="flex items-center gap-1"><Server class="w-3 h-3 text-slate-400" />Stash</span>
          <ExternalLink class="w-2.5 h-2.5 text-slate-500 group-hover:text-blue-400 ml-0.5" />
        </button>
        <button @click="reloadWindow()" class="p-1.5 rounded-lg bg-slate-900/80 hover:bg-slate-800 text-slate-300 hover:text-white border border-slate-800/80 text-xs flex items-center gap-1.5 cursor-pointer" title="再読み込み (Ctrl+R)">
          <RefreshCw class="w-3.5 h-3.5" /><span class="hidden md:inline text-[11px]">再読込</span>
        </button>
        <button @click="emit('openAdmin')" class="flex items-center gap-1.5 px-2.5 py-1.5 bg-gradient-to-r from-blue-600 to-indigo-600 hover:from-blue-500 text-white rounded-lg text-xs font-semibold shadow-md cursor-pointer" title="設定・Admin Board">
          <Settings class="w-3.5 h-3.5" /><span>設定・管理</span>
        </button>
      </div>

      <div class="hidden sm:flex items-center ml-1 border-l border-slate-800/80 pl-1.5 wails-no-drag">
        <button @click="() => { try { WindowMinimise(); } catch {} }" class="p-1.5 rounded-md hover:bg-slate-800 text-slate-400 hover:text-slate-200" title="最小化"><Minus class="w-3.5 h-3.5" /></button>
        <button @click="() => { try { WindowToggleMaximise(); } catch {} }" class="p-1.5 rounded-md hover:bg-slate-800 text-slate-400 hover:text-slate-200" title="最大化"><Square class="w-3.5 h-3.5" /></button>
        <button @click="() => { try { Quit(); } catch {} }" class="p-1.5 rounded-md hover:bg-red-600/80 text-slate-400 hover:text-white" title="閉じる"><X class="w-3.5 h-3.5" /></button>
      </div>
    </div>

    <GlobalMobileMenu :is-open="isMobileMenuOpen" :is-online="isOnline" :active-article-id="activeArticleId" @close="isMobileMenuOpen = false" @open-admin="emit('openAdmin')" @back-to-timeline="emit('backToTimeline')" />
  </nav>
</template>
