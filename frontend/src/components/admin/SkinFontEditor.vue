<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue';
import { models } from '../../../wailsjs/go/models';

const props = defineProps<{
  skinCSS: string;
  loadingSkin: boolean;
  savingSkin: boolean;
  skinStatus: { success: boolean; message: string } | null;
  selectedPlatform: string;
  fontPresets: {
    ja: Array<{ label: string; value: string }>;
    en: Array<{ label: string; value: string }>;
    zh: Array<{ label: string; value: string }>;
  };
  config: models.AppConfig | null;
  savingConfig: boolean;
  saveStatus: { success: boolean; message: string } | null;
}>();

const emit = defineEmits<{
  (e: 'fetchSkin', platform: string): void;
  (e: 'saveSkin', platform: string, css: string): void;
  (e: 'applyDynamicSkin', css: string): void;
  (e: 'applyFontVariables', cfg?: models.AppConfig | null): void;
  (e: 'saveConfig'): void;
}>();

// エディタのローカル編集状態
const localCSS = ref<string>('');
const isDirty = ref(false);
const autoApplyPreview = ref(true);

// フォントのローカル編集状態
const localFontJa = ref<string>('');
const localFontEn = ref<string>('');
const localFontZh = ref<string>('');

// プラットフォーム一覧
const platforms = [
  { id: 'twitter', name: 'Twitter / X (標準)' },
  { id: 'bsky', name: 'Bluesky (拡張予定)' },
];

const currentPlatform = ref(props.selectedPlatform || 'twitter');

// props の skinCSS 変化時にローカル状態へ同期
watch(
  () => props.skinCSS,
  (val) => {
    localCSS.value = val || '';
    isDirty.value = false;
  },
  { immediate: true }
);

// props の config 変化時にフォントローカル状態へ同期
watch(
  () => props.config,
  (cfg) => {
    if (cfg?.appearance) {
      localFontJa.value = cfg.appearance.font_family_ja || '';
      localFontEn.value = cfg.appearance.font_family_en || '';
      localFontZh.value = cfg.appearance.font_family_zh || '';
    }
  },
  { immediate: true }
);

// CSS行数カウント（SPEC-PLUGIN-001 100行以下厳守チェック）
const lineCount = computed(() => {
  if (!localCSS.value) return 1;
  return localCSS.value.split('\n').length;
});

const onCodeInput = () => {
  isDirty.value = true;
  if (autoApplyPreview.value) {
    emit('applyDynamicSkin', localCSS.value);
  }
};

const handlePlatformChange = () => {
  if (isDirty.value) {
    if (!confirm('未保存の変更があります。プラットフォームを切り替えますか？')) {
      return;
    }
  }
  emit('fetchSkin', currentPlatform.value);
};

const handleSaveSkin = () => {
  emit('saveSkin', currentPlatform.value, localCSS.value);
  isDirty.value = false;
};

const handleResetSkin = () => {
  if (isDirty.value && !confirm('編集内容を破棄して最新のファイルを再読み込みしますか？')) {
    return;
  }
  emit('fetchSkin', currentPlatform.value);
};

// タブキー入力対応
const handleTabKey = (e: KeyboardEvent) => {
  if (e.key === 'Tab') {
    e.preventDefault();
    const textarea = e.target as HTMLTextAreaElement;
    const start = textarea.selectionStart;
    const end = textarea.selectionEnd;

    localCSS.value = localCSS.value.substring(0, start) + '  ' + localCSS.value.substring(end);
    isDirty.value = true;
    if (autoApplyPreview.value) {
      emit('applyDynamicSkin', localCSS.value);
    }

    setTimeout(() => {
      textarea.selectionStart = textarea.selectionEnd = start + 2;
    }, 0);
  } else if ((e.ctrlKey || e.metaKey) && e.key === 's') {
    e.preventDefault();
    handleSaveSkin();
  }
};

// フォントのプリセット適用
const applyPreset = (lang: 'ja' | 'en' | 'zh', value: string) => {
  if (lang === 'ja') localFontJa.value = value;
  if (lang === 'en') localFontEn.value = value;
  if (lang === 'zh') localFontZh.value = value;
  syncFontToConfig();
};

// フォント入力を config と CSS 変数へ同期
const syncFontToConfig = () => {
  if (props.config?.appearance) {
    props.config.appearance.font_family_ja = localFontJa.value;
    props.config.appearance.font_family_en = localFontEn.value;
    props.config.appearance.font_family_zh = localFontZh.value;
    emit('applyFontVariables', props.config);
  }
};

onMounted(() => {
  if (!props.skinCSS) {
    emit('fetchSkin', currentPlatform.value);
  }
});
</script>

