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
  <div v-if="media && media.length > 0" class="mt-3 rounded-xl overflow-hidden border border-slate-800 bg-slate-950/80">
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
            v-if="item.type === 'image' || item.type === 'gif'"
            :src="item.urls.image || item.urls.thumbnail"
            :alt="item.id"
            :class="[
              'cursor-pointer transition-opacity hover:opacity-95',
              media.length === 1
                ? 'w-full h-auto max-h-[580px] object-contain mx-auto'
                : 'w-full h-full object-cover'
            ]"
            loading="lazy"
            @click="emit('clickMedia', item, media)"
          />
          <div
            v-else-if="item.type === 'video'"
            :class="[
              'w-full flex items-center justify-center',
              media.length === 1 ? 'max-h-[580px]' : 'h-full'
            ]"
          >
            <StashPlayer
              :src="item.urls.stream"
              :poster="item.urls.thumbnail"
              :stashSceneId="item.stash_scene_id"
              :controls="true"
            />
          </div>
        </template>
        <!-- 未完了/実体未紐付け/失敗時: SVGプレースホルダー・フィラー -->
        <MediaFiller v-else :media="item" class="w-full h-full" @retry="(id) => emit('retry', id)" />
      </div>
    </div>
  </div>
</template>

