<!-- frontend/src/components/article/Avatar.vue (100行以下) -->
<script setup lang="ts">
import { ref, computed, watch } from 'vue';
import type { RenderAuthor } from '../../models/RenderTree';
import { resolveAvatarUrl, getAvatarInitial } from '../../utils/avatar';

const props = defineProps<{
  avatarUrl?: string;
  handle?: string;
  author?: RenderAuthor;
  sizeClass?: string;
}>();

const imgError = ref(false);
const resolvedUrl = computed(() => props.avatarUrl || resolveAvatarUrl(props.author) || '');
const resolvedHandle = computed(() => props.handle || props.author?.handle || props.author?.display_name || '');
const initial = computed(() => getAvatarInitial(resolvedHandle.value));

watch(resolvedUrl, () => { imgError.value = false; });
</script>

<template>
  <div class="relative flex-shrink-0 flex items-center justify-center overflow-hidden" :class="sizeClass || 'w-10 h-10'">
    <img
      v-if="resolvedUrl && !imgError"
      :src="resolvedUrl"
      :alt="resolvedHandle"
      @error="imgError = true"
      class="w-full h-full rounded-full border border-slate-700 object-cover bg-slate-800"
      loading="lazy"
    />
    <!-- 人型シルエットSVGフォールバック -->
    <svg
      v-else
      class="w-full h-full rounded-full border border-slate-700 bg-slate-800/80 text-slate-500 shadow-inner"
      viewBox="0 0 64 64"
      fill="none"
      xmlns="http://www.w3.org/2000/svg"
    >
      <rect width="64" height="64" fill="#1e293b" />
      <circle cx="32" cy="24" r="12" fill="#64748b" />
      <path d="M16 54c0-8.837 7.163-16 16-16s16 7.163 16 16" fill="#64748b" />
    </svg>
  </div>
</template>