<template>
  <div class="space-y-6">
    <!-- ヘッダーガイド -->
    <div class="flex items-center justify-between bg-slate-900/60 p-3 rounded-xl border border-slate-800">
      <div class="flex items-center gap-2">
        <span class="text-base">🎨</span>
        <div>
          <h3 class="text-xs font-bold text-slate-200">
            CSSスキンエディタ ＆ フォント微調整パネル
            <span class="ml-1 text-[10px] font-mono text-blue-400 bg-blue-950/60 px-1.5 py-0.5 rounded border border-blue-800/60">
              SPEC-ADMINBOARD-001 (第5・第6ピース)
            </span>
          </h3>
          <p class="text-[11px] text-slate-400">
            プラットフォーム固有のCSS直接物理保存と、日・英・中の多言語フォントシグナル同期
          </p>
        </div>
      </div>
    </div>

    <!-- 1. CSS スキンエディタ (第5ピース) -->
    <div class="bg-slate-900/80 border border-slate-800 rounded-xl p-4 space-y-4">
      <div class="flex flex-wrap items-center justify-between gap-3 pb-3 border-b border-slate-800">
        <div class="flex items-center gap-3">
          <h4 class="text-sm font-bold text-slate-200 flex items-center gap-1.5">
            <span>📝</span> デフォルトCSSコードエディタ
          </h4>
          <!-- プラットフォーム選択 -->
          <div class="flex items-center gap-1.5 text-xs text-slate-400">
            <span>対象:</span>
            <select
              v-model="currentPlatform"
              @change="handlePlatformChange"
              class="bg-slate-950 border border-slate-700 rounded px-2 py-1 text-xs text-slate-200 font-mono focus:outline-none focus:border-blue-500"
            >
              <option v-for="p in platforms" :key="p.id" :value="p.id">
                {{ p.name }}
              </option>
            </select>
          </div>
        </div>

        <div class="flex items-center gap-2">
          <!-- プレビュー即時適用トグル -->
          <label class="flex items-center gap-1.5 text-[11px] text-slate-400 cursor-pointer select-none">
            <input
              type="checkbox"
              v-model="autoApplyPreview"
              class="rounded bg-slate-950 border-slate-700 text-blue-600 focus:ring-0"
            />
            <span>即時プレビュー反映</span>
          </label>

          <!-- 行数制限インジケータ -->
          <span
            class="text-[10px] font-mono px-2 py-0.5 rounded border"
            :class="
              lineCount <= 100
                ? 'bg-emerald-950/60 text-emerald-300 border-emerald-800/60'
                : 'bg-amber-950/60 text-amber-300 border-amber-800/60'
            "
            title="SPEC-PLUGIN-001 原則: プラグインCSSは100行以下厳守"
          >
            {{ lineCount }} / 100 行 {{ lineCount > 100 ? '(推奨超過)' : '' }}
          </span>

          <button
            type="button"
            @click="handleResetSkin"
            :disabled="loadingSkin || savingSkin"
            class="px-2.5 py-1 bg-slate-800 hover:bg-slate-700 text-slate-300 rounded text-xs transition-colors flex items-center gap-1 font-medium"
            title="最新のファイルを再読み込み"
          >
            <span>🔄</span> 再取得
          </button>

          <button
            type="button"
            @click="handleSaveSkin"
            :disabled="savingSkin"
            class="px-3.5 py-1 bg-blue-600 hover:bg-blue-500 disabled:bg-slate-700 text-white rounded text-xs font-bold transition-colors shadow-sm flex items-center gap-1 cursor-pointer"
            title="ファイルを上書き保存 (Ctrl+S)"
          >
            <span v-if="savingSkin" class="animate-spin">⏳</span>
            <span v-else>💾</span>
            <span>保存 (Ctrl+S)</span>
          </button>
        </div>
      </div>

      <!-- ステータス通知 -->
      <div
        v-if="skinStatus"
        class="p-2.5 rounded-lg text-xs font-semibold flex items-center justify-between transition-all"
        :class="
          skinStatus.success
            ? 'bg-emerald-950/80 border border-emerald-500/50 text-emerald-300'
            : 'bg-rose-950/80 border border-rose-500/50 text-rose-300'
        "
      >
        <span>{{ skinStatus.message }}</span>
        <span class="text-[10px] font-mono opacity-75">plugins/{{ currentPlatform }}/skin/design.css</span>
      </div>

      <!-- コードエディタ領域 -->
      <div class="relative rounded-lg border border-slate-800 bg-slate-950 font-mono text-xs overflow-hidden">
        <!-- エディタヘッダーバナー -->
        <div class="px-3 py-1.5 bg-slate-900 border-b border-slate-800 flex items-center justify-between text-[11px] text-slate-400">
          <span class="flex items-center gap-1.5">
            <span class="w-2 h-2 rounded-full" :class="isDirty ? 'bg-amber-400' : 'bg-emerald-400'"></span>
            <span>plugins/{{ currentPlatform }}/skin/design.css</span>
            <span v-if="isDirty" class="text-amber-400 font-bold text-[10px]">(未保存の変更あり)</span>
          </span>
          <span class="text-[10px] text-slate-500">CSS3 • UTF-8</span>
        </div>

        <div v-if="loadingSkin" class="py-20 text-center text-slate-500">
          <span class="animate-spin text-lg inline-block mr-2">⏳</span> スキンCSSを読み込み中...
        </div>

        <textarea
          v-else
          v-model="localCSS"
          @input="onCodeInput"
          @keydown="handleTabKey"
          spellcheck="false"
          placeholder="/* CSSスタイルをここに入力... */"
          class="w-full h-72 p-3.5 bg-transparent text-slate-200 font-mono text-xs leading-relaxed resize-y focus:outline-none focus:ring-1 focus:ring-blue-500/50 selection:bg-blue-900 selection:text-white"
        ></textarea>
      </div>

      <p class="text-[11px] text-slate-400 flex items-center gap-1">
        <span>💡</span>
        <span>エディタ内で編集すると即座に画面へプレビュー反映され、<strong>Ctrl+S</strong> または保存ボタンで実ファイルへ物理保存されます。</span>
      </p>
    </div>

    <!-- 2. フォント微調整パネル (第6ピース) -->
    <div class="bg-slate-900/80 border border-slate-800 rounded-xl p-4 space-y-4">
      <div class="flex items-center justify-between pb-3 border-b border-slate-800">
        <h4 class="text-sm font-bold text-slate-200 flex items-center gap-1.5">
          <span>🔤</span> 多言語フォント微調整パネル
          <span class="text-[10px] font-mono font-normal text-slate-400 bg-slate-800 px-1.5 py-0.5 rounded">
            CSSシグナル同期
          </span>
        </h4>

        <div class="flex items-center gap-2">
          <button
            type="button"
            @click="syncFontToConfig"
            class="px-2.5 py-1 bg-slate-800 hover:bg-slate-700 text-slate-300 rounded text-xs transition-colors font-medium flex items-center gap-1"
            title="CSSカスタム変数へ即時反映"
          >
            <span>⚡</span> 即時シグナル同期
          </button>
          <button
            type="button"
            @click="emit('saveConfig')"
            :disabled="savingConfig"
            class="px-3.5 py-1 bg-blue-600 hover:bg-blue-500 disabled:bg-slate-700 text-white rounded text-xs font-bold transition-colors shadow-sm flex items-center gap-1 cursor-pointer"
            title="config.json にフォント設定を永続保存"
          >
            <span v-if="savingConfig" class="animate-spin">⏳</span>
            <span v-else>💾</span>
            <span>設定を保存 (config.json)</span>
          </button>
        </div>
      </div>

      <!-- ステータス通知 -->
      <div
        v-if="saveStatus"
        class="p-2.5 rounded-lg text-xs font-semibold flex items-center justify-between transition-all"
        :class="
          saveStatus.success
            ? 'bg-emerald-950/80 border border-emerald-500/50 text-emerald-300'
            : 'bg-rose-950/80 border border-rose-500/50 text-rose-300'
        "
      >
        <span>{{ saveStatus.message }}</span>
      </div>

      <!-- 3言語のフォント設定グリッド -->
      <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
        <!-- 1. 日本語 (JA) -->
        <div class="bg-slate-950/70 border border-slate-800/90 rounded-lg p-3 space-y-2.5">
          <div class="flex items-center justify-between">
            <span class="text-xs font-bold text-slate-200 flex items-center gap-1">
              <span>🇯🇵</span> 日本語フォント (JA)
            </span>
            <span class="text-[10px] font-mono text-slate-500">--font-family-ja</span>
          </div>

          <!-- プリセット選択 -->
          <div>
            <label class="block text-[10px] text-slate-400 mb-1">プリセット選択:</label>
            <div class="flex flex-wrap gap-1">
              <button
                v-for="p in fontPresets.ja"
                :key="p.label"
                type="button"
                @click="applyPreset('ja', p.value)"
                class="text-[10px] px-2 py-0.5 rounded border transition-colors cursor-pointer"
                :class="
                  localFontJa === p.value
                    ? 'bg-blue-900/60 text-blue-200 border-blue-600 font-semibold'
                    : 'bg-slate-900 text-slate-400 border-slate-700 hover:text-slate-200'
                "
              >
                {{ p.label.split(' ')[0] }}
              </button>
            </div>
          </div>

          <!-- フォントファミリ入力 -->
          <div>
            <label class="block text-[10px] text-slate-400 mb-1">font-family 定義:</label>
            <input
              v-model="localFontJa"
              @input="syncFontToConfig"
              type="text"
              class="w-full bg-slate-900 border border-slate-700 rounded px-2.5 py-1 text-xs text-slate-200 font-mono focus:outline-none focus:border-blue-500"
            />
          </div>

          <!-- プレビューボックス -->
          <div class="p-2.5 bg-slate-900/90 rounded border border-slate-800 text-xs space-y-1">
            <div class="text-[10px] text-slate-500 font-mono">レンダリング・プレビュー:</div>
            <div :style="{ fontFamily: localFontJa }" class="text-slate-100 font-normal leading-relaxed text-sm">
              土蔵・型抜き 2A03 APU 矩形波チャンネルの動態保存テスト。
            </div>
          </div>
        </div>

        <!-- 2. 英語 (EN) -->
        <div class="bg-slate-950/70 border border-slate-800/90 rounded-lg p-3 space-y-2.5">
          <div class="flex items-center justify-between">
            <span class="text-xs font-bold text-slate-200 flex items-center gap-1">
              <span>🇺🇸</span> 英語フォント (EN)
            </span>
            <span class="text-[10px] font-mono text-slate-500">--font-family-en</span>
          </div>

          <!-- プリセット選択 -->
          <div>
            <label class="block text-[10px] text-slate-400 mb-1">プリセット選択:</label>
            <div class="flex flex-wrap gap-1">
              <button
                v-for="p in fontPresets.en"
                :key="p.label"
                type="button"
                @click="applyPreset('en', p.value)"
                class="text-[10px] px-2 py-0.5 rounded border transition-colors cursor-pointer"
                :class="
                  localFontEn === p.value
                    ? 'bg-blue-900/60 text-blue-200 border-blue-600 font-semibold'
                    : 'bg-slate-900 text-slate-400 border-slate-700 hover:text-slate-200'
                "
              >
                {{ p.label.split(' ')[0] }}
              </button>
            </div>
          </div>

          <!-- フォントファミリ入力 -->
          <div>
            <label class="block text-[10px] text-slate-400 mb-1">font-family 定義:</label>
            <input
              v-model="localFontEn"
              @input="syncFontToConfig"
              type="text"
              class="w-full bg-slate-900 border border-slate-700 rounded px-2.5 py-1 text-xs text-slate-200 font-mono focus:outline-none focus:border-blue-500"
            />
          </div>

          <!-- プレビューボックス -->
          <div class="p-2.5 bg-slate-900/90 rounded border border-slate-800 text-xs space-y-1">
            <div class="text-[10px] text-slate-500 font-mono">レンダリング・プレビュー:</div>
            <div :style="{ fontFamily: localFontEn }" class="text-slate-100 font-normal leading-relaxed text-sm">
              Dozou Katanuki v4.0 Wails & Stash integration is running.
            </div>
          </div>
        </div>

        <!-- 3. 中国語 (ZH) -->
        <div class="bg-slate-950/70 border border-slate-800/90 rounded-lg p-3 space-y-2.5">
          <div class="flex items-center justify-between">
            <span class="text-xs font-bold text-slate-200 flex items-center gap-1">
              <span>🇨🇳</span> 中国語フォント (ZH)
            </span>
            <span class="text-[10px] font-mono text-slate-500">--font-family-zh</span>
          </div>

          <!-- プリセット選択 -->
          <div>
            <label class="block text-[10px] text-slate-400 mb-1">プリセット選択:</label>
            <div class="flex flex-wrap gap-1">
              <button
                v-for="p in fontPresets.zh"
                :key="p.label"
                type="button"
                @click="applyPreset('zh', p.value)"
                class="text-[10px] px-2 py-0.5 rounded border transition-colors cursor-pointer"
                :class="
                  localFontZh === p.value
                    ? 'bg-blue-900/60 text-blue-200 border-blue-600 font-semibold'
                    : 'bg-slate-900 text-slate-400 border-slate-700 hover:text-slate-200'
                "
              >
                {{ p.label.split(' ')[0] }}
              </button>
            </div>
          </div>

          <!-- フォントファミリ入力 -->
          <div>
            <label class="block text-[10px] text-slate-400 mb-1">font-family 定義:</label>
            <input
              v-model="localFontZh"
              @input="syncFontToConfig"
              type="text"
              class="w-full bg-slate-900 border border-slate-700 rounded px-2.5 py-1 text-xs text-slate-200 font-mono focus:outline-none focus:border-blue-500"
            />
          </div>

          <!-- プレビューボックス -->
          <div class="p-2.5 bg-slate-900/90 rounded border border-slate-800 text-xs space-y-1">
            <div class="text-[10px] text-slate-500 font-mono">レンダリング・プレビュー:</div>
            <div :style="{ fontFamily: localFontZh }" class="text-slate-100 font-normal leading-relaxed text-sm">
              红白机 2A03 APU 寄存器脉冲通道配置测试记录。
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
