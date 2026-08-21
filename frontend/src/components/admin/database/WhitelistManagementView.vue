<!-- frontend/src/components/admin/database/WhitelistManagementView.vue (100行以下) -->
<script setup lang="ts">
import { reactive, onMounted } from 'vue';
import DatabaseSpreadsheet from './DatabaseSpreadsheet.vue';

const props = defineProps<{
  whitelistList: any[];
  loading: boolean;
}>();

const emit = defineEmits<{
  (e: 'fetch'): void;
  (e: 'add', type: string, value: string): void;
  (e: 'toggle', id: number): void;
  (e: 'delete', id: number): void;
}>();

const newForm = reactive({ type: 'account', value: '' });
onMounted(() => { emit('fetch'); });

const handleAdd = () => {
  if (!newForm.value.trim()) return;
  emit('add', newForm.type, newForm.value.trim());
  newForm.value = '';
};

const cols = [
  { key: 'id', label: 'ID', width: '60px' },
  { key: 'type', label: 'Type', width: '100px' },
  { key: 'value', label: 'Target Value (ScreenName/Handle)', width: '240px' },
  { key: 'is_active', label: 'Is Active', width: '100px' },
];
</script>

<template>
  <div class="space-y-4 flex flex-col h-full">
    <!-- 新規追加バー -->
    <div class="flex gap-2 items-center bg-slate-950 p-3 rounded-xl border border-slate-800 text-xs">
      <select v-model="newForm.type" class="bg-slate-900 border border-slate-700 rounded-lg px-3 py-1.5 text-slate-200">
        <option value="account">アカウント (@handle)</option>
        <option value="keyword">キーワード</option>
      </select>
      <input v-model="newForm.value" @keyup.enter="handleAdd" type="text" placeholder="追加するハンドル名..." class="flex-1 bg-slate-900 border border-slate-700 rounded-lg px-3 py-1.5 text-slate-200" />
      <button @click="handleAdd" class="px-4 py-1.5 bg-blue-600 hover:bg-blue-500 text-white font-bold rounded-lg">+ 追加</button>
    </div>

    <!-- Excel風スプレッドシート ＆ アクション -->
    <div class="flex-1 min-h-[360px]">
      <DatabaseSpreadsheet :columns="cols" :rows="whitelistList" />
    </div>
  </div>
</template>
