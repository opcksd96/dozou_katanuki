<script setup lang="ts">
import type { RenderMedia } from '../../models/RenderTree';
import MediaFiller from './MediaFiller.vue';
import StashPlayer from './StashPlayer.vue';

defineProps<{
  media: RenderMedia[];
}>();

const emit = defineEmits<{
  (e: 'retry', mediaId: string): void;
  (e: 'clickMedia', media: RenderMedia, list: RenderMedia[]): void;
}>();
</script>

<template>
  <div v-if="media && media.length > 0" class="mt-3 rounded-lg overflow-hidden border border-slate-800 bg-black">
    <div :class="['grid gap-1', media.length === 1 ? 'grid-cols-1' : 'grid-cols-2']">
      <div v-for="item in media" :key="item.id" class="relative group">
        <!-- ローカル確保完了時: 描画 -->
        <template v-if="item.download_status === 'COMPLETED'">
          <img
            v-if="item.type === 'image' || item.type === 'gif'"
            :src="item.urls.image || item.urls.thumbnail || item.urls.original"
            :alt="item.id"
            class="w-full h-auto max-h-[400px] object-cover hover:opacity-95 transition-opacity cursor-pointer"
            loading="lazy"
            @click="emit('clickMedia', item, media)"
          />
          <div v-else-if="item.type === 'video'" class="w-full max-h-[400px]">
            <StashPlayer
              :src="item.urls.stream || item.urls.original"
              :poster="item.urls.thumbnail"
              :controls="true"
            />
          </div>
        </template>
        <!-- 未完了/失敗時: SVGプレースホルダー・フィラー -->
        <MediaFiller v-else :media="item" @retry="(id) => emit('retry', id)" />
      </div>
    </div>
  </div>
</template>
