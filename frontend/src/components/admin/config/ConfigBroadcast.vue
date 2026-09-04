<!-- frontend/src/components/admin/config/ConfigBroadcast.vue (100行以下) -->
<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { models } from '../../../../wailsjs/go/models';
import { useAdminBroadcast } from '../../../composables/admin/useAdminBroadcast';
import { useStashResolver } from '../../../composables/useStashResolver';

const props = defineProps<{ config: models.AppConfig }>();
const { broadcastStatus, isBroadcastLoading, fetchBroadcastStatus, toggleBroadcast } = useAdminBroadcast();
const { stashPort } = useStashResolver();
const copiedIndex = ref<number | null>(null);

const isEnabled = computed({
  get: () => props.config?.broadcast?.enabled ?? false,
  set: async (val: boolean) => {
    if (!props.config.broadcast) {
      props.config.broadcast = { enabled: val, allowed_networks: [] } as unknown as models.BroadcastConfig;
    } else {
      props.config.broadcast.enabled = val;
    }
    await toggleBroadcast(val);
  },
});

const allowedNetworksText = computed({
  get: () => props.config?.broadcast?.allowed_networks?.join(', ') || '',
  set: (val: string) => {
    if (!props.config.broadcast) {
      props.config.broadcast = { enabled: false, allowed_networks: [] } as unknown as models.BroadcastConfig;
    }
    props.config.broadcast.allowed_networks = val.split(',').map(s => s.trim()).filter(Boolean);
  },
});

const copyUrl = (url: string, index: number) => {
  navigator.clipboard.writeText(url);
  copiedIndex.value = index;
  setTimeout(() => { copiedIndex.value = null; }, 2000);
};

const addSubnet = (subnet: string) => {
  const current = props.config?.broadcast?.allowed_networks || [];
  if (!current.includes(subnet)) props.config.broadcast.allowed_networks = [...current, subnet];
};

onMounted(() => { fetchBroadcastStatus(); });
</script>

<template>
  <div class="bg-slate-900/80 border border-slate-800 rounded-xl p-4 space-y-4">
    <div class="flex items-center justify-between">
      <div class="flex items-center gap-2">
        <h3 class="text-sm font-bold text-slate-200 flex items-center gap-2"><span>📡</span> LAN Broadcast ＆ キャスト配信設定</h3>
        <span class="text-[10px] font-mono px-2 py-0.5 rounded" :class="broadcastStatus?.running ? 'bg-emerald-950/80 text-emerald-400 border border-emerald-500/40' : 'bg-slate-800 text-slate-400 border border-slate-700'">
          {{ broadcastStatus?.running ? '● 配信中 (HTTPS)' : '○ 停止中' }}
        </span>
      </div>
      <label class="flex items-center gap-2 cursor-pointer select-none">
        <input type="checkbox" v-model="isEnabled" :disabled="isBroadcastLoading" class="rounded bg-slate-950 border-slate-700 text-blue-600 focus:ring-0" />
        <span class="text-xs font-semibold text-slate-300">LAN配信を有効化</span>
      </label>
    </div>

    <div class="space-y-3 pt-2">
      <div>
        <label class="block text-xs text-slate-400 mb-1">許可するLANサブネット (カンマ区切りCIDR)</label>
        <input v-model="allowedNetworksText" type="text" placeholder="192.168.10.0/24, 192.168.3.0/24, 127.0.0.1/32" class="w-full bg-slate-950 border border-slate-700 rounded-lg px-3 py-1.5 text-xs text-slate-200 font-mono" />
      </div>

      <div v-if="broadcastStatus?.detected_subnets?.length" class="flex flex-wrap items-center gap-2 text-xs">
        <span class="text-slate-500">検出されたサブネット:</span>
        <button v-for="sub in broadcastStatus.detected_subnets" :key="sub" type="button" @click="addSubnet(sub)" class="px-2 py-0.5 rounded bg-blue-950/60 border border-blue-700/50 hover:bg-blue-900 text-blue-300 text-[11px] font-mono transition-colors" title="クリックして許可リストに追加">
          + {{ sub }}
        </button>
      </div>

      <div class="space-y-2">
        <div class="text-xs text-slate-400 flex items-center justify-between">
          <span>配信アクセスURL (LAN内スマホ・PCからアクセス):</span>
          <span v-if="broadcastStatus?.running" class="text-[11px] text-amber-400/90">※初回「詳細設定→続行」をタップ</span>
        </div>
        <div v-if="!broadcastStatus?.running" class="p-3 bg-slate-950/40 border border-dashed border-slate-800 rounded-lg text-xs text-slate-500 text-center">
          現在LAN配信は停止されています（上の「LAN配信を有効化」で開始できます）
        </div>
        <div v-else v-for="(ip, idx) in (broadcastStatus?.local_ips?.length ? broadcastStatus.local_ips : ['127.0.0.1'])" :key="ip" class="p-2.5 bg-slate-950/80 border border-slate-800/80 rounded-lg flex items-center justify-between">
          <div class="text-xs font-mono text-slate-300 truncate">
            <span class="text-slate-500 mr-2">{{ ip.startsWith('192.168.10.') ? '🌐 [10系]' : ip.startsWith('192.168.3.') ? '🌐 [3系]' : '🏠 [Local]' }}</span>
            <span class="text-emerald-400 font-bold">http://{{ ip }}:{{ config?.network?.middleware_port || 5175 }}</span>
          </div>
          <div class="flex items-center gap-1.5 ml-2">
            <button type="button" @click="copyUrl(`http://${ip}:${config?.network?.middleware_port || 5175}`, idx)" class="px-2.5 py-1 bg-slate-800 hover:bg-slate-700 text-slate-300 text-xs rounded transition-colors whitespace-nowrap">
              {{ copiedIndex === idx ? 'コピー完了!' : '土蔵URL' }}
            </button>
            <button type="button" @click="copyUrl(`http://${ip}:${stashPort}`, idx + 100)" class="px-2.5 py-1 bg-purple-950/80 hover:bg-purple-900 border border-purple-700/50 text-purple-300 text-xs rounded transition-colors whitespace-nowrap" title="Stash WebUI のアクセスURLをコピー">
              {{ copiedIndex === (idx + 100) ? 'コピー完了!' : 'Stash URL' }}
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
