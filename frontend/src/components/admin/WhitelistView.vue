<!-- frontend/src/components/admin/WhitelistView.vue (100行以下) -->
<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import WhitelistAddForm from './whitelist/WhitelistAddForm.vue';
import WhitelistTable from './whitelist/WhitelistTable.vue';

interface WhitelistItem { id: number; type: string; value: string; is_active: boolean; }

const props = defineProps<{
  whitelistList: WhitelistItem[];
  loading: boolean;
  statusMessage: { success: boolean; message: string } | null;
}>();

const emit = defineEmits<{
  (e: 'fetch'): void;
  (e: 'add', type: string, value: string): void;
  (e: 'update', id: number, type: string, value: string, isActive: boolean): void;
  (e: 'delete', id: number): void;
  (e: 'toggle', id: number): void;
}>();

const filterType = ref<'all' | 'account' | 'keyword'>('all');
const searchQuery = ref('');

const filteredList = computed(() => {
  return props.whitelistList.filter((item) => {
    const matchType = filterType.value === 'all' || item.type === filterType.value;
    const matchQuery = !searchQuery.value.trim() || item.value.toLowerCase().includes(searchQuery.value.trim().toLowerCase());
    return matchType && matchQuery;
  });
});

onMounted(() => { emit('fetch'); });
</script>

<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h3 class="text-base font-bold text-slate-100 flex items-center gap-2">
          <span>📋</span> Whitelist ガバナンス管理
          <span class="text-[10px] font-mono bg-blue-900/40 text-blue-300 border border-blue-700/50 px-2 py-0.5 rounded">whitelists テーブル連携</span>
        </h3>
        <p class="text-xs text-slate-400 mt-0.5">自動サルベージおよびタイムラインの対象アカウント・検索キーワードを統治します。</p>
      </div>
      <button @click="emit('fetch')" :disabled="loading" class="px-3 py-1.5 bg-slate-800 hover:bg-slate-700 text-slate-300 text-xs rounded-lg flex items-center gap-1.5 border border-slate-700 disabled:opacity-50">
        <span :class="{ 'animate-spin': loading }">🔄</span> 更新
      </button>
    </div>

    <div v-if="statusMessage" class="p-3 rounded-lg text-xs font-semibold flex items-center gap-2" :class="statusMessage.success ? 'bg-emerald-950/60 border border-emerald-500/30 text-emerald-300' : 'bg-rose-950/60 border border-rose-500/30 text-rose-300'">
      <span>{{ statusMessage.success ? '✅' : '⚠️' }}</span><span>{{ statusMessage.message }}</span>
    </div>

    <WhitelistAddForm @add="(t, v) => emit('add', t, v)" />

    <div class="flex flex-wrap items-center justify-between gap-3 pt-2">
      <div class="flex items-center gap-1 bg-slate-900/80 p-1 rounded-lg border border-slate-800 text-xs">
        <button @click="filterType = 'all'" class="px-3 py-1 rounded-md font-semibold" :class="filterType === 'all' ? 'bg-blue-600 text-white' : 'text-slate-400'">すべて ({{ whitelistList.length }})</button>
        <button @click="filterType = 'account'" class="px-3 py-1 rounded-md font-semibold" :class="filterType === 'account' ? 'bg-blue-600 text-white' : 'text-slate-400'">👤 アカウント ({{ whitelistList.filter((i) => i.type === 'account').length }})</button>
        <button @click="filterType = 'keyword'" class="px-3 py-1 rounded-md font-semibold" :class="filterType === 'keyword' ? 'bg-purple-600 text-white' : 'text-slate-400'">🔍 キーワード ({{ whitelistList.filter((i) => i.type === 'keyword').length }})</button>
      </div>
      <div class="w-64">
        <input v-model="searchQuery" type="text" placeholder="リスト内を検索..." class="w-full bg-slate-900 border border-slate-800 rounded-lg px-3 py-1.5 text-xs text-slate-200 placeholder-slate-500" />
      </div>
    </div>

    <WhitelistTable
      :items="filteredList"
      :loading="loading"
      @update="(id, t, v, act) => emit('update', id, t, v, act)"
      @delete="(id) => emit('delete', id)"
      @toggle="(id) => emit('toggle', id)"
    />
  </div>
</template>
