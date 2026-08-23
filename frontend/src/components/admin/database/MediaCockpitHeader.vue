<!-- frontend/src/components/admin/database/MediaCockpitHeader.vue (100行以下 - SPEC-PRINCIPLE-001) -->
<script setup lang="ts">
import { computed } from 'vue';

const props = defineProps<{
  activeJob: any;
}>();

const emit = defineEmits<{
  (e: 'cancelJob', id: string): void;
}>();

const isJobActive = computed(() => {
  return props.activeJob && props.activeJob.status === 'RUNNING';
});
</script>

<template>
  <!-- 実行中ジョブがある時だけ薄型インジケーターとしてスッと展開 -->
  <div v-if="activeJob && isJobActive" class="bg-slate-900/95 border border-blue-600/60 rounded-xl px-3 py-1.5 flex items-center justify-between gap-3 text-xs font-mono shadow-md animate-pulse">
    <div class="flex items-center gap-2 min-w-0">
      <span class="text-blue-400 font-bold">● {{ activeJob.type }}</span>
      <span class="text-slate-300 truncate max-w-[400px]">{{ activeJob.message }}</span>
    </div>
    <div class="flex items-center gap-3 shrink-0">
      <div class="w-28 bg-slate-800 rounded-full h-2 overflow-hidden">
        <div class="bg-blue-500 h-full transition-all duration-300" :style="{ width: `${activeJob.percentage}%` }" />
      </div>
      <span class="text-blue-400 font-bold">{{ activeJob.current }}/{{ activeJob.total }} ({{ activeJob.percentage }}%)</span>
      <button @click="emit('cancelJob', activeJob.id)" class="px-2 py-0.5 bg-rose-900 hover:bg-rose-800 text-rose-200 rounded text-xs">中断</button>
    </div>
  </div>
</template>
