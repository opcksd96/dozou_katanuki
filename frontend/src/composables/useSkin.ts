// frontend/src/composables/useSkin.ts (100行以下 - SPEC-PLUGIN-001)
import { ref } from 'vue';
import type { SkinController, SkinContext } from '../models/SkinController';

export interface SkinLayoutConfig {
  name?: string; version?: string; platform?: string; display_name?: string;
  styles?: Record<string, string>; shortcuts?: Record<string, any>;
}

export function useSkin() {
  const currentPlatform = ref('twitter'), layoutConfig = ref<SkinLayoutConfig | null>(null);
  const controller = ref<SkinController | null>(null), isLoaded = ref(false);

  const parseSimpleYAML = (yamlStr: string): SkinLayoutConfig => {
    const res: any = { styles: {}, shortcuts: {} };
    let section = '';
    for (const raw of yamlStr.split('\n')) {
      const line = raw.trim();
      if (!line || line.startsWith('#')) continue;
      if (line.endsWith(':') && !line.includes(' ')) { section = line.slice(0, -1).trim(); continue; }
      const idx = line.indexOf(':');
      if (idx > 0) {
        const k = line.slice(0, idx).trim();
        let v = line.slice(idx + 1).trim();
        if ((v.startsWith('"') && v.endsWith('"')) || (v.startsWith("'") && v.endsWith("'"))) v = v.slice(1, -1);
        if (section === 'styles') res.styles[k] = v;
        else if (section === 'shortcuts') res.shortcuts[k] = v;
        else res[k] = v;
      }
    }
    return res;
  };

  const applyCSS = (css: string) => {
    let el = document.getElementById('dynamic-skin-css') as HTMLStyleElement | null;
    if (!el) {
      el = document.createElement('style');
      el.id = 'dynamic-skin-css';
      document.head.appendChild(el);
    }
    el.textContent = css;
  };

  const loadSkin = async (platform: string = 'twitter', ctx?: SkinContext) => {
    currentPlatform.value = platform;
    try {
      const app = (window as any)?.go?.main?.App;
      if (!app) return;
      let pkg: any = null;
      if (typeof app.GetSkinPackage === 'function') pkg = await app.GetSkinPackage(platform);
      if (pkg?.design_css) applyCSS(pkg.design_css);
      else if (typeof app.GetSkinCSS === 'function') {
        const css = await app.GetSkinCSS(platform);
        if (css) applyCSS(css);
      }
      if (pkg?.layout_yaml) layoutConfig.value = parseSimpleYAML(pkg.layout_yaml);
      if (pkg?.controller_js) {
        try {
          const blob = new Blob([pkg.controller_js], { type: 'application/javascript' });
          const url = URL.createObjectURL(blob);
          const mod = await import(/* @vite-ignore */ url);
          URL.revokeObjectURL(url);
          const ctrl = mod.default || mod;
          if (ctrl?.init) { ctrl.init(ctx || { platform }); controller.value = ctrl; }
        } catch (_) {}
      }
      isLoaded.value = true;
    } catch (_) {}
  };

  return { currentPlatform, layoutConfig, controller, isLoaded, loadSkin, applyCSS };
}
