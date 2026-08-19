<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch } from 'vue';
import { useAdmin } from '../../composables/useAdmin';
import JobController from './JobController.vue';
import ConfigPortal from './ConfigPortal.vue';
import StashStatusView from './StashStatusView.vue';
import WhitelistView from './WhitelistView.vue';
import DatabaseView from './DatabaseView.vue';

const props = defineProps<{
  isOpen: boolean;
}>();

const emit = defineEmits<{
  (e: 'close'): void;
}>();

const activeTab = ref<'jobs' | 'config' | 'stash' | 'whitelist' | 'db'>('jobs');

const {
  activeJob,
  jobList,
  logs,
  config,
  loadingJobs,
  loadingConfig,
  savingConfig,
  actionLoading,
  saveStatus,
  salvageForm,
  importForm,
  fetchJobs,
  loadConfig,
  saveConfig,
  startSalvage,
  startImport,
  cancelJob,
  clearLogs,
  setupEventListeners,
  cleanupEventListeners,
  // Whitelist
  whitelistList,
  loadingWhitelist,
  whitelistStatus,
  fetchWhitelists,
  addWhitelist,
  updateWhitelist,
  deleteWhitelist,
  toggleWhitelist,
  // Database
  dbArticles,
  totalArticles,
  selectedArticle,
  loadingArticles,
  savingTranslation,
  translationStatus,
  dbSearchParams,
  searchArticles,
  selectArticle,
  saveArticleTranslations,
} = useAdmin();

// モーダルオープン時にデータ取得とイベントリスナー設定
watch(
  () => props.isOpen,
  (val) => {
    if (val) {
      setupEventListeners();
      fetchJobs();
      loadConfig();
      fetchWhitelists();
      searchArticles();
    } else {
      cleanupEventListeners();
    }
  },
  { immediate: true }
);

// タブ切り替え時のデータ最新化
watch(
  () => activeTab.value,
  (tab) => {
    if (!props.isOpen) return;
    if (tab === 'whitelist') {
      fetchWhitelists();
    } else if (tab === 'db') {
      searchArticles();
    }
  }
);

const handleKeyDown = (e: KeyboardEvent) => {
  if (e.key === 'Escape' && props.isOpen) {
    emit('close');
  }
};

onMounted(() => {
  window.addEventListener('keydown', handleKeyDown);
});

onUnmounted(() => {
  window.removeEventListener('keydown', handleKeyDown);
  cleanupEventListeners();
});
</script>

