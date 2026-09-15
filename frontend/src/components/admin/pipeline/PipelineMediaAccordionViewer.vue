<!-- frontend/src/components/admin/pipeline/PipelineMediaAccordionViewer.vue (100行以下 - SPEC-PRINCIPLE-001) -->
<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { usePipelineMediaTasks } from '../../../composables/admin/usePipelineMediaTasks';
import PipelineMediaAccordionItem from './PipelineMediaAccordionItem.vue';

const { items, loading, fetchItems } = usePipelineMediaTasks();
const expandAll = ref(false);

const toggleExpandAll = () => {
  expandAll.value = !expandAll.value;
};

onMounted(() => { fetchItems('ESCALATED', 50); });
</script>

<template>
  <div class="p-3.5 bg-slate-900/90 border border-slate-800 rounded-2xl flex flex-col space-y-3 shadow-lg w-full">
    <div class="flex items-center justify-between border-b border-slate-800 pb-2">
      <div class="flex items-center gap-2">
        <span class="text-xs font-bold text-slate-200 font-mono">作品別・候補URLコラプス探査</span>
        <span class="text-[10px] font-mono text-purple-400 bg-purple-950/80 px-2 py-0.5 rounded border border-purple-800">
          迅雷探索中 {{ items.length }} 作品
        </span>
      </div>
      <div class="flex items-center gap-1.5">
        <button
          @click="toggleExpandAll"
          class="px-2.5 py-1 rounded bg-slate-800 hover:bg-slate-700 text-slate-300 text-xs font-mono cursor-pointer transition shadow-sm"
        >
          {{ expandAll ? 'すべて閉じる' : 'すべて開く' }}
        </button>
        <button
          @click="fetchItems('ESCALATED', 50)"
          :disabled="loading"
          class="px-2.5 py-1 rounded bg-slate-800 hover:bg-slate-700 text-slate-300 text-xs font-mono cursor-pointer shadow-sm"
        >
          再取得
        </button>
      </div>
    </div>

    <div class="space-y-2.5 w-full pt-1">
      <div v-if="items.length === 0" class="text-center py-12 text-slate-600 text-xs font-mono">
        探索中の作品はありません
      </div>
      <PipelineMediaAccordionItem
        v-for="(m, idx) in items"
        :key="m.media_id"
        :media="m"
        :is-initially-open="expandAll || idx === 0"
      />
    </div>
  </div>
</template>
