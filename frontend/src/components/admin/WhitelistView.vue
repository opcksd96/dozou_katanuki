<script setup lang="ts">
import { ref, computed, reactive, onMounted } from 'vue';

interface WhitelistItem {
  id: number;
  type: string;
  value: string;
  is_active: boolean;
}

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

// 新規作成フォーム
const newForm = reactive({
  type: 'account',
  value: '',
});

// インライン編集状態
const editingId = ref<number | null>(null);
const editForm = reactive({
  type: 'account',
  value: '',
  is_active: true,
});

// フィルタリングされたリスト
const filteredList = computed(() => {
  return props.whitelistList.filter((item) => {
    const matchType = filterType.value === 'all' || item.type === filterType.value;
    const matchQuery =
      !searchQuery.value.trim() ||
      item.value.toLowerCase().includes(searchQuery.value.trim().toLowerCase());
    return matchType && matchQuery;
  });
});

const handleAdd = () => {
  if (!newForm.value.trim()) return;
  emit('add', newForm.type, newForm.value.trim());
  newForm.value = '';
};

const startEdit = (item: WhitelistItem) => {
  editingId.value = item.id;
  editForm.type = item.type;
  editForm.value = item.value;
  editForm.is_active = item.is_active;
};

const cancelEdit = () => {
  editingId.value = null;
};

const saveEdit = (id: number) => {
  if (!editForm.value.trim()) return;
  emit('update', id, editForm.type, editForm.value.trim(), editForm.is_active);
  editingId.value = null;
};

const handleDelete = (item: WhitelistItem) => {
  if (confirm(`Whitelist 項目「${item.value}」(${item.type}) を削除しますか？`)) {
    emit('delete', item.id);
  }
};

onMounted(() => {
  emit('fetch');
});
</script>

