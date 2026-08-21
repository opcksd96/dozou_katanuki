<!-- frontend/src/components/article/ArticleDetailItem.vue (100行以下) -->
<script setup lang="ts">
import type { RenderTree, RenderMedia } from '../../models/RenderTree';
import type { LanguageCode } from '../../composables/useTimeline';
import Avatar from './Avatar.vue';
import MediaGrid from '../media/MediaGrid.vue';
import ArticleBody from './ArticleBody.vue';

defineProps<{
  item: RenderTree;
  targetLang: LanguageCode;
  isFocused?: boolean;
}>();

defineEmits<{
  (e: 'selectArticle', id: string): void;
  (e: 'clickMedia', media: RenderMedia, list?: RenderMedia[], article?: RenderTree): void;
  (e: 'clickTag', tag: string): void;
  (e: 'clickMention', handle: string): void;
}>();
</script>

<template>
  <div
    @click="$emit('selectArticle', item.id)"
    class="p-4 border-b border-slate-800/60 hover:bg-slate-900/40 cursor-pointer transition-colors"
    :class="{ 'bg-blue-950/20 border-l-2 border-blue-500': isFocused }"
  >
    <div class="flex gap-3">
      <Avatar :author="item.author" class="w-10 h-10 flex-shrink-0" />
      <div class="flex-1 min-w-0">
        <div class="flex items-center gap-2 mb-1">
          <span class="font-bold text-slate-100 text-xs truncate">{{ item.author.display_name }}</span>
          <span class="text-slate-500 text-xs font-mono">@{{ item.author.handle }}</span>
          <span class="text-slate-600 text-xs">・</span>
          <span class="text-slate-500 text-xs">{{ new Date(item.created_at).toLocaleDateString() }}</span>
        </div>
        <ArticleBody :content="item.content" :target-lang="targetLang" class="text-xs text-slate-200 mb-2" @click-tag="(t) => $emit('clickTag', t)" @click-mention="(m) => $emit('clickMention', m)" />
        <MediaGrid v-if="item.media?.length" :media="item.media" @clickMedia="(m, l) => $emit('clickMedia', m, l, item)" />
      </div>
    </div>
  </div>
</template>
