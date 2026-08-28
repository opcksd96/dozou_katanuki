<!-- frontend/src/components/admin/database/MediaPaginationBar.vue (100行以下 - SPEC-PRINCIPLE-001) -->
<script setup lang="ts">
import { computed, ref, watch } from 'vue';
import { ChevronLeft, ChevronRight, ChevronsLeft, ChevronsRight } from 'lucide-vue-next';

const props = defineProps<{ page: number; limit: number; total: number; }>();
const emit = defineEmits<{ (e: 'update:page', p: number): void; (e: 'update:limit', l: number): void; }>();

const totalPages = computed(() => Math.max(1, Math.ceil(props.total / props.limit)));
const startRange = computed(() => props.total === 0 ? 0 : (props.page - 1) * props.limit + 1);
const endRange = computed(() => Math.min(props.page * props.limit, props.total));

const jumpPage = ref(props.page);
watch(() => props.page, (p) => { jumpPage.value = p; });

const handleJump = () => {
  const target = Math.max(1, Math.min(totalPages.value, Number(jumpPage.value) || 1));
  jumpPage.value = target;
  emit('update:page', target);
};
</script>

<template>
  <div class="p-2.5 sm:p-3 bg-slate-900/90 border border-slate-800 rounded-xl flex flex-wrap items-center justify-start gap-2 sm:gap-3 text-xs sm:text-sm text-slate-300 font-sans shadow-md">
    <!-- 件数カウンター -->
    <div class="font-mono text-xs text-slate-400 bg-slate-950 px-2.5 py-1.5 rounded-lg border border-slate-800">
      <span class="text-blue-400 font-bold">{{ startRange }}-{{ endRange }}</span> / {{ total }} 件
    </div>

    <!-- ページ送りナビゲーション -->
    <div class="flex items-center gap-1 bg-slate-950 p-0.5 rounded-lg border border-slate-800">
      <button :disabled="page <= 1" @click="emit('update:page', 1)" class="p-1.5 rounded hover:bg-slate-800 disabled:opacity-30 text-slate-300 cursor-pointer disabled:cursor-not-allowed" title="先頭ページ">
        <ChevronsLeft class="w-4 h-4" />
      </button>
      <button :disabled="page <= 1" @click="emit('update:page', page - 1)" class="px-2.5 py-1.5 rounded hover:bg-slate-800 disabled:opacity-30 text-slate-200 font-bold flex items-center gap-1 cursor-pointer disabled:cursor-not-allowed" title="前のページ">
        <ChevronLeft class="w-4 h-4" /><span>前へ</span>
      </button>
      <div class="px-2 font-mono font-bold text-slate-100 text-xs sm:text-sm">
        <span class="text-blue-400">{{ page }}</span> / {{ totalPages }}
      </div>
      <button :disabled="page >= totalPages" @click="emit('update:page', page + 1)" class="px-2.5 py-1.5 rounded hover:bg-slate-800 disabled:opacity-30 text-slate-200 font-bold flex items-center gap-1 cursor-pointer disabled:cursor-not-allowed" title="次のページ">
        <span>次へ</span><ChevronRight class="w-4 h-4" />
      </button>
      <button :disabled="page >= totalPages" @click="emit('update:page', totalPages)" class="p-1.5 rounded hover:bg-slate-800 disabled:opacity-30 text-slate-300 cursor-pointer disabled:cursor-not-allowed" title="最終ページ">
        <ChevronsRight class="w-4 h-4" />
      </button>
    </div>

    <!-- ダイレクトジャンプフォーム -->
    <form @submit.prevent="handleJump" class="flex items-center gap-1.5 bg-slate-950 px-2 py-1 rounded-lg border border-slate-800">
      <span class="text-xs text-slate-400">移動:</span>
      <input v-model.number="jumpPage" type="number" min="1" :max="totalPages" class="w-14 bg-slate-900 border border-slate-700 rounded px-1.5 py-0.5 text-center font-mono text-slate-100 text-xs sm:text-sm font-bold focus:border-blue-500 focus:outline-none" />
      <span class="text-xs text-slate-500 font-mono">頁</span>
      <button type="submit" class="px-2 py-0.5 bg-blue-600 hover:bg-blue-500 text-white rounded text-xs font-bold cursor-pointer active:scale-95">Go</button>
    </form>

    <!-- 表示件数セレクタ -->
    <div class="flex items-center gap-1.5 bg-slate-950 px-2.5 py-1 rounded-lg border border-slate-800 text-xs text-slate-400">
      <span>表示:</span>
      <select :value="limit" @change="emit('update:limit', Number(($event.target as HTMLSelectElement).value))" class="bg-slate-900 border border-slate-700 rounded px-2 py-0.5 text-slate-200 text-xs font-mono font-bold cursor-pointer focus:border-blue-500 focus:outline-none">
        <option :value="12">12 件</option>
        <option :value="24">24 件</option>
        <option :value="48">48 件</option>
        <option :value="96">96 件</option>
      </select>
    </div>
  </div>
</template>
