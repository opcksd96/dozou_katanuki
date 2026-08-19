<script setup lang="ts">
import { computed } from 'vue';
import type { TranslatedText } from '../types/article';
import type { SupportedLang } from '../composables/useTimeline';

const props = defineProps<{
  content: TranslatedText;
  currentLang: SupportedLang;
}>();

const displayContent = computed(() => {
  if (props.currentLang === 'original') return props.content.original;
  return props.content[props.currentLang] || props.content.original;
});
</script>

<template>
  <div class="my-2 text-slate-200 text-sm leading-relaxed whitespace-pre-wrap">
    <p>{{ displayContent }}</p>
    <div
      v-if="currentLang !== 'original' && content[currentLang]"
      class="mt-1 text-[11px] text-emerald-400 font-mono"
    >
      [翻訳: {{ currentLang.toUpperCase() }}]
    </div>
  </div>
</template>
