import { ref, reactive, onMounted, onUnmounted } from 'vue';
import { models } from '../../wailsjs/go/models';
import * as WailsApp from '../../wailsjs/go/main/App';
import { EventsOn, EventsOff } from '../../wailsjs/runtime/runtime';

// Wails RPC の安全な呼び出しヘルパー
const getApp = () => (window as any)?.go?.main?.App;

// デフォルト設定モック（ブラウザ単体テスト用）
const defaultMockConfig = {
  system: { language: 'ja', default_framework: 'twitter', env: 'production' },
  storage: { db_path: './archive.db', local_media_dir: './blobs', stash_dir: './stash', dumps_dir: './dumps', stash_enabled: true },
  network: { stash_port: 9999, frontend_port: 5173, internal_bind_address: '127.0.0.1' },
  scheduler: { poll_interval_sec: 60, backup_interval_hours: 24, max_backup_generations: 7 },
  appearance: { font_family_ja: 'Hiragino Sans, Meiryo, sans-serif', font_family_en: 'Nunito, sans-serif', font_family_zh: 'Microsoft YaHei, SimHei, sans-serif' },
};

let localMockConfig = JSON.parse(JSON.stringify(defaultMockConfig));

const callGetConfig = async (): Promise<any> => {
  const app = getApp();
  if (app && typeof app.GetConfig === 'function') {
    return await app.GetConfig();
  }
  return localMockConfig;
};

const callSaveConfig = async (cfg: any): Promise<any> => {
  const app = getApp();
  if (app && typeof app.SaveConfig === 'function') {
    return await app.SaveConfig(cfg);
  }
  localMockConfig = JSON.parse(JSON.stringify(cfg));
  return true;
};

const callGetActiveJob = async (): Promise<any> => {
  const app = getApp();
  if (app && typeof app.GetActiveJob === 'function') {
    return await app.GetActiveJob();
  }
  return null;
};

const callListJobs = async (): Promise<any> => {
  const app = getApp();
  if (app && typeof app.ListJobs === 'function') {
    return await app.ListJobs();
  }
  return [];
};

const callStartSalvageJob = async (platform: string, account: string, limit: number): Promise<any> => {
  const app = getApp();
  if (app && typeof app.StartSalvageJob === 'function') {
    return await app.StartSalvageJob(platform, account, limit);
  }
  return {
    id: `mock-salvage-${Date.now()}`,
    type: 'salvage',
    status: 'running',
    current: 10,
    total: limit,
    percentage: (10 / limit) * 100,
    started_at: new Date().toISOString(),
    logs: ['[Mock] Starting salvage job for ' + account, '[Mock] Fetching tweets...'],
  };
};

const callStartManualImportJob = async (warcPath: string, offline: boolean): Promise<any> => {
  const app = getApp();
  if (app && typeof app.StartManualImportJob === 'function') {
    return await app.StartManualImportJob(warcPath, offline);
  }
  return {
    id: `mock-import-${Date.now()}`,
    type: 'import',
    status: 'running',
    current: 5,
    total: 20,
    percentage: 25.0,
    started_at: new Date().toISOString(),
    logs: ['[Mock] Starting WARC import: ' + warcPath],
  };
};

const callCancelJob = async (jobId: string): Promise<any> => {
  const app = getApp();
  if (app && typeof app.CancelJob === 'function') {
    return await app.CancelJob(jobId);
  }
  return true;
};

const callTriggerRestore = async (dumpsDir: string, resetDB: boolean): Promise<any> => {
  const app = getApp();
  if (app && typeof app.TriggerRestore === 'function') {
    return await app.TriggerRestore(dumpsDir, resetDB);
  }
  return {
    id: `mock-restore-${Date.now()}`,
    type: 'restore',
    status: 'running',
    current: 1,
    total: 10,
    percentage: 10.0,
    started_at: new Date().toISOString(),
    logs: ['[Mock] Starting offline recovery restore from ' + (dumpsDir || './backups/dumps')],
  };
};

// --- Whitelist RPC & Mock ---
let localMockWhitelists = [
  { id: 1, type: 'account', value: 'mashu_dev', is_active: true },
  { id: 2, type: 'keyword', value: 'NES_APU', is_active: true },
  { id: 3, type: 'account', value: 'famicom_history', is_active: false },
];

