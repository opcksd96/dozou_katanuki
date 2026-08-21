<!-- frontend/src/components/admin/database/DatabaseArticleList.vue (100行以下) -->
<script setup lang="ts">
import { models } from '../../../../wailsjs/go/models';
const props = defineProps<{
  articles: models.RenderTree[];
  total: number;
  selectedArticle: models.RenderTree | null;
  loading: boolean;
  searchParams: { query: string; accountID: string; filter: string; page: number; limit: number };
}>();
const emit = defineEmits<{ (e: 'search'): void; (e: 'select', art: models.RenderTree): void }>();

const handleSearch = () => { props.searchParams.page = 1; emit('search'); };
const handlePageChange = (p: number) => { if (p >= 1) { props.searchParams.page = p; emit('search'); } };
</script>

<template>
  <div class="space-y-3 flex flex-col h-full">
    <div class="flex gap-2">
      <input v-model="searchParams.query" @keyup.enter="handleSearch" type="text" placeholder="キーワード・翻訳文を検索..." class="flex-1 bg-slate-950 border border-slate-700 rounded-lg px-3 py-1.5 text-xs text-slate-200" />
      <button @click="handleSearch" :disabled="loading" class="px-3 py-1.5 bg-blue-600 hover:bg-blue-500 text-white text-xs font-bold rounded-lg">{{ loading ? '⏳' : '検索' }}</button>
    </div>
    <div class="text-[11px] text-slate-400 font-mono flex justify-between">
      <span>全 {{ total }} 件</span>
      <span v-if="articles.length > 0">ページ {{ searchParams.page }} / {{ Math.ceil(total / searchParams.limit) || 1 }}</span>
    </div>
    <div class="flex-1 overflow-y-auto border border-slate-800 rounded-xl bg-slate-900/40 divide-y divide-slate-800/60 max-h-[420px]">
      <div v-if="articles.length === 0" class="p-8 text-center text-slate-500 text-xs">{{ loading ? '読込中...' : '記事がありません' }}</div>
      <div v-for="art in articles" :key="art.id" @click="$emit('select', art)" class="p-3 hover:bg-slate-800/40 cursor-pointer transition-colors" :class="{ 'bg-blue-950/30 border-l-2 border-blue-500': selectedArticle?.id === art.id }">
        <div class="flex items-center justify-between text-[11px] text-slate-400 mb-1">
          <span class="font-bold text-slate-200 font-mono">@{{ art.author.handle }}</span>
          <span>{{ new Date(art.created_at).toLocaleDateString() }}</span>
        </div>
        <div class="text-xs text-slate-300 line-clamp-2">{{ art.content.original }}</div>
        <div class="flex gap-1.5 mt-1.5 text-[10px] font-mono">
          <span :class="art.content.ja ? 'text-emerald-400' : 'text-slate-600'">JA:{{ art.content.ja ? '✓' : '-' }}</span>
          <span :class="art.content.en ? 'text-blue-400' : 'text-slate-600'">EN:{{ art.content.en ? '✓' : '-' }}</span>
          <span :class="art.content.zh ? 'text-purple-400' : 'text-slate-600'">ZH:{{ art.content.zh ? '✓' : '-' }}</span>
        </div>
      </div>
    </div>
    <div class="flex justify-between items-center text-xs">
      <button @click="handlePageChange(searchParams.page - 1)" :disabled="searchParams.page <= 1" class="px-2.5 py-1 bg-slate-800 rounded disabled:opacity-40">前へ</button>
      <button @click="handlePageChange(searchParams.page + 1)" :disabled="searchParams.page * searchParams.limit >= total" class="px-2.5 py-1 bg-slate-800 rounded disabled:opacity-40">次へ</button>
    </div>
  </div>
</template>
