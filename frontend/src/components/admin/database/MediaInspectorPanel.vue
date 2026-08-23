<script setup lang="ts">
import { ref, computed, watch } from 'vue';
import { BrowserOpenURL } from '../../../../wailsjs/runtime/runtime';
import Avatar from '../../article/Avatar.vue';
import { useAdminDatabase } from '../../../composables/admin/useAdminDatabase';

const props = defineProps<{
  media: any;
}>();

const emit = defineEmits<{
  (e: 'saveMetadata', payload: { mediaId: string; downloadStatus: string; stashSceneId: string; stashImageId: string; failedReason: string }): void;
  (e: 'retry', mediaId: string): void;
  (e: 'purge', mediaId: string): void;
  (e: 'viewPost', articleId: string): void;
}>();

const { fetchStashMetadata, updateStashMetadata } = useAdminDatabase();

// SQLite メタデータ編集用状態
const editStatus = ref(props.media.download_status || props.media.raw_status || 'QUEUED');
const editReason = ref(props.media.failed_reason || '');

// Stash GraphQL メタデータ
const stashData = ref<any>(null);
const isLoadingStash = ref(false);
const editTitle = ref('');
const editDetails = ref('');
const editRating = ref(0);

// 安全なミューテーション管理（確認モーダル・Undoスナップショット）
const showConfirmModal = ref(false);
const isMutating = ref(false);
const undoSnapshot = ref<{ title: string; details: string; rating100: number } | null>(null);
const copiedField = ref<string | null>(null);

const currentHostname = computed(() => (typeof window !== 'undefined' && window.location?.hostname) ? window.location.hostname : '127.0.0.1');
const stashDirectUrl = computed(() => {
  if (props.media.stash_scene_id) return `http://${currentHostname.value}:9999/scenes/${props.media.stash_scene_id}`;
  if (props.media.stash_image_id) return `http://${currentHostname.value}:9999/images/${props.media.stash_image_id}`;
  return null;
});

const isStashModified = computed(() => {
  if (!stashData.value) return false;
  return editTitle.value !== (stashData.value.title || '') ||
         editDetails.value !== (stashData.value.details || '') ||
         editRating.value !== (stashData.value.rating100 || 0);
});

const loadStashInfo = async () => {
  if (!props.media.stash_scene_id && !props.media.stash_image_id) {
    stashData.value = null;
    return;
  }
  isLoadingStash.value = true;
  const res = await fetchStashMetadata(props.media.stash_scene_id || '', props.media.stash_image_id || '');
  if (res) {
    stashData.value = res;
    editTitle.value = res.title || '';
    editDetails.value = res.details || '';
    editRating.value = res.rating100 || 0;
  }
  isLoadingStash.value = false;
};

watch(() => props.media, (m) => {
  if (m) {
    editStatus.value = m.download_status || m.raw_status || 'QUEUED';
    editReason.value = m.failed_reason || '';
    undoSnapshot.value = null;
    loadStashInfo();
  }
}, { immediate: true });

const copyToClipboard = (text: string, field: string) => {
  if (navigator.clipboard) {
    navigator.clipboard.writeText(text);
    copiedField.value = field;
    setTimeout(() => { copiedField.value = null; }, 2000);
  }
};

const openStash = () => {
  if (!stashDirectUrl.value) return;
  try { BrowserOpenURL(stashDirectUrl.value); } catch { window.open(stashDirectUrl.value, '_blank', 'noopener,noreferrer'); }
};

const handleSaveSqlite = () => {
  emit('saveMetadata', {
    mediaId: props.media.media_id || props.media.id,
    downloadStatus: editStatus.value,
    stashSceneId: props.media.stash_scene_id || '',
    stashImageId: props.media.stash_image_id || '',
    failedReason: editReason.value.trim(),
  });
};

