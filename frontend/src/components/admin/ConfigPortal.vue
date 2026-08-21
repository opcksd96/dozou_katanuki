<!-- frontend/src/components/admin/ConfigPortal.vue (100行以下) -->
<script setup lang="ts">
import { models } from '../../../wailsjs/go/models';
import ConfigSystemStorage from './config/ConfigSystemStorage.vue';
import ConfigNetworkScheduler from './config/ConfigNetworkScheduler.vue';
import ConfigBroadcast from './config/ConfigBroadcast.vue';

const props = defineProps<{
  config: models.AppConfig | null;
  loadingConfig: boolean;
  savingConfig: boolean;
  saveStatus: { success: boolean; message: string } | null;
}>();

const emit = defineEmits<{
  (e: 'saveConfig'): void;
  (e: 'loadConfig'): void;
}>();
</script>

<template>
  <div class="space-y-6">
    <div v-if="loadingConfig && !config" class="py-12 text-center text-slate-400">
      <span class="animate-spin text-xl inline-block mr-2">⏳</span> 設定を読み込み中...
    </div>

    <form v-else-if="config" @submit.prevent="emit('saveConfig')" class="space-y-6">
      <div v-if="saveStatus" class="p-3 rounded-lg text-xs font-semibold flex items-center justify-between" :class="saveStatus.success ? 'bg-emerald-950/80 border border-emerald-500/50 text-emerald-300' : 'bg-rose-950/80 border border-rose-500/50 text-rose-300'">
        <span>{{ saveStatus.message }}</span>
      </div>

      <ConfigSystemStorage :config="config" />
      <ConfigNetworkScheduler :config="config" />
      <ConfigBroadcast :config="config" />

      <!-- 保存ボタンバー -->
      <div class="sticky bottom-0 bg-slate-950/90 backdrop-blur-md p-4 -mx-4 -mb-4 border-t border-slate-800 flex items-center justify-between z-10">
        <button type="button" @click="emit('loadConfig')" :disabled="loadingConfig || savingConfig" class="px-4 py-2 bg-slate-800 hover:bg-slate-700 text-slate-300 text-xs rounded-lg border border-slate-700 disabled:opacity-50">
          🔄 設定を再読込
        </button>
        <button type="submit" :disabled="savingConfig" class="px-6 py-2 bg-blue-600 hover:bg-blue-500 text-white text-xs font-bold rounded-lg shadow-lg disabled:opacity-50 flex items-center gap-1.5">
          <span v-if="savingConfig" class="animate-spin">⏳</span>
          <span>💾 config.json に保存</span>
        </button>
      </div>
    </form>
  </div>
</template>
