<!-- frontend/src/components/admin/jobs/JobTerminalView.vue (100行以下) -->
<script setup lang="ts">
import { ref, watch, nextTick } from 'vue';
const props = defineProps<{ logs: string[] }>();
defineEmits<{ (e: 'clearLogs'): void }>();
const autoScroll = ref(true);
const terminalRef = ref<HTMLElement | null>(null);

watch(() => props.logs.length, () => {
  if (autoScroll.value && terminalRef.value) {
    nextTick(() => { if (terminalRef.value) terminalRef.value.scrollTop = terminalRef.value.scrollHeight; });
  }
});
const copyLogs = () => {
  navigator.clipboard.writeText(props.logs.join('\n')).then(() => alert('ログをコピーしました'));
};
</script>

<template>
  <div class="bg-black/90 border border-slate-800 rounded-xl overflow-hidden">
    <div class="bg-slate-900 px-4 py-2 border-b border-slate-800 flex items-center justify-between text-xs">
      <div class="flex items-center gap-2">
        <span class="w-2.5 h-2.5 rounded-full bg-red-500/80"></span>
        <span class="w-2.5 h-2.5 rounded-full bg-yellow-500/80"></span>
        <span class="w-2.5 h-2.5 rounded-full bg-green-500/80"></span>
        <span class="ml-2 font-mono text-slate-400 font-semibold">Scraper View Terminal</span>
      </div>
      <div class="flex items-center gap-3 text-slate-400">
        <label class="flex items-center gap-1 cursor-pointer text-[11px]">
          <input type="checkbox" v-model="autoScroll" class="rounded bg-slate-950 border-slate-700" />
          <span>自動追従</span>
        </label>
        <button @click="copyLogs" class="hover:text-slate-200">📋 コピー</button>
        <button @click="$emit('clearLogs')" class="hover:text-slate-200">🧹 クリア</button>
      </div>
    </div>
    <div ref="terminalRef" class="p-3 h-52 overflow-y-auto font-mono text-xs text-emerald-400 space-y-0.5 bg-black/60">
      <div v-if="logs.length === 0" class="text-slate-600 italic">--- 待機中 / ログはありません ---</div>
      <div v-for="(line, idx) in logs" :key="idx" class="whitespace-pre-wrap break-all px-1 rounded">
        <span class="text-slate-600 select-none mr-2 text-[10px]">{{ idx + 1 }}</span>
        <span>{{ line }}</span>
      </div>
    </div>
  </div>
</template>
