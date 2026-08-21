<!-- frontend/src/components/admin/DatabaseView.vue (100行以下) -->
<script setup lang="ts">
import { models } from '../../../wailsjs/go/models';
import DatabaseArticleList from './database/DatabaseArticleList.vue';
import DatabaseTranslationEditor from './database/DatabaseTranslationEditor.vue';

const props = defineProps<{
  articles: models.RenderTree[];
  total: number;
  selectedArticle: models.RenderTree | null;
  loading: boolean;
  saving: boolean;
  statusMessage: { success: boolean; message: string } | null;
  searchParams: { query: string; accountID: string; filter: string; page: number; limit: number };
}>();

const emit = defineEmits<{
  (e: 'search'): void;
  (e: 'select', article: models.RenderTree): void;
  (e: 'saveTranslations', id: string, ja: string, en: string, zh: string): void;
}>();
</script>

<template>
  <div class="space-y-4">
    <div class="flex items-center justify-between">
      <div>
        <h3 class="text-base font-bold text-slate-100 flex items-center gap-2">
          <span>🗄️</span> データベース閲覧 ＆ 翻訳編集
        </h3>
        <p class="text-xs text-slate-400 mt-0.5">保存された記事の検索・抽出および 3言語（日・英・中）の翻訳文の直接編集が可能です。</p>
      </div>
    </div>

    <div v-if="statusMessage" class="p-2.5 rounded-lg text-xs font-semibold" :class="statusMessage.success ? 'bg-emerald-950/60 border border-emerald-500/30 text-emerald-300' : 'bg-rose-950/60 border border-rose-500/30 text-rose-300'">
      {{ statusMessage.success ? '✅' : '⚠️' }} {{ statusMessage.message }}
    </div>

    <div class="grid grid-cols-1 md:grid-cols-2 gap-6 min-h-[500px]">
      <DatabaseArticleList
        :articles="articles"
        :total="total"
        :selected-article="selectedArticle"
        :loading="loading"
        :search-params="searchParams"
        @search="$emit('search')"
        @select="(art) => $emit('select', art)"
      />

      <DatabaseTranslationEditor
        :article="selectedArticle"
        :saving="saving"
        @save="(ja, en, zh) => { if (selectedArticle) $emit('saveTranslations', selectedArticle.id, ja, en, zh); }"
      />
    </div>
  </div>
</template>
