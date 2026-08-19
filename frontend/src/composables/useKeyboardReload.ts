import { onMounted, onUnmounted } from 'vue';
import { WindowReload } from '../../wailsjs/runtime/runtime';

/**
 * F5 や Ctrl+R / Cmd+R によるリロードを捕捉し、
 * Wails の WindowReload またはブラウザ標準のリロードを実行します。
 */
export function useKeyboardReload() {
  const handleKeyDown = (e: KeyboardEvent) => {
    const isCtrlR = (e.ctrlKey || e.metaKey) && (e.key === 'r' || e.key === 'R' || e.code === 'KeyR');
    if (isCtrlR || e.key === 'F5' || e.code === 'F5') {
      e.preventDefault();
      try {
        if (typeof WindowReload === 'function') WindowReload();
        else window.location.reload();
      } catch (_) {
        window.location.reload();
      }
    }
  };

  onMounted(() => {
    window.addEventListener('keydown', handleKeyDown, true);
  });

  onUnmounted(() => {
    window.removeEventListener('keydown', handleKeyDown, true);
  });
}
