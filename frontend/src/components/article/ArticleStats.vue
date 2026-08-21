<!-- frontend/src/components/article/ArticleStats.vue (100行以下) -->
<script setup lang="ts">
import { computed } from 'vue';
import type { RenderMetrics } from '../../models/RenderTree';
import { formatStatNumber } from '../../utils/formatters';

const props = withDefaults(
  defineProps<{
    metrics?: RenderMetrics;
    stats?: RenderMetrics;
    isLiked?: boolean;
  }>(),
  { isLiked: false }
);

const effectiveMetrics = computed<RenderMetrics>(() => {
  return props.metrics || props.stats || { replies: 0, retweets: 0, likes: 0 };
});

const emit = defineEmits<{
  (e: 'toggleLike'): void;
}>();
</script>

<template>
  <div class="flex items-center justify-between max-w-sm mt-2 text-xs text-slate-400 font-mono select-none">
    <!-- リプライ -->
    <div class="flex items-center gap-1.5 hover:text-sky-400 transition-colors cursor-pointer group" title="返信">
      <span class="group-hover:scale-110 transition-transform">💬</span>
      <span>{{ formatStatNumber(effectiveMetrics.replies || 0) }}</span>
    </div>

    <!-- リツイート -->
    <div class="flex items-center gap-1.5 hover:text-emerald-400 transition-colors cursor-pointer group" title="リポスト">
      <span class="group-hover:scale-110 transition-transform">🔁</span>
      <span>{{ formatStatNumber(effectiveMetrics.retweets || 0) }}</span>
    </div>

    <!-- いいね -->
    <button
      @click="emit('toggleLike')"
      class="flex items-center gap-1.5 transition-colors cursor-pointer group"
      :class="isLiked ? 'text-rose-500 font-bold' : 'hover:text-rose-400'"
      title="いいね"
    >
      <span class="group-hover:scale-125 transition-transform">{{ isLiked ? '❤️' : '🤍' }}</span>
      <span>{{ formatStatNumber(effectiveMetrics.likes || 0) }}</span>
    </button>

    <!-- 表示回数 -->
    <span v-if="effectiveMetrics.views" class="text-slate-500 text-[11px]" title="インプレッション">
      📊 {{ formatStatNumber(effectiveMetrics.views) }}
    </span>
  </div>
</template>