const triggerStashMutation = async () => {
  const isScene = !!props.media.stash_scene_id;
  const targetId = props.media.stash_scene_id || props.media.stash_image_id;
  if (!targetId) return;

  isMutating.value = true;
  // Undo用スナップショットを退避
  if (stashData.value) {
    undoSnapshot.value = {
      title: stashData.value.title || '',
      details: stashData.value.details || '',
      rating100: stashData.value.rating100 || 0,
    };
  }

  const updated = await updateStashMetadata(isScene, targetId, editTitle.value.trim(), editDetails.value.trim(), editRating.value);
  if (updated) {
    stashData.value = updated;
    editTitle.value = updated.title || '';
    editDetails.value = updated.details || '';
    editRating.value = updated.rating100 || 0;
  }
  isMutating.value = false;
  showConfirmModal.value = false;
};

const handleUndoStash = async () => {
  if (!undoSnapshot.value) return;
  const isScene = !!props.media.stash_scene_id;
  const targetId = props.media.stash_scene_id || props.media.stash_image_id;
  if (!targetId) return;

  isMutating.value = true;
  const snap = undoSnapshot.value;
  const restored = await updateStashMetadata(isScene, targetId, snap.title, snap.details, snap.rating100);
  if (restored) {
    stashData.value = restored;
    editTitle.value = restored.title || '';
    editDetails.value = restored.details || '';
    editRating.value = restored.rating100 || 0;
    undoSnapshot.value = null;
  }
  isMutating.value = false;
};
</script>

