<script setup lang="ts">
import type { RenderTree, RenderMedia } from '../../models/RenderTree';
import type { LanguageCode } from '../../composables/useTimeline';
import Avatar from './Avatar.vue';
import ArticleHeader from './ArticleHeader.vue';
import ArticleBody from './ArticleBody.vue';
import MediaGrid from '../media/MediaGrid.vue';
import ArticleStats from './ArticleStats.vue';

defineProps<{
  article: RenderTree;
  targetLang: LanguageCode;
}>();

const emit = defineEmits<{
  (e: 'toggleLike', id: string): void;
  (e: 'retryMedia', mediaId: string): void;
  (e: 'clickMedia', media: RenderMedia, list?: RenderMedia[], article?: RenderTree): void;
}>();
</script>

<template>
  <article class="twitter-card flex items-start gap-3 p-4 bg-slate-950/80 border-b border-slate-800 hover:bg-slate-900/40 transition-colors text-left">
    <!-- 左カラム: アバター (40px) -->
    <div class="twitter-avatar-col flex-shrink-0 pt-0.5">
      <Avatar :avatarUrl="article.author.avatar_url" :handle="article.author.handle" />
    </div>

    <!-- 右カラム: ヘッダー、本文、メディア、アクション -->
    <div class="twitter-content-col flex-1 min-w-0 space-y-1.5">
      <ArticleHeader
        :author="article.author"
        :createdAt="article.created_at"
        :sourceUrl="article.source_url"
        :isPinned="article.is_pinned"
      />
      <ArticleBody :content="article.content" :targetLang="targetLang" />
      <MediaGrid
        :media="article.media"
        @retry="(mediaId) => emit('retryMedia', mediaId)"
        @clickMedia="(media, list) => emit('clickMedia', media, list, article)"
      />
      <ArticleStats
        :metrics="article.metrics"
        :isLiked="article.is_liked"
        @toggleLike="emit('toggleLike', article.id)"
      />
    </div>
  </article>
</template>
