<!-- frontend/src/components/article/Avatar.vue (100行以下) -->
<script setup lang="ts">
import { computed } from 'vue';
import type { RenderAuthor } from '../../models/RenderTree';

const props = defineProps<{
  avatarUrl?: string;
  handle?: string;
  author?: RenderAuthor;
}>();

const resolvedUrl = computed(() => props.avatarUrl || props.author?.avatar_url || '');
const resolvedHandle = computed(() => props.handle || props.author?.handle || '');
</script>

<template>
  <div class="relative flex-shrink-0">
    <img
      :src="resolvedUrl"
      :alt="resolvedHandle"
      class="w-10 h-10 rounded-full border border-slate-700 object-cover bg-slate-800"
      loading="lazy"
    />
  </div>
</template>