<template>
  <div class="w-80 md:w-[420px] bg-slate-900/95 border-l border-slate-800 flex flex-col p-4 space-y-3.5 overflow-y-auto text-xs font-mono text-slate-300 shrink-0">
    <!-- 1. アカウント情報 -->
    <div class="p-3 bg-slate-950/80 rounded-xl border border-slate-800 space-y-2">
      <div class="text-[10px] text-slate-400 font-bold uppercase tracking-wider">👤 投稿者アカウント</div>
      <div class="flex items-center gap-2.5">
        <Avatar :avatar-url="media.avatar_url" :handle="media.username" size-class="w-10 h-10" />
        <div class="min-w-0 flex-1 leading-tight">
          <div class="text-slate-100 font-bold truncate text-[13px]">{{ media.display_name || media.username || 'Unknown' }}</div>
          <div class="text-slate-400 text-[11px] truncate">@{{ media.username }} <span class="text-slate-500">({{ media.account_id || '-' }})</span></div>
        </div>
      </div>
      <div v-if="media.created_at || media.tweet_date" class="text-[10px] text-slate-500 pt-1 border-t border-slate-850 flex items-center justify-between">
        <span>投稿日: {{ media.tweet_date || (media.created_at ? new Date(media.created_at).toLocaleString() : '-') }}</span>
      </div>
    </div>

    <!-- 1.5. 親ツイート本文 (Parent Tweet Details) -->
    <div v-if="media.full_text || media.full_text_ja" class="p-3 bg-slate-950/80 rounded-xl border border-slate-800 space-y-2">
      <div class="flex items-center justify-between text-[10px]">
        <span class="text-slate-400 font-bold uppercase tracking-wider">📝 親ツイート本文</span>
        <button @click="copyToClipboard(media.full_text_ja || media.full_text, 'tweet')" class="text-blue-400 hover:text-blue-300">
          {{ copiedField === 'tweet' ? '✓ コピー済' : '📋 コピー' }}
        </button>
      </div>
      <div class="p-2.5 bg-slate-900 rounded-lg border border-slate-800 space-y-1.5 select-text">
        <div v-if="media.full_text_ja" class="text-xs text-slate-100 leading-relaxed font-sans">
          {{ media.full_text_ja }}
        </div>
        <div v-if="media.full_text && media.full_text !== media.full_text_ja" class="text-[11px] text-slate-400 leading-relaxed font-sans border-t border-slate-800 pt-1">
          {{ media.full_text }}
        </div>
      </div>
      <button 
        v-if="stashDirectUrl"
        @click="editDetails = (media.full_text_ja ? media.full_text_ja + '\n\n' + media.full_text : media.full_text)"
        class="w-full py-1 bg-slate-800 hover:bg-slate-700 text-purple-300 border border-purple-700/40 rounded text-[10px] font-mono flex items-center justify-center gap-1 transition-colors"
      >
        <span>⬇️</span> 本文を Stash 詳細メモへ転記
      </button>
    </div>

    <!-- 2. Stash 連携・GraphQL メタデータインスペクタ -->
    <div class="p-3 bg-slate-950/80 rounded-xl border border-slate-800 space-y-2.5">
      <div class="flex items-center justify-between text-[10px]">
        <span class="text-slate-400 font-bold uppercase tracking-wider">📦 Stash 連携情報 (GraphQL)</span>
        <span :class="stashDirectUrl ? 'text-emerald-400 font-bold' : 'text-slate-500'">{{ stashDirectUrl ? '● 連携中' : '○ 未連携' }}</span>
      </div>

      <!-- ID表示（読み取り専用・変更不可・クリックでコピー） -->
      <div class="p-2 bg-slate-900 rounded-lg border border-slate-800 space-y-1 text-[11px]">
        <div class="flex items-center justify-between">
          <span class="text-slate-500">Stash {{ media.stash_scene_id ? 'Scene' : 'Image' }} ID (不可逆):</span>
          <button @click="copyToClipboard(media.stash_scene_id || media.stash_image_id || '', 'id')" class="text-[10px] text-blue-400 hover:text-blue-300">
            {{ copiedField === 'id' ? '✓ コピー済' : 'コピー' }}
          </button>
        </div>
        <div class="font-bold text-purple-300 truncate select-all bg-black/40 px-2 py-1 rounded border border-purple-900/50">
          {{ media.stash_scene_id || media.stash_image_id || '未連携' }}
        </div>
      </div>

      <!-- Stash ファイル・詳細メタデータ -->
      <div v-if="stashData" class="space-y-2 text-xs text-slate-300 bg-slate-900/80 p-3 rounded-lg border border-slate-800">
        <!-- Stash 側の詳細テキスト (親ツイート本文) -->
        <div v-if="stashData.details" class="space-y-1">
          <div class="text-[10px] text-purple-400 font-bold uppercase flex items-center justify-between">
            <span>📝 Stash 詳細テキスト (Parent Tweet)</span>
            <button @click="copyToClipboard(stashData.details, 'stash_details')" class="text-blue-400 hover:text-blue-300">
              {{ copiedField === 'stash_details' ? '✓ コピー済' : '📋 コピー' }}
            </button>
          </div>
          <div class="p-2 bg-slate-950/80 rounded border border-slate-800/80 text-slate-200 text-xs leading-relaxed whitespace-pre-wrap font-sans select-text">
            {{ stashData.details }}
          </div>
        </div>

        <div class="grid grid-cols-2 gap-2 text-[11px] text-slate-400 pt-1 border-t border-slate-800/60">
          <div v-if="stashData.date"><span class="text-slate-500">日付:</span> <span class="text-slate-200 font-semibold">{{ stashData.date }}</span></div>
          <div v-if="stashData.studio"><span class="text-slate-500">スタジオ:</span> <span class="text-slate-200 font-semibold">{{ stashData.studio }}</span></div>
          <div v-if="stashData.files?.length"><span class="text-slate-500">解像度:</span> <span class="text-slate-200 font-bold">{{ stashData.files[0].width }}x{{ stashData.files[0].height }}</span></div>
          <div v-if="stashData.files?.[0]?.duration"><span class="text-slate-500">再生時間:</span> <span class="text-slate-200">{{ Math.round(stashData.files[0].duration) }}秒</span></div>
        </div>
        <div v-if="stashData.files?.[0]?.path" class="truncate text-[10px] text-slate-500 font-mono" :title="stashData.files[0].path">
          <span>パス:</span> <span class="text-slate-400">{{ stashData.files[0].path }}</span>
        </div>
      </div>

      <button v-if="stashDirectUrl" @click="openStash" class="w-full py-1.5 bg-purple-950 hover:bg-purple-900 text-purple-200 border border-purple-700/60 font-bold rounded-lg transition-colors flex items-center justify-center gap-1.5 text-xs">
        🎛️ Stash WebUI で開く ↗
      </button>
    </div>

    <!-- 3. Stash GraphQL メタデータ安全ミューテーション -->
    <div v-if="stashDirectUrl" class="p-3 bg-slate-950/80 rounded-xl border border-slate-800 space-y-2.5">
      <div class="flex items-center justify-between text-[10px]">
        <span class="text-slate-400 font-bold uppercase tracking-wider">✏️ Stash メタデータ更新</span>
        <span v-if="isStashModified" class="text-amber-400 font-bold animate-pulse">● 未保存の変更あり</span>
      </div>

      <div class="space-y-1">
        <label class="text-[10px] text-slate-400">タイトル (Title)</label>
        <input v-model="editTitle" type="text" placeholder="タイトルを入力..." class="w-full bg-slate-900 border border-slate-700 rounded px-2 py-1 text-slate-200 text-[11px]" />
      </div>

      <div class="space-y-1">
        <label class="text-[10px] text-slate-400">詳細・メモ (Details)</label>
        <textarea v-model="editDetails" rows="2" placeholder="詳細メモを入力..." class="w-full bg-slate-900 border border-slate-700 rounded px-2 py-1 text-slate-200 text-[11px] resize-none"></textarea>
      </div>

      <div class="space-y-1">
        <label class="text-[10px] text-slate-400">評価 (Rating: {{ editRating ? editRating + '/100' : '未設定' }})</label>
        <div class="flex items-center gap-1">
          <button v-for="r in [20, 40, 60, 80, 100]" :key="r" @click="editRating = (editRating === r ? 0 : r)" type="button" class="px-2 py-0.5 rounded text-[10px] border transition-colors" :class="editRating >= r ? 'bg-amber-500/20 text-amber-300 border-amber-500/60 font-bold' : 'bg-slate-900 text-slate-500 border-slate-800'">
            ★ {{ r / 20 }}
          </button>
        </div>
      </div>

      <div class="flex gap-2 pt-1">
        <button @click="showConfirmModal = true" :disabled="!isStashModified || isMutating" class="flex-1 py-1.5 bg-purple-600 hover:bg-purple-500 disabled:bg-slate-800 disabled:text-slate-600 text-white font-bold rounded-lg transition-colors text-[11px] shadow flex items-center justify-center gap-1">
          <span>💾</span> Stash 更新 (確認へ)
        </button>
        <button v-if="undoSnapshot" @click="handleUndoStash" :disabled="isMutating" class="px-3 py-1.5 bg-amber-950/80 hover:bg-amber-900 text-amber-200 border border-amber-700/60 font-bold rounded-lg transition-colors text-[11px] flex items-center gap-1" title="直前の状態にロールバック">
          <span>↩️</span> Undo
        </button>
      </div>
    </div>

    <!-- 4. SQLite ダウンロードステータス管理 -->
    <div class="p-3 bg-slate-950/80 rounded-xl border border-slate-800 space-y-2">
      <div class="text-[10px] text-slate-400 font-bold uppercase tracking-wider">🗄️ アーカイブ DB 状態管理</div>
      <div class="grid grid-cols-2 gap-2">
        <div>
          <label class="text-[9px] text-slate-500 block mb-0.5">ダウンロード状態</label>
          <select v-model="editStatus" class="w-full bg-slate-900 border border-slate-700 rounded px-2 py-1 text-slate-200 text-[10px]">
            <option value="COMPLETED">COMPLETED</option>
            <option value="QUEUED">QUEUED</option>
            <option value="EXCLUDED">EXCLUDED</option>
            <option value="DEAD_404">DEAD_404</option>
          </select>
        </div>
        <div>
          <label class="text-[9px] text-slate-500 block mb-0.5">失敗/保留理由</label>
          <input v-model="editReason" type="text" placeholder="理由..." class="w-full bg-slate-900 border border-slate-700 rounded px-2 py-1 text-slate-200 text-[10px]" />
        </div>
      </div>
      <button @click="handleSaveSqlite" class="w-full py-1 bg-slate-800 hover:bg-slate-700 text-slate-200 font-bold rounded text-[10px] transition-colors border border-slate-700">
        DB ステータスを更新
      </button>
    </div>

    <!-- 5. アクションボタン -->
    <div class="pt-2 border-t border-slate-800 space-y-1.5">
      <button v-if="media.article_id" @click="emit('viewPost', media.article_id)" class="w-full py-1.5 bg-slate-800 hover:bg-slate-700 text-blue-300 font-bold rounded-lg transition-colors text-[11px] flex items-center justify-center gap-1">
        📝 親記事を見る
      </button>
      <div class="grid grid-cols-2 gap-2">
        <button @click="emit('retry', media.media_id || media.id)" class="py-1.5 bg-blue-600 hover:bg-blue-500 text-white font-bold rounded-lg transition-colors text-[11px]">
          🔄 再取得
        </button>
        <button @click="emit('purge', media.media_id || media.id)" class="py-1.5 bg-rose-950 hover:bg-rose-800 text-rose-200 font-bold rounded-lg border border-rose-700/60 transition-colors text-[11px]">
          🗑️ パージ
        </button>
      </div>
    </div>

    <!-- 変更確認モーダル (Safe Mutation Modal) -->
    <div v-if="showConfirmModal" class="fixed inset-0 z-50 bg-black/80 backdrop-blur-sm flex items-center justify-center p-4" @click.self="showConfirmModal = false">
      <div class="bg-slate-900 border border-slate-700 rounded-2xl max-w-md w-full p-5 space-y-4 shadow-2xl font-mono text-xs">
        <div class="flex items-center gap-2 text-purple-400 font-bold text-sm">
          <span>⚠️</span> Stash メタデータ更新の確認
        </div>
        <p class="text-slate-300 text-[11px] leading-relaxed">
          Stash GraphQL API を通じてメタデータを書き換えます。よろしいですか？（実行後も Undo で直前の状態に戻せます）
        </p>

        <!-- 差分プレビュー -->
        <div class="p-3 bg-slate-950 rounded-lg border border-slate-800 space-y-2 text-[11px]">
          <div>
            <span class="text-slate-500 block">タイトル:</span>
            <div class="flex items-center gap-2">
              <span class="text-slate-400 line-through truncate max-w-[140px]">{{ stashData?.title || '(未設定)' }}</span>
              <span>➔</span>
              <span class="text-emerald-400 font-bold truncate max-w-[140px]">{{ editTitle || '(未設定)' }}</span>
            </div>
          </div>
          <div>
            <span class="text-slate-500 block">詳細メモ:</span>
            <div class="text-emerald-400 font-bold truncate">{{ editDetails || '(未設定)' }}</div>
          </div>
          <div>
            <span class="text-slate-500 block">評価:</span>
            <span class="text-amber-400 font-bold">{{ editRating ? editRating + ' / 100' : '未設定' }}</span>
          </div>
        </div>

        <div class="flex justify-end gap-2 pt-2">
          <button @click="showConfirmModal = false" class="px-3 py-1.5 bg-slate-800 hover:bg-slate-700 text-slate-300 rounded-lg">
            キャンセル
          </button>
          <button @click="triggerStashMutation" :disabled="isMutating" class="px-4 py-1.5 bg-purple-600 hover:bg-purple-500 text-white font-bold rounded-lg shadow">
            {{ isMutating ? '更新中...' : '実行する' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
