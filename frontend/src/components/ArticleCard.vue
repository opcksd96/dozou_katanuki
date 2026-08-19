<script setup lang="ts">
import type { ArticleRenderTree } from '../types/article';
import type { SupportedLang } from '../composables/useTimeline';
import ArticleAuthor from './ArticleAuthor.vue';
import ArticleBody from './ArticleBody.vue';
import MediaGrid from './MediaGrid.vue';

defineProps<{
  article: ArticleRenderTree;
  currentLang: SupportedLang;
}>();
</script>

<template>
  <article class="p-4 bg-slate-900 border border-slate-800 rounded-xl mb-4 hover:border-slate-700 transition-colors">
    <ArticleAuthor
      :author="article.author"
      :createdAt="article.createdAt"
      :originalUrl="article.originalUrl"
    />
    <ArticleBody
      :content="article.content"
      :currentLang="currentLang"
    />
    <MediaGrid :media="article.media" />
  </article>
</template>
