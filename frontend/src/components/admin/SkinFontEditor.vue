<!-- frontend/src/components/admin/SkinFontEditor.vue (100行以下) -->
<script setup lang="ts">
import { ref, watch } from 'vue';
import { models } from '../../../wailsjs/go/models';
import SkinCssEditor from './skin/SkinCssEditor.vue';
import SkinFontSettings from './skin/SkinFontSettings.vue';
import SkinCardPreview from './skin/SkinCardPreview.vue';

const props = defineProps<{
  skinCSS: string;
  loadingSkin: boolean;
  savingSkin: boolean;
  skinStatus: { success: boolean; message: string } | null;
  selectedPlatform: string;
  fontPresets: { ja: Array<{ label: string; value: string }>; en: Array<{ label: string; value: string }>; zh: Array<{ label: string; value: string }> };
  config: models.AppConfig | null;
  savingConfig: boolean;
  saveStatus: { success: boolean; message: string } | null;
}>();

const emit = defineEmits<{
  (e: 'fetchSkin', platform: string): void;
  (e: 'saveSkin', platform: string, css: string): void;
  (e: 'applyDynamicSkin', css: string): void;
  (e: 'saveConfig'): void;
}>();

const localCSS = ref(props.skinCSS || '');
const isDirty = ref(false);

watch(() => props.skinCSS, (val) => { localCSS.value = val || ''; isDirty.value = false; }, { immediate: true });

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
          <span>🎨</span> スキンCSS ＆ グローバルフォント設定
          <span class="text-[10px] font-mono bg-purple-950/80 text-purple-400 border border-purple-700/50 px-2 py-0.5 rounded">SPEC-PLUGIN-001</span>
        </h3>
        <p class="text-xs text-slate-400 mt-0.5">タイムラインのカードデザイン (design.css) と各言語の表示フォントをカスタマイズします。</p>
      </div>
    </div>

    <div v-if="skinStatus || saveStatus" class="p-2.5 rounded-lg text-xs font-semibold" :class="(skinStatus || saveStatus)?.success ? 'bg-emerald-950/70 border border-emerald-500/40 text-emerald-300' : 'bg-rose-950/70 border border-rose-500/40 text-rose-300'">
      {{ (skinStatus || saveStatus)?.message }}
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

    <SkinFontSettings
      :config="config"
      :font-presets="fontPresets"
      :saving-config="savingConfig"
      @save="() => emit('saveConfig')"
    />
  </div>
</template>
