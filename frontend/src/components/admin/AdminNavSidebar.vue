<!-- frontend/src/components/admin/AdminNavSidebar.vue (100行以下) -->
<script setup lang="ts">
export type AdminTabId = 'jobs' | 'audit' | 'stash' | 'posts' | 'media' | 'accounts' | 'whitelist' | 'config' | 'skin';

defineProps<{ activeTab: AdminTabId }>();
const emit = defineEmits<{ (e: 'select', tab: AdminTabId): void }>();

const groups = [
  {
    title: '運用・実行',
    items: [
      { id: 'jobs', icon: '🚀', label: 'ジョブ制御' },
      { id: 'audit', icon: '🩺', label: '整合性監査' },
      { id: 'stash', icon: '🎛️', label: 'Stash状態' },
    ],
  },
  {
    title: 'データ管理',
    items: [
      { id: 'posts', icon: '📝', label: '投稿・翻訳' },
      { id: 'media', icon: '🖼️', label: 'メディア' },
      { id: 'accounts', icon: '👤', label: 'アカウント' },
      { id: 'whitelist', icon: '🛡️', label: 'Whitelist' },
    ],
  },
  {
    title: 'システム設定',
    items: [
      { id: 'config', icon: '🌐', label: 'システム設定' },
      { id: 'skin', icon: '🎨', label: 'スキン・フォント' },
    ],
  },
];
</script>

<template>
  <aside class="w-44 bg-slate-900/80 border-r border-slate-800 flex flex-col py-3 px-2 gap-4 shrink-0 select-none overflow-y-auto">
    <div v-for="g in groups" :key="g.title" class="space-y-1">
      <div class="px-2 text-[10px] font-bold text-slate-500 uppercase tracking-wider">{{ g.title }}</div>
      <button
        v-for="item in g.items"
        :key="item.id"
        @click="emit('select', item.id as AdminTabId)"
        :class="[
          'w-full flex items-center gap-2 px-2.5 py-1.5 rounded-lg text-xs font-semibold transition-all text-left cursor-pointer',
          activeTab === item.id
            ? 'bg-blue-600 text-white shadow-md shadow-blue-900/30'
            : 'text-slate-400 hover:text-slate-200 hover:bg-slate-800/60'
        ]"
      >
        <span class="text-sm shrink-0">{{ item.icon }}</span>
        <span class="truncate">{{ item.label }}</span>
      </button>
    </div>
  </aside>
</template>
