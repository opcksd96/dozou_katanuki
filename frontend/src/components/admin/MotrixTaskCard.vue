<!-- frontend/src/components/admin/MotrixTaskCard.vue (100行以下 - SPEC-PRINCIPLE-001) -->
<script setup lang="ts">
defineProps<{
  task: any;
  isSelected: boolean;
}>();
const emit = defineEmits<{
  (e: 'toggle', gid: string): void;
}>();

const formatSpeed = (b: number) => {
  if (!b || b === 0) return '0 B/s';
  return b < 1024 * 1024 ? `${(b / 1024).toFixed(1)} KB/s` : `${(b / (1024 * 1024)).toFixed(2)} MB/s`;
};
</script>

<template>
  <div class="p-2.5 bg-slate-950 rounded-xl border border-slate-800/80 space-y-1.5 flex items-start gap-2.5">
    <input type="checkbox" :checked="isSelected" @change="emit('toggle', task.gid)" class="mt-1 rounded bg-slate-800 border-slate-700 text-blue-500 cursor-pointer" />
    <div class="flex-1 min-w-0 space-y-1">
      <div class="flex items-center justify-between text-xs gap-2">
        <span class="font-mono font-bold text-slate-200 truncate flex-1">{{ task.file_name || task.gid }}</span>
        <div class="flex items-center gap-2 shrink-0">
          <span class="font-mono text-purple-400 text-[11px]">{{ formatSpeed(task.download_speed) }}</span>
          <span :class="task.status === 'active' ? 'bg-blue-950 text-blue-300 border-blue-800' : task.status === 'paused' ? 'bg-amber-950 text-amber-300 border-amber-800' : task.status === 'error' ? 'bg-rose-950 text-rose-300 border-rose-800' : 'bg-emerald-950 text-emerald-300 border-emerald-800'" class="px-1.5 py-0.2 rounded text-[9px] font-mono border">
            {{ task.status.toUpperCase() }}
          </span>
        </div>
      </div>
      <div class="w-full bg-slate-900 rounded-full h-1.5 overflow-hidden">
        <div class="bg-gradient-to-r from-blue-500 to-purple-500 h-full rounded-full" :style="{ width: `${task.progress}%` }"></div>
      </div>
      <div class="flex items-center justify-between text-[10px] text-slate-500 font-mono">
        <span class="truncate max-w-[280px]">GID: {{ task.gid }}</span>
        <span>{{ task.progress.toFixed(1) }}% ({{ (task.completed_length / (1024*1024)).toFixed(1) }} / {{ (task.total_length / (1024*1024)).toFixed(1) }} MB)</span>
      </div>
      <div v-if="task.error_message" class="text-[10px] text-rose-400 font-mono bg-rose-950/40 p-1 rounded border border-rose-900/40">
        ⚠️ {{ task.error_message }}
      </div>
    </div>
  </div>
</template>
