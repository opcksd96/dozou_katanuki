<script setup lang="ts">
import { ref, computed } from 'vue';
import type { RenderMedia, RenderTree } from '../../models/RenderTree';
import type { LanguageCode } from '../../composables/useTimeline';
import StashPlayer from './StashPlayer.vue';

const props = defineProps<{
  media: RenderMedia | null;
  article?: RenderTree | null;
  targetLang?: LanguageCode;
  hasNext?: boolean;
  hasPrev?: boolean;
}>();

const emit = defineEmits<{
  (e: 'close'): void;
  (e: 'next'): void;
  (e: 'prev'): void;
  (e: 'toggleLike', id: string): void;
}>();

// 翻訳言語のローカル切り替え
const selectedLang = ref<LanguageCode>(props.targetLang || 'ja');

// 表示テキスト
const displayText = computed(() => {
  if (!props.article) return '';
  const content = props.article.content;
  if (!content) return '';

  if (selectedLang.value === 'ja' && content.ja) return content.ja;
  if (selectedLang.value === 'en' && content.en) return content.en;
  if (selectedLang.value === 'zh' && content.zh) return content.zh;
  return content.original;
});

// Stash への直接導線 URL を算出
const stashDirectUrl = computed(() => {
  if (!props.media) return null;

  if (props.media.stash_scene_id) {
    return `http://127.0.0.1:9999/scenes/${props.media.stash_scene_id}`;
  }
  if (props.media.stash_image_id) {
    return `http://127.0.0.1:9999/images/${props.media.stash_image_id}`;
  }
  if (props.media.urls?.stream) {
    const match = props.media.urls.stream.match(/\/stash-proxy\/scene\/([^/]+)/);
    if (match && match[1]) {
      return `http://127.0.0.1:9999/scenes/${match[1]}`;
    }
  }
  if (props.media.urls?.image) {
    const match = props.media.urls.image.match(/\/stash-proxy\/image\/([^/]+)/);
    if (match && match[1]) {
      return `http://127.0.0.1:9999/images/${match[1]}`;
    }
  }
  return 'http://127.0.0.1:9999';
});
</script>

