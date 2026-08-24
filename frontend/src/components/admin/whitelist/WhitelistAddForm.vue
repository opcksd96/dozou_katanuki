<!-- frontend/src/components/admin/whitelist/WhitelistAddForm.vue (100行以下) -->
<script setup lang="ts">
import { reactive } from 'vue';
const emit = defineEmits<{ (e: 'add', type: string, value: string, groupName: string, aliasOf: string): void }>();
const newForm = reactive({ type: 'account', value: '', groupName: '', aliasOf: '' });

const handleAdd = () => {
  if (!newForm.value.trim()) return;
  emit('add', newForm.type, newForm.value.trim(), newForm.groupName.trim(), newForm.aliasOf.trim());
  newForm.value = ''; newForm.groupName = ''; newForm.aliasOf = '';
};
</script>

<template>
  <div class="p-4 bg-slate-900/60 border border-slate-800 rounded-xl space-y-3">
    <h4 class="text-xs font-bold text-slate-300 flex items-center gap-1.5"><span>➕</span> 新規 Whitelist ルール・裏垢エイリアス追加</h4>
    <form @submit.prevent="handleAdd" class="space-y-2.5">
      <div class="flex flex-wrap items-center gap-3">
        <div class="w-36">
          <select v-model="newForm.type" class="w-full bg-slate-950 border border-slate-700 rounded-lg px-3 py-2 text-xs text-slate-200">
            <option value="account">👤 アカウント</option>
            <option value="keyword">🔍 キーワード</option>
          </select>
        </div>
        <div class="flex-1 min-w-[180px]">
          <input
            v-model="newForm.value"
            type="text"
            :placeholder="newForm.type === 'account' ? 'アカウント名 (例: sub_mashu)' : 'キーワード (例: famicom)'"
            class="w-full bg-slate-950 border border-slate-700 rounded-lg px-3 py-2 text-xs text-slate-200 placeholder-slate-500"
          />
        </div>
      </div>

      <div v-if="newForm.type === 'account'" class="flex flex-wrap items-center gap-3 pt-1 border-t border-slate-800/60">
        <div class="flex-1 min-w-[150px]">
          <input
            v-model="newForm.groupName"
            type="text"
            placeholder="所属グループ名 (任意: 例: カルデア)"
            class="w-full bg-slate-950 border border-slate-700 rounded-lg px-3 py-1.5 text-xs text-slate-300 placeholder-slate-600"
          />
        </div>
        <div class="flex-1 min-w-[150px]">
          <input
            v-model="newForm.aliasOf"
            type="text"
            placeholder="裏垢の場合の本垢名 (任意: 例: main_mashu)"
            class="w-full bg-slate-950 border border-slate-700 rounded-lg px-3 py-1.5 text-xs text-slate-300 placeholder-slate-600"
          />
        </div>
        <button
          type="submit"
          :disabled="!newForm.value.trim()"
          class="px-4 py-2 bg-blue-600 hover:bg-blue-500 text-white text-xs font-bold rounded-lg disabled:opacity-50 transition-colors ml-auto"
        >
          追加する
        </button>
      </div>
      <div v-else class="flex justify-end">
        <button
          type="submit"
          :disabled="!newForm.value.trim()"
          class="px-4 py-2 bg-blue-600 hover:bg-blue-500 text-white text-xs font-bold rounded-lg disabled:opacity-50 transition-colors"
        >
          追加する
        </button>
      </div>
    </form>
  </div>
</template>

