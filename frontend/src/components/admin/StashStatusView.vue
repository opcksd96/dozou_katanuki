<script setup lang="ts">
import { computed, ref } from 'vue';
import { models } from '../../../wailsjs/go/models';

const props = withDefaults(defineProps<{
  config?: models.AppConfig | null;
}>(), {
  config: null,
});

const stashPort = computed(() => props.config?.network?.stash_port || 9999);
const isStashEnabled = computed(() => props.config?.storage?.stash_enabled ?? true);
const currentHostname = computed(() => {
  if (typeof window !== 'undefined' && window.location?.hostname) {
    return window.location.hostname;
  }
  return '127.0.0.1';
});
const stashUrl = computed(() => `http://${currentHostname.value}:${stashPort.value}`);

const copied = ref(false);
const copyStashUrl = () => {
  if (typeof navigator !== 'undefined' && navigator.clipboard) {
    navigator.clipboard.writeText(stashUrl.value);
    copied.value = true;
    setTimeout(() => { copied.value = false; }, 2000);
  }
};

const openStash = () => {
  window.open(stashUrl.value, '_blank');
};
</script>

<template>
  <div class="space-y-6">
    <div class="bg-slate-900/80 border border-slate-800 rounded-xl p-5 space-y-4">
      <div class="flex items-center justify-between">
        <div class="flex items-center gap-3">
          <div class="w-10 h-10 rounded-xl bg-purple-950/60 border border-purple-500/30 flex items-center justify-center text-lg">
            📦
          </div>
          <div>
            <h3 class="text-sm font-bold text-slate-200">
              Stashapp プロセス連携ステータス (SPEC-ADMINBOARD-001)
            </h3>
            <p class="text-xs text-slate-400">
              バックエンド起動時に自動連動する動画メディア管理サーバー (0.0.0.0:{{ stashPort }} LAN透過待受)
            </p>
          </div>
        </div>

        <span
          class="px-2.5 py-1 rounded-full text-xs font-semibold font-mono"
          :class="isStashEnabled ? 'bg-emerald-500/20 text-emerald-400 border border-emerald-500/30' : 'bg-slate-800 text-slate-400 border border-slate-700'"
        >
          {{ isStashEnabled ? '● 連動モード' : '○ 「Stash使わんし！」' }}
        </span>
      </div>

      <div class="p-4 bg-slate-950/60 rounded-lg border border-slate-800/80 space-y-3">
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-3 text-xs font-mono">
          <div>
            <span class="text-slate-500 block mb-0.5">Stash WebUI Endpoint (接続中ホスト準拠):</span>
            <div class="flex items-center gap-2">
              <span class="text-blue-400 font-semibold">{{ stashUrl }}</span>
              <button
                type="button"
                @click="copyStashUrl"
                class="px-2 py-0.5 bg-slate-800 hover:bg-slate-700 text-slate-300 text-[11px] rounded transition-colors"
                title="URLをコピー"
              >
                {{ copied ? '✓ コピー済' : 'コピー' }}
              </button>
            </div>
          </div>
          <div>
            <span class="text-slate-500 block mb-0.5">Stash Data Directory:</span>
            <span class="text-slate-300">{{ config?.storage?.stash_dir || './stash' }}</span>
          </div>
        </div>

        <div class="pt-2 border-t border-slate-800/60 flex items-center justify-between">
          <span class="text-xs text-slate-400">
            {{ isStashEnabled ? 'Stash サーバーの管理画面をブラウザ別窓で開きます（LAN端末からも同IPで開けます）' : '現在 Stashapp の自動起動は停止されています' }}
          </span>
          <button
            @click="openStash"
            :disabled="!isStashEnabled"
            class="px-4 py-2 bg-purple-600 hover:bg-purple-500 disabled:bg-slate-800 disabled:text-slate-600 disabled:cursor-not-allowed text-white text-xs font-bold rounded-lg transition-colors shadow-sm flex items-center gap-1.5"
          >
            <span>🔗</span> Stash WebUI を開く
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
