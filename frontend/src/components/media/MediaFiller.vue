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
    <!-- 1. 画像用 SVG フィラー (カメラ / フォトフレーム) -->
    <svg v-if="media.type === 'image'" class="w-10 h-10 text-slate-500 mb-2 drop-shadow" fill="none" stroke="currentColor" viewBox="0 0 24 24">
      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M3 9a2 2 0 012-2h.93a2 2 0 001.664-.89l.812-1.22A2 2 0 0110.07 4h3.86a2 2 0 011.664.89l.812 1.22A2 2 0 0018.07 7H19a2 2 0 012 2v9a2 2 0 01-2 2H5a2 2 0 01-2-2V9z" />
      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M15 13a3 3 0 11-6 0 3 3 0 016 0z" />
    </svg>

    <!-- 2. 動画用 SVG フィラー (Play シンボル) -->
    <svg v-else-if="media.type === 'video'" class="w-12 h-12 text-slate-500 mb-2 drop-shadow" fill="currentColor" viewBox="0 0 24 24">
      <path d="M8 5v14l11-7z"/>
    </svg>

    <!-- 3. GIF用 ピルバッジ -->
    <div v-else-if="media.type === 'gif'" class="bg-indigo-600 text-white font-bold text-xs px-2.5 py-0.5 rounded-full font-mono mb-2 shadow">GIF</div>

    <!-- ステータスオーバーレイ ＆ リトライボタン (半透明ブラーシート) -->
    <div class="z-10 flex flex-col items-center gap-1 bg-slate-900/70 backdrop-blur-sm px-3 py-1.5 rounded-lg border border-slate-800/80">
      <span class="text-[11px] font-mono font-semibold text-amber-300">
        [{{ media.download_status }}]
      </span>
      <p v-if="media.failed_reason" class="text-[10px] text-slate-400 max-w-xs truncate">{{ media.failed_reason }}</p>
      <button
        v-if="['DEAD_404', 'OUTSOURCED', 'RETAINED'].includes(media.download_status)"
        @click="emit('retry', media.id)"
        class="mt-1 text-xs text-blue-400 hover:text-blue-300 underline font-mono cursor-pointer flex items-center gap-1"
      >
        <svg class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
        </svg>
        再試行 (Retry)
      </button>
    </div>

    <!-- 動画下部シークバープレースホルダーライン -->
    <div v-if="media.type === 'video'" class="absolute bottom-0 left-0 right-0 h-1 bg-slate-800">
      <div class="h-full bg-slate-600 w-1/3"></div>
    </div>
  </div>
</template>
