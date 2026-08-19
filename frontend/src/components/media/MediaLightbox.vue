<script setup lang="ts">
import type { RenderMedia } from '../../models/RenderTree';
import MediaOverlay from './MediaOverlay.vue';

defineProps<{
  media: RenderMedia | null;
  hasNext?: boolean;
  hasPrev?: boolean;
}>();

const emit = defineEmits<{
  (e: 'close'): void;
  (e: 'next'): void;
  (e: 'prev'): void;
}>();
</script>

<template>
  <MediaOverlay
    :media="media"
    :hasNext="hasNext"
    :hasPrev="hasPrev"
    @close="emit('close')"
    @next="emit('next')"
    @prev="emit('prev')"
  />
</template>
