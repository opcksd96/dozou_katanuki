<!-- frontend/src/components/admin/pipeline/PipelineTaskCard.vue (100行以下 - SPEC-PRINCIPLE-001) -->
<script setup lang="ts">
import { BrowserOpenURL } from '../../../../wailsjs/runtime/runtime';

const props = defineProps<{ task: any }>();
const emit = defineEmits<{ (e: 'jumpToMedia', mediaId: string): void }>();

const openInStash = (stashSceneId?: string, stashImageId?: string) => {
  let url = 'http://127.0.0.1:9999/';
  if (stashSceneId) url += `scenes/${stashSceneId}`;
  else if (stashImageId) url += `images/${stashImageId}`;
  try { BrowserOpenURL(url); } catch { window.open(url, '_blank'); }
};
</script>

<template>
  <div class="p-2.5 bg-slate-950 rounded-xl border border-slate-800/80 flex flex-col sm:flex-row sm:items-center justify-between gap-3 hover:border-slate-700 transition">
    <div class="space-y-1 truncate flex-1 min-w-0">
      <div class="flex items-center gap-2">
        <span class="text-slate-200 font-bold truncate text-xs">{{ task.media_id }}</span>
        <span class="px-1.5 py-0.2 rounded text-[9px] font-bold bg-indigo-950 text-indigo-300 border border-indigo-800 shrink-0">
          {{ task.resolution_type || 'orig' }}
        </span>
        <span v-if="task.status === 'COMPLETED'" class="px-1.5 py-0.2 rounded text-[9px] font-bold bg-emerald-950 text-emerald-300 border border-emerald-800 shrink-0">
          ✅ 救出完了
        </span>
        <span v-else-if="task.status === 'REAPED'" class="px-1.5 py-0.2 rounded text-[9px] font-bold bg-slate-800 text-slate-400 border border-slate-700 shrink-0">
          ♻️ 重複解消
        </span>
        <span v-else-if="task.status === 'DEPLETED'" class="px-1.5 py-0.2 rounded text-[9px] font-bold bg-rose-950 text-rose-300 border border-rose-800 shrink-0">
          ⚠️ 枯渇
        </span>
        <span v-else class="px-1.5 py-0.2 rounded text-[9px] font-bold bg-purple-950 text-purple-300 border border-purple-800 shrink-0 animate-pulse">
          ⚡ {{ task.stage || 'THUNDER' }} 稼働中
        </span>
      </div>
      <div class="text-[10px] text-slate-400 truncate font-mono">{{ task.url || task.file_name }}</div>
    </div>

    <!-- メディア管理 ＆ Stash へのワンクリック動線ボタン群 -->
    <div class="flex items-center gap-1.5 shrink-0">
      <button
        @click="emit('jumpToMedia', task.media_id)"
        class="px-2.5 py-1 bg-slate-800 hover:bg-slate-700 text-slate-300 rounded-lg text-xs font-semibold flex items-center gap-1 cursor-pointer transition active:scale-95"
        title="メディア管理ビューでこのメディアを開く"
      >
        🖼️ メディア管理
      </button>
      <button
        @click="openInStash(task.stash_scene_id, task.stash_image_id)"
        class="px-2.5 py-1 bg-indigo-950 hover:bg-indigo-900 border border-indigo-800 text-indigo-200 rounded-lg text-xs font-semibold flex items-center gap-1 cursor-pointer transition active:scale-95"
        title="Stash Web UI で開く"
      >
        🎬 Stash
      </button>
    </div>
  </div>
</template>
