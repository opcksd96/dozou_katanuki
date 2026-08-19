<script setup lang="ts">
import type { Author } from '../types/article';
import { formatDate } from '../utils/formatters';

defineProps<{
  author: Author;
  createdAt: string;
  originalUrl: string;
}>();
</script>

<template>
  <div class="flex items-center justify-between gap-3 mb-2">
    <div class="flex items-center gap-3">
      <img
        :src="author.avatarUrl"
        :alt="author.username"
        class="w-10 h-10 rounded-full border border-slate-700 object-cover bg-slate-800"
      />
      <div>
        <div class="font-bold text-slate-100 leading-tight">
          {{ author.displayName }}
        </div>
        <div class="text-xs text-slate-400">
          @{{ author.username }} · {{ formatDate(createdAt) }}
        </div>
      </div>
    </div>
    <a
      :href="originalUrl"
      target="_blank"
      rel="noopener noreferrer"
      class="text-xs text-blue-400 hover:underline"
    >
      原本リンク
    </a>
  </div>
</template>
