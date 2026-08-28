<!-- frontend/src/components/admin/ThunderConfirmModal.vue (100行以下 - SPEC-PRINCIPLE-001) -->
<script setup lang="ts">
import { ref } from 'vue';

const props = defineProps<{ isOpen: boolean; jobCount: number }>();
const emit = defineEmits<{ (e: 'close'): void; (e: 'confirm', intervalSec: number): void }>();

const confirmed = ref(false);
const intervalSec = ref(3);

const handleConfirm = () => {
  if (confirmed.value) {
    emit('confirm', intervalSec.value);
    emit('close');
  }
};
</script>

<template>
  <div v-if="isOpen" class="fixed inset-0 z-50 flex items-center justify-center bg-black/75 backdrop-blur-sm p-4 animate-fade-in font-sans">
    <div class="bg-slate-900 border border-amber-500/50 rounded-2xl max-w-lg w-full p-5 space-y-4 shadow-2xl text-slate-100">
      <div class="flex items-center gap-2 text-amber-400 border-b border-slate-800 pb-3">
        <span class="text-2xl">⚠️</span>
        <h3 class="text-base font-bold">迅雷エスカレーション投入の免責・注意事項</h3>
      </div>

      <div class="text-xs text-slate-300 space-y-2 leading-relaxed bg-slate-950/80 p-3 rounded-xl border border-slate-800 font-mono">
        <p class="text-amber-300 font-bold">【大量キュー投入に関する重要確認】</p>
        <p>1. これから対象メディア（約 {{ jobCount }} ジョブ）を迅雷アプリへ順次ディスパッチします。</p>
        <p>2. 迅雷側のリストが大量になりますので、事前に迅雷本体のリストをクリアするか整理しておくことを推奨します（自己責任）。</p>
        <p>3. 投入後30秒以内に <span class="text-purple-400 font-bold">*.xltd</span>（ダウンロードファイル）が生えたものは <span class="text-purple-400 font-bold">ESCALATED</span> で継続し、生えなかったものは <span class="text-amber-400 font-bold">RETAINED（長期保留）</span> へ自動移行します。</p>
      </div>

      <div class="flex items-center justify-between gap-3 text-xs bg-slate-950 p-2.5 rounded-xl border border-slate-800">
        <span class="text-slate-300 font-bold">投入間隔 (秒):</span>
        <select v-model="intervalSec" class="bg-slate-900 border border-slate-700 rounded px-2.5 py-1 text-slate-200 text-xs font-mono cursor-pointer">
          <option :value="2">2 秒おき (高速)</option>
          <option :value="3">3 秒おき (標準・推奨)</option>
          <option :value="5">5 秒おき (低速・安全)</option>
        </select>
      </div>

      <label class="flex items-center gap-2 text-xs text-slate-200 cursor-pointer pt-1">
        <input type="checkbox" v-model="confirmed" class="rounded bg-slate-800 border-slate-700 text-purple-600 focus:ring-0" />
        <span>上記免責事項と迅雷のリスト状況を確認しました（自己責任）</span>
      </label>

      <div class="flex justify-end gap-2 pt-2 border-t border-slate-800">
        <button @click="emit('close')" class="px-3 py-1.5 bg-slate-800 hover:bg-slate-700 text-slate-300 rounded-lg text-xs font-bold cursor-pointer">キャンセル</button>
        <button @click="handleConfirm" :disabled="!confirmed" class="px-4 py-1.5 bg-gradient-to-r from-purple-600 to-indigo-600 hover:from-purple-500 text-white rounded-lg text-xs font-bold shadow-lg disabled:opacity-40 cursor-pointer">
          🚀 承諾してディスパッチ開始
        </button>
      </div>
    </div>
  </div>
</template>
