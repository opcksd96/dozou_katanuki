<!-- frontend/src/components/admin/MotrixReserveCard.vue (100行以下 - SPEC-PRINCIPLE-001) -->
<script setup lang="ts">
defineProps<{
  reserve: any;
}>();
const emit = defineEmits<{
  (e: 'restore', id: number): void;
}>();
</script>

<template>
  <div class="p-2.5 bg-slate-950 rounded-xl border border-slate-800/80 space-y-1.5 flex items-center justify-between gap-2">
    <div class="min-w-0 flex-1 space-y-0.5">
      <div class="flex items-center gap-2">
        <span class="font-mono font-bold text-xs text-slate-200 truncate">{{ reserve.file_name || reserve.gid }}</span>
        <span class="px-1.5 py-0.2 rounded text-[9px] font-mono bg-emerald-950 text-emerald-300 border border-emerald-800">{{ reserve.status }}</span>
      </div>
      <div class="text-[10px] text-slate-400 font-mono truncate">{{ reserve.url }}</div>
      <div class="text-[9px] text-slate-500 font-mono">{{ reserve.reason }} | 退避日時: {{ new Date(reserve.created_at).toLocaleString() }}</div>
    </div>
    <button @click="emit('restore', reserve.id)" class="px-2.5 py-1 bg-blue-600 hover:bg-blue-500 text-white rounded text-xs font-bold shrink-0 cursor-pointer active:scale-95">
      🚀 再投入
    </button>
  </div>
</template>
