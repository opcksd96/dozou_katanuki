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
  <div class="flex items-center gap-2 overflow-x-auto pb-2 scrollbar-none">
    <button
      @click="emit('select', 'all')"
      :class="[
        'flex items-center gap-2 px-3 py-1.5 rounded-full text-xs font-medium whitespace-nowrap transition-all border',
        selectedId === 'all'
          ? 'bg-blue-600 border-blue-500 text-white shadow-sm'
          : 'bg-slate-900 border-slate-800 text-slate-400 hover:text-slate-200 hover:border-slate-700'
      ]"
    >
      <span>🌐 All Accounts</span>
    </button>

    <button
      v-for="acc in accounts"
      :key="acc.numeric_id"
      @click="emit('select', acc.numeric_id)"
      :class="[
        'flex items-center gap-2 px-3 py-1.5 rounded-full text-xs font-medium whitespace-nowrap transition-all border',
        selectedId === acc.numeric_id
          ? 'bg-blue-600 border-blue-500 text-white shadow-sm'
          : 'bg-slate-900 border-slate-800 text-slate-400 hover:text-slate-200 hover:border-slate-700'
      ]"
    >
      <img
        :src="acc.avatar_url"
        :alt="acc.handle"
        class="w-4 h-4 rounded-full object-cover bg-slate-800"
        @error="($event.target as HTMLElement).style.display = 'none'"
      />
      <span>@{{ acc.handle }}</span>
    </button>
  </div>
</template>
