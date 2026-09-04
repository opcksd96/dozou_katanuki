<!-- frontend/src/components/media/MediaOverlayBottomCard.vue (100行以下) -->
<script setup lang="ts">
import { ref, computed } from 'vue';
import type { RenderMedia, RenderTree } from '../../models/RenderTree';
import type { LanguageCode } from '../../composables/useTimeline';
import { useStashResolver } from '../../composables/useStashResolver';
import Avatar from '../article/Avatar.vue';

const props = defineProps<{ media: RenderMedia; article: RenderTree; targetLang?: LanguageCode }>();
defineEmits<{ (e: 'toggleLike', id: string): void }>();

const isExpanded = ref(false), selectedLang = ref<LanguageCode>(props.targetLang || 'ja');
const displayText = computed(() => {
  const c = props.article.content;
  if (!c) return '';
  if (selectedLang.value === 'ja' && c.ja) return c.ja;
  if (selectedLang.value === 'en' && c.en) return c.en;
  if (selectedLang.value === 'zh' && c.zh) return c.zh;
  return c.original;
});

const { getStashSceneUrl, getStashImageUrl, stashBaseUrl } = useStashResolver();
const stashDirectUrl = computed(() => {
  if (props.media.stash_scene_id) return getStashSceneUrl(props.media.stash_scene_id);
  if (props.media.stash_image_id) return getStashImageUrl(props.media.stash_image_id);
  return stashBaseUrl.value;
});
</script>

<template>
  <div class="w-full flex items-end justify-between gap-4 p-4 md:p-8 select-text pointer-events-none">
    <!-- 左側: 文字情報オーバーレイ -->
    <div class="flex-1 max-w-xl space-y-2 pointer-events-auto text-left">
      <div class="flex items-center gap-2.5">
        <Avatar :avatar-url="article.author.avatar_url" :handle="article.author.handle" :author="article.author" class="w-8 h-8 md:w-9 md:h-9 rounded-full shadow-lg border border-white/40" />
        <div class="flex items-center gap-2 leading-tight">
          <span class="text-base font-bold text-white drop-shadow-[0_2px_4px_rgba(0,0,0,0.9)]">@{{ article.author.display_name || article.author.handle }}</span>
          <span class="text-xs text-white/70 font-mono drop-shadow-[0_1px_3px_rgba(0,0,0,0.8)]">@{{ article.author.handle }}</span>
        </div>
      </div>

      <div class="relative">
        <p class="text-xs md:text-sm text-white/95 leading-relaxed font-sans drop-shadow-[0_2px_4px_rgba(0,0,0,0.9)] whitespace-pre-wrap" :class="isExpanded ? 'max-h-60 overflow-y-auto pr-2' : 'line-clamp-2'">
          {{ displayText }}
        </p>
        <button v-if="displayText && displayText.length > 50" @click="isExpanded = !isExpanded" class="text-xs text-amber-300 font-bold drop-shadow-[0_1px_3px_rgba(0,0,0,0.9)] mt-1 cursor-pointer">
          {{ isExpanded ? '折りたたむ' : '展開...' }}
        </button>
      </div>

      <div class="flex items-center gap-1.5 pt-1">
        <button v-for="l in [{ id: 'original', l: '原文' }, { id: 'ja', l: 'JA' }, { id: 'en', l: 'EN' }, { id: 'zh', l: 'ZH' }]" :key="l.id" @click="selectedLang = l.id as any" class="px-2 py-0.5 rounded-full text-[10px] font-mono backdrop-blur-md cursor-pointer" :class="selectedLang === l.id ? 'bg-white/30 text-white font-bold border border-white/40' : 'bg-black/40 text-white/70'">
          {{ l.l }}
        </button>
      </div>
    </div>

    <!-- 右側: アクション列 -->
    <div class="flex flex-col items-center gap-4 pb-2 pointer-events-auto shrink-0">
      <Avatar :avatar-url="article.author.avatar_url" :handle="article.author.handle" :author="article.author" class="w-11 h-11 rounded-full border-2 border-white/80 shadow-lg" />
      <button @click="$emit('toggleLike', article.id)" class="flex flex-col items-center gap-1 cursor-pointer active:scale-125 transition-transform">
        <div class="w-10 h-10 rounded-full flex items-center justify-center text-xl backdrop-blur-md shadow-lg" :class="article.is_liked ? 'bg-rose-600/80 text-white' : 'bg-black/40 text-white border border-white/20'">{{ article.is_liked ? '❤️' : '🤍' }}</div>
        <span class="text-[10px] text-white font-mono">{{ article.metrics?.like_count || (article.is_liked ? 1 : 0) }}</span>
      </button>
      <a :href="stashDirectUrl" target="_blank" rel="noopener noreferrer" class="flex flex-col items-center gap-1 cursor-pointer">
        <div class="w-10 h-10 rounded-full bg-black/40 hover:bg-purple-600/80 text-white border border-white/20 flex items-center justify-center text-lg backdrop-blur-md shadow-lg">📦</div>
        <span class="text-[10px] text-white font-mono">Stash ↗</span>
      </a>
    </div>
  </div>
</template>
