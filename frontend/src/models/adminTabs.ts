// frontend/src/models/adminTabs.ts (100行以下 - SPEC-PRINCIPLE-001)

export type AdminTabId = 
  | 'plugins' | 'accounts' | 'posts' | 'media' | 'explorer' | 'downloaders' | 'audit' | 'console' | 'config';

export interface AdminTabItem { id: AdminTabId; icon: string; label: string; }
export interface AdminTabGroup { title: string; items: AdminTabItem[]; }

export const adminTabGroups: AdminTabGroup[] = [
  {
    title: 'サイドカー・採取',
    items: [
      { id: 'plugins', icon: '🧩', label: 'プラグイン (サイドカー)' },
    ],
  },
  {
    title: 'データ管理 (SSOT)',
    items: [
      { id: 'accounts', icon: '👤', label: 'アカウント' },
      { id: 'posts', icon: '📝', label: '投稿・翻訳' },
      { id: 'media', icon: '🖼️', label: 'メディア' },
      { id: 'explorer', icon: '🧭', label: 'リレーション探査' },
    ],
  },
  {
    title: 'システム運用・インフラ',
    items: [
      { id: 'downloaders', icon: '🚀', label: 'ダウンローダー遠隔管理' },
      { id: 'audit', icon: '🩺', label: 'データベース＆Stash監査' },
      { id: 'console', icon: '📜', label: 'システムコンソール' },
      { id: 'config', icon: '⚙️', label: 'システム環境設定' },
    ],
  },
];
