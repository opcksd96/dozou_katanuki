<!-- frontend/src/components/article/ArticleDetailView.vue (100行以下) -->
<script setup lang="ts">
import { computed } from 'vue';
import type { RenderTree, RenderMedia } from '../../models/RenderTree';
import type { LanguageCode } from '../../composables/useTimeline';
import ArticleDetailNav from './ArticleDetailNav.vue';
import ArticleDetailItem from './ArticleDetailItem.vue';
import ArticleDetailMain from './ArticleDetailMain.vue';

const props = defineProps<{
  article: RenderTree;
  thread: RenderTree[];
  targetLang: LanguageCode;
  loading?: boolean;
  focusedArticleId?: string;
}>();

const emit = defineEmits<{
  (e: 'back'): void;
  (e: 'selectArticle', id: string): void;
  (e: 'toggleLike', id: string): void;
  (e: 'retryMedia', mediaId: string): void;
  (e: 'clickMedia', media: RenderMedia, list?: RenderMedia[], article?: RenderTree): void;
  (e: 'clickTag', tag: string): void;
  (e: 'clickMention', handle: string): void;
}>();

const parentChain = computed(() => {
  if (!props.thread) return [];
  const mainTime = new Date(props.article.created_at).getTime();
  return props.thread.filter((t) => t.id !== props.article.id && new Date(t.created_at).getTime() < mainTime)
    .sort((a, b) => new Date(a.created_at).getTime() - new Date(b.created_at).getTime());
});

const childReplies = computed(() => {
  if (!props.thread) return [];
  const mainTime = new Date(props.article.created_at).getTime();
  return props.thread.filter((t) => t.id !== props.article.id && new Date(t.created_at).getTime() >= mainTime)
    .sort((a, b) => new Date(a.created_at).getTime() - new Date(b.created_at).getTime());
});
</script>

<template>
  <div class="w-full bg-slate-950 min-h-screen text-left">
    <ArticleDetailNav :article="article" @back="$emit('back')" />

    <div>
      <!-- 親スレッド -->
      <div v-if="parentChain.length > 0" class="border-b border-slate-800">
        <ArticleDetailItem
          v-for="parent in parentChain"
          :key="parent.id"
          :item="parent"
          :target-lang="targetLang"
          :is-focused="focusedArticleId === parent.id"
          @select-article="(id) => $emit('selectArticle', id)"
          @click-media="(m, l, a) => $emit('clickMedia', m, l, a)"
          @click-tag="(t) => $emit('clickTag', t)"
          @click-mention="(h) => $emit('clickMention', h)"
        />
      </div>

      <!-- メイン記事 -->
      <ArticleDetailMain
        :article="article"
        :target-lang="targetLang"
        @toggle-like="(id) => $emit('toggleLike', id)"
        @retry-media="(id) => $emit('retryMedia', id)"
        @click-media="(m, l, a) => $emit('clickMedia', m, l, a)"
        @click-tag="(t) => $emit('clickTag', t)"
        @click-mention="(h) => $emit('clickMention', h)"
      />

      <!-- 子スレッド -->
      <div v-if="childReplies.length > 0">
        <ArticleDetailItem
          v-for="reply in childReplies"
          :key="reply.id"
          :item="reply"
          :target-lang="targetLang"
          :is-focused="focusedArticleId === reply.id"
          @select-article="(id) => $emit('selectArticle', id)"
          @click-media="(m, l, a) => $emit('clickMedia', m, l, a)"
          @click-tag="(t) => $emit('clickTag', t)"
          @click-mention="(h) => $emit('clickMention', h)"
        />
      </div>
    </div>
  </div>
</template>