const callGetWhitelists = async (): Promise<any[]> => {
  const app = getApp();
  if (app && typeof app.GetWhitelists === 'function') {
    return await app.GetWhitelists();
  }
  return [...localMockWhitelists];
};

const callAddWhitelist = async (type: string, value: string): Promise<any> => {
  const app = getApp();
  if (app && typeof app.AddWhitelist === 'function') {
    return await app.AddWhitelist(type, value);
  }
  const newItem = {
    id: localMockWhitelists.length > 0 ? Math.max(...localMockWhitelists.map((w) => w.id)) + 1 : 1,
    type,
    value,
    is_active: true,
  };
  localMockWhitelists.push(newItem);
  return newItem;
};

const callUpdateWhitelist = async (id: number, type: string, value: string, isActive: boolean): Promise<any> => {
  const app = getApp();
  if (app && typeof app.UpdateWhitelist === 'function') {
    return await app.UpdateWhitelist(id, type, value, isActive);
  }
  const idx = localMockWhitelists.findIndex((w) => w.id === id);
  if (idx !== -1) {
    localMockWhitelists[idx] = { id, type, value, is_active: isActive };
  }
  return true;
};

const callDeleteWhitelist = async (id: number): Promise<any> => {
  const app = getApp();
  if (app && typeof app.DeleteWhitelist === 'function') {
    return await app.DeleteWhitelist(id);
  }
  localMockWhitelists = localMockWhitelists.filter((w) => w.id !== id);
  return true;
};

const callToggleWhitelist = async (id: number): Promise<any> => {
  const app = getApp();
  if (app && typeof app.ToggleWhitelist === 'function') {
    return await app.ToggleWhitelist(id);
  }
  const item = localMockWhitelists.find((w) => w.id === id);
  if (item) {
    item.is_active = !item.is_active;
  }
  return true;
};

// --- Database & Article Search / Translation RPC & Mock ---
let localMockArticles: any[] = [
  {
    id: 'mock_art_1',
    conversation_id: 'mock_art_1',
    created_at: new Date().toISOString(),
    content: {
      original: 'Famicom 2A03 APU register $4000 pulse channel configuration test note.',
      ja: 'ファミコン 2A03 APU レジスタ $4000 矩形波チャンネル設定テスト記録。',
      en: 'Famicom 2A03 APU register $4000 pulse channel configuration test note.',
      zh: '红白机 2A03 APU 寄存器 $4000 脉冲通道配置测试记录。',
    },
    author: {
      numeric_id: '1001',
      handle: 'senpai_apu',
      display_name: '先輩マスター',
      avatar_url: '',
    },
    media: [],
    metrics: { replies: 2, retweets: 5, likes: 18 },
    is_liked: true,
    is_pinned: false,
    source_url: 'https://web.archive.org/web/mock_art_1',
  },
  {
    id: 'mock_art_2',
    conversation_id: 'mock_art_2',
    created_at: new Date(Date.now() - 3600000).toISOString(),
    content: {
      original: 'Dozou Katanuki v4.0 Wails & Stash integration is progressing smoothly!',
      ja: '土蔵・型抜き v4.0 Wails ＆ Stash 統合が極めて順調に進んでいます！',
      en: 'Dozou Katanuki v4.0 Wails & Stash integration is progressing smoothly!',
      zh: '',
    },
    author: {
      numeric_id: '1002',
      handle: 'mashu_dev',
      display_name: 'マシュ・キリエライト',
      avatar_url: '',
    },
    media: [],
    metrics: { replies: 0, retweets: 3, likes: 12 },
    is_liked: false,
    is_pinned: false,
    source_url: 'https://web.archive.org/web/mock_art_2',
  },
];

