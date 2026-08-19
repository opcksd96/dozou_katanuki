# 1. 必要なディレクトリの作成
New-Item -ItemType Directory -Force -Path "frontend/src/models", "frontend/src/mock", "frontend/src/utils", "frontend/src/composables", "frontend/src/components/article", "frontend/src/components/media"

# 2. models/RenderTree.ts (仕様書 2.3節 完全準拠)
@'
export type DownloadStatus = 'QUEUED' | 'COMPLETED' | 'DEAD_404' | 'OUTSOURCED' | 'RETAINED';

export interface RenderMetrics {
  replies: number;
  retweets: number;
  likes: number;
  views?: number;
}

export interface RenderAuthor {
  numeric_id: string;
  handle: string;
  display_name: string;
  avatar_url: string;
  bio: string;
}

export interface RenderMedia {
  id: string;
  type: 'image' | 'video' | 'gif';
  download_status: DownloadStatus;
  failed_reason?: string;
  urls: {
    stream: string;
    image: string;
    thumbnail: string;
    original: string;
  };
  width?: number;
  height?: number;
}

export interface RenderTree {
  id: string;
  conversation_id: string;
  created_at: string;
  content: {
    original: string;
    ja?: string;
    en?: string;
    zh?: string;
  };
  author: RenderAuthor;
  media: RenderMedia[];
  metrics: RenderMetrics;
  source_url: string;
  is_liked: boolean;
  is_pinned: boolean;
  parent_id?: string;
}
'@ | Out-File -Encoding utf8 "frontend/src/models/RenderTree.ts"

# 3. mock/sample_render_trees.json (COMPLETED, QUEUED, DEAD_404 検証用)
@'
[
  {
    "id": "1879382757924868404",
    "conversation_id": "1879382757924868404",
    "created_at": "2024-01-15T12:30:00Z",
    "source_url": "https://web.archive.org/web/https://twitter.com/msluo14/status/1879382757924868404",
    "is_liked": false,
    "is_pinned": true,
    "author": {
      "numeric_id": "1001",
      "handle": "msluo14",
      "display_name": "Yike Luo",
      "avatar_url": "https://api.dicebear.com/7.x/bottts/svg?seed=msluo14",
      "bio": "Archival preservation target."
    },
    "content": {
      "original": "Archived snapshot from Wayback Machine. Deep persistence verified.",
      "ja": "Wayback Machineからの魚拓スナップショット。完全オフライン動態保存の検証に成功しました。",
      "en": "Archived snapshot from Wayback Machine. Deep persistence verified.",
      "zh": "来自 Wayback Machine 的归档快照。深度持久化验证成功。"
    },
    "media": [
      {
        "id": "F8wZ1abXYAAY7kL.jpg",
        "type": "image",
        "download_status": "COMPLETED",
        "urls": {
          "stream": "",
          "image": "https://picsum.photos/800/600?random=10",
          "thumbnail": "https://picsum.photos/400/300?random=10",
          "original": "https://web.archive.org/..."
        }
      },
      {
        "id": "sample_video_clip.mp4",
        "type": "video",
        "download_status": "OUTSOURCED",
        "failed_reason": "外部Motrix(aria2)にて並列ダウンロード中...",
        "urls": {
          "stream": "",
          "image": "",
          "thumbnail": "",
          "original": "https://web.archive.org/..."
        }
      }
    ],
    "metrics": {
      "replies": 12,
      "retweets": 45,
      "likes": 128,
      "views": 3200
    }
  },
  {
    "id": "1879382757924868405",
    "conversation_id": "1879382757924868405",
    "created_at": "2024-01-16T08:15:00Z",
    "source_url": "https://web.archive.org/...",
    "is_liked": true,
    "is_pinned": false,
    "author": {
      "numeric_id": "1001",
      "handle": "msluo14",
      "display_name": "Yike Luo",
      "avatar_url": "https://api.dicebear.com/7.x/bottts/svg?seed=msluo14",
      "bio": "Archival preservation target."
    },
    "content": {
      "original": "Testing dead link SVG filler fallback and animated GIF.",
      "ja": "404消失メディア用SVGフィラーとGIFアニメーションの描画フォールバック検証。"
    },
    "media": [
      {
        "id": "dead_asset.jpg",
        "type": "image",
        "download_status": "DEAD_404",
        "failed_reason": "Wayback CDN 404: 元サーバーから消滅しています",
        "urls": { "stream": "", "image": "", "thumbnail": "", "original": "" }
      },
      {
        "id": "reaction.gif",
        "type": "gif",
        "download_status": "QUEUED",
        "urls": { "stream": "", "image": "", "thumbnail": "", "original": "" }
      }
    ],
    "metrics": {
      "replies": 3,
      "retweets": 8,
      "likes": 56
    }
  }
]
'@ | Out-File -Encoding utf8 "frontend/src/mock/sample_render_trees.json"

