<!-- frontend/src/components/admin/database/AccountHistoryTimeline.vue (100行以下) -->
<script setup lang="ts">
import { resolveHistoryAvatarUrl } from '../../../utils/avatar';

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
  <div class="p-4 bg-slate-900 border border-slate-800 rounded-xl space-y-3">
    <div class="flex items-center justify-between">
      <h4 class="text-xs font-bold uppercase tracking-wider text-slate-400">🕒 歴代アバター世代履歴 ({{ histories?.length || 0 }}件)</h4>
    </div>
    <div v-if="!histories?.length" class="text-xs text-slate-500 py-2">世代履歴レコードはありません</div>
    <div v-else class="space-y-2 max-h-64 overflow-y-auto pr-1">
      <div v-for="h in histories" :key="h.avatar_seq" class="p-2.5 bg-slate-950/80 rounded-lg border border-slate-800 flex items-center justify-between text-xs font-mono">
        <div class="flex items-center gap-3">
          <img :src="resolveHistoryAvatarUrl(h, 'twitter')" class="w-9 h-9 rounded-full object-cover border border-slate-700" />
          <div>
            <div class="text-slate-200 font-bold">第{{ h.avatar_seq }}世代: {{ h.display_name }}</div>
            <div class="text-[11px] text-slate-500">{{ h.observed_at ? new Date(h.observed_at).toLocaleDateString() : '-' }}</div>
          </div>
        </div>
        <div class="flex items-center gap-2">
          <label
            class="px-2.5 py-1 bg-slate-800 hover:bg-slate-700 text-purple-300 border border-purple-800/50 rounded cursor-pointer text-[11px]"
            @dragover.prevent @drop.prevent="handleDrop($event, `${username}_gen${h.avatar_seq}`)"
          >
            <span>📁 ドロップ / 選択</span>
            <input type="file" accept="image/*" class="hidden" @change="handleFileInput($event, `${username}_gen${h.avatar_seq}`)" />
          </label>
        </div>
      </div>
    </div>
  </div>
</template>
