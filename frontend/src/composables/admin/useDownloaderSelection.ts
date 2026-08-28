// frontend/src/composables/admin/useDownloaderSelection.ts (100行以下)
import { ref, computed } from 'vue';

export function useDownloaderSelection() {
  const selectedGids = ref<Set<string>>(new Set());

  const selectedCount = computed(() => selectedGids.value.size);
  const isSelected = (gid: string) => selectedGids.value.has(gid);

  const toggleSelect = (gid: string) => {
    const next = new Set(selectedGids.value);
    if (next.has(gid)) next.delete(gid);
    else next.add(gid);
    selectedGids.value = next;
  };

  const selectAll = (tasks: Array<{ gid: string }>) => {
    const next = new Set(selectedGids.value);
    tasks.forEach((t) => { if (t.gid) next.add(t.gid); });
    selectedGids.value = next;
  };

  const selectOnlyErrors = (tasks: Array<{ gid: string; status: string }>) => {
    const next = new Set<string>();
    tasks.filter((t) => t.status === 'error' || t.status === 'paused').forEach((t) => { if (t.gid) next.add(t.gid); });
    selectedGids.value = next;
  };

  const clearSelection = () => { selectedGids.value = new Set(); };

  return { selectedGids, selectedCount, isSelected, toggleSelect, selectAll, selectOnlyErrors, clearSelection };
}
