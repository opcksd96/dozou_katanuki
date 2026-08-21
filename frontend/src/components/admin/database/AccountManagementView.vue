<!-- frontend/src/components/admin/database/AccountManagementView.vue (100行以下) -->
<script setup lang="ts">
import { ref, watch, onMounted } from 'vue';
import DatabaseSpreadsheet from './DatabaseSpreadsheet.vue';
import Avatar from '../../article/Avatar.vue';
import { resolveAvatarUrl, resolveHistoryAvatarUrl } from '../../../utils/avatar';

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
const isSaving = ref(false);
const showImagePicker = ref(false);

const isHistoryAvatarFound = (h: any) => Boolean(h.avatar_base64);

const editForm = ref({
  displayName: '',
  username: '',
  avatarUrl: '',
  description: '',
});

const syncFormWithDetail = () => {
  if (props.selectedDetail?.account) {
    const acc = props.selectedDetail.account;
    editForm.value = {
      displayName: acc.display_name || '',
      username: acc.username || '',
      avatarUrl: acc.avatar_url || '',
      description: acc.description || '',
    };
  }
};

watch(() => props.selectedDetail, () => {
  syncFormWithDetail();
  isEditing.value = false;
  showImagePicker.value = false;
}, { immediate: true });

onMounted(() => { emit('refresh'); });

const handleStartEdit = () => {
  syncFormWithDetail();
  isEditing.value = true;
};

const handleCancelEdit = () => {
  syncFormWithDetail();
  isEditing.value = false;
  showImagePicker.value = false;
};

const handleSelectPresetAvatar = (url: string) => {
  editForm.value.avatarUrl = url;
  showImagePicker.value = false;
};

// ファイルを Base64 に変換してアップロード
const processFile = (file: File, virtualKey: string) => {
  if (!file || !file.type.startsWith('image/')) return;
  const reader = new FileReader();
  reader.onload = (e) => {
    const base64 = e.target?.result as string;
    if (base64) {
      emit('uploadAvatar', { virtualKey, base64Data: base64 });
    }
  };
  reader.readAsDataURL(file);
};

const handleHistoryDrop = (e: DragEvent, virtualKey: string) => {
  const files = e.dataTransfer?.files;
  if (files && files.length > 0) {
    processFile(files[0], virtualKey);
  }
};

const handleHistoryFileInput = (e: Event, virtualKey: string) => {
  const input = e.target as HTMLInputElement;
  if (input.files && input.files.length > 0) {
    processFile(input.files[0], virtualKey);
    input.value = '';
  }
};

const handleEditAvatarDrop = (e: DragEvent) => {
  const files = e.dataTransfer?.files;
  if (files && files.length > 0) {
    const username = editForm.value.username || props.selectedDetail?.account?.username || 'user';
    const key = `${username}_avatar_${Date.now().toString().slice(-4)}`;
    processFile(files[0], key);
    editForm.value.avatarUrl = `/avatars/twitter/${key}.jpg`;
  }
};

