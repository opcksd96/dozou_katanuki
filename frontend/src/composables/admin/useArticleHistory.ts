// frontend/src/composables/admin/useArticleHistory.ts (100行以下 - SPEC-PRINCIPLE-001)
import { ref, computed } from 'vue';

export type ArticleAction =
  | { type: 'TRASH'; ids: string[]; trashedBy: string; reason: string }
  | { type: 'RESTORE'; ids: string[] }
  | { type: 'RESET_TRANSLATIONS'; items: Array<{ id: string; ja: string; en: string; zh: string }> };

export function useArticleHistory() {
  const undoStack = ref<ArticleAction[]>([]), redoStack = ref<ArticleAction[]>([]);

  const pushAction = (action: ArticleAction) => {
    undoStack.value.push(action);
    if (undoStack.value.length > 50) undoStack.value.shift();
    redoStack.value = [];
  };

  const undo = async (api: any, refresh: () => Promise<void>) => {
    const action = undoStack.value.pop();
    if (!action) return null;
    try {
      if (action.type === 'TRASH') {
        await api.batchRestoreArticles(action.ids);
        redoStack.value.push(action);
      } else if (action.type === 'RESTORE') {
        await api.batchTrashArticles(action.ids, 'admin', 'Undo Restore');
        redoStack.value.push(action);
      } else if (action.type === 'RESET_TRANSLATIONS') {
        for (const it of action.items) { await api.saveTranslation(it.id, it.ja, it.en, it.zh); }
        redoStack.value.push(action);
      }
      await refresh();
      return action;
    } catch (e) {
      undoStack.value.push(action);
      throw e;
    }
  };

  const redo = async (api: any, refresh: () => Promise<void>) => {
    const action = redoStack.value.pop();
    if (!action) return null;
    try {
      if (action.type === 'TRASH') {
        await api.batchTrashArticles(action.ids, action.trashedBy, action.reason);
        undoStack.value.push(action);
      } else if (action.type === 'RESTORE') {
        await api.batchRestoreArticles(action.ids);
        undoStack.value.push(action);
      } else if (action.type === 'RESET_TRANSLATIONS') {
        await api.batchResetTranslations(action.items.map(i => i.id));
        undoStack.value.push(action);
      }
      await refresh();
      return action;
    } catch (e) {
      redoStack.value.push(action);
      throw e;
    }
  };

  return {
    pushAction, undo, redo,
    canUndo: computed(() => undoStack.value.length > 0),
    canRedo: computed(() => redoStack.value.length > 0),
    undoCount: computed(() => undoStack.value.length),
    redoCount: computed(() => redoStack.value.length),
  };
}
