<!-- frontend/src/components/article/ArticleBody.vue (100行以下) -->
<script setup lang="ts">
import { ref, computed } from 'vue';
import type { LanguageCode } from '../../composables/useTimeline';
import { decorateText } from '../../utils/decorator';
import { BrowserOpenURL } from '../../../wailsjs/runtime/runtime';

const props = defineProps<{
  content: { original: string; ja?: string; en?: string; zh?: string };
  targetLang: LanguageCode;
}>();

const emit = defineEmits<{ (e: 'clickTag', tag: string): void; (e: 'clickMention', handle: string): void }>();

const isTranslated = ref(false);
const langLabel = computed(() => {
  const map: Record<string, string> = { ja: '日本語', en: '英語', zh: '中国語' };
  return map[props.targetLang] || props.targetLang.toUpperCase();
});

const hasTranslation = computed(() => {
  const trans = props.content[props.targetLang as 'ja' | 'en' | 'zh'];
  return Boolean(trans && trans.trim() !== '' && trans !== props.content.original);
});

const rawDisplayText = computed(() => {
  if (isTranslated.value) {
    const trans = props.content[props.targetLang as 'ja' | 'en' | 'zh'];
    return (trans && trans.trim() !== '') ? trans : props.content.original || '';
  }
  return props.content.original || '';
});

const decoratedHtml = computed(() => decorateText(rawDisplayText.value));

const handleBodyClick = (e: MouseEvent) => {
  const link = (e.target as HTMLElement).closest('a');
  if (!link) return;
  e.preventDefault(); e.stopPropagation();

  if (link.classList.contains('hashtag-link')) {
    const tag = link.getAttribute('data-tag') || link.innerText.replace(/^#/, '');
    if (tag) emit('clickTag', tag);
  } else if (link.classList.contains('mention-link')) {
    const mention = link.getAttribute('data-mention') || link.innerText.replace(/^@/, '');
    if (mention) emit('clickMention', mention);
  } else if (link.classList.contains('external-link')) {
    const url = link.getAttribute('data-url') || link.getAttribute('href');
    if (url) { try { BrowserOpenURL(url); } catch { window.open(url, '_blank', 'noopener,noreferrer'); } }
  }
};
</script>

<template>
  <div class="my-2.5 text-slate-200 text-sm leading-relaxed">
    <div @click="handleBodyClick" class="whitespace-pre-line break-words leading-relaxed select-text" v-html="decoratedHtml"></div>
    <div class="mt-2 flex items-center gap-2 text-xs font-mono select-none">
      <button @click="isTranslated = !isTranslated" type="button" class="text-blue-400 hover:text-blue-300 flex items-center gap-1.5 py-0.5">
        <span>🌐</span>
        <span v-if="!isTranslated" class="hover:underline">{{ langLabel }} に翻訳</span>
        <span v-else class="flex items-center gap-1.5">
          <span class="text-emerald-400 bg-emerald-950/60 border border-emerald-800/60 px-1.5 py-0.5 rounded text-[10px]">{{ langLabel }} 翻訳</span>
          <span class="text-slate-400 hover:text-slate-200 hover:underline">原文を表示</span>
        </span>
      </button>
      <span v-if="isTranslated && !hasTranslation" class="text-[10px] text-amber-400/80 font-sans">（※ 翻訳未登録のため原文）</span>
    </div>
  </div>
</template>
