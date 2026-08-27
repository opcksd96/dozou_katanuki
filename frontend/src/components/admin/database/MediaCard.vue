<!-- frontend/src/components/admin/database/MediaCard.vue (100行以下 - SPEC-PRINCIPLE-001) -->
<script setup lang="ts">
import { ref, computed } from 'vue';
import { BrowserOpenURL } from '../../../../wailsjs/runtime/runtime';
import Avatar from '../../article/Avatar.vue';

const props = withDefaults(defineProps<{ media: any; compact?: boolean }>(), { compact: false });
const emit = defineEmits<{
  (e: 'click', m: any): void; (e: 'retry', id: string): void; (e: 'purge', id: string): void;
  (e: 'openExplorer', id: string): void; (e: 'openDefault', id: string): void; (e: 'toggleBookmark', id: string): void;
  (e: 'viewPost', articleId: string): void; (e: 'viewPostTimeline', articleId: string): void;
}>();

const isHovered = ref(false), imgFailed = ref(false);
const isVideo = computed(() => { const t = props.media.type?.toLowerCase(); return t === 'video' || t === 'gif' || !!props.media.stash_scene_id; });
const formattedTitle = computed(() => props.media.title || `X (@${props.media.username || 'unknown'}): Tweet ${props.media.article_id || props.media.media_id || props.media.id}`);
const stashDirectUrl = computed(() => {
  const host = typeof window !== 'undefined' && window.location?.hostname ? window.location.hostname : '127.0.0.1';
  if (props.media.stash_scene_id) return `http://${host}:9999/scenes/${props.media.stash_scene_id}`;
  if (props.media.stash_image_id) return `http://${host}:9999/images/${props.media.stash_image_id}`;
  return null;
});

const openStash = () => {
  if (!stashDirectUrl.value) return;
  try { BrowserOpenURL(stashDirectUrl.value); } catch { window.open(stashDirectUrl.value, '_blank', 'noopener,noreferrer'); }
};
</script>

<template>
  <div class="bg-slate-900/95 rounded-xl p-3 flex flex-col space-y-2 group shadow-lg hover:shadow-2xl transition-all cursor-pointer border border-slate-800 hover:border-slate-600 select-none" @click="emit('click', media)" @mouseenter="isHovered = true" @mouseleave="isHovered = false">
    <!-- プレビュー枠 -->
    <div :class="[compact ? 'h-32' : 'h-48', 'bg-black/80 rounded-lg overflow-hidden flex items-center justify-center relative select-none']">
      <video v-if="isVideo && media.urls?.preview && isHovered" :src="media.urls.preview" autoplay muted loop playsinline class="w-full h-full object-contain" />
      <img v-else-if="media.urls?.thumbnail && !imgFailed" :src="media.urls.thumbnail" :alt="media.media_id" class="w-full h-full object-contain group-hover:scale-105 transition-transform" loading="lazy" @error="imgFailed = true" />
      <div v-else class="text-slate-500 text-xs font-mono font-bold">{{ isVideo ? '🎬 VIDEO' : '🖼️ IMAGE' }}</div>

      <button @click.stop="emit('toggleBookmark', media.media_id || media.id)" class="absolute top-2 left-2 p-1 rounded bg-black/80 text-xs cursor-pointer">{{ media.is_bookmarked ? '⭐' : '☆' }}</button>
      <span class="absolute top-2 right-2 px-2 py-0.5 rounded text-[10px] font-mono font-bold" :class="{
        'bg-emerald-950 text-emerald-300 border border-emerald-700/50': media.download_status === 'COMPLETED',
        'bg-purple-950 text-purple-300 border border-purple-700/50': media.download_status === 'OUTSOURCED',
        'bg-amber-950 text-amber-300 border border-amber-600/60': media.download_status === 'RETAINED',
        'bg-rose-950 text-rose-300 border border-rose-700/50': media.download_status === 'DEAD_404',
        'bg-red-950 text-red-300 border border-red-700/50': media.download_status === 'FAILED',
        'bg-slate-800 text-slate-300': media.download_status === 'QUEUED'
      }">{{ media.download_status }}</span>
      <div class="absolute bottom-2 right-2 flex items-center gap-1">
        <span v-if="media.media_quality" class="px-1.5 py-0.5 rounded bg-blue-950/90 text-blue-300 border border-blue-600/50 text-[9px] font-mono font-bold uppercase">{{ media.media_quality }}</span>
        <span v-if="media.width && media.height" class="px-2 py-0.5 rounded bg-black/80 text-slate-200 text-[10px] font-mono font-bold">{{ media.width }}x{{ media.height }}</span>
      </div>
      <span v-if="isVideo" class="absolute bottom-2 left-2 px-2 py-0.5 rounded bg-black/80 text-slate-200 text-[10px] font-mono font-bold">▶ VIDEO</span>
    </div>

    <!-- メディア詳細 -->
    <div class="space-y-1 flex-1 min-w-0 font-sans text-xs">
      <div class="font-bold text-slate-100 truncate font-mono" :title="formattedTitle">{{ formattedTitle }}</div>
      <div v-if="!compact && (media.full_text || media.full_text_ja)" class="text-slate-200 line-clamp-2 leading-relaxed bg-slate-950 p-1.5 rounded border border-slate-800 select-text">{{ media.full_text_ja || media.full_text }}</div>
      <div class="flex items-center justify-between text-slate-400 font-mono text-[11px] pt-1">
        <div class="flex items-center gap-1 truncate"><Avatar :avatar-url="media.avatar_url" :handle="media.username" size-class="w-4 h-4" /><span class="truncate">@{{ media.username }}</span></div>
        <span>{{ media.tweet_date || '-' }}</span>
      </div>
    </div>

    <!-- アクションボタン -->
    <div class="pt-1.5 border-t border-slate-800 flex items-center justify-between text-xs" @click.stop>
      <div class="flex gap-1">
        <button v-if="media.article_id" @click.stop="emit('viewPostTimeline', media.article_id)" class="px-2 py-0.5 bg-blue-600 hover:bg-blue-500 text-white font-bold rounded text-[11px] shadow-sm active:scale-95 cursor-pointer" title="タイムラインで該当ツイートを表示">📱 タイムライン</button>
        <button @click.stop="emit('openExplorer', media.media_id || media.id)" class="p-1 bg-slate-800 hover:bg-slate-700 rounded active:scale-95 cursor-pointer" title="フォルダ">📂</button>
        <button v-if="stashDirectUrl" @click.stop="openStash" class="p-1 bg-purple-950 hover:bg-purple-900 text-purple-300 rounded active:scale-95 cursor-pointer" title="Stash WebUI">🎛️</button>
      </div>
      <div class="flex gap-1">
        <button @click="emit('retry', media.media_id || media.id)" class="p-1 bg-blue-950 hover:bg-blue-900 text-blue-300 rounded active:scale-95 cursor-pointer" title="リトライ">🔄</button>
        <button @click="emit('purge', media.media_id || media.id)" class="p-1 bg-rose-950 hover:bg-rose-900 text-rose-300 rounded active:scale-95 cursor-pointer" title="削除">🗑️</button>
      </div>
    </div>
  </div>
</template>
