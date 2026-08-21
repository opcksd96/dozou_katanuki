<!-- frontend/src/components/admin/skin/SkinFontSettings.vue (100行以下) -->
<script setup lang="ts">
import { models } from '../../../../wailsjs/go/models';
defineProps<{
  config: models.AppConfig | null;
  fontPresets: { ja: Array<{ label: string; value: string }>; en: Array<{ label: string; value: string }>; zh: Array<{ label: string; value: string }> };
  savingConfig: boolean;
}>();
defineEmits<{ (e: 'save'): void }>();
</script>

<template>
  <div v-if="config" class="bg-slate-900/60 p-4 border border-slate-800 rounded-xl space-y-3">
    <div class="flex items-center justify-between">
      <h4 class="text-xs font-bold text-slate-200">🔤 フォント設定 (config.json)</h4>
      <button @click="$emit('save')" :disabled="savingConfig" class="px-3.5 py-1 bg-blue-600 hover:bg-blue-500 text-white text-xs font-bold rounded disabled:opacity-50">
        {{ savingConfig ? '保存中...' : '💾 フォント保存' }}
      </button>
    </div>
    <div class="grid grid-cols-1 sm:grid-cols-3 gap-3">
      <div>
        <label class="block text-[11px] text-slate-400 mb-1">🇯🇵 日本語フォント</label>
        <input v-model="config.appearance.font_family_ja" type="text" class="w-full bg-slate-950 border border-slate-700 rounded p-1.5 text-xs text-slate-200 font-mono" />
      </div>
      <div>
        <label class="block text-[11px] text-slate-400 mb-1">🇺🇸 英語フォント</label>
        <input v-model="config.appearance.font_family_en" type="text" class="w-full bg-slate-950 border border-slate-700 rounded p-1.5 text-xs text-slate-200 font-mono" />
      </div>
      <div>
        <label class="block text-[11px] text-slate-400 mb-1">🇨🇳 中国語フォント</label>
        <input v-model="config.appearance.font_family_zh" type="text" class="w-full bg-slate-950 border border-slate-700 rounded p-1.5 text-xs text-slate-200 font-mono" />
      </div>
    </div>
  </div>
</template>
