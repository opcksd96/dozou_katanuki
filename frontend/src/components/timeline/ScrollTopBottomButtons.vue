<!-- frontend/src/components/timeline/ScrollTopBottomButtons.vue (100行以下 - SPEC-PRINCIPLE-001) -->
<script setup lang="ts">
import { ArrowUp, ArrowDown } from 'lucide-vue-next';

const scrollToTop = () => {
  window.scrollTo({ top: 0, behavior: 'smooth' });
};
const scrollToBottom = () => {
  window.scrollTo({ top: document.body.scrollHeight || document.documentElement.scrollHeight, behavior: 'smooth' });
};
</script>

<template>
  <div class="fixed bottom-6 right-5 flex flex-col gap-2.5 z-40 select-none touch-manipulation">
    <button
      @click="scrollToTop"
      class="w-11 h-11 bg-slate-900/90 hover:bg-blue-600 active:bg-blue-500 text-slate-200 hover:text-white rounded-full border border-slate-700/80 shadow-xl flex items-center justify-center transition-all active:scale-90 cursor-pointer backdrop-blur-md"
      title="ページ最上部へスクロール"
      aria-label="Scroll to top"
    >
      <ArrowUp class="w-5 h-5" />
    </button>
    <button
      @click="scrollToBottom"
      class="w-11 h-11 bg-slate-900/90 hover:bg-blue-600 active:bg-blue-500 text-slate-200 hover:text-white rounded-full border border-slate-700/80 shadow-xl flex items-center justify-center transition-all active:scale-90 cursor-pointer backdrop-blur-md"
      title="ページ最下部へスクロール"
      aria-label="Scroll to bottom"
    >
      <ArrowDown class="w-5 h-5" />
    </button>
  </div>
</template>

