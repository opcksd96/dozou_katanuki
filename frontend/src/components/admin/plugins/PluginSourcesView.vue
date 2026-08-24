<!-- frontend/src/components/admin/plugins/PluginSourcesView.vue (100行以下) -->
<script setup lang="ts">
defineProps<{ platform: string }>();

const sources = [
  { id: 'all', name: '全ソース自動フォールバック (Auto Orchestrator)', priority: 1, type: 'Orchestrator', desc: '優先度順に多重フォールバック走査し、ID重複を完全排除して収集します。', status: 'Active 🟢' },
  { id: 'sotwe', name: 'Sotwe Mirror API', priority: 20, type: 'Archive Mirror', desc: 'sotwe.com からスレッド・高解像度メディア・リプライツリーをまるごと抽出します。', status: 'Active 🟢' },
  { id: 'twistalker', name: 'Twistalker Engine', priority: 30, type: 'Archive Mirror', desc: 'twistalker.com のHTML構造から凍結・消滅投稿をサルベージします。', status: 'Active 🟢' },
  { id: 'nitter', name: 'Nitter Distributed Cluster', priority: 40, type: 'Public Instance', desc: '分散Nitterインスタンス群（nitter.net等）を走査して収集します。', status: 'Active 🟢' },
  { id: 'wayback', name: 'Wayback Machine CDX API', priority: 50, type: 'Web Archive', desc: 'Internet Archive の CDX prefix走査と原本WARCキャッシュ保存を行います。', status: 'Active 🟢' },
  { id: 'official', name: 'X Syndication API', priority: 10, type: 'Official Web API', desc: 'X公式の公開Syndicationエンドポイントから直接取得を試みます。', status: 'Standby 🟡' },
];
</script>

<template>
  <div class="space-y-4">
    <div class="flex items-center justify-between">
      <div>
        <h4 class="text-xs font-bold text-slate-200 flex items-center gap-2">
          <span>📡</span> {{ platform.toUpperCase() }} サルベージソース・エンジン一覧
        </h4>
        <p class="text-[11px] text-slate-400 mt-0.5">プラグインコアに登録されている多重サルベージソースと優先順位（Priority Chain）です。</p>
      </div>
    </div>

    <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
      <div v-for="s in sources" :key="s.id" class="bg-slate-900/80 border border-slate-800 rounded-xl p-3 space-y-2 hover:border-slate-700 transition-colors">
        <div class="flex items-center justify-between">
          <div class="flex items-center gap-2">
            <span class="text-xs font-bold text-slate-200">{{ s.name }}</span>
            <span class="text-[9px] px-1.5 py-0.5 rounded bg-blue-950/80 text-blue-400 border border-blue-800/50 font-mono">{{ s.type }}</span>
          </div>
          <span class="text-[10px] font-mono font-semibold text-slate-300">{{ s.status }}</span>
        </div>
        <p class="text-[11px] text-slate-400 leading-relaxed">{{ s.desc }}</p>
        <div class="flex items-center justify-between text-[10px] text-slate-500 pt-1 border-t border-slate-800/60 font-mono">
          <span>Priority: #{{ s.priority }}</span>
          <span>Source ID: {{ s.id }}</span>
        </div>
      </div>
    </div>
  </div>
</template>