# 4. composables/useTimeline.ts
@'
import { ref, computed } from 'vue';
import type { RenderTree } from '../models/RenderTree';
import mockData from '../mock/sample_render_trees.json';

export type SupportedLang = 'original' | 'ja' | 'en' | 'zh';

export function useTimeline() {
  const articles = ref<RenderTree[]>(mockData as RenderTree[]);
  const currentLang = ref<SupportedLang>('ja');
  const isLoading = ref<boolean>(false);

  const setLanguage = (lang: SupportedLang) => {
    currentLang.value = lang;
  };

  const toggleLike = (id: string) => {
    const target = articles.value.find((a) => a.id === id);
    if (target) {
      target.is_liked = !target.is_liked;
      target.metrics.likes += target.is_liked ? 1 : -1;
    }
  };

  const retryDownload = (mediaId: string) => {
    console.log(`[UDF Action] メディア再取得要求: ${mediaId}`);
  };

  return {
    articles,
    currentLang,
    isLoading,
    articleCount: computed(() => articles.value.length),
    setLanguage,
    toggleLike,
    retryDownload,
  };
}
'@ | Out-File -Encoding utf8 "frontend/src/composables/useTimeline.ts"

# 5. components/article/Avatar.vue
@'
<script setup lang="ts">
defineProps<{
  avatarUrl: string;
  handle: string;
}>();
</script>

<template>
  <div class="relative flex-shrink-0">
    <img
      :src="avatarUrl"
      :alt="handle"
      class="w-10 h-10 rounded-full border border-slate-700 object-cover bg-slate-800"
      loading="lazy"
    />
  </div>
</template>
'@ | Out-File -Encoding utf8 "frontend/src/components/article/Avatar.vue"

# 6. components/article/ArticleHeader.vue
@'
<script setup lang="ts">
import type { RenderAuthor } from '../../models/RenderTree';
import { formatDate } from '../../utils/formatters';
import Avatar from './Avatar.vue';

defineProps<{
  author: RenderAuthor;
  createdAt: string;
  sourceUrl: string;
  isPinned?: boolean;
}>();
</script>

<template>
  <div class="flex items-start justify-between gap-3 mb-2">
    <div class="flex items-center gap-3">
      <Avatar :avatarUrl="author.avatar_url" :handle="author.handle" />
      <div>
        <div class="flex items-center gap-2">
          <span class="font-bold text-slate-100 leading-tight">{{ author.display_name }}</span>
          <span v-if="isPinned" class="text-[10px] bg-amber-500/20 text-amber-300 px-1.5 py-0.5 rounded font-mono">PINNED</span>
        </div>
        <div class="text-xs text-slate-400">@{{ author.handle }} · {{ formatDate(createdAt) }}</div>
      </div>
    </div>
    <a :href="sourceUrl" target="_blank" rel="noopener noreferrer" class="text-xs text-blue-400 hover:underline">原本リンク</a>
  </div>
</template>
'@ | Out-File -Encoding utf8 "frontend/src/components/article/ArticleHeader.vue"

# 7. components/article/ArticleBody.vue
@'
<script setup lang="ts">
import { computed } from 'vue';
import type { SupportedLang } from '../../composables/useTimeline';

const props = defineProps<{
  content: { original: string; ja?: string; en?: string; zh?: string };
  currentLang: SupportedLang;
}>();

const text = computed(() => {
  if (props.currentLang === 'original') return props.content.original;
  return props.content[props.currentLang] || props.content.original;
});
</script>

<template>
  <div class="my-2 text-slate-200 text-sm leading-relaxed whitespace-pre-wrap">
    <p>{{ text }}</p>
    <div v-if="currentLang !== 'original' && content[currentLang]" class="mt-1 text-[11px] text-emerald-400 font-mono">
      [翻訳: {{ currentLang.toUpperCase() }}]
    </div>
  </div>
</template>
'@ | Out-File -Encoding utf8 "frontend/src/components/article/ArticleBody.vue"

