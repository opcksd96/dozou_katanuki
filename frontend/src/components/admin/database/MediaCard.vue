<!-- frontend/src/components/admin/database/MediaCard.vue (100行以下) -->
<script setup lang="ts">
import { ref, computed } from 'vue';

const props = defineProps<{
  media: any;
}>();

const emit = defineEmits<{
  (e: 'click', m: any): void;
  (e: 'retry', mediaId: string): void;
}>();

const isHovered = ref(false);
const imgFailed = ref(false);

const isVideo = computed(() => {
  const t = props.media.type?.toLowerCase();
  return t === 'video' || t === 'gif' || t === 'animated_gif' || !!props.media.stash_scene_id;
});

const thumbnailUrl = computed(() => {
  if (props.media.stash_scene_id) return `/stash-proxy/scene/${props.media.stash_scene_id}/screenshot`;
  if (props.media.stash_image_id) return `/stash-proxy/image/${props.media.stash_image_id}/thumbnail`;
  return props.media.download_url || '';
});

const previewVideoUrl = computed(() => {
  if (props.media.stash_scene_id) return `/stash-proxy/scene/${props.media.stash_scene_id}/preview`;
  return '';
});

const handleImgError = () => { imgFailed.value = true; };
</script>

<template>
  <div 
    class="bg-slate-950 border border-slate-800 rounded-xl p-2.5 flex flex-col space-y-2 group shadow hover:border-slate-600 transition-all cursor-pointer"
    @click="emit('click', media)"
    @mouseenter="isHovered = true"
    @mouseleave="isHovered = false"
  >
    <!-- プレビュー枠 -->
    <div class="h-28 bg-slate-900 rounded-lg overflow-hidden flex items-center justify-center relative select-none">
      <!-- ホバー時動画プレビュー (Stash連携動画) -->
      <video
        v-if="isVideo && previewVideoUrl && isHovered"
        :src="previewVideoUrl"
        autoplay
        muted
        loop
        playsinline
        class="w-full h-full object-cover"
      />
      <!-- 通常時: サムネイル画像 -->
      <img
        v-else-if="thumbnailUrl && !imgFailed"
        :src="thumbnailUrl"
        :alt="media.media_id"
        class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300"
        loading="lazy"
        @error="handleImgError"
      />
      <!-- 画像取得不可/エラー時フォールバック -->
      <div v-else class="flex flex-col items-center justify-center text-slate-500 gap-1 text-[11px]">
        <span class="text-lg">{{ isVideo ? '🎬' : '🖼️' }}</span>
        <span class="text-[9px] font-mono">{{ isVideo ? 'VIDEO' : 'IMAGE' }}</span>
      </div>

      <!-- タイプバッジ (動画/GIF) -->
      <span v-if="isVideo" class="absolute bottom-1.5 left-1.5 px-1.5 py-0.5 rounded bg-black/75 text-slate-200 text-[9px] font-mono font-bold flex items-center gap-1 backdrop-blur-xs">
        ▶ {{ media.type?.toUpperCase() || 'VIDEO' }}
      </span>

      <!-- ダウンロードステータスバッジ -->
      <span 
        class="absolute top-1.5 right-1.5 px-1.5 py-0.5 rounded text-[9px] font-mono font-bold backdrop-blur-xs" 
        :class="{
          'bg-emerald-950/90 text-emerald-300 border border-emerald-600/70': media.download_status === 'COMPLETED',
          'bg-amber-950/90 text-amber-300 border border-amber-600/70': media.download_status === 'QUEUED',
          'bg-rose-950/90 text-rose-300 border border-rose-600/70': media.download_status === 'DEAD_404' || media.failed_reason
        }"
      >
        {{ media.download_status }}
      </span>
    </div>

    <!-- メディア情報 -->
    <div class="text-[11px] font-mono space-y-0.5 flex-1 min-w-0">
      <div class="text-slate-200 font-bold truncate" :title="media.media_id">{{ media.media_id }}</div>
      <div class="text-slate-400 text-[10px] flex items-center justify-between">
        <span class="truncate">@{{ media.username }}</span>
        <span v-if="media.width && media.height" class="text-slate-500 shrink-0">{{ media.width }}x{{ media.height }}</span>
      </div>
      <div v-if="media.failed_reason" class="text-[10px] text-rose-400 truncate" :title="media.failed_reason">
        ⚠️ {{ media.failed_reason }}
      </div>
    </div>

    <!-- Stash 動線 ＆ アクション -->
    <div class="pt-1.5 border-t border-slate-850 flex items-center justify-between text-[10px] font-mono" @click.stop>
      <span v-if="media.stash_scene_id || media.stash_image_id" class="text-emerald-400 truncate max-w-[130px]" :title="media.stash_scene_id || media.stash_image_id">
        🎛️ Stash: {{ (media.stash_scene_id || media.stash_image_id)?.slice(0, 8) }}...
      </span>
      <span v-else class="text-slate-500">Stash未連携</span>
      <button @click="emit('retry', media.media_id)" class="px-2 py-0.5 bg-slate-800 hover:bg-slate-700 text-blue-400 rounded text-[10px] transition-colors">
        再取得
      </button>
    </div>
  </div>
</template>
