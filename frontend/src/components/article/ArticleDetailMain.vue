<!-- frontend/src/components/article/ArticleDetailMain.vue (100行以下) -->
<script setup lang="ts">
import type { RenderTree, RenderMedia } from '../../models/RenderTree';
import type { LanguageCode } from '../../composables/useTimeline';
import Avatar from './Avatar.vue';
import MediaGrid from '../media/MediaGrid.vue';
import ArticleStats from './ArticleStats.vue';
import ArticleBody from './ArticleBody.vue';

defineProps<{
  article: RenderTree;
  targetLang: LanguageCode;
}>();

defineEmits<{
  (e: 'toggleLike', id: string): void;
  (e: 'retryMedia', mediaId: string): void;
  (e: 'clickMedia', media: RenderMedia, list?: RenderMedia[], article?: RenderTree): void;
  (e: 'clickTag', tag: string): void;
  (e: 'clickMention', handle: string): void;
}>();
</script>

<template>
  <div class="p-5 bg-slate-950 border-b border-slate-800 space-y-4">
    <div class="flex items-center gap-3">
      <Avatar :author="article.author" class="w-12 h-12 flex-shrink-0" />
      <div>
        <div class="text-base font-bold text-white">{{ article.author.display_name }}</div>
        <div class="text-xs text-slate-400 font-mono">@{{ article.author.handle }}</div>
      </div>
    </div>

    <ArticleBody
      :content="article.content"
      :target-lang="targetLang"
      class="text-sm text-slate-100 leading-relaxed select-text"
      @click-tag="(t) => $emit('clickTag', t)"
      @click-mention="(m) => $emit('clickMention', m)"
    />

    <MediaGrid
      v-if="article.media?.length"
      :media="article.media"
      @clickMedia="(m, l) => $emit('clickMedia', m, l, article)"
      @retry="(id) => $emit('retryMedia', id)"
    />

    <div class="text-xs text-slate-500 font-mono pt-2 border-t border-slate-900 flex items-center justify-between">
      <span>{{ new Date(article.created_at).toLocaleString() }}</span>
      <span v-if="article.via">via {{ article.via }}</span>
    </div>

    <ArticleStats :metrics="article.metrics" :is-liked="article.is_liked" @toggle-like="$emit('toggleLike', article.id)" @toggleLike="$emit('toggleLike', article.id)" />
  </div>
</template>
