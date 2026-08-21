<!-- frontend/src/components/layout/ToastContainer.vue (100行以下) -->
<script setup lang="ts">
import { useToast } from '../../composables/useToast';

const { toasts, removeToast } = useToast();
</script>

<template>
  <div class="fixed top-5 right-5 z-50 flex flex-col gap-2 pointer-events-none max-w-sm w-full">
    <TransitionGroup
      enter-active-class="transform ease-out duration-300 transition"
      enter-from-class="translate-y-[-10px] opacity-0 scale-95"
      enter-to-class="translate-y-0 opacity-100 scale-100"
      leave-active-class="transition ease-in duration-200"
      leave-from-class="opacity-100 scale-100"
      leave-to-class="opacity-0 scale-90"
    >
      <div
        v-for="toast in toasts"
        :key="toast.id"
        @click="removeToast(toast.id)"
        class="pointer-events-auto flex items-center justify-between gap-3 px-4 py-2.5 rounded-xl shadow-2xl border text-xs font-medium backdrop-blur-md cursor-pointer transition-all hover:scale-[1.02]"
        :class="{
          'bg-emerald-950/90 border-emerald-500/50 text-emerald-200 shadow-emerald-950/50': toast.type === 'success',
          'bg-blue-950/90 border-blue-500/50 text-blue-200 shadow-blue-950/50': toast.type === 'info',
          'bg-amber-950/90 border-amber-500/50 text-amber-200 shadow-amber-950/50': toast.type === 'warning',
          'bg-rose-950/90 border-rose-500/50 text-rose-200 shadow-rose-950/50': toast.type === 'error'
        }"
      >
        <span class="truncate">{{ toast.message }}</span>
        <button type="button" class="text-white/60 hover:text-white text-xs">✕</button>
      </div>
    </TransitionGroup>
  </div>
</template>