const callSearchArticles = async (query: string, accountId: string, filter: string, limit: number, offset: number): Promise<{ items: any[]; total: number }> => {
  const app = getApp();
  if (app && typeof app.SearchArticles === 'function') {
    return await app.SearchArticles(query, accountId, filter, limit, offset);
  }
  let filtered = [...localMockArticles];
  if (query) {
    const q = query.toLowerCase();
    filtered = filtered.filter(
      (a) =>
        a.id.toLowerCase().includes(q) ||
        a.content.original.toLowerCase().includes(q) ||
        (a.content.ja && a.content.ja.toLowerCase().includes(q)) ||
        (a.content.en && a.content.en.toLowerCase().includes(q)) ||
        (a.content.zh && a.content.zh.toLowerCase().includes(q))
    );
  }
  if (accountId && accountId !== 'all') {
    filtered = filtered.filter((a) => a.author.numeric_id === accountId || a.author.handle === accountId);
  }
  const total = filtered.length;
  const items = filtered.slice(offset, offset + limit);
  return { items, total };
};

const callUpdateArticleTranslations = async (id: string, ja: string, en: string, zh: string): Promise<any> => {
  const app = getApp();
  if (app && typeof app.UpdateArticleTranslations === 'function') {
    return await app.UpdateArticleTranslations(id, ja, en, zh);
  }
  const art = localMockArticles.find((a) => a.id === id);
  if (art) {
    art.content.ja = ja;
    art.content.en = en;
    art.content.zh = zh;
  }
  return true;
};

// --- Audit RPC & Mock (SPEC-AUDIT-001) ---
let localMockAuditReport: any = {
  executed_at: new Date().toISOString(),
  integrity_ok: true,
  integrity_errors: ['ok'],
  foreign_key_ok: true,
  foreign_key_errors: [],
  orphan_db_media: [
    {
      media_id: 'mock_orphan_media_1',
      article_id: 'deleted_art_99',
      type: 'image',
      download_url: 'https://example.com/img_orphan.jpg',
      status: 'QUEUED',
      reason: '親記事 (deleted_art_99) が存在しません',
    },
  ],
  orphan_files: [
    {
      path: './stash/scenes/mock_orphan_scene_101.mp4',
      file_name: 'mock_orphan_scene_101.mp4',
      file_size: 15420000,
      category: 'stash_scene',
    },
    {
      path: './blobs/mock_orphan_blob_202.jpg',
      file_name: 'mock_orphan_blob_202.jpg',
      file_size: 245000,
      category: 'blob',
    },
  ],
  purged_file_count: 0,
  purged_db_media_count: 0,
  summary: '要クレンジング (孤立DB:1件, 孤立ファイル:2件)',
};

const callRunAudit = async (purgeFiles: boolean, purgeDB: boolean): Promise<any> => {
  const app = getApp();
  if (app && typeof app.RunAudit === 'function') {
    return await app.RunAudit(purgeFiles, purgeDB);
  }
  const report = JSON.parse(JSON.stringify(localMockAuditReport));
  report.executed_at = new Date().toISOString();
  if (purgeFiles) {
    report.purged_file_count = report.orphan_files.length;
    report.orphan_files = [];
    localMockAuditReport.orphan_files = [];
  }
  if (purgeDB) {
    report.purged_db_media_count = report.orphan_db_media.length;
    report.orphan_db_media = [];
    localMockAuditReport.orphan_db_media = [];
  }
  if (report.orphan_files.length === 0 && report.orphan_db_media.length === 0) {
    report.summary = '健全';
  }
  return report;
};

const callPurgeOrphanFiles = async (paths: string[]): Promise<number> => {
  const app = getApp();
  if (app && typeof app.PurgeOrphanFiles === 'function') {
    return await app.PurgeOrphanFiles(paths);
  }
  localMockAuditReport.orphan_files = localMockAuditReport.orphan_files.filter(
    (f: any) => !paths.includes(f.path)
  );
  return paths.length;
};

const callPurgeOrphanDBMedia = async (mediaIDs: string[]): Promise<number> => {
  const app = getApp();
  if (app && typeof app.PurgeOrphanDBMedia === 'function') {
    return await app.PurgeOrphanDBMedia(mediaIDs);
  }
  localMockAuditReport.orphan_db_media = localMockAuditReport.orphan_db_media.filter(
    (m: any) => !mediaIDs.includes(m.media_id)
  );
  return mediaIDs.length;
};


