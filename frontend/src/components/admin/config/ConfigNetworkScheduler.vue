<!-- frontend/src/components/admin/config/ConfigNetworkScheduler.vue (100行以下) -->
<script setup lang="ts">
import { models } from '../../../../wailsjs/go/models';
defineProps<{ config: models.AppConfig }>();
</script>

<template>
  <div class="space-y-4">
    <!-- Network 設定 -->
    <div class="bg-slate-900/80 border border-slate-800 rounded-xl p-4 space-y-4">
      <h3 class="text-sm font-bold text-slate-200 flex items-center gap-2"><span>🔌</span> ネットワーク・ポート (Network)</h3>
      <div class="grid grid-cols-1 sm:grid-cols-3 gap-4">
        <div>
          <label class="block text-xs text-slate-400 mb-1">Stashapp ポート</label>
          <input v-model.number="config.network.stash_port" type="number" class="w-full bg-slate-950 border border-slate-700 rounded-lg px-3 py-1.5 text-xs text-slate-200 font-mono" />
        </div>
        <div>
          <label class="block text-xs text-slate-400 mb-1">中継サーバーポート</label>
          <input v-model.number="config.network.middleware_port" type="number" class="w-full bg-slate-950 border border-slate-700 rounded-lg px-3 py-1.5 text-xs text-slate-200 font-mono" />
        </div>
        <div>
          <label class="block text-xs text-slate-400 mb-1">バインドアドレス</label>
          <input v-model="config.network.internal_bind_address" type="text" class="w-full bg-slate-950 border border-slate-700 rounded-lg px-3 py-1.5 text-xs text-slate-200 font-mono" />
        </div>
      </div>
    </div>

    <!-- Scheduler 設定 -->
    <div class="bg-slate-900/80 border border-slate-800 rounded-xl p-4 space-y-4">
      <h3 class="text-sm font-bold text-slate-200 flex items-center gap-2"><span>⏱️</span> 自動ポーリング・世代管理 (Scheduler)</h3>
      <div class="grid grid-cols-1 sm:grid-cols-3 gap-4">
        <div>
          <label class="block text-xs text-slate-400 mb-1">回収ポーリング間隔 (秒)</label>
          <input v-model.number="config.scheduler.poll_interval_sec" type="number" class="w-full bg-slate-950 border border-slate-700 rounded-lg px-3 py-1.5 text-xs text-slate-200 font-mono" />
        </div>
        <div>
          <label class="block text-xs text-slate-400 mb-1">バックアップ周期 (時間)</label>
          <input v-model.number="config.scheduler.backup_interval_hours" type="number" class="w-full bg-slate-950 border border-slate-700 rounded-lg px-3 py-1.5 text-xs text-slate-200 font-mono" />
        </div>
        <div>
          <label class="block text-xs text-slate-400 mb-1">バックアップ保持世代数</label>
          <input v-model.number="config.scheduler.max_backup_generations" type="number" class="w-full bg-slate-950 border border-slate-700 rounded-lg px-3 py-1.5 text-xs text-slate-200 font-mono" />
        </div>
      </div>
    </div>
  </div>
</template>
