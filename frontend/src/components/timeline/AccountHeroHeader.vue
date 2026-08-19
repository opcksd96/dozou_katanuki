<script setup lang="ts">
import { computed } from 'vue';
import type { RenderAuthor } from '../../models/RenderTree';

const props = defineProps<{
  account: RenderAuthor;
  totalArticles?: number;
}>();

const emit = defineEmits<{
  (e: 'backToAll'): void;
  (e: 'refresh'): void;
}>();
</script>

<template>
  <div class="mb-4 bg-slate-900/90 border border-slate-800 rounded-2xl overflow-hidden shadow-lg transition-all">
    <!-- カバーバナー背景 -->
    <div class="h-28 bg-gradient-to-r from-blue-900 via-indigo-900 to-slate-900 relative">
      <div class="absolute inset-0 bg-[radial-gradient(ellipse_at_top,_var(--tw-gradient-stops))] from-blue-500/20 via-transparent to-transparent"></div>
      
      <!-- 戻るボタン -->
      <button
        @click="emit('backToAll')"
        title="全アカウントのタイムラインに戻る"
        class="absolute top-3 left-3 px-3 py-1.5 bg-slate-950/70 hover:bg-slate-900/90 text-slate-200 hover:text-white rounded-full text-xs font-semibold backdrop-blur border border-slate-700/60 transition-all flex items-center gap-1.5 cursor-pointer shadow"
      >
        <span>←</span>
        <span>All Accounts</span>
      </button>

      <!-- プラットフォームバッジ -->
      <div class="absolute top-3 right-3 px-2.5 py-1 bg-slate-950/60 rounded-full text-[10px] font-mono text-slate-300 backdrop-blur border border-slate-800 flex items-center gap-1">
        <span class="w-1.5 h-1.5 rounded-full bg-emerald-400 animate-pulse"></span>
        <span>Local Archive</span>
      </div>
    </div>

    <!-- プロフィール詳細領域 -->
    <div class="px-5 pb-5 pt-0 relative">
      <div class="flex justify-between items-end -mt-12 mb-3">
        <!-- アバター -->
        <div class="relative">
          <img
            :src="account.avatar_url"
            :alt="account.handle"
            class="w-20 h-20 rounded-full object-cover border-4 border-slate-900 bg-slate-800 shadow-md"
            @error="($event.target as HTMLElement).style.display = 'none'"
          />
        </div>

        <!-- アクションボタン -->
        <div class="flex items-center gap-2">
          <button
            @click="emit('refresh')"
            title="このユーザーのアーカイブを再読み込み"
            class="p-2 bg-slate-800 hover:bg-slate-700 text-slate-300 hover:text-white border border-slate-700 rounded-lg text-xs transition-colors cursor-pointer"
          >
            🔄
          </button>
        </div>
      </div>

      <!-- 名前 & ハンドル -->
      <div class="space-y-1">
        <h2 class="text-lg font-bold text-white tracking-tight flex items-center gap-2">
          {{ account.display_name || account.handle }}
        </h2>
        <p class="text-xs text-slate-400 font-mono">@{{ account.handle }}</p>
      </div>

      <!-- 自己紹介文 (Bio) -->
      <p v-if="account.bio" class="text-xs text-slate-300 mt-2.5 leading-relaxed break-words whitespace-pre-line bg-slate-950/40 p-2.5 rounded-lg border border-slate-800/60">
        {{ account.bio }}
      </p>

      <!-- アーカイブステータスバー -->
      <div class="flex items-center gap-4 mt-3 pt-3 border-t border-slate-800/80 text-xs text-slate-400 font-mono">
        <div class="flex items-center gap-1.5">
          <span class="text-slate-500">ID:</span>
          <span class="text-slate-300 font-semibold">{{ account.numeric_id }}</span>
        </div>
        <div v-if="totalArticles !== undefined" class="flex items-center gap-1.5">
          <span class="text-slate-500">Posts:</span>
          <span class="text-blue-400 font-semibold">{{ totalArticles }} items</span>
        </div>
      </div>
    </div>
  </div>
</template>
