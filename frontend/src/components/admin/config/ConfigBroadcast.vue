<!-- frontend/src/components/admin/config/ConfigBroadcast.vue (100行以下) -->
<script setup lang="ts">
import { ref, computed } from 'vue';
import { models } from '../../../../wailsjs/go/models';

const props = defineProps<{ config: models.AppConfig }>();
const copied = ref(false);

const allowedNetworksText = computed({
  get: () => props.config?.broadcast?.allowed_networks?.join(', ') || '',
  set: (val: string) => {
    if (!props.config.broadcast) {
      props.config.broadcast = { enabled: false, allowed_networks: [] } as unknown as models.BroadcastConfig;
    }
    props.config.broadcast.allowed_networks = val.split(',').map(s => s.trim()).filter(Boolean);
  },
});

const sampleCastURL = computed(() => `http://<あなたのPCのIPアドレス>:${props.config?.network?.middleware_port || 5175}`);
const copyCastURL = () => {
  navigator.clipboard.writeText(sampleCastURL.value);
  copied.value = true;
  setTimeout(() => { copied.value = false; }, 2000);
};
</script>

<template>
  <div class="bg-slate-900/80 border border-slate-800 rounded-xl p-4 space-y-4">
    <div class="flex items-center justify-between">
      <h3 class="text-sm font-bold text-slate-200 flex items-center gap-2"><span>📡</span> LAN Broadcast ＆ キャスト配信設定 (SPEC-SCHEDULER-001)</h3>
      <label class="flex items-center gap-2 cursor-pointer select-none">
        <input type="checkbox" v-model="config.broadcast.enabled" class="rounded bg-slate-950 border-slate-700 text-blue-600" />
        <span class="text-xs font-semibold text-slate-300">LAN配信を有効化</span>
      </label>
    </div>
    <div class="space-y-3 pt-2">
      <div>
        <label class="block text-xs text-slate-400 mb-1">許可するLANサブネット (カンマ区切りCIDR)</label>
        <input v-model="allowedNetworksText" type="text" placeholder="192.168.0.0/16, 10.0.0.0/8, 127.0.0.1/32" class="w-full bg-slate-950 border border-slate-700 rounded-lg px-3 py-1.5 text-xs text-slate-200 font-mono" />
      </div>
      <div class="p-3 bg-slate-950/80 border border-slate-800/80 rounded-lg flex items-center justify-between">
        <div class="text-xs font-mono text-slate-400 truncate">
          <span class="text-slate-500">サンプルURL: </span>
          <span class="text-emerald-400">{{ sampleCastURL }}</span>
        </div>
        <button type="button" @click="copyCastURL" class="px-2.5 py-1 bg-slate-800 hover:bg-slate-700 text-slate-300 text-xs rounded transition-colors whitespace-nowrap ml-2">
          {{ copied ? 'コピー完了!' : 'URLコピー' }}
        </button>
      </div>
    </div>
  </div>
</template>
