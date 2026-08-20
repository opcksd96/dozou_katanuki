<script setup lang="ts">
defineProps<{
  isOpen: boolean;
}>();

const emit = defineEmits<{
  (e: 'close'): void;
}>();

const shortcuts = [
  { key: 'J / K', desc: '次のポスト / 前のポストに移動' },
  { key: 'Enter / O', desc: '選択中のポストの詳細 / 会話ツリーを開く' },
  { key: 'L', desc: '選択中のポストのいいね（ブックマーク）をトグル' },
  { key: 'M / X', desc: '選択中のポストの添付メディアを展開' },
  { key: '← / →', desc: 'メディアビューアで前後スライド切り替え' },
  { key: 'Esc', desc: '詳細ビュー / メディア / モーダルを閉じる' },
  { key: 'Ctrl + ,', desc: '管理者設定ボードを開く' },
  { key: '?', desc: 'このキーボードショートカット一覧を表示' },
];
</script>

<template>
  <div
    v-if="isOpen"
    @click.self="emit('close')"
    class="fixed inset-0 z-50 bg-black/70 backdrop-blur-sm flex items-center justify-center p-4 animate-in fade-in duration-150"
  >
    <div class="bg-slate-900 border border-slate-700/80 rounded-2xl w-full max-w-md p-6 shadow-2xl space-y-4 text-left">
      <div class="flex items-center justify-between border-b border-slate-800 pb-3">
        <h3 class="text-base font-bold text-white flex items-center gap-2">
          <span>⌨️</span> キーボードショートカット
        </h3>
        <button
          @click="emit('close')"
          class="p-1 rounded-full hover:bg-slate-800 text-slate-400 hover:text-white transition-colors cursor-pointer"
        >
          ✕
        </button>
      </div>

      <div class="space-y-2.5">
        <div
          v-for="s in shortcuts"
          :key="s.key"
          class="flex items-center justify-between text-xs py-1 px-2 rounded hover:bg-slate-800/40 transition-colors"
        >
          <span class="text-slate-300">{{ s.desc }}</span>
          <span class="twitter-shortcut-key">{{ s.key }}</span>
        </div>
      </div>

      <div class="pt-2 border-t border-slate-800 text-center">
        <button
          @click="emit('close')"
          class="px-4 py-1.5 rounded-lg bg-blue-600 hover:bg-blue-500 text-white text-xs font-semibold transition-colors cursor-pointer"
        >
          閉じる (Esc)
        </button>
      </div>
    </div>
  </div>
</template>
