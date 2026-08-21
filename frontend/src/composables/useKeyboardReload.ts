// frontend/src/composables/useKeyboardReload.ts (100行以下)
import { onMounted, onUnmounted } from 'vue';
import { WindowReload } from '../../wailsjs/runtime/runtime';

/**
 * 画面（WebView / DOM）の完全リロードを実行する統一関数
 */
export function reloadWindow() {
  try {
    if (typeof WindowReload === 'function') {
      WindowReload();
      return;
    }
  } catch (_) {}
  window.location.reload();
}

/**
 * F5 や Ctrl+R / Cmd+R によるリロードショートカットをハンドリングし、
 * 共通の reloadWindow を実行します。
 */
export function useKeyboardReload() {
  const handleKeyDown = (e: KeyboardEvent) => {
    const isCtrlR = (e.ctrlKey || e.metaKey) && (e.key === 'r' || e.key === 'R' || e.code === 'KeyR');
    const isF5 = e.key === 'F5' || e.code === 'F5';
    if (isCtrlR || isF5) {
      e.preventDefault();
      reloadWindow();
    }
  };

  onMounted(() => {
    window.addEventListener('keydown', handleKeyDown, true);
  });

  onUnmounted(() => {
    window.removeEventListener('keydown', handleKeyDown, true);
  });
}


