<script setup lang="ts">
import type { RenderMedia } from '../../models/RenderTree';
import MediaFiller from './MediaFiller.vue';
import InlineVideoPlayer from './InlineVideoPlayer.vue';

defineProps<{
  media: RenderMedia[];
}>();

const emit = defineEmits<{
  (e: 'retry', mediaId: string): void;
  (e: 'clickMedia', media: RenderMedia, list: RenderMedia[]): void;
}>();

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
        <template v-if="item.download_status === 'COMPLETED' && ((item.type === 'video' && item.urls.stream) || (item.type !== 'video' && (item.urls.image || item.urls.thumbnail)))">
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
          <!-- 動画: 軽量HLS/HTML5プレイヤーで直接再生可能 ＆ 全画面詳細ボタンでモーダル展開 -->
          <div
            v-else-if="item.type === 'video'"
            @click.stop
            :class="[
              'w-full flex items-center justify-center bg-black',
              media.length === 1 ? 'max-h-[580px]' : 'h-full'
            ]"
          >
            <InlineVideoPlayer
              :src="item.urls.stream"
              :poster="item.urls.thumbnail"
              :stashSceneId="item.stash_scene_id"
              @expand="emit('clickMedia', item, media)"
            />
          </div>
        </template>
        <!-- 未完了/実体未紐付け/失敗時: SVGプレースホルダー・フィラー -->
        <MediaFiller v-else :media="item" class="w-full h-full" @click.stop @retry="(id) => emit('retry', id)" />
      </div>
    </div>
  </div>
</template>


