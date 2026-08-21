// frontend/src/composables/useToast.ts (100行以下)
import { ref, onMounted, onUnmounted } from 'vue';
import { EventsOn } from '../../wailsjs/runtime/runtime';

export interface ToastItem {
  id: string;
  type: 'success' | 'info' | 'warning' | 'error';
  message: string;
}

const toasts = ref<ToastItem[]>([]);

export function useToast() {
  const addToast = (message: string, type: 'success' | 'info' | 'warning' | 'error' = 'info', duration = 4000) => {
    const id = `${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;
    toasts.value.push({ id, type, message });
    if (duration > 0) {
      setTimeout(() => {
        removeToast(id);
      }, duration);
    }
  };

  const removeToast = (id: string) => {
    toasts.value = toasts.value.filter((t) => t.id !== id);
  };

  const unoffs: (() => void)[] = [];

  onMounted(() => {
    try {
      if ((window as any)?.runtime?.EventsOnMultiple) {
        unoffs.push(
          EventsOn('toast:notify', (data: { type?: string; message?: string } | string) => {
            if (typeof data === 'string') {
              addToast(data, 'info');
            } else if (data && data.message) {
              const t = (data.type as any) || 'info';
              addToast(data.message, t);
            }
          })
        );
      }
    } catch (_) {}
  });

  onUnmounted(() => {
    unoffs.forEach((fn) => {
      try { fn(); } catch (_) {}
    });
  });

  return {
    toasts,
    addToast,
    removeToast,
  };
}
