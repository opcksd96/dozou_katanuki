<script setup lang="ts">
import type { RenderAuthor } from '../../models/RenderTree';

defineProps<{
  accounts: RenderAuthor[];
  selectedId: string;
}>();

const emit = defineEmits<{
  (e: 'select', id: string): void;
}>();
</script>

<template>
  <div class="flex flex-wrap items-center gap-2 pt-1 pb-1">
    <button
      @click="emit('select', 'all')"
      :class="[
        'flex items-center gap-2 px-3 py-1.5 rounded-full text-xs font-medium whitespace-nowrap transition-all border cursor-pointer shrink-0 shadow-sm',
        selectedId === 'all'
          ? 'bg-blue-600 border-blue-500 text-white shadow-blue-500/20'
          : 'bg-slate-900/90 border-slate-800 text-slate-300 hover:text-white hover:bg-slate-800 hover:border-slate-700'
      ]"
    >
      <span>🌐 全てのアカウント</span>
    </button>

    <button
      v-for="acc in accounts"
      :key="acc.numeric_id"
      @click="emit('select', acc.numeric_id)"
      :class="[
        'flex items-center gap-2 px-3 py-1.5 rounded-full text-xs font-medium whitespace-nowrap transition-all border cursor-pointer shrink-0 shadow-sm',
        selectedId === acc.numeric_id
          ? 'bg-blue-600 border-blue-500 text-white shadow-blue-500/20'
          : 'bg-slate-900/90 border-slate-800 text-slate-300 hover:text-white hover:bg-slate-800 hover:border-slate-700'
      ]"
    >
      <img
        v-if="acc.avatar_url"
        :src="acc.avatar_url"
        :alt="acc.handle"
        class="w-4 h-4 rounded-full object-cover bg-slate-800"
        @error="($event.target as HTMLElement).style.display = 'none'"
      />
      <span>@{{ acc.handle }}</span>
    </button>
  </div>
</template>
