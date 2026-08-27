<!-- frontend/src/components/admin/database/MediaPaginationBar.vue (100行以下) -->
<script setup lang="ts">
import { computed } from 'vue';

const props = defineProps<{
  page: number;
  limit: number;
  total: number;
}>();

const emit = defineEmits<{
  (e: 'update:page', p: number): void;
  (e: 'update:limit', l: number): void;
}>();

const totalPages = computed(() => Math.max(1, Math.ceil(props.total / props.limit)));
const startRange = computed(() => props.total === 0 ? 0 : (props.page - 1) * props.limit + 1);
const endRange = computed(() => Math.min(props.page * props.limit, props.total));
</script>

<template>
  <div class="p-2 bg-slate-900 border border-slate-800 rounded-xl flex items-center justify-between text-xs text-slate-400">
    <div>{{ startRange }} - {{ endRange }} / {{ total }} 件</div>
    <div class="flex items-center gap-2">
      <button :disabled="page <= 1" @click="emit('update:page', page - 1)" class="px-2.5 py-1 bg-slate-950 border border-slate-700 rounded disabled:opacity-40 hover:bg-slate-800 text-slate-200">
        ◀ 前へ
      </button>
      <span class="text-slate-200 font-mono">{{ page }} / {{ totalPages }}</span>
      <button :disabled="page >= totalPages" @click="emit('update:page', page + 1)" class="px-2.5 py-1 bg-slate-950 border border-slate-700 rounded disabled:opacity-40 hover:bg-slate-800 text-slate-200">
        次へ ▶
      </button>
      <div class="flex items-center gap-1 ml-2">
        <span class="text-[11px] text-slate-400">表示:</span>
        <select :value="limit" @change="emit('update:limit', Number(($event.target as HTMLSelectElement).value))" class="bg-slate-950 border border-slate-700 rounded px-1.5 py-1 text-slate-200 text-xs cursor-pointer">
          <option :value="12">12件</option>
          <option :value="24">24件</option>
          <option :value="48">48件</option>
          <option :value="96">96件</option>
        </select>
      </div>
    </div>
  </div>
</template>
