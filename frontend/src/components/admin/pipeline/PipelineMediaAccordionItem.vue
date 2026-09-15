<!-- frontend/src/components/admin/pipeline/PipelineMediaAccordionItem.vue (100行以下 - SPEC-PRINCIPLE-001) -->
<script setup lang="ts">
import { ref, computed } from 'vue';
import type { MediaWithTasks } from '../../../types/pipeline_media_tasks';
import PipelineCandidateTaskRow from './PipelineCandidateTaskRow.vue';

const props = defineProps<{ media: MediaWithTasks; isInitiallyOpen?: boolean }>();
const isOpen = ref(props.isInitiallyOpen ?? false);
const showMediaJson = ref(false);
const toggle = () => { isOpen.value = !isOpen.value; };

const displayStatus = computed(() => {
  const tasks = props.media.candidate_tasks || [];
  if (tasks.length === 0) return props.media.failed_reason || '待機中';
  const active = tasks.find((t) => t.status === 'RUNNING' || t.status === 'ONBOARDED' || t.status === 'HOLDING');
  if (active) {
    const res = active.resolution_type || 'orig';
    return active.status === 'HOLDING' ? `保留・接続中 (${res})` : `迅雷投入中 (${res})`;
  }
  const allDead = tasks.every((t) => t.status === 'REAPED' || t.status === 'RETIRED' || t.status === 'DEPLETED');
  if (allDead) return '🔴 全候補枯渇 (退避待機)';
  return props.media.failed_reason || '待機中';
});

const copyText = (txt?: string) => {
  if (txt) { navigator.clipboard.writeText(txt); alert('コピーしました'); }
};
</script>

<template>
  <div class="border border-slate-800 rounded-2xl overflow-hidden bg-slate-950/70 shadow-md">
    <!-- 親ヘッダー: 作品メタ情報（解像度・タイプ・ステータス・アカウント） -->
    <div
      @click="toggle"
      class="p-3 flex flex-col space-y-1.5 cursor-pointer hover:bg-slate-900/80 transition select-none"
      :class="isOpen ? 'bg-slate-900/60 border-b border-slate-800/80' : ''"
    >
      <div class="flex items-center justify-between gap-2 flex-wrap">
        <div class="flex items-center gap-2 flex-wrap min-w-0">
          <span class="text-xs text-slate-400 font-mono transition-transform duration-200" :class="isOpen ? 'rotate-90' : ''">▶</span>
          <span class="text-xs font-bold text-indigo-300 font-mono">{{ media.media_id }}</span>
          <span v-if="media.type" class="px-1.5 py-0.2 rounded text-[9px] font-bold bg-slate-800 text-slate-300 uppercase border border-slate-700">
            {{ media.type }}
          </span>
          <span v-if="media.width && media.height" class="px-1.5 py-0.2 rounded text-[9.5px] font-mono bg-slate-900 text-purple-300 border border-purple-900/60">
            {{ media.width }} × {{ media.height }}
          </span>
          <span class="px-2 py-0.5 rounded text-[10px] font-mono font-bold bg-purple-950 text-purple-300 border border-purple-800">
            {{ media.download_status }}
          </span>
        </div>
        <div class="flex items-center gap-1.5 shrink-0 font-mono text-[10px]">
          <span class="px-2 py-0.5 rounded-full bg-slate-900 text-slate-300 border border-slate-700">
            候補 {{ media.candidate_tasks?.length ?? 0 }} 本
          </span>
          <button @click.stop="showMediaJson = !showMediaJson" class="px-2 py-0.5 rounded bg-slate-800 hover:bg-slate-700 text-slate-400 hover:text-slate-200 cursor-pointer">
            { } 作品JSON
          </button>
        </div>
      </div>

      <!-- サブ情報行: アカウント、投稿ID、元URL、失敗理由 -->
      <div class="flex flex-wrap items-center gap-x-3 gap-y-1 text-[10.5px] font-mono text-slate-400">
        <span v-if="media.account_id" class="text-indigo-400/90 truncate max-w-[200px]" :title="media.account_id">👤 {{ media.account_id }}</span>
        <span v-if="media.article_id" class="text-slate-500">PostID: {{ media.article_id }}</span>
        <span v-if="displayStatus" class="text-amber-400">状態: {{ displayStatus }}</span>
        <span v-if="media.original_url" class="text-slate-500 truncate max-w-[240px]" :title="media.original_url">URL: {{ media.original_url }}</span>
        <button v-if="media.original_url" @click.stop="copyText(media.original_url)" class="text-[9px] text-slate-400 hover:text-slate-200 underline cursor-pointer">URLコピー</button>
      </div>

      <!-- 作品全体の生JSON展開 -->
      <div v-if="showMediaJson" @click.stop class="mt-2 p-2 bg-slate-950 rounded-xl border border-slate-800 text-[10px] font-mono text-slate-300 overflow-x-auto select-text">
        <pre>{{ JSON.stringify(media, null, 2) }}</pre>
      </div>
    </div>

    <!-- 子要素: 候補URLタスク一覧 -->
    <div v-show="isOpen" class="p-3 space-y-2 bg-slate-950/40">
      <div v-if="!media.candidate_tasks || media.candidate_tasks.length === 0" class="text-xs text-slate-600 font-mono py-2 text-center">
        候補タスクはまだ登録されていません
      </div>
      <PipelineCandidateTaskRow
        v-for="(t, idx) in media.candidate_tasks"
        :key="t.id || idx"
        :task="t"
        :index="idx + 1"
      />
    </div>
  </div>
</template>
