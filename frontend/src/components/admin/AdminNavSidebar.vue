<!-- frontend/src/components/admin/AdminNavSidebar.vue (100行以下 - SPEC-PRINCIPLE-001) -->
<script setup lang="ts">
import { AdminTabId, adminTabGroups } from '../../models/adminTabs';

defineProps<{ activeTab: AdminTabId }>();
const emit = defineEmits<{ (e: 'select', tab: AdminTabId): void }>();
</script>

<template>
  <aside class="hidden md:flex w-44 bg-slate-900/80 border-r border-slate-800 flex-col py-3 px-2 gap-4 shrink-0 select-none overflow-y-auto font-sans">
    <div v-for="g in adminTabGroups" :key="g.title" class="space-y-1">
      <div class="px-2 text-[10px] font-bold text-slate-500 uppercase tracking-wider">{{ g.title }}</div>
      <button
        v-for="item in g.items"
        :key="item.id"
        @click="emit('select', item.id)"
        :class="[
          'w-full flex items-center gap-2 px-2.5 py-1.5 rounded-lg text-xs font-semibold transition-all text-left cursor-pointer active:scale-95',
          activeTab === item.id ? 'bg-blue-600 text-white shadow-md shadow-blue-900/30' : 'text-slate-400 hover:text-slate-200 hover:bg-slate-800/60'
        ]"
      >
        <span>{{ item.icon }}</span>
        <span class="truncate">{{ item.label }}</span>
      </button>
    </div>
  </aside>
</template>
