<!-- frontend/src/components/admin/database/AccountManagementView.vue (100行以下 - SPEC-PRINCIPLE-001) -->
<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import DatabaseSpreadsheet from './DatabaseSpreadsheet.vue';
import AccountDetailCard from './AccountDetailCard.vue';
import AccountHistoryTimeline from './AccountHistoryTimeline.vue';
import AccountTrashModal from './AccountTrashModal.vue';

const props = defineProps<{ accounts: any[]; selectedDetail: any; loading: boolean; availableAvatars?: string[]; showTrash?: boolean; }>();
const emit = defineEmits<{
  (e: 'selectAccount', numericId: string): void; (e: 'saveAccount', payload: any): void; (e: 'mergeAccounts', sourceId: string, targetId: string): void;
  (e: 'uploadAvatar', payload: { virtualKey: string; base64Data: string }): void; (e: 'toggleWhitelist', numericId: string, isWhitelist: boolean): void;
  (e: 'viewPosts', numericId: string): void; (e: 'viewMedia', numericId: string): void; (e: 'refresh'): void;
  (e: 'trashAccount', payload: { numericId: string; reason: string }): void; (e: 'restoreAccount', id: string): void; (e: 'toggleShowTrash'): void;
}>();

const isEditing = ref(false), showTrashModal = ref(false);
const cols = [
  { key: 'is_whitelist', label: '🛡️ 巡回', width: '65px' }, { key: 'username', label: 'Handle', width: '95px' },
  { key: 'post_count', label: '投稿数', width: '65px' }, { key: 'display_name', label: 'Name', width: '95px' }, { key: 'actions', label: '動線', width: '80px' },
];

onMounted(() => emit('refresh'));
const handleSave = (p: any) => { if (!props.selectedDetail?.account?.numeric_id) return; emit('saveAccount', { numericId: props.selectedDetail.account.numeric_id, ...p }); isEditing.value = false; };
const mergeCandidates = computed(() => {
  const acc = props.selectedDetail?.account; if (!acc || !props.accounts) return [];
  const n = (u: string) => (u || '').toLowerCase().replace(/^@/, ''), me = n(acc.username);
  return props.accounts.filter(a => !a.is_trash && a.numeric_id !== acc.numeric_id && me && n(a.username) && (me.startsWith(n(a.username)) || n(a.username).startsWith(me)));
});
const confirmMerge = (cand: any) => {
  const t = props.selectedDetail?.account, pCount = props.selectedDetail?.post_count ?? t?.post_count ?? 0;
  if (!t || !cand) return;
  const msg = `【アカウント名寄せ・統合の確認】\n\n🔗 統合先 (親・存続):\n  ・@${t.username} (${t.display_name || '未設定'})\n  ・投稿数: ${pCount}件 (ID: ${t.numeric_id})\n\n📥 合併元 (子・吸収):\n  ・@${cand.username} (${cand.display_name || '未設定'})\n  ・投稿数: ${cand.post_count || 0}件 (ID: ${cand.numeric_id})\n\n合併元の全データを統合先へ移管・合体します。実行しますか？`;
  if (window.confirm(msg)) emit('mergeAccounts', cand.numeric_id, t.numeric_id);
};
</script>

