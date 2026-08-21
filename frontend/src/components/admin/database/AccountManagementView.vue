<!-- frontend/src/components/admin/database/AccountManagementView.vue (100行以下) -->
<script setup lang="ts">
import { onMounted } from 'vue';
import DatabaseSpreadsheet from './DatabaseSpreadsheet.vue';

const props = defineProps<{ accounts: any[]; selectedDetail: any; loading: boolean }>();
const emit = defineEmits<{
  (e: 'selectAccount', numericId: string): void;
  (e: 'viewPosts', numericId: string): void;
  (e: 'viewMedia', numericId: string): void;
  (e: 'refresh'): void;
}>();

onMounted(() => { emit('refresh'); });
const cols = [
  { key: 'numeric_id', label: 'Numeric ID', width: '130px' },
  { key: 'username', label: 'Username', width: '110px' },
  { key: 'display_name', label: 'Display Name', width: '130px' },
  { key: 'updated_at', label: 'Updated At', width: '140px' },
];
</script>

<template>
  <div class="grid grid-cols-1 lg:grid-cols-12 gap-5 h-full">
    <!-- 左側: アカウント選択データテーブル -->
    <div class="lg:col-span-6 flex flex-col space-y-2">
      <div class="flex justify-between items-center text-xs text-slate-300 font-semibold px-1">
        <span>📋 登録アカウント一覧 (全 {{ accounts.length }} 件)</span>
        <button @click="emit('refresh')" class="text-blue-400 hover:text-blue-300">🔄 更新</button>
      </div>
      <div class="flex-1 min-h-[350px]">
        <DatabaseSpreadsheet :columns="cols" :rows="accounts" :selected-row-id="selectedDetail?.account?.numeric_id" id-key="numeric_id" @select-row="(r) => emit('selectAccount', r.numeric_id)" />
      </div>
    </div>

    <!-- 右側: アカウントメディアカード ＆ 世代アバター履歴 -->
    <div class="lg:col-span-6 flex flex-col space-y-3 bg-slate-900/40 border border-slate-800 rounded-xl p-4 overflow-y-auto max-h-[500px]">
      <div v-if="selectedDetail" class="space-y-4">
        <!-- アカウント概要カード -->
        <div class="bg-slate-950 p-4 rounded-xl border border-slate-800 shadow-md space-y-3">
          <div class="flex items-center gap-4">
            <img :src="selectedDetail.account.avatar_url || '/assets/default-avatar.png'" class="w-16 h-16 rounded-full border-2 border-blue-500/60 object-cover bg-slate-800" />
            <div class="space-y-1 flex-1">
              <h4 class="text-sm font-bold text-slate-100">{{ selectedDetail.account.display_name }}</h4>
              <div class="text-xs font-mono text-blue-400">@{{ selectedDetail.account.username }}</div>
              <div class="text-[11px] font-mono text-slate-400">ID: {{ selectedDetail.account.numeric_id }} | 投稿数: {{ selectedDetail.post_count }} 件</div>
            </div>
          </div>
          <!-- クイックアクション導線 -->
          <div class="flex gap-2 pt-2 border-t border-slate-850">
            <button @click="emit('viewPosts', selectedDetail.account.numeric_id)" class="flex-1 py-1.5 px-3 bg-blue-600/80 hover:bg-blue-600 text-white rounded-lg text-xs font-semibold flex items-center justify-center gap-1">
              <span>📝</span> 投稿一覧を表示
            </button>
            <button @click="emit('viewMedia', selectedDetail.account.numeric_id)" class="flex-1 py-1.5 px-3 bg-indigo-600/80 hover:bg-indigo-600 text-white rounded-lg text-xs font-semibold flex items-center justify-center gap-1">
              <span>🖼️</span> メディア一覧を表示
            </button>
          </div>
        </div>

        <!-- 世代別アバター変更履歴 -->
        <div class="space-y-2">
          <h5 class="text-xs font-bold text-slate-200 flex items-center gap-1.5"><span>🖼️</span> アバター世代履歴 (account_profile_histories)</h5>
          <div v-if="!selectedDetail.histories || selectedDetail.histories.length === 0" class="p-6 text-center text-xs text-slate-500 bg-slate-950/60 rounded-xl border border-slate-800/60">世代履歴データはありません</div>
          <div v-else class="grid grid-cols-1 sm:grid-cols-2 gap-3">
            <div v-for="h in selectedDetail.histories" :key="h.id" class="p-3 bg-slate-950 border border-slate-800 rounded-lg flex items-center gap-3">
              <img :src="h.avatar_original_url || selectedDetail.account.avatar_url" class="w-10 h-10 rounded-full border border-slate-700 object-cover bg-slate-800" />
              <div class="text-[11px] font-mono space-y-0.5 truncate">
                <div class="text-emerald-400 font-bold">Seq: {{ h.avatar_seq }}</div>
                <div class="text-slate-300 truncate" :title="h.avatar_virtual_key">{{ h.avatar_virtual_key }}</div>
                <div class="text-[10px] text-slate-500">{{ new Date(h.observed_at).toLocaleDateString() }}</div>
              </div>
            </div>
          </div>
        </div>
      </div>
      <div v-else class="h-full flex items-center justify-center text-xs text-slate-500">左側の一覧からアカウントを選択してください</div>
    </div>
  </div>
</template>
