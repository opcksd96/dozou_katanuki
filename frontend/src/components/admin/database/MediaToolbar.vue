<script setup lang="ts">
import { computed } from 'vue';

const props = defineProps<{
  searchQuery: string; accountFilter: string; statusFilter: string; typeFilter: 'all' | 'image' | 'video';
  viewMode: 'large' | 'compact' | 'table'; onlyBookmarked: boolean; accounts: any[]; activeJob?: any;
  stats?: { total_count: number; image_count: number; video_count: number };
}>();

const emit = defineEmits<{
  (e: 'update:searchQuery', v: string): void; (e: 'update:accountFilter', v: string): void;
  (e: 'update:statusFilter', v: string): void; (e: 'update:typeFilter', v: 'all' | 'image' | 'video'): void;
  (e: 'update:viewMode', v: 'large' | 'compact' | 'table'): void; (e: 'update:onlyBookmarked', v: boolean): void;
  (e: 'openStash'): void; (e: 'startSmartRecovery'): void; (e: 'startThunder'): void; (e: 'reconcileStash'): void;
}>();

const isSmartRunning = computed(() => props.activeJob?.status === 'RUNNING' && props.activeJob?.type?.includes('SMART'));
const isThunderRunning = computed(() => props.activeJob?.status === 'RUNNING' && props.activeJob?.type?.includes('THUNDER'));
const isReconcileRunning = computed(() => props.activeJob?.status === 'RUNNING' && props.activeJob?.type?.includes('STASH'));
</script>

<template>
  <div class="p-2.5 bg-slate-900 border border-slate-800 rounded-xl flex flex-wrap items-center justify-between gap-2 text-xs">
    <!-- フィルタ・検索群 -->
    <div class="flex flex-wrap items-center gap-2">
      <input :value="searchQuery" @input="emit('update:searchQuery', ($event.target as HTMLInputElement).value)" placeholder="🔍 メディア/URL検索..." class="w-40 bg-slate-950 border border-slate-700 rounded px-2 py-1 text-slate-200" />
      <select :value="accountFilter" @change="emit('update:accountFilter', ($event.target as HTMLSelectElement).value)" class="bg-slate-950 border border-slate-700 rounded px-2 py-1 text-slate-200">
        <option value="all">全アカウント</option>
        <option v-for="a in accounts" :key="a.numeric_id" :value="a.numeric_id">@{{ a.username || a.handle || a.numeric_id }}</option>
      </select>
      <select :value="statusFilter" @change="emit('update:statusFilter', ($event.target as HTMLSelectElement).value)" class="bg-slate-950 border border-slate-700 rounded px-2 py-1 text-slate-200" :class="statusFilter === 'TRASH' ? 'border-rose-600 text-rose-300 font-bold' : ''">
        <option value="all">全ステータス</option>
        <option value="COMPLETED">COMPLETED</option>
        <option value="QUEUED">QUEUED</option>
        <option value="OUTSOURCED">OUTSOURCED</option>
        <option value="ESCALATED">ESCALATED</option>
        <option value="RETAINED">RETAINED</option>
        <option value="DEAD_404">DEAD_404</option>
        <option value="FAILED">FAILED</option>
        <option value="EXCLUDED">EXCLUDED</option>
        <option value="TRASH">🗑️ ゴミ箱</option>
      </select>
      <div class="flex rounded border border-slate-700 overflow-hidden text-[11px]">
        <button v-for="t in ['all', 'image', 'video'] as const" :key="t" @click="emit('update:typeFilter', t)" :class="typeFilter === t ? 'bg-purple-600 text-white' : 'bg-slate-950 text-slate-400'" class="px-2 py-1 uppercase cursor-pointer transition-colors">
          {{ t }}
        </button>
      </div>
      <button @click="emit('update:onlyBookmarked', !onlyBookmarked)" :class="onlyBookmarked ? 'bg-amber-600 text-white' : 'bg-slate-950 text-slate-400 border border-slate-700'" class="px-2 py-1 rounded text-[11px] cursor-pointer transition-colors active:scale-95">
        ★ ブックマーク
      </button>
    </div>

    <!-- 統合アクションボタン群 (明確なアイコンと実行中フィードバック) -->
    <div class="flex items-center gap-2">
      <button @click="emit('startSmartRecovery')" :disabled="isSmartRunning" :class="isSmartRunning ? 'opacity-80 ring-2 ring-blue-400 animate-pulse' : 'hover:from-blue-500 hover:to-indigo-500 active:scale-95 shadow-md shadow-blue-900/30'" class="px-3.5 py-1.5 bg-gradient-to-r from-blue-600 to-indigo-600 text-white rounded-lg text-xs font-bold flex items-center gap-1.5 transition-all duration-150 cursor-pointer disabled:cursor-not-allowed" title="自律回収: 直接取得 ➔ Motrix外注 ➔ Stash照合 ➔ 完了バインドを一括自律実行">
        <span :class="isSmartRunning ? 'animate-spin inline-block' : ''">{{ isSmartRunning ? '⏳' : '🪄' }}</span>
        <span>{{ isSmartRunning ? '回収実行中...' : 'スマート一括回収' }}</span>
      </button>
      <button @click="emit('startThunder')" :disabled="isThunderRunning" :class="isThunderRunning ? 'opacity-80 ring-2 ring-amber-400 animate-pulse' : 'hover:bg-amber-500 active:scale-95 shadow-sm'" class="px-2.5 py-1.5 bg-amber-600 text-white rounded-lg text-[11px] font-bold flex items-center gap-1.5 transition-all duration-150 cursor-pointer disabled:cursor-not-allowed" title="迅雷投入: 保留(RETAINED)メディアを Thunder (迅雷) へ一括転送">
        <span :class="isThunderRunning ? 'animate-spin inline-block' : ''">{{ isThunderRunning ? '⏳' : '⚡' }}</span>
        <span>{{ isThunderRunning ? '迅雷投入中...' : '迅雷一括投入' }}</span>
      </button>
      <button @click="emit('reconcileStash')" :disabled="isReconcileRunning" :class="isReconcileRunning ? 'opacity-80 ring-2 ring-purple-400 animate-pulse' : 'hover:bg-purple-500 active:scale-95 shadow-sm'" class="px-2.5 py-1.5 bg-purple-600 text-white rounded-lg text-[11px] font-bold flex items-center gap-1 transition-all duration-150 cursor-pointer disabled:cursor-not-allowed" title="Stash照合: ローカル保存ファイルをStashと照合してCOMPLETEDへ昇格">
        <span :class="isReconcileRunning ? 'animate-spin inline-block' : ''">{{ isReconcileRunning ? '⏳' : '📦' }}</span>
        <span>{{ isReconcileRunning ? '照合中...' : 'Stash照合' }}</span>
      </button>
      <button @click="emit('openStash')" class="px-2 py-1.5 bg-purple-950/90 hover:bg-purple-900 border border-purple-500/70 text-purple-200 hover:text-white font-bold rounded-lg text-[11px] shadow-sm transition-all duration-150 active:scale-95 cursor-pointer flex items-center gap-1" title="Stash WebUIをブラウザで開く">
        <span>🎛️ WebUI</span><span class="text-purple-400">↗</span>
      </button>
    </div>
  </div>
</template>