<template>
  <div class="flex-1 flex flex-col md:flex-row min-h-0 bg-slate-950 overflow-y-auto md:overflow-hidden font-sans">
    <!-- 左側: アカウント一覧 -->
    <div class="w-full md:w-1/2 flex flex-col border-b md:border-b-0 md:border-r border-slate-800 h-80 md:h-full shrink-0 md:shrink min-h-0">
      <div class="p-2.5 bg-slate-900 border-b border-slate-800 flex items-center justify-between">
        <h3 class="text-xs font-bold uppercase text-slate-300">{{ showTrash ? '🗑️ ゴミ箱アカウント' : '👥 アカウント一覧' }} ({{ accounts?.length || 0 }}件)</h3>
        <div class="flex items-center gap-1.5">
          <button @click="emit('toggleShowTrash')" :class="showTrash ? 'bg-rose-950 text-rose-300 border-rose-800' : 'bg-slate-800 text-slate-400'" class="px-2 py-1 rounded text-xs font-bold border cursor-pointer active:scale-95" title="ゴミ箱アカウントの表示切替">{{ showTrash ? '🗑️ ゴミ箱表示中' : '🗑️ ゴミ箱' }}</button>
          <button @click="emit('refresh')" class="px-2.5 py-1 bg-slate-800 hover:bg-slate-700 text-slate-300 rounded text-xs cursor-pointer active:scale-95">🔄 更新</button>
        </div>
      </div>
      <div class="flex-1 min-h-0 overflow-y-auto">
        <DatabaseSpreadsheet :columns="cols" :rows="accounts" :selected-row-id="selectedDetail?.account?.numeric_id" :id-key="'numeric_id'" @select-row="emit('selectAccount', $event?.numeric_id || $event)">
          <template #cell-is_whitelist="{ row }">
            <button @click.stop="emit('toggleWhitelist', row.numeric_id, !row.is_whitelist)" class="px-2 py-0.5 rounded text-[10px] font-bold font-mono cursor-pointer shadow-sm" :class="row.is_whitelist ? 'bg-emerald-950 text-emerald-300 border border-emerald-700/60' : 'bg-slate-800 text-slate-400 hover:text-slate-200'">{{ row.is_whitelist ? '🛡️ ON' : '⚪ OFF' }}</button>
          </template>
          <template #cell-post_count="{ row }"><span class="font-mono text-blue-400 font-bold text-[11px]">{{ row.post_count || 0 }}件</span></template>
          <template #cell-actions="{ row }">
            <div class="flex items-center gap-1" @click.stop>
              <button @click="emit('viewPosts', row.numeric_id)" class="px-1.5 py-0.5 bg-blue-900/60 hover:bg-blue-700 text-blue-200 rounded text-[10px] font-bold cursor-pointer" title="投稿一覧">📝</button>
              <button @click="emit('viewMedia', row.numeric_id)" class="px-1.5 py-0.5 bg-indigo-900/60 hover:bg-indigo-700 text-indigo-200 rounded text-[10px] font-bold cursor-pointer" title="メディア一覧">🖼️</button>
            </div>
          </template>
        </DatabaseSpreadsheet>
      </div>
    </div>

    <!-- 右側: アカウント詳細 -->
    <div class="w-full md:w-1/2 flex flex-col min-h-0 overflow-y-auto p-3 sm:p-4 space-y-4 bg-slate-950">
      <div v-if="loading" class="text-center py-12 text-xs text-slate-500">読み込み中...</div>
      <template v-else-if="selectedDetail?.account">
        <AccountDetailCard :account="selectedDetail.account" :post-count="selectedDetail.post_count ?? selectedDetail.PostCount ?? selectedDetail.account?.post_count ?? 0" :is-editing="isEditing" :available-avatars="availableAvatars" @start-edit="isEditing = true" @cancel-edit="isEditing = false" @save="handleSave" @toggle-whitelist="emit('toggleWhitelist', selectedDetail.account.numeric_id, !selectedDetail.account.is_whitelist)" @view-posts="emit('viewPosts', selectedDetail.account.numeric_id)" @view-media="emit('viewMedia', selectedDetail.account.numeric_id)" @trash-account="showTrashModal = true" @restore-account="emit('restoreAccount', $event)" />
        <div v-if="mergeCandidates.length" class="p-3 bg-amber-950/40 border border-amber-800/60 rounded-xl space-y-2 text-xs">
          <div class="font-bold text-amber-300 flex items-center justify-between"><span>💡 名寄せ候補アカウント</span><span class="text-[10px] text-amber-400 font-mono">統合先: @{{ selectedDetail.account.username }} ({{ selectedDetail.post_count ?? 0 }}件)</span></div>
          <div class="space-y-1.5">
            <div v-for="cand in mergeCandidates" :key="cand.numeric_id" class="flex items-center justify-between bg-black/50 p-2.5 rounded-lg border border-amber-900/40 gap-2">
              <div class="min-w-0 flex-1"><div class="flex items-center gap-1.5 flex-wrap"><span class="text-amber-200 font-mono font-bold">@{{ cand.username }}</span><span class="text-slate-300 text-[11px] truncate">({{ cand.display_name }})</span><span class="px-1.5 py-0.2 bg-blue-950 text-blue-300 border border-blue-800/60 rounded text-[10px] font-mono">{{ cand.post_count || 0 }}件</span></div><div class="text-[10px] text-slate-500 font-mono truncate">ID: {{ cand.numeric_id }}</div></div>
              <button @click="confirmMerge(cand)" class="px-2.5 py-1 bg-amber-600 hover:bg-amber-500 text-white rounded font-bold text-[11px] shrink-0 active:scale-95 cursor-pointer shadow-sm">📥 合併</button>
            </div>
          </div>
        </div>
        <AccountHistoryTimeline :histories="selectedDetail.histories || selectedDetail.profile_histories || []" :username="selectedDetail.account.username" @upload-avatar="emit('uploadAvatar', $event)" />
      </template>
      <div v-else class="text-center py-12 text-xs text-slate-500">アカウントを選択してください</div>
    </div>
    <AccountTrashModal :show="showTrashModal" :account="selectedDetail?.account" @close="showTrashModal = false" @confirm="(p) => emit('trashAccount', p)" />
  </div>
</template>
