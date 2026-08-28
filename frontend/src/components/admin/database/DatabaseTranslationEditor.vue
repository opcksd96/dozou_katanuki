<!-- frontend/src/components/admin/database/DatabaseTranslationEditor.vue (100行以下 - SPEC-PRINCIPLE-001) -->
<script setup lang="ts">
import { ref, reactive, watch } from 'vue';
import { models } from '../../../../wailsjs/go/models';
import PostTrashModal from './PostTrashModal.vue';
import ArticleUrlBadges from './ArticleUrlBadges.vue';

const props = defineProps<{ article: models.RenderTree | null; saving: boolean; translating?: boolean }>();
const emit = defineEmits<{
  (e: 'save', ja: string, en: string, zh: string): void;
  (e: 'autoTranslate', autoSave: boolean): void;
  (e: 'trash', id: string, reason: string): void;
}>();

const autoSave = ref(false), showTrashModal = ref(false);
const activeLang = ref<'ja' | 'en' | 'zh' | 'all'>('ja'), editBuffer = reactive({ ja: '', en: '', zh: '' });

const stripHtml = (s: string) => (!s ? '' : new DOMParser().parseFromString(s, 'text/html').body.textContent || '');

watch(() => [props.article, props.article?.content?.ja, props.article?.content?.en, props.article?.content?.zh], () => {
  if (props.article) {
    editBuffer.ja = stripHtml(props.article.content?.ja || '');
    editBuffer.en = stripHtml(props.article.content?.en || '');
    editBuffer.zh = stripHtml(props.article.content?.zh || '');
  }
}, { immediate: true, deep: true });
</script>

<template>
  <div v-if="article" class="space-y-2 bg-slate-900/60 border border-slate-800 rounded-xl p-3 flex flex-col h-full">
    <!-- ヘッダー -->
    <div class="flex justify-between items-center border-b border-slate-800 pb-1.5 flex-wrap gap-1.5">
      <h4 class="text-xs font-bold text-slate-200">記事 ID: <span class="font-mono text-blue-400">{{ article.id }}</span></h4>
      <div class="flex items-center gap-1 flex-wrap">
        <label class="text-[10px] text-slate-400 flex items-center gap-1 cursor-pointer select-none mr-1">
          <input type="checkbox" v-model="autoSave" class="rounded border-slate-700 bg-slate-950 text-blue-600 focus:ring-0 w-3 h-3">
          <span>自動保存</span>
        </label>
        <button @click="emit('autoTranslate', autoSave)" :disabled="translating || saving" class="px-2 py-0.5 bg-emerald-600 hover:bg-emerald-500 text-white text-[11px] font-bold rounded disabled:opacity-50 cursor-pointer">🤖 自動翻訳</button>
        <button @click="emit('save', editBuffer.ja, editBuffer.en, editBuffer.zh)" :disabled="saving || translating" class="px-2 py-0.5 bg-blue-600 hover:bg-blue-500 text-white text-[11px] font-bold rounded disabled:opacity-50 cursor-pointer">💾 保存</button>
        <button @click="showTrashModal = true" class="px-2 py-0.5 bg-rose-950/80 hover:bg-rose-900 border border-rose-700/60 text-rose-300 text-[11px] font-bold rounded cursor-pointer" title="ゴミ箱へ移動">🗑️ 削除</button>
      </div>
    </div>

    <!-- URL参照バッジ -->
    <ArticleUrlBadges :domain="article?.source_domain" :original-url="article?.original_url" :wayback-url="article?.wayback_url || article?.source_url" :sotwe-url="article?.sotwe_url" :nitter-url="article?.nitter_url" :twistalker-url="article?.twistalker_url" />

    <!-- 原文プレビュー -->
    <div class="space-y-0.5">
      <label class="text-[10px] text-slate-400 font-semibold flex justify-between">
        <span>原文 ({{ article?.author?.handle ? '@' + article.author.handle : '' }})</span>
      </label>
      <div class="p-1.5 bg-slate-950/80 rounded text-[11px] text-slate-300 font-sans select-text max-h-12 overflow-y-auto border border-slate-800/80 leading-relaxed" v-html="article?.content?.original || ''"></div>
    </div>

    <!-- 言語切り替えタブ -->
    <div class="flex items-center justify-between border-b border-slate-800/80 pb-0.5 text-[11px]">
      <div class="flex gap-1">
        <button v-for="l in [{ id: 'ja', label: '🇯🇵 JA', has: !!editBuffer.ja }, { id: 'en', label: '🇺🇸 EN', has: !!editBuffer.en }, { id: 'zh', label: '🇨🇳 ZH', has: !!editBuffer.zh }, { id: 'all', label: '📋 並列', has: true }]" :key="l.id" @click="activeLang = l.id as any" class="px-2 py-0.5 rounded font-bold cursor-pointer" :class="activeLang === l.id ? 'bg-blue-600/30 text-blue-300 border border-blue-500/40' : 'text-slate-400 hover:text-slate-200'">
          <span>{{ l.label }}</span>
          <span v-if="l.id !== 'all' && !l.has" class="w-1.5 h-1.5 rounded-full bg-amber-500 inline-block ml-0.5"></span>
        </button>
      </div>
    </div>

    <!-- 翻訳入力エリア -->
    <div class="flex-1 flex flex-col min-h-[120px] overflow-y-auto gap-1.5">
      <div v-if="activeLang === 'ja' || activeLang === 'all'" class="flex-1 flex flex-col min-h-[60px]">
        <label class="text-[10px] text-emerald-400 font-semibold mb-0.5">🇯🇵 日本語翻訳 (JA)</label>
        <textarea v-model="editBuffer.ja" placeholder="日本語訳" class="w-full flex-1 min-h-[50px] bg-slate-950 border border-slate-700 rounded p-2 text-xs text-slate-100 font-sans leading-relaxed"></textarea>
      </div>
      <div v-if="activeLang === 'en' || activeLang === 'all'" class="flex-1 flex flex-col min-h-[60px]">
        <label class="text-[10px] text-blue-400 font-semibold mb-0.5">🇺🇸 英語翻訳 (EN)</label>
        <textarea v-model="editBuffer.en" placeholder="英語訳" class="w-full flex-1 min-h-[50px] bg-slate-950 border border-slate-700 rounded p-2 text-xs text-slate-100 font-sans leading-relaxed"></textarea>
      </div>
      <div v-if="activeLang === 'zh' || activeLang === 'all'" class="flex-1 flex flex-col min-h-[60px]">
        <label class="text-[10px] text-purple-400 font-semibold mb-0.5">🇨🇳 中国語翻訳 (ZH)</label>
        <textarea v-model="editBuffer.zh" placeholder="中国語訳" class="w-full flex-1 min-h-[50px] bg-slate-950 border border-slate-700 rounded p-2 text-xs text-slate-100 font-sans leading-relaxed"></textarea>
      </div>
    </div>

    <PostTrashModal :show="showTrashModal" :article-id="article?.id || ''" @close="showTrashModal = false" @confirm="(r) => emit('trash', article?.id || '', r)" />
  </div>
  <div v-else class="h-full flex items-center justify-center border border-slate-800 rounded-xl bg-slate-900/30 text-xs text-slate-500">
    左側の一覧から記事を選択して翻訳を編集してください
  </div>
</template>