<template>
  <div class="space-y-6">
    <!-- ヘッダー説明 ＆ ステータス通知 -->
    <div class="flex items-center justify-between">
      <div>
        <h3 class="text-base font-bold text-slate-100 flex items-center gap-2">
          <span>📋</span> Whitelist ガバナンス管理
          <span class="text-[10px] font-mono bg-blue-900/40 text-blue-300 border border-blue-700/50 px-2 py-0.5 rounded">
            whitelists テーブル連携
          </span>
        </h3>
        <p class="text-xs text-slate-400 mt-0.5">
          自動サルベージおよびタイムラインフィルタリングの対象アカウント・検索キーワードを統治します。
        </p>
      </div>
      <button
        @click="emit('fetch')"
        :disabled="loading"
        class="px-3 py-1.5 bg-slate-800 hover:bg-slate-700 text-slate-300 text-xs rounded-lg transition-colors flex items-center gap-1.5 border border-slate-700 disabled:opacity-50"
      >
        <span :class="{ 'animate-spin': loading }">🔄</span>
        更新
      </button>
    </div>

    <!-- ステータスフィードバック -->
    <div
      v-if="statusMessage"
      class="p-3 rounded-lg text-xs font-semibold flex items-center gap-2 transition-all"
      :class="
        statusMessage.success
          ? 'bg-emerald-950/60 border border-emerald-500/30 text-emerald-300'
          : 'bg-rose-950/60 border border-rose-500/30 text-rose-300'
      "
    >
      <span>{{ statusMessage.success ? '✅' : '⚠️' }}</span>
      <span>{{ statusMessage.message }}</span>
    </div>

    <!-- 新規登録カード -->
    <div class="p-4 bg-slate-900/60 border border-slate-800 rounded-xl space-y-3">
      <h4 class="text-xs font-bold text-slate-300 flex items-center gap-1.5">
        <span>➕</span> 新規 Whitelist ルールの追加
      </h4>
      <form @submit.prevent="handleAdd" class="flex flex-wrap items-center gap-3">
        <!-- 種別選択 -->
        <div class="w-36">
          <select
            v-model="newForm.type"
            class="w-full bg-slate-950 border border-slate-700 rounded-lg px-3 py-2 text-xs text-slate-200 focus:outline-none focus:border-blue-500"
          >
            <option value="account">👤 アカウント (Handle)</option>
            <option value="keyword">🔍 キーワード (Keyword)</option>
          </select>
        </div>

        <!-- 値入力 -->
        <div class="flex-1 min-w-[200px]">
          <input
            v-model="newForm.value"
            type="text"
            :placeholder="newForm.type === 'account' ? '例: msluo14 (@なし)' : '例: famicom, apu'"
            class="w-full bg-slate-950 border border-slate-700 rounded-lg px-3 py-2 text-xs text-slate-200 placeholder-slate-500 focus:outline-none focus:border-blue-500"
          />
        </div>

        <!-- 追加ボタン -->
        <button
          type="submit"
          :disabled="!newForm.value.trim()"
          class="px-4 py-2 bg-blue-600 hover:bg-blue-500 text-white text-xs font-bold rounded-lg transition-colors shadow-lg shadow-blue-600/20 disabled:opacity-50 disabled:cursor-not-allowed"
        >
          追加する
        </button>
      </form>
    </div>

    <!-- 検索 ＆ フィルタバー -->
    <div class="flex flex-wrap items-center justify-between gap-3 pt-2">
      <!-- 種別フィルタ -->
      <div class="flex items-center gap-1 bg-slate-900/80 p-1 rounded-lg border border-slate-800 text-xs">
        <button
          @click="filterType = 'all'"
          class="px-3 py-1 rounded-md transition-all font-semibold"
          :class="filterType === 'all' ? 'bg-blue-600 text-white' : 'text-slate-400 hover:text-slate-200'"
        >
          すべて ({{ whitelistList.length }})
        </button>
        <button
          @click="filterType = 'account'"
          class="px-3 py-1 rounded-md transition-all font-semibold flex items-center gap-1"
          :class="filterType === 'account' ? 'bg-blue-600 text-white' : 'text-slate-400 hover:text-slate-200'"
        >
          <span>👤</span> アカウント ({{ whitelistList.filter((i) => i.type === 'account').length }})
        </button>
        <button
          @click="filterType = 'keyword'"
          class="px-3 py-1 rounded-md transition-all font-semibold flex items-center gap-1"
          :class="filterType === 'keyword' ? 'bg-purple-600 text-white' : 'text-slate-400 hover:text-slate-200'"
        >
          <span>🔍</span> キーワード ({{ whitelistList.filter((i) => i.type === 'keyword').length }})
        </button>
      </div>

      <!-- 検索入力 -->
      <div class="w-64">
        <input
          v-model="searchQuery"
          type="text"
          placeholder="リスト内を検索..."
          class="w-full bg-slate-900 border border-slate-800 rounded-lg px-3 py-1.5 text-xs text-slate-200 placeholder-slate-500 focus:outline-none focus:border-blue-500"
        />
      </div>
    </div>

    <!-- リストテーブル -->
    <div class="border border-slate-800 rounded-xl overflow-hidden bg-slate-900/40">
      <table class="w-full text-left text-xs border-collapse">
        <thead>
          <tr class="border-b border-slate-800 bg-slate-900/80 text-slate-400 font-mono">
            <th class="py-2.5 px-4 w-14">ID</th>
            <th class="py-2.5 px-4 w-28">種別</th>
            <th class="py-2.5 px-4">ルール値 / ハンドル</th>
            <th class="py-2.5 px-4 w-28 text-center">稼働状態</th>
            <th class="py-2.5 px-4 w-36 text-right">アクション</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-slate-800/60">
          <tr v-if="filteredList.length === 0">
            <td colspan="5" class="py-8 text-center text-slate-500 text-xs">
              <span v-if="loading">データを読み込み中...</span>
              <span v-else>該当する Whitelist 項目がありません。</span>
            </td>
          </tr>

          <tr
            v-for="item in filteredList"
            :key="item.id"
            class="hover:bg-slate-800/30 transition-colors"
            :class="{ 'bg-blue-950/20': editingId === item.id }"
          >
            <!-- ID -->
            <td class="py-3 px-4 font-mono text-slate-500">
              #{{ item.id }}
            </td>

            <!-- 種別 -->
            <td class="py-3 px-4">
              <span
                v-if="editingId !== item.id"
                class="px-2 py-0.5 rounded text-[11px] font-semibold border flex items-center gap-1 w-fit"
                :class="
                  item.type === 'account'
                    ? 'bg-blue-950/60 text-blue-300 border-blue-800/50'
                    : 'bg-purple-950/60 text-purple-300 border-purple-800/50'
                "
              >
                <span>{{ item.type === 'account' ? '👤' : '🔍' }}</span>
                {{ item.type }}
              </span>
              <select
                v-else
                v-model="editForm.type"
                class="bg-slate-950 border border-slate-700 rounded px-2 py-1 text-xs text-slate-200"
              >
                <option value="account">account</option>
                <option value="keyword">keyword</option>
              </select>
            </td>

            <!-- ルール値 -->
            <td class="py-3 px-4 font-mono">
              <div v-if="editingId !== item.id" class="text-slate-200 font-semibold flex items-center gap-1.5">
                <span v-if="item.type === 'account'" class="text-slate-500">@</span>
                <span>{{ item.value }}</span>
              </div>
              <input
                v-else
                v-model="editForm.value"
                type="text"
                class="w-full bg-slate-950 border border-blue-500 rounded px-2 py-1 text-xs text-white focus:outline-none"
              />
            </td>

            <!-- 有効/無効トグル -->
            <td class="py-3 px-4 text-center">
              <button
                @click="emit('toggle', item.id)"
                class="px-2.5 py-1 rounded-full text-[11px] font-semibold border transition-all flex items-center justify-center gap-1 mx-auto"
                :class="
                  item.is_active
                    ? 'bg-emerald-950/60 text-emerald-300 border-emerald-700/50 hover:bg-emerald-900/60'
                    : 'bg-slate-800/80 text-slate-500 border-slate-700 hover:text-slate-300'
                "
              >
                <span class="w-1.5 h-1.5 rounded-full" :class="item.is_active ? 'bg-emerald-400' : 'bg-slate-500'"></span>
                {{ item.is_active ? '有効' : '停止中' }}
              </button>
            </td>

            <!-- アクション -->
            <td class="py-3 px-4 text-right">
              <div v-if="editingId !== item.id" class="flex items-center justify-end gap-1.5">
                <button
                  @click="startEdit(item)"
                  class="px-2 py-1 bg-slate-800 hover:bg-slate-700 text-slate-300 rounded text-[11px] transition-colors"
                >
                  編集
                </button>
                <button
                  @click="handleDelete(item)"
                  class="px-2 py-1 bg-rose-950/40 hover:bg-rose-900/60 text-rose-300 border border-rose-800/40 rounded text-[11px] transition-colors"
                >
                  削除
                </button>
              </div>
              <div v-else class="flex items-center justify-end gap-1.5">
                <button
                  @click="saveEdit(item.id)"
                  class="px-2 py-1 bg-blue-600 hover:bg-blue-500 text-white rounded text-[11px] font-bold transition-colors"
                >
                  保存
                </button>
                <button
                  @click="cancelEdit"
                  class="px-2 py-1 bg-slate-800 hover:bg-slate-700 text-slate-400 rounded text-[11px] transition-colors"
                >
                  取消
                </button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>
