<script setup lang="ts">
import type { FilterType } from '../../composables/useTimeline';

defineProps<{
  currentFilter: FilterType;
}>();

const emit = defineEmits<{
  (e: 'filter', filter: FilterType): void;
}>();

const tabs: { label: string; value: FilterType }[] = [
  { label: 'All', value: 'all' },
  { label: 'Media', value: 'media' },
  { label: 'Reposts', value: 'reposts' },
  { label: 'Bookmarks', value: 'bookmarks' },
];
</script>

<template>
  <div class="flex items-center justify-around border-b border-slate-800 bg-slate-950/80 backdrop-blur sticky top-0 z-10 mb-4">
    <button
      v-for="tab in tabs"
      :key="tab.value"
      @click="emit('filter', tab.value)"
      :class="[
        'flex-1 py-3 text-xs font-semibold transition-colors relative text-center',
        currentFilter === tab.value ? 'text-blue-400' : 'text-slate-400 hover:text-slate-200'
      ]"
    >
      {{ tab.label }}
      <div
        v-if="currentFilter === tab.value"
        class="absolute bottom-0 left-1/4 right-1/4 h-0.5 bg-blue-500 rounded-full"
      />
    </button>
  </div>
</template>
