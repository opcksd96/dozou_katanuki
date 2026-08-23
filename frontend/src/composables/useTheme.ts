// frontend/src/composables/useTheme.ts
import { ref } from 'vue';
import {
  WindowSetDarkTheme,
  WindowSetLightTheme,
  WindowSetSystemDefaultTheme,
  WindowSetBackgroundColour,
} from '../../wailsjs/runtime/runtime';

export type ThemeMode = 'system' | 'dark' | 'light';

const currentTheme = ref<ThemeMode>('dark');
const fontScale = ref<number>(1.0);
let mediaQueryListener: ((e: MediaQueryListEvent) => void) | null = null;

export function useTheme() {
  const applyDOMTheme = (isDark: boolean) => {
    if (isDark) {
      document.documentElement.classList.add('dark');
      document.documentElement.classList.remove('light');
      document.documentElement.setAttribute('data-theme', 'dark');
      try { WindowSetBackgroundColour(2, 6, 23, 220); } catch (_) {}
    } else {
      document.documentElement.classList.add('light');
      document.documentElement.classList.remove('dark');
      document.documentElement.setAttribute('data-theme', 'light');
      try { WindowSetBackgroundColour(248, 250, 252, 220); } catch (_) {}
    }
  };

  const applyWailsTheme = (theme: ThemeMode) => {
    try {
      if (theme === 'dark') {
        WindowSetDarkTheme();
      } else if (theme === 'light') {
        WindowSetLightTheme();
      } else {
        WindowSetSystemDefaultTheme();
      }
    } catch (_) {}
  };

  const applyFontScale = (scale: number) => {
    const validScale = Math.max(0.75, Math.min(1.5, scale || 1.0));
    document.documentElement.style.fontSize = `${validScale * 100}%`;
    document.documentElement.style.setProperty('--app-font-scale', validScale.toString());
  };

  const setFontScale = (scale: number) => {
    fontScale.value = scale;
    localStorage.setItem('dozou_font_scale', scale.toString());
    applyFontScale(scale);
  };

  const setTheme = (theme: ThemeMode) => {
    currentTheme.value = theme;
    localStorage.setItem('dozou_theme', theme);
    applyWailsTheme(theme);

    const mql = window.matchMedia('(prefers-color-scheme: dark)');

    if (mediaQueryListener) {
      mql.removeEventListener('change', mediaQueryListener);
      mediaQueryListener = null;
    }

    if (theme === 'system') {
      applyDOMTheme(mql.matches);
      mediaQueryListener = (e: MediaQueryListEvent) => {
        if (currentTheme.value === 'system') {
          applyDOMTheme(e.matches);
        }
      };
      mql.addEventListener('change', mediaQueryListener);
    } else {
      applyDOMTheme(theme === 'dark');
    }
  };

  const initTheme = (initialTheme?: string, initialScale?: number) => {
    const savedTheme = initialTheme || (localStorage.getItem('dozou_theme') as ThemeMode) || 'dark';
    setTheme(savedTheme as ThemeMode);

    const savedScaleStr = localStorage.getItem('dozou_font_scale');
    const savedScale = initialScale || (savedScaleStr ? parseFloat(savedScaleStr) : 1.0);
    fontScale.value = savedScale;
    applyFontScale(savedScale);
  };

  return {
    currentTheme,
    fontScale,
    setTheme,
    setFontScale,
    initTheme,
  };
}
