// frontend/src/models/adminTabs.ts (100行以下 - SPEC-PRINCIPLE-001)

export type AdminTabId = 
  | 'plugins' | 'explorer' | 'audit' | 'stash' | 'downloaders' | 'posts' | 'media' | 'accounts' | 'config';

export interface AdminTabItem { id: AdminTabId; icon: string; label: string; }
export interface AdminTabGroup { title: string; items: AdminTabItem[]; }

export const adminTabGroups: AdminTabGroup[] = [
  {
    title: '運用・実行',
    items: [
      { id: 'plugins', icon: '🧩', label: 'プラグイン＆ジョブ' },
      { id: 'explorer', icon: '🧭', label: 'リレーション探査' },
      { id: 'audit', icon: '🩺', label: '整合性監査' },
      { id: 'stash', icon: '🎛️', label: 'Stash状態' },
      { id: 'downloaders', icon: '⚡', label: 'Motrix & Thunder' },
    ],
  },
  {
    title: 'データ管理',
    items: [
      { id: 'accounts', icon: '👤', label: 'アカウント' },
      { id: 'posts', icon: '📝', label: '投稿・翻訳' },
      { id: 'media', icon: '🖼️', label: 'メディア' },
    ],
  },
  {
    title: '設定・環境',
    items: [
      { id: 'config', icon: '⚙️', label: 'システム設定' },
    ],
  },
];
