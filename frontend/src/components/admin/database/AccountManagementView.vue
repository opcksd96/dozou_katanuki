<!-- frontend/src/components/admin/database/AccountManagementView.vue (100行以下 - SPEC-PRINCIPLE-001) -->
<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import DatabaseSpreadsheet from './DatabaseSpreadsheet.vue';
import AccountDetailCard from './AccountDetailCard.vue';
import AccountHistoryTimeline from './AccountHistoryTimeline.vue';

const props = defineProps<{ accounts: any[]; selectedDetail: any; loading: boolean; availableAvatars?: string[]; }>();
const emit = defineEmits<{
  (e: 'selectAccount', numericId: string): void;
  (e: 'saveAccount', payload: { numericId: string; displayName: string; username: string; avatarUrl: string; description: string; aliasOf: string; groupName: string }): void;
  (e: 'mergeAccounts', sourceId: string, targetId: string): void;
  (e: 'uploadAvatar', payload: { virtualKey: string; base64Data: string }): void;
  (e: 'toggleWhitelist', numericId: string, isWhitelist: boolean): void;
  (e: 'viewPosts', numericId: string): void; (e: 'viewMedia', numericId: string): void; (e: 'refresh'): void;
}>();

const isEditing = ref(false);
const cols = [
  { key: 'is_whitelist', label: '🛡️ 巡回', width: '65px' },
  { key: 'username', label: 'Handle', width: '110px' }, { key: 'display_name', label: 'Name', width: '110px' },
  { key: 'numeric_id', label: 'ID', width: '90px' },
];

onMounted(() => emit('refresh'));
const handleSave = (payload: any) => {
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
  <div class="flex-1 flex flex-col md:flex-row min-h-0 bg-slate-950 overflow-y-auto md:overflow-hidden font-sans">
    <!-- 左側 / 上段: アカウント一覧スプレッドシート -->
    <div class="w-full md:w-1/2 flex flex-col border-b md:border-b-0 md:border-r border-slate-800 h-80 md:h-full shrink-0 md:shrink min-h-0">
      <div class="p-2.5 bg-slate-900 border-b border-slate-800 flex items-center justify-between">
        <h3 class="text-xs font-bold uppercase text-slate-300">👥 アカウント・巡回一覧 ({{ accounts?.length || 0 }}件)</h3>
        <button @click="emit('refresh')" class="px-2.5 py-1 bg-slate-800 hover:bg-slate-700 text-slate-300 rounded text-xs cursor-pointer active:scale-95">🔄 更新</button>
      </div>
      <div class="flex-1 min-h-0 overflow-y-auto">
        <DatabaseSpreadsheet :columns="cols" :rows="accounts" :selected-row-id="selectedDetail?.account?.numeric_id" :id-key="'numeric_id'" @select-row="emit('selectAccount', $event?.numeric_id || $event)">
          <template #cell-is_whitelist="{ row }">
            <button @click.stop="emit('toggleWhitelist', row.numeric_id, !row.is_whitelist)" class="px-2 py-0.5 rounded text-[10px] font-bold font-mono transition-all cursor-pointer shadow-sm" :class="row.is_whitelist ? 'bg-emerald-950 text-emerald-300 border border-emerald-700/60' : 'bg-slate-800 text-slate-400 hover:text-slate-200'" :title="row.is_whitelist ? '巡回・保存対象 (ON)' : '外部参照のみ (OFF)'">
              {{ row.is_whitelist ? '🛡️ ON' : '⚪ OFF' }}
            </button>
          </template>
        </DatabaseSpreadsheet>
      </div>
    </div>

    <!-- 右側 / 下段: アカウント詳細・名寄せ・歴代変遷史 -->
    <div class="w-full md:w-1/2 flex flex-col min-h-0 overflow-y-auto p-3 sm:p-4 space-y-4 bg-slate-950">
      <div v-if="loading" class="text-center py-12 text-xs text-slate-500">アカウント詳細を読み込み中...</div>
      <template v-else-if="selectedDetail?.account">
        <AccountDetailCard :account="selectedDetail.account" :post-count="selectedDetail.post_count || 0" :is-editing="isEditing" :available-avatars="availableAvatars" @start-edit="isEditing = true" @cancel-edit="isEditing = false" @save="handleSave" @toggle-whitelist="emit('toggleWhitelist', selectedDetail.account.numeric_id, !selectedDetail.account.is_whitelist)" @view-posts="emit('viewPosts', selectedDetail.account.numeric_id)" @view-media="emit('viewMedia', selectedDetail.account.numeric_id)" />
        <div v-if="mergeCandidates.length" class="p-3 bg-amber-950/40 border border-amber-800/60 rounded-xl space-y-2 text-xs">
          <div class="font-bold text-amber-300">💡 名寄せ候補アカウント</div>
          <div class="space-y-1.5">
            <div v-for="cand in mergeCandidates" :key="cand.numeric_id" class="flex items-center justify-between bg-black/40 p-2 rounded">
              <span class="text-slate-300 font-mono">@{{ cand.username }} ({{ cand.display_name }})</span>
              <button @click="confirmMerge(cand.numeric_id)" class="px-2.5 py-1 bg-amber-600 hover:bg-amber-500 text-white rounded font-bold text-[11px]">🔗 合併する</button>
            </div>
          </div>
        </div>
        <AccountHistoryTimeline :histories="selectedDetail.histories || selectedDetail.profile_histories || []" :username="selectedDetail.account.username" @upload-avatar="emit('uploadAvatar', $event)" />
      </template>
      <div v-else class="text-center py-12 text-xs text-slate-500">アカウントを選択してください</div>
    </div>
  </div>
</template>

