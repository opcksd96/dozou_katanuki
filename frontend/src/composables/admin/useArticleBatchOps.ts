// frontend/src/composables/admin/useArticleBatchOps.ts (100行以下 - SPEC-PRINCIPLE-001)
import { ref, computed } from 'vue';

export function useArticleBatchOps() {
  const selectedIds = ref<Set<string>>(new Set());

  const toggleSelect = (id: string) => {
    const next = new Set(selectedIds.value);
    if (next.has(id)) next.delete(id);
    else next.add(id);
    selectedIds.value = next;
  };

  const selectAll = (articles: Array<{ id?: string }>) => {
    const next = new Set<string>();
    articles.forEach(a => { if (a?.id) next.add(a.id); });
    selectedIds.value = next;
  };

  const clearSelection = () => {
    selectedIds.value = new Set();
  };

  const isSelected = (id: string) => selectedIds.value.has(id);

  return {
    selectedIds,
    selectedCount: computed(() => selectedIds.value.size),
    toggleSelect, selectAll, clearSelection, isSelected,
  };
}
