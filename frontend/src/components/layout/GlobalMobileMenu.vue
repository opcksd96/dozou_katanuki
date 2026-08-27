<!-- frontend/src/components/layout/GlobalMobileMenu.vue (100行以下 - SPEC-PRINCIPLE-001) -->
<script setup lang="ts">
import { Settings, RefreshCw, Layers, Server, X } from 'lucide-vue-next';
import { reloadWindow } from '../../composables/useKeyboardReload';

defineProps<{
  isOpen: boolean;
  isOnline: boolean;
  activeArticleId?: string | null;
}>();

const emit = defineEmits<{
  (e: 'close'): void;
  (e: 'openAdmin'): void;
  (e: 'backToTimeline'): void;
}>();
</script>

<template>
  <div v-if="isOpen" class="fixed inset-0 z-50 flex flex-col justify-start items-end bg-black/60 backdrop-blur-sm select-none p-3 pt-14 font-sans">
    <div class="fixed inset-0" @click="emit('close')"></div>
    <div class="relative w-full max-w-xs bg-slate-900/95 border border-slate-700/80 rounded-2xl p-4 shadow-2xl space-y-3 backdrop-blur-xl">
      <div class="flex items-center justify-between border-b border-slate-800 pb-2">
        <span class="text-xs font-bold text-slate-300">クイックメニュー</span>
        <button @click="emit('close')" class="p-1 rounded-md text-slate-400 hover:text-white" aria-label="Close">
          <X class="w-4 h-4" />
        </button>
      </div>

      <div class="space-y-2">
        <!-- 設定・管理 ボタン (デスクトップと同一のスタイル) -->
        <button
          @click="emit('openAdmin'); emit('close');"
          class="w-full flex items-center justify-between px-3.5 py-2.5 bg-gradient-to-r from-blue-600 to-indigo-600 hover:from-blue-500 hover:to-indigo-500 active:scale-98 text-white rounded-xl text-xs font-bold shadow-lg shadow-blue-500/20 transition-all cursor-pointer"
        >
          <div class="flex items-center gap-2">
            <Settings class="w-4 h-4" />
            <span>設定・管理 (Admin Board)</span>
          </div>
          <span class="text-[10px] bg-white/20 px-1.5 py-0.5 rounded font-mono">開く</span>
        </button>

        <!-- 再読込 ボタン -->
        <button
          @click="reloadWindow(); emit('close');"
          class="w-full flex items-center gap-2 px-3.5 py-2 bg-slate-800/90 hover:bg-slate-700 active:scale-98 text-slate-200 rounded-xl text-xs font-semibold border border-slate-700/80 transition-all cursor-pointer"
        >
          <RefreshCw class="w-4 h-4 text-slate-300" />
          <span>再読込 (Reload)</span>
        </button>

        <!-- Stash 状態バッジ -->
        <div class="flex items-center justify-between px-3.5 py-2 rounded-xl bg-slate-950/80 border border-slate-800 text-xs font-mono">
          <div class="flex items-center gap-2">
            <Server class="w-4 h-4 text-slate-400" />
            <span class="text-slate-300">Stash サーバー</span>
          </div>
          <div class="flex items-center gap-1.5">
            <span :class="['w-2 h-2 rounded-full', isOnline ? 'bg-emerald-400 shadow-[0_0_8px_rgba(52,211,153,0.8)]' : 'bg-amber-500']"></span>
            <span :class="isOnline ? 'text-emerald-400 font-bold' : 'text-amber-400'">{{ isOnline ? 'ONLINE' : 'OFFLINE' }}</span>
          </div>
        </div>

        <!-- タイムラインに戻る (記事詳細閲覧時のみ) -->
        <button
          v-if="activeArticleId"
          @click="emit('backToTimeline'); emit('close');"
          class="w-full flex items-center gap-2 px-3.5 py-2 bg-slate-800/90 hover:bg-slate-700 active:scale-98 text-blue-400 rounded-xl text-xs font-semibold border border-slate-700/80 transition-all cursor-pointer"
        >
          <Layers class="w-4 h-4 text-blue-400" />
          <span>タイムラインに戻る</span>
        </button>
      </div>
    </div>
  </div>
</template>
