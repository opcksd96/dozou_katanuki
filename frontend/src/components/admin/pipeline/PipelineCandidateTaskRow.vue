<!-- frontend/src/components/admin/pipeline/PipelineCandidateTaskRow.vue (100行以下 - SPEC-PRINCIPLE-001) -->
<script setup lang="ts">
import { ref, computed } from 'vue';
import type { CandidateTask, ActiveThunderTask } from '../../../types/pipeline_media_tasks';

const props = defineProps<{ task: CandidateTask; index?: number }>();
const showRaw = ref(false);

const cdp = computed<Partial<ActiveThunderTask>>(() => {
  if (props.task.cdp) return props.task.cdp;
  if (props.task.raw_json) {
    try {
      const p = JSON.parse(props.task.raw_json);
      return {
        task_id: p.taskId, file_size: p.fileSize, download_size: p.downloadSize,
        download_speed: p.downloadSpeed, vip_speed: p.vipSpeed,
        task_status_code: p.taskStatus, error_code: p.errorCode,
        save_path: p.savePath, create_time: p.createTime, cid: p.cid, gcid: p.gcid,
        origin: p.origin, raw_json: props.task.raw_json,
      };
    } catch (_) {}
  }
  return {
    task_id: props.task.thunder_task_id, file_size: props.task.file_size,
    download_size: props.task.download_size, save_path: props.task.save_path,
    task_status_code: props.task.task_status_code, error_code: props.task.error_code,
    detail_text: props.task.detail_text || props.task.error_reason,
  };
});

const isReaped = computed(() => props.task.status === 'REAPED');
const isRetired = computed(() => props.task.status === 'RETIRED');
const fmtB = (b?: number) => {
  if (!b || b <= 0) return '0 B';
  const k = 1024, s = ['B', 'KB', 'MB', 'GB', 'TB'];
  const i = Math.floor(Math.log(b) / Math.log(k));
  return (b / Math.pow(k, i)).toFixed(1) + ' ' + s[i];
};
const fmtRaw = (n?: number) => n ? n.toLocaleString() + ' B' : '0 B';
const copy = (txt?: string) => { if (txt) { navigator.clipboard.writeText(txt); alert('コピーしました'); } };
</script>