# 8. components/article/ArticleStats.vue
@'
<script setup lang="ts">
import type { RenderMetrics } from '../../models/RenderTree';
import { formatStatNumber } from '../../utils/formatters';

defineProps<{
  metrics: RenderMetrics;
  isLiked: boolean;
}>();

const emit = defineEmits<{
  (e: 'toggleLike'): void;
}>();
</script>

<template>
  <div class="flex items-center gap-6 mt-3 text-xs text-slate-400 pt-2 border-t border-slate-800/60 font-mono">
    <span>💬 {{ formatStatNumber(metrics.replies) }}</span>
    <span>🔁 {{ formatStatNumber(metrics.retweets) }}</span>
    <button
      @click="emit('toggleLike')"
      class="flex items-center gap-1.5 transition-colors"
      :class="isLiked ? 'text-rose-500 font-bold' : 'hover:text-rose-400'"
    >
      <span>{{ isLiked ? '❤️' : '🤍' }}</span>
      <span>{{ formatStatNumber(metrics.likes) }}</span>
    </button>
    <span v-if="metrics.views" class="text-slate-500 ml-auto">👁️ {{ formatStatNumber(metrics.views) }}</span>
  </div>
</template>
'@ | Out-File -Encoding utf8 "frontend/src/components/article/ArticleStats.vue"

# 9. components/media/MediaFiller.vue (仕様書 2.6.2節 SVGフィラー＆リトライ)
@'
<script setup lang="ts">
import type { RenderMedia } from '../../models/RenderTree';

defineProps<{
  media: RenderMedia;
}>();

const emit = defineEmits<{
  (e: 'retry', mediaId: string): void;
}>();
</script>

<template>
  <div
    class="relative w-full h-48 flex flex-col items-center justify-center p-4 text-center rounded-lg overflow-hidden border border-slate-800 select-none"
    :class="{
      'bg-slate-900 animate-pulse': media.type === 'image',
      'bg-slate-950 animate-pulse': media.type === 'video',
      'bg-slate-900': media.type === 'gif'
    }"
  >
    <!-- 1. 画像用 SVG フィラー (カメラ) -->
    <svg v-if="media.type === 'image'" class="w-10 h-10 text-slate-600 mb-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M3 9a2 2 0 012-2h.93a2 2 0 001.664-.89l.812-1.22A2 2 0 0110.07 4h3.86a2 2 0 011.664.89l.812 1.22A2 2 0 0018.07 7H19a2 2 0 012 2v9a2 2 0 01-2 2H5a2 2 0 01-2-2V9z" />
      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M15 13a3 3 0 11-6 0 3 3 0 016 0z" />
    </svg>

    <!-- 2. 動画用 SVG フィラー (Play) -->
    <svg v-else-if="media.type === 'video'" class="w-12 h-12 text-slate-600 mb-2" fill="currentColor" viewBox="0 0 24 24">
      <path d="M8 5v14l11-7z"/>
    </svg>

    <!-- 3. GIF用 ピルバッジ -->
    <div v-else-if="media.type === 'gif'" class="bg-indigo-600 text-white font-bold text-xs px-2 py-0.5 rounded font-mono mb-2">GIF</div>

    <!-- ステータスオーバーレイ ＆ リトライボタン -->
    <div class="z-10 flex flex-col items-center gap-1">
      <span class="text-xs font-mono font-bold px-2 py-0.5 rounded bg-slate-800/80 text-amber-300 border border-slate-700">
        [{{ media.download_status }}]
      </span>
      <p v-if="media.failed_reason" class="text-[11px] text-slate-400 max-w-xs truncate">{{ media.failed_reason }}</p>
      <button
        v-if="['DEAD_404', 'OUTSOURCED', 'RETAINED'].includes(media.download_status)"
        @click="emit('retry', media.id)"
        class="mt-1 text-xs text-blue-400 hover:text-blue-300 underline font-mono cursor-pointer"
      >
        ↻ 再試行 (Retry)
      </button>
    </div>

    <!-- 動画下部シークバープレースホルダー -->
    <div v-if="media.type === 'video'" class="absolute bottom-0 left-0 right-0 h-1 bg-slate-800">
      <div class="h-full bg-slate-700 w-1/3"></div>
    </div>
  </div>
</template>
'@ | Out-File -Encoding utf8 "frontend/src/components/media/MediaFiller.vue"

# 10. components/media/MediaGrid.vue
@'
<script setup lang="ts">
import type { RenderMedia } from '../../models/RenderTree';
import MediaFiller from './MediaFiller.vue';

