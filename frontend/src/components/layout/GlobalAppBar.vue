<!-- frontend/src/components/layout/GlobalAppBar.vue (100行以下) -->
<script setup lang="ts">
import { computed, ref, onMounted, onUnmounted } from 'vue';
import { reloadWindow } from '../../composables/useKeyboardReload';
import { EventsOn, WindowMinimise, WindowToggleMaximise, Quit } from '../../../wailsjs/runtime/runtime';
import { RefreshCw, Settings, Minus, Square, X, Server, Layers } from 'lucide-vue-next';

const props = defineProps<{ activeArticleHandle?: string; activeArticleId?: string | null; isStashOnline?: boolean }>();
const emit = defineEmits<{ (e: 'openAdmin'): void; (e: 'backToTimeline'): void }>();

const localOnline = ref(false);
const isOnline = computed(() => props.isStashOnline ?? localOnline.value);
let unoff: (() => void) | null = null;

const checkStash = async () => {
  try {
    const res = await fetch('/stash-proxy/', { method: 'HEAD' });
    localOnline.value = res.ok || res.status === 401 || res.status === 404;
  } catch { localOnline.value = false; }
};

const handleMinimise = () => { try { WindowMinimise(); } catch {} };
const handleToggleMaximise = () => { try { WindowToggleMaximise(); } catch {} };
const handleQuit = () => { try { Quit(); } catch {} };

onMounted(() => {
  try { if ((window as any)?.runtime?.EventsOnMultiple) unoff = EventsOn('stash:ready', (ready: boolean) => { localOnline.value = !!ready; }); } catch {}
  checkStash();
});
onUnmounted(() => { if (unoff) try { unoff(); } catch {} });
</script>

<template>
  <nav class="w-full bg-slate-950/70 backdrop-blur-xl border-b border-white/10 sticky top-0 z-40 px-3 py-2 flex items-center justify-between select-none wails-drag transition-colors duration-200">
    <!-- Left: App Brand & Navigation Breadcrumb -->
    <div class="flex items-center gap-3 wails-no-drag">
      <div @click="emit('backToTimeline')" class="flex items-center gap-2 cursor-pointer group" title="トップに戻る">
        <div class="w-6 h-6 rounded-lg bg-gradient-to-br from-blue-500 to-indigo-600 flex items-center justify-center text-white font-bold text-xs shadow-md shadow-blue-500/20 group-hover:scale-105 transition-transform">
          <Layers class="w-3.5 h-3.5 text-white" />
        </div>
        <div class="flex items-center gap-1.5">
          <span class="text-xs font-bold text-slate-100 tracking-tight group-hover:text-blue-400 transition-colors">dozou_katanuki</span>
          <span class="text-[9px] px-1.5 py-0.2 rounded-full bg-slate-800/80 border border-slate-700/50 text-slate-400 font-mono">v4.0</span>
        </div>
      </div>

      <div v-if="activeArticleId" class="hidden sm:flex items-center gap-1.5 text-xs text-slate-400 font-mono pl-3 border-l border-slate-800/80">
        <button @click="emit('backToTimeline')" class="hover:text-blue-400 transition-colors cursor-pointer">タイムライン</button>
        <span class="text-slate-600">/</span>
        <span class="text-slate-300">@{{ activeArticleHandle }}</span>
        <span class="text-slate-600">/</span>
        <span class="text-blue-400 font-semibold">ポスト</span>
      </div>
    </div>

    <!-- Right: Status, Actions & Window Controls -->
    <div class="flex items-center gap-2">
      <!-- App Status & Actions -->
      <div class="flex items-center gap-1.5 wails-no-drag">
        <div class="hidden sm:flex items-center gap-1.5 px-2 py-1 rounded-lg bg-slate-900/80 border border-slate-800/80 text-[11px] font-mono shadow-sm">
          <span :class="['w-1.5 h-1.5 rounded-full', isOnline ? 'bg-emerald-400 shadow-[0_0_8px_rgba(52,211,153,0.8)]' : 'bg-amber-500']"></span>
          <span :class="isOnline ? 'text-slate-300' : 'text-slate-500'" class="flex items-center gap-1">
            <Server class="w-3 h-3 text-slate-400 inline" />
            Stash
          </span>
        </div>

        <button @click="reloadWindow()" class="p-1.5 rounded-lg bg-slate-900/80 hover:bg-slate-800 text-slate-300 hover:text-white border border-slate-800/80 hover:border-slate-700 text-xs flex items-center gap-1.5 transition-all shadow-sm cursor-pointer" title="再読み込み (Ctrl+R)">
          <RefreshCw class="w-3.5 h-3.5" />
          <span class="hidden md:inline text-[11px]">再読み込み</span>
        </button>

        <button @click="emit('openAdmin')" class="flex items-center gap-1.5 px-2.5 py-1.5 bg-gradient-to-r from-blue-600 to-indigo-600 hover:from-blue-500 hover:to-indigo-500 text-white rounded-lg text-xs font-semibold shadow-md shadow-blue-500/20 transition-all hover:scale-[1.02] cursor-pointer" title="設定・Admin Board">
          <Settings class="w-3.5 h-3.5" />
          <span class="hidden sm:inline text-[11px]">設定・管理</span>
        </button>
      </div>

      <!-- Window Controls (Minimize / Maximize / Close) -->
      <div class="flex items-center ml-2 border-l border-slate-800/80 pl-2 wails-no-drag">
        <button @click="handleMinimise" class="p-1.5 rounded-md hover:bg-slate-800 text-slate-400 hover:text-slate-200 transition-colors" title="最小化">
          <Minus class="w-3.5 h-3.5" />
        </button>
        <button @click="handleToggleMaximise" class="p-1.5 rounded-md hover:bg-slate-800 text-slate-400 hover:text-slate-200 transition-colors" title="最大化 / 復元">
          <Square class="w-3 h-3" />
        </button>
        <button @click="handleQuit" class="p-1.5 rounded-md hover:bg-red-600/80 text-slate-400 hover:text-white transition-colors" title="閉じる">
          <X class="w-3.5 h-3.5" />
        </button>
      </div>
    </div>
  </nav>
</template>