<template>
  <div class="p-3 rounded-xl bg-slate-900/95 border border-slate-800 space-y-2 font-mono text-xs shadow-md">
    <!-- ヘッダー: No, FileName, Status, スキーム, TaskID, JSON切替 -->
    <div class="flex flex-wrap items-center justify-between gap-1.5 border-b border-slate-800/80 pb-2">
      <div class="flex items-center gap-1.5 flex-wrap">
        <span class="px-1.5 py-0.5 rounded text-[10px] font-bold bg-slate-800 text-slate-400 border border-slate-700">#{{ index ?? '-' }}</span>
        <span class="px-2 py-0.5 rounded text-[10px] font-bold bg-indigo-950 text-indigo-300 border border-indigo-800">{{ task.resolution_type || 'orig' }}</span>
        <span
          class="px-2 py-0.5 rounded text-[10px] font-bold border"
          :class="isReaped ? 'bg-rose-950 text-rose-300 border-rose-800' : (isRetired ? 'bg-amber-950 text-amber-300 border-amber-800' : (cdp.task_id ? 'bg-emerald-950 text-emerald-300 border-emerald-800' : (task.status === 'ONBOARDED' ? 'bg-amber-950 text-amber-300 border-amber-800' : 'bg-slate-800 text-slate-200 border-slate-700')))"
        >
          {{ isReaped ? '🔴 失敗 (REAPED)' : (isRetired ? '🛑 終了 (RETIRED)' : (cdp.status || (task.status === 'ONBOARDED' && !cdp.task_id ? '⚠️ 迅雷不在' : task.status))) }}
        </span>
        <span v-if="cdp.task_id" class="px-1.5 py-0.5 rounded text-[9.5px] bg-purple-950 text-purple-300 border border-purple-800">⚡ ID:{{ cdp.task_id }}</span>
        <span v-if="cdp.error_code" class="px-1.5 py-0.5 rounded text-[9.5px] bg-rose-950 text-rose-300 border border-rose-800">Err:{{ cdp.error_code }}</span>
      </div>
      <div class="flex items-center gap-2">
        <span class="text-[11px] text-slate-200 font-bold truncate max-w-[280px]" :title="task.file_name">{{ task.file_name }}</span>
        <button v-if="task.raw_json || cdp.raw_json" @click="showRaw = !showRaw" class="text-[9.5px] px-2 py-0.5 rounded bg-slate-800 hover:bg-slate-700 text-slate-300 cursor-pointer">
          {{ showRaw ? '閉じる' : '{ } 生JSON' }}
        </button>
      </div>
    </div>

    <!-- CSV 12項目マッピング グリッド -->
    <div class="grid grid-cols-2 sm:grid-cols-4 gap-1.5 p-2 bg-slate-950/70 rounded-lg border border-slate-800/80 text-[10px]">
      <div>Progress: <span class="text-emerald-300 font-bold">{{ cdp.progress_text || (isReaped ? '0 B / 0 B' : (cdp.file_size ? `${fmtB(cdp.download_size)} / ${fmtB(cdp.file_size)}` : '-')) }}</span></div>
      <div>Speed: <span class="text-sky-300 font-bold">{{ cdp.speed_text || (cdp.download_speed ? fmtB(cdp.download_speed) + '/s' : '0 B/s') }}</span></div>
      <div>VIPSpeed: <span class="text-amber-300">{{ cdp.vip_speed ? fmtB(cdp.vip_speed) + '/s' : '0 B/s' }}</span></div>
      <div>FileSize: <span class="text-slate-300 font-bold">{{ fmtRaw(cdp.file_size) }}</span></div>
      <div>DownloadSize: <span class="text-slate-300">{{ fmtRaw(cdp.download_size) }}</span></div>
      <div>DownloadSpeed: <span class="text-slate-400">{{ cdp.download_speed ? cdp.download_speed.toLocaleString() + ' B/s' : '0 B/s' }}</span></div>
      <div>VIPSpeed(raw): <span class="text-slate-400">{{ cdp.vip_speed ? cdp.vip_speed.toLocaleString() + ' B/s' : '0 B/s' }}</span></div>
      <div>StatusCode: <span class="text-slate-400 font-bold" :class="isReaped ? 'text-rose-400' : ''">{{ cdp.task_status_code ?? (isReaped ? 9 : '-') }}</span></div>
      <div v-if="cdp.save_path" class="col-span-2 sm:col-span-4 truncate text-slate-400" :title="cdp.save_path">SavePath: <span class="text-slate-300">{{ cdp.save_path }}</span></div>
    </div>

    <!-- Target_URL -->
    <div class="flex items-center justify-between gap-2 text-[10.5px] bg-slate-950/90 p-1.5 rounded-lg border border-slate-800/80">
      <span class="text-slate-400 truncate select-all flex-1" :title="task.url">Target_URL: <span class="text-slate-200">{{ task.url }}</span></span>
      <button @click.stop="copy(task.url)" class="text-[9.5px] text-indigo-400 hover:text-indigo-200 px-2 py-0.5 bg-slate-800 rounded cursor-pointer shrink-0">コピー</button>
    </div>

    <!-- Detail (迅雷DOMメッセージ) -->
    <div v-if="cdp.detail_text || task.error_reason || cdp.raw_text" class="text-[10.5px] text-rose-300 bg-rose-950/40 p-2 rounded-lg border border-rose-900/60 leading-relaxed">
      Detail: {{ cdp.detail_text || task.error_reason || cdp.raw_text }}
    </div>

    <!-- 生JSONプレビュー (本物の生JSONが存在する場合のみ表示) -->
    <div v-if="showRaw && (task.raw_json || cdp.raw_json)" class="p-2.5 bg-slate-950 rounded-xl border border-slate-800 text-[10px] text-slate-300 overflow-x-auto select-text">
      <pre>{{ JSON.stringify(cdp.raw_json ? JSON.parse(cdp.raw_json) : JSON.parse(task.raw_json || '{}'), null, 2) }}</pre>
    </div>
  </div>
</template>
