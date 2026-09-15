// frontend/src/composables/admin/usePipelineMediaTasks.ts (100行以下 - SPEC-PRINCIPLE-001)
import { ref, computed } from 'vue';
import type { CandidateTask, MediaWithTasks } from '../../types/pipeline_media_tasks';
export type { CandidateTask, MediaWithTasks };

export function usePipelineMediaTasks() {
  const items = ref<MediaWithTasks[]>([]);
  const selectedMediaId = ref<string>('');
  const carouselIndex = ref<number>(0);
  const loading = ref<boolean>(false);

  const isWails = () => typeof (window as any).go !== 'undefined';

  const fetchItems = async (status = 'ESCALATED', limit = 50) => {
    loading.value = true;
    try {
      if (isWails()) {
        const app = (window as any).go?.app?.App;
        if (app?.GetMediaWithCandidateTasks) {
          const res = await app.GetMediaWithCandidateTasks(status, limit);
          if (res) items.value = res;
        }
      } else {
        const res = await fetch(`/api/admin/pipeline/media-tasks?status=${status}&limit=${limit}`);
        if (res.ok) {
          const data = await res.json();
          if (data) items.value = data;
        }
      }
      if (items.value.length > 0 && !selectedMediaId.value) {
        selectMedia(items.value[0].media_id);
      }
    } catch (_) {} finally {
      loading.value = false;
    }
  };

  const selectMedia = (id: string) => {
    selectedMediaId.value = id;
    carouselIndex.value = 0;
  };

  const activeMedia = computed(() => {
    return items.value.find((m) => m.media_id === selectedMediaId.value) || items.value[0] || null;
  });

  const activeTasks = computed(() => activeMedia.value?.candidate_tasks || []);

  const currentTask = computed(() => {
    const tasks = activeTasks.value;
    if (tasks.length === 0) return null;
    const idx = Math.min(Math.max(carouselIndex.value, 0), tasks.length - 1);
    return tasks[idx];
  });

  const nextTask = () => {
    if (carouselIndex.value < activeTasks.value.length - 1) carouselIndex.value++;
    else carouselIndex.value = 0;
  };

  const prevTask = () => {
    if (carouselIndex.value > 0) carouselIndex.value--;
    else carouselIndex.value = Math.max(0, activeTasks.value.length - 1);
  };

  return {
    items, selectedMediaId, carouselIndex, loading,
    activeMedia, activeTasks, currentTask,
    fetchItems, selectMedia, nextTask, prevTask,
  };
}
