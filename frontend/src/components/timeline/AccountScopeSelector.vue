<!-- frontend/src/components/timeline/AccountScopeSelector.vue (100行以下 - SPEC-PRINCIPLE-001) -->
<script setup lang="ts">
import { computed } from 'vue';
import type { RenderAuthor } from '../../models/RenderTree';
import AccountGroupSection from './AccountGroupSection.vue';

const props = defineProps<{ accounts: RenderAuthor[]; selectedId: string; }>();
const emit = defineEmits<{ (e: 'select', id: string): void; }>();

interface AccountGroup { name: string; accounts: RenderAuthor[]; }

const groupedAccounts = computed(() => {
  const groups = new Map<string, RenderAuthor[]>();
  const ungrouped: RenderAuthor[] = [];
  for (const acc of props.accounts) {
    if (acc.group_name) {
      if (!groups.has(acc.group_name)) groups.set(acc.group_name, []);
      groups.get(acc.group_name)!.push(acc);
    } else { ungrouped.push(acc); }
  }
  const result: AccountGroup[] = [];
  for (const [name, accs] of groups) result.push({ name, accounts: accs });
  result.sort((a, b) => a.name.localeCompare(b.name));
  if (ungrouped.length > 0) result.push({ name: '', accounts: ungrouped });
  return result;
});

const definedGroups = computed(() => groupedAccounts.value.filter(g => g.name !== ''));
const hasGroups = computed(() => definedGroups.value.length > 0);
const isGroupSelected = (g: string) => props.selectedId === `group:${g}`;

const pillClass = (active: boolean, isGroup = false) => [
  'flex items-center gap-1.5 px-3 py-1.5 rounded-full text-xs font-medium whitespace-nowrap transition-all border cursor-pointer shrink-0 shadow-sm active:scale-95 touch-manipulation',
  active
    ? (isGroup ? 'bg-amber-600 border-amber-500 text-white shadow-amber-500/20 font-bold' : 'bg-blue-600 border-blue-500 text-white shadow-blue-500/20 font-bold')
    : (isGroup ? 'bg-slate-950 border-amber-900/40 text-amber-300/80 hover:text-amber-200' : 'bg-slate-950 border-slate-800 text-slate-400 hover:text-slate-200')
];
</script>

<template>
  <div class="bg-slate-900/60 border-y sm:border sm:rounded-2xl border-slate-800/80 p-3 sm:p-4 mb-2 sm:mb-4 space-y-3 shadow-xl backdrop-blur-sm">
    <div class="flex items-center justify-between text-[11px] font-mono text-slate-400 border-b border-slate-800/60 pb-2">
      <span class="flex items-center gap-1.5 font-semibold text-slate-200"><span>📁</span> スコープ (Scope):</span>
      <div class="flex items-center gap-2 text-[10px]">
        <span v-if="hasGroups" class="text-amber-400 font-semibold">{{ definedGroups.length }} Groups</span>
        <span class="text-slate-600">•</span>
        <span class="text-slate-400">{{ accounts.length }} Accounts</span>
      </div>
    </div>

    <!-- 最上段：全アカウント ＆ グループ一括選択ピルボタン一覧 -->
    <div class="space-y-1.5">
      <div class="flex flex-wrap items-center gap-1.5 sm:gap-2">
        <button @click="emit('select', 'all')" :class="pillClass(selectedId === 'all')">
          <span>🌐 全てのアカウント</span>
          <span class="text-[10px] opacity-75 font-mono">({{ accounts.length }})</span>
        </button>

        <button
          v-for="group in definedGroups"
          :key="`btn-grp-${group.name}`"
          @click="emit('select', `group:${group.name}`)"
          :class="pillClass(isGroupSelected(group.name), true)"
          :title="`グループ「${group.name}」所属全${group.accounts.length}件を一括表示`"
        >
          <span>🏷️ {{ group.name }}</span>
          <span class="text-[10px] opacity-75 font-mono">({{ group.accounts.length }})</span>
        </button>
      </div>
    </div>

    <!-- グループ別アカウント個別ピル一覧 -->
    <template v-if="hasGroups">
      <div class="border-t border-slate-800/60 pt-2 space-y-2.5">
        <AccountGroupSection
          v-for="group in groupedAccounts"
          :key="group.name || '__ungrouped__'"
          :group-name="group.name"
          :accounts="group.accounts"
          :selected-id="selectedId"
          :is-group-selected="isGroupSelected(group.name)"
          @select="(id) => emit('select', id)"
        />
      </div>
    </template>
  </div>
</template>

