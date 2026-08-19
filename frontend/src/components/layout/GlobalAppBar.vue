<script setup lang="ts">
import { ref, onMounted } from 'vue';

defineProps<{
  activeArticleHandle?: string;
  activeArticleId?: string | null;
}>();

const emit = defineEmits<{
  (e: 'openAdmin'): void;
  (e: 'reloadAll'): void;
  (e: 'backToTimeline'): void;
}>();

const isStashOnline = ref(true);

const checkStash = async () => {
  try {
    const res = await fetch('/stash-proxy/', { method: 'HEAD' });
    isStashOnline.value = res.ok || res.status === 401 || res.status === 404;
  } catch {
    isStashOnline.value = false;
  }
};

onMounted(() => {
  checkStash();
  setInterval(checkStash, 15000);
});
</script>

<template>
  <nav class="w-full bg-slate-950/95 backdrop-blur border-b border-slate-800 sticky top-0 z-40 px-4 py-2.5">
    <div class="max-w-4xl mx-auto flex items-center justify-between gap-3">
      <!-- 左側: ブランドロゴ & パンくずリスト -->
      <div class="flex items-center gap-3">
        <div
          @click="emit('backToTimeline')"
          class="flex items-center gap-2 cursor-pointer group"
          title="トップ・タイムラインに戻る"
        >
          <div class="w-7 h-7 rounded-lg bg-blue-600 flex items-center justify-center text-white font-bold text-sm shadow-md shadow-blue-600/30 group-hover:bg-blue-500 transition-colors">
            蔵
          </div>
          <div>
            <div class="flex items-center gap-1.5">
              <span class="text-sm font-bold text-white tracking-tight group-hover:text-blue-400 transition-colors">dozou_katanuki</span>
              <span class="text-[10px] px-1.5 py-0.5 rounded bg-slate-800 text-slate-400 font-mono">v4.0</span>
            </div>
          </div>
        </div>

        <!-- 個別詳細時のパンくずナビゲーション -->
        <div v-if="activeArticleId" class="hidden sm:flex items-center gap-1.5 text-xs text-slate-400 font-mono pl-3 border-l border-slate-800">
          <button @click="emit('backToTimeline')" class="hover:text-blue-400 transition-colors cursor-pointer">タイムライン</button>
          <span>/</span>
          <span class="text-slate-300">@{{ activeArticleHandle }}</span>
          <span>/</span>
          <span class="text-blue-400 font-semibold">ポスト</span>
        </div>
      </div>

      <!-- 右側: グローバルコントロール -->
      <div class="flex items-center gap-2">
        <!-- Stash 稼働インジケーター -->
        <div
          :title="isStashOnline ? 'Stash サーバー接続完了 (ポート:9999)' : 'Stash 未接続 / 待機中'"
          class="hidden sm:flex items-center gap-1.5 px-2.5 py-1 rounded-md bg-slate-900 border border-slate-800 text-[11px] font-mono"
        >
          <span :class="['w-2 h-2 rounded-full', isStashOnline ? 'bg-emerald-400 animate-pulse' : 'bg-amber-500']"></span>
          <span :class="isStashOnline ? 'text-slate-300' : 'text-slate-400'">Stash</span>
        </div>

        <!-- 再読み込みボタン -->
        <button
          @click="emit('reloadAll')"
          title="最新のアーカイブデータを再取得 (Ctrl+R / F5)"
          class="p-1.5 rounded-lg bg-slate-900 hover:bg-slate-800 text-slate-300 hover:text-white border border-slate-800 text-xs transition-all cursor-pointer flex items-center gap-1"
        >
          <span>🔄</span>
          <span class="hidden md:inline text-[11px] font-medium">更新</span>
        </button>

        <!-- 設定・Admin Board ボタン -->
        <button
          @click="emit('openAdmin')"
          title="管理ダッシュボード ＆ 設定を開く (Ctrl+,)"
          class="flex items-center gap-1.5 px-3 py-1.5 bg-blue-600 hover:bg-blue-500 text-white rounded-lg text-xs font-semibold transition-all shadow-md shadow-blue-600/20 cursor-pointer"
        >
          <span>⚙️</span>
          <span>設定・管理</span>
        </button>
      </div>
    </div>
  </nav>
</template>
