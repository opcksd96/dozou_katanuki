<!-- frontend/src/components/admin/whitelist/WhitelistAddForm.vue (100行以下) -->
<script setup lang="ts">
import { reactive } from 'vue';
const emit = defineEmits<{ (e: 'add', type: string, value: string): void }>();
const newForm = reactive({ type: 'account', value: '' });

const handleAdd = () => {
  if (!newForm.value.trim()) return;
  emit('add', newForm.type, newForm.value.trim());
  newForm.value = '';
};
</script>

<template>
  <div class="p-4 bg-slate-900/60 border border-slate-800 rounded-xl space-y-3">
    <h4 class="text-xs font-bold text-slate-300 flex items-center gap-1.5"><span>➕</span> 新規 Whitelist ルールの追加</h4>
    <form @submit.prevent="handleAdd" class="flex flex-wrap items-center gap-3">
      <div class="w-36">
        <select v-model="newForm.type" class="w-full bg-slate-950 border border-slate-700 rounded-lg px-3 py-2 text-xs text-slate-200">
          <option value="account">👤 アカウント</option>
          <option value="keyword">🔍 キーワード</option>
        </select>
      </div>
      <div class="flex-1 min-w-[200px]">
        <input
          v-model="newForm.value"
          type="text"
          :placeholder="newForm.type === 'account' ? '例: msluo14 (@なし)' : '例: famicom, apu'"
          class="w-full bg-slate-950 border border-slate-700 rounded-lg px-3 py-2 text-xs text-slate-200 placeholder-slate-500"
        />
      </div>
      <button
        type="submit"
        :disabled="!newForm.value.trim()"
        class="px-4 py-2 bg-blue-600 hover:bg-blue-500 text-white text-xs font-bold rounded-lg disabled:opacity-50"
      >
        追加する
      </button>
    </form>
  </div>
</template>
