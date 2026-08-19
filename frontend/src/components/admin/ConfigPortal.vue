<script setup lang="ts">
import { models } from '../../../wailsjs/go/models';

defineProps<{
  config: models.AppConfig | null;
  loadingConfig: boolean;
  savingConfig: boolean;
  saveStatus: { success: boolean; message: string } | null;
}>();

const emit = defineEmits<{
  (e: 'saveConfig'): void;
  (e: 'loadConfig'): void;
}>();
</script>

<template>
  <div class="space-y-6">
    <!-- ローディング状態 -->
    <div v-if="loadingConfig && !config" class="py-12 text-center text-slate-400">
      <span class="animate-spin text-xl inline-block mr-2">⏳</span> 設定を読み込み中...
    </div>

    <form v-else-if="config" @submit.prevent="emit('saveConfig')" class="space-y-6">
      <!-- ステータス通知 -->
      <div
        v-if="saveStatus"
        class="p-3 rounded-lg text-xs font-semibold flex items-center justify-between transition-all"
        :class="saveStatus.success ? 'bg-emerald-950/80 border border-emerald-500/50 text-emerald-300' : 'bg-rose-950/80 border border-rose-500/50 text-rose-300'"
      >
        <span>{{ saveStatus.message }}</span>
      </div>

      <!-- 1. System 設定 -->
      <div class="bg-slate-900/80 border border-slate-800 rounded-xl p-4 space-y-4">
        <h3 class="text-sm font-bold text-slate-200 flex items-center gap-2">
          <span>🌐</span> システム基本設定 (System)
        </h3>
        <div class="grid grid-cols-1 sm:grid-cols-3 gap-4">
          <div>
            <label class="block text-xs text-slate-400 mb-1">システム言語 (UI/翻訳)</label>
            <select
              v-model="config.system.language"
              class="w-full bg-slate-950 border border-slate-700 rounded-lg px-3 py-1.5 text-xs text-slate-200 focus:outline-none focus:border-blue-500 font-medium"
            >
              <option value="ja">日本語 (ja)</option>
              <option value="en">English (en)</option>
              <option value="zh">简体中文 (zh)</option>
            </select>
          </div>
          <div>
            <label class="block text-xs text-slate-400 mb-1">デフォルトプラットフォーム</label>
            <select
              v-model="config.system.default_framework"
              class="w-full bg-slate-950 border border-slate-700 rounded-lg px-3 py-1.5 text-xs text-slate-200 focus:outline-none focus:border-blue-500"
            >
              <option value="twitter">Twitter / X</option>
              <option value="bsky">Bluesky</option>
            </select>
          </div>
          <div>
            <label class="block text-xs text-slate-400 mb-1">環境モード (env)</label>
            <input
              v-model="config.system.env"
              type="text"
              class="w-full bg-slate-950 border border-slate-700 rounded-lg px-3 py-1.5 text-xs text-slate-200 focus:outline-none focus:border-blue-500 font-mono"
            />
          </div>
        </div>
      </div>

      <!-- 2. Storage ＆ 「Stash使わんし！」モード -->
      <div class="bg-slate-900/80 border border-slate-800 rounded-xl p-4 space-y-4">
        <div class="flex items-center justify-between">
          <h3 class="text-sm font-bold text-slate-200 flex items-center gap-2">
            <span>💾</span> ストレージ ＆ 保存先設定 (Storage)
          </h3>
          <!-- 「Stash使わんし！」モードトグル -->
          <label class="flex items-center gap-2 cursor-pointer select-none">
            <span class="text-xs font-semibold" :class="config.storage.stash_enabled ? 'text-blue-400' : 'text-amber-400'">
              {{ config.storage.stash_enabled ? 'Stash連携: 有効' : '「Stash使わんし！」モード: 有効' }}
            </span>
            <input
              type="checkbox"
              v-model="config.storage.stash_enabled"
              class="sr-only peer"
            />
            <div class="w-9 h-5 bg-slate-700 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-slate-300 after:border after:rounded-full after:h-4 after:w-4 after:transition-all peer-checked:bg-blue-600 relative"></div>
          </label>
        </div>

        <p class="text-[11px] text-slate-400 leading-relaxed bg-slate-950/60 p-2.5 rounded-lg border border-slate-800/80">
          <strong class="text-slate-300">💡 「Stash使わんし！」モード:</strong>
          無効（OFF）にすると、外部 Stashapp プロセスを起動せず、<code>/media-local/</code> 経由で物理保存ディレクトリから軽量かつダイレクトにメディアを配信します。
        </p>

        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <div>
            <label class="block text-xs text-slate-400 mb-1">SQLite データベースパス</label>
            <input
              v-model="config.storage.db_path"
              type="text"
              required
              class="w-full bg-slate-950 border border-slate-700 rounded-lg px-3 py-1.5 text-xs text-slate-200 focus:outline-none focus:border-blue-500 font-mono"
            />
          </div>
          <div>
            <label class="block text-xs text-slate-400 mb-1">物理メディア保存ディレクトリ</label>
            <input
              v-model="config.storage.local_media_dir"
              type="text"
              required
              class="w-full bg-slate-950 border border-slate-700 rounded-lg px-3 py-1.5 text-xs text-slate-200 focus:outline-none focus:border-blue-500 font-mono"
            />
          </div>
          <div>
            <label class="block text-xs text-slate-400 mb-1">Stash ディレクトリパス</label>
            <input
              v-model="config.storage.stash_dir"
              type="text"
              class="w-full bg-slate-950 border border-slate-700 rounded-lg px-3 py-1.5 text-xs text-slate-200 focus:outline-none focus:border-blue-500 font-mono"
            />
          </div>
          <div>
            <label class="block text-xs text-slate-400 mb-1">バックアップダンプ先ディレクトリ</label>
            <input
              v-model="config.storage.dumps_dir"
              type="text"
              class="w-full bg-slate-950 border border-slate-700 rounded-lg px-3 py-1.5 text-xs text-slate-200 focus:outline-none focus:border-blue-500 font-mono"
            />
          </div>
        </div>
      </div>

      <!-- 3. Network ＆ ポート設定 -->
      <div class="bg-slate-900/80 border border-slate-800 rounded-xl p-4 space-y-4">
        <h3 class="text-sm font-bold text-slate-200 flex items-center gap-2">
          <span>🔌</span> ネットワーク設定 (Network)
        </h3>
        <p class="text-[11px] text-slate-400">
          ※ Wails In-Memory プロキシ統合により、Middleware/Backend 間通信はメモリ内で直接完結しています。
        </p>
        <div class="grid grid-cols-1 sm:grid-cols-3 gap-4">
          <div>
            <label class="block text-xs text-slate-400 mb-1">Stash 外部ポート (stash_port)</label>
            <input
              v-model.number="config.network.stash_port"
              type="number"
              class="w-full bg-slate-950 border border-slate-700 rounded-lg px-3 py-1.5 text-xs text-slate-200 focus:outline-none focus:border-blue-500 font-mono"
            />
          </div>
          <div>
            <label class="block text-xs text-slate-400 mb-1">フロントエンドポート (dev)</label>
            <input
              v-model.number="config.network.frontend_port"
              type="number"
              class="w-full bg-slate-950 border border-slate-700 rounded-lg px-3 py-1.5 text-xs text-slate-200 focus:outline-none focus:border-blue-500 font-mono"
            />
          </div>
          <div>
            <label class="block text-xs text-slate-400 mb-1">内部バインドアドレス</label>
            <input
              v-model="config.network.internal_bind_address"
              type="text"
              class="w-full bg-slate-950 border border-slate-700 rounded-lg px-3 py-1.5 text-xs text-slate-200 focus:outline-none focus:border-blue-500 font-mono"
            />
          </div>
        </div>
      </div>

      <!-- 4. Scheduler 設定 -->
      <div class="bg-slate-900/80 border border-slate-800 rounded-xl p-4 space-y-4">
        <h3 class="text-sm font-bold text-slate-200 flex items-center gap-2">
          <span>⏰</span> スケジューラ ＆ 自動運用設定 (Scheduler)
        </h3>
        <div class="grid grid-cols-1 sm:grid-cols-3 gap-4">
          <div>
            <label class="block text-xs text-slate-400 mb-1">ポーリング間隔 (秒)</label>
            <input
              v-model.number="config.scheduler.poll_interval_sec"
              type="number"
              min="10"
              class="w-full bg-slate-950 border border-slate-700 rounded-lg px-3 py-1.5 text-xs text-slate-200 focus:outline-none focus:border-blue-500 font-mono"
            />
          </div>
          <div>
            <label class="block text-xs text-slate-400 mb-1">バックアップ周期 (時間)</label>
            <input
              v-model.number="config.scheduler.backup_interval_hours"
              type="number"
              min="1"
              class="w-full bg-slate-950 border border-slate-700 rounded-lg px-3 py-1.5 text-xs text-slate-200 focus:outline-none focus:border-blue-500 font-mono"
            />
          </div>
          <div>
            <label class="block text-xs text-slate-400 mb-1">最大バックアップ保持世代数</label>
            <input
              v-model.number="config.scheduler.max_backup_generations"
              type="number"
              min="1"
              class="w-full bg-slate-950 border border-slate-700 rounded-lg px-3 py-1.5 text-xs text-slate-200 focus:outline-none focus:border-blue-500 font-mono"
            />
          </div>
        </div>
      </div>

      <!-- 5. Appearance フォント設定 -->
      <div class="bg-slate-900/80 border border-slate-800 rounded-xl p-4 space-y-4">
        <h3 class="text-sm font-bold text-slate-200 flex items-center gap-2">
          <span>🎨</span> 外観 ＆ 多言語フォント設定 (Appearance)
        </h3>
        <div class="grid grid-cols-1 sm:grid-cols-3 gap-4">
          <div>
            <label class="block text-xs text-slate-400 mb-1">日本語フォント (JA)</label>
            <input
              v-model="config.appearance.font_family_ja"
              type="text"
              class="w-full bg-slate-950 border border-slate-700 rounded-lg px-3 py-1.5 text-xs text-slate-200 focus:outline-none focus:border-blue-500 font-mono text-[11px]"
            />
          </div>
          <div>
            <label class="block text-xs text-slate-400 mb-1">英語フォント (EN)</label>
            <input
              v-model="config.appearance.font_family_en"
              type="text"
              class="w-full bg-slate-950 border border-slate-700 rounded-lg px-3 py-1.5 text-xs text-slate-200 focus:outline-none focus:border-blue-500 font-mono text-[11px]"
            />
          </div>
          <div>
            <label class="block text-xs text-slate-400 mb-1">中国語フォント (ZH)</label>
            <input
              v-model="config.appearance.font_family_zh"
              type="text"
              class="w-full bg-slate-950 border border-slate-700 rounded-lg px-3 py-1.5 text-xs text-slate-200 focus:outline-none focus:border-blue-500 font-mono text-[11px]"
            />
          </div>
        </div>
      </div>

      <!-- 下部アクションバー -->
      <div class="flex items-center justify-end gap-3 pt-2">
        <button
          type="button"
          @click="emit('loadConfig')"
          :disabled="loadingConfig || savingConfig"
          class="px-4 py-2 bg-slate-800 hover:bg-slate-700 text-slate-300 rounded-lg text-xs font-semibold transition-colors"
        >
          リセット
        </button>
        <button
          type="submit"
          :disabled="savingConfig"
          class="px-5 py-2 bg-blue-600 hover:bg-blue-500 disabled:bg-slate-700 text-white rounded-lg text-xs font-bold transition-colors shadow-sm flex items-center gap-1.5"
        >
          <span v-if="savingConfig" class="animate-spin">⏳</span>
          <span v-else>💾</span>
          <span>設定を保存 (config.json)</span>
        </button>
      </div>
    </form>
  </div>
</template>
