<template>
  <div class="space-y-4 p-4 text-neutral-800 dark:text-neutral-200">
    <!-- 1. 対象ツイートコンテキスト -->
    <div class="p-3 rounded-lg border border-neutral-200 dark:border-neutral-800 bg-neutral-50 dark:bg-neutral-900/60">
      <div class="flex justify-between items-center mb-1">
        <span class="text-xs font-bold text-sky-500 uppercase">🎯 Target Tweet</span>
        <span class="text-xs font-mono text-neutral-400">ID: {{ targetTweetId || 'None' }}</span>
      </div>
      <p class="text-xs text-neutral-600 dark:text-neutral-300 line-clamp-2">
        {{ targetTweetText || 'タイムラインまたは個別記事からツイートを選択してください。' }}
      </p>
    </div>

    <!-- 2. アカウント・エンジン選択 -->
    <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
      <div class="p-3 rounded-lg border border-neutral-200 dark:border-neutral-800 space-y-2">
        <div class="text-xs font-bold text-neutral-500">🧭 検出メンション/リプライ先</div>
        <div class="flex flex-wrap gap-1.5 max-h-24 overflow-y-auto">
          <button
            v-for="h in detectedHandles"
            :key="h"
            @click="toggleHandle(h)"
            :class="selectedHandles.includes(h) ? 'bg-sky-500 text-white' : 'bg-neutral-200 dark:bg-neutral-800 text-neutral-500'"
            class="px-2 py-0.5 rounded text-xs font-mono transition"
          >
            @{{ h }}
          </button>
          <span v-if="!detectedHandles.length" class="text-xs text-neutral-400">未検出</span>
        </div>
      </div>

      <div class="p-3 rounded-lg border border-neutral-200 dark:border-neutral-800 space-y-2.5">
        <div class="text-xs font-bold text-neutral-500">⚙️ 探索設定</div>
        <div class="flex items-center gap-2">
          <select v-model="selectedEngine" class="px-2 py-1 rounded border border-neutral-300 dark:border-neutral-700 bg-white dark:bg-neutral-800 text-xs">
            <option value="sotwe">⚡ Sotwe (JSON API)</option>
            <option value="nitter" disabled>🕸️ Nitter (HTML)</option>
          </select>
          <label class="flex items-center gap-1 text-xs text-neutral-500 cursor-pointer">
            <input type="checkbox" v-model="enableWayback" class="rounded text-sky-500" />
            魚拓巻き込み
          </label>
        </div>
        <button
          @click="startExploration"
          :disabled="isExploring || !selectedHandles.length"
          class="w-full py-1.5 bg-sky-500 hover:bg-sky-600 disabled:opacity-50 text-white font-bold text-xs rounded transition"
        >
          {{ isExploring ? '⏳ 探索中...' : '🚀 Search & Salvage Now' }}
        </button>
      </div>
    </div>

    <!-- 3. 疑似ターミナルログ -->
    <div class="p-3 rounded-lg border border-neutral-800 bg-neutral-950 font-mono text-xs text-neutral-300 space-y-1">
      <div class="text-neutral-500 border-b border-neutral-800 pb-1 flex justify-between">
        <span>🖥️ Live Console</span>
        <span>{{ logLines.length }} lines</span>
      </div>
      <div class="h-36 overflow-y-auto space-y-0.5">
        <div v-for="(l, i) in logLines" :key="i"><span class="text-emerald-400">&gt;</span> {{ l }}</div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useAdminRelationExplorer } from '../../composables/admin/useAdminRelationExplorer';

const {
  targetTweetId,
  targetTweetText,
  detectedHandles,
  selectedHandles,
  selectedEngine,
  enableWayback,
  isExploring,
  logLines,
  toggleHandle,
  startExploration
} = useAdminRelationExplorer();
</script>
