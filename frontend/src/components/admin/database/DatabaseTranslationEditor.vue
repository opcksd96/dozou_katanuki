<!-- frontend/src/components/admin/database/DatabaseTranslationEditor.vue (100行以下) -->
<script setup lang="ts">
import { ref, reactive, watch } from 'vue';
import { models } from '../../../../wailsjs/go/models';

const props = defineProps<{
  article: models.RenderTree | null; saving: boolean; translating?: boolean;
}>();
const emit = defineEmits<{
  (e: 'save', ja: string, en: string, zh: string): void;
  (e: 'autoTranslate', autoSave: boolean): void;
}>();

const autoSave = ref(false);
const activeLang = ref<'ja' | 'en' | 'zh' | 'all'>('ja');
const editBuffer = reactive({ ja: '', en: '', zh: '' });

const stripHtml = (htmlStr: string) => {
  if (!htmlStr) return '';
  const doc = new DOMParser().parseFromString(htmlStr, 'text/html');
  return doc.body.textContent || '';
};

watch(() => [props.article, props.article?.content?.ja, props.article?.content?.en, props.article?.content?.zh], () => {
  if (props.article) {
    editBuffer.ja = stripHtml(props.article.content?.ja || '');
    editBuffer.en = stripHtml(props.article.content?.en || '');
    editBuffer.zh = stripHtml(props.article.content?.zh || '');
  }
}, { immediate: true, deep: true });
</script>

<template>
  <div v-if="article" class="space-y-2.5 bg-slate-900/60 border border-slate-800 rounded-xl p-3.5 flex flex-col h-full">
    <!-- ヘッダー -->
    <div class="flex justify-between items-center border-b border-slate-800 pb-2 flex-wrap gap-2">
      <h4 class="text-xs font-bold text-slate-200">記事 ID: <span class="font-mono text-blue-400">{{ article.id }}</span></h4>
      <div class="flex items-center gap-2">
        <label class="text-[11px] text-slate-400 flex items-center gap-1 cursor-pointer select-none hover:text-slate-300">
          <input type="checkbox" v-model="autoSave" class="rounded border-slate-700 bg-slate-950 text-blue-600 focus:ring-0 w-3.5 h-3.5">
          <span>翻訳時に自動保存</span>
        </label>
        <button @click="emit('autoTranslate', autoSave)" :disabled="translating || saving" class="px-2.5 py-1 bg-emerald-600 hover:bg-emerald-500 text-white text-xs font-bold rounded-lg disabled:opacity-50 flex items-center gap-1">
          <span>🤖</span> {{ translating ? '翻訳中...' : '自動翻訳' }}
        </button>
        <button @click="emit('save', editBuffer.ja, editBuffer.en, editBuffer.zh)" :disabled="saving || translating" class="px-2.5 py-1 bg-blue-600 hover:bg-blue-500 text-white text-xs font-bold rounded-lg disabled:opacity-50">
          {{ saving ? '保存中...' : '💾 保存' }}
        </button>
      </div>
    </div>

    <!-- 原文プレビュー（コンパクトな保護ビュー：最大2行程度に制限して下を圧迫しない） -->
    <div class="space-y-1">
      <label class="text-[10px] text-slate-400 font-semibold flex justify-between">
        <span>原文プレビュー (保護ビュー)</span>
        <span class="text-slate-500 font-mono text-[9px] truncate max-w-[200px]">{{ article.author?.handle ? '@' + article.author.handle : '' }}</span>
      </label>
      <div class="p-2 bg-slate-950/80 rounded-lg text-[11px] text-slate-300 font-sans select-text max-h-12 overflow-y-auto border border-slate-800/80 leading-relaxed" v-html="article.content.original"></div>
    </div>

    <!-- 言語切り替えタブ -->
    <div class="flex items-center justify-between border-b border-slate-800/80 pb-1 text-[11px]">
      <div class="flex gap-1">
        <button v-for="l in [
          { id: 'ja', label: '🇯🇵 日本語 (JA)', has: !!editBuffer.ja },
          { id: 'en', label: '🇺🇸 英語 (EN)', has: !!editBuffer.en },
          { id: 'zh', label: '🇨🇳 中国語 (ZH)', has: !!editBuffer.zh },
          { id: 'all', label: '📋 全て並列', has: true }
        ]" :key="l.id" @click="activeLang = l.id as any" class="px-2 py-1 rounded font-bold transition-colors flex items-center gap-1" :class="activeLang === l.id ? 'bg-blue-600/30 text-blue-300 border border-blue-500/40' : 'text-slate-400 hover:text-slate-200'">
          <span>{{ l.label }}</span>
          <span v-if="l.id !== 'all' && !l.has" class="w-1.5 h-1.5 rounded-full bg-amber-500" title="未入力"></span>
        </button>
      </div>
    </div>

    <!-- 翻訳入力エリア（広々とした編集領域） -->
    <div class="flex-1 flex flex-col min-h-[140px] overflow-y-auto gap-2">
      <div v-if="activeLang === 'ja' || activeLang === 'all'" class="flex-1 flex flex-col min-h-[90px]">
        <label class="text-[10px] text-emerald-400 font-semibold mb-0.5">🇯🇵 日本語翻訳 (JA)</label>
        <textarea v-model="editBuffer.ja" placeholder="日本語訳を入力（自動翻訳ボタンで取得可能）" class="w-full flex-1 min-h-[70px] bg-slate-950 border border-slate-700 rounded-lg p-2.5 text-xs text-slate-100 focus:border-blue-500 font-sans leading-relaxed"></textarea>
      </div>
      <div v-if="activeLang === 'en' || activeLang === 'all'" class="flex-1 flex flex-col min-h-[90px]">
        <label class="text-[10px] text-blue-400 font-semibold mb-0.5">🇺🇸 英語翻訳 (EN)</label>
        <textarea v-model="editBuffer.en" placeholder="英語訳を入力" class="w-full flex-1 min-h-[70px] bg-slate-950 border border-slate-700 rounded-lg p-2.5 text-xs text-slate-100 focus:border-blue-500 font-sans leading-relaxed"></textarea>
      </div>
      <div v-if="activeLang === 'zh' || activeLang === 'all'" class="flex-1 flex flex-col min-h-[90px]">
        <label class="text-[10px] text-purple-400 font-semibold mb-0.5">🇨🇳 中国語翻訳 (ZH)</label>
        <textarea v-model="editBuffer.zh" placeholder="中国語訳を入力" class="w-full flex-1 min-h-[70px] bg-slate-950 border border-slate-700 rounded-lg p-2.5 text-xs text-slate-100 focus:border-blue-500 font-sans leading-relaxed"></textarea>
      </div>
    </div>
  </div>
  <div v-else class="h-full flex items-center justify-center border border-slate-800 rounded-xl bg-slate-900/30 text-xs text-slate-500">
    左側の一覧から記事を選択して翻訳を編集してください
  </div>
</template>
