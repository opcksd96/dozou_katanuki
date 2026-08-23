<!-- frontend/src/components/admin/database/DatabaseSpreadsheet.vue (100行以下) -->
<script setup lang="ts">
import { ref } from 'vue';

const props = defineProps<{
  columns: { key: string; label: string; width?: string }[];
  rows: any[];
  selectedRowId?: string;
  idKey?: string;
}>();

const emit = defineEmits<{ (e: 'selectRow', row: any): void }>();
const copiedCell = ref<string | null>(null);

const copyVal = (val: any, cellKey: string) => {
  if (val === null || val === undefined) return;
  const str = typeof val === 'object' ? JSON.stringify(val) : String(val);
  navigator.clipboard.writeText(str);
  copiedCell.value = cellKey;
  setTimeout(() => { if (copiedCell.value === cellKey) copiedCell.value = null; }, 1200);
};
</script>

<template>
  <div class="border border-slate-700/80 rounded-xl overflow-hidden bg-slate-950 flex flex-col h-full min-h-0 font-mono text-xs select-text shadow-inner">
    <div class="overflow-x-auto overflow-y-auto flex-1 min-h-0">
      <table class="w-full border-collapse text-left">
        <thead class="bg-slate-900 sticky top-0 z-10 border-b border-slate-700 text-slate-400 select-none shadow-sm">
          <tr>
            <th class="w-10 px-2 py-2 text-center border-r border-slate-800 text-[10px] bg-slate-900/90">#</th>
            <th v-for="col in columns" :key="col.key" class="px-3 py-2 border-r border-slate-800 text-[11px] font-bold tracking-wider" :style="{ minWidth: col.width || '120px' }">
              {{ col.label }}
            </th>
          </tr>
        </thead>
        <tbody class="divide-y divide-slate-800/80 text-slate-300">
          <tr v-if="!rows || rows.length === 0">
            <td :colspan="columns.length + 1" class="text-center py-10 text-slate-500 font-sans">データがありません</td>
          </tr>
          <tr v-for="(r, idx) in (rows || []).filter(Boolean)" :key="r?.[idKey || 'id'] || idx" @click="emit('selectRow', r)" class="hover:bg-blue-900/20 cursor-pointer transition-colors" :class="{ 'bg-blue-950/40 text-white font-medium': selectedRowId && (r?.[idKey || 'id'] === selectedRowId) }">
            <td class="px-2 py-1.5 text-center border-r border-slate-800 text-[10px] text-slate-500 bg-slate-900/40 select-none">{{ idx + 1 }}</td>
            <td v-for="col in columns" :key="col.key" @dblclick="copyVal(r?.[col.key], `${idx}-${col.key}`)" class="px-3 py-1.5 border-r border-slate-800/60 truncate max-w-xs relative group" :title="String(r?.[col.key] ?? '')">
              <span v-if="copiedCell === `${idx}-${col.key}`" class="text-emerald-400 font-sans text-[10px]">コピー完了!</span>
              <span v-else>{{ r?.[col.key] !== null && r?.[col.key] !== undefined ? (typeof r[col.key] === 'boolean' ? (r[col.key] ? 'TRUE' : 'FALSE') : r[col.key]) : '-' }}</span>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

