<script setup lang="ts">
import type { MediaItem } from '../types/article';

defineProps<{
  media: MediaItem[];
}>();
</script>

<template>
  <div v-if="media && media.length > 0" class="mt-3 rounded-xl overflow-hidden border border-slate-800 bg-slate-950/80">
    <div
      :class="[
        'grid gap-1',
        media.length === 1 ? 'grid-cols-1' : 'grid-cols-2'
      ]"
    >
      <div
        v-for="item in media"
        :key="item.mediaId"
        class="relative group flex items-center justify-center bg-black/40 overflow-hidden"
        :class="media.length > 1 ? 'h-52 md:h-64' : 'max-h-[580px]'"
      >
        <template v-if="item.type === 'image'">
          <img
            :src="item.url"
            :alt="item.mediaId"
            :class="[
              'transition-opacity hover:opacity-95',
              media.length === 1
                ? 'w-full h-auto max-h-[580px] object-contain mx-auto'
                : 'w-full h-full object-cover'
            ]"
            loading="lazy"
          />
        </template>
        <template v-else-if="item.type === 'video'">
          <video
            :src="item.url"
            controls
            preload="metadata"
            class="max-w-full max-h-[580px] object-contain rounded"
          />
        </template>
        <div class="absolute bottom-1 right-1 bg-black/60 px-1.5 py-0.5 rounded text-[10px] text-slate-400 font-mono">
          {{ item.mediaId }}
        </div>
      </div>
    </div>
  </div>
</template>

