<script setup lang="ts">
import { onMounted } from 'vue';

const props = defineProps<{
  report: any | null;
  loading: boolean;
  purgingFiles: boolean;
  purgingDB: boolean;
  statusMessage: { success: boolean; message: string } | null;
}>();

const emit = defineEmits<{
  (e: 'runAudit', purgeFiles?: boolean, purgeDB?: boolean): void;
  (e: 'purgeOrphanFiles', paths?: string[]): void;
  (e: 'purgeOrphanDBMedia', mediaIDs?: string[]): void;
}>();

const formatBytes = (bytes: number): string => {
  if (bytes === 0) return '0 B';
  const k = 1024;
  const sizes = ['B', 'KB', 'MB', 'GB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
};

onMounted(() => {
  if (!props.report) {
    emit('runAudit', false, false);
  }
});
</script>

<template>
  <div class="space-y-4">
    <!-- ヘッダー ＆ アクションバー -->
    <div class="bg-slate-900/60 p-4 border border-slate-800 rounded-xl space-y-3">
      <div class="flex flex-col md:flex-row md:items-center justify-between gap-3">
        <div>
          <h3 class="text-base font-bold text-slate-100 flex items-center gap-2">
            <span>🩺</span> SQLite3 整合性監査 ＆ 孤立ファイルパージ
            <span class="text-[10px] font-mono bg-emerald-950/80 text-emerald-400 border border-emerald-700/50 px-2 py-0.5 rounded">
              SPEC-AUDIT-001
            </span>
          </h3>
          <p class="text-xs text-slate-400 mt-0.5">
            B-Tree・インデックス破損、外部キー整合性、および DB と Stash / Blobs ストレージ間の孤立ファイルを完全監査します。
          </p>
        </div>

        <!-- 監査実行ボタン群 -->
        <div class="flex items-center gap-2">
          <button
            @click="emit('runAudit', false, false)"
            :disabled="loading"
            class="px-4 py-2 bg-blue-600 hover:bg-blue-500 text-white text-xs font-bold rounded-lg transition-all shadow-lg shadow-blue-600/20 flex items-center gap-1.5 disabled:opacity-50"
          >
            <span :class="{ 'animate-spin': loading }">🔍</span>
            整合性監査を実行
          </button>

          <button
            @click="emit('runAudit', true, false)"
            :disabled="loading || purgingFiles"
            class="px-3.5 py-2 bg-amber-600/90 hover:bg-amber-500 text-white text-xs font-bold rounded-lg transition-all shadow-lg shadow-amber-600/20 flex items-center gap-1.5 disabled:opacity-50"
          >
            <span :class="{ 'animate-spin': loading || purgingFiles }">🧹</span>
            自動パージ付き監査
          </button>
        </div>
      </div>

      <!-- ステータスフィードバックメッセージ -->
      <div
        v-if="statusMessage"
        class="p-2.5 rounded-lg text-xs font-semibold flex items-center gap-2 transition-all"
        :class="
          statusMessage.success
            ? 'bg-emerald-950/70 border border-emerald-500/40 text-emerald-300'
            : 'bg-rose-950/70 border border-rose-500/40 text-rose-300'
        "
      >
        <span>{{ statusMessage.success ? '✅' : '⚠️' }}</span>
        <span>{{ statusMessage.message }}</span>
      </div>

      <!-- 監査サマリーバー -->
      <div v-if="report" class="grid grid-cols-2 md:grid-cols-4 gap-2 pt-1 font-mono text-xs">
        <div class="bg-slate-950/80 p-2.5 rounded-lg border border-slate-800">
          <div class="text-[10px] text-slate-500">B-Tree / インデックス</div>
          <div class="font-bold flex items-center gap-1.5 mt-0.5" :class="report.integrity_ok ? 'text-emerald-400' : 'text-rose-400'">
            <span>{{ report.integrity_ok ? '🛡️ 健全 (OK)' : '❌ 破損検知' }}</span>
          </div>
        </div>

        <div class="bg-slate-950/80 p-2.5 rounded-lg border border-slate-800">
          <div class="text-[10px] text-slate-500">外部キー整合性 (FK)</div>
          <div class="font-bold flex items-center gap-1.5 mt-0.5" :class="report.foreign_key_ok ? 'text-emerald-400' : 'text-amber-400'">
            <span>{{ report.foreign_key_ok ? '🔗 正常 (0件)' : `⚠️ 違反 ${report.foreign_key_errors?.length || 0}件` }}</span>
          </div>
        </div>

        <div class="bg-slate-950/80 p-2.5 rounded-lg border border-slate-800">
          <div class="text-[10px] text-slate-500">孤立 DB メディア</div>
          <div class="font-bold flex items-center gap-1.5 mt-0.5" :class="(report.orphan_db_media?.length || 0) === 0 ? 'text-slate-300' : 'text-amber-400'">
            <span>📑 {{ report.orphan_db_media?.length || 0 }} 件</span>
          </div>
        </div>

        <div class="bg-slate-950/80 p-2.5 rounded-lg border border-slate-800">
          <div class="text-[10px] text-slate-500">ストレージ孤立ファイル</div>
          <div class="font-bold flex items-center gap-1.5 mt-0.5" :class="(report.orphan_files?.length || 0) === 0 ? 'text-slate-300' : 'text-amber-400'">
            <span>📁 {{ report.orphan_files?.length || 0 }} 件</span>
          </div>
        </div>
      </div>
    </div>

    <!-- ローディング表示 -->
    <div v-if="loading && !report" class="py-20 text-center text-slate-400 space-y-3">
      <div class="text-3xl animate-bounce">🩺</div>
      <p class="text-sm font-semibold">SQLite3 データベースおよびストレージの整合性を監査中...</p>
    </div>

    <!-- メインコンテンツ領域 -->
    <div v-else-if="report" class="space-y-4">
      <!-- 1. PRAGMA integrity_check 結果詳細 -->
      <div class="bg-slate-900/40 border border-slate-800 rounded-xl p-4 space-y-2">
        <div class="flex items-center justify-between">
          <h4 class="text-xs font-bold text-slate-200 flex items-center gap-1.5">
            <span>🛡️</span> PRAGMA integrity_check 診断結果
          </h4>
          <span
            class="text-[10px] font-mono px-2 py-0.5 rounded border"
            :class="
              report.integrity_ok
                ? 'bg-emerald-950 text-emerald-400 border-emerald-800'
                : 'bg-rose-950 text-rose-400 border-rose-800'
            "
          >
            {{ report.integrity_ok ? 'PASSED (健全)' : 'FAILED (破損)' }}
          </span>
        </div>
        <div class="bg-slate-950 p-2.5 rounded-lg font-mono text-xs border border-slate-800/80 max-h-24 overflow-y-auto">
          <div v-for="(msg, idx) in report.integrity_errors" :key="idx" :class="report.integrity_ok ? 'text-slate-400' : 'text-rose-400 font-semibold'">
            {{ msg }}
          </div>
        </div>
      </div>

      <!-- 2. PRAGMA foreign_key_check 結果詳細 (存在する場合) -->
      <div v-if="!report.foreign_key_ok && report.foreign_key_errors?.length > 0" class="bg-slate-900/40 border border-amber-800/50 rounded-xl p-4 space-y-2">
        <div class="flex items-center justify-between">
          <h4 class="text-xs font-bold text-amber-300 flex items-center gap-1.5">
            <span>⚠️</span> PRAGMA foreign_key_check 外部キー制約違反 ({{ report.foreign_key_errors.length }}件)
          </h4>
        </div>
        <div class="bg-slate-950 rounded-lg border border-slate-800 overflow-x-auto">
          <table class="w-full text-left text-xs font-mono">
            <thead class="bg-slate-900/80 text-slate-400 border-b border-slate-800">
              <tr>
                <th class="p-2">テーブル</th>
                <th class="p-2">Row ID</th>
                <th class="p-2">参照先テーブル</th>
                <th class="p-2">FK ID</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-800/60 text-slate-300">
              <tr v-for="(fk, idx) in report.foreign_key_errors" :key="idx">
                <td class="p-2 font-bold text-amber-400">{{ fk.table }}</td>
                <td class="p-2">{{ fk.row_id }}</td>
                <td class="p-2 text-slate-400">{{ fk.parent_table }}</td>
                <td class="p-2">{{ fk.fk_id }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- 3. 孤立 DB メディアレコード一覧 -->
      <div class="bg-slate-900/40 border border-slate-800 rounded-xl p-4 space-y-3">
        <div class="flex items-center justify-between">
          <div>
            <h4 class="text-xs font-bold text-slate-200 flex items-center gap-1.5">
              <span>📑</span> 孤立 DB メディアレコード ({{ report.orphan_db_media?.length || 0 }}件)
            </h4>
            <p class="text-[11px] text-slate-400">親記事が存在しない、または実体ファイルが消失している DB レコードです。</p>
          </div>

          <button
            v-if="report.orphan_db_media?.length > 0"
            @click="emit('purgeOrphanDBMedia')"
            :disabled="purgingDB"
            class="px-3 py-1.5 bg-rose-600/80 hover:bg-rose-500 text-white text-xs font-bold rounded-lg transition-colors flex items-center gap-1 disabled:opacity-50"
          >
            <span :class="{ 'animate-spin': purgingDB }">🗑️</span>
            孤立レコードを一括削除
          </button>
        </div>

        <div v-if="(report.orphan_db_media?.length || 0) === 0" class="py-6 text-center text-slate-500 text-xs font-mono">
          ✅ 孤立した DB メディアレコードはありません (完全整合)
        </div>

        <div v-else class="bg-slate-950 rounded-lg border border-slate-800 overflow-x-auto max-h-48 overflow-y-auto">
          <table class="w-full text-left text-xs font-mono">
            <thead class="bg-slate-900/80 text-slate-400 border-b border-slate-800 sticky top-0">
              <tr>
                <th class="p-2">Media ID</th>
                <th class="p-2">Article ID</th>
                <th class="p-2">種別</th>
                <th class="p-2">理由</th>
                <th class="p-2 text-right">操作</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-800/60 text-slate-300">
              <tr v-for="m in report.orphan_db_media" :key="m.media_id">
                <td class="p-2 font-bold text-blue-400 truncate max-w-[140px]">{{ m.media_id }}</td>
                <td class="p-2 text-slate-400 truncate max-w-[120px]">{{ m.article_id }}</td>
                <td class="p-2">
                  <span class="px-1.5 py-0.5 rounded bg-slate-800 text-[10px]">{{ m.type }}</span>
                </td>
                <td class="p-2 text-amber-400/90 text-[11px]">{{ m.reason }}</td>
                <td class="p-2 text-right">
                  <button
                    @click="emit('purgeOrphanDBMedia', [m.media_id])"
                    class="px-2 py-0.5 bg-slate-800 hover:bg-rose-600 text-slate-300 hover:text-white rounded transition-colors text-[11px]"
                  >
                    削除
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- 4. ストレージ孤立ファイル一覧 ＆ ごみ箱退避 -->
      <div class="bg-slate-900/40 border border-slate-800 rounded-xl p-4 space-y-3">
        <div class="flex items-center justify-between">
          <div>
            <h4 class="text-xs font-bold text-slate-200 flex items-center gap-1.5">
              <span>📁</span> ストレージ孤立ファイル ({{ report.orphan_files?.length || 0 }}件)
            </h4>
            <p class="text-[11px] text-slate-400">
              <code>stash/scenes/</code>, <code>stash/images/</code>, <code>blobs/</code> 内で DB に紐付いていないゾンビファイルです。
            </p>
          </div>

          <button
            v-if="report.orphan_files?.length > 0"
            @click="emit('purgeOrphanFiles')"
            :disabled="purgingFiles"
            class="px-3.5 py-1.5 bg-amber-600 hover:bg-amber-500 text-white text-xs font-bold rounded-lg transition-colors shadow-lg shadow-amber-600/20 flex items-center gap-1.5 disabled:opacity-50"
          >
            <span :class="{ 'animate-spin': purgingFiles }">🗑️</span>
            ごみ箱へ一括退避パージ
          </button>
        </div>

        <div v-if="(report.orphan_files?.length || 0) === 0" class="py-6 text-center text-slate-500 text-xs font-mono">
          ✅ ストレージ上に未紐付けの孤立ファイルはありません (クリーン)
        </div>

        <div v-else class="bg-slate-950 rounded-lg border border-slate-800 overflow-x-auto max-h-56 overflow-y-auto">
          <table class="w-full text-left text-xs font-mono">
            <thead class="bg-slate-900/80 text-slate-400 border-b border-slate-800 sticky top-0">
              <tr>
                <th class="p-2">ファイル名</th>
                <th class="p-2">カテゴリ</th>
                <th class="p-2">サイズ</th>
                <th class="p-2">パス</th>
                <th class="p-2 text-right">操作</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-800/60 text-slate-300">
              <tr v-for="f in report.orphan_files" :key="f.path">
                <td class="p-2 font-bold text-amber-400 truncate max-w-[160px]">{{ f.file_name }}</td>
                <td class="p-2">
                  <span class="px-1.5 py-0.5 rounded bg-slate-800 text-[10px]">{{ f.category }}</span>
                </td>
                <td class="p-2 text-slate-400">{{ formatBytes(f.file_size) }}</td>
                <td class="p-2 text-slate-500 truncate max-w-[200px]" :title="f.path">{{ f.path }}</td>
                <td class="p-2 text-right">
                  <button
                    @click="emit('purgeOrphanFiles', [f.path])"
                    class="px-2 py-0.5 bg-slate-800 hover:bg-amber-600 text-slate-300 hover:text-white rounded transition-colors text-[11px]"
                  >
                    ごみ箱退避
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
  </div>
</template>
