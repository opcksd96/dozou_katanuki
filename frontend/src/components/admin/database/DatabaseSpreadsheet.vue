<!-- frontend/src/components/admin/database/DatabaseSpreadsheet.vue (100行以下 - SPEC-PRINCIPLE-001) -->
<script setup lang="ts">
import { ref } from 'vue';

const props = defineProps<{
  columns: { key: string; label: string; width?: string }[];
  rows: any[]; selectedRowId?: string; idKey?: string; selectedIds?: Set<string>;
}>();

const emit = defineEmits<{
  (e: 'selectRow', row: any): void; (e: 'toggleSelect', id: string): void;
  (e: 'toggleSelectAll'): void; (e: 'scrollBottom'): void;
}>();

const copiedCell = ref<string | null>(null);

const copyVal = (val: any, cellKey: string) => {
  if (val === null || val === undefined) return;
  navigator.clipboard.writeText(typeof val === 'object' ? JSON.stringify(val) : String(val));
  copiedCell.value = cellKey;
  setTimeout(() => { if (copiedCell.value === cellKey) copiedCell.value = null; }, 1200);
};

const handleKeydown = (e: KeyboardEvent) => {
  if (e.key !== 'ArrowDown' && e.key !== 'ArrowUp') return;
  const valid = (props.rows || []).filter(Boolean);
  if (valid.length === 0) return;
  e.preventDefault();
  const cur = valid.findIndex(r => (r?.[props.idKey || 'id'] || '') === props.selectedRowId);
  const next = Math.max(0, Math.min(valid.length - 1, cur + (e.key === 'ArrowDown' ? 1 : -1)));
  emit('selectRow', valid[next]);
  document.getElementById(`row-${valid[next]?.[props.idKey || 'id']}`)?.scrollIntoView({ block: 'nearest' });
};

const handleScroll = (e: Event) => {
  const t = e.target as HTMLElement;
  if (t.scrollTop + t.clientHeight >= t.scrollHeight - 50) emit('scrollBottom');
};
</script>

<template>
  <div @keydown="handleKeydown" tabindex="0" class="border border-slate-700/80 rounded-xl overflow-hidden bg-slate-950 flex flex-col h-full min-h-0 font-mono text-xs select-text shadow-inner focus:outline-none focus:ring-1 focus:ring-blue-500/50">
    <div @scroll="handleScroll" class="overflow-x-auto overflow-y-auto flex-1 min-h-0">
      <table class="w-full border-collapse text-left">
        <thead class="bg-slate-900 sticky top-0 z-10 border-b border-slate-700 text-slate-400 select-none shadow-sm">
          <tr>
            <th class="w-8 px-1.5 py-2 text-center border-r border-slate-800 bg-slate-900/90">
              <input type="checkbox" :checked="(selectedIds?.size || 0) > 0 && selectedIds?.size === rows.length" @change="emit('toggleSelectAll')" class="rounded border-slate-700 bg-slate-950 text-blue-600 focus:ring-0 cursor-pointer w-3.5 h-3.5" />
            </th>
            <th class="w-10 px-2 py-2 text-center border-r border-slate-800 text-[10px] bg-slate-900/90">#</th>
            <th v-for="col in columns" :key="col.key" class="px-3 py-2 border-r border-slate-800 text-[11px] font-bold tracking-wider" :style="{ minWidth: col.width || '120px' }">{{ col.label }}</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-slate-800/80 text-slate-300">
          <tr v-if="!rows || rows.length === 0"><td :colspan="columns.length + 2" class="text-center py-10 text-slate-500 font-sans">データがありません</td></tr>
          <tr v-for="(r, idx) in (rows || []).filter(Boolean)" :key="r?.[idKey || 'id'] || idx" :id="`row-${r?.[idKey || 'id']}`" @click="emit('selectRow', r)" class="hover:bg-blue-900/20 cursor-pointer transition-colors" :class="{ 'bg-blue-950/40 text-white font-medium': selectedRowId && (r?.[idKey || 'id'] === selectedRowId), 'opacity-60 bg-rose-950/10 line-through decoration-rose-500/50': r?._raw?.is_trash || r?.is_trash }">
            <td class="px-1.5 py-1.5 text-center border-r border-slate-800/60 bg-slate-900/20" @click.stop>
              <input type="checkbox" :checked="selectedIds?.has(r?.[idKey || 'id'])" @change="emit('toggleSelect', r?.[idKey || 'id'])" class="rounded border-slate-700 bg-slate-950 text-blue-600 focus:ring-0 cursor-pointer w-3.5 h-3.5" />
            </td>
            <td class="px-2 py-1.5 text-center border-r border-slate-800 text-[10px] text-slate-500 bg-slate-900/40 select-none">{{ idx + 1 }}</td>
            <td v-for="col in columns" :key="col.key" @dblclick="copyVal(r?.[col.key], `${idx}-${col.key}`)" class="px-3 py-1.5 border-r border-slate-800/60 truncate max-w-xs relative group" :title="String(r?.[col.key] ?? '')">
              <slot :name="`cell-${col.key}`" :row="r" :value="r?.[col.key]">
                <span v-if="copiedCell === `${idx}-${col.key}`" class="text-emerald-400 font-sans text-[10px]">コピー完了!</span>
                <span v-else-if="typeof r?.[col.key] === 'string' && r[col.key].startsWith('http')">
                  <a :href="r[col.key]" target="_blank" rel="noopener noreferrer" class="text-sky-400 hover:underline truncate inline-block max-w-full" @click.stop>{{ r[col.key] }}</a>
                </span>
                <span v-else>{{ r?.[col.key] !== null && r?.[col.key] !== undefined ? (typeof r[col.key] === 'boolean' ? (r[col.key] ? 'TRUE' : 'FALSE') : r[col.key]) : '-' }}</span>
              </slot>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>
