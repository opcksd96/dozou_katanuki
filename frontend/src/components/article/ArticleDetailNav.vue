<!-- frontend/src/components/article/ArticleDetailNav.vue (100行以下) -->
<script setup lang="ts">
import { ref } from 'vue';
import type { RenderTree } from '../../models/RenderTree';

const props = defineProps<{ article: RenderTree }>();
defineEmits<{ (e: 'back'): void }>();

const copySuccess = ref(false);
const copyArticleUrl = async () => {
  const url = props.article.source_url || `https://twitter.com/${props.article.author.handle}/status/${props.article.id}`;
  try {
    await navigator.clipboard.writeText(url);
    copySuccess.value = true;
    setTimeout(() => { copySuccess.value = false; }, 2000);
  } catch (_) {}
};
</script>

<template>
  <div class="sticky top-0 z-30 bg-slate-950/85 backdrop-blur border-b border-slate-800 px-4 py-3 flex items-center justify-between">
    <div class="flex items-center gap-3">
      <button @click="$emit('back')" class="p-2 rounded-full hover:bg-slate-800 text-slate-300 hover:text-white transition-colors">
        <span class="text-base font-bold">←</span>
      </button>
      <div>
        <h2 class="text-sm font-bold text-white">ポスト</h2>
        <p class="text-[10px] text-slate-400 font-mono">ID: {{ article.id }}</p>
      </div>
    </div>
    <div class="flex items-center gap-2">
      <button @click="copyArticleUrl" class="px-2.5 py-1 rounded bg-slate-900 hover:bg-slate-800 text-slate-300 border border-slate-700 text-xs font-mono">
        {{ copySuccess ? '✓ Copied' : '🔗 Link' }}
      </button>
      <a v-if="article.source_url" :href="article.source_url" target="_blank" rel="noopener noreferrer" class="px-2.5 py-1 rounded bg-blue-900/40 hover:bg-blue-800/60 text-blue-300 border border-blue-700/50 text-xs font-mono">
        🏛️ Wayback
      </a>
    </div>
  </div>
</template>
