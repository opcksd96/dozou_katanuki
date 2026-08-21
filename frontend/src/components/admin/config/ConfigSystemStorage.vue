<!-- frontend/src/components/admin/config/ConfigSystemStorage.vue (100行以下) -->
<script setup lang="ts">
import { models } from '../../../../wailsjs/go/models';
defineProps<{ config: models.AppConfig }>();
</script>

<template>
  <div class="space-y-4">
    <!-- 1. System 設定 -->
    <div class="bg-slate-900/80 border border-slate-800 rounded-xl p-4 space-y-4">
      <h3 class="text-sm font-bold text-slate-200 flex items-center gap-2"><span>🌐</span> システム基本設定 (System)</h3>
      <div class="grid grid-cols-1 sm:grid-cols-3 gap-4">
        <div>
          <label class="block text-xs text-slate-400 mb-1">システム言語</label>
          <select v-model="config.system.language" class="w-full bg-slate-950 border border-slate-700 rounded-lg px-3 py-1.5 text-xs text-slate-200">
            <option value="ja">日本語 (ja)</option><option value="en">English (en)</option><option value="zh">简体中文 (zh)</option>
          </select>
        </div>
        <div>
          <label class="block text-xs text-slate-400 mb-1">デフォルトプラットフォーム</label>
          <select v-model="config.system.default_framework" class="w-full bg-slate-950 border border-slate-700 rounded-lg px-3 py-1.5 text-xs text-slate-200">
            <option value="twitter">Twitter / X</option><option value="bsky">Bluesky</option>
          </select>
        </div>
        <div>
          <label class="block text-xs text-slate-400 mb-1">環境モード</label>
          <input v-model="config.system.env" type="text" class="w-full bg-slate-950 border border-slate-700 rounded-lg px-3 py-1.5 text-xs text-slate-200 font-mono" />
        </div>
      </div>
    </div>

    <!-- 2. Storage 設定 -->
    <div class="bg-slate-900/80 border border-slate-800 rounded-xl p-4 space-y-4">
      <h3 class="text-sm font-bold text-slate-200 flex items-center gap-2"><span>💾</span> ストレージ・保全パス (Storage)</h3>
      <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
        <div>
          <label class="block text-xs text-slate-400 mb-1">データベースパス (db_path)</label>
          <input v-model="config.storage.db_path" type="text" class="w-full bg-slate-950 border border-slate-700 rounded-lg px-3 py-1.5 text-xs text-slate-200 font-mono" />
        </div>
        <div>
          <label class="block text-xs text-slate-400 mb-1">ローカル画像保存先 (local_media_dir)</label>
          <input v-model="config.storage.local_media_dir" type="text" class="w-full bg-slate-950 border border-slate-700 rounded-lg px-3 py-1.5 text-xs text-slate-200 font-mono" />
        </div>
        <div>
          <label class="block text-xs text-slate-400 mb-1">Stashapp ディレクトリ (stash_dir)</label>
          <input v-model="config.storage.stash_dir" type="text" class="w-full bg-slate-950 border border-slate-700 rounded-lg px-3 py-1.5 text-xs text-slate-200 font-mono" />
        </div>
        <div>
          <label class="block text-xs text-slate-400 mb-1">バックアップダンプ先 (dumps_dir)</label>
          <input v-model="config.storage.dumps_dir" type="text" class="w-full bg-slate-950 border border-slate-700 rounded-lg px-3 py-1.5 text-xs text-slate-200 font-mono" />
        </div>
      </div>
    </div>
  </div>
</template>
