<!-- frontend/src/components/admin/database/AccountManagementView.vue (100行以下) -->
<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import DatabaseSpreadsheet from './DatabaseSpreadsheet.vue';
import AccountDetailCard from './AccountDetailCard.vue';
import AccountHistoryTimeline from './AccountHistoryTimeline.vue';

const props = defineProps<{
  accounts: any[];
  selectedDetail: any;
  loading: boolean;
  availableAvatars?: string[];
}>();

const emit = defineEmits<{
  (e: 'selectAccount', numericId: string): void;
  (e: 'saveAccount', payload: { numericId: string; displayName: string; username: string; avatarUrl: string; description: string; aliasOf: string; groupName: string }): void;
  (e: 'mergeAccounts', sourceId: string, targetId: string): void;
  (e: 'uploadAvatar', payload: { virtualKey: string; base64Data: string }): void;
  (e: 'viewPosts', numericId: string): void;
  (e: 'viewMedia', numericId: string): void;
  (e: 'refresh'): void;
}>();

const isEditing = ref(false);
const cols = [
  { key: 'numeric_id', label: 'ID', width: '110px' }, { key: 'username', label: 'Handle', width: '110px' },
  { key: 'display_name', label: 'Name', width: '130px' }, { key: 'created_at', label: 'Created At', width: '120px' },
];

onMounted(() => { emit('refresh'); });

const handleSave = (payload: { displayName: string; username: string; avatarUrl: string; description: string; aliasOf: string; groupName: string }) => {
  if (!props.selectedDetail?.account?.numeric_id) return;
  emit('saveAccount', { numericId: props.selectedDetail.account.numeric_id, ...payload });
  isEditing.value = false;
};

const mergeCandidates = computed(() => {
  const acc = props.selectedDetail?.account;
  if (!acc || !props.accounts) return [];
  const n = (u: string) => (u || '').toLowerCase().replace(/^@/, '');
  const me = n(acc.username);
  return props.accounts.filter(a => a.numeric_id !== acc.numeric_id && me && n(a.username) && (me.startsWith(n(a.username)) || n(a.username).startsWith(me)));
});

const confirmMerge = (sourceId: string) => {
  const t = props.selectedDetail?.account;
  if (!t || !window.confirm(`[@${t.username}] へ ${sourceId} を合併しますか？`)) return;
  emit('mergeAccounts', sourceId, t.numeric_id);
};
</script>

<template>
  <div class="flex-1 flex flex-col md:flex-row min-h-0 bg-slate-950 overflow-hidden">
    <!-- 左側: アカウント一覧スプレッドシート -->
    <div class="w-full md:w-1/2 flex flex-col border-r border-slate-800 min-h-0">
      <div class="p-3 bg-slate-900 border-b border-slate-800 flex items-center justify-between">
        <h3 class="text-xs font-bold uppercase tracking-wider text-slate-300">👥 アカウント一覧 ({{ accounts?.length || 0 }}件)</h3>
        <button @click="emit('refresh')" class="px-2.5 py-1 bg-slate-800 hover:bg-slate-700 text-slate-300 rounded text-xs cursor-pointer">🔄 更新</button>
      </div>
      <div class="flex-1 overflow-y-auto">
        <DatabaseSpreadsheet
          :columns="cols"
          :rows="accounts || []"
          :selected-row-id="selectedDetail?.account?.numeric_id"
          id-key="numeric_id"
          @select-row="emit('selectAccount', $event.numeric_id)"
        />
      </div>
    </div>

    <!-- 右側: 選択中アカウントの詳細カード・世代履歴 -->
    <div class="w-full md:w-1/2 flex flex-col p-4 space-y-4 overflow-y-auto bg-slate-950">
      <div v-if="loading && !selectedDetail" class="text-xs text-slate-500 py-8 text-center">アカウント情報を読み込み中...</div>
      <div v-else-if="!selectedDetail || !selectedDetail.account" class="text-xs text-slate-500 py-8 text-center">左側の一覧からアカウントを選択してください</div>
      <template v-else>
        <AccountDetailCard
          :account="selectedDetail.account"
          :post-count="selectedDetail.post_count || 0"
          :is-editing="isEditing"
          :available-avatars="availableAvatars"
          @start-edit="isEditing = true"
          @cancel-edit="isEditing = false"
          @save="handleSave"
          @view-posts="emit('viewPosts', selectedDetail.account.numeric_id)"
          @view-media="emit('viewMedia', selectedDetail.account.numeric_id)"
         />
        <div v-if="mergeCandidates.length" class="p-2 bg-slate-900 border border-slate-700 rounded-lg text-xs">
          <div class="text-xs font-bold text-slate-300 mb-1">🔗 名寄せ候補 ({{ mergeCandidates.length }})</div>
          <div v-for="c in mergeCandidates" :key="c.numeric_id" class="flex justify-between items-center py-0.5 border-t border-slate-800 first:border-t-0">
            <span class="text-slate-400">@{{ c.username }} （{{ c.display_name }}） {{ c.numeric_id }}</span>
            <button @click="confirmMerge(c.numeric_id)" class="px-2 py-0.5 bg-amber-900/40 text-amber-300 rounded text-[10px]">合併</button>
          </div>
        </div>
        <AccountHistoryTimeline
          :histories="selectedDetail.histories || []"
          :username="selectedDetail.account?.username || ''"
          @upload-avatar="emit('uploadAvatar', $event)"
        />
      </template>
    </div>

  </div>
</template>
