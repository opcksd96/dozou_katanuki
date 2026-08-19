// frontend/src/composables/useMediaOverlay.ts (100行以下)
import { ref, computed, onMounted, onUnmounted } from 'vue';
import type { RenderMedia } from '../models/RenderTree';

export function useMediaOverlay() {
  const isOpen = ref(false);
  const mediaList = ref<RenderMedia[]>([]);
  const currentIndex = ref(0);

  const activeMedia = computed(() => {
    if (!isOpen.value || mediaList.value.length === 0) return null;
    return mediaList.value[currentIndex.value] || null;
  });

  const hasNext = computed(() => currentIndex.value < mediaList.value.length - 1);
  const hasPrev = computed(() => currentIndex.value > 0);

  const openMedia = (media: RenderMedia, list?: RenderMedia[]) => {
    if (list && list.length > 0) {
      mediaList.value = list;
      const idx = list.findIndex(m => m.id === media.id);
      currentIndex.value = idx >= 0 ? idx : 0;
    } else {
      mediaList.value = [media];
      currentIndex.value = 0;
    }
    isOpen.value = true;
  };

  const closeMedia = () => {
    isOpen.value = false;
    mediaList.value = [];
    currentIndex.value = 0;
  };

  const nextMedia = () => {
    if (hasNext.value) {
      currentIndex.value++;
    }
  };

  const prevMedia = () => {
    if (hasPrev.value) {
      currentIndex.value--;
    }
  };

  const handleKeyDown = (e: KeyboardEvent) => {
    if (!isOpen.value) return;
    if (e.key === 'Escape') {
      closeMedia();
    } else if (e.key === 'ArrowRight') {
      nextMedia();
    } else if (e.key === 'ArrowLeft') {
      prevMedia();
    }
  };

  onMounted(() => {
    window.addEventListener('keydown', handleKeyDown);
  });

  onUnmounted(() => {
    window.removeEventListener('keydown', handleKeyDown);
  });

  return {
    isOpen,
    activeMedia,
    currentIndex,
    mediaList,
    hasNext,
    hasPrev,
    openMedia,
    closeMedia,
    nextMedia,
    prevMedia,
  };
}