<template>
  <Teleport to="body">
    <div
      v-if="isOpen"
      class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/90 transition-opacity"
      @click.self="emit('close')"
    >
      <div
        class="bg-slate-950 border border-slate-800 rounded-2xl w-full max-w-5xl max-h-[92vh] flex flex-col shadow-2xl overflow-hidden"
      >
        <!-- モーダルヘッダー -->
        <div class="px-6 py-4 border-b border-slate-800 flex items-center justify-between bg-slate-900/50">
          <div class="flex items-center gap-3">
            <div class="w-8 h-8 rounded-lg bg-blue-600/20 border border-blue-500/30 flex items-center justify-center text-sm">
              🛠️
            </div>
            <div>
              <h2 class="text-base font-bold text-white flex items-center gap-2">
                Admin Board ＆ 管理コンソール
                <span class="text-[10px] font-mono font-normal bg-slate-800 text-slate-400 px-2 py-0.5 rounded border border-slate-700">
                  SPEC-ADMINBOARD-001
                </span>
              </h2>
              <p class="text-xs text-slate-400">システム設定・ジョブ監視・Whitelist統治・DB翻訳エディタ</p>
            </div>
          </div>
          <button
            @click="emit('close')"
            class="w-8 h-8 rounded-lg bg-slate-800/80 hover:bg-slate-700 text-slate-400 hover:text-slate-200 flex items-center justify-center transition-colors text-sm"
          >
            ✕
          </button>
        </div>

        <!-- タブバー -->
        <div class="px-6 border-b border-slate-800 bg-slate-900/30 flex items-center gap-2 overflow-x-auto text-xs">
          <button
            @click="activeTab = 'jobs'"
            class="py-3 px-3.5 border-b-2 font-semibold transition-all flex items-center gap-1.5 whitespace-nowrap"
            :class="
              activeTab === 'jobs'
                ? 'border-blue-500 text-blue-400 bg-blue-500/5'
                : 'border-transparent text-slate-400 hover:text-slate-200'
            "
          >
            <span>🚀</span> ジョブ ＆ ログ
            <span
              v-if="activeJob && activeJob.status === 'running'"
              class="w-2 h-2 rounded-full bg-blue-400 animate-pulse ml-1"
            ></span>
          </button>

          <button
            @click="activeTab = 'config'"
            class="py-3 px-3.5 border-b-2 font-semibold transition-all flex items-center gap-1.5 whitespace-nowrap"
            :class="
              activeTab === 'config'
                ? 'border-blue-500 text-blue-400 bg-blue-500/5'
                : 'border-transparent text-slate-400 hover:text-slate-200'
            "
          >
            <span>⚙️</span> システム設定 (Config)
          </button>

          <button
            @click="activeTab = 'stash'"
            class="py-3 px-3.5 border-b-2 font-semibold transition-all flex items-center gap-1.5 whitespace-nowrap"
            :class="
              activeTab === 'stash'
                ? 'border-purple-500 text-purple-400 bg-purple-500/5'
                : 'border-transparent text-slate-400 hover:text-slate-200'
            "
          >
            <span>📦</span> Stash 連携
          </button>

          <!-- 実稼働 Whitelist 管理 -->
          <button
            @click="activeTab = 'whitelist'"
            class="py-3 px-3.5 border-b-2 font-semibold transition-all flex items-center gap-1.5 whitespace-nowrap"
            :class="
              activeTab === 'whitelist'
                ? 'border-blue-500 text-blue-400 bg-blue-500/5'
                : 'border-transparent text-slate-400 hover:text-slate-200'
            "
          >
            <span>📋</span> Whitelist 管理
            <span
              v-if="whitelistList.length > 0"
              class="text-[10px] bg-slate-800 text-slate-300 px-1.5 py-0.2 rounded-full border border-slate-700 font-mono"
            >
              {{ whitelistList.length }}
            </span>
          </button>

          <!-- 実稼働 データベース閲覧 ＆ 翻訳エディタ -->
          <button
            @click="activeTab = 'db'"
            class="py-3 px-3.5 border-b-2 font-semibold transition-all flex items-center gap-1.5 whitespace-nowrap"
            :class="
              activeTab === 'db'
                ? 'border-blue-500 text-blue-400 bg-blue-500/5'
                : 'border-transparent text-slate-400 hover:text-slate-200'
            "
          >
            <span>🗄️</span> データベース閲覧 ＆ 翻訳
          </button>
        </div>

        <!-- モーダルボディ（スクロールエリア） -->
        <div class="p-6 overflow-y-auto flex-1 space-y-6">
          <!-- 1. ジョブ＆ログ -->
          <div v-show="activeTab === 'jobs'">
            <JobController
              :activeJob="activeJob"
              :jobList="jobList"
              :logs="logs"
              :salvageForm="salvageForm"
              :importForm="importForm"
              :actionLoading="actionLoading"
              :loadingJobs="loadingJobs"
              @startSalvage="startSalvage"
              @startImport="startImport"
              @cancelJob="cancelJob"
              @fetchJobs="fetchJobs"
              @clearLogs="clearLogs"
            />
          </div>

          <!-- 2. システム設定 -->
          <div v-show="activeTab === 'config'">
            <ConfigPortal
              :config="config"
              :loadingConfig="loadingConfig"
              :savingConfig="savingConfig"
              :saveStatus="saveStatus"
              @saveConfig="saveConfig"
              @loadConfig="loadConfig"
            />
          </div>

          <!-- 3. Stash連携 -->
          <div v-show="activeTab === 'stash'">
            <StashStatusView :config="config" />
          </div>

          <!-- 4. Whitelist 管理（実稼働） -->
          <div v-show="activeTab === 'whitelist'">
            <WhitelistView
              :whitelistList="whitelistList"
              :loading="loadingWhitelist"
              :statusMessage="whitelistStatus"
              @fetch="fetchWhitelists"
              @add="addWhitelist"
              @update="updateWhitelist"
              @delete="deleteWhitelist"
              @toggle="toggleWhitelist"
            />
          </div>

          <!-- 5. データベース閲覧 ＆ 3言語翻訳エディタ（実稼働） -->
          <div v-show="activeTab === 'db'">
            <DatabaseView
              :articles="dbArticles"
              :total="totalArticles"
              :selectedArticle="selectedArticle"
              :loading="loadingArticles"
              :saving="savingTranslation"
              :statusMessage="translationStatus"
              :searchParams="dbSearchParams"
              @search="searchArticles"
              @select="selectArticle"
              @saveTranslations="saveArticleTranslations"
            />
          </div>
        </div>

        <!-- モーダルフッター -->
        <div class="px-6 py-3 border-t border-slate-800 bg-slate-900/40 flex items-center justify-between text-xs text-slate-500 font-mono">
          <span>dozou_katanuki Admin v1.0 • SPEC-ADMINBOARD-001</span>
          <button
            @click="emit('close')"
            class="px-4 py-1.5 bg-slate-800 hover:bg-slate-700 text-slate-300 rounded-lg transition-colors font-semibold"
          >
            閉じる (ESC)
          </button>
        </div>
      </div>
    </div>
  </Teleport>
</template>

