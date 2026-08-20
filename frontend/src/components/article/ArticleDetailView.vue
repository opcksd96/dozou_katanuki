<script setup lang="ts">
import { ref, computed } from 'vue';
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

// 親スレッド（メイン記事より前に投稿された関連ポスト、時系列昇順）
const parentChain = computed(() => {
  if (!props.thread || props.thread.length === 0) return [];
  const mainTime = new Date(props.article.created_at).getTime();
  return props.thread
    .filter((t) => t.id !== props.article.id && new Date(t.created_at).getTime() < mainTime)
    .sort((a, b) => new Date(a.created_at).getTime() - new Date(b.created_at).getTime());
});

// 子スレッド（メイン記事の後に投稿された関連ポスト、時系列昇順）
const childReplies = computed(() => {
  if (!props.thread || props.thread.length === 0) return [];
  const mainTime = new Date(props.article.created_at).getTime();
  return props.thread
    .filter((t) => t.id !== props.article.id && new Date(t.created_at).getTime() >= mainTime)
    .sort((a, b) => new Date(a.created_at).getTime() - new Date(b.created_at).getTime());
});
</script>

<template>
  <div class="w-full bg-slate-950 min-h-screen text-left">
    <!-- 上部ナビゲーションバー -->
    <div class="sticky top-0 z-30 bg-slate-950/85 backdrop-blur border-b border-slate-800 px-4 py-3 flex items-center justify-between">
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

    <div>
      <!-- 1. 親スレッドチェーン（上位リプライ群） -->
      <div v-if="parentChain.length > 0" class="divide-y divide-slate-800/40">
        <div
          v-for="(t, idx) in parentChain"
          :key="t.id"
          @click="emit('selectArticle', t.id)"
          :class="[
            'twitter-card relative p-4 bg-slate-950/60 hover:bg-slate-900/40 transition-colors cursor-pointer',
            focusedArticleId === t.id ? 'is-focused' : ''
          ]"
        >
          <!-- 連続スレッド接続線 -->
          <div class="absolute left-8 top-12 bottom-0 w-0.5 bg-slate-700"></div>

          <div class="twitter-avatar-col relative z-10 flex-shrink-0 pt-0.5">
            <Avatar :avatarUrl="t.author.avatar_url" :handle="t.author.handle" />
          </div>

          <div class="twitter-content-col relative z-10 flex-1 min-w-0 space-y-1.5">
            <div class="flex items-center justify-between gap-2 text-xs">
              <div class="flex items-center gap-1.5 truncate">
                <span class="font-bold text-white truncate">{{ t.author.display_name || t.author.handle }}</span>
                <span class="text-slate-400 font-mono">@{{ t.author.handle }}</span>
              </div>
              <span class="text-slate-500 font-mono text-[11px] whitespace-nowrap">{{ t.created_at }}</span>
            </div>

            <div v-if="t.reply_to_handle" class="twitter-reply-badge">
              <span>返信先:</span>
              <a @click.stop="emit('clickMention', t.reply_to_handle)">@{{ t.reply_to_handle }}</a>
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

      <!-- 2. フォーカス中のメインポスト（フルサイズ表示） -->
      <div
        :class="[
          'twitter-main-focus-card p-5 bg-slate-950 border-b border-slate-800 space-y-4 relative',
          focusedArticleId === article.id ? 'is-focused' : ''
        ]"
      >
        <!-- 投稿者ヘッダー -->
        <div class="flex items-center gap-3">
          <Avatar :avatarUrl="article.author.avatar_url" :handle="article.author.handle" class="w-12 h-12" />
          <div class="flex-1 min-w-0">
            <h3 class="font-bold text-white text-base leading-tight truncate">{{ article.author.display_name || article.author.handle }}</h3>
            <p class="text-xs text-slate-400 font-mono">@{{ article.author.handle }}</p>
          </div>
          <button
            @click="emit('toggleLike', article.id)"
            class="p-2 rounded-full hover:bg-slate-800 transition-colors cursor-pointer"
            :title="article.is_liked ? 'ブックマーク解除 (L)' : 'ブックマーク (L)'"
          >
            <span :class="article.is_liked ? 'text-pink-500' : 'text-slate-500'">{{ article.is_liked ? '❤️' : '🤍' }}</span>
          </button>
        </div>

        <!-- リプライ先バッジ -->
        <div v-if="article.reply_to_handle" class="twitter-reply-badge text-sm">
          <span>返信先:</span>
          <a @click.stop="emit('clickMention', article.reply_to_handle)">@{{ article.reply_to_handle }}</a>
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
      <div v-if="childReplies.length > 0" class="divide-y divide-slate-800/40">
        <div
          v-for="(t, idx) in childReplies"
          :key="t.id"
          @click="emit('selectArticle', t.id)"
          :class="[
            'twitter-card relative p-4 bg-slate-950/60 hover:bg-slate-900/40 transition-colors cursor-pointer',
            focusedArticleId === t.id ? 'is-focused' : ''
          ]"
        >
          <!-- 子スレッド接続線（複数ある場合） -->
          <div v-if="idx < childReplies.length - 1" class="absolute left-8 top-12 bottom-0 w-0.5 bg-slate-800"></div>

          <div class="twitter-avatar-col relative z-10 flex-shrink-0 pt-0.5">
            <Avatar :avatarUrl="t.author.avatar_url" :handle="t.author.handle" />
          </div>

          <div class="twitter-content-col relative z-10 flex-1 min-w-0 space-y-1.5">
            <div class="flex items-center justify-between gap-2 text-xs">
              <div class="flex items-center gap-1.5 truncate">
                <span class="font-bold text-white truncate">{{ t.author.display_name || t.author.handle }}</span>
                <span class="text-slate-400 font-mono">@{{ t.author.handle }}</span>
              </div>
              <span class="text-slate-500 font-mono text-[11px] whitespace-nowrap">{{ t.created_at }}</span>
            </div>

            <div v-if="t.reply_to_handle" class="twitter-reply-badge">
              <span>返信先:</span>
              <a @click.stop="emit('clickMention', t.reply_to_handle)">@{{ t.reply_to_handle }}</a>
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
    </div>
  </div>
</template>
