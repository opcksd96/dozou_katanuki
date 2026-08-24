<!-- frontend/src/components/timeline/AccountGroupSection.vue (100行以下) -->
<script setup lang="ts">
import type { RenderAuthor } from '../../models/RenderTree';

defineProps<{
  groupName: string;
  accounts: RenderAuthor[];
  selectedId: string;
  isGroupSelected: boolean;
}>();

const emit = defineEmits<{
  (e: 'select', id: string): void;
}>();

const pillClass = (active: boolean) => [
  'flex items-center gap-2 px-3 py-1.5 rounded-full text-xs font-medium whitespace-nowrap transition-all border cursor-pointer shrink-0 shadow-sm',
  active
    ? 'bg-blue-600 border-blue-500 text-white shadow-blue-500/20 font-bold'
    : 'bg-slate-950 border-slate-800 text-slate-400 hover:text-slate-200 hover:border-slate-700'
];
</script>

<template>
  <div class="space-y-1.5">
    <div class="flex items-center justify-between">
      <div class="flex items-center gap-2">
        <button
          v-if="groupName"
          @click="emit('select', `group:${groupName}`)"
          class="text-[11px] font-bold uppercase tracking-wider flex items-center gap-1.5 px-2 py-0.5 rounded transition-all cursor-pointer"
          :class="isGroupSelected
            ? 'bg-amber-500/20 text-amber-300 border border-amber-500/40'
            : 'text-amber-400/90 hover:text-amber-200 hover:bg-amber-950/40'"
          :title="`グループ「${groupName}」全体を選択`"
        >
          <span>🏷️ {{ groupName }}</span>
          <span class="text-[9px] text-amber-500/80 font-normal">一括</span>
        </button>
        <span v-else class="text-[11px] font-bold text-slate-400 uppercase tracking-wider">未分類アカウント</span>
      </div>
      <span class="text-[10px] text-slate-400 font-mono">{{ accounts.length }} 件</span>
    </div>

    <div class="flex flex-wrap items-center gap-2 pl-1">
      <button
        v-for="acc in accounts"
        :key="acc.numeric_id"
        @click="emit('select', acc.numeric_id)"
        :class="pillClass(selectedId === acc.numeric_id)"
        :title="acc.alias_of ? `裏垢 → 本垢: @${acc.alias_of}` : acc.display_name"
      >
        <img
          v-if="acc.avatar_url"
          :src="acc.avatar_url"
          :alt="acc.handle"
          class="w-4 h-4 rounded-full object-cover bg-slate-800"
          @error="($event.target as HTMLElement).style.display = 'none'"
        />
        <span>@{{ acc.handle }}</span>
        <span v-if="acc.alias_of" class="text-[9px] text-teal-400 font-mono opacity-80" title="裏垢">🔗</span>
      </button>
    </div>
  </div>
</template>
