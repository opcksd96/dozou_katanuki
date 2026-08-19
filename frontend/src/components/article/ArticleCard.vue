<script setup lang="ts">
import type { RenderTree } from '../../models/RenderTree';
import type { SupportedLang } from '../../composables/useTimeline';
import ArticleHeader from './ArticleHeader.vue';
import ArticleBody from './ArticleBody.vue';
import MediaGrid from '../media/MediaGrid.vue';
import ArticleStats from './ArticleStats.vue';

defineProps<{
  article: RenderTree;
  currentLang: SupportedLang;
}>();

const emit = defineEmits<{
  (e: 'toggleLike', id: string): void;
  (e: 'retryMedia', mediaId: string): void;
}>();
</script>

<template>
  <article class="p-4 bg-slate-900 border border-slate-800 rounded-xl mb-4 hover:border-slate-700/80 transition-colors">
    <ArticleHeader
      :author="article.author"
      :createdAt="article.created_at"
      :sourceUrl="article.source_url"
      :isPinned="article.is_pinned"
    />
    <ArticleBody :content="article.content" :currentLang="currentLang" />
    <MediaGrid :media="article.media" @retry="(mediaId) => emit('retryMedia', mediaId)" />
    <ArticleStats :metrics="article.metrics" :isLiked="article.is_liked" @toggleLike="emit('toggleLike', article.id)" />
  </article>
</template>
