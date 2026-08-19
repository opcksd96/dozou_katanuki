<script setup lang="ts">
import { ref, watch, onMounted, onUnmounted } from 'vue';
import type { RenderTree, RenderMedia } from '../../models/RenderTree';
import type { LanguageCode } from '../../composables/useTimeline';
import Avatar from './Avatar.vue';
import MediaGrid from '../media/MediaGrid.vue';
import ArticleStats from './ArticleStats.vue';
import ArticleBody from './ArticleBody.vue';

const props = defineProps<{
  article: RenderTree;
  thread: RenderTree[];
  targetLang: LanguageCode;
  loading?: boolean;
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

const copySuccess = ref(false);

const copyArticleUrl = async () => {
  const url = props.article.source_url || `https://twitter.com/${props.article.author.handle}/status/${props.article.id}`;
  try {
    await navigator.clipboard.writeText(url);
    copySuccess.value = true;
    setTimeout(() => { copySuccess.value = false; }, 2000);
  } catch (err) {
    console.error('Copy failed', err);
  }
};

const handleKeyDown = (e: KeyboardEvent) => {
  if (e.key === 'Escape') {
    emit('back');
  }
};

onMounted(() => {
  window.addEventListener('keydown', handleKeyDown);
});

onUnmounted(() => {
  window.removeEventListener('keydown', handleKeyDown);
});
</script>

<template>
  <div class="w-full bg-slate-950 min-h-screen">
    <!-- 上部ナビゲーションバー -->
    <div class="sticky top-0 z-20 bg-slate-950/80 backdrop-blur border-b border-slate-800 px-4 py-3 flex items-center justify-between">
      <div class="flex items-center gap-3">
        <button
          @click="emit('back')"
          title="タイムラインへ戻る (Esc)"
          class="p-2 rounded-full hover:bg-slate-800 text-slate-300 hover:text-white transition-colors cursor-pointer flex items-center justify-center"
        >
          <span class="text-base font-bold">←</span>
        </button>
        <div>
          <h2 class="text-sm font-bold text-white tracking-tight">ポスト</h2>
          <p class="text-[10px] text-slate-400 font-mono">ID: {{ article.id }}</p>
        </div>
      </div>

      <div class="flex items-center gap-2">
        <button
          @click="copyArticleUrl"
          title="原本URLをコピー"
          class="px-2.5 py-1 rounded bg-slate-900 hover:bg-slate-800 text-slate-300 hover:text-white border border-slate-700/60 text-xs font-mono transition-colors cursor-pointer flex items-center gap-1"
        >
          <span>{{ copySuccess ? '✓ Copied' : '🔗 Link' }}</span>
        </button>
        <a
          v-if="article.source_url"
          :href="article.source_url"
          target="_blank"
          rel="noopener noreferrer"
          title="Wayback Machine で原本を開く"
          class="px-2.5 py-1 rounded bg-blue-900/40 hover:bg-blue-800/60 text-blue-300 hover:text-blue-200 border border-blue-700/50 text-xs font-mono transition-colors cursor-pointer"
        >
          🏛️ Wayback
        </a>
      </div>
    </div>

    <div class="divide-y divide-slate-800">
      <!-- 1. 親スレッド (もし本ポストの前に発言がある場合) -->
      <div v-if="thread && thread.length > 0" class="divide-y divide-slate-800/60">
        <template v-for="t in thread" :key="t.id">
          <div
            v-if="t.id !== article.id && new Date(t.created_at).getTime() < new Date(article.created_at).getTime()"
            @click="emit('selectArticle', t.id)"
            class="p-4 bg-slate-950/40 hover:bg-slate-900/30 transition-colors cursor-pointer relative"
          >
            <!-- スレッド接続線 -->
            <div class="absolute left-9 top-14 bottom-0 w-0.5 bg-slate-800"></div>

            <div class="flex items-start gap-3">
              <Avatar :avatarUrl="t.author.avatar_url" :handle="t.author.handle" />
              <div class="flex-1 min-w-0 space-y-1">
                <div class="flex items-center gap-2 text-xs">
                  <span class="font-bold text-white">{{ t.author.display_name || t.author.handle }}</span>
                  <span class="text-slate-400 font-mono">@{{ t.author.handle }}</span>
                  <span class="text-slate-500">·</span>
                  <span class="text-slate-500 font-mono">{{ t.created_at }}</span>
                </div>
                <ArticleBody
                  :content="t.content"
                  :targetLang="targetLang"
                  @clickTag="(tag) => emit('clickTag', tag)"
                  @clickMention="(handle) => emit('clickMention', handle)"
                />
                <MediaGrid
                  v-if="t.media && t.media.length > 0"
                  :media="t.media"
                  @retry="(mId) => emit('retryMedia', mId)"
                  @clickMedia="(m, list) => emit('clickMedia', m, list, t)"
                />
              </div>
            </div>
          </div>
        </template>
      </div>

      <!-- 2. フォーカス中のメインポスト (フルサイズ表示) -->
      <div class="p-5 bg-slate-950 space-y-4">
        <!-- 投稿者ヘッダー -->
        <div class="flex items-center gap-3">
          <img
            :src="article.author.avatar_url"
            :alt="article.author.handle"
            class="w-12 h-12 rounded-full border border-slate-700 object-cover bg-slate-800 shadow"
          />
          <div class="flex-1 min-w-0">
            <h3 class="font-bold text-white text-base leading-tight">{{ article.author.display_name || article.author.handle }}</h3>
            <p class="text-xs text-slate-400 font-mono">@{{ article.author.handle }}</p>
          </div>
          <button
            @click="emit('toggleLike', article.id)"
            class="p-2 rounded-full hover:bg-slate-800 transition-colors cursor-pointer"
            :title="article.is_liked ? 'ブックマーク解除' : 'ブックマーク'"
          >
            <span :class="article.is_liked ? 'text-pink-500' : 'text-slate-500'">{{ article.is_liked ? '❤️' : '🤍' }}</span>
          </button>
        </div>

        <!-- 本文 (大きなフォント) -->
        <div class="text-base text-slate-100 leading-relaxed font-normal whitespace-pre-line break-words pt-1">
          <ArticleBody
            :content="article.content"
            :targetLang="targetLang"
            @clickTag="(tag) => emit('clickTag', tag)"
            @clickMention="(handle) => emit('clickMention', handle)"
          />
        </div>

        <!-- メディアグリッド (フル解像度) -->
        <div v-if="article.media && article.media.length > 0" class="pt-2">
          <MediaGrid
            :media="article.media"
            @retry="(mId) => emit('retryMedia', mId)"
            @clickMedia="(m, list) => emit('clickMedia', m, list, article)"
          />
        </div>

        <!-- タイムスタンプ & エンゲージメント情報 -->
        <div class="pt-3 border-t border-slate-800/80 flex flex-wrap items-center justify-between gap-2 text-xs text-slate-400 font-mono">
          <time>{{ article.created_at }}</time>
          <div class="flex items-center gap-4 text-slate-300">
            <span v-if="article.metrics.retweets"><strong>{{ article.metrics.retweets }}</strong> <span class="text-slate-500">Reposts</span></span>
            <span v-if="article.metrics.likes"><strong>{{ article.metrics.likes }}</strong> <span class="text-slate-500">Likes</span></span>
          </div>
        </div>

        <!-- アクションバー -->
        <div class="pt-3 border-t border-slate-800/80">
          <ArticleStats
            :metrics="article.metrics"
            :isLiked="article.is_liked"
            @toggleLike="emit('toggleLike', article.id)"
          />
        </div>
      </div>

      <!-- 3. 子リプライ・後続スレッド -->
      <div v-if="thread && thread.length > 0" class="divide-y divide-slate-800/60">
        <template v-for="t in thread" :key="t.id">
          <div
            v-if="t.id !== article.id && new Date(t.created_at).getTime() >= new Date(article.created_at).getTime()"
            @click="emit('selectArticle', t.id)"
            class="p-4 bg-slate-950/40 hover:bg-slate-900/30 transition-colors cursor-pointer"
          >
            <div class="flex items-start gap-3">
              <Avatar :avatarUrl="t.author.avatar_url" :handle="t.author.handle" />
              <div class="flex-1 min-w-0 space-y-1">
                <div class="flex items-center gap-2 text-xs">
                  <span class="font-bold text-white">{{ t.author.display_name || t.author.handle }}</span>
                  <span class="text-slate-400 font-mono">@{{ t.author.handle }}</span>
                  <span class="text-slate-500">·</span>
                  <span class="text-slate-500 font-mono">{{ t.created_at }}</span>
                </div>
                <ArticleBody
                  :content="t.content"
                  :targetLang="targetLang"
                  @clickTag="(tag) => emit('clickTag', tag)"
                  @clickMention="(handle) => emit('clickMention', handle)"
                />
                <MediaGrid
                  v-if="t.media && t.media.length > 0"
                  :media="t.media"
                  @retry="(mId) => emit('retryMedia', mId)"
                  @clickMedia="(m, list) => emit('clickMedia', m, list, t)"
                />
              </div>
            </div>
          </div>
        </template>
      </div>
    </div>
  </div>
</template>