const handleSave = async () => {
  if (!props.selectedDetail?.account?.numeric_id) return;
  isSaving.value = true;
  try {
    emit('saveAccount', {
      numericId: props.selectedDetail.account.numeric_id,
      displayName: editForm.value.displayName.trim(),
      username: editForm.value.username.trim(),
      avatarUrl: editForm.value.avatarUrl.trim(),
      description: editForm.value.description.trim(),
    });
    isEditing.value = false;
  } finally {
    isSaving.value = false;
  }
};

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
    <div class="lg:col-span-5 flex flex-col space-y-2">
      <div class="flex justify-between items-center text-xs text-slate-300 font-semibold px-1">
        <span>📋 登録アカウント一覧 (全 {{ accounts.length }} 件)</span>
        <button @click="emit('refresh')" class="text-blue-400 hover:text-blue-300 cursor-pointer">🔄 更新</button>
      </div>
      <div class="flex-1 min-h-[350px]">
        <DatabaseSpreadsheet :columns="cols" :rows="accounts" :selected-row-id="selectedDetail?.account?.numeric_id" id-key="numeric_id" @select-row="(r) => emit('selectAccount', r.numeric_id)" />
      </div>
    </div>

    <!-- 右側: アカウント詳細・編集カード ＆ 世代アバター履歴 -->
    <div class="lg:col-span-7 flex flex-col space-y-3 bg-slate-900/40 border border-slate-800 rounded-xl p-4 overflow-y-auto max-h-[550px]">
      <div v-if="selectedDetail" class="space-y-4">
        <!-- アカウント概要・編集カード -->
        <div class="bg-slate-950 p-4 rounded-xl border border-slate-800 shadow-md space-y-3">
          <div class="flex items-start justify-between">
            <div class="flex items-center gap-3.5 flex-1 min-w-0">
              <Avatar
                :avatar-url="isEditing ? editForm.avatarUrl : resolveAvatarUrl(selectedDetail.account, selectedDetail.histories)"
                :handle="isEditing ? editForm.username : selectedDetail.account.username"
                size-class="w-14 h-14"
                class="rounded-full border-2 border-blue-500/60 shadow-sm"
              />
              <div class="space-y-0.5 flex-1 min-w-0">
                <h4 class="text-sm font-bold text-slate-100 truncate">
                  {{ isEditing ? editForm.displayName || '(未設定)' : selectedDetail.account.display_name }}
                </h4>
                <div class="text-xs font-mono text-blue-400 truncate">
                  @{{ isEditing ? editForm.username : selectedDetail.account.username }}
                </div>
                <div class="text-[11px] font-mono text-slate-400">
                  ID: {{ selectedDetail.account.numeric_id }} | 投稿数: {{ selectedDetail.post_count }} 件
                </div>
              </div>
            </div>
            <button
              v-if="!isEditing"
              @click="handleStartEdit"
              class="px-2.5 py-1 bg-slate-800 hover:bg-slate-700 text-slate-200 border border-slate-700 rounded-lg text-xs font-semibold flex items-center gap-1 transition-colors cursor-pointer"
            >
              <span>✏️</span> 編集
            </button>
          </div>

          <!-- 編集フォーム表示 -->
          <div v-if="isEditing" class="pt-3 border-t border-slate-800 space-y-3 bg-slate-900/60 p-3 rounded-lg border border-blue-500/30">
            <div class="text-xs font-bold text-blue-400 flex items-center gap-1.5">
              <span>✏️</span> アカウント情報・アバターの編集
            </div>
            <div class="grid grid-cols-1 sm:grid-cols-2 gap-2.5">
              <div class="space-y-1">
                <label class="text-[11px] font-bold text-slate-300">表示名 (Display Name)</label>
                <input v-model="editForm.displayName" type="text" class="w-full px-2.5 py-1.5 bg-slate-950 border border-slate-700 rounded-lg text-xs text-slate-100 focus:border-blue-500 outline-none" placeholder="例: マシュ・キリエライト" />
              </div>
              <div class="space-y-1">
                <label class="text-[11px] font-bold text-slate-300">ユーザー名 (Username)</label>
                <div class="relative">
                  <span class="absolute left-2.5 top-1.5 text-xs text-slate-500">@</span>
                  <input v-model="editForm.username" type="text" class="w-full pl-6 pr-2.5 py-1.5 bg-slate-950 border border-slate-700 rounded-lg text-xs text-slate-100 focus:border-blue-500 outline-none font-mono" placeholder="username" />
                </div>
              </div>
            </div>

            <!-- ③ イメージセレクタ ＆ アバター指定エリア -->
            <div class="space-y-2">
              <div class="flex justify-between items-center">
                <label class="text-[11px] font-bold text-slate-300">アバター設定 (Avatar Selector)</label>
                <button
                  type="button"
                  @click="showImagePicker = !showImagePicker"
                  class="text-[11px] px-2 py-0.5 bg-blue-600/80 hover:bg-blue-600 text-white rounded font-semibold flex items-center gap-1 cursor-pointer"
                >
                  <span>🖼️</span> {{ showImagePicker ? 'ピッカーを閉じる' : '既存アセットから選択' }}
                </button>
              </div>

              <!-- イメージセレクタ (プレビュー一覧) -->
              <div v-if="showImagePicker" class="p-2.5 bg-slate-950 border border-slate-800 rounded-lg space-y-2">
                <div class="text-[11px] font-semibold text-slate-400 flex items-center justify-between">
                  <span>📁 assets/ 内の利用可能アバター (クリックで選択)</span>
                  <span class="text-slate-500 text-[10px]">全 {{ availableAvatars?.length || 0 }} 件</span>
                </div>
                <div v-if="!availableAvatars || availableAvatars.length === 0" class="text-xs text-slate-500 p-3 text-center">
                  利用可能なアセット画像がありません
                </div>
                <div v-else class="grid grid-cols-4 sm:grid-cols-6 gap-2 max-h-40 overflow-y-auto p-1">
                  <div
                    v-for="imgUrl in availableAvatars"
                    :key="imgUrl"
                    @click="handleSelectPresetAvatar(imgUrl)"
                    :class="editForm.avatarUrl === imgUrl ? 'border-blue-500 ring-2 ring-blue-500/40 bg-blue-950/40' : 'border-slate-800 hover:border-slate-600 bg-slate-900'"
                    class="p-1 rounded-lg border flex flex-col items-center gap-1 cursor-pointer transition-all group"
                    :title="imgUrl"
                  >
                    <img :src="imgUrl" class="w-9 h-9 rounded-full object-cover bg-slate-800" />
                    <span class="text-[9px] font-mono text-slate-400 truncate w-full text-center group-hover:text-slate-200">
                      {{ imgUrl.split('/').pop() }}
                    </span>
                  </div>
                </div>
              </div>

              <!-- 画像ドラッグ＆ドロップ登録ゾーン -->
              <div
                @dragover.prevent
                @drop.prevent="handleEditAvatarDrop"
                class="p-3 border-2 border-dashed border-slate-700 hover:border-blue-500 rounded-lg bg-slate-950/60 flex items-center justify-between gap-3 transition-colors"
              >
                <div class="flex items-center gap-2.5 flex-1 min-w-0">
                  <Avatar :avatar-url="editForm.avatarUrl" :handle="editForm.username" class="w-8 h-8 rounded-full flex-shrink-0" />
                  <div class="text-[11px] font-mono text-slate-300 truncate">
                    <span class="text-slate-500">選択中: </span>{{ editForm.avatarUrl || '(未設定 / デフォルト)' }}
                  </div>
                </div>
                <label class="px-2.5 py-1 bg-slate-800 hover:bg-slate-700 text-slate-200 rounded text-[11px] font-semibold flex items-center gap-1 cursor-pointer flex-shrink-0">
                  <span>📥</span> ドロップまたは選択
                  <input type="file" accept="image/*" class="hidden" @change="handleEditAvatarDrop($event as any)" />
                </label>
              </div>
            </div>

            <div class="space-y-1">
              <label class="text-[11px] font-bold text-slate-300">一言コメント / 自己紹介 (Bio)</label>
              <textarea v-model="editForm.description" rows="2" class="w-full px-2.5 py-1.5 bg-slate-950 border border-slate-700 rounded-lg text-xs text-slate-100 focus:border-blue-500 outline-none leading-relaxed" placeholder="アカウントの一言コメントやプロフィール自己紹介文を入力..." />
            </div>
            <div class="flex justify-end gap-2 pt-1">
              <button @click="handleCancelEdit" :disabled="isSaving" class="px-3 py-1.5 bg-slate-800 hover:bg-slate-700 text-slate-300 rounded-lg text-xs font-semibold cursor-pointer">
                キャンセル
              </button>
              <button @click="handleSave" :disabled="isSaving" class="px-4 py-1.5 bg-blue-600 hover:bg-blue-500 text-white rounded-lg text-xs font-semibold flex items-center gap-1.5 cursor-pointer shadow">
                <span v-if="isSaving" class="animate-spin">⌛</span>
                <span v-else>💾</span>
                <span>保存する</span>
              </button>
            </div>
          </div>

          <!-- 通常表示時: 一言コメント (Bio) -->
          <div v-else-if="selectedDetail.account.description" class="pt-2 border-t border-slate-850">
            <div class="text-[11px] font-semibold text-slate-400 mb-1">💬 一言コメント (Bio):</div>
            <p class="text-xs text-slate-200 bg-slate-900/60 p-2.5 rounded-lg border border-slate-800/80 leading-relaxed whitespace-pre-line break-words">
              {{ selectedDetail.account.description }}
            </p>
          </div>

          <!-- クイックアクション導線 -->
          <div class="flex gap-2 pt-2 border-t border-slate-850">
            <button @click="emit('viewPosts', selectedDetail.account.numeric_id)" class="flex-1 py-1.5 px-3 bg-blue-600/80 hover:bg-blue-600 text-white rounded-lg text-xs font-semibold flex items-center justify-center gap-1 cursor-pointer">
              <span>📝</span> 投稿一覧を表示
            </button>
            <button @click="emit('viewMedia', selectedDetail.account.numeric_id)" class="flex-1 py-1.5 px-3 bg-indigo-600/80 hover:bg-indigo-600 text-white rounded-lg text-xs font-semibold flex items-center justify-center gap-1 cursor-pointer">
              <span>🖼️</span> メディア一覧を表示
            </button>
          </div>
        </div>

        <!-- ② 世代別アバター変更履歴 (ドラッグ＆ドロップで画像取り込み可能) -->
        <div class="space-y-2">
          <div class="flex justify-between items-center">
            <h5 class="text-xs font-bold text-slate-200 flex items-center gap-1.5">
              <span>🖼️</span> アバター世代履歴 (account_profile_histories)
            </h5>
            <span class="text-[10px] text-slate-400">💡 画像をカードに直接ドロップで登録</span>
          </div>

          <div v-if="!selectedDetail.histories || selectedDetail.histories.length === 0" class="p-5 text-center text-xs text-slate-500 bg-slate-950/60 rounded-xl border border-slate-800/60">
            世代履歴データはありません
          </div>
          <div v-else class="grid grid-cols-1 sm:grid-cols-2 gap-2.5">
            <div
              v-for="h in selectedDetail.histories"
              :key="h.id"
              @dragover.prevent
              @drop.prevent="handleHistoryDrop($event, h.avatar_virtual_key)"
              class="p-2.5 bg-slate-950 border border-slate-800 hover:border-blue-500/60 rounded-lg flex items-center gap-3 transition-colors relative group"
            >
              <!-- 共通 Avatar モジュールを使用 -->
              <label class="relative cursor-pointer flex-shrink-0" :title="h.avatar_virtual_key + ' に画像をドロップ/クリックで登録'">
                <Avatar
                  :avatar-url="resolveHistoryAvatarUrl(h)"
                  :handle="selectedDetail.account.username"
                  class="w-10 h-10 rounded-full border border-slate-700 group-hover:border-blue-400 transition-colors"
                />
                <div class="absolute inset-0 rounded-full bg-blue-600/0 group-hover:bg-blue-600/30 flex items-center justify-center text-[10px] text-white opacity-0 group-hover:opacity-100 transition-all font-bold">
                  📥
                </div>
                <input type="file" accept="image/*" class="hidden" @change="handleHistoryFileInput($event, h.avatar_virtual_key)" />
              </label>

              <div class="text-[11px] font-mono space-y-0.5 truncate flex-1 min-w-0">
                <div class="flex items-center justify-between gap-1">
                  <span class="text-emerald-400 font-bold">Seq: {{ h.avatar_seq }}</span>
                  <span class="text-[10px] text-slate-500">{{ new Date(h.observed_at).toLocaleDateString() }}</span>
                </div>
                <div class="text-slate-300 truncate text-[11px] flex items-center flex-wrap" :title="h.avatar_virtual_key">
                  <span>{{ h.avatar_virtual_key }}</span>
                  <span v-if="!isHistoryAvatarFound(h)" class="text-slate-500 text-[10px] ml-1.5 font-sans font-normal">
                    （見つかりませんでした）
                  </span>
                  <span v-else class="text-emerald-400/80 text-[10px] ml-1.5 font-sans font-normal">
                    （配置済）
                  </span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
      <div v-else class="h-full flex items-center justify-center text-xs text-slate-500">左側の一覧からアカウントを選択してください</div>
    </div>
  </div>
</template>
