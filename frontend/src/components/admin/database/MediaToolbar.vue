<!-- frontend/src/components/admin/database/MediaToolbar.vue (100行以下) -->
<script setup lang="ts">
defineProps<{
  searchQuery: string;
  accountFilter: string;
  statusFilter: string;
  typeFilter: 'all' | 'image' | 'video';
  viewMode: 'large' | 'compact' | 'table';
  onlyBookmarked: boolean;
  accounts: any[];
  stats?: { total_count: number; image_count: number; video_count: number };
}>();

const emit = defineEmits<{
  (e: 'update:searchQuery', v: string): void;
  (e: 'update:accountFilter', v: string): void;
  (e: 'update:statusFilter', v: string): void;
  (e: 'update:typeFilter', v: 'all' | 'image' | 'video'): void;
  (e: 'update:viewMode', v: 'large' | 'compact' | 'table'): void;
  (e: 'update:onlyBookmarked', v: boolean): void;
  (e: 'openStash'): void;
  (e: 'startDownload'): void;
  (e: 'startPoll'): void;
  (e: 'startEscalate'): void;
  (e: 'requeueFailed'): void;
  (e: 'reconcileStash'): void;
}>();
</script>

<template>
  <div class="p-2.5 bg-slate-900 border border-slate-800 rounded-xl flex flex-wrap items-center justify-between gap-2 text-xs">
    <!-- フィルタ・検索群 -->
    <div class="flex flex-wrap items-center gap-2">
      <input :value="searchQuery" @input="emit('update:searchQuery', ($event.target as HTMLInputElement).value)" placeholder="🔍 メディア/URL検索..." class="w-40 bg-slate-950 border border-slate-700 rounded px-2 py-1 text-slate-200" />
      <select :value="accountFilter" @change="emit('update:accountFilter', ($event.target as HTMLSelectElement).value)" class="bg-slate-950 border border-slate-700 rounded px-2 py-1 text-slate-200">
        <option value="all">全アカウント</option>
        <option v-for="a in accounts" :key="a.numeric_id" :value="a.numeric_id">@{{ a.username }}</option>
      </select>
      <select :value="statusFilter" @change="emit('update:statusFilter', ($event.target as HTMLSelectElement).value)" class="bg-slate-950 border border-slate-700 rounded px-2 py-1 text-slate-200">
        <option value="all">全ステータス</option>
        <option value="COMPLETED">COMPLETED</option>
        <option value="QUEUED">QUEUED</option>
        <option value="FAILED">FAILED</option>
        <option value="DEAD_404">DEAD_404</option>
        <option value="EXCLUDED">EXCLUDED</option>
      </select>
      <div class="flex rounded border border-slate-700 overflow-hidden text-[11px]">
        <button v-for="t in ['all', 'image', 'video'] as const" :key="t" @click="emit('update:typeFilter', t)" :class="typeFilter === t ? 'bg-purple-600 text-white' : 'bg-slate-950 text-slate-400'" class="px-2 py-1 uppercase">
          {{ t }}
        </button>
      </div>
      <button @click="emit('update:onlyBookmarked', !onlyBookmarked)" :class="onlyBookmarked ? 'bg-amber-600 text-white' : 'bg-slate-950 text-slate-400 border border-slate-700'" class="px-2 py-1 rounded text-[11px]">
        ★ ブックマーク
      </button>
    </div>

    <!-- アクションボタン群 -->
    <div class="flex items-center gap-1.5">
      <button @click="emit('startDownload')" class="px-2.5 py-1 bg-blue-600 hover:bg-blue-500 text-white rounded text-[11px] font-bold">📥 回収開始</button>
      <button @click="emit('startEscalate')" class="px-2.5 py-1 bg-indigo-600 hover:bg-indigo-500 text-white rounded text-[11px]" title="DEAD_404をMotrix/Aria2へ外注委託">🚀 Motrix外注</button>
      <button @click="emit('requeueFailed')" class="px-2.5 py-1 bg-amber-600 hover:bg-amber-500 text-white rounded text-[11px]">🔄 404再試行</button>
      <button @click="emit('reconcileStash')" class="px-2.5 py-1 bg-purple-600 hover:bg-purple-500 text-white rounded text-[11px]">📦 Stash同期</button>
      <button @click="emit('openStash')" class="px-2.5 py-1 bg-slate-800 hover:bg-slate-700 text-purple-300 rounded text-[11px]">WebUI ↗</button>
    </div>
  </div>
</template>
