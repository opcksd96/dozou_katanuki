<!-- frontend/src/components/admin/database/MediaCard.vue (100行以下 - SPEC-PRINCIPLE-001) -->
<script setup lang="ts">
import { ref, computed } from 'vue';
import { BrowserOpenURL } from '../../../../wailsjs/runtime/runtime';
import Avatar from '../../article/Avatar.vue';

const props = withDefaults(
  defineProps<{
    media: any;
    compact?: boolean;
  }>(),
  { compact: false }
);

const emit = defineEmits<{
  (e: 'click', m: any): void;
  (e: 'retry', mediaId: string): void;
  (e: 'purge', mediaId: string): void;
  (e: 'viewPost', articleId: string): void;
  (e: 'openExplorer', mediaId: string): void;
  (e: 'openDefault', mediaId: string): void;
  (e: 'copy', media: any): void;
  (e: 'toggleBookmark', mediaId: string): void;
}>();

const isHovered = ref(false);
const imgFailed = ref(false);

const isVideo = computed(() => {
  const t = props.media.type?.toLowerCase();
  return t === 'video' || t === 'gif' || t === 'animated_gif' || !!props.media.stash_scene_id;
});

const isExcluded = computed(() => {
  return props.media.download_status === 'EXCLUDED' || props.media.raw_status === 'EXCLUDED' || props.media.failed_reason?.includes('Whitelist外');
});

const formattedTitle = computed(() => {
  if (props.media.title) return props.media.title;
  return `X (@${props.media.username || 'unknown'}): Tweet ${props.media.article_id || props.media.media_id || props.media.id}`;
});

const currentHostname = computed(() => (typeof window !== 'undefined' && window.location?.hostname) ? window.location.hostname : '127.0.0.1');

const stashDirectUrl = computed(() => {
  if (props.media.stash_scene_id) return `http://${currentHostname.value}:9999/scenes/${props.media.stash_scene_id}`;
  if (props.media.stash_image_id) return `http://${currentHostname.value}:9999/images/${props.media.stash_image_id}`;
  return null;
});

const openStash = () => {
  if (!stashDirectUrl.value) return;
  try { BrowserOpenURL(stashDirectUrl.value); } catch { window.open(stashDirectUrl.value, '_blank', 'noopener,noreferrer'); }
};
</script>

