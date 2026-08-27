<!-- frontend/src/components/article/ArticleDetailMain.vue (100行以下 - SPEC-PRINCIPLE-001) -->
<script setup lang="ts">
import type { RenderTree, RenderMedia } from '../../models/RenderTree';
import type { LanguageCode } from '../../composables/useTimeline';
import Avatar from './Avatar.vue';
import MediaGrid from '../media/MediaGrid.vue';
import ArticleStats from './ArticleStats.vue';
import ArticleBody from './ArticleBody.vue';
import { formatTweetDetailDate } from '../../utils/formatters';

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
  <div class="p-4 md:p-5 bg-slate-950 border-b border-slate-800/80 space-y-3 font-sans">
    <!-- ① 投稿者情報 (アバター + 表示名 + ハンドル) -->
    <div class="flex items-center gap-3 select-none">
      <Avatar :author="article.author" :avatar-url="article.author.avatar_url" :handle="article.author.handle" class="w-10 h-10 md:w-11 md:h-11 flex-shrink-0" />
      <div class="min-w-0 flex-1">
        <div class="flex items-center gap-1.5 leading-tight">
          <span class="text-sm md:text-base font-bold text-slate-100 truncate">{{ article.author.display_name || article.author.handle || 'User' }}</span>
          <span v-if="article.author.group_name" class="px-1.5 py-0.2 rounded-full bg-amber-500/20 text-amber-300 text-[9px] font-bold border border-amber-500/30 shrink-0">{{ article.author.group_name }}</span>
        </div>
        <div class="text-xs text-slate-400 font-mono leading-tight">@{{ article.author.handle || 'unknown' }}</div>
      </div>
    </div>

    <!-- 本文 -->
    <ArticleBody
      :content="article.content"
      :target-lang="targetLang"
      class="text-sm md:text-[15px] text-slate-100 leading-relaxed select-text py-1"
      @click-tag="(t) => $emit('clickTag', t)"
      @click-mention="(m) => $emit('clickMention', m)"
    />

    <!-- メディア一覧 -->
    <MediaGrid
      v-if="article.media?.length"
      :media="article.media"
      @clickMedia="(m, l) => $emit('clickMedia', m, l, article)"
      @retry="(id) => $emit('retryMedia', id)"
    />

    <!-- ③ X公式ライクな詳細日時表示 (例: 午後6:03 · 2026年8月25日) -->
    <div class="text-xs text-slate-500 font-sans pt-3 pb-2 border-t border-slate-800/60 flex items-center justify-between select-none">
      <div class="flex items-center gap-2">
        <span class="text-slate-400 hover:underline cursor-default">{{ formatTweetDetailDate(article.created_at) }}</span>
        <span v-if="article.metrics?.views" class="text-slate-500">· <strong class="text-slate-300 font-mono">{{ article.metrics.views }}</strong> 件の表示</span>
      </div>
      <span v-if="article.via" class="text-[11px] text-slate-500">via {{ article.via }}</span>
    </div>

    <!-- ④ X完全互換エンゲージメントバー (返信・リポスト・いいね・ブックマーク・共有) -->
    <div class="border-t border-slate-800/80 pt-1">
      <ArticleStats :metrics="article.metrics" :is-liked="article.is_liked" @toggle-like="$emit('toggleLike', article.id)" />
    </div>
  </div>
</template>
