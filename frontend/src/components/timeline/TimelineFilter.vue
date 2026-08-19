<script setup lang="ts">
import type { FilterType } from '../../composables/useTimeline';

defineProps<{
  currentFilter: FilterType;
}>();

const emit = defineEmits<{
  (e: 'filter', filter: FilterType): void;
}>();

const tabs: { label: string; subLabel: string; value: FilterType }[] = [
  { label: 'すべて', subLabel: 'All', value: 'all' },
  { label: 'メディア', subLabel: 'Media', value: 'media' },
  { label: 'リポスト', subLabel: 'Reposts', value: 'reposts' },
  { label: 'ブックマーク', subLabel: 'Saved', value: 'bookmarks' },
];
</script>

<template>
  <div class="flex items-center justify-around border-b border-slate-800 bg-slate-950/95 backdrop-blur">
    <button
      v-for="tab in tabs"
      :key="tab.value"
      @click="emit('filter', tab.value)"
      :class="[
        'flex-1 py-3 text-xs font-semibold transition-all relative text-center cursor-pointer flex flex-col items-center gap-0.5 group',
        currentFilter === tab.value ? 'text-blue-400 font-bold' : 'text-slate-400 hover:text-slate-200'
      ]"
    >
      <span class="text-xs">{{ tab.label }}</span>
      <span class="text-[10px] font-mono opacity-60">{{ tab.subLabel }}</span>
      <div
        v-if="currentFilter === tab.value"
        class="absolute bottom-0 left-1/4 right-1/4 h-0.5 bg-blue-500 rounded-full shadow-sm shadow-blue-500/50"
      />
    </button>
  </div>
</template>