<template>
  <div 
    class="bg-slate-900/95 rounded-xl p-3 flex flex-col space-y-2.5 group shadow-lg hover:shadow-2xl transition-all cursor-pointer relative border border-slate-800/70 hover:border-slate-600"
    @click="emit('click', media)"
    @mouseenter="isHovered = true"
    @mouseleave="isHovered = false"
  >
    <!-- プレビュー枠 (Stash風のダイナミック＆アスペクト比対応エリア) -->
    <div :class="[compact ? 'h-36' : 'h-56', 'bg-black/80 rounded-lg overflow-hidden flex items-center justify-center relative select-none']">
      <video
        v-if="isVideo && media.urls?.preview && isHovered"
        :src="media.urls.preview"
        autoplay
        muted
        loop
        playsinline
        class="w-full h-full object-contain"
      />
      <img
        v-else-if="media.urls?.thumbnail && !imgFailed"
        :src="media.urls.thumbnail"
        :alt="media.media_id || media.id"
        class="w-full h-full object-contain group-hover:scale-105 transition-transform duration-300"
        loading="lazy"
        @error="imgFailed = true"
      />
      <div v-else class="flex flex-col items-center justify-center text-slate-500 gap-1.5 text-xs">
        <span :class="compact ? 'text-2xl' : 'text-4xl'">{{ isVideo ? '🎬' : '🖼️' }}</span>
        <span class="text-xs font-mono text-slate-400 font-bold">{{ isVideo ? 'VIDEO' : 'IMAGE' }}</span>
      </div>

      <!-- ブックマーク星印 -->
      <button 
        @click.stop="emit('toggleBookmark', media.media_id || media.id)" 
        class="absolute top-2 left-2 p-1.5 rounded-lg bg-black/80 hover:bg-black text-sm transition-colors backdrop-blur-xs"
        :title="media.is_bookmarked ? 'ブックマーク解除' : 'ブックマークに追加'"
      >
        {{ media.is_bookmarked ? '⭐' : '☆' }}
      </button>

      <!-- ステータスバッジ (12px基準) -->
      <span 
        v-if="isExcluded"
        class="absolute top-2 right-2 px-2.5 py-0.5 rounded text-xs font-mono font-bold bg-slate-800/90 text-slate-300 border border-slate-700"
      >
        EXCLUDED
      </span>
      <span 
        v-else
        class="absolute top-2 right-2 px-2.5 py-0.5 rounded text-xs font-mono font-bold backdrop-blur-xs shadow" 
        :class="{
          'bg-emerald-950/90 text-emerald-300 border border-emerald-700': media.download_status === 'COMPLETED',
          'bg-amber-950/90 text-amber-300 border border-amber-700': media.download_status === 'QUEUED',
          'bg-rose-950/90 text-rose-300 border border-rose-700': media.download_status === 'DEAD_404' || media.failed_reason
        }"
      >
        {{ media.download_status }}
      </span>

      <!-- 動画バッジ -->
      <span v-if="isVideo" class="absolute bottom-2 left-2 px-2.5 py-0.5 rounded bg-black/85 text-slate-200 text-xs font-mono font-bold flex items-center gap-1 backdrop-blur-xs">
        ▶ VIDEO
      </span>
    </div>

    <!-- メディア詳細情報 (Stashスタイルのカード情報) -->
    <div class="space-y-1.5 flex-1 min-w-0 font-sans">
      <!-- タイトル (Stash標準形式: X (@user): Tweet ID) - 13px太字 -->
      <div class="text-xs sm:text-sm font-bold text-slate-100 truncate hover:text-blue-400 transition-colors font-mono" :title="formattedTitle">
        {{ formattedTitle }}
      </div>

      <!-- 日付 ＆ ファイルメタデータ (12px) -->
      <div class="flex items-center justify-between text-xs text-slate-400 font-mono">
        <span v-if="media.tweet_date" class="text-slate-300 font-semibold">{{ media.tweet_date }}</span>
        <span v-if="media.width && media.height" class="text-slate-400">{{ media.width }}x{{ media.height }}</span>
      </div>

      <!-- 詳細テキスト表示 (ツイート本文プレビュー: 3行自然折り返し、下側切れ防止) -->
      <div 
        v-if="!compact && (media.full_text || media.full_text_ja)" 
        class="text-xs text-slate-200 line-clamp-3 leading-relaxed bg-slate-950/80 p-2 rounded-lg border border-slate-800/80 select-text overflow-hidden"
        :title="media.full_text_ja || media.full_text"
      >
        {{ media.full_text_ja || media.full_text }}
      </div>

      <!-- アカウント / パフォーマー情報 (12px) -->
      <div class="flex items-center justify-between text-xs text-slate-300 pt-0.5 font-mono">
        <div class="flex items-center gap-1.5 min-w-0 truncate">
          <Avatar :avatar-url="media.avatar_url" :handle="media.username" size-class="w-4 h-4" />
          <span class="text-slate-200 font-semibold truncate">@{{ media.username }}</span>
        </div>
        <span class="text-slate-400 text-xs truncate max-w-[120px]" :title="media.media_id">{{ media.media_id }}</span>
      </div>

      <!-- エラー理由 (ある場合) -->
      <div v-if="media.failed_reason" class="text-xs text-rose-400 truncate bg-rose-950/50 p-1 rounded border border-rose-900/60 font-mono" :title="media.failed_reason">
        ⚠️ {{ media.failed_reason }}
      </div>
    </div>

    <!-- クイックアクションバー (統一アイコンスタイル) -->
    <div class="pt-2 border-t border-slate-800/80 flex items-center justify-between text-xs font-mono" @click.stop>
      <div class="flex items-center gap-1.5">
        <button @click.stop="emit('openExplorer', media.media_id || media.id)" class="p-1.5 bg-slate-800 hover:bg-slate-700 text-slate-200 rounded-lg transition-colors" title="エクスプローラーで開く">📂</button>
        <button @click.stop="emit('openDefault', media.media_id || media.id)" class="p-1.5 bg-slate-800 hover:bg-slate-700 text-slate-200 rounded-lg transition-colors" title="既定アプリで開く">🚀</button>
        <button @click.stop="emit('copy', media)" class="p-1.5 bg-slate-800 hover:bg-slate-700 text-slate-200 rounded-lg transition-colors" title="クリップボードにコピー">📋</button>
        <button v-if="stashDirectUrl" @click.stop="openStash" class="p-1.5 bg-purple-950/90 hover:bg-purple-900 border border-purple-700/60 text-purple-300 rounded-lg flex items-center gap-1 transition-colors" title="Stash WebUIで開く">
          <span>🎛️</span>
        </button>
      </div>
      <div class="flex items-center gap-1.5">
        <button @click="emit('retry', media.media_id || media.id)" class="p-1.5 bg-blue-950/80 hover:bg-blue-900 border border-blue-700/60 text-blue-300 rounded-lg transition-colors" title="再取得 (リトライ)">
          <span>🔄</span>
        </button>
        <button @click="emit('purge', media.media_id || media.id)" class="p-1.5 bg-rose-950/80 hover:bg-rose-900 border border-rose-700/60 text-rose-300 rounded-lg transition-colors" title="DBからパージ (削除)">
          <span>🗑️</span>
        </button>
      </div>
    </div>
  </div>
</template>
