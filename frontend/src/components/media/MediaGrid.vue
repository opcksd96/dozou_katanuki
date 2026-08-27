<!-- frontend/src/components/media/MediaGrid.vue (100行以下 - SPEC-PRINCIPLE-001) -->
<script setup lang="ts">
import type { RenderMedia } from '../../models/RenderTree';
import MediaFiller from './MediaFiller.vue';
import InlineVideoPlayer from './InlineVideoPlayer.vue';

defineProps<{ media: RenderMedia[] }>();
const emit = defineEmits<{ (e: 'retry', mediaId: string): void; (e: 'clickMedia', media: RenderMedia, list: RenderMedia[]): void; }>();

const handleImageError = (e: Event) => {
  const img = e.target as HTMLImageElement;
  if (!img?.src) return;
  try {
    const u = new URL(img.src, window.location.href);
    if (!u.pathname || u.pathname === '/') return;
    const count = Number(img.dataset.retry || '0');
    if (count < 2) {
      img.dataset.retry = String(count + 1);
      setTimeout(() => { u.searchParams.set('_r', String(Date.now())); img.src = u.pathname + u.search; }, 1500 * (count + 1));
    }
  } catch (_) {}
};
</script>

<template>
  <div v-if="media && media.length > 0" class="media-grid-container mt-2.5 rounded-xl overflow-hidden border border-slate-800/80 bg-black/60">
    <div :class="['grid gap-1', media.length === 1 ? 'grid-cols-1' : (media.length === 3 ? 'grid-cols-2' : 'grid-cols-2')]">
      <div
        v-for="(item, idx) in media"
        :key="item.id"
        class="relative group flex items-center justify-center bg-black/80 overflow-hidden"
        :class="[
          media.length === 1 ? 'max-h-[560px]' : (media.length === 3 && idx === 0 ? 'row-span-2 min-h-[280px]' : 'h-48 sm:h-60')
        ]"
      >
        <template v-if="item.download_status === 'COMPLETED' && ((item.type === 'video' && item.urls.stream) || (item.type !== 'video' && (item.urls.image || item.urls.thumbnail)))">
          <img
            v-if="(item.type === 'image' || item.type === 'gif') && (item.urls.image || item.urls.thumbnail)"
            :src="item.urls.image || item.urls.thumbnail"
            :alt="item.id"
            :class="[
              'cursor-pointer transition-transform hover:scale-[1.01] w-full h-full object-contain mx-auto'
            ]"
            loading="lazy"
            @error="handleImageError"
            @click.stop="emit('clickMedia', item, media)"
          />
          <div v-else-if="item.type === 'video'" @click.stop :class="['w-full flex items-center justify-center bg-black', media.length === 1 ? 'max-h-[560px]' : 'h-full']">
            <InlineVideoPlayer :src="item.urls.stream" :poster="item.urls.thumbnail" :stashSceneId="item.stash_scene_id" @expand="emit('clickMedia', item, media)" />
          </div>
        </template>
        <MediaFiller v-else :media="item" class="w-full h-full" @click.stop @retry="(id) => emit('retry', id)" />
      </div>
    </div>
  </div>
</template>
