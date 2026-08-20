<script setup lang="ts">
import { ref, computed } from 'vue';
import { models } from '../../../wailsjs/go/models';

const props = defineProps<{
  config: models.AppConfig | null;
  loadingConfig: boolean;
  savingConfig: boolean;
  saveStatus: { success: boolean; message: string } | null;
}>();

const emit = defineEmits<{
  (e: 'saveConfig'): void;
  (e: 'loadConfig'): void;
}>();

const copied = ref(false);

const allowedNetworksText = computed({
  get: () => {
    if (!props.config?.broadcast?.allowed_networks) return '';
    return props.config.broadcast.allowed_networks.join(', ');
  },
  set: (val: string) => {
    if (!props.config) return;
    if (!props.config.broadcast) {
      props.config.broadcast = {
        enabled: false,
        allowed_networks: [],
      } as unknown as models.BroadcastConfig;
    }
    props.config.broadcast.allowed_networks = val
      .split(',')
      .map(s => s.trim())
      .filter(s => s.length > 0);
  },
});

const sampleCastURL = computed(() => {
  const port = props.config?.network?.middleware_port || 5175;
  return `http://<あなたのPCのIPアドレス>:${port}`;
});

const copyCastURL = () => {
  navigator.clipboard.writeText(sampleCastURL.value);
  copied.value = true;
  setTimeout(() => {
    copied.value = false;
  }, 2000);
};
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

      <!-- 3. LAN Broadcast ＆ キャスト配信 (SPEC-SCHEDULER-001) -->
      <div class="bg-slate-900/80 border border-slate-800 rounded-xl p-4 space-y-4">
        <div class="flex items-center justify-between">
          <h3 class="text-sm font-bold text-slate-200 flex items-center gap-2">
            <span>📡</span> 家庭内LAN Broadcast ＆ キャスト配信 (SPEC-SCHEDULER-001)
          </h3>
          <label class="flex items-center gap-2 cursor-pointer select-none">
            <span class="text-xs font-semibold" :class="config.broadcast.enabled ? 'text-emerald-400' : 'text-slate-500'">
              {{ config.broadcast.enabled ? '配信: オンライン' : '配信: 停止中' }}
            </span>
            <input
              type="checkbox"
              v-model="config.broadcast.enabled"
              class="sr-only peer"
            />
            <div class="w-9 h-5 bg-slate-700 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-slate-300 after:border after:rounded-full after:h-4 after:w-4 after:transition-all peer-checked:bg-emerald-600 relative"></div>
          </label>
        </div>

        <p class="text-[11px] text-slate-400 leading-relaxed bg-slate-950/60 p-2.5 rounded-lg border border-slate-800/80">
          <strong class="text-slate-300">💡 キャスト配信機能:</strong>
          PCのデスクを離れても、同一LAN内のスマホ・タブレット・スマートTV等のブラウザからメディアストリーミング・動画・タイムラインを直接鑑賞できます。
          CORS制限を自動中和し、許可ネットワーク外からのアクセスは <span class="text-rose-400 font-mono font-bold">403 Forbidden</span> で確実に遮断します。
        </p>

        <div class="grid grid-cols-1 sm:grid-cols-3 gap-4">
          <div>
            <label class="block text-xs text-slate-400 mb-1">配信ポート (middleware_port)</label>
            <input
              v-model.number="config.network.middleware_port"
              type="number"
              class="w-full bg-slate-950 border border-slate-700 rounded-lg px-3 py-1.5 text-xs text-slate-200 focus:outline-none focus:border-blue-500 font-mono"
            />
          </div>
          <div>
            <label class="block text-xs text-slate-400 mb-1">パブリックバインドアドレス</label>
            <input
              v-model="config.network.public_bind_address"
              type="text"
              class="w-full bg-slate-950 border border-slate-700 rounded-lg px-3 py-1.5 text-xs text-slate-200 focus:outline-none focus:border-blue-500 font-mono"
            />
          </div>
          <div>
            <label class="block text-xs text-slate-400 mb-1">セキュリティゲートウェイ状態</label>
            <div class="bg-slate-950 border border-slate-800 rounded-lg px-3 py-1.5 text-xs flex items-center gap-1.5 text-emerald-400 font-medium">
              <span>🛡️</span> IP/CIDR フィルタリング有効
            </div>
          </div>
        </div>

        <div>
          <label class="block text-xs text-slate-400 mb-1">許可ネットワーク / CIDR サブネット (カンマ区切り)</label>
          <input
            v-model="allowedNetworksText"
            type="text"
            placeholder="192.168.0.0/16, 10.0.0.0/8, 172.16.0.0/12, 127.0.0.1/32"
            class="w-full bg-slate-950 border border-slate-700 rounded-lg px-3 py-1.5 text-xs text-slate-200 focus:outline-none focus:border-blue-500 font-mono"
          />
          <p class="text-[10px] text-slate-500 mt-1">※ 指定したCIDR/IP以外からのアクセスはすべて自動的に即座に遮断（403）されます。</p>
        </div>

        <div v-if="config.broadcast.enabled" class="bg-emerald-950/30 border border-emerald-500/30 rounded-lg p-3 flex flex-col sm:flex-row sm:items-center justify-between gap-2">
          <div class="text-xs text-slate-300 flex items-center gap-2">
            <span class="text-emerald-400">📱 スマホ等からのアクセスURL:</span>
            <code class="bg-slate-950 px-2 py-0.5 rounded text-emerald-300 font-mono text-[11px] border border-emerald-800/50">
              {{ sampleCastURL }}
            </code>
          </div>
          <button
            type="button"
            @click="copyCastURL"
            class="px-2.5 py-1 bg-emerald-800/80 hover:bg-emerald-700 text-white rounded text-[11px] font-semibold transition-colors flex items-center gap-1 self-start sm:self-auto"
          >
            <span>{{ copied ? '✅ コピー完了' : '📋 URLをコピー' }}</span>
          </button>
        </div>
      </div>

      <!-- 4. Network ＆ 内部ポート設定 -->
      <div class="bg-slate-900/80 border border-slate-800 rounded-xl p-4 space-y-4">
        <h3 class="text-sm font-bold text-slate-200 flex items-center gap-2">
          <span>🔌</span> 内部ネットワーク設定 (Network)
        </h3>
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

      <!-- 5. Scheduler 設定 -->
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

      <!-- 6. Appearance フォント設定 -->
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

