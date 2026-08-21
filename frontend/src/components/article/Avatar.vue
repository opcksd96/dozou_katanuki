<!-- frontend/src/components/article/Avatar.vue (100行以下) -->
<script setup lang="ts">
import { ref, computed, watch } from 'vue';
import type { RenderAuthor } from '../../models/RenderTree';

const props = defineProps<{
  avatarUrl?: string;
  handle?: string;
  author?: RenderAuthor;
}>();

const imgError = ref(false);
const resolvedUrl = computed(() => props.avatarUrl || props.author?.avatar_url || '');
const resolvedHandle = computed(() => props.handle || props.author?.handle || props.author?.display_name || '');
const initial = computed(() => (resolvedHandle.value.charAt(0) || 'A').toUpperCase());

watch(resolvedUrl, () => { imgError.value = false; });
</script>

<template>
  <div class="relative flex-shrink-0 flex items-center justify-center">
    <img
      v-if="resolvedUrl && !imgError"
      :src="resolvedUrl"
      :alt="resolvedHandle"
      @error="imgError = true"
      class="w-10 h-10 rounded-full border border-slate-700 object-cover bg-slate-800"
      loading="lazy"
    />
    <div
      v-else
      class="w-10 h-10 rounded-full border border-slate-700 bg-slate-800 text-slate-200 font-bold flex items-center justify-center text-sm shadow-inner"
    >
      {{ initial }}
    </div>
  </div>
</template>