<template>
  <Transition name="fade">
    <div
      v-if="media"
      class="fixed inset-0 z-50 bg-black/95 backdrop-blur-md flex flex-col lg:flex-row items-stretch justify-between overflow-hidden select-none cursor-pointer"
      @click="emit('close')"
    >
      <!-- 左/中央: メディア表示エリア (背景クリックでオーバーレイ閉じる) -->
      <div
        class="relative flex-1 flex items-center justify-center p-4 min-h-[60vh] lg:min-h-full overflow-hidden"
        @click="emit('close')"
      >
        <!-- 前へ送りボタン -->
        <button
          v-if="hasPrev"
          @click.stop="emit('prev')"
          title="前のメディア (←)"
          class="absolute left-4 z-30 text-white/70 hover:text-white bg-black/60 hover:bg-black/80 border border-white/10 p-3 rounded-full transition-colors cursor-pointer"
        >
          <svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
          </svg>
        </button>

        <!-- 次へ送りボタン -->
        <button
          v-if="hasNext"
          @click.stop="emit('next')"
          title="次のメディア (→)"
          class="absolute right-4 z-30 text-white/70 hover:text-white bg-black/60 hover:bg-black/80 border border-white/10 p-3 rounded-full transition-colors cursor-pointer"
        >
          <svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
          </svg>
        </button>

        <!-- メディア実体 (画像・動画クリック時は閉じないよう @click.stop) -->
        <div class="relative max-w-full max-h-[88vh] flex items-center justify-center cursor-default">
          <img
            v-if="media.type === 'image' || media.type === 'gif'"
            :src="media.urls.image || media.urls.original"
            :alt="media.id"
            class="max-w-full max-h-[88vh] object-contain rounded shadow-2xl transition-all select-none"
            @click.stop
          />
          <div
            v-else-if="media.type === 'video'"
            class="w-full max-w-4xl max-h-[85vh] flex items-center justify-center"
            @click.stop
          >
            <StashPlayer
              :src="media.urls.stream || media.urls.original"
              :poster="media.urls.thumbnail"
              :stashSceneId="media.stash_scene_id"
              :autoplay="true"
              :controls="true"
            />
          </div>
        </div>
      </div>

      <!-- 右側: 元ツイート詳細サイドバー (クリックイベントを stop して操作可能に) -->
      <div
        class="w-full lg:w-96 bg-slate-950/95 border-t lg:border-t-0 lg:border-l border-slate-800 flex flex-col justify-between max-h-[40vh] lg:max-h-full overflow-y-auto cursor-default z-40 shadow-2xl backdrop-blur-xl"
        @click.stop
      >
        <!-- サイドバーヘッダー ＆ 投稿者情報 -->
        <div class="p-5 space-y-4 border-b border-slate-800/80">
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-3 truncate">
              <!-- アバター -->
              <div class="w-10 h-10 rounded-full overflow-hidden bg-slate-800 border border-slate-700 flex-shrink-0">
                <img
                  v-if="article?.author?.avatar_url"
                  :src="article.author.avatar_url"
                  :alt="article.author.display_name"
                  class="w-full h-full object-cover"
                />
                <div v-else class="w-full h-full flex items-center justify-center text-sm font-bold text-slate-400">
                  {{ article?.author?.display_name?.charAt(0) || '👤' }}
                </div>
              </div>

              <!-- ユーザー名 -->
              <div class="truncate">
                <div class="text-sm font-bold text-white truncate">
                  {{ article?.author?.display_name || 'Anonymous' }}
                </div>
                <div class="text-xs text-slate-400 font-mono">
                  @{{ article?.author?.handle || 'user' }}
                </div>
              </div>
            </div>

            <!-- 閉じるボタン -->
            <button
              @click="emit('close')"
              title="閉じる (Esc)"
              class="text-slate-400 hover:text-white bg-slate-800/80 hover:bg-slate-700 p-2 rounded-full transition-colors flex-shrink-0 cursor-pointer"
            >
              <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>

          <!-- 投稿日時 ＆ 外部リンク導線 -->
          <div class="flex items-center justify-between text-xs text-slate-500 font-mono">
            <span>{{ article?.created_at ? new Date(article.created_at).toLocaleString() : '' }}</span>
            <div class="flex items-center gap-2">
              <a
                v-if="stashDirectUrl"
                :href="stashDirectUrl"
                target="_blank"
                rel="noopener noreferrer"
                class="px-2 py-0.5 rounded bg-purple-950/70 hover:bg-purple-900 text-purple-300 border border-purple-700/50 text-[11px] font-semibold transition-colors flex items-center gap-1 cursor-pointer"
              >
                <span>📦</span> Stash
              </a>
              <a
                v-if="article?.source_url || media.urls?.original"
                :href="article?.source_url || media.urls?.original"
                target="_blank"
                rel="noopener noreferrer"
                class="px-2 py-0.5 rounded bg-slate-900 hover:bg-slate-800 text-blue-400 border border-slate-700 text-[11px] font-semibold transition-colors flex items-center gap-1 cursor-pointer"
              >
                <span>🌐</span> 原本
              </a>
            </div>
          </div>
        </div>

        <!-- 本文エリア ＆ 翻訳セレクタ -->
        <div class="p-5 flex-1 space-y-3 overflow-y-auto">
          <!-- 翻訳言語切り替えタブ -->
          <div v-if="article?.content?.ja || article?.content?.en || article?.content?.zh" class="flex items-center gap-1 bg-slate-900/80 p-1 rounded-lg border border-slate-800 text-xs w-fit">
            <button
              @click="selectedLang = 'ja'"
              class="px-2.5 py-0.5 rounded transition-all font-semibold text-[11px]"
              :class="selectedLang === 'ja' ? 'bg-blue-600 text-white' : 'text-slate-400 hover:text-slate-200'"
            >
              🇯🇵 日本語
            </button>
            <button
              v-if="article.content.en"
              @click="selectedLang = 'en'"
              class="px-2.5 py-0.5 rounded transition-all font-semibold text-[11px]"
              :class="selectedLang === 'en' ? 'bg-blue-600 text-white' : 'text-slate-400 hover:text-slate-200'"
            >
              🇺🇸 EN
            </button>
            <button
              v-if="article.content.zh"
              @click="selectedLang = 'zh'"
              class="px-2.5 py-0.5 rounded transition-all font-semibold text-[11px]"
              :class="selectedLang === 'zh' ? 'bg-blue-600 text-white' : 'text-slate-400 hover:text-slate-200'"
            >
              🇨🇳 中文
            </button>
            <button
              @click="selectedLang = 'original' as any"
              class="px-2.5 py-0.5 rounded transition-all font-semibold text-[11px]"
              :class="selectedLang === ('original' as any) ? 'bg-slate-700 text-white' : 'text-slate-400 hover:text-slate-200'"
            >
              原文
            </button>
          </div>

          <!-- 本文テキスト -->
          <p class="text-sm text-slate-100 leading-relaxed font-sans whitespace-pre-wrap select-text">
            {{ displayText || article?.content?.original }}
          </p>
        </div>

        <!-- フッター: エンゲージメント統計 ＆ アクション -->
        <div class="p-5 border-t border-slate-800/80 bg-slate-900/40 flex items-center justify-between text-xs text-slate-400 font-mono">
          <div class="flex items-center gap-4">
            <span class="flex items-center gap-1" title="リプライ">
              <span>💬</span> {{ article?.metrics?.replies || 0 }}
            </span>
            <span class="flex items-center gap-1" title="リツイート">
              <span>🔁</span> {{ article?.metrics?.retweets || 0 }}
            </span>
            <button
              v-if="article"
              @click="emit('toggleLike', article.id)"
              class="flex items-center gap-1 transition-colors cursor-pointer"
              :class="article.is_liked ? 'text-rose-400 font-bold' : 'hover:text-rose-400'"
              title="いいね"
            >
              <span>{{ article.is_liked ? '❤️' : '🤍' }}</span>
              <span>{{ article.metrics?.likes || 0 }}</span>
            </button>
          </div>

          <div class="text-[10px] text-slate-500">
            {{ media.id }}
          </div>
        </div>
      </div>
    </div>
  </Transition>
</template>

<style scoped>
.fade-enter-active, .fade-leave-active { transition: opacity 0.2s ease; }
.fade-enter-from, .fade-leave-to { opacity: 0; }
</style>
