<!-- frontend/src/components/admin/database/MediaInspectorStash.vue (100行以下) -->
<script setup lang="ts">
import { computed } from 'vue';
import { BrowserOpenURL } from '../../../../wailsjs/runtime/runtime';

const props = defineProps<{
  media: any;
  stashData: any;
  editTitle: string;
  editDetails: string;
  editRating: number;
  isModified: boolean;
  isMutating: boolean;
  hasUndo: boolean;
}>();

const emit = defineEmits<{
  (e: 'update:editTitle', val: string): void;
  (e: 'update:editDetails', val: string): void;
  (e: 'update:editRating', val: number): void;
  (e: 'requestSave'): void;
  (e: 'undo'): void;
}>();

const currentHostname = computed(() => (typeof window !== 'undefined' && window.location?.hostname) ? window.location.hostname : '127.0.0.1');
const stashDirectUrl = computed(() => {
  if (props.media.stash_scene_id) return `http://${currentHostname.value}:9999/scenes/${props.media.stash_scene_id}`;
  if (props.media.stash_image_id) return `http://${currentHostname.value}:9999/images/${props.media.stash_image_id}`;
  return null;
});

const openStash = () => {
  if (!stashDirectUrl.value) return;
  try { BrowserOpenURL(stashDirectUrl.value); } catch { window.open(stashDirectUrl.value, '_blank', 'noopener,noreferrer'); }
};
</script>

<template>
  <div v-if="stashDirectUrl" class="space-y-3">
    <!-- Stash 情報表示 -->
    <div class="p-3 bg-slate-950/80 rounded-xl border border-slate-800 space-y-2.5">
      <div class="flex items-center justify-between text-[10px]">
        <span class="text-slate-400 font-bold uppercase tracking-wider">📦 Stash 連携情報</span>
        <span class="text-emerald-400 font-bold">● 連携中</span>
      </div>
      <div class="p-2 bg-slate-900 rounded-lg border border-slate-800 text-[11px] text-purple-300 font-bold truncate">
        ID: {{ media.stash_scene_id || media.stash_image_id }}
      </div>
      <button @click="openStash" class="w-full py-1.5 bg-purple-950 hover:bg-purple-900 text-purple-200 border border-purple-700/60 font-bold rounded-lg transition-colors text-xs">
        🎛️ Stash WebUI で開く ↗
      </button>
    </div>

    <!-- Stash メタデータ更新 -->
    <div class="p-3 bg-slate-950/80 rounded-xl border border-slate-800 space-y-2.5">
      <div class="flex items-center justify-between text-[10px]">
        <span class="text-slate-400 font-bold uppercase tracking-wider">✏️ Stash メタデータ更新</span>
        <span v-if="isModified" class="text-amber-400 font-bold animate-pulse">● 未保存</span>
      </div>
      <div class="space-y-1">
        <label class="text-[10px] text-slate-400">タイトル</label>
        <input :value="editTitle" @input="emit('update:editTitle', ($event.target as HTMLInputElement).value)" type="text" class="w-full bg-slate-900 border border-slate-700 rounded px-2 py-1 text-slate-200 text-[11px]" />
      </div>
      <div class="space-y-1">
        <label class="text-[10px] text-slate-400">詳細メモ</label>
        <textarea :value="editDetails" @input="emit('update:editDetails', ($event.target as HTMLTextAreaElement).value)" rows="3" class="w-full bg-slate-900 border border-slate-700 rounded px-2 py-1 text-slate-200 text-[11px]"></textarea>
      </div>
      <div class="flex gap-2">
        <button @click="emit('requestSave')" :disabled="!isModified || isMutating" class="flex-1 py-1.5 bg-purple-600 hover:bg-purple-500 disabled:opacity-40 text-white font-bold rounded text-xs">
          {{ isMutating ? '更新中...' : 'Stash へ反映' }}
        </button>
        <button v-if="hasUndo" @click="emit('undo')" class="px-2 py-1.5 bg-slate-800 hover:bg-slate-700 text-amber-300 rounded text-xs">
          元に戻す
        </button>
      </div>
    </div>
  </div>
</template>
