<script setup lang="ts">
import type { RenderMetrics } from '../../models/RenderTree';
import { formatStatNumber } from '../../utils/formatters';

defineProps<{
  metrics: RenderMetrics;
  isLiked: boolean;
}>();

const emit = defineEmits<{
  (e: 'toggleLike'): void;
}>();
</script>

<template>
  <div class="flex items-center gap-6 mt-3 text-xs text-slate-400 pt-2 border-t border-slate-800/60 font-mono">
    <span>💬 {{ formatStatNumber(metrics.replies) }}</span>
    <span>🔁 {{ formatStatNumber(metrics.retweets) }}</span>
    <button
      @click="emit('toggleLike')"
      class="flex items-center gap-1.5 transition-colors"
      :class="isLiked ? 'text-rose-500 font-bold' : 'hover:text-rose-400'"
    >
      <span>{{ isLiked ? '❤️' : '🤍' }}</span>
      <span>{{ formatStatNumber(metrics.likes) }}</span>
    </button>
    <span v-if="metrics.views" class="text-slate-500 ml-auto">👁️ {{ formatStatNumber(metrics.views) }}</span>
  </div>
</template>
