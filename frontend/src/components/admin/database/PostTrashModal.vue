<!-- frontend/src/components/admin/database/PostTrashModal.vue (100行以下 - SPEC-PRINCIPLE-001) -->
<script setup lang="ts">
import { ref, watch } from 'vue';

const props = defineProps<{ show: boolean; articleId: string }>();
const emit = defineEmits<{ (e: 'close'): void; (e: 'confirm', reason: string): void }>();

const reason = ref('');

watch(() => props.show, (val) => {
  if (val) reason.value = '';
});

const handleConfirm = () => {
  const trimmed = reason.value.trim();
  if (!trimmed) return;
  emit('confirm', trimmed);
  emit('close');
};
</script>

<template>
  <div v-if="show" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/70 backdrop-blur-sm">
    <div class="bg-slate-900 border border-slate-700 rounded-2xl w-full max-w-md p-5 shadow-2xl space-y-4 text-slate-200">
      <!-- ヘッダー -->
      <div class="flex items-center justify-between border-b border-slate-800 pb-3">
        <div class="flex items-center gap-2">
          <span class="text-rose-400 text-lg">🗑️</span>
          <h3 class="text-sm font-bold text-slate-100">投稿をゴミ箱へ移動</h3>
        </div>
        <button @click="emit('close')" class="text-slate-400 hover:text-slate-200 text-lg leading-none cursor-pointer">✕</button>
      </div>

      <!-- 説明 & 対象ID -->
      <div class="space-y-1 text-xs">
        <p class="text-slate-400">対象記事をタイムライン・通常検索から除外し、ゴミ箱へ退避します。</p>
        <p class="text-slate-300 font-mono text-[11px] bg-slate-950 p-2 rounded border border-slate-800">
          記事 ID: <span class="text-blue-400 font-bold">{{ articleId }}</span>
        </p>
      </div>

      <!-- 削除理由入力欄（必須） -->
      <div class="space-y-1.5">
        <label class="text-[11px] font-bold text-rose-300 flex items-center justify-between">
          <span>削除理由の記入（必須）</span>
          <span class="text-[10px] text-slate-500 font-normal">※理由なしの削除は不可</span>
        </label>
        <textarea
          v-model="reason"
          rows="3"
          placeholder="理由を入力してください（例: 災害対応・デバッグデータ破棄、手動整理、破損など）..."
          class="w-full bg-slate-950 border border-slate-700 focus:border-rose-500 rounded-lg p-2.5 text-xs text-slate-100 placeholder:text-slate-600 focus:outline-none"
          @keydown.enter.ctrl="handleConfirm"
        ></textarea>
      </div>

      <!-- アクションボタン -->
      <div class="flex justify-end gap-2 pt-2 border-t border-slate-800">
        <button
          @click="emit('close')"
          class="px-3.5 py-1.5 bg-slate-800 hover:bg-slate-700 text-slate-300 text-xs font-bold rounded-lg cursor-pointer transition-colors"
        >
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
