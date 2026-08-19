<script setup lang="ts">
import { ref, reactive, watch, onMounted } from 'vue';
import { models } from '../../../wailsjs/go/models';

const props = defineProps<{
  articles: models.RenderTree[];
  total: number;
  selectedArticle: models.RenderTree | null;
  loading: boolean;
  saving: boolean;
  statusMessage: { success: boolean; message: string } | null;
  searchParams: {
    query: string;
    accountID: string;
    filter: string;
    page: number;
    limit: number;
  };
}>();

const emit = defineEmits<{
  (e: 'search'): void;
  (e: 'select', article: models.RenderTree): void;
  (e: 'saveTranslations', id: string, ja: string, en: string, zh: string): void;
}>();

// 編集中の3言語翻訳バッファ
const editBuffer = reactive({
  ja: '',
  en: '',
  zh: '',
});

// 選択された記事が変わったら編集バッファを同期
watch(
  () => props.selectedArticle,
  (art) => {
    if (art) {
      editBuffer.ja = art.content?.ja || '';
      editBuffer.en = art.content?.en || '';
      editBuffer.zh = art.content?.zh || '';
    } else {
      editBuffer.ja = '';
      editBuffer.en = '';
      editBuffer.zh = '';
    }
  },
  { immediate: true }
);

const handleSearch = () => {
  props.searchParams.page = 1;
  emit('search');
};

const handlePageChange = (newPage: number) => {
  if (newPage < 1) return;
  props.searchParams.page = newPage;
  emit('search');
};

const handleSave = () => {
  if (!props.selectedArticle) return;
  emit('saveTranslations', props.selectedArticle.id, editBuffer.ja, editBuffer.en, editBuffer.zh);
};

// クイックコピー機能
const copyOriginalToJA = () => {
  if (props.selectedArticle) {
    editBuffer.ja = props.selectedArticle.content.original;
  }
};

const copyOriginalToEN = () => {
  if (props.selectedArticle) {
    editBuffer.en = props.selectedArticle.content.original;
  }
};

const copyOriginalToZH = () => {
  if (props.selectedArticle) {
    editBuffer.zh = props.selectedArticle.content.original;
  }
};

const copyToClipboard = async (text: string) => {
  try {
    await navigator.clipboard.writeText(text);
    alert('クリップボードにコピーしました');
  } catch (err) {
    console.error('Copy failed:', err);
  }
};

// キーボードショートカット (Ctrl+S / Cmd+S で保存)
const handleKeyDown = (e: KeyboardEvent) => {
  if ((e.ctrlKey || e.metaKey) && e.key === 's') {
    e.preventDefault();
    handleSave();
  }
};

onMounted(() => {
  emit('search');
});
</script>

