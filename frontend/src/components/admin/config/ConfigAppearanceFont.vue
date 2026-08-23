<!-- frontend/src/components/admin/config/ConfigAppearanceFont.vue (100行以下) -->
<script setup lang="ts">
import { models } from '../../../../wailsjs/go/models';
import { useTheme, ThemeMode } from '../../../composables/useTheme';

const props = defineProps<{ config: models.AppConfig }>();
const { currentTheme, fontScale, setTheme, setFontScale } = useTheme();

const scalePresets = [
  { label: '85%', value: 0.85 },
  { label: '100% (標準)', value: 1.0 },
  { label: '115%', value: 1.15 },
  { label: '130% (拡大)', value: 1.3 },
];

const handleThemeSelect = (theme: ThemeMode) => {
  if (props.config.appearance) {
    props.config.appearance.theme = theme;
  }
  setTheme(theme);
};

const handleScaleSelect = (scale: number) => {
  if (props.config.appearance) {
    props.config.appearance.font_scale = scale;
  }
  setFontScale(scale);
};

const handleSliderChange = (e: Event) => {
  const val = parseFloat((e.target as HTMLInputElement).value);
  handleScaleSelect(val);
};
</script>

<template>
  <div class="bg-slate-900/80 border border-slate-800 rounded-xl p-4 space-y-4">
    <h3 class="text-sm font-bold text-slate-200 flex items-center gap-2">
      <span>🎨</span> 外観・文字サイズ・フォント設定 (Appearance & Typography)
    </h3>

    <!-- 1. 外観モード (Theme) & フォントサイズ倍率 (Font Scale) -->
    <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
      <!-- テーマモード -->
      <div class="space-y-1.5">
        <label class="block text-xs font-semibold text-slate-300">🌙 外観テーマモード (Theme)</label>
        <div class="grid grid-cols-3 gap-2">
          <button
            type="button"
            @click="handleThemeSelect('system')"
            :class="[
              'py-1.5 px-2 rounded-lg text-xs font-medium border transition-all cursor-pointer text-center',
              (config.appearance?.theme || currentTheme) === 'system'
                ? 'bg-blue-600 border-blue-500 text-white font-bold shadow-sm'
                : 'bg-slate-950/80 border-slate-700/80 text-slate-400 hover:text-slate-200'
            ]"
          >
            🖥️ システム
          </button>
          <button
            type="button"
            @click="handleThemeSelect('dark')"
            :class="[
              'py-1.5 px-2 rounded-lg text-xs font-medium border transition-all cursor-pointer text-center',
              (config.appearance?.theme || currentTheme) === 'dark'
                ? 'bg-blue-600 border-blue-500 text-white font-bold shadow-sm'
                : 'bg-slate-950/80 border-slate-700/80 text-slate-400 hover:text-slate-200'
            ]"
          >
            🌙 ダーク
          </button>
          <button
            type="button"
            @click="handleThemeSelect('light')"
            :class="[
              'py-1.5 px-2 rounded-lg text-xs font-medium border transition-all cursor-pointer text-center',
              (config.appearance?.theme || currentTheme) === 'light'
                ? 'bg-blue-600 border-blue-500 text-white font-bold shadow-sm'
                : 'bg-slate-950/80 border-slate-700/80 text-slate-400 hover:text-slate-200'
            ]"
          >
            ☀️ ライト
          </button>
        </div>
      </div>

      <!-- フォントサイズ倍率 (Font Scale) -->
      <div class="space-y-1.5">
        <div class="flex items-center justify-between">
          <label class="block text-xs font-semibold text-slate-300">🔍 文字サイズ倍率 (Font Scale)</label>
          <span class="text-xs font-mono font-bold text-blue-400">{{ Math.round(fontScale * 100) }}%</span>
        </div>
        <div class="flex items-center gap-2">
          <input
            type="range"
            min="0.8"
            max="1.4"
            step="0.05"
            :value="config.appearance?.font_scale || fontScale"
            @input="handleSliderChange"
            class="w-full accent-blue-500 cursor-pointer h-1.5 bg-slate-950 rounded-lg"
          />
        </div>
        <div class="flex items-center justify-between gap-1 pt-0.5">
          <button
            v-for="p in scalePresets"
            :key="p.value"
            type="button"
            @click="handleScaleSelect(p.value)"
            :class="[
              'px-2 py-0.5 rounded text-[10px] font-mono transition-colors cursor-pointer border',
              Math.abs(fontScale - p.value) < 0.03
                ? 'bg-blue-600 border-blue-500 text-white font-bold'
                : 'bg-slate-950 border-slate-800 text-slate-400 hover:text-slate-200'
            ]"
          >
            {{ p.label }}
          </button>
        </div>
      </div>
    </div>

    <!-- 2. 多言語フォント設定 -->
    <div class="space-y-2 pt-2 border-t border-slate-800/80">
      <label class="block text-xs font-semibold text-slate-300">🔤 各言語のフォントファミリー (Font Families)</label>
      <div class="grid grid-cols-1 sm:grid-cols-3 gap-3">
        <div>
          <label class="block text-[10px] text-slate-400 mb-1">🇯🇵 日本語フォント</label>
          <input v-model="config.appearance.font_family_ja" type="text" class="w-full bg-slate-950 border border-slate-700 rounded-lg p-1.5 text-xs text-slate-200 font-mono" />
        </div>
        <div>
          <label class="block text-[10px] text-slate-400 mb-1">🇺🇸 英語フォント</label>
          <input v-model="config.appearance.font_family_en" type="text" class="w-full bg-slate-950 border border-slate-700 rounded-lg p-1.5 text-xs text-slate-200 font-mono" />
        </div>
        <div>
          <label class="block text-[10px] text-slate-400 mb-1">🇨🇳 中国語フォント</label>
          <input v-model="config.appearance.font_family_zh" type="text" class="w-full bg-slate-950 border border-slate-700 rounded-lg p-1.5 text-xs text-slate-200 font-mono" />
        </div>
      </div>
    </div>
  </div>
</template>
