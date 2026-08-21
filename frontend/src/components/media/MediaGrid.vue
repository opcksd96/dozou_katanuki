<script setup lang="ts">
import { ref } from 'vue';
import type { RenderMedia } from '../../models/RenderTree';
import MediaFiller from './MediaFiller.vue';

defineProps<{
  media: RenderMedia[];
}>();

const emit = defineEmits<{
  (e: 'retry', mediaId: string): void;
  (e: 'clickMedia', media: RenderMedia, list: RenderMedia[]): void;
}>();

const hoveredId = ref<string | null>(null);

const handleImageError = (e: Event) => {
  const img = e.target as HTMLImageElement;
  if (!img || !img.src) return;
  try {
    const u = new URL(img.src, window.location.href);
    if (!u.pathname || u.pathname === '/' || u.pathname === '') return;
    const count = Number(img.dataset.retry || '0');
    if (count < 2) {
      img.dataset.retry = String(count + 1);
      setTimeout(() => {
        u.searchParams.set('_r', String(Date.now()));
        img.src = u.pathname + u.search;
      }, 1500 * (count + 1));
    }
  } catch (_) {}
};
</script>

<template>
  <div v-if="media && media.length > 0" class="media-grid-container mt-3 rounded-xl overflow-hidden border border-slate-800 bg-slate-950/80">
    <div :class="['grid gap-1', media.length === 1 ? 'grid-cols-1' : 'grid-cols-2']">
      <div
        v-for="item in media"
        :key="item.id"
        class="relative group flex items-center justify-center bg-black/40 overflow-hidden"
        :class="media.length > 1 ? 'h-52 md:h-64' : 'max-h-[580px]'"
      >
        <!-- ローカル確保完了時: 描画 (有効なローカルURLが存在する場合のみ) -->
        <template v-if="item.download_status === 'COMPLETED' && ((item.type === 'video' && (item.urls.stream || item.urls.preview || item.urls.thumbnail)) || (item.type !== 'video' && (item.urls.image || item.urls.thumbnail)))">
          <img
            v-if="(item.type === 'image' || item.type === 'gif') && (item.urls.image || item.urls.thumbnail)"
            :src="item.urls.image || item.urls.thumbnail"
            :alt="item.id"
            :class="[
              'cursor-pointer transition-opacity hover:opacity-95',
              media.length === 1
                ? 'w-full h-auto max-h-[580px] object-contain mx-auto'
                : 'w-full h-full object-cover'
            ]"
            loading="lazy"
            @error="handleImageError"
            @click.stop="emit('clickMedia', item, media)"
          />
          <div
            v-else-if="item.type === 'video'"
            @mouseenter="hoveredId = item.id"
            @mouseleave="hoveredId = null"
            @click.stop="emit('clickMedia', item, media)"
            :class="[
              'w-full relative flex items-center justify-center cursor-pointer group/video bg-black/60',
              media.length === 1 ? 'max-h-[580px] min-h-[220px]' : 'h-full'
            ]"
          >
            <!-- ホバープレビュー動画 (マウスオーバー時に自動再生) -->
            <video
              v-if="hoveredId === item.id && (item.urls.preview || item.urls.stream)"
              :src="item.urls.preview || item.urls.stream"
              autoplay
              loop
              muted
              playsinline
              preload="auto"
              class="w-full h-full object-cover"
            ></video>

            <!-- 通常時: 静止画サムネイル (サムネイルURLが存在する場合のみ) -->
            <img
              v-else-if="item.urls.thumbnail || item.urls.image"
              :src="item.urls.thumbnail || item.urls.image"
              :alt="item.id"
              class="w-full h-full object-cover transition-transform duration-300 group-hover/video:scale-[1.02]"
              loading="lazy"
              @error="handleImageError"
            />

            <!-- サムネイル不在時: スタイリッシュなダーク背景プレースホルダー -->
            <div
              v-else
              class="w-full h-full min-h-[200px] flex items-center justify-center bg-gradient-to-br from-slate-900 via-slate-950 to-purple-950/40"
            ></div>

            <!-- 再生バッジオーバーレイ (Twitter風) -->
            <div
              v-show="hoveredId !== item.id"
              class="absolute inset-0 flex items-center justify-center bg-black/20 group-hover/video:bg-black/40 transition-colors"
            >
              <div class="w-12 h-12 rounded-full bg-blue-600/90 group-hover/video:bg-blue-500 text-white flex items-center justify-center shadow-lg shadow-blue-600/40 transition-transform group-hover/video:scale-110">
                <span class="text-lg pl-0.5">▶</span>
              </div>
            </div>

            <!-- 動画バッジ -->
            <div class="absolute bottom-2 right-2 px-2 py-0.5 rounded bg-black/70 text-white text-[10px] font-mono border border-white/10">
              {{ hoveredId === item.id ? 'PREVIEW' : 'VIDEO' }}
            </div>
          </div>
        </template>
        <!-- 未完了/実体未紐付け/失敗時: SVGプレースホルダー・フィラー -->
        <MediaFiller v-else :media="item" class="w-full h-full" @click.stop @retry="(id) => emit('retry', id)" />
      </div>
    </div>
  </div>
</template>


