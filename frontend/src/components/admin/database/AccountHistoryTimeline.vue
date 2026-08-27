<!-- frontend/src/components/admin/database/AccountHistoryTimeline.vue (100行以下 - SPEC-PRINCIPLE-001) -->
<script setup lang="ts">
import { resolveHistoryAvatarUrl } from '../../../utils/avatar';
import { formatDate } from '../../../utils/formatters';

defineProps<{
  histories: any[];
  username: string;
}>();

const emit = defineEmits<{
  (e: 'uploadAvatar', payload: { virtualKey: string; base64Data: string }): void;
}>();

const processFile = (file: File, virtualKey: string) => {
  if (!file || !file.type.startsWith('image/')) return;
  const reader = new FileReader();
  reader.onload = (e) => {
    const base64 = e.target?.result as string;
    if (base64) emit('uploadAvatar', { virtualKey, base64Data: base64 });
  };
  reader.readAsDataURL(file);
};

const handleDrop = (e: DragEvent, virtualKey: string) => {
  const files = e.dataTransfer?.files;
  if (files && files.length > 0) processFile(files[0], virtualKey);
};

const handleFileInput = (e: Event, virtualKey: string) => {
  const input = e.target as HTMLInputElement;
  if (input.files && input.files.length > 0) {
    processFile(input.files[0], virtualKey);
    input.value = '';
  }
};
</script>

<template>
  <div class="p-4 bg-slate-900/90 border border-slate-800 rounded-xl space-y-3 font-sans">
    <div class="flex items-center justify-between">
      <h4 class="text-xs font-bold uppercase tracking-wider text-slate-400">🕒 歴代プロファイル変遷史 ({{ histories?.length || 0 }}世代)</h4>
    </div>
    <div v-if="!histories?.length" class="text-xs text-slate-500 py-2">世代履歴レコードはありません</div>
    <div v-else class="space-y-2.5 max-h-72 overflow-y-auto pr-1">
      <div v-for="h in histories" :key="h.id || h.avatar_seq" class="p-3 bg-slate-950/85 rounded-lg border border-slate-800/80 space-y-2 text-xs">
        <div class="flex items-center justify-between">
          <div class="flex items-center gap-3 min-w-0">
            <img :src="resolveHistoryAvatarUrl(h, 'twitter')" class="w-10 h-10 rounded-full object-cover border border-slate-700 shrink-0 shadow-md" />
            <div class="min-w-0">
              <div class="text-slate-100 font-bold flex items-center gap-1.5 truncate">
                <span class="px-1.5 py-0.2 rounded bg-indigo-950 text-indigo-300 font-mono text-[10px] font-semibold border border-indigo-800/60 shrink-0">第{{ h.avatar_seq }}世代</span>
                <span class="truncate">{{ h.display_name || username }}</span>
              </div>
              <div class="text-[11px] text-slate-400 font-mono flex items-center gap-1 pt-0.5">
                <span>📅 投稿時期:</span>
                <span class="text-slate-300">{{ formatDate(h.observed_at) || '-' }}</span>
              </div>
            </div>
          </div>

          <!-- アバター手動登録 -->
          <label class="px-2.5 py-1 bg-slate-800 hover:bg-slate-700 text-purple-300 border border-purple-800/50 rounded cursor-pointer text-[11px] font-mono shrink-0 active:scale-95 transition-transform" @dragover.prevent @drop.prevent="handleDrop($event, `${username}_gen${h.avatar_seq}`)">
            <span>📁 画像登録</span>
            <input type="file" accept="image/*" class="hidden" @change="handleFileInput($event, `${username}_gen${h.avatar_seq}`)" />
          </label>
        </div>

        <!-- 一言コメント / 自己紹介文 (bio) の変遷プレビュー -->
        <div v-if="h.description" class="text-[11px] text-slate-300 bg-slate-900/80 px-2.5 py-1.5 rounded border border-slate-800/60 leading-relaxed font-sans select-text">
          <span class="text-slate-500 font-mono mr-1">💬</span>{{ h.description }}
        </div>
      </div>
    </div>
  </div>
</template>
