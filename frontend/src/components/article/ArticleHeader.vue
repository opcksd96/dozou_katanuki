<script setup lang="ts">
import type { RenderAuthor } from '../../models/RenderTree';
import { formatDate } from '../../utils/formatters';

defineProps<{
  author: RenderAuthor;
  createdAt: string;
  sourceUrl: string;
  isPinned?: boolean;
}>();
</script>

<template>
  <div class="flex items-center justify-between gap-2 text-left">
    <div class="flex items-center gap-1.5 min-w-0 truncate">
      <span class="font-bold text-slate-100 text-sm hover:underline cursor-pointer truncate leading-snug">
        {{ author.display_name }}
      </span>
      <span class="text-xs text-slate-400 font-mono truncate leading-snug">
        @{{ author.handle }}
      </span>
      <span class="text-xs text-slate-500 leading-snug">·</span>
      <span class="text-xs text-slate-400 font-mono whitespace-nowrap leading-snug">
        {{ formatDate(createdAt) }}
      </span>
      <span v-if="isPinned" class="text-[9px] bg-amber-500/20 text-amber-300 px-1 py-0.2 rounded font-mono ml-1">
        PINNED
      </span>
    </div>

    <a
      v-if="sourceUrl"
      :href="sourceUrl"
      target="_blank"
      rel="noopener noreferrer"
      class="text-[11px] text-blue-400 hover:text-blue-300 font-mono flex-shrink-0 hover:underline"
      title="原本（Wayback/CDN）を開く"
    >
      原本 ↗
    </a>
  </div>
</template>