<template>
  <div class="space-y-4" @keydown="handleKeyDown">
    <!-- ヘッダー ＆ 検索バー -->
    <div class="space-y-3 bg-slate-900/60 p-4 border border-slate-800 rounded-xl">
      <div class="flex items-center justify-between">
        <div>
          <h3 class="text-base font-bold text-slate-100 flex items-center gap-2">
            <span>🗄️</span> データベース閲覧 ＆ 3言語翻訳エディタ
            <span class="text-[10px] font-mono bg-blue-900/40 text-blue-300 border border-blue-700/50 px-2 py-0.5 rounded">
              articles テーブル連携
            </span>
          </h3>
          <p class="text-xs text-slate-400 mt-0.5">
            保存済みアーカイブ記事の直接全文検索・閲覧および日英中の翻訳テキスト手動微調整が可能です。
          </p>
        </div>
        <div class="text-xs text-slate-400 font-mono">
          該当件数: <span class="text-blue-400 font-bold">{{ total }}</span> 件
        </div>
      </div>

      <!-- 検索・絞り込みフォーム -->
      <div class="grid grid-cols-1 md:grid-cols-12 gap-2 pt-1">
        <!-- 全文検索キーワード -->
        <div class="md:col-span-5">
          <div class="relative">
            <input
              v-model="searchParams.query"
              type="text"
              placeholder="本文・翻訳・記事IDで検索 (Enterで実行)..."
              @keyup.enter="handleSearch"
              class="w-full bg-slate-950 border border-slate-700 rounded-lg pl-8 pr-3 py-1.5 text-xs text-slate-200 placeholder-slate-500 focus:outline-none focus:border-blue-500"
            />
            <span class="absolute left-2.5 top-2 text-slate-500 text-xs">🔍</span>
          </div>
        </div>

        <!-- フィルタ種別 -->
        <div class="md:col-span-3">
          <select
            v-model="searchParams.filter"
            @change="handleSearch"
            class="w-full bg-slate-950 border border-slate-700 rounded-lg px-3 py-1.5 text-xs text-slate-200 focus:outline-none focus:border-blue-500"
          >
            <option value="all">すべての投稿</option>
            <option value="media">🖼️ メディア付きのみ</option>
            <option value="reposts">🔁 リツイートのみ</option>
            <option value="bookmarks">⭐ ブックマークのみ</option>
          </select>
        </div>

        <!-- アクションボタン -->
        <div class="md:col-span-4 flex items-center gap-2">
          <button
            @click="handleSearch"
            :disabled="loading"
            class="flex-1 px-4 py-1.5 bg-blue-600 hover:bg-blue-500 text-white text-xs font-bold rounded-lg transition-colors shadow-md shadow-blue-600/20 flex items-center justify-center gap-1.5 disabled:opacity-50"
          >
            <span :class="{ 'animate-spin': loading }">🔍</span>
            検索
          </button>
        </div>
      </div>
    </div>

    <!-- メイン2カラム領域 -->
    <div class="grid grid-cols-1 md:grid-cols-12 gap-4 items-start">
      <!-- 左ペイン: 記事一覧リスト (5 cols) -->
      <div class="md:col-span-5 border border-slate-800 bg-slate-900/40 rounded-xl overflow-hidden flex flex-col max-h-[600px]">
        <!-- リストヘッダー -->
        <div class="px-4 py-2.5 bg-slate-900/80 border-b border-slate-800 flex items-center justify-between text-xs text-slate-400">
          <span class="font-bold">記事一覧</span>
          <span class="font-mono text-[11px]">ページ {{ searchParams.page }} ({{ articles.length }}件表示)</span>
        </div>

        <!-- 記事リスト -->
        <div class="overflow-y-auto flex-1 divide-y divide-slate-800/60 p-1">
          <div
            v-if="articles.length === 0"
            class="py-12 text-center text-slate-500 text-xs"
          >
            <span v-if="loading">検索中...</span>
            <span v-else>一致する記事が見つかりませんでした</span>
          </div>

          <div
            v-for="art in articles"
            :key="art.id"
            @click="emit('select', art)"
            class="p-3 rounded-lg cursor-pointer transition-all space-y-1.5"
            :class="
              selectedArticle?.id === art.id
                ? 'bg-blue-950/40 border border-blue-600/50 text-white'
                : 'hover:bg-slate-800/40 text-slate-300'
            "
          >
            <!-- ユーザー ＆ 投稿日時 -->
            <div class="flex items-center justify-between text-[11px]">
              <div class="flex items-center gap-1.5 truncate">
                <span class="font-bold truncate">{{ art.author?.display_name || art.author?.handle }}</span>
                <span class="text-slate-500 font-mono">@{{ art.author?.handle }}</span>
              </div>
              <span class="text-[10px] text-slate-500 font-mono whitespace-nowrap">
                {{ art.created_at ? new Date(art.created_at).toLocaleDateString() : '' }}
              </span>
            </div>

            <!-- 本文抜粋 -->
            <p class="text-xs line-clamp-2 text-slate-300 leading-relaxed font-sans">
              {{ art.content?.original }}
            </p>

            <!-- バッジ群 (メディア・翻訳状況) -->
            <div class="flex items-center gap-1.5 pt-0.5">
              <span
                v-if="art.media && art.media.length > 0"
                class="text-[10px] bg-slate-800 text-slate-300 px-1.5 py-0.2 rounded border border-slate-700"
              >
                🖼️ {{ art.media.length }}
              </span>
              <span
                v-if="art.is_liked"
                class="text-[10px] bg-amber-950/40 text-amber-300 px-1.5 py-0.2 rounded border border-amber-800/40"
              >
                ⭐
              </span>

              <div class="flex items-center gap-1 ml-auto text-[9px] font-mono">
                <span
                  class="px-1 py-0.2 rounded"
                  :class="art.content?.ja ? 'bg-emerald-950 text-emerald-400 border border-emerald-800' : 'bg-slate-800 text-slate-600'"
                >
                  JA
                </span>
                <span
                  class="px-1 py-0.2 rounded"
                  :class="art.content?.en ? 'bg-blue-950 text-blue-400 border border-blue-800' : 'bg-slate-800 text-slate-600'"
                >
                  EN
                </span>
                <span
                  class="px-1 py-0.2 rounded"
                  :class="art.content?.zh ? 'bg-purple-950 text-purple-400 border border-purple-800' : 'bg-slate-800 text-slate-600'"
                >
                  ZH
                </span>
              </div>
            </div>
          </div>
        </div>

        <!-- ページネーションコントロール -->
        <div class="p-2.5 bg-slate-900/80 border-t border-slate-800 flex items-center justify-between text-xs">
          <button
            @click="handlePageChange(searchParams.page - 1)"
            :disabled="searchParams.page <= 1 || loading"
            class="px-3 py-1 bg-slate-800 hover:bg-slate-700 text-slate-300 rounded transition-colors disabled:opacity-40"
          >
            ◀ 前へ
          </button>
          <span class="text-slate-400 font-mono text-[11px]">
            ページ {{ searchParams.page }} / {{ Math.max(1, Math.ceil(total / searchParams.limit)) }}
          </span>
          <button
            @click="handlePageChange(searchParams.page + 1)"
            :disabled="searchParams.page * searchParams.limit >= total || loading"
            class="px-3 py-1 bg-slate-800 hover:bg-slate-700 text-slate-300 rounded transition-colors disabled:opacity-40"
          >
            次へ ▶
          </button>
        </div>
      </div>

      <!-- 右ペイン: 記事詳細 ＆ 3言語翻訳エディタ (7 cols) -->
      <div class="md:col-span-7 border border-slate-800 bg-slate-900/50 rounded-xl p-4 space-y-4 max-h-[600px] overflow-y-auto">
        <div v-if="!selectedArticle" class="py-24 text-center text-slate-500 text-xs">
          左の記事一覧から編集したい記事を選択してください。
        </div>

        <div v-else class="space-y-4">
          <!-- 記事ヘッダー -->
          <div class="flex items-start justify-between border-b border-slate-800 pb-3">
            <div>
              <div class="flex items-center gap-2">
                <span class="text-sm font-bold text-white">{{ selectedArticle.author?.display_name }}</span>
                <span class="text-xs text-slate-400 font-mono">@{{ selectedArticle.author?.handle }}</span>
              </div>
              <div class="text-[10px] text-slate-500 font-mono flex items-center gap-2 mt-0.5">
                <span>ID: {{ selectedArticle.id }}</span>
                <span>•</span>
                <span>{{ selectedArticle.created_at ? new Date(selectedArticle.created_at).toLocaleString() : '' }}</span>
              </div>
            </div>

            <!-- 原本 Wayback リンク -->
            <a
              v-if="selectedArticle.source_url"
              :href="selectedArticle.source_url"
              target="_blank"
              class="px-2.5 py-1 bg-slate-800 hover:bg-slate-700 text-blue-400 hover:text-blue-300 border border-slate-700 rounded text-xs flex items-center gap-1 transition-colors"
            >
              <span>🌐</span> 原本 (Wayback)
            </a>
          </div>

          <!-- ステータスフィードバック -->
          <div
            v-if="statusMessage"
            class="p-2.5 rounded-lg text-xs font-semibold flex items-center gap-2"
            :class="
              statusMessage.success
                ? 'bg-emerald-950/60 border border-emerald-500/30 text-emerald-300'
                : 'bg-rose-950/60 border border-rose-500/30 text-rose-300'
            "
          >
            <span>{{ statusMessage.success ? '✅' : '⚠️' }}</span>
            <span>{{ statusMessage.message }}</span>
          </div>

          <!-- 原文表示エリア -->
          <div class="p-3 bg-slate-950/80 border border-slate-800 rounded-lg space-y-1.5">
            <div class="flex items-center justify-between text-xs text-slate-400 font-bold">
              <span class="flex items-center gap-1">
                <span>📄</span> 原文テキスト (Original Full Text)
              </span>
              <button
                @click="copyToClipboard(selectedArticle.content.original)"
                class="text-[11px] text-blue-400 hover:text-blue-300 font-normal"
              >
                コピー
              </button>
            </div>
            <p class="text-xs text-slate-200 leading-relaxed font-sans whitespace-pre-wrap select-text">
              {{ selectedArticle.content.original }}
            </p>
          </div>

          <!-- 添付メディアサムネイル (存在する場合) -->
          <div v-if="selectedArticle.media && selectedArticle.media.length > 0" class="space-y-1.5">
            <span class="text-xs font-bold text-slate-400 flex items-center gap-1">
              <span>🖼️</span> 添付メディア ({{ selectedArticle.media.length }}件)
            </span>
            <div class="grid grid-cols-4 gap-2">
              <div
                v-for="m in selectedArticle.media"
                :key="m.id"
                class="border border-slate-800 rounded-lg overflow-hidden bg-slate-950 relative h-24 flex items-center justify-center p-0.5"
              >
                <img
                  v-if="m.urls?.thumbnail || m.urls?.image || m.urls?.original"
                  :src="m.urls?.thumbnail || m.urls?.image || m.urls?.original"
                  class="max-w-full max-h-full object-contain rounded"
                  loading="lazy"
                />
                <span v-else class="text-[10px] text-slate-500">{{ m.type }}</span>
                <span class="absolute bottom-1 right-1 text-[9px] bg-black/70 px-1 py-0.2 rounded font-mono text-slate-300">
                  {{ m.type }}
                </span>
              </div>
            </div>
          </div>

          <!-- 3言語翻訳エディタ -->
          <div class="space-y-3 pt-2">
            <div class="flex items-center justify-between">
              <h4 class="text-xs font-bold text-slate-200 flex items-center gap-1.5">
                <span>🌐</span> 3言語翻訳手動微調整
              </h4>
              <span class="text-[10px] text-slate-500 font-mono">ショートカット: Ctrl + S で保存</span>
            </div>

            <!-- 1. 日本語 (JA) -->
            <div class="space-y-1">
              <div class="flex items-center justify-between text-xs">
                <label class="text-slate-300 font-semibold flex items-center gap-1">
                  <span>🇯🇵</span> 日本語翻訳 (full_text_ja)
                </label>
                <button
                  @click="copyOriginalToJA"
                  class="text-[11px] text-blue-400 hover:text-blue-300 font-medium"
                >
                  原文からコピー
                </button>
              </div>
              <textarea
                v-model="editBuffer.ja"
                rows="2"
                placeholder="日本語の翻訳テキストを入力..."
                class="w-full bg-slate-950 border border-slate-700 rounded-lg p-2.5 text-xs text-slate-100 placeholder-slate-600 focus:outline-none focus:border-blue-500 leading-relaxed font-sans"
              ></textarea>
            </div>

            <!-- 2. 英語 (EN) -->
            <div class="space-y-1">
              <div class="flex items-center justify-between text-xs">
                <label class="text-slate-300 font-semibold flex items-center gap-1">
                  <span>🇺🇸</span> 英語翻訳 (full_text_en)
                </label>
                <button
                  @click="copyOriginalToEN"
                  class="text-[11px] text-blue-400 hover:text-blue-300 font-medium"
                >
                  原文からコピー
                </button>
              </div>
              <textarea
                v-model="editBuffer.en"
                rows="2"
                placeholder="English translation text..."
                class="w-full bg-slate-950 border border-slate-700 rounded-lg p-2.5 text-xs text-slate-100 placeholder-slate-600 focus:outline-none focus:border-blue-500 leading-relaxed font-sans"
              ></textarea>
            </div>

            <!-- 3. 中国語 (ZH) -->
            <div class="space-y-1">
              <div class="flex items-center justify-between text-xs">
                <label class="text-slate-300 font-semibold flex items-center gap-1">
                  <span>🇨🇳</span> 中国語翻訳 (full_text_zh)
                </label>
                <button
                  @click="copyOriginalToZH"
                  class="text-[11px] text-blue-400 hover:text-blue-300 font-medium"
                >
                  原文からコピー
                </button>
              </div>
              <textarea
                v-model="editBuffer.zh"
                rows="2"
                placeholder="中文翻译文本..."
                class="w-full bg-slate-950 border border-slate-700 rounded-lg p-2.5 text-xs text-slate-100 placeholder-slate-600 focus:outline-none focus:border-blue-500 leading-relaxed font-sans"
              ></textarea>
            </div>

            <!-- 保存ボタン -->
            <div class="pt-2 flex justify-end">
              <button
                @click="handleSave"
                :disabled="saving"
                class="px-6 py-2 bg-blue-600 hover:bg-blue-500 text-white text-xs font-bold rounded-lg transition-colors shadow-lg shadow-blue-600/30 flex items-center gap-2 disabled:opacity-50"
              >
                <span v-if="saving" class="animate-spin">⏳</span>
                <span>💾 翻訳テキストを保存する</span>
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
