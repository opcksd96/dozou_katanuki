<script setup lang="ts">
import type { RenderTree, RenderMedia } from '../../models/RenderTree';
import type { LanguageCode } from '../../composables/useTimeline';
import Avatar from './Avatar.vue';
import ArticleHeader from './ArticleHeader.vue';
import ArticleBody from './ArticleBody.vue';
import MediaGrid from '../media/MediaGrid.vue';
import ArticleStats from './ArticleStats.vue';
import ArticleInlineThread from './ArticleInlineThread.vue';

const props = defineProps<{
  article: RenderTree;
  targetLang: LanguageCode;
  isFocused?: boolean;
  hasParentLine?: boolean;
  hasChildLine?: boolean;
  isExpanded?: boolean;
  parentArticles?: RenderTree[];
  loadingThread?: boolean;
}>();

const emit = defineEmits<{
  (e: 'clickArticle', id: string): void;
  (e: 'toggleLike', id: string): void;
  (e: 'retryMedia', mediaId: string): void;
  (e: 'clickMedia', media: RenderMedia, list?: RenderMedia[], article?: RenderTree): void;
  (e: 'clickTag', tag: string): void;
  (e: 'clickMention', handle: string): void;
  (e: 'toggleExpandThread', id: string): void;
}>();

const handleCardClick = (e: MouseEvent) => {
  const t = e.target as HTMLElement;
  if (t.closest('button, a, .media-click-target, .media-grid-container, .plyr, .group\\/player, .twitter-inline-thread')) return;
  emit('clickArticle', props.article.id);
};
</script>

<template>
  <article @click="handleCardClick" :class="['twitter-card relative flex items-start gap-3 p-4 bg-slate-950/80 border-b border-slate-800 hover:bg-slate-900/40 transition-all text-left cursor-pointer select-text', isFocused ? 'is-focused' : '']">
    <div v-if="hasParentLine" class="twitter-thread-line-top"></div>
    <div v-if="hasChildLine" class="twitter-thread-line-bottom"></div>

    <div class="twitter-avatar-col relative z-10 flex-shrink-0 pt-0.5">
      <Avatar :avatarUrl="article.author.avatar_url" :handle="article.author.handle" />
    </div>

    <div class="twitter-content-col relative z-10 flex-1 min-w-0 space-y-1.5">
      <ArticleHeader :author="article.author" :createdAt="article.created_at" :sourceUrl="article.source_url" :isPinned="article.is_pinned" />

      <!-- リプライ先バッジ ＆ インライン会話トグル -->
      <div v-if="article.reply_to_handle || article.parent_id" class="twitter-reply-badge flex items-center justify-between">
        <div class="flex items-center gap-1">
          <span>返信先:</span>
          <a @click.stop="emit('clickMention', article.reply_to_handle || '')">@{{ article.reply_to_handle || 'thread' }}</a>
        </div>
        <button v-if="!hasParentLine && article.parent_id" @click.stop="emit('toggleExpandThread', article.id)" class="text-[11px] px-1.5 py-0.5 rounded bg-sky-950/60 hover:bg-sky-900/80 text-sky-400 border border-sky-800/50 transition-colors">
          {{ loadingThread ? '⏳ 読込中...' : (isExpanded ? '▲ 会話を閉じる' : '💬 会話を表示') }}
        </button>
      </div>

      <!-- インライン展開された親ツイート -->
      <div v-if="isExpanded && parentArticles && parentArticles.length > 0" class="space-y-1 my-1">
        <ArticleInlineThread v-for="p in parentArticles" :key="p.id" :article="p" :target-lang="targetLang" @click-article="(id) => emit('clickArticle', id)" @open-media="(m, l, a) => emit('clickMedia', m, l, a)" @click-tag="(t) => emit('clickTag', t)" @click-mention="(h) => emit('clickMention', h)" />
      </div>

      <ArticleBody :content="article.content" :targetLang="targetLang" @clickTag="(t) => emit('clickTag', t)" @clickMention="(h) => emit('clickMention', h)" />
      <MediaGrid :media="article.media" @retry="(mId) => emit('retryMedia', mId)" @clickMedia="(m, l) => emit('clickMedia', m, l, article)" />
      <ArticleStats :metrics="article.metrics" :isLiked="article.is_liked" @toggleLike="emit('toggleLike', article.id)" />
    </div>
  </article>
</template>

