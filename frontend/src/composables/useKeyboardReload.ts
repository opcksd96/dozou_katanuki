// frontend/src/composables/useKeyboardReload.ts (100行以下)
import { onMounted, onUnmounted } from 'vue';
import { EventsOn, EventsOff } from '../../wailsjs/runtime/runtime';

/**
 * F5 や Ctrl+R / Cmd+R によるリロードを捕捉し、
 * WebView をブラックアウトさせずに、安全にデータリフレッシュ (reloadCallback) を実行します。
 */
export function useKeyboardReload(onReload?: () => void | Promise<void>) {
  const triggerReload = () => {
    if (onReload) {
      onReload();
    }
  };

  const handleKeyDown = (e: KeyboardEvent) => {
    const isCtrlR = (e.ctrlKey || e.metaKey) && (e.key === 'r' || e.key === 'R' || e.code === 'KeyR');
    if (isCtrlR || e.key === 'F5' || e.code === 'F5') {
      e.preventDefault();
      triggerReload();
    }
  };

  onMounted(() => {
    window.addEventListener('keydown', handleKeyDown, true);
    try {
      EventsOn('app:refresh', triggerReload);
    } catch (_) {}
  });

  onUnmounted(() => {
    window.removeEventListener('keydown', handleKeyDown, true);
    try {
      EventsOff('app:refresh');
    } catch (_) {}
  });
}
