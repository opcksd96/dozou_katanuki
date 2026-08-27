<!-- frontend/src/components/article/ArticleDetailNav.vue (100行以下) -->
<script setup lang="ts">
import { ref, computed } from 'vue';
import { BrowserOpenURL } from '../../../wailsjs/runtime/runtime';
import type { RenderTree } from '../../models/RenderTree';
import { useToast } from '../../composables/useToast';

const props = defineProps<{ article: RenderTree }>();
defineEmits<{ (e: 'back'): void }>();

const { addToast } = useToast();
const copySuccess = ref(false);

const originalTweetUrl = computed(() => `https://twitter.com/${props.article.author.handle || 'i'}/status/${props.article.id}`);
const waybackUrl = computed(() => props.article.source_url || `https://web.archive.org/web/20250000000000/https://twitter.com/${props.article.author.handle}/status/${props.article.id}`);

const copyArticleUrl = async () => {
  try {
    await navigator.clipboard.writeText(originalTweetUrl.value);
    copySuccess.value = true;
    addToast('📋 元ツイートURLをコピーしました', 'info', 2500);
    setTimeout(() => { copySuccess.value = false; }, 2000);
  } catch (_) {}
};

const openOriginal = () => {
  try { BrowserOpenURL(originalTweetUrl.value); } catch { window.open(originalTweetUrl.value, '_blank', 'noopener,noreferrer'); }
};

const openWayback = () => {
  if (!waybackUrl.value) return;
  try { BrowserOpenURL(waybackUrl.value); } catch { window.open(waybackUrl.value, '_blank', 'noopener,noreferrer'); }
};
</script>

<template>
  <div class="sticky top-0 z-30 bg-slate-950/90 backdrop-blur-xl border-b border-slate-800/80 px-4 py-2.5 flex items-center justify-between select-none">
    <div class="flex items-center gap-4 min-w-0">
      <button @click="$emit('back')" class="w-8 h-8 rounded-full hover:bg-slate-800/80 text-slate-200 hover:text-white flex items-center justify-center transition-colors active:scale-95 cursor-pointer" title="タイムラインに戻る (Esc)">
        <span class="text-lg font-bold">←</span>
      </button>
      <div class="flex items-center gap-2 min-w-0">
        <h2 class="text-base font-bold text-white tracking-wide">ポスト</h2>
      </div>
    </div>

    <!-- 右側アクション (意図が明確で洗練されたスマートバッジ) -->
    <div class="flex items-center gap-1.5 shrink-0">
      <button @click="copyArticleUrl" class="px-2.5 py-1 rounded-full bg-slate-900/90 hover:bg-slate-800 text-slate-300 hover:text-white border border-slate-700/80 text-xs font-semibold flex items-center gap-1 transition-all active:scale-95 cursor-pointer" title="元ツイートURLをコピー">
        <span>{{ copySuccess ? '✓' : '🔗' }}</span>
        <span>{{ copySuccess ? 'コピー済' : '原本' }}</span>
      </button>
      <button @click="openWayback" class="px-2.5 py-1 rounded-full bg-blue-950/80 hover:bg-blue-900/80 text-blue-200 hover:text-white border border-blue-700/60 text-xs font-semibold flex items-center gap-1 transition-all active:scale-95 cursor-pointer" title="Wayback 魚拓アーカイブを開く">
        <span>🏛️</span>
        <span>魚拓</span>
      </button>
    </div>
  </div>
</template>
