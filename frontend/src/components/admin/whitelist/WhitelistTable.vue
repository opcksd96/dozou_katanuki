<!-- frontend/src/components/admin/whitelist/WhitelistTable.vue (100行以下) -->
<script setup lang="ts">
import { ref, reactive } from 'vue';
interface WhitelistItem { id: number; type: string; value: string; group_name: string; alias_of: string; is_active: boolean; }
const props = defineProps<{ items: WhitelistItem[]; loading: boolean }>();
const emit = defineEmits<{
  (e: 'update', id: number, type: string, value: string, groupName: string, aliasOf: string, isActive: boolean): void;
  (e: 'delete', id: number): void;
  (e: 'toggle', id: number): void;
}>();

const editingId = ref<number | null>(null);
const editForm = reactive({ type: 'account', value: '', groupName: '', aliasOf: '', is_active: true });

const startEdit = (item: WhitelistItem) => {
  editingId.value = item.id;
  editForm.type = item.type; editForm.value = item.value;
  editForm.groupName = item.group_name || ''; editForm.aliasOf = item.alias_of || '';
  editForm.is_active = item.is_active;
};
const cancelEdit = () => { editingId.value = null; };
const saveEdit = (id: number) => {
  if (!editForm.value.trim()) return;
  emit('update', id, editForm.type, editForm.value.trim(), editForm.groupName.trim(), editForm.aliasOf.trim(), editForm.is_active);
  editingId.value = null;
};
const handleDelete = (item: WhitelistItem) => {
  if (confirm(`Whitelist「${item.value}」を削除しますか？`)) emit('delete', item.id);
};
</script>

<template>
  <div class="border border-slate-800 rounded-xl overflow-hidden bg-slate-900/40">
    <table class="w-full text-left text-xs border-collapse">
      <thead>
        <tr class="border-b border-slate-800 bg-slate-900/80 text-slate-400 font-mono">
          <th class="py-2.5 px-4 w-14">ID</th><th class="py-2.5 px-4 w-24">種別</th>
          <th class="py-2.5 px-4">ルール値</th>
          <th class="py-2.5 px-4 w-28">グループ</th>
          <th class="py-2.5 px-4 w-28">エイリアス</th>
          <th class="py-2.5 px-4 w-20 text-center">稼働</th>
          <th class="py-2.5 px-4 w-32 text-right">アクション</th>
        </tr>
      </thead>
      <tbody class="divide-y divide-slate-800/60">
        <tr v-if="items.length === 0">
          <td colspan="7" class="py-8 text-center text-slate-500">{{ loading ? '読込中...' : '該当項目なし' }}</td>
        </tr>
        <tr v-for="item in items" :key="item.id" class="hover:bg-slate-800/30" :class="{ 'bg-blue-950/20': editingId === item.id }">
          <td class="py-3 px-4 font-mono text-slate-500">#{{ item.id }}</td>
          <td class="py-3 px-4">
            <span v-if="editingId !== item.id" class="px-2 py-0.5 rounded text-[11px] font-semibold border flex items-center gap-1 w-fit" :class="item.type === 'account' ? 'bg-blue-950/60 text-blue-300 border-blue-800/50' : 'bg-purple-950/60 text-purple-300 border-purple-800/50'">
              {{ item.type === 'account' ? '👤' : '🔍' }} {{ item.type }}
            </span>
            <select v-else v-model="editForm.type" class="bg-slate-950 border border-slate-700 rounded px-2 py-1 text-xs text-slate-200">
              <option value="account">account</option><option value="keyword">keyword</option>
            </select>
          </td>
          <td class="py-3 px-4 font-mono">
            <span v-if="editingId !== item.id" class="text-slate-200 font-semibold">{{ item.type === 'account' ? '@' : '' }}{{ item.value }}</span>
            <input v-else v-model="editForm.value" type="text" class="w-full bg-slate-950 border border-blue-500 rounded px-2 py-1 text-xs text-white" />
          </td>
          <td class="py-3 px-4">
            <template v-if="editingId !== item.id">
              <span v-if="item.group_name" class="px-2 py-0.5 rounded-full text-[10px] font-semibold bg-amber-950/50 text-amber-300 border border-amber-700/40">🏷️ {{ item.group_name }}</span>
              <span v-else class="text-slate-600 text-[10px]">—</span>
            </template>
            <input v-else v-model="editForm.groupName" type="text" placeholder="グループ名" class="w-full bg-slate-950 border border-slate-700 rounded px-2 py-1 text-xs text-slate-300 placeholder-slate-600" />
          </td>
          <td class="py-3 px-4">
            <template v-if="editingId !== item.id">
              <span v-if="item.alias_of" class="px-2 py-0.5 rounded-full text-[10px] font-semibold bg-teal-950/50 text-teal-300 border border-teal-700/40">🔗 @{{ item.alias_of }}</span>
              <span v-else class="text-slate-600 text-[10px]">—</span>
            </template>
            <input v-else v-model="editForm.aliasOf" type="text" placeholder="本垢名" class="w-full bg-slate-950 border border-slate-700 rounded px-2 py-1 text-xs text-slate-300 placeholder-slate-600" />
          </td>
          <td class="py-3 px-4 text-center">
            <button @click="$emit('toggle', item.id)" class="px-2.5 py-1 rounded-full text-[11px] font-semibold border transition-all mx-auto" :class="item.is_active ? 'bg-emerald-950/60 text-emerald-300 border-emerald-700/50' : 'bg-slate-800/80 text-slate-500 border-slate-700'">
              {{ item.is_active ? '有効' : '停止中' }}
            </button>
          </td>
          <td class="py-3 px-4 text-right">
            <div v-if="editingId !== item.id" class="flex items-center justify-end gap-1.5">
              <button @click="startEdit(item)" class="px-2 py-1 bg-slate-800 hover:bg-slate-700 text-slate-300 rounded text-[11px]">編集</button>
              <button @click="handleDelete(item)" class="px-2 py-1 bg-rose-950/40 hover:bg-rose-900/60 text-rose-300 border border-rose-800/40 rounded text-[11px]">削除</button>
            </div>
            <div v-else class="flex items-center justify-end gap-1.5">
              <button @click="saveEdit(item.id)" class="px-2 py-1 bg-blue-600 hover:bg-blue-500 text-white rounded text-[11px] font-bold">保存</button>
              <button @click="cancelEdit" class="px-2 py-1 bg-slate-800 hover:bg-slate-700 text-slate-400 rounded text-[11px]">取消</button>
            </div>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

