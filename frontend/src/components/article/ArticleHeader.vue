<script setup lang="ts">
import type { RenderAuthor } from '../../models/RenderTree';
import { formatDate } from '../../utils/formatters';
import Avatar from './Avatar.vue';

defineProps<{
  author: RenderAuthor;
  createdAt: string;
  sourceUrl: string;
  isPinned?: boolean;
}>();
</script>

<template>
  <div class="flex items-start justify-between gap-3 mb-2">
    <div class="flex items-center gap-3">
      <Avatar :avatarUrl="author.avatar_url" :handle="author.handle" />
      <div>
        <div class="flex items-center gap-2">
          <span class="font-bold text-slate-100 leading-tight">{{ author.display_name }}</span>
          <span v-if="isPinned" class="text-[10px] bg-amber-500/20 text-amber-300 px-1.5 py-0.5 rounded font-mono">PINNED</span>
        </div>
        <div class="text-xs text-slate-400">@{{ author.handle }} · {{ formatDate(createdAt) }}</div>
      </div>
    </div>
    <a :href="sourceUrl" target="_blank" rel="noopener noreferrer" class="text-xs text-blue-400 hover:underline">原本リンク</a>
  </div>
</template>
