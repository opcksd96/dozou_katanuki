<!-- frontend/src/components/admin/ThunderSlotCard.vue (100行以下 - SPEC-PRINCIPLE-001) -->
<script setup lang="ts">
defineProps<{
  slot: {
    index: number;
    is_occupied: boolean;
    current_task?: {
      media_id: string;
      resolution_type: string;
    };
  };
}>();
</script>

<template>
  <div :class="slot.is_occupied ? 'border-purple-500/80 bg-purple-950/30' : 'border-slate-800 bg-slate-950/60'" class="p-2 rounded-xl border flex flex-col justify-between min-h-[70px] transition-all">
    <div class="flex items-center justify-between">
      <span class="text-[9px] font-mono font-bold text-slate-400">Slot #{{ slot.index + 1 }}</span>
      <span :class="slot.is_occupied ? 'bg-purple-900 text-purple-200' : 'bg-slate-800 text-slate-500'" class="px-1 py-0.2 rounded text-[8px] font-mono font-bold">
        {{ slot.is_occupied ? 'ACTIVE' : 'IDLE' }}
      </span>
    </div>
    <div v-if="slot.is_occupied && slot.current_task" class="space-y-0.5">
      <div class="text-[10px] font-mono font-bold text-purple-300 truncate">{{ slot.current_task.media_id }}</div>
      <div class="text-[8px] font-mono text-amber-300">{{ slot.current_task.resolution_type }}</div>
    </div>
    <div v-else class="text-[9px] text-slate-600 font-mono text-center py-2">待機中</div>
  </div>
</template>
