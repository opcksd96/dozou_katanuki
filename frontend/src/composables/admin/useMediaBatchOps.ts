// frontend/src/composables/admin/useMediaBatchOps.ts (100行以下 - SPEC-PRINCIPLE-001)
import { ref, computed } from 'vue';

export function useMediaBatchOps() {
  const selectedIds = ref<Set<string>>(new Set());

  const toggleSelect = (id: string) => {
    if (!id) return;
    const next = new Set(selectedIds.value);
    if (next.has(id)) next.delete(id);
    else next.add(id);
    selectedIds.value = next;
  };

  const selectAll = (items: Array<{ media_id?: string; id?: string }>) => {
    const next = new Set<string>();
    items.forEach(m => {
      const id = m?.media_id || m?.id;
      if (id) next.add(id);
    });
    selectedIds.value = next;
  };

  const clearSelection = () => {
    selectedIds.value = new Set();
  };

  const isSelected = (id: string) => selectedIds.value.has(id);

  return {
    selectedIds,
    selectedCount: computed(() => selectedIds.value.size),
    toggleSelect,
    selectAll,
    clearSelection,
    isSelected,
  };
}
