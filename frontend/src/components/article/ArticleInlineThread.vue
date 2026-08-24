<!-- frontend/src/components/article/ArticleInlineThread.vue (100行以下) -->
<script setup lang="ts">
import type { RenderTree, RenderMedia } from '../../models/RenderTree';
import type { LanguageCode } from '../../composables/useTimeline';
import Avatar from './Avatar.vue';
import ArticleHeader from './ArticleHeader.vue';
import ArticleBody from './ArticleBody.vue';
import MediaGrid from '../media/MediaGrid.vue';

defineProps<{
  article: RenderTree;
  targetLang: LanguageCode;
  loading?: boolean;
}>();

const emit = defineEmits<{
  (e: 'clickArticle', id: string): void;
  (e: 'openMedia', m: RenderMedia, list: RenderMedia[], a: RenderTree): void;
  (e: 'clickTag', tag: string): void;
  (e: 'clickMention', handle: string): void;
}>();
</script>

<template>
  <div class="twitter-inline-thread relative pl-3 pr-2 py-2.5 my-1.5 bg-slate-900/50 hover:bg-slate-900/80 border-l-2 border-sky-500/60 rounded-r-lg text-left transition-colors cursor-pointer"
       @click.stop="emit('clickArticle', article.id)">
    <div class="flex items-start gap-2.5">
      <div class="flex-shrink-0 pt-0.5">
        <Avatar :avatarUrl="article.author.avatar_url" :handle="article.author.handle" class="w-8 h-8" />
      </div>
      <div class="flex-1 min-w-0 space-y-1">
        <ArticleHeader
          :author="article.author"
          :createdAt="article.created_at"
          :sourceUrl="article.source_url"
          :isPinned="false"
        />
        <ArticleBody
          :content="article.content"
          :targetLang="targetLang"
          @clickTag="(tag) => emit('clickTag', tag)"
          @clickMention="(handle) => emit('clickMention', handle)"
        />
        <MediaGrid
          v-if="article.media && article.media.length > 0"
          :media="article.media"
          @clickMedia="(m, l) => emit('openMedia', m, l || [], article)"
        />
      </div>
    </div>
  </div>
</template>
