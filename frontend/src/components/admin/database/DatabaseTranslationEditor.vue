<!-- frontend/src/components/admin/database/DatabaseTranslationEditor.vue (100行以下) -->
<script setup lang="ts">
import { reactive, watch } from 'vue';
import { models } from '../../../../wailsjs/go/models';

const props = defineProps<{
  article: models.RenderTree | null;
  saving: boolean;
}>();
const emit = defineEmits<{ (e: 'save', ja: string, en: string, zh: string): void }>();

const editBuffer = reactive({ ja: '', en: '', zh: '' });

watch(() => props.article, (art) => {
  if (art) {
    editBuffer.ja = art.content?.ja || '';
    editBuffer.en = art.content?.en || '';
    editBuffer.zh = art.content?.zh || '';
  }
}, { immediate: true });
</script>

<template>
  <div v-if="article" class="space-y-4 bg-slate-900/60 border border-slate-800 rounded-xl p-4 flex flex-col h-full">
    <div class="flex justify-between items-center border-b border-slate-800 pb-2">
      <h4 class="text-xs font-bold text-slate-200">記事 ID: <span class="font-mono text-blue-400">{{ article.id }}</span></h4>
      <button @click="emit('save', editBuffer.ja, editBuffer.en, editBuffer.zh)" :disabled="saving" class="px-4 py-1.5 bg-blue-600 hover:bg-blue-500 text-white text-xs font-bold rounded-lg disabled:opacity-50">
        {{ saving ? '保存中...' : '💾 翻訳保存 (Ctrl+S)' }}
      </button>
    </div>
    <div class="space-y-1">
      <label class="text-[11px] text-slate-400 font-semibold">原文 (Original Text)</label>
      <div class="p-2.5 bg-slate-950 rounded-lg text-xs text-slate-300 font-mono select-text max-h-24 overflow-y-auto">{{ article.content.original }}</div>
    </div>
    <div class="grid grid-cols-1 gap-3 flex-1 overflow-y-auto">
      <div>
        <label class="text-[11px] text-emerald-400 font-semibold mb-1 block">🇯🇵 日本語翻訳 (JA)</label>
        <textarea v-model="editBuffer.ja" rows="3" class="w-full bg-slate-950 border border-slate-700 rounded-lg p-2 text-xs text-slate-200 focus:border-blue-500"></textarea>
      </div>
      <div>
        <label class="text-[11px] text-blue-400 font-semibold mb-1 block">🇺🇸 英語翻訳 (EN)</label>
        <textarea v-model="editBuffer.en" rows="3" class="w-full bg-slate-950 border border-slate-700 rounded-lg p-2 text-xs text-slate-200 focus:border-blue-500"></textarea>
      </div>
      <div>
        <label class="text-[11px] text-purple-400 font-semibold mb-1 block">🇨🇳 中国語翻訳 (ZH)</label>
        <textarea v-model="editBuffer.zh" rows="3" class="w-full bg-slate-950 border border-slate-700 rounded-lg p-2 text-xs text-slate-200 focus:border-blue-500"></textarea>
      </div>
    </div>
  </div>
  <div v-else class="h-full flex items-center justify-center border border-slate-800 rounded-xl bg-slate-900/30 text-xs text-slate-500">
    左側の一覧から記事を選択して翻訳を編集してください
  </div>
</template>
