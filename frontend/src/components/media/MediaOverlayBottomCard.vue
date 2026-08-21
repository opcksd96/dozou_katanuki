<!-- frontend/src/components/media/MediaOverlayBottomCard.vue -->
<script setup lang="ts">
import { ref, computed } from 'vue';
import type { RenderMedia, RenderTree } from '../../models/RenderTree';
import type { LanguageCode } from '../../composables/useTimeline';
import Avatar from '../article/Avatar.vue';

const props = defineProps<{
  media: RenderMedia;
  article: RenderTree;
  targetLang?: LanguageCode;
}>();

defineEmits<{ (e: 'toggleLike', id: string): void }>();

const isExpanded = ref(false);
const selectedLang = ref<LanguageCode>(props.targetLang || 'ja');

const displayText = computed(() => {
  const c = props.article.content;
  if (!c) return '';
  if (selectedLang.value === 'ja' && c.ja) return c.ja;
  if (selectedLang.value === 'en' && c.en) return c.en;
  if (selectedLang.value === 'zh' && c.zh) return c.zh;
  return c.original;
});

const currentHostname = computed(() => {
  if (typeof window !== 'undefined' && window.location?.hostname) {
    return window.location.hostname;
  }
  return '127.0.0.1';
});

const stashDirectUrl = computed(() => {
  if (props.media.stash_scene_id) return `http://${currentHostname.value}:9999/scenes/${props.media.stash_scene_id}`;
  if (props.media.stash_image_id) return `http://${currentHostname.value}:9999/images/${props.media.stash_image_id}`;
  return `http://${currentHostname.value}:9999`;
});
</script>

<template>
  <div class="w-full flex items-end justify-between gap-4 p-4 md:p-8 select-text pointer-events-none">
    <!-- 左側: 文字情報オーバーレイ (ボックスなし、ダイレクト描画) -->
    <div class="flex-1 max-w-xl space-y-2 pointer-events-auto text-left">
      <!-- 投稿者ヘッダー (アバター + 名前) -->
      <div class="flex items-center gap-2.5">
        <Avatar
          :avatar-url="article.author.avatar_url"
          :handle="article.author.handle"
          :author="article.author"
          class="w-8 h-8 md:w-9 md:h-9 rounded-full shadow-lg border border-white/40 drop-shadow-md"
        />
        <div class="flex flex-col sm:flex-row sm:items-center sm:gap-2 leading-tight">
          <span class="text-base md:text-lg font-bold text-white drop-shadow-[0_2px_4px_rgba(0,0,0,0.9)]">
            @{{ article.author.display_name || article.author.handle }}
          </span>
          <span class="text-xs text-white/70 font-mono drop-shadow-[0_1px_3px_rgba(0,0,0,0.8)]">
            @{{ article.author.handle }}
          </span>
        </div>
      </div>

      <!-- 本文テキスト (ボックスなし、ドロップシャドウでくっきり) -->
      <div class="relative">
        <p
          class="text-xs md:text-sm text-white/95 leading-relaxed font-sans drop-shadow-[0_2px_4px_rgba(0,0,0,0.9)] whitespace-pre-wrap transition-all"
          :class="isExpanded ? 'max-h-60 overflow-y-auto pr-2' : 'line-clamp-2'"
        >
          {{ displayText }}
        </p>

        <!-- 展開 / 折りたたみトグル -->
        <button
          v-if="displayText && displayText.length > 50"
          @click="isExpanded = !isExpanded"
          class="text-xs text-amber-300 hover:text-amber-200 font-bold drop-shadow-[0_1px_3px_rgba(0,0,0,0.9)] mt-1 inline-flex items-center gap-0.5 cursor-pointer"
        >
          {{ isExpanded ? '折りたたむ' : '展開...' }}
        </button>
      </div>

      <!-- 言語切り替えピル (半透明ミニマル) -->
      <div class="flex items-center gap-1.5 pt-1">
        <button
          v-for="l in [{ id: 'original', label: '原文' }, { id: 'ja', label: 'JA' }, { id: 'en', label: 'EN' }, { id: 'zh', label: 'ZH' }]"
          :key="l.id"
          @click="selectedLang = l.id as any"
          class="px-2 py-0.5 rounded-full text-[10px] font-mono transition-all backdrop-blur-md cursor-pointer drop-shadow-sm"
          :class="selectedLang === l.id ? 'bg-white/30 text-white font-bold border border-white/40 shadow-sm' : 'bg-black/40 text-white/70 hover:bg-black/60 border border-white/10'"
        >
          {{ l.label }}
        </button>
      </div>
    </div>

    <!-- 右側: TikTok/Douyin風 縦並びアクション列 -->
    <div class="flex flex-col items-center gap-4 pb-2 pointer-events-auto flex-shrink-0">
      <!-- アバター (右側アクション列トップ) -->
      <div class="relative group/avatar cursor-pointer">
        <Avatar
          :avatar-url="article.author.avatar_url"
          :handle="article.author.handle"
          :author="article.author"
          class="w-11 h-11 rounded-full border-2 border-white/80 shadow-lg"
        />
      </div>

      <!-- いいねボタン -->
      <button
        @click="$emit('toggleLike', article.id)"
        class="flex flex-col items-center gap-1 group/like cursor-pointer transition-transform active:scale-125"
      >
        <div
          class="w-10 h-10 rounded-full flex items-center justify-center text-xl transition-colors backdrop-blur-md shadow-lg"
          :class="article.is_liked ? 'bg-rose-600/80 text-white' : 'bg-black/40 hover:bg-black/60 text-white border border-white/20'"
        >
          {{ article.is_liked ? '❤️' : '🤍' }}
        </div>
        <span class="text-[10px] text-white/90 font-mono drop-shadow-[0_1px_3px_rgba(0,0,0,0.9)]">
          {{ article.metrics?.like_count || (article.is_liked ? 1 : 0) }}
        </span>
      </button>

      <!-- Stash導線ボタン -->
      <a
        :href="stashDirectUrl"
        target="_blank"
        rel="noopener noreferrer"
        title="Stash で開く"
        class="flex flex-col items-center gap-1 group/stash cursor-pointer"
      >
        <div class="w-10 h-10 rounded-full bg-black/40 hover:bg-purple-600/80 text-white border border-white/20 flex items-center justify-center text-lg transition-colors backdrop-blur-md shadow-lg">
          📦
        </div>
        <span class="text-[10px] text-white/90 font-mono drop-shadow-[0_1px_3px_rgba(0,0,0,0.9)]">
          Stash ↗
        </span>
      </a>
    </div>
  </div>
</template>
