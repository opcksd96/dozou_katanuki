<!-- frontend/src/components/admin/database/MediaPreviewModal.vue (100行以下) -->
<script setup lang="ts">
import { computed, onMounted, onUnmounted } from 'vue';
import StashPlayer from '../../media/StashPlayer.vue';

const props = defineProps<{
  media: any;
  hasPrev?: boolean;
  hasNext?: boolean;
}>();

const emit = defineEmits<{
  (e: 'close'): void;
  (e: 'prev'): void;
  (e: 'next'): void;
  (e: 'retry', mediaId: string): void;
  (e: 'purge', mediaId: string): void;
  (e: 'viewPost', articleId: string): void;
}>();

const isVideo = computed(() => {
  const t = props.media.type?.toLowerCase();
  return t === 'video' || t === 'gif' || t === 'animated_gif' || !!props.media.stash_scene_id;
});

const handleKeyDown = (e: KeyboardEvent) => {
  if (e.key === 'Escape') { e.preventDefault(); emit('close'); }
  else if ((e.key === 'ArrowLeft' || e.key === 'a' || e.key === 'A') && props.hasPrev) { e.preventDefault(); emit('prev'); }
  else if ((e.key === 'ArrowRight' || e.key === 'd' || e.key === 'D') && props.hasNext) { e.preventDefault(); emit('next'); }
};

onMounted(() => window.addEventListener('keydown', handleKeyDown));
onUnmounted(() => window.removeEventListener('keydown', handleKeyDown));
</script>

<template>
  <div class="fixed inset-0 z-50 flex items-center justify-center bg-black/85 backdrop-blur-md p-2 md:p-6 select-none" @click.self="emit('close')">
    <!-- 前へボタン -->
    <button v-if="hasPrev" @click.stop="emit('prev')" class="absolute left-3 z-30 w-11 h-11 rounded-full bg-black/60 hover:bg-black/90 text-white flex items-center justify-center text-xl transition-all shadow-lg backdrop-blur-sm border border-white/10 active:scale-95">‹</button>

    <!-- モーダル本体 -->
    <div class="bg-slate-950/95 border border-slate-700/80 rounded-2xl w-full max-w-6xl max-h-[92vh] flex flex-col overflow-hidden shadow-2xl">
      <!-- ヘッダー -->
      <div class="px-4 py-2.5 border-b border-slate-800 flex items-center justify-between bg-slate-900/70 shrink-0">
        <div class="flex items-center gap-2 font-mono text-xs text-slate-200 min-w-0">
          <span class="px-2 py-0.5 rounded bg-blue-950 text-blue-300 font-bold border border-blue-800 shrink-0">{{ media.type?.toUpperCase() || 'MEDIA' }}</span>
          <span class="font-bold truncate max-w-md">{{ media.media_id || media.id }}</span>
          <span v-if="media.width && media.height" class="text-slate-400 text-[11px] shrink-0">({{ media.width }}x{{ media.height }})</span>
        </div>
        <button @click="emit('close')" class="w-7 h-7 rounded-lg bg-slate-800 hover:bg-slate-700 text-slate-400 hover:text-white flex items-center justify-center text-sm font-mono transition-colors">✕</button>
      </div>

      <!-- メディアプレビュー本体 (アスペクト比完全維持・画面内フィット) -->
      <div class="flex-1 min-h-[320px] max-h-[70vh] bg-black flex items-center justify-center p-2 relative overflow-hidden" @click.self="emit('close')">
        <!-- 動画プレイヤー (StashPlayer / HLS / 直接再生) -->
        <div v-if="isVideo && media.urls?.stream" class="w-full h-full max-h-full flex items-center justify-center">
          <StashPlayer
            :src="media.urls.stream"
            :poster="media.urls.thumbnail"
            :stashSceneId="media.stash_scene_id"
            :autoplay="true"
            :show-expand-button="false"
            class="max-w-full max-h-full"
          />
        </div>
        <!-- 画像ビューア (object-contain で自然なアスペクト比フィット) -->
        <img
          v-else-if="media.urls?.image || media.urls?.thumbnail"
          :src="media.urls.image || media.urls.thumbnail"
          :alt="media.media_id || media.id"
          class="max-w-full max-h-full object-contain rounded shadow-2xl select-none mx-auto"
        />
        <div v-else class="text-slate-500 font-mono text-xs">メディア実体を表示できません</div>
      </div>

      <!-- 詳細情報・フッター -->
      <div class="p-3 border-t border-slate-800 bg-slate-900/60 text-xs font-mono flex flex-wrap items-center justify-between gap-2 shrink-0">
        <div class="flex items-center gap-4 text-[11px]">
          <div><span class="text-slate-500">アカウント:</span> <span class="text-slate-200 font-bold">@{{ media.username }}</span></div>
          <div><span class="text-slate-500">状態:</span> <span :class="media.download_status === 'COMPLETED' ? 'text-emerald-400 font-bold' : 'text-amber-400 font-bold'">{{ media.download_status }}</span></div>
          <div v-if="media.stash_scene_id || media.stash_image_id" class="text-emerald-400">
            🎛️ Stash ID: {{ media.stash_scene_id || media.stash_image_id }}
          </div>
        </div>
        <div class="flex items-center gap-2">
          <button v-if="media.article_id" @click="emit('viewPost', media.article_id)" class="px-3 py-1 bg-slate-800 hover:bg-slate-700 text-blue-300 font-bold rounded-lg transition-colors flex items-center gap-1">
            📝 親記事を見る
          </button>
          <button @click="emit('retry', media.media_id || media.id)" class="px-3 py-1 bg-blue-600 hover:bg-blue-500 text-white font-bold rounded-lg transition-colors">
            再取得
          </button>
          <button @click="emit('purge', media.media_id || media.id)" class="px-3 py-1 bg-rose-950 hover:bg-rose-800 text-rose-200 font-bold rounded-lg border border-rose-700/60 transition-colors">
            🗑️ パージ
          </button>
          <button @click="emit('close')" class="px-3 py-1 bg-slate-800 hover:bg-slate-700 text-slate-300 rounded-lg transition-colors">
            閉じる
          </button>
        </div>
      </div>
    </div>

    <!-- 次へボタン -->
    <button v-if="hasNext" @click.stop="emit('next')" class="absolute right-3 z-30 w-11 h-11 rounded-full bg-black/60 hover:bg-black/90 text-white flex items-center justify-center text-xl transition-all shadow-lg backdrop-blur-sm border border-white/10 active:scale-95">›</button>
  </div>
</template>
