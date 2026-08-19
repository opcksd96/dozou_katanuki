<script setup lang="ts">
import type { MediaItem } from '../types/article';

defineProps<{
  media: MediaItem[];
}>();
</script>

<template>
  <div v-if="media && media.length > 0" class="mt-3 rounded-lg overflow-hidden border border-slate-800 bg-black">
    <div
      :class="[
        'grid gap-1',
        media.length === 1 ? 'grid-cols-1' : 'grid-cols-2'
      ]"
    >
      <div v-for="item in media" :key="item.mediaId" class="relative group">
        <template v-if="item.type === 'image'">
          <img
            :src="item.url"
            :alt="item.mediaId"
            class="w-full h-auto max-h-[400px] object-cover hover:opacity-90 transition-opacity"
            loading="lazy"
          />
        </template>
        <template v-else-if="item.type === 'video'">
          <video
            :src="item.url"
            controls
            preload="metadata"
            class="w-full max-h-[400px]"
          />
        </template>
        <div class="absolute bottom-1 right-1 bg-black/60 px-1.5 py-0.5 rounded text-[10px] text-slate-400 font-mono">
          {{ item.mediaId }}
        </div>
      </div>
    </div>
  </div>
</template>
