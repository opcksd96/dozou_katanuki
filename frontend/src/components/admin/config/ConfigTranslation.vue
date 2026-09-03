<!-- frontend/src/components/admin/config/ConfigTranslation.vue (100行以下) -->
<script setup lang="ts">
import { ref } from 'vue';
import { models } from '../../../../wailsjs/go/models';

const props = defineProps<{ config: models.AppConfig }>();

if (!props.config.translation) {
  (props.config as any).translation = {
    provider: 'deepl',
    deepl_api_key: '',
    google_translate_api_key: ''
  };
}

const showDeeplKey = ref(false);
const showGoogleKey = ref(false);
</script>

<template>
  <div class="bg-slate-900/80 border border-slate-800 rounded-xl p-4 space-y-4">
    <div class="flex items-center justify-between">
      <h3 class="text-sm font-bold text-slate-200 flex items-center gap-2">
        <span>🌐</span> 翻訳API設定 (Translation API)
      </h3>
      <span class="text-[11px] px-2 py-0.5 rounded bg-blue-950/60 text-blue-400 border border-blue-800/50">
        Single Source of Truth
      </span>
    </div>

    <p class="text-xs text-slate-400">
      記事翻訳（ja / en / zh）に使用する外部翻訳サービスの設定です。APIキーは <code class="text-blue-300 font-mono">config.json</code> に保存され、実行時に安全に渡されます。
    </p>

    <div class="space-y-4">
      <!-- 1. プロバイダー選択 -->
      <div>
        <label class="block text-xs font-medium text-slate-300 mb-1">優先翻訳エンジン (Provider)</label>
        <select v-model="config.translation.provider" class="w-full bg-slate-950 border border-slate-700 rounded-lg px-3 py-1.5 text-xs text-slate-200 focus:border-blue-500 focus:outline-none">
          <option value="deepl">DeepL API（推奨・自然な翻訳）</option>
          <option value="google">Google Cloud Translation API</option>
          <option value="none">無効（翻訳を行わない）</option>
        </select>
      </div>

      <!-- 2. DeepL API Key -->
      <div v-if="config.translation.provider === 'deepl' || !config.translation.provider">
        <div class="flex items-center justify-between mb-1">
          <label class="block text-xs text-slate-400">DeepL API Key (Free版 / Pro版)</label>
          <button type="button" @click="showDeeplKey = !showDeeplKey" class="text-[11px] text-blue-400 hover:underline">
            {{ showDeeplKey ? '🙈 非表示' : '👁️ 表示' }}
          </button>
        </div>
        <input
          v-model="config.translation.deepl_api_key"
          :type="showDeeplKey ? 'text' : 'password'"
          placeholder="例: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx:fx"
          class="w-full bg-slate-950 border border-slate-700 rounded-lg px-3 py-1.5 text-xs text-slate-200 font-mono placeholder:text-slate-600 focus:border-blue-500 focus:outline-none"
        />
        <p class="text-[10px] text-slate-500 mt-1">※ Free版キー末尾の <code class="text-slate-400">:fx</code> に自動対応します。</p>
      </div>

      <!-- 3. Google Translate API Key -->
      <div v-if="config.translation.provider === 'google'">
        <div class="flex items-center justify-between mb-1">
          <label class="block text-xs text-slate-400">Google Cloud Translation API Key</label>
          <button type="button" @click="showGoogleKey = !showGoogleKey" class="text-[11px] text-blue-400 hover:underline">
            {{ showGoogleKey ? '🙈 非表示' : '👁️ 表示' }}
          </button>
        </div>
        <input
          v-model="config.translation.google_translate_api_key"
          :type="showGoogleKey ? 'text' : 'password'"
          placeholder="例: AIzaSy..."
          class="w-full bg-slate-950 border border-slate-700 rounded-lg px-3 py-1.5 text-xs text-slate-200 font-mono placeholder:text-slate-600 focus:border-blue-500 focus:outline-none"
        />
      </div>
    </div>
  </div>
</template>
