<!-- frontend/src/components/admin/SkinFontEditor.vue (100行以下 - SPEC-PLUGIN-001) -->
<script setup lang="ts">
import { ref, watch, onMounted } from 'vue';
import SkinCssEditor from './skin/SkinCssEditor.vue';
import SkinCardPreview from './skin/SkinCardPreview.vue';

const props = defineProps<{
  skinCSS: string;
  loadingSkin: boolean;
  savingSkin: boolean;
  skinStatus: { success: boolean; message: string } | null;
  selectedPlatform: string;
}>();

const emit = defineEmits<{
  (e: 'fetchSkin', platform: string): void;
  (e: 'saveSkin', platform: string, css: string): void;
  (e: 'applyDynamicSkin', css: string): void;
}>();

const localCSS = ref(props.skinCSS || '');
const isDirty = ref(false);

watch(() => props.skinCSS, (val) => { localCSS.value = val || ''; isDirty.value = false; }, { immediate: true });

onMounted(() => {
  if (!props.skinCSS) emit('fetchSkin', props.selectedPlatform || 'twitter');
});

const onUpdateCSS = (val: string) => {
  localCSS.value = val;
  isDirty.value = true;
  emit('applyDynamicSkin', val);
};
</script>

<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h3 class="text-base font-bold text-slate-100 flex items-center gap-2">
          <span>🧩</span> プラグイン・スキン設定 (Plugins & Skin)
          <span class="text-[10px] font-mono bg-purple-950/80 text-purple-400 border border-purple-700/50 px-2 py-0.5 rounded">SPEC-PLUGIN-001</span>
        </h3>
        <p class="text-xs text-slate-400 mt-0.5">プラットフォーム（Twitter/Bluesky等）ごとの表示スキン (design.css) のカスタマイズとプレビューを行います。</p>
      </div>
    </div>

    <div v-if="skinStatus" class="p-2.5 rounded-lg text-xs font-semibold" :class="skinStatus.success ? 'bg-emerald-950/70 border border-emerald-500/40 text-emerald-300' : 'bg-rose-950/70 border border-rose-500/40 text-rose-300'">
      {{ skinStatus.message }}
    </div>

    <SkinCssEditor
      :css="localCSS"
      :saving="savingSkin"
      :is-dirty="isDirty"
      @update:css="onUpdateCSS"
      @save="() => { emit('saveSkin', selectedPlatform || 'twitter', localCSS); isDirty = false; }"
      @reset="() => emit('fetchSkin', selectedPlatform || 'twitter')"
    />

    <SkinCardPreview :platform="selectedPlatform || 'twitter'" />
  </div>
</template>

