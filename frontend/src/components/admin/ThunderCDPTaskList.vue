<!-- frontend/src/components/admin/ThunderCDPTaskList.vue (100行以下 - SPEC-PRINCIPLE-001) -->
<script setup lang="ts">
defineProps<{ status: any; loading?: boolean; tab?: 'running' | 'all' }>();
</script>

<template>
  <div class="flex-1 bg-slate-900/90 border border-slate-800 rounded-2xl p-3 flex flex-col space-y-2 shadow-lg overflow-hidden min-h-[300px]">
    <div class="flex-1 overflow-y-auto space-y-2 pr-1 font-mono text-xs">
      <div v-if="!status?.recent_tasks || status.recent_tasks.length === 0" class="text-center py-12 text-slate-500 text-xs font-mono">
        現在投入中のタスクはありません。「救出開始」でエスカレーションを開始できます
      </div>

      <!-- MotrixTaskCard と完全に統一されたカードデザイン -->
      <div
        v-for="t in status?.recent_tasks || []"
        :key="t.id"
        class="p-2.5 bg-slate-950 rounded-xl border border-slate-800/80 flex items-center justify-between gap-3 hover:border-purple-500/50 hover:bg-slate-900/40 transition"
      >
        <div class="space-y-1 truncate flex-1 min-w-0">
          <div class="flex items-center gap-2">
            <span class="text-slate-200 font-bold truncate text-xs">{{ t.media_id }}</span>
            <span class="px-1.5 py-0.2 rounded text-[9px] font-bold bg-indigo-950 text-indigo-300 border border-indigo-800 shrink-0">
              {{ t.resolution_type }}
            </span>
          </div>
          <div class="text-[10px] text-slate-400 truncate font-sans">{{ t.url }}</div>
        </div>

        <div class="flex items-center gap-2 shrink-0">
          <span v-if="t.status === 'downloaded' || t.status === 'completed'" class="px-2.5 py-1 rounded-lg text-[10px] font-bold bg-emerald-950 text-emerald-300 border border-emerald-800 flex items-center gap-1">
            ✅ 救出完了
          </span>
          <span v-else-if="t.status === 'reaped'" class="px-2.5 py-1 rounded-lg text-[10px] font-bold bg-slate-800 text-slate-400 border border-slate-700 flex items-center gap-1">
            ♻️ 重複解消
          </span>
          <span v-else-if="t.status === 'depleted'" class="px-2.5 py-1 rounded-lg text-[10px] font-bold bg-rose-950 text-rose-300 border border-rose-800 flex items-center gap-1">
            ⚠️ 枯渇退避
          </span>
          <span v-else class="px-2.5 py-1 rounded-lg text-[10px] font-bold bg-purple-950 text-purple-300 border border-purple-800 flex items-center gap-1 animate-pulse">
            ⚡ 投入中 (Slot {{ t.slot_index >= 0 ? t.slot_index + 1 : '-' }})
          </span>
        </div>
      </div>
    </div>
  </div>
</template>
