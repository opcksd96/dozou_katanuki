<!-- frontend/src/components/article/ArticleStats.vue (100行以下 - SPEC-PRINCIPLE-001) -->
<script setup lang="ts">
import { computed } from 'vue';
import type { RenderMetrics } from '../../models/RenderTree';
import { formatStatNumber } from '../../utils/formatters';
import { useToast } from '../../composables/useToast';

const props = withDefaults(
  defineProps<{
    metrics?: RenderMetrics;
    stats?: RenderMetrics;
    isLiked?: boolean;
  }>(),
  { isLiked: false }
);

const emit = defineEmits<{
  (e: 'toggleLike'): void;
}>();

const { addToast } = useToast();
const effectiveMetrics = computed<RenderMetrics>(() => props.metrics || props.stats || { replies: 0, retweets: 0, likes: 0 });

const handleShare = () => {
  addToast('📤 共有メニューを展開しました', 'info', 2000);
};
</script>

<template>
  <div class="w-full flex items-center justify-between py-2 text-slate-400 font-sans text-xs select-none">
    <!-- 返信 (Reply) -->
    <div class="flex items-center gap-1.5 hover:text-sky-400 transition-colors cursor-pointer group" title="返信">
      <div class="w-8 h-8 rounded-full flex items-center justify-center group-hover:bg-sky-500/10 transition-colors">
        <span class="text-sm group-hover:scale-110 transition-transform">💬</span>
      </div>
      <span class="font-mono text-xs">{{ formatStatNumber(effectiveMetrics.replies || 0) }}</span>
    </div>

    <!-- リポスト (Repost) -->
    <div class="flex items-center gap-1.5 hover:text-emerald-400 transition-colors cursor-pointer group" title="リポスト">
      <div class="w-8 h-8 rounded-full flex items-center justify-center group-hover:bg-emerald-500/10 transition-colors">
        <span class="text-sm group-hover:scale-110 transition-transform">🔁</span>
      </div>
      <span class="font-mono text-xs">{{ formatStatNumber(effectiveMetrics.retweets || 0) }}</span>
    </div>

    <!-- いいね (Like) -->
    <button
      @click="emit('toggleLike')"
      class="flex items-center gap-1.5 transition-colors cursor-pointer group"
      :class="isLiked ? 'text-rose-500 font-bold' : 'hover:text-rose-400'"
      title="いいね"
    >
      <div class="w-8 h-8 rounded-full flex items-center justify-center group-hover:bg-rose-500/10 transition-colors">
        <span class="text-sm group-hover:scale-125 transition-transform">{{ isLiked ? '❤️' : '🤍' }}</span>
      </div>
      <span class="font-mono text-xs">{{ formatStatNumber(effectiveMetrics.likes || 0) }}</span>
    </button>

    <!-- ブックマーク (Bookmark) -->
    <div class="flex items-center gap-1.5 hover:text-blue-400 transition-colors cursor-pointer group" title="ブックマーク">
      <div class="w-8 h-8 rounded-full flex items-center justify-center group-hover:bg-blue-500/10 transition-colors">
        <span class="text-sm group-hover:scale-110 transition-transform">🔖</span>
      </div>
      <span class="font-mono text-xs">{{ formatStatNumber((effectiveMetrics as any).bookmarks || 0) }}</span>
    </div>

    <!-- 共有 (Share) -->
    <div @click="handleShare" class="flex items-center hover:text-sky-400 transition-colors cursor-pointer group" title="共有">
      <div class="w-8 h-8 rounded-full flex items-center justify-center group-hover:bg-sky-500/10 transition-colors">
        <span class="text-sm group-hover:scale-110 transition-transform">📤</span>
      </div>
    </div>
  </div>
</template>
