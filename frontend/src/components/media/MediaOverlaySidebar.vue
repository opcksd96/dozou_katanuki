<!-- frontend/src/components/media/MediaOverlaySidebar.vue (100行以下) -->
<script setup lang="ts">
import { ref, computed } from 'vue';
import type { RenderMedia, RenderTree } from '../../models/RenderTree';
import type { LanguageCode } from '../../composables/useTimeline';

const props = defineProps<{
  media: RenderMedia;
  article: RenderTree;
  targetLang?: LanguageCode;
}>();
defineEmits<{ (e: 'toggleLike', id: string): void }>();

const selectedLang = ref<LanguageCode>(props.targetLang || 'ja');

const displayText = computed(() => {
  const c = props.article.content;
  if (!c) return '';
  if (selectedLang.value === 'ja' && c.ja) return c.ja;
  if (selectedLang.value === 'en' && c.en) return c.en;
  if (selectedLang.value === 'zh' && c.zh) return c.zh;
  return c.original;
});

const stashDirectUrl = computed(() => {
  if (props.media.stash_scene_id) return `http://127.0.0.1:9999/scenes/${props.media.stash_scene_id}`;
  if (props.media.stash_image_id) return `http://127.0.0.1:9999/images/${props.media.stash_image_id}`;
  return 'http://127.0.0.1:9999';
});
</script>

<template>
  <div class="w-full md:w-80 bg-slate-900 border-t md:border-t-0 md:border-l border-slate-800 p-5 flex flex-col justify-between overflow-y-auto max-h-[40vh] md:max-h-full">
    <div class="space-y-4">
      <div class="flex items-center gap-3">
        <div class="w-10 h-10 rounded-full bg-slate-800 border border-slate-700 flex items-center justify-center font-bold text-slate-200">
          {{ article.author.display_name?.charAt(0) || 'A' }}
        </div>
        <div>
          <div class="text-sm font-bold text-slate-100">{{ article.author.display_name }}</div>
          <div class="text-xs text-slate-400 font-mono">@{{ article.author.handle }}</div>
        </div>
      </div>

      <!-- 言語タブ -->
      <div class="flex gap-1 bg-slate-950 p-1 rounded-lg text-[10px] font-mono">
        <button v-for="l in [{ id: 'original', label: '原文' }, { id: 'ja', label: '🇯🇵 JA' }, { id: 'en', label: '🇺🇸 EN' }, { id: 'zh', label: '🇨🇳 ZH' }]" :key="l.id" @click="selectedLang = l.id as any" class="flex-1 py-1 rounded transition-colors" :class="selectedLang === l.id ? 'bg-blue-600 text-white font-bold' : 'text-slate-400'">
          {{ l.label }}
        </button>
      </div>

      <div class="text-xs text-slate-200 leading-relaxed font-sans select-text whitespace-pre-wrap">{{ displayText }}</div>
    </div>

    <!-- フッターアクション -->
    <div class="pt-4 border-t border-slate-800 flex items-center justify-between">
      <a :href="stashDirectUrl" target="_blank" class="text-xs text-blue-400 hover:underline flex items-center gap-1 font-mono">
        <span>🎛️</span> Stashで開く ↗
      </a>
      <button @click="$emit('toggleLike', article.id)" class="text-lg transition-transform active:scale-125" :class="article.is_liked ? 'text-rose-500' : 'text-slate-500'">
        {{ article.is_liked ? '❤️' : '🤍' }}
      </button>
    </div>
  </div>
</template>