defineProps<{
  media: RenderMedia[];
}>();

const emit = defineEmits<{
  (e: 'retry', mediaId: string): void;
}>();
</script>

<template>
  <div v-if="media && media.length > 0" class="mt-3 rounded-lg overflow-hidden border border-slate-800 bg-black">
    <div :class="['grid gap-1', media.length === 1 ? 'grid-cols-1' : 'grid-cols-2']">
      <div v-for="item in media" :key="item.id" class="relative group">
        <!-- ローカル確保完了時: 描画 -->
        <template v-if="item.download_status === 'COMPLETED'">
          <img
            v-if="item.type === 'image' || item.type === 'gif'"
            :src="item.urls.image || item.urls.thumbnail"
            :alt="item.id"
            class="w-full h-auto max-h-[400px] object-cover hover:opacity-95 transition-opacity"
            loading="lazy"
          />
          <video
            v-else-if="item.type === 'video'"
            :src="item.urls.stream"
            controls
            preload="metadata"
            class="w-full max-h-[400px]"
          />
        </template>
        <!-- 未完了/失敗時: SVGプレースホルダー・フィラー -->
        <MediaFiller v-else :media="item" @retry="(id) => emit('retry', id)" />
      </div>
    </div>
  </div>
</template>
'@ | Out-File -Encoding utf8 "frontend/src/components/media/MediaGrid.vue"

# 11. components/article/ArticleCard.vue
@'
<script setup lang="ts">
import type { RenderTree } from '../../models/RenderTree';
import type { SupportedLang } from '../../composables/useTimeline';
import ArticleHeader from './ArticleHeader.vue';
import ArticleBody from './ArticleBody.vue';
import MediaGrid from '../media/MediaGrid.vue';
import ArticleStats from './ArticleStats.vue';

defineProps<{
  article: RenderTree;
  currentLang: SupportedLang;
}>();

const emit = defineEmits<{
  (e: 'toggleLike', id: string): void;
  (e: 'retryMedia', mediaId: string): void;
}>();
</script>

<template>
  <article class="p-4 bg-slate-900 border border-slate-800 rounded-xl mb-4 hover:border-slate-700/80 transition-colors">
    <ArticleHeader
      :author="article.author"
      :createdAt="article.created_at"
      :sourceUrl="article.source_url"
      :isPinned="article.is_pinned"
    />
    <ArticleBody :content="article.content" :currentLang="currentLang" />
    <MediaGrid :media="article.media" @retry="(mediaId) => emit('retryMedia', mediaId)" />
    <ArticleStats :metrics="article.metrics" :isLiked="article.is_liked" @toggleLike="emit('toggleLike', article.id)" />
  </article>
</template>
'@ | Out-File -Encoding utf8 "frontend/src/components/article/ArticleCard.vue"

# 12. App.vue (ルートコンテナ統合)
@'
<script setup lang="ts">
import { useTimeline } from './composables/useTimeline';
import ArticleCard from './components/article/ArticleCard.vue';

const { articles, currentLang, setLanguage, toggleLike, retryDownload } = useTimeline();
</script>

<template>
  <div class="min-h-screen bg-slate-950 text-slate-100 flex flex-col items-center py-6 px-4">
    <header class="w-full max-w-2xl flex items-center justify-between pb-4 border-b border-slate-800 mb-6">
      <div>
        <h1 class="text-xl font-bold tracking-tight text-white">dozou_katanuki</h1>
        <p class="text-xs text-slate-400 font-mono">Dynamic Archival System (SPEC-FRONTEND-001)</p>
      </div>
      <div class="flex bg-slate-900 border border-slate-800 rounded-lg p-1 text-xs">
        <button
          v-for="lang in (['original', 'ja', 'en', 'zh'] as const)"
          :key="lang"
          @click="setLanguage(lang)"
          :class="[
            'px-2.5 py-1 rounded transition-colors',
            currentLang === lang ? 'bg-blue-600 text-white font-bold' : 'text-slate-400 hover:text-slate-200'
          ]"
        >
          {{ lang.toUpperCase() }}
        </button>
      </div>
    </header>

    <main class="w-full max-w-2xl">
      <ArticleCard
        v-for="article in articles"
        :key="article.id"
        :article="article"
        :currentLang="currentLang"
        @toggleLike="toggleLike"
        @retryMedia="retryDownload"
      />
    </main>
  </div>
</template>
'@ | Out-File -Encoding utf8 "frontend/src/App.vue"
