<script setup lang="ts">
import { ref, computed } from 'vue';
import type { LanguageCode } from '../../composables/useTimeline';
import { decorateText } from '../../utils/decorator';
import { BrowserOpenURL } from '../../../wailsjs/runtime/runtime';

const props = defineProps<{
  content: { original: string; ja?: string; en?: string; zh?: string };
  targetLang: LanguageCode; // アプリ設定 (config.json) のシステム言語
}>();

const emit = defineEmits<{
  (e: 'clickTag', tag: string): void;
  (e: 'clickMention', handle: string): void;
}>();

// 投稿個別の翻訳状態（デフォルトは原文表示）
const isTranslated = ref(false);

const langLabel = computed(() => {
  switch (props.targetLang) {
    case 'ja': return '日本語';
    case 'en': return '英語';
    case 'zh': return '中国語';
    default: return props.targetLang.toUpperCase();
  }
});

// アプリ設定言語での翻訳データが存在するか
const hasTranslation = computed(() => {
  const trans = props.content[props.targetLang as 'ja' | 'en' | 'zh'];
  return Boolean(trans && trans.trim() !== '' && trans !== props.content.original);
});

// 表示するテキスト
const rawDisplayText = computed(() => {
  if (isTranslated.value) {
    const trans = props.content[props.targetLang as 'ja' | 'en' | 'zh'];
    return trans && trans.trim() !== '' ? trans : props.content.original || '';
  }
  return props.content.original || '';
});

// 装飾済みHTML (DOMリンク化 & 改行展開)
const decoratedHtml = computed(() => {
  return decorateText(rawDisplayText.value);
});

const toggleTranslate = () => {
  isTranslated.value = !isTranslated.value;
};

// 本文クリック時のリンク・タグ・メンションのハンドリング
const handleBodyClick = (e: MouseEvent) => {
  const target = e.target as HTMLElement;
  const link = target.closest('a');
  if (!link) return;

  if (link.classList.contains('hashtag-link')) {
    e.preventDefault();
    e.stopPropagation();
    const tag = link.getAttribute('data-tag') || link.innerText.replace(/^#/, '');
    if (tag) emit('clickTag', tag);
  } else if (link.classList.contains('mention-link')) {
    e.preventDefault();
    e.stopPropagation();
    const mention = link.getAttribute('data-mention') || link.innerText.replace(/^@/, '');
    if (mention) emit('clickMention', mention);
  } else if (link.classList.contains('external-link')) {
    e.preventDefault();
    e.stopPropagation();
    const url = link.getAttribute('data-url') || link.getAttribute('href');
    if (url) {
      try {
        BrowserOpenURL(url);
      } catch {
        window.open(url, '_blank', 'noopener,noreferrer');
      }
    }
  }
};
</script>

<template>
  <div class="my-2.5 text-slate-200 text-sm leading-relaxed">
    <!-- 装飾済み本文 (DOMリンク化 ＆ 改行展開 ＆ クリック委譲) -->
    <div
      @click="handleBodyClick"
      class="whitespace-pre-line break-words leading-relaxed select-text"
      v-html="decoratedHtml"
    ></div>

    <!-- X (Twitter) ライクな翻訳切り替えリンク -->
    <div class="mt-2 flex items-center gap-2 text-xs font-mono select-none">
      <button
        @click="toggleTranslate"
        type="button"
        class="text-blue-400 hover:text-blue-300 flex items-center gap-1.5 transition-colors cursor-pointer group py-0.5"
      >
        <span class="text-xs group-hover:scale-110 transition-transform">🌐</span>
        <span v-if="!isTranslated" class="hover:underline">
          {{ langLabel }} に翻訳
        </span>
        <span v-else class="flex items-center gap-1.5">
          <span class="text-emerald-400 font-semibold bg-emerald-950/60 border border-emerald-800/60 px-1.5 py-0.5 rounded text-[10px]">
            {{ langLabel }} 翻訳
          </span>
          <span class="text-slate-400 hover:text-slate-200 hover:underline">原文を表示</span>
        </span>
      </button>

      <!-- 翻訳データ未登録時の控えめな注記 -->
      <span
        v-if="isTranslated && !hasTranslation"
        class="text-[10px] text-amber-400/80 font-sans"
      >
        （※ {{ langLabel }} 翻訳データ未登録のため原文を表示中）
      </span>
    </div>
  </div>
</template>
