// frontend/src/composables/useToast.ts (100行以下 - SPEC-PRINCIPLE-001)
import { ref, onMounted, onUnmounted } from 'vue';
import { EventsOn } from '../../wailsjs/runtime/runtime';

export interface ToastItem {
  id: string;
  type: 'success' | 'info' | 'warning' | 'error';
  message: string;
}

const toasts = ref<ToastItem[]>([]);
let lastToastMsg = '', lastToastAt = 0;

export function useToast() {
  const addToast = (message: string, type: 'success' | 'info' | 'warning' | 'error' = 'info', duration = 4000) => {
    const now = Date.now();
    if (message === lastToastMsg && now - lastToastAt < 3000) return; // 3秒以内の同一メッセージ重複ラッシュを抑制
    lastToastMsg = message; lastToastAt = now;

    const id = `${now}-${Math.random().toString(36).substring(2, 9)}`;
    toasts.value.push({ id, type, message });
    if (duration > 0) setTimeout(() => removeToast(id), duration);
  };

  const removeToast = (id: string) => { toasts.value = toasts.value.filter((t) => t.id !== id); };
  const unoffs: (() => void)[] = [];

  onMounted(() => {
    try {
      if ((window as any)?.runtime?.EventsOnMultiple) {
        unoffs.push(
          EventsOn('toast:notify', (data: { type?: string; message?: string } | string) => {
            if (typeof data === 'string') addToast(data, 'info');
            else if (data && data.message) addToast(data.message, (data.type as any) || 'info');
          })
        );
      }
    } catch (_) {}
  });

  onUnmounted(() => { unoffs.forEach((fn) => { try { fn(); } catch (_) {} }); });
  return { toasts, addToast, removeToast };
}
