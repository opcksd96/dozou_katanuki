<!-- frontend/src/components/admin/database/MediaTrashModal.vue (100行以下 - SPEC-PRINCIPLE-001) -->
<script setup lang="ts">
import { ref, watch } from 'vue';

const props = defineProps<{ show: boolean; mediaItem: any }>();
const emit = defineEmits<{ (e: 'close'): void; (e: 'confirm', payload: { mediaId: string; reason: string }): void }>();

const reason = ref('');

watch(() => props.show, (val) => {
  if (val) reason.value = '';
});

const handleConfirm = () => {
  const trimmed = reason.value.trim();
  const mId = props.mediaItem?.media_id || props.mediaItem?.id || props.mediaItem?.mediaId;
  if (!trimmed || !mId) return;
  emit('confirm', { mediaId: mId, reason: trimmed });
  emit('close');
};
</script>

<template>
  <div v-if="show" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/70 backdrop-blur-sm font-sans">
    <div class="bg-slate-900 border border-slate-700 rounded-2xl w-full max-w-md p-5 shadow-2xl space-y-4 text-slate-200">
      <!-- ヘッダー -->
      <div class="flex items-center justify-between border-b border-slate-800 pb-3">
        <div class="flex items-center gap-2">
          <span class="text-rose-400 text-lg">🗑️</span>
          <h3 class="text-sm font-bold text-slate-100">メディアをゴミ箱へ移動</h3>
        </div>
        <button @click="emit('close')" class="text-slate-400 hover:text-slate-200 text-lg leading-none cursor-pointer">✕</button>
      </div>

      <!-- 説明 & 対象メディア情報 -->
      <div class="space-y-1 text-xs">
        <p class="text-slate-400">対象メディアを通常タイムライン・一覧から除外し、ゴミ箱へ退避します。</p>
        <div class="bg-slate-950 p-2.5 rounded-lg border border-slate-800 space-y-0.5 font-mono text-[11px]">
          <p class="text-slate-200 font-bold">メディア ID: <span class="text-blue-400">{{ mediaItem?.media_id || mediaItem?.id || '' }}</span></p>
          <p class="text-slate-500 text-[10px]">種別: {{ mediaItem?.type || 'unknown' }} | 状態: {{ mediaItem?.download_status || mediaItem?.raw_status || 'QUEUED' }}</p>
        </div>
      </div>

      <!-- 削除理由入力欄（必須） -->
      <div class="space-y-1.5">
        <label class="text-[11px] font-bold text-rose-300 flex items-center justify-between">
          <span>退避・削除理由の記入（必須）</span>
          <span class="text-[10px] text-slate-500 font-normal">※理由なしの退避は不可</span>
        </label>
        <textarea
          v-model="reason"
          rows="3"
          placeholder="理由を入力してください（例: 破損メディア、重複ダウンロード、不要なファイルなど）..."
          class="w-full bg-slate-950 border border-slate-700 focus:border-rose-500 rounded-lg p-2.5 text-xs text-slate-100 placeholder:text-slate-600 focus:outline-none"
          @keydown.enter.ctrl="handleConfirm"
        ></textarea>
      </div>

      <!-- アクションボタン -->
      <div class="flex justify-end gap-2 pt-2 border-t border-slate-800">
        <button @click="emit('close')" class="px-3.5 py-1.5 bg-slate-800 hover:bg-slate-700 text-slate-300 text-xs font-bold rounded-lg cursor-pointer transition-colors">
          キャンセル
        </button>
        <button
          @click="handleConfirm"
          :disabled="!reason.trim()"
          class="px-4 py-1.5 bg-rose-600 hover:bg-rose-500 disabled:opacity-40 disabled:cursor-not-allowed text-white text-xs font-bold rounded-lg cursor-pointer transition-colors shadow flex items-center gap-1"
        >
          <span>🗑️</span> ゴミ箱へ移動
        </button>
      </div>
    </div>
  </div>
</template>
