<!-- frontend/src/components/admin/jobs/MissionReportModal.vue (100行以下 - SPEC-PRINCIPLE-001) -->
<script setup lang="ts">
import { ref } from 'vue';
import { FileText, Copy, Check, X, Download } from 'lucide-vue-next';

const props = defineProps<{ isOpen: boolean; report: any }>();
const emit = defineEmits<{ (e: 'close'): void }>();
const copied = ref(false);

const copyText = async () => {
  if (!props.report?.markdown_text) return;
  try {
    await navigator.clipboard.writeText(props.report.markdown_text);
    copied.value = true;
    setTimeout(() => { copied.value = false; }, 2000);
  } catch (_) {}
};

const downloadFile = () => {
  if (!props.report?.markdown_text) return;
  const blob = new Blob([props.report.markdown_text], { type: 'text/markdown;charset=utf-8;' });
  const url = URL.createObjectURL(blob);
  const link = document.createElement('a');
  link.href = url;
  link.download = props.report.file_name || 'mission_report.md';
  link.click();
  URL.revokeObjectURL(url);
};
</script>

<template>
  <div v-if="isOpen && report" class="fixed inset-0 z-50 flex items-center justify-center bg-black/80 backdrop-blur-sm p-3 font-sans">
    <div class="bg-slate-900 border border-slate-700 rounded-2xl w-full max-w-3xl max-h-[85vh] flex flex-col overflow-hidden shadow-2xl">
      <!-- ヘッダー -->
      <div class="px-4 py-3 bg-slate-950 border-b border-slate-800 flex items-center justify-between">
        <div class="flex items-center gap-2">
          <FileText class="w-4 h-4 text-blue-400" />
          <h3 class="text-xs font-bold text-slate-100">📋 ミッション完遂レポート (5W1H 証跡)</h3>
          <span class="text-[10px] font-mono text-slate-400 bg-slate-800 px-2 py-0.5 rounded">{{ report.file_name || report.id }}</span>
        </div>
        <div class="flex items-center gap-2">
          <button @click="copyText" class="px-2.5 py-1 bg-slate-800 hover:bg-slate-700 text-slate-200 text-xs font-semibold rounded flex items-center gap-1.5 cursor-pointer active:scale-95">
            <Check v-if="copied" class="w-3.5 h-3.5 text-emerald-400" /><Copy v-else class="w-3.5 h-3.5" />
            <span>{{ copied ? 'コピー完了' : 'コピー' }}</span>
          </button>
          <button @click="downloadFile" class="px-2.5 py-1 bg-blue-600 hover:bg-blue-500 text-white text-xs font-semibold rounded flex items-center gap-1.5 cursor-pointer active:scale-95">
            <Download class="w-3.5 h-3.5" /><span>保存 (.md)</span>
          </button>
          <button @click="emit('close')" class="p-1 rounded text-slate-400 hover:text-white cursor-pointer"><X class="w-4 h-4" /></button>
        </div>
      </div>

      <!-- Markdown テキスト本文 -->
      <div class="flex-1 min-h-0 overflow-y-auto p-4 bg-slate-950/80 font-mono text-xs text-slate-300 whitespace-pre-wrap leading-relaxed select-text">
        {{ report.markdown_text || 'レポート内容がありません。' }}
      </div>
    </div>
  </div>
</template>