export function useAdmin() {
  const activeJob = ref<models.JobProgress | null>(null);
  const jobList = ref<models.JobProgress[]>([]);
  const logs = ref<string[]>([]);
  const config = ref<models.AppConfig | null>(null);

  const loadingJobs = ref(false);
  const loadingConfig = ref(false);
  const savingConfig = ref(false);
  const actionLoading = ref(false);
  const saveStatus = ref<{ success: boolean; message: string } | null>(null);

  const salvageForm = reactive({
    platform: 'twitter',
    account: '',
    limit: 50,
  });

  const importForm = reactive({
    warcPath: '',
    offline: false,
  });

  // ジョブ一覧および実行中ジョブの取得
  const fetchJobs = async () => {
    loadingJobs.value = true;
    try {
      const [active, list] = await Promise.all([
        callGetActiveJob().catch(() => null),
        callListJobs().catch(() => []),
      ]);
      activeJob.value = active && active.id ? active : null;
      jobList.value = (list || []).sort((a, b) => {
        const timeA = a.started_at ? new Date(a.started_at).getTime() : 0;
        const timeB = b.started_at ? new Date(b.started_at).getTime() : 0;
        return timeB - timeA;
      });

      if (activeJob.value && activeJob.value.logs) {
        logs.value = [...activeJob.value.logs];
      }
    } catch (err: any) {
      console.error('[useAdmin] Failed to fetch jobs:', err);
    } finally {
      loadingJobs.value = false;
    }
  };

  // 設定のロード
  const loadConfig = async () => {
    loadingConfig.value = true;
    saveStatus.value = null;
    try {
      const res = await callGetConfig();
      config.value = models.AppConfig.createFrom(res);
    } catch (err: any) {
      console.error('[useAdmin] Failed to load config:', err);
      saveStatus.value = {
        success: false,
        message: `設定の読み込みに失敗しました: ${err?.message || err}`,
      };
    } finally {
      loadingConfig.value = false;
    }
  };

  // 設定の保存
  const saveConfig = async () => {
    if (!config.value) return;
    savingConfig.value = true;
    saveStatus.value = null;
    try {
      await callSaveConfig(config.value);
      saveStatus.value = {
        success: true,
        message: '設定を正常に保存しました (config.json)',
      };
      setTimeout(() => {
        if (saveStatus.value?.success) {
          saveStatus.value = null;
        }
      }, 4000);
    } catch (err: any) {
      console.error('[useAdmin] Failed to save config:', err);
      saveStatus.value = {
        success: false,
        message: `保存に失敗しました: ${err?.message || err}`,
      };
    } finally {
      savingConfig.value = false;
    }
  };

  // サルベージジョブ開始
  const startSalvage = async () => {
    if (!salvageForm.account.trim()) {
      alert('アカウント名/ID を入力してください');
      return;
    }
    actionLoading.value = true;
    logs.value = [];
    try {
      const job = await callStartSalvageJob(
        salvageForm.platform,
        salvageForm.account.trim(),
        Number(salvageForm.limit) || 50
      );
      activeJob.value = job;
      if (job && job.logs) {
        logs.value = [...job.logs];
      }
      await fetchJobs();
    } catch (err: any) {
      alert(`サルベージジョブの開始に失敗しました: ${err?.message || err}`);
    } finally {
      actionLoading.value = false;
    }
  };

  // 手動インポート開始
  const startImport = async () => {
    if (!importForm.warcPath.trim()) {
      alert('WARC ファイルパスを入力してください');
      return;
    }
    actionLoading.value = true;
    logs.value = [];
    try {
      const job = await callStartManualImportJob(
        importForm.warcPath.trim(),
        importForm.offline
      );
      activeJob.value = job;
      await fetchJobs();
    } catch (err: any) {
      alert(`インポートジョブの開始に失敗しました: ${err?.message || err}`);
    } finally {
      actionLoading.value = false;
    }
  };

  // ジョブキャンセル
  const cancelJob = async (jobId: string) => {
    if (!confirm(`ジョブ [${jobId}] を中止しますか？`)) return;
    try {
      await callCancelJob(jobId);
      await fetchJobs();
    } catch (err: any) {
      alert(`ジョブのキャンセルに失敗しました: ${err?.message || err}`);
    }
  };

  // ログクリア
  const clearLogs = () => {
    logs.value = [];
  };

  // イベント購読ハンドラ
  const handleJobProgress = (progress: models.JobProgress) => {
    if (activeJob.value && activeJob.value.id === progress.id) {
      activeJob.value = progress;
    } else if (progress.status === 'running') {
      activeJob.value = progress;
    }
    if (progress.logs && progress.logs.length > logs.value.length) {
      logs.value = [...progress.logs];
    }
    // 履歴リストの該当アイテムも更新
    const idx = jobList.value.findIndex((j) => j.id === progress.id);
    if (idx !== -1) {
      jobList.value[idx] = progress;
    } else {
      jobList.value.unshift(progress);
    }
  };

  const handleJobFinished = (progress: models.JobProgress) => {
    handleJobProgress(progress);
    if (activeJob.value && activeJob.value.id === progress.id) {
      activeJob.value = null;
    }
    fetchJobs();
  };

  const handleJobLog = (logLine: string) => {
    logs.value.push(logLine);
    if (logs.value.length > 1000) {
      logs.value.shift();
    }
  };

  // ライフサイクル管理
  let isMounted = false;
  const setupEventListeners = () => {
    if (isMounted) return;
    isMounted = true;
    try {
      if (typeof EventsOn === 'function') {
        EventsOn('job:progress', handleJobProgress);
        EventsOn('job:started', handleJobProgress);
        EventsOn('job:finished', handleJobFinished);
        EventsOn('job:log', handleJobLog);
      }
    } catch (_) {}
  };

  const cleanupEventListeners = () => {
    isMounted = false;
    try {
      if (typeof EventsOff === 'function') {
        EventsOff('job:progress');
        EventsOff('job:started');
        EventsOff('job:finished');
        EventsOff('job:log');
      }
    } catch (_) {}
  };

  // --- Whitelist 管理状態 ---
  const whitelistList = ref<Array<{ id: number; type: string; value: string; is_active: boolean }>>([]);
  const loadingWhitelist = ref(false);
  const whitelistStatus = ref<{ success: boolean; message: string } | null>(null);

  const fetchWhitelists = async () => {
    loadingWhitelist.value = true;
    try {
      const list = await callGetWhitelists();
      whitelistList.value = list || [];
    } catch (err: any) {
      console.error('[useAdmin] Failed to fetch whitelists:', err);
    } finally {
      loadingWhitelist.value = false;
    }
  };

  const addWhitelist = async (type: string, value: string) => {
    if (!value.trim()) return;
    try {
      await callAddWhitelist(type, value.trim());
      await fetchWhitelists();
      whitelistStatus.value = { success: true, message: `[${type}] ${value} を追加しました` };
      setTimeout(() => (whitelistStatus.value = null), 3000);
    } catch (err: any) {
      console.error('[useAdmin] Failed to add whitelist:', err);
      whitelistStatus.value = { success: false, message: `追加に失敗しました: ${err?.message || err}` };
    }
  };

  const updateWhitelist = async (id: number, type: string, value: string, isActive: boolean) => {
    try {
      await callUpdateWhitelist(id, type, value.trim(), isActive);
      await fetchWhitelists();
      whitelistStatus.value = { success: true, message: `ID:${id} を更新しました` };
      setTimeout(() => (whitelistStatus.value = null), 3000);
    } catch (err: any) {
      console.error('[useAdmin] Failed to update whitelist:', err);
      whitelistStatus.value = { success: false, message: `更新に失敗しました: ${err?.message || err}` };
    }
  };

  const deleteWhitelist = async (id: number) => {
    try {
      await callDeleteWhitelist(id);
      await fetchWhitelists();
      whitelistStatus.value = { success: true, message: `ID:${id} を削除しました` };
      setTimeout(() => (whitelistStatus.value = null), 3000);
    } catch (err: any) {
      console.error('[useAdmin] Failed to delete whitelist:', err);
      whitelistStatus.value = { success: false, message: `削除に失敗しました: ${err?.message || err}` };
    }
  };

  const toggleWhitelist = async (id: number) => {
    try {
      await callToggleWhitelist(id);
      await fetchWhitelists();
    } catch (err: any) {
      console.error('[useAdmin] Failed to toggle whitelist:', err);
    }
  };

  // --- Database 閲覧 ＆ 翻訳エディタ状態 ---
  const dbArticles = ref<models.RenderTree[]>([]);
  const totalArticles = ref(0);
  const selectedArticle = ref<models.RenderTree | null>(null);
  const loadingArticles = ref(false);
  const savingTranslation = ref(false);
  const translationStatus = ref<{ success: boolean; message: string } | null>(null);

  const dbSearchParams = reactive({
    query: '',
    accountID: 'all',
    filter: 'all',
    page: 1,
    limit: 10,
  });

  const searchArticles = async () => {
    loadingArticles.value = true;
    try {
      const offset = (dbSearchParams.page - 1) * dbSearchParams.limit;
      const res = await callSearchArticles(
        dbSearchParams.query.trim(),
        dbSearchParams.accountID,
        dbSearchParams.filter,
        dbSearchParams.limit,
        offset
      );
      dbArticles.value = (res.items || []).map((item: any) => models.RenderTree.createFrom(item));
      totalArticles.value = res.total || 0;

      // 選択中記事があれば更新または最新を選択
      if (selectedArticle.value) {
        const found = dbArticles.value.find((a) => a.id === selectedArticle.value?.id);
        if (found) {
          selectedArticle.value = found;
        }
      } else if (dbArticles.value.length > 0) {
        selectedArticle.value = dbArticles.value[0];
      }
    } catch (err: any) {
      console.error('[useAdmin] Failed to search articles:', err);
    } finally {
      loadingArticles.value = false;
    }
  };

  const selectArticle = (art: models.RenderTree) => {
    selectedArticle.value = art;
    translationStatus.value = null;
  };

  const saveArticleTranslations = async (id: string, ja: string, en: string, zh: string) => {
    savingTranslation.value = true;
    translationStatus.value = null;
    try {
      await callUpdateArticleTranslations(id, ja, en, zh);
      if (selectedArticle.value && selectedArticle.value.id === id) {
        selectedArticle.value.content.ja = ja;
        selectedArticle.value.content.en = en;
        selectedArticle.value.content.zh = zh;
      }
      // リスト側の該当記事も同期
      const idx = dbArticles.value.findIndex((a) => a.id === id);
      if (idx !== -1) {
        dbArticles.value[idx].content.ja = ja;
        dbArticles.value[idx].content.en = en;
        dbArticles.value[idx].content.zh = zh;
      }
      translationStatus.value = { success: true, message: '3言語翻訳テキストを保存しました' };
      setTimeout(() => (translationStatus.value = null), 4000);
    } catch (err: any) {
      console.error('[useAdmin] Failed to save translations:', err);
      translationStatus.value = { success: false, message: `翻訳の保存に失敗しました: ${err?.message || err}` };
    } finally {
      savingTranslation.value = false;
    }
  };

  // --- Audit 整合性監査 ＆ パージ状態 (SPEC-AUDIT-001) ---
  const auditReport = ref<any | null>(null);
  const loadingAudit = ref(false);
  const purgingFiles = ref(false);
  const purgingDB = ref(false);
  const auditStatus = ref<{ success: boolean; message: string } | null>(null);

  const runAudit = async (purgeFiles = false, purgeDB = false) => {
    loadingAudit.value = true;
    auditStatus.value = null;
    try {
      const report = await callRunAudit(purgeFiles, purgeDB);
      auditReport.value = report;
      if (purgeFiles || purgeDB) {
        auditStatus.value = {
          success: true,
          message: `監査とパージが完了しました (退避ファイル: ${report.purged_file_count || 0}件, 削除DBレコード: ${report.purged_db_media_count || 0}件)`,
        };
      } else {
        auditStatus.value = {
          success: true,
          message: `整合性監査が完了しました: ${report.summary || '完了'}`,
        };
      }
      setTimeout(() => (auditStatus.value = null), 5000);
    } catch (err: any) {
      console.error('[useAdmin] Failed to run audit:', err);
      auditStatus.value = {
        success: false,
        message: `監査の実行に失敗しました: ${err?.message || err}`,
      };
    } finally {
      loadingAudit.value = false;
    }
  };

  const purgeOrphanFiles = async (paths?: string[]) => {
    if (!auditReport.value) return;
    const targetPaths = paths || (auditReport.value.orphan_files || []).map((f: any) => f.path);
    if (targetPaths.length === 0) return;

    if (!confirm(`${targetPaths.length} 件の孤立ファイルをOSのごみ箱へ退避しますか？`)) {
      return;
    }

    purgingFiles.value = true;
    try {
      const count = await callPurgeOrphanFiles(targetPaths);
      auditReport.value.orphan_files = (auditReport.value.orphan_files || []).filter(
        (f: any) => !targetPaths.includes(f.path)
      );
      auditStatus.value = {
        success: true,
        message: `${count} 件の孤立ファイルをOSのごみ箱へ退避しました`,
      };
      setTimeout(() => (auditStatus.value = null), 4000);
    } catch (err: any) {
      console.error('[useAdmin] Failed to purge files:', err);
      auditStatus.value = {
        success: false,
        message: `ファイル退避に失敗しました: ${err?.message || err}`,
      };
    } finally {
      purgingFiles.value = false;
    }
  };

  const purgeOrphanDBMedia = async (mediaIDs?: string[]) => {
    if (!auditReport.value) return;
    const targetIDs = mediaIDs || (auditReport.value.orphan_db_media || []).map((m: any) => m.media_id);
    if (targetIDs.length === 0) return;

    if (!confirm(`${targetIDs.length} 件の孤立DBメディアレコードを削除しますか？`)) {
      return;
    }

    purgingDB.value = true;
    try {
      const count = await callPurgeOrphanDBMedia(targetIDs);
      auditReport.value.orphan_db_media = (auditReport.value.orphan_db_media || []).filter(
        (m: any) => !targetIDs.includes(m.media_id)
      );
      auditStatus.value = {
        success: true,
        message: `${count} 件の孤立DBレコードを削除しました`,
      };
      setTimeout(() => (auditStatus.value = null), 4000);
    } catch (err: any) {
      console.error('[useAdmin] Failed to purge DB records:', err);
      auditStatus.value = {
        success: false,
        message: `DBレコード削除に失敗しました: ${err?.message || err}`,
      };
    } finally {
      purgingDB.value = false;
    }
  };

  // --- Disaster Recovery (SPEC-RECOVERY-001) ---
  const restoringDB = ref(false);
  const restoreStatus = ref<{ success: boolean; message: string } | null>(null);

  const triggerRestore = async (dumpsDir?: string, resetDB: boolean = false) => {
    const confirmMsg = resetDB
      ? '🚨【完全再構築】現在のDBテーブルを初期化し、全ダンプファイルからオフライン再構築を開始します。よろしいですか？'
      : '🚨【災害復旧】全ダンプファイルからDBのオフライン再構築・メディア再同期を開始します。よろしいですか？';

    if (!confirm(confirmMsg)) {
      return;
    }

    restoringDB.value = true;
    try {
      const targetDir = dumpsDir || config.value?.storage?.dumps_dir || './backups/dumps';
      const job = await callTriggerRestore(targetDir, resetDB);
      if (job) {
        activeJob.value = job;
        if (!jobList.value.some((j: any) => j.id === job.id)) {
          jobList.value.unshift(job);
        }
      }
      restoreStatus.value = {
        success: true,
        message: '災害復旧（全ダンプからDB再構築）ジョブを開始しました',
      };
      setTimeout(() => (restoreStatus.value = null), 5000);
    } catch (err: any) {
      console.error('[useAdmin] Failed to trigger restore:', err);
      restoreStatus.value = {
        success: false,
        message: `災害復旧ジョブの開始に失敗しました: ${err?.message || err}`,
      };
    } finally {
      restoringDB.value = false;
    }
  };

  return {
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
    // Audit (SPEC-AUDIT-001)
    auditReport,
    loadingAudit,
    purgingFiles,
    purgingDB,
    auditStatus,
    runAudit,
    purgeOrphanFiles,
    purgeOrphanDBMedia,
    // Disaster Recovery (SPEC-RECOVERY-001)
    restoringDB,
    restoreStatus,
    triggerRestore,
  };
}


