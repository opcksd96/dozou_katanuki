<script setup lang="ts">
import { computed } from 'vue';
import type { SupportedLang } from '../../composables/useTimeline';

const props = defineProps<{
  content: { original: string; ja?: string; en?: string; zh?: string };
  currentLang: SupportedLang;
}>();

const text = computed(() => {
  if (props.currentLang === 'original') return props.content.original;
  return props.content[props.currentLang] || props.content.original;
});
</script>

<template>
  <div class="my-2 text-slate-200 text-sm leading-relaxed whitespace-pre-wrap">
    <p>{{ text }}</p>
    <div v-if="currentLang !== 'original' && content[currentLang]" class="mt-1 text-[11px] text-emerald-400 font-mono">
      [翻訳: {{ currentLang.toUpperCase() }}]
    </div>
  </div>
</template>
