<!-- frontend/src/components/admin/database/AccountManagementView.vue (100行以下) -->
<script setup lang="ts">
import { ref, onMounted } from 'vue';
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
  (e: 'saveAccount', payload: { numericId: string; displayName: string; username: string; avatarUrl: string; description: string }): void;
  (e: 'uploadAvatar', payload: { virtualKey: string; base64Data: string }): void;
  (e: 'viewPosts', numericId: string): void;
  (e: 'viewMedia', numericId: string): void;
  (e: 'refresh'): void;
}>();

const isEditing = ref(false);
const cols = [
  { key: 'numeric_id', label: 'ID', width: '110px' },
  { key: 'username', label: 'Handle', width: '110px' },
  { key: 'display_name', label: 'Name', width: '130px' },
  { key: 'created_at', label: 'Created At', width: '120px' },
];

onMounted(() => { emit('refresh'); });

const handleSave = (payload: { displayName: string; username: string; avatarUrl: string; description: string }) => {
  if (!props.selectedDetail?.account?.numeric_id) return;
  emit('saveAccount', { numericId: props.selectedDetail.account.numeric_id, ...payload });
  isEditing.value = false;
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
        <AccountHistoryTimeline
          :histories="selectedDetail.histories || []"
          :username="selectedDetail.account?.username || ''"
          @upload-avatar="emit('uploadAvatar', $event)"
        />
      </template>
    </div>

  </div>
</template>
