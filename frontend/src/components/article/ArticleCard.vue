<script setup lang="ts">
import type { RenderTree, RenderMedia } from '../../models/RenderTree';
import type { LanguageCode } from '../../composables/useTimeline';
import Avatar from './Avatar.vue';
import ArticleHeader from './ArticleHeader.vue';
import ArticleBody from './ArticleBody.vue';
import MediaGrid from '../media/MediaGrid.vue';
import ArticleStats from './ArticleStats.vue';

const props = defineProps<{
  article: RenderTree;
  targetLang: LanguageCode;
  isFocused?: boolean;
  hasThreadLine?: boolean;
}>();

const emit = defineEmits<{
  (e: 'clickArticle', id: string): void;
  (e: 'toggleLike', id: string): void;
  (e: 'retryMedia', mediaId: string): void;
  (e: 'clickMedia', media: RenderMedia, list?: RenderMedia[], article?: RenderTree): void;
  (e: 'clickTag', tag: string): void;
  (e: 'clickMention', handle: string): void;
}>();

const handleCardClick = (e: MouseEvent) => {
  const target = e.target as HTMLElement;
  if (target.closest('button') || target.closest('a') || target.closest('.media-click-target')) {
    return;
  }
  emit('clickArticle', props.article.id);
};
</script>

<template>
  <article
    @click="handleCardClick"
    :class="[
      'twitter-card relative flex items-start gap-3 p-4 bg-slate-950/80 border-b border-slate-800 hover:bg-slate-900/40 transition-all text-left cursor-pointer select-text',
      isFocused ? 'is-focused' : ''
    ]"
  >
    <!-- スレッド接続線（親への継続線） -->
    <div v-if="hasThreadLine" class="twitter-thread-line-vertical"></div>

    <!-- 左カラム: アバター (40px) -->
    <div class="twitter-avatar-col relative z-10 flex-shrink-0 pt-0.5">
      <Avatar :avatarUrl="article.author.avatar_url" :handle="article.author.handle" />
    </div>

    <!-- 右カラム: ヘッダー、リプライ先、本文、メディア、アクション -->
    <div class="twitter-content-col relative z-10 flex-1 min-w-0 space-y-1.5">
      <ArticleHeader
        :author="article.author"
        :createdAt="article.created_at"
        :sourceUrl="article.source_url"
        :isPinned="article.is_pinned"
      />

      <!-- リプライ先バッジ -->
      <div v-if="article.reply_to_handle" class="twitter-reply-badge">
        <span>返信先:</span>
        <a @click.stop="emit('clickMention', article.reply_to_handle)">@{{ article.reply_to_handle }}</a>
      </div>

      <ArticleBody
        :content="article.content"
        :targetLang="targetLang"
        @clickTag="(tag) => emit('clickTag', tag)"
        @clickMention="(handle) => emit('clickMention', handle)"
      />
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
