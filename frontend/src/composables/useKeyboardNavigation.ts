// frontend/src/composables/useKeyboardNavigation.ts (100行以下 - SPEC-FRONTEND-001)
import { ref, onMounted, onUnmounted } from 'vue';
import type { RenderTree } from '../models/RenderTree';

export interface KeyboardNavOptions {
  getItems: () => RenderTree[];
  onSelectArticle: (id: string) => void;
  onToggleLike: (id: string) => void;
  onOpenMedia: (article: RenderTree) => void;
  onBack: () => void;
  isDetailView: () => boolean;
  isOverlayOpen: () => boolean;
  isAdminOpen: () => boolean;
  openAdmin: () => void;
}

export function useKeyboardNavigation(options: KeyboardNavOptions) {
  const focusedIndex = ref<number>(-1), isHelpOpen = ref(false);

  const scrollToFocused = () => {
    setTimeout(() => {
      const el = document.querySelector('.twitter-card.is-focused, .twitter-detail-item.is-focused');
      if (el) el.scrollIntoView({ behavior: 'smooth', block: 'nearest' });
    }, 10);
  };

  const handleKeyDown = (e: KeyboardEvent) => {
    const t = e.target as HTMLElement;
    if (t && (t.tagName === 'INPUT' || t.tagName === 'TEXTAREA' || t.isContentEditable)) return;

    if ((e.ctrlKey || e.metaKey) && e.key === ',') { e.preventDefault(); options.openAdmin(); return; }
    if (e.key === '?' || (e.shiftKey && e.key === '/')) { e.preventDefault(); isHelpOpen.value = !isHelpOpen.value; return; }
    if (isHelpOpen.value) { if (e.key === 'Escape') isHelpOpen.value = false; return; }
    if (options.isAdminOpen() || options.isOverlayOpen()) return;

    const items = options.getItems();
    if (!items || items.length === 0) return;

    if (e.key === 'j' || e.key === 'J') {
      e.preventDefault();
      focusedIndex.value = (focusedIndex.value < items.length - 1) ? focusedIndex.value + 1 : 0;
      scrollToFocused();
    } else if (e.key === 'k' || e.key === 'K') {
      e.preventDefault();
      focusedIndex.value = (focusedIndex.value > 0) ? focusedIndex.value - 1 : items.length - 1;
      scrollToFocused();
    } else if (e.key === 'Enter' || e.key === 'o' || e.key === 'O') {
      if (focusedIndex.value >= 0 && focusedIndex.value < items.length) {
        e.preventDefault(); options.onSelectArticle(items[focusedIndex.value].id);
      }
    } else if (e.key === 'l' || e.key === 'L') {
      if (focusedIndex.value >= 0 && focusedIndex.value < items.length) {
        e.preventDefault(); options.onToggleLike(items[focusedIndex.value].id);
      }
    } else if (e.key === 'm' || e.key === 'M' || e.key === 'x' || e.key === 'X') {
      if (focusedIndex.value >= 0 && focusedIndex.value < items.length) {
        const it = items[focusedIndex.value];
        if (it.media && it.media.length > 0) { e.preventDefault(); options.onOpenMedia(it); }
      }
    } else if (e.key === 'Escape') {
      if (options.isDetailView()) { e.preventDefault(); options.onBack(); }
    }
  };

  onMounted(() => window.addEventListener('keydown', handleKeyDown));
  onUnmounted(() => window.removeEventListener('keydown', handleKeyDown));

  return { focusedIndex, isHelpOpen };
}
