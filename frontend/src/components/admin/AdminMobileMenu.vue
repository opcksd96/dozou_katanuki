<!-- frontend/src/components/admin/AdminMobileMenu.vue (100行以下 - SPEC-PRINCIPLE-001) -->
<script setup lang="ts">
import { X, ArrowLeft, RefreshCw, Server, ExternalLink } from 'lucide-vue-next';
import { AdminTabId, adminTabGroups } from '../../models/adminTabs';

defineProps<{ isOpen: boolean; activeTab: AdminTabId; isStashOnline?: boolean; }>();
const emit = defineEmits<{
  (e: 'close'): void;
  (e: 'select', tab: AdminTabId): void;
  (e: 'backToTimeline'): void;
  (e: 'openStashWeb'): void;
  (e: 'hardReload'): void;
}>();
</script>

<template>
  <div v-if="isOpen" class="fixed inset-0 z-50 flex justify-end bg-black/60 backdrop-blur-sm select-none">
    <div class="fixed inset-0" @click="emit('close')"></div>
    <div class="relative w-72 max-w-[85vw] h-full bg-slate-900 border-l border-slate-700/80 p-4 shadow-2xl flex flex-col overflow-y-auto space-y-3.5">
      <div class="flex items-center justify-between border-b border-slate-800 pb-2.5">
        <span class="text-xs font-bold text-slate-200">管理ナビゲーション</span>
        <button @click="emit('close')" class="p-1 rounded-md text-slate-400 hover:text-white" aria-label="Close">
          <X class="w-4 h-4" />
        </button>
      </div>

      <!-- クイックアクション -->
      <div class="grid grid-cols-2 gap-2">
        <button @click="emit('openStashWeb'); emit('close');" class="flex items-center justify-center gap-1.5 py-2 rounded-xl bg-slate-800 hover:bg-slate-700 active:scale-98 border border-slate-700 text-xs font-mono transition-all cursor-pointer">
          <span :class="['w-1.5 h-1.5 rounded-full', isStashOnline ? 'bg-emerald-400 shadow-[0_0_8px_rgba(52,211,153,0.8)]' : 'bg-amber-500']"></span>
          <span class="text-slate-300 flex items-center gap-1"><Server class="w-3 h-3" />Stash</span>
          <ExternalLink class="w-2.5 h-2.5 text-slate-500 ml-0.5" />
        </button>
        <button @click="emit('hardReload'); emit('close');" class="flex items-center justify-center gap-1.5 py-2 rounded-xl bg-slate-800 hover:bg-slate-700 active:scale-98 text-slate-300 hover:text-white text-xs font-semibold border border-slate-700 cursor-pointer transition-all">
          <RefreshCw class="w-3.5 h-3.5" /><span>再読込</span>
        </button>
      </div>

      <!-- 最上部：タイムラインへ戻るボタン -->
      <button
        @click="emit('backToTimeline'); emit('close');"
        class="w-full flex items-center justify-center gap-2 py-2.5 px-3 rounded-xl bg-slate-800 hover:bg-slate-700 active:scale-98 text-slate-100 text-xs font-bold border border-slate-700 cursor-pointer shadow-md transition-all"
      >
        <ArrowLeft class="w-4 h-4 text-blue-400" />
        <span>タイムラインへ戻る</span>
      </button>

      <!-- タブグループ一覧 -->
      <div class="space-y-3 pt-1 border-t border-slate-800/60">
        <div v-for="group in adminTabGroups" :key="group.title" class="space-y-1">
          <div class="text-[10px] font-mono text-slate-400 px-2 font-semibold uppercase tracking-wider">
            {{ group.title }}
          </div>
          <div class="space-y-0.5">
            <button
              v-for="item in group.items"
              :key="item.id"
              @click="emit('select', item.id); emit('close');"
              :class="[
                'w-full flex items-center gap-2.5 px-3 py-2 rounded-xl text-xs font-medium transition-all text-left cursor-pointer active:scale-98',
                activeTab === item.id
                  ? 'bg-blue-600 text-white font-bold shadow-md shadow-blue-500/20'
                  : 'text-slate-300 hover:bg-slate-800/80 hover:text-white'
              ]"
            >
              <span class="text-sm">{{ item.icon }}</span>
              <span>{{ item.label }}</span>
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

