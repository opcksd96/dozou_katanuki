// frontend/src/composables/useSkin.ts (SPEC-PLUGIN-001)
import { ref } from 'vue';
import type { SkinController, SkinContext } from '../models/SkinController';

export interface SkinLayoutConfig {
  name?: string;
  version?: string;
  platform?: string;
  display_name?: string;
  styles?: Record<string, string>;
  shortcuts?: Record<string, any>;
}

export function useSkin() {
  const currentPlatform = ref('twitter');
  const layoutConfig = ref<SkinLayoutConfig | null>(null);
  const controller = ref<SkinController | null>(null);
  const isLoaded = ref(false);

  // 簡易 YAML / Key-Value パーサー（ブラウザ完結）
  const parseSimpleYAML = (yamlStr: string): SkinLayoutConfig => {
    const res: any = { styles: {}, shortcuts: {} };
    let currentSection = '';
    const lines = yamlStr.split('\n');
    for (const rawLine of lines) {
      const line = rawLine.trim();
      if (!line || line.startsWith('#')) continue;
      if (line.endsWith(':') && !line.includes(' ')) {
        currentSection = line.slice(0, -1).trim();
        continue;
      }
      const colonIdx = line.indexOf(':');
      if (colonIdx > 0) {
        const key = line.slice(0, colonIdx).trim();
        let val = line.slice(colonIdx + 1).trim();
        if ((val.startsWith('"') && val.endsWith('"')) || (val.startsWith("'") && val.endsWith("'"))) {
          val = val.slice(1, -1);
        }
        if (currentSection === 'styles') res.styles[key] = val;
        else if (currentSection === 'shortcuts') res.shortcuts[key] = val;
        else res[key] = val;
      }
    }
    return res;
  };

  const applyCSS = (css: string) => {
    let styleEl = document.getElementById('dynamic-skin-css') as HTMLStyleElement | null;
    if (!styleEl) {
      styleEl = document.createElement('style');
      styleEl.id = 'dynamic-skin-css';
      document.head.appendChild(styleEl);
    }
    styleEl.textContent = css;
  };

  const loadSkin = async (platform: string = 'twitter', ctx?: SkinContext) => {
    currentPlatform.value = platform;
    try {
      const app = (window as any)?.go?.main?.App;
      if (!app) return;

      let pkg: any = null;
      if (typeof app.GetSkinPackage === 'function') {
        pkg = await app.GetSkinPackage(platform);
      }

      if (pkg?.design_css) {
        applyCSS(pkg.design_css);
      } else if (typeof app.GetSkinCSS === 'function') {
        const css = await app.GetSkinCSS(platform);
        if (css) applyCSS(css);
      }

      if (pkg?.layout_yaml) {
        layoutConfig.value = parseSimpleYAML(pkg.layout_yaml);
      }

      if (pkg?.controller_js) {
        try {
          const blob = new Blob([pkg.controller_js], { type: 'application/javascript' });
          const url = URL.createObjectURL(blob);
          const mod = await import(/* @vite-ignore */ url);
          URL.revokeObjectURL(url);
          const ctrl = mod.default || mod;
          if (ctrl && typeof ctrl.init === 'function') {
            ctrl.init(ctx || { platform });
            controller.value = ctrl;
          }
        } catch (ctrlErr) {
          console.warn('[useSkin] Controller dynamic import fallback:', ctrlErr);
        }
      }

      isLoaded.value = true;
    } catch (e) {
      console.warn('[useSkin] Failed to load skin:', e);
    }
  };

  return {
    currentPlatform,
    layoutConfig,
    controller,
    isLoaded,
    loadSkin,
    applyCSS,
  };
}
