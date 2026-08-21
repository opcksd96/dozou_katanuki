<!-- frontend/src/components/admin/database/MediaManagementView.vue (100行以下) -->
<script setup lang="ts">
import { onMounted } from 'vue';

const props = defineProps<{
  mediaItems: any[];
  total: number;
  accounts: any[];
  accountFilter: string;
  statusFilter: string;
  loading: boolean;
  page: number;
  limit: number;
}>();

const emit = defineEmits<{
  (e: 'fetch'): void;
  (e: 'update:accountFilter', v: string): void;
  (e: 'update:statusFilter', v: string): void;
  (e: 'update:page', p: number): void;
  (e: 'retryMedia', mediaId: string): void;
}>();

const onAccountChange = (val: string) => {
  emit('update:accountFilter', val);
  emit('update:page', 1);
  emit('fetch');
};

const onStatusChange = (val: string) => {
  emit('update:statusFilter', val);
  emit('update:page', 1);
  emit('fetch');
};

onMounted(() => { emit('fetch'); });
</script>

<template>
  <div class="space-y-3 flex flex-col h-full">
    <!-- ツールバー: アカウント切替 & ステータスフィルター -->
    <div class="flex flex-wrap items-center justify-between gap-2 text-xs">
      <div class="flex items-center gap-2">
        <select :value="accountFilter" @change="onAccountChange(($event.target as HTMLSelectElement).value)" class="bg-slate-950 border border-slate-700 rounded-lg px-2.5 py-1 text-slate-200 font-mono">
          <option value="all">全アカウント (ALL)</option>
          <option v-for="acc in accounts" :key="acc.numeric_id" :value="acc.numeric_id">@{{ acc.username }}</option>
        </select>
        <select :value="statusFilter" @change="onStatusChange(($event.target as HTMLSelectElement).value)" class="bg-slate-950 border border-slate-700 rounded-lg px-2.5 py-1 text-slate-200 font-mono">
          <option value="all">全ステータス</option>
          <option value="COMPLETED">COMPLETED (確保済)</option>
          <option value="QUEUED">QUEUED (待機中)</option>
          <option value="DEAD_404">DEAD_404 (消失)</option>
        </select>
      </div>
      <div class="text-[11px] text-slate-400 font-mono flex items-center gap-2">
        <span>全 {{ total }} 件</span>
        <button @click="if (page > 1) { emit('update:page', page - 1); emit('fetch'); }" :disabled="page <= 1" class="px-2 py-1 bg-slate-800 rounded disabled:opacity-40">◀</button>
        <span>{{ page }}</span>
        <button @click="if (page * limit < total) { emit('update:page', page + 1); emit('fetch'); }" :disabled="page * limit >= total" class="px-2 py-1 bg-slate-800 rounded disabled:opacity-40">▶</button>
      </div>
    </div>

    <!-- メディアカードグリッド -->
    <div class="flex-1 overflow-y-auto max-h-[440px] pr-1">
      <div v-if="loading" class="p-12 text-center text-slate-400 text-xs">読み込み中...</div>
      <div v-else-if="!mediaItems || mediaItems.length === 0" class="p-12 text-center text-slate-500 text-xs">メディアデータがありません</div>
      <div v-else class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-3">
        <div v-for="m in mediaItems" :key="m.media_id" class="bg-slate-950 border border-slate-800 rounded-xl p-2.5 flex flex-col space-y-2 group shadow hover:border-slate-700 transition-colors">
          <!-- プレビュー枠 -->
          <div class="h-28 bg-slate-900 rounded-lg overflow-hidden flex items-center justify-center relative">
            <img v-if="m.type === 'image'" :src="m.download_url" class="w-full h-full object-cover" onerror="this.src='/assets/default-avatar.png'" />
            <div v-else class="flex flex-col items-center justify-center text-slate-400 gap-1 text-[11px]"><span>🎬 動画</span><span class="text-[9px] font-mono text-slate-500">MP4 / Video</span></div>
            <span class="absolute top-1.5 right-1.5 px-1.5 py-0.5 rounded text-[9px] font-mono font-bold" :class="m.download_status === 'COMPLETED' ? 'bg-emerald-950 text-emerald-300 border border-emerald-700/60' : 'bg-amber-950 text-amber-300 border border-amber-700/60'">{{ m.download_status }}</span>
          </div>
          <!-- メディア情報 -->
          <div class="text-[11px] font-mono space-y-0.5 flex-1">
            <div class="text-slate-200 font-bold truncate" :title="m.media_id">{{ m.media_id }}</div>
            <div class="text-slate-400 text-[10px]">@{{ m.username }} | {{ m.width }}x{{ m.height }}</div>
            <div v-if="m.failed_reason" class="text-[10px] text-rose-400 truncate" :title="m.failed_reason">⚠️ {{ m.failed_reason }}</div>
          </div>
          <!-- Stash 動線 ＆ アクション -->
          <div class="pt-1.5 border-t border-slate-850 flex items-center justify-between text-[10px] font-mono">
            <span v-if="m.stash_scene_id || m.stash_image_id" class="text-emerald-400 truncate max-w-[120px]" :title="m.stash_scene_id || m.stash_image_id">🎛️ Stash: {{ (m.stash_scene_id || m.stash_image_id)?.slice(0, 8) }}...</span>
            <span v-else class="text-slate-500">Stash未連携</span>
            <button @click="emit('retryMedia', m.media_id)" class="px-2 py-0.5 bg-slate-800 hover:bg-slate-700 text-blue-400 rounded text-[10px]">再取得</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
