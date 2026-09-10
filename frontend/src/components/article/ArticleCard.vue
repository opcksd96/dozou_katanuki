<script setup lang="ts">
import { computed } from 'vue';
import type { RenderTree, RenderMedia } from '../../models/RenderTree';
import type { LanguageCode } from '../../composables/useTimeline';
import Avatar from './Avatar.vue';
import ArticleHeader from './ArticleHeader.vue';
import ArticleBody from './ArticleBody.vue';
import MediaGrid from '../media/MediaGrid.vue';
import ArticleStats from './ArticleStats.vue';
import ArticleInlineThread from './ArticleInlineThread.vue';

const props = defineProps<{
  article: RenderTree; targetLang: LanguageCode; isFocused?: boolean;
  hasParentLine?: boolean; hasChildLine?: boolean; isExpanded?: boolean;
  parentArticles?: RenderTree[]; loadingThread?: boolean;
}>();

const emit = defineEmits<{
  (e: 'clickArticle', id: string): void; (e: 'toggleLike', id: string): void;
  (e: 'retryMedia', mediaId: string): void; (e: 'clickMedia', media: RenderMedia, list?: RenderMedia[], article?: RenderTree): void;
  (e: 'clickTag', tag: string): void; (e: 'clickMention', handle: string): void; (e: 'toggleExpandThread', id: string): void;
}>();

const targetPost = computed(() => props.article.retweeted_article || props.article);
const handleCardClick = (e: MouseEvent) => {
  const t = e.target as HTMLElement;
  if (t.closest('button, a, .media-click-target, .media-grid-container, .plyr, .group\\/player, .twitter-inline-thread')) return;
  emit('clickArticle', props.article.id);
};
</script>

<template>
  <article @click="handleCardClick" :class="['twitter-card relative flex flex-col gap-1 p-4 bg-slate-950/80 border-b border-slate-800 hover:bg-slate-900/40 transition-all text-left cursor-pointer select-text', isFocused ? 'is-focused' : '']">
    <div v-if="hasParentLine" class="twitter-thread-line-top"></div>
    <div v-if="hasChildLine" class="twitter-thread-line-bottom"></div>

    <!-- リポストヘッダー -->
    <div v-if="article.is_repost" class="flex items-center gap-1.5 text-xs text-slate-400 font-medium pl-8 pb-1">
      <span class="text-emerald-400 font-bold">🔁</span>
      <span>{{ article.author.display_name || article.author.handle }} さんがリポスト</span>
    </div>

    <div class="flex items-start gap-3 w-full">
      <div class="twitter-avatar-col relative z-10 flex-shrink-0 pt-0.5">
        <Avatar :avatarUrl="targetPost.author.avatar_url" :handle="targetPost.author.handle" />
      </div>

      <div class="twitter-content-col relative z-10 flex-1 min-w-0 space-y-1.5">
        <ArticleHeader :author="targetPost.author" :createdAt="targetPost.created_at" :sourceUrl="targetPost.source_url" :isPinned="targetPost.is_pinned" />

        <div v-if="targetPost.reply_to_handle || targetPost.parent_id" class="twitter-reply-badge flex items-center justify-between">
          <div class="flex items-center gap-1">
            <span>返信先:</span>
            <a @click.stop="emit('clickMention', targetPost.reply_to_handle || '')">@{{ targetPost.reply_to_handle || 'thread' }}</a>
          </div>
          <button v-if="!hasParentLine && targetPost.parent_id" @click.stop="emit('toggleExpandThread', targetPost.id)" class="text-[11px] px-1.5 py-0.5 rounded bg-sky-950/60 hover:bg-sky-900/80 text-sky-400 border border-sky-800/50 transition-colors">
            {{ loadingThread ? '⏳ 読込中...' : (isExpanded ? '▲ 会話を閉じる' : '💬 会話を表示') }}
          </button>
        </div>

        <div v-if="isExpanded && parentArticles && parentArticles.length > 0" class="space-y-1 my-1">
          <ArticleInlineThread v-for="p in parentArticles" :key="p.id" :article="p" :target-lang="targetLang" @click-article="(id) => emit('clickArticle', id)" @open-media="(m, l, a) => emit('clickMedia', m, l, a)" @click-tag="(t) => emit('clickTag', t)" @click-mention="(h) => emit('clickMention', h)" />
        </div>

        <ArticleBody :content="targetPost.content" :targetLang="targetLang" @clickTag="(t) => emit('clickTag', t)" @clickMention="(h) => emit('clickMention', h)" />
        <MediaGrid :media="targetPost.media" @retry="(mId) => emit('retryMedia', mId)" @clickMedia="(m, l) => emit('clickMedia', m, l, targetPost)" />
        <ArticleStats :metrics="targetPost.metrics" :isLiked="article.is_liked" @toggleLike="emit('toggleLike', article.id)" />
      </div>
    </div>
  </article>
</template>

