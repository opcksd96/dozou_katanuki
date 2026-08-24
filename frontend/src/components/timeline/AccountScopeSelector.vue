<script setup lang="ts">
import { computed } from 'vue';
import type { RenderAuthor } from '../../models/RenderTree';

const props = defineProps<{
  accounts: RenderAuthor[];
  selectedId: string;
}>();

const emit = defineEmits<{
  (e: 'select', id: string): void;
}>();

interface AccountGroup {
  name: string;
  accounts: RenderAuthor[];
}

const groupedAccounts = computed(() => {
  const groups = new Map<string, RenderAuthor[]>();
  const ungrouped: RenderAuthor[] = [];
  for (const acc of props.accounts) {
    if (acc.group_name) {
      if (!groups.has(acc.group_name)) groups.set(acc.group_name, []);
      groups.get(acc.group_name)!.push(acc);
    } else {
      ungrouped.push(acc);
    }
  }
  const result: AccountGroup[] = [];
  for (const [name, accs] of groups) result.push({ name, accounts: accs });
  result.sort((a, b) => a.name.localeCompare(b.name));
  if (ungrouped.length > 0) result.push({ name: '', accounts: ungrouped });
  return result;
});

const hasGroups = computed(() => groupedAccounts.value.some(g => g.name !== ''));

const pillClass = (id: string) => [
  'flex items-center gap-2 px-3 py-1.5 rounded-full text-xs font-medium whitespace-nowrap transition-all border cursor-pointer shrink-0 shadow-sm',
  props.selectedId === id
    ? 'bg-blue-600 border-blue-500 text-white shadow-blue-500/20 font-bold'
    : 'bg-slate-950 border-slate-800 text-slate-400 hover:text-slate-200 hover:border-slate-700'
];
</script>

<template>
  <div class="bg-slate-900/40 border border-slate-800/80 rounded-xl p-3 mb-4 space-y-2">
    <div class="flex items-center justify-between text-[11px] font-mono text-slate-400">
      <span class="flex items-center gap-1.5 font-semibold text-slate-300">
        <span>📁</span> アーカイブ対象アカウント (Scope):
      </span>
      <span class="text-slate-500">{{ accounts.length }} Accounts</span>
    </div>

    <!-- 全アカウントボタン -->
    <div class="flex flex-wrap items-center gap-2 pt-1">
      <button @click="emit('select', 'all')" :class="pillClass('all')">
        <span>🌐 全てのアカウント</span>
      </button>
    </div>

    <!-- グループ化されたアカウント -->
    <template v-if="hasGroups">
      <div v-for="group in groupedAccounts" :key="group.name || '__ungrouped__'" class="space-y-1.5">
        <div v-if="group.name" class="flex items-center gap-2 pt-2">
          <span class="text-[10px] font-bold text-amber-400/80 uppercase tracking-wider flex items-center gap-1">
            🏷️ {{ group.name }}
          </span>
          <div class="flex-1 border-t border-slate-800/60"></div>
          <span class="text-[10px] text-slate-600">{{ group.accounts.length }}</span>
        </div>
        <div v-else-if="groupedAccounts.length > 1" class="flex items-center gap-2 pt-2">
          <span class="text-[10px] font-bold text-slate-500 uppercase tracking-wider">未分類</span>
          <div class="flex-1 border-t border-slate-800/60"></div>
          <span class="text-[10px] text-slate-600">{{ group.accounts.length }}</span>
        </div>
        <div class="flex flex-wrap items-center gap-2">
          <button
            v-for="acc in group.accounts"
            :key="acc.numeric_id"
            @click="emit('select', acc.numeric_id)"
            :class="pillClass(acc.numeric_id)"
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
            <span v-if="acc.alias_of" class="text-[9px] text-teal-400 opacity-80" title="裏垢">🔗</span>
          </button>
        </div>
      </div>
    </template>

    <!-- グループがない場合はフラット表示 -->
    <div v-else class="flex flex-wrap items-center gap-2">
      <button
        v-for="acc in accounts"
        :key="acc.numeric_id"
        @click="emit('select', acc.numeric_id)"
        :class="pillClass(acc.numeric_id)"
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
        <span v-if="acc.alias_of" class="text-[9px] text-teal-400 opacity-80" title="裏垢">🔗</span>
      </button>
    </div>
  </div>
</template>

