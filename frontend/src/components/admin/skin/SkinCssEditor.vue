<!-- frontend/src/components/admin/skin/SkinCssEditor.vue (100行以下) -->
<script setup lang="ts">
import { computed } from 'vue';
const props = defineProps<{
  css: string;
  saving: boolean;
  isDirty: boolean;
}>();
const emit = defineEmits<{
  (e: 'update:css', val: string): void;
  (e: 'save'): void;
  (e: 'reset'): void;
}>();

const lineCount = computed(() => (props.css ? props.css.split('\n').length : 1));
</script>

<template>
  <div class="bg-slate-900/60 p-4 border border-slate-800 rounded-xl space-y-3">
    <div class="flex items-center justify-between">
      <div class="flex items-center gap-2">
        <h4 class="text-xs font-bold text-slate-200">🎨 design.css エディタ</h4>
        <span class="text-[10px] font-mono px-2 py-0.5 rounded border" :class="lineCount <= 100 ? 'bg-emerald-950/80 text-emerald-400 border-emerald-700/50' : 'bg-rose-950/80 text-rose-400 border-rose-700/50'">
          {{ lineCount }} / 100行 {{ lineCount <= 100 ? '(健全)' : '(超過)' }}
        </span>
      </div>
      <div class="flex items-center gap-2">
        <button @click="$emit('reset')" class="px-2.5 py-1 bg-slate-800 hover:bg-slate-700 text-slate-400 text-xs rounded">再読込</button>
        <button @click="$emit('save')" :disabled="saving || !isDirty" class="px-3.5 py-1 bg-blue-600 hover:bg-blue-500 text-white text-xs font-bold rounded disabled:opacity-50">
          {{ saving ? '保存中...' : '💾 保存' }}
        </button>
      </div>
    </div>
    <textarea
      :value="css"
      @input="$emit('update:css', ($event.target as HTMLTextAreaElement).value)"
      rows="12"
      class="w-full bg-slate-950 border border-slate-700 rounded-lg p-3 text-xs font-mono text-emerald-300 leading-relaxed focus:border-blue-500 focus:outline-none"
    ></textarea>
  </div>
</template>
